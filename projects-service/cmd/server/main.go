package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vsuaiqq/cicd/projects-service/internal/config"
	projectdb "github.com/vsuaiqq/cicd/projects-service/internal/db"
	"github.com/vsuaiqq/cicd/projects-service/internal/grpchandler"
	"github.com/vsuaiqq/cicd/projects-service/internal/repository"
	"github.com/vsuaiqq/cicd/projects-service/internal/service"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	"github.com/vsuaiqq/cicd/shared/logger"
	sharedPostgres "github.com/vsuaiqq/cicd/shared/postgres"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/projects"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

func main() {
	flag.Parse()
	logger.Init(logger.Config{ServiceName: "projects-service", Level: os.Getenv("LOG_LEVEL"), Pretty: os.Getenv("LOG_PRETTY") == "1"})
	log := logger.L()

	ctx := context.Background()
	cfg, err := sharedConfig.LoadService(ctx, sharedConfig.ServiceLoader[config.Config]{
		ServiceName: "projects-service",
		ConfigPath:  *configPath,
		ApplyEnv:    config.ApplyEnv,
		Validate:    config.Validate,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("config failed")
	}

	_, shutdownTracer, err := telemetry.TracerProvider(context.Background(), telemetry.Config{
		ServiceName:  "projects-service",
		OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		SampleRatio:  1,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry init failed")
	}
	defer shutdownTracer()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Postgres.ConnectTimeout)
	defer cancel()

	db, err := sharedPostgres.New(ctx, sharedPostgres.Config{
		DSN:             cfg.Postgres.DSN,
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Postgres.ConnMaxIdleTime,
		ConnectTimeout:  cfg.Postgres.ConnectTimeout,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer db.Close()
	if err := projectdb.RunMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}
	log.Info().Msg("Database connected successfully")

	encrypter, err := service.NewEncrypter(cfg.Projects.EncryptionKey)
	if err != nil {
		log.Fatal().Err(err).Msg("encrypter init failed")
	}

	repo := repository.NewProjectRepository(db)
	projectSvc := service.NewProjectService(repo, encrypter, cfg.Projects.WebhookBaseURL)
	projectsServer := grpchandler.NewProjectsServer(projectSvc)

	grpcServer := grpc.NewServer()
	pb.RegisterProjectsServiceServer(grpcServer, projectsServer)
	if os.Getenv("ENABLE_GRPC_REFLECTION") == "1" {
		reflection.Register(grpcServer)
		log.Info().Msg("gRPC reflection enabled (ENABLE_GRPC_REFLECTION=1)")
	}

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on gRPC port")
	}

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler:      httpMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Int("port", cfg.Server.GRPCPort).Msg("gRPC server starting")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	go func() {
		log.Info().Int("port", cfg.Server.HTTPPort).Msg("HTTP server (healthz) starting")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down servers...")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := httpServer.Shutdown(ctxShutdown); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}
	select {
	case <-stopped:
		log.Info().Msg("gRPC server stopped")
	case <-time.After(5 * time.Second):
		grpcServer.Stop()
	}
	log.Info().Msg("Server stopped")
}
