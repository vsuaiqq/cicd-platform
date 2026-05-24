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

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/vsuaiqq/cicd/ai-service/internal/config"
	"github.com/vsuaiqq/cicd/ai-service/internal/grpchandler"
	"github.com/vsuaiqq/cicd/ai-service/internal/llm"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	"github.com/vsuaiqq/cicd/shared/logger"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/ai"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

func main() {
	flag.Parse()
	logger.Init(logger.Config{ServiceName: "ai-service", Level: os.Getenv("LOG_LEVEL"), Pretty: os.Getenv("LOG_PRETTY") == "1"})
	log := logger.L()

	ctx := context.Background()
	cfg, err := sharedConfig.LoadService(ctx, sharedConfig.ServiceLoader[config.Config]{
		ServiceName:   "ai-service",
		ConfigPath:    *configPath,
		ApplyEnv:      config.ApplyEnv,
		ApplyDefaults: config.ApplyDefaults,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("config failed")
	}

	_, shutdownTracer, err := telemetry.TracerProvider(context.Background(), telemetry.Config{
		ServiceName:  "ai-service",
		OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		SampleRatio:  1,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry init failed")
	}
	defer shutdownTracer()

	llmClient := llm.New(
		cfg.LLM.BaseURL,
		cfg.LLM.APIKey,
		cfg.LLM.Model,
		cfg.LLM.MaxTokens,
		cfg.LLM.Timeout,
	)

	grpcServer := grpc.NewServer()
	pb.RegisterAIServiceServer(grpcServer, grpchandler.NewAIServer(llmClient))
	reflection.Register(grpcServer)

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatal().Err(err).Int("port", cfg.Server.GRPCPort).Msg("failed to listen on gRPC port")
	}

	go func() {
		log.Info().Int("port", cfg.Server.GRPCPort).Str("model", cfg.LLM.Model).Msg("gRPC listening")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server error")
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
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
	grpcServer.GracefulStop()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("HTTP shutdown error")
	}
	log.Info().Msg("stopped")
}
