package config

import sharedConfig "github.com/vsuaiqq/cicd/shared/config"

func ApplyEnv(cfg *Config) {
	sharedConfig.EnvOverrideBrokers("KAFKA_BROKERS", &cfg.Kafka.Brokers)
	sharedConfig.EnvOverride("RUNNER_ARTIFACT_DIR", &cfg.Runner.ArtifactDir)
	sharedConfig.EnvOverride("RUNNER_CACHE_DIR", &cfg.Runner.CacheDir)
}

func ApplyDefaults(cfg *Config) {
	if cfg.Runner.WorkDir == "" {
		cfg.Runner.WorkDir = "/tmp/cicd-runner"
	}
	if cfg.Runner.MaxConcurrent <= 0 {
		cfg.Runner.MaxConcurrent = 4
	}
	if cfg.Runner.HTTPPort == 0 {
		cfg.Runner.HTTPPort = 8085
	}
	if cfg.Runner.DockerSocket == "" {
		cfg.Runner.DockerSocket = "/var/run/docker.sock"
	}
	if cfg.Runner.ArtifactDir == "" {
		cfg.Runner.ArtifactDir = "/var/cicd/artifacts"
	}
	if cfg.Runner.CacheDir == "" {
		cfg.Runner.CacheDir = "/var/cicd/cache"
	}
	if cfg.Kafka.CancelJobsTopic == "" {
		cfg.Kafka.CancelJobsTopic = "pipeline.cancel-jobs"
	}
	if cfg.Kafka.JobsTopic == "" {
		cfg.Kafka.JobsTopic = "pipeline.jobs"
	}
	if cfg.Kafka.JobResultsTopic == "" {
		cfg.Kafka.JobResultsTopic = "pipeline.job-results"
	}
	if cfg.Kafka.GroupID == "" {
		cfg.Kafka.GroupID = "runner"
	}
	if cfg.Kafka.ClientID == "" {
		cfg.Kafka.ClientID = "runner"
	}
	if len(cfg.Kafka.Brokers) == 0 {
		cfg.Kafka.Brokers = []string{"kafka:9092"}
	}
}
