package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/vsuaiqq/cicd/auth-service/internal/config"
	authdb "github.com/vsuaiqq/cicd/auth-service/internal/db"
	"github.com/vsuaiqq/cicd/auth-service/internal/handler"
	"github.com/vsuaiqq/cicd/auth-service/internal/repository"
	"github.com/vsuaiqq/cicd/auth-service/internal/service"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	"github.com/vsuaiqq/cicd/shared/logger"
	sharedPostgres "github.com/vsuaiqq/cicd/shared/postgres"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/auth"
	sharedRedis "github.com/vsuaiqq/cicd/shared/redis"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

func main() {
	flag.Parse()
	logger.Init(logger.Config{ServiceName: "auth-service", Level: os.Getenv("LOG_LEVEL"), Pretty: os.Getenv("LOG_PRETTY") == "1"})
	log := logger.L()

	ctx := context.Background()
	cfg, err := sharedConfig.LoadService(ctx, sharedConfig.ServiceLoader[config.Config]{
		ServiceName: "auth-service",
		ConfigPath:  *configPath,
		ApplyEnv:    config.ApplyEnv,
		Validate:    config.Validate,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("config failed")
	}

	_, shutdownTracer, err := telemetry.TracerProvider(context.Background(), telemetry.Config{
		ServiceName:  "auth-service",
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
	if err := authdb.RunMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}
	log.Info().Msg("Database connected successfully")

	redisClient, err := sharedRedis.New(ctx, sharedRedis.Config{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
		PoolTimeout:  cfg.Redis.PoolTimeout,
		MaxRetries:   cfg.Redis.MaxRetries,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("redis connection failed")
	}
	defer redisClient.Close()
	log.Info().Msg("Redis connected successfully")

	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(redisClient)

	jwtService := service.NewJWTService(cfg.JWT)
	authService := service.NewAuthService(userRepo, refreshTokenRepo, jwtService)

	authHandler := handler.NewAuthHandler(authService)

	grpcServer := setupGRPCServer(cfg, authHandler)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on gRPC port")
	}

	httpServer := setupHTTPServer(cfg, db, redisClient)

	go func() {
		log.Info().Int("port", cfg.Server.GRPCPort).Msg("gRPC server starting")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	go func() {
		log.Info().Int("port", cfg.Server.HTTPPort).Msg("HTTP server starting")
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

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
	}
	select {
	case <-stopped:
		log.Info().Msg("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		log.Info().Msg("gRPC server shutdown timeout, forcing stop")
		grpcServer.Stop()
	}
	log.Info().Msg("Servers stopped")
}

func setupGRPCServer(cfg *config.Config, authHandler *handler.AuthHandler) *grpc.Server {

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	reflection.Register(grpcServer)

	return grpcServer
}

func setupHTTPServer(cfg *config.Config, db *sql.DB, redisClient *redis.Client) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready: postgres"))
			return
		}

		if err := redisClient.Ping(ctx).Err(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready: redis"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
