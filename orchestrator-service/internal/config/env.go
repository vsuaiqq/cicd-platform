package config

import (
	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
)

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverride("POSTGRES_DSN", &cfg.Postgres.DSN)
	sharedConfig.EnvOverride("ORCHESTRATOR_INTERNAL_API_KEY", &cfg.InternalAPIKey)
	sharedConfig.EnvOverrideBrokers("KAFKA_BROKERS", &cfg.Kafka.Brokers)
	sharedConfig.EnvOverride("PROJECTS_GRPC_ADDRESS", &cfg.Projects.GRPCAddress)
}

func Validate(cfg *Config) error {
	return sharedConfig.RequireNonEmpty(map[string]string{
		"POSTGRES_DSN":                  cfg.Postgres.DSN,
		"ORCHESTRATOR_INTERNAL_API_KEY": cfg.InternalAPIKey,
	})
}
