package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vsuaiqq/cicd/git-gateway-service/internal/config"
	httpHandler "github.com/vsuaiqq/cicd/git-gateway-service/internal/http/handler"
	"github.com/vsuaiqq/cicd/git-gateway-service/internal/kafka"
	"github.com/vsuaiqq/cicd/git-gateway-service/internal/projects"
	sharedChi "github.com/vsuaiqq/cicd/shared/chi"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

func main() {
	flag.Parse()
	logger.Init(logger.Config{ServiceName: "git-gateway-service", Level: os.Getenv("LOG_LEVEL"), Pretty: os.Getenv("LOG_PRETTY") == "1"})
	log := logger.L()

	cfg, err := sharedConfig.LoadService(context.Background(), sharedConfig.ServiceLoader[config.Config]{
		ConfigPath: *configPath,
		ApplyEnv:   config.ApplyEnv,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("config failed")
	}

	_, shutdownTracer, err := telemetry.TracerProvider(context.Background(), telemetry.Config{
		ServiceName:  "git-gateway-service",
		OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		SampleRatio:  1,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry init failed")
	}
	defer shutdownTracer()

	timeout := cfg.Projects.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	projectsClient, err := projects.NewClient(cfg.Projects.GRPCAddress, timeout)
	if err != nil {
		log.Fatal().Err(err).Msg("projects grpc client failed")
	}
	defer projectsClient.Close()
	log.Info().Msg("connected to projects-service")

	producer, err := sharedKafka.NewProducer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("kafka producer failed")
	}
	defer producer.Close()

	publisher := kafka.NewGitEventPublisher(producer, cfg.Kafka.OutputTopic)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      setupRouter(projectsClient, publisher),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Int("port", cfg.Server.Port).Msg("starting")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
	}
	log.Info().Msg("stopped")
}

func setupRouter(projectsClient *projects.Client, publisher *kafka.GitEventPublisher) *chi.Mux {
	router := chi.NewRouter()
	router.Use(telemetry.Middleware("git-gateway-service"))
	for _, mw := range sharedChi.DefaultStack(30 * time.Second) {
		router.Use(mw)
	}

	router.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	router.Route("/webhook", func(r chi.Router) {
		r.Post("/github/{projectID}", httpHandler.WebhookHandler(
			projectsClient, publisher, httpHandler.NewGitHubProcessor(),
		))

	})

	return router
}
