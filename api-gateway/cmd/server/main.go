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
	"github.com/go-chi/cors"

	"github.com/vsuaiqq/cicd/api-gateway/internal/client"
	"github.com/vsuaiqq/cicd/api-gateway/internal/config"
	"github.com/vsuaiqq/cicd/api-gateway/internal/handler"
	sharedChi "github.com/vsuaiqq/cicd/shared/chi"
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
	sharedHTTP "github.com/vsuaiqq/cicd/shared/httputil"
	"github.com/vsuaiqq/cicd/shared/logger"
	"github.com/vsuaiqq/cicd/shared/telemetry"
)

var configPath = flag.String("config", "./configs/config.yaml", "config file path")

const serviceName = "api-gateway"

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
		ServiceName: "api-gateway",
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

	authClient, err := client.NewAuthClient(cfg.AuthService.GRPCAddress, cfg.AuthService.Timeout)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to auth service")
	}
	defer authClient.Close()
	log.Info().Msg("connected to auth service")

	projectsClient, err := client.NewProjectsClient(cfg.ProjectsService.GRPCAddress, cfg.ProjectsService.Timeout)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to projects service")
	}
	defer projectsClient.Close()
	log.Info().Msg("connected to projects service")

	orchTimeout := cfg.OrchestratorService.Timeout
	if orchTimeout == 0 {
		orchTimeout = 15 * time.Second
	}
	pipelineClient := client.NewPipelineClient(cfg.OrchestratorService.BaseURL, orchTimeout, cfg.InternalAPIKey)

	aiTimeout := cfg.AIService.Timeout
	if aiTimeout == 0 {
		aiTimeout = 120 * time.Second
	}
	aiClient, err := client.NewAIClient(cfg.AIService.GRPCAddress, aiTimeout)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to ai service")
	}
	defer aiClient.Close()
	log.Info().Msg("connected to ai service")

	analyticsTimeout := cfg.AnalyticsService.Timeout
	if analyticsTimeout == 0 {
		analyticsTimeout = 10 * time.Second
	}
	analyticsClient, err := client.NewAnalyticsClient(cfg.AnalyticsService.GRPCAddress, analyticsTimeout)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to analytics service")
	}
	defer analyticsClient.Close()
	log.Info().Msg("connected to analytics service")

	allowedOrigins := cfg.Server.AllowedOrigins
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}

	authHandler := handler.NewAuthHandler(authClient)
	projectsHandler := handler.NewProjectsHandler(authClient, projectsClient)
	pipelineHandler := handler.NewPipelineHandler(authClient, pipelineClient, projectsClient)
	aiHandler := handler.NewAIHandler(authClient, pipelineClient, projectsClient, aiClient)
	analyticsHandler := handler.NewAnalyticsHandler(authClient, analyticsClient, projectsClient)
	wsProxyHandler := handler.NewWSProxyHandler(pipelineClient, authClient, projectsClient, allowedOrigins)

	router := setupRouter(authHandler, projectsHandler, pipelineHandler, aiHandler, analyticsHandler, wsProxyHandler, allowedOrigins)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info().Int("port", cfg.Server.Port).Msg("starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
	}

	log.Info().Msg("stopped")
}

func setupRouter(
	authHandler *handler.AuthHandler,
	projectsHandler *handler.ProjectsHandler,
	pipelineHandler *handler.PipelineHandler,
	aiHandler *handler.AIHandler,
	analyticsHandler *handler.AnalyticsHandler,
	wsProxyHandler *handler.WSProxyHandler,
	allowedOrigins []string,
) *chi.Mux {
	router := chi.NewRouter()

	router.Use(telemetry.Middleware(serviceName))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Group(func(r chi.Router) {
		r.Get("/api/v1/pipeline/ws/runs/{runID}", wsProxyHandler.ServeWS)
		r.Get("/api/v1/pipeline/ws/events", wsProxyHandler.ServeGlobalWS)
	})

	router.Group(func(r chi.Router) {
		for _, mw := range sharedChi.DefaultStack(30 * time.Second) {
			r.Use(mw)
		}

		r.Get("/healthz", healthCheck)
		r.Get("/readyz", readinessCheck(authHandler))

		r.Route("/api/v1", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
			projectsHandler.RegisterRoutes(r)
			pipelineHandler.RegisterRoutes(r)
			analyticsHandler.RegisterRoutes(r)
		})
	})

	router.Group(func(r chi.Router) {
		for _, mw := range sharedChi.DefaultStack(120 * time.Second) {
			r.Use(mw)
		}
		r.Route("/api/v1/ai", func(r chi.Router) {
			aiHandler.RegisterRoutes(r)
		})
	})

	router.NotFound(notFoundHandler)
	router.MethodNotAllowed(methodNotAllowedHandler)

	return router
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func readinessCheck(authHandler *handler.AuthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if err := authHandler.Healthcheck(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("not ready"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	}
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	sharedHTTP.JSON(w, http.StatusNotFound, sharedHTTP.ErrorResponse{
		Error:   "Not Found",
		Message: "the requested resource was not found",
	})
}

func methodNotAllowedHandler(w http.ResponseWriter, _ *http.Request) {
	sharedHTTP.JSON(w, http.StatusMethodNotAllowed, sharedHTTP.ErrorResponse{
		Error:   "Method Not Allowed",
		Message: "the requested method is not allowed for this resource",
	})
}
