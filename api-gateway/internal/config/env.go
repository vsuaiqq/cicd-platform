package config

import sharedConfig "github.com/vsuaiqq/cicd/shared/config"

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverride("ORCHESTRATOR_INTERNAL_API_KEY", &cfg.InternalAPIKey)
	sharedConfig.EnvOverride("AUTH_SERVICE_GRPC_ADDRESS", &cfg.AuthService.GRPCAddress)
	sharedConfig.EnvOverride("PROJECTS_SERVICE_GRPC_ADDRESS", &cfg.ProjectsService.GRPCAddress)
	sharedConfig.EnvOverride("ORCHESTRATOR_BASE_URL", &cfg.OrchestratorService.BaseURL)
	sharedConfig.EnvOverride("AI_SERVICE_GRPC_ADDRESS", &cfg.AIService.GRPCAddress)
	sharedConfig.EnvOverride("ANALYTICS_SERVICE_GRPC_ADDRESS", &cfg.AnalyticsService.GRPCAddress)
}

func Validate(cfg *Config) error {
	return sharedConfig.RequireNonEmpty(map[string]string{
		"ORCHESTRATOR_INTERNAL_API_KEY": cfg.InternalAPIKey,
	})
}
