package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vsuaiqq/cicd/analytics-service/internal/config"
	"github.com/vsuaiqq/cicd/analytics-service/internal/db"
	"github.com/vsuaiqq/cicd/analytics-service/internal/grpchandler"
	analyticsKafka "github.com/vsuaiqq/cicd/analytics-service/internal/kafka"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/analytics"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

func main() {
	flag.Parse()
	logger.Init(logger.Config{ServiceName: "analytics-service", Level: os.Getenv("LOG_LEVEL"), Pretty: os.Getenv("LOG_PRETTY") == "1"})
	log := logger.L()

	ctx := context.Background()
	cfg, err := sharedConfig.LoadService(ctx, sharedConfig.ServiceLoader[config.Config]{
		ServiceName:   "analytics-service",
		ConfigPath:    *configPath,
		ApplyEnv:      config.ApplyEnv,
		ApplyDefaults: config.ApplyDefaults,
		Validate:      config.Validate,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("config failed")
	}

	_, shutdownTracer, err := telemetry.TracerProvider(context.Background(), telemetry.Config{
		ServiceName:  "analytics-service",
		OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		SampleRatio:  1,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry init failed")
	}
	defer shutdownTracer()

	conn, err := openClickHouse(cfg.ClickHouse)
	if err != nil {
		log.Fatal().Err(err).Msg("clickhouse connection failed")
	}
	defer conn.Close()
	log.Info().Str("addr", cfg.ClickHouse.Addr).Str("database", cfg.ClickHouse.Database).Msg("clickhouse connected")

	repo := db.NewRepository(conn)
	if err := repo.Migrate(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}
	log.Info().Msg("clickhouse schema ready")

	runConsumerClient, err := sharedKafka.NewConsumer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID + "-runs",
		GroupID:  cfg.Kafka.RunsGroupID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("run consumer failed")
	}
	runConsumer := analyticsKafka.NewRunEventConsumer(runConsumerClient, repo)

	jobConsumerClient, err := sharedKafka.NewConsumer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID + "-jobs",
		GroupID:  cfg.Kafka.JobsGroupID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("job consumer failed")
	}
	jobConsumer := analyticsKafka.NewJobResultConsumer(jobConsumerClient, repo)

	grpcServer := grpc.NewServer()
	pb.RegisterAnalyticsServiceServer(grpcServer, grpchandler.NewAnalyticsServer(repo))
	reflection.Register(grpcServer)

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC listen failed")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := conn.Ping(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("clickhouse unreachable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runConsumer.Start(ctx, []string{cfg.Kafka.RunEventsTopic}); err != nil && ctx.Err() == nil {
			log.Error().Err(err).Msg("run consumer error")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := jobConsumer.Start(ctx, []string{cfg.Kafka.JobResultsTopic}); err != nil && ctx.Err() == nil {
			log.Error().Err(err).Msg("job consumer error")
		}
	}()

	go func() {
		log.Info().Int("port", cfg.Server.GRPCPort).Msg("gRPC listening")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	go func() {
		log.Info().Int("port", cfg.Server.HTTPPort).Msg("HTTP healthz listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down")
	cancel()
	wg.Wait()
	_ = runConsumer.Close()
	_ = jobConsumer.Close()
	grpcServer.GracefulStop()
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutCancel()
	if err := httpServer.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("HTTP shutdown error")
	}
	log.Info().Msg("stopped")
}

func openClickHouse(cfg config.ClickHouseConfig) (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.Addr},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		DialTimeout:     cfg.DialTimeout,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Settings: clickhouse.Settings{
			"max_execution_time":             60,
			"max_bytes_before_external_sort": 1 << 28,
		},
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()
	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("clickhouse ping: %w", err)
	}
	return conn, nil
}
