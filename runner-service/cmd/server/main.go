package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vsuaiqq/cicd/runner-service/internal/config"
	runnerKafka "github.com/vsuaiqq/cicd/runner-service/internal/kafka"
	"github.com/vsuaiqq/cicd/runner-service/internal/runner"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	"github.com/vsuaiqq/cicd/shared/events"
	sharedKafka "github.com/vsuaiqq/cicd/shared/kafka"
	"github.com/vsuaiqq/cicd/shared/logger"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var _ runnerKafka.CancelHandler = (*jobHandlerImpl)(nil)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

func main() {
	flag.Parse()
	logger.Init(logger.Config{
		ServiceName: "runner-service",
		Level:       os.Getenv("LOG_LEVEL"),
		Pretty:      os.Getenv("LOG_PRETTY") == "1",
	})
	log := logger.L()

	cfg, err := sharedConfig.LoadService(context.Background(), sharedConfig.ServiceLoader[config.Config]{
		ConfigPath:    *configPath,
		ApplyEnv:      config.ApplyEnv,
		ApplyDefaults: config.ApplyDefaults,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("config failed")
	}

	_, shutdownTracer, err := telemetry.TracerProvider(context.Background(), telemetry.Config{
		ServiceName:  "runner-service",
		OTLPEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
		SampleRatio:  1,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry init failed")
	}
	defer shutdownTracer()

	for _, dir := range []string{cfg.Runner.WorkDir, cfg.Runner.ArtifactDir, cfg.Runner.CacheDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal().Err(err).Str("dir", dir).Msg("create directory failed")
		}
	}

	producer, err := sharedKafka.NewProducer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("kafka producer init failed")
	}
	defer producer.Close()

	resultPub := runnerKafka.NewResultPublisher(producer, cfg.Kafka.JobResultsTopic)

	jobConsumerClient, err := sharedKafka.NewConsumer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID,
		GroupID:  cfg.Kafka.GroupID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("kafka job consumer init failed")
	}

	execCfg := &runner.ExecConfig{
		DockerSocket: cfg.Runner.DockerSocket,
		ArtifactDir:  cfg.Runner.ArtifactDir,
		CacheDir:     cfg.Runner.CacheDir,
	}

	var wg sync.WaitGroup
	cancelReg := &runner.CancelRegistry{}
	jobHandler := &jobHandlerImpl{
		sem:       make(chan struct{}, cfg.Runner.MaxConcurrent),
		resultPub: resultPub,
		workDir:   cfg.Runner.WorkDir,
		execCfg:   execCfg,
		wg:        &wg,
		cancelReg: cancelReg,
	}
	jobConsumer := runnerKafka.NewJobConsumer(jobConsumerClient, jobHandler)

	hostname, err := os.Hostname()
	if err != nil {
		log.Warn().Err(err).Msg("cannot determine hostname; falling back to pid for cancel group")
		hostname = fmt.Sprintf("pid-%d", os.Getpid())
	}
	cancelGroupID := cfg.Kafka.ClientID + "-cancel-" + hostname

	cancelConsumerClient, err := sharedKafka.NewConsumer(sharedKafka.Config{
		Brokers:  cfg.Kafka.Brokers,
		ClientID: cfg.Kafka.ClientID,
		GroupID:  cancelGroupID,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("kafka cancel consumer init failed")
	}
	cancelConsumer := runnerKafka.NewCancelConsumer(cancelConsumerClient, jobHandler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Runner.HTTPPort),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	go func() {
		log.Info().Int("port", cfg.Runner.HTTPPort).Msg("HTTP healthz listening")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server error")
		}
	}()

	go func() {
		log.Info().Str("topic", cfg.Kafka.JobsTopic).Str("group", cfg.Kafka.GroupID).Msg("job consumer starting")
		if err := jobConsumer.Start(ctx, []string{cfg.Kafka.JobsTopic}); err != nil && ctx.Err() == nil {
			log.Error().Err(err).Msg("job consumer error")
		}
	}()

	go func() {
		log.Info().Str("topic", cfg.Kafka.CancelJobsTopic).Str("group", cancelGroupID).Msg("cancel consumer starting")
		if err := cancelConsumer.Start(ctx, []string{cfg.Kafka.CancelJobsTopic}); err != nil && ctx.Err() == nil {
			log.Error().Err(err).Msg("cancel consumer error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down: waiting for in-progress jobs")
	cancel()
	wg.Wait()
	_ = jobConsumer.Close()
	_ = cancelConsumer.Close()

	shutCtx, shutCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutCancel()
	_ = httpServer.Shutdown(shutCtx)

	log.Info().Msg("runner stopped")
}

type jobHandlerImpl struct {
	sem       chan struct{}
	resultPub *runnerKafka.ResultPublisher
	workDir   string
	execCfg   *runner.ExecConfig
	wg        *sync.WaitGroup
	cancelReg *runner.CancelRegistry
}

func (h *jobHandlerImpl) HandleJob(ctx context.Context, task *events.PipelineJobTask) error {
	h.sem <- struct{}{}
	h.wg.Add(1)

	jobCtx, jobCancel := context.WithCancel(ctx)
	h.cancelReg.Register(task.JobID, jobCancel)

	go func() {
		defer func() { <-h.sem }()
		defer h.wg.Done()
		defer h.cancelReg.Unregister(task.JobID)
		defer jobCancel()

		logger.L().Info().Str("job", task.JobName).Str("run_id", task.RunID).Msg("job started")
		result := runner.Execute(jobCtx, task, h.workDir, h.execCfg)

		pubCtx, pubCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer pubCancel()
		if err := h.resultPub.Publish(pubCtx, result); err != nil {
			logger.L().Error().Err(err).Str("job_id", task.JobID).Msg("failed to publish result")
		}
		logger.L().Info().Str("job", task.JobName).Str("status", result.Status).Msg("job finished")
	}()

	return nil
}

func (h *jobHandlerImpl) HandleCancel(evt *events.CancelJobEvent) {
	if evt.JobID == "" {
		logger.L().Debug().Str("run_id", evt.RunID).Msg("run-level cancel received; no per-run index on runner")
		return
	}
	if h.cancelReg.Cancel(evt.JobID) {
		logger.L().Info().Str("job_id", evt.JobID).Str("run_id", evt.RunID).Msg("job cancelled")
	}
}
