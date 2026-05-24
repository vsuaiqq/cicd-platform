package config

import sharedConfig "github.com/vsuaiqq/cicd/shared/config"

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverrideBrokers("KAFKA_BROKERS", &cfg.Kafka.Brokers)
	sharedConfig.EnvOverride("PROJECTS_GRPC_ADDRESS", &cfg.Projects.GRPCAddress)
}
