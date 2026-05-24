package config

import (
	"time"

	sharedConfig "github.com/vsuaiqq/cicd/shared/config"
)

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverride("CLICKHOUSE_ADDR", &cfg.ClickHouse.Addr)
	sharedConfig.EnvOverride("CLICKHOUSE_USER", &cfg.ClickHouse.Username)
	sharedConfig.EnvOverride("CLICKHOUSE_PASSWORD", &cfg.ClickHouse.Password)
	sharedConfig.EnvOverride("CLICKHOUSE_DATABASE", &cfg.ClickHouse.Database)
	sharedConfig.EnvOverrideBrokers("KAFKA_BROKERS", &cfg.Kafka.Brokers)
}

func ApplyDefaults(cfg *Config) {
	if cfg.Server.GRPCPort == 0 {
		cfg.Server.GRPCPort = 50054
	}
	if cfg.Server.HTTPPort == 0 {
		cfg.Server.HTTPPort = 8087
	}
	if cfg.Kafka.RunEventsTopic == "" {
		cfg.Kafka.RunEventsTopic = "pipeline.run-events"
	}
	if cfg.Kafka.JobResultsTopic == "" {
		cfg.Kafka.JobResultsTopic = "pipeline.job-results"
	}
	if cfg.Kafka.RunsGroupID == "" {
		cfg.Kafka.RunsGroupID = "analytics-runs"
	}
	if cfg.Kafka.JobsGroupID == "" {
		cfg.Kafka.JobsGroupID = "analytics-jobs"
	}
	if cfg.ClickHouse.Addr == "" {
		cfg.ClickHouse.Addr = "clickhouse:9000"
	}
	if cfg.ClickHouse.Database == "" {
		cfg.ClickHouse.Database = "default"
	}
	if cfg.ClickHouse.Username == "" {
		cfg.ClickHouse.Username = "default"
	}
	if cfg.ClickHouse.MaxOpenConns == 0 {
		cfg.ClickHouse.MaxOpenConns = 10
	}
	if cfg.ClickHouse.MaxIdleConns == 0 {
		cfg.ClickHouse.MaxIdleConns = 5
	}
	if cfg.ClickHouse.ConnMaxLifetime == 0 {
		cfg.ClickHouse.ConnMaxLifetime = 5 * time.Minute
	}
	if cfg.ClickHouse.DialTimeout == 0 {
		cfg.ClickHouse.DialTimeout = 10 * time.Second
	}
}

func Validate(cfg *Config) error {
	return sharedConfig.RequireNonEmpty(map[string]string{
		"CLICKHOUSE_PASSWORD": cfg.ClickHouse.Password,
	})
}
