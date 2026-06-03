package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orchApi "github.com/vsuaiqq/cicd/orchestrator-service/internal/api"
	"github.com/vsuaiqq/cicd/orchestrator-service/internal/config"
	orchDb "github.com/vsuaiqq/cicd/orchestrator-service/internal/db"
	orchInt "github.com/vsuaiqq/cicd/orchestrator-service/internal/internalauth"
	orchKafka "github.com/vsuaiqq/cicd/orchestrator-service/internal/kafka"
	"github.com/vsuaiqq/cicd/orchestrator-service/internal/orchestrator"
	orchWs "github.com/vsuaiqq/cicd/orchestrator-service/internal/ws"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/projects"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

const serviceName = "orchestrator-service"

func main() {
	flag.Parse()

	logger.Init(logger.Config{
		ServiceName: serviceName,
		Level:       os.Getenv("LOG_LEVEL"),
		Pretty:      os.Getenv("LOG_PRETTY") == "1",
	})
	log := logger.L()

	ctx := context.Background()
	cfg, err := sharedConfig.LoadService(ctx, sharedConfig.ServiceLoader[config.Config]{
		ServiceName: "orchestrator-service",
		ConfigPath:  *configPath,
		ApplyEnv:    config.ApplyEnv,
		Validate:    config.Validate,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("config failed")
	}

	tp, shutdown, err := telemetry.TracerProvider(ctx, telemetry.Config{
		ServiceName:  serviceName,
		OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		SampleRatio:  1,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry init failed")
	}
	defer shutdown()
	_ = tp

	dsn := cfg.Postgres.DSN
	db, err := openDB(dsn, cfg.Postgres)
	if err != nil {
		log.Fatal().Err(err).Msg("database connection failed")
	}
	defer db.Close()
	log.Info().Msg("database connected")

	if err := orchDb.RunMigrations(db); err != nil {
		log.Fatal().Err(err).Msg("migrations failed")
	}

	repo := orchDb.NewRepository(db)

	producer, err := sharedKafka.NewProducer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("kafka producer failed")
	}
	defer producer.Close()

	cancelTopic := cfg.Kafka.CancelJobsTopic
	if cancelTopic == "" {
		cancelTopic = "pipeline.cancel-jobs"
	}
	runEventsTopic := cfg.Kafka.RunEventsTopic
	if runEventsTopic == "" {
		runEventsTopic = "pipeline.run-events"
	}
	jobPub := orchKafka.NewJobPublisher(producer, cfg.Kafka.JobsTopic)
	cancelPub := orchKafka.NewCancelPublisher(producer, cancelTopic)
	runEventPub := orchKafka.NewRunEventPublisher(producer, runEventsTopic)

	projTimeout := cfg.Projects.Timeout
	if projTimeout == 0 {
		projTimeout = 10 * time.Second
	}
	projConn, err := grpc.NewClient(
		cfg.Projects.GRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("projects grpc client failed")
	}
	defer projConn.Close()
	projectsClient := pb.NewProjectsServiceClient(projConn)
	log.Info().Msg("connected to projects-service")

	analyticsTimeout := cfg.Analytics.Timeout
	if analyticsTimeout == 0 {
		analyticsTimeout = 15 * time.Second
	}
	perfGate, err := orchestrator.NewAnalyticsGateClient(cfg.Analytics.GRPCAddress, analyticsTimeout)
	if err != nil {
		log.Fatal().Err(err).Msg("analytics gate client failed")
	}
	log.Info().Msg("connected to analytics-service (performance gate)")

	hub := orchWs.NewHub()

	orch := orchestrator.New(repo, jobPub, cancelPub, runEventPub, projectsClient, perfGate, hub)

	gitConsumerClient, err := sharedKafka.NewConsumer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID + "-git",
		GroupID:  cfg.Kafka.GroupID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("git kafka consumer failed")
	}
	gitConsumer := orchKafka.NewGitEventConsumer(gitConsumerClient, orch)

	resultConsumerClient, err := sharedKafka.NewConsumer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID + "-results",
		GroupID:  cfg.Kafka.ResultsGroupID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("results kafka consumer failed")
	}
	resultConsumer := orchKafka.NewJobResultConsumer(resultConsumerClient, orch)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      setupRouter(repo, hub, orch, cfg.InternalAPIKey),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hub.Start(ctx)

	var consumerWg sync.WaitGroup

	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		if err := gitConsumer.Start(ctx, []string{cfg.Kafka.InputTopic}); err != nil && ctx.Err() == nil {
			log.Error().Err(err).Msg("git consumer error")
		}
	}()

	consumerWg.Add(1)
	go func() {
		defer consumerWg.Done()
		if err := resultConsumer.Start(ctx, []string{cfg.Kafka.JobResultsTopic}); err != nil && ctx.Err() == nil {
			log.Error().Err(err).Msg("result consumer error")
		}
	}()

	go func() {
		log.Info().Int("port", cfg.Server.Port).Msg("HTTP server starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")

	cancel()
	consumerWg.Wait()
	_ = gitConsumer.Close()
	_ = resultConsumer.Close()

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutCancel()
	if err := server.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("HTTP shutdown error")
	}

	log.Info().Msg("stopped")
}

func setupRouter(repo *orchDb.Repository, hub *orchWs.Hub, orch orchApi.RunController, internalAPIKey string) *chi.Mux {
	router := chi.NewRouter()
	router.Use(telemetry.Middleware(serviceName))
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	keyMW := orchInt.APIKeyMiddleware(internalAPIKey)

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	router.Group(func(r chi.Router) {
		r.Use(keyMW)
		wsHandler := orchWs.NewHandler(hub)
		r.Get("/ws/runs/{runID}", wsHandler.ServeWS)
		r.Get("/ws/events", wsHandler.ServeGlobalWS)
	})

	router.Group(func(r chi.Router) {
		r.Use(keyMW)
		r.Use(middleware.Timeout(30 * time.Second))
		orchApi.NewHandler(repo, orch).RegisterRoutes(r)
	})

	return router
}

func openDB(dsn string, cfg config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx/v5", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
