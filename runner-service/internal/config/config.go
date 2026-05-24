package config

type Config struct {
	Kafka  KafkaConfig  `yaml:"kafka"`
	Runner RunnerConfig `yaml:"runner"`
}

type KafkaConfig struct {
	Brokers         []string `yaml:"brokers"`
	JobsTopic       string   `yaml:"jobs_topic"`
	JobResultsTopic string   `yaml:"job_results_topic"`
	CancelJobsTopic string   `yaml:"cancel_jobs_topic"`
	ClientID        string   `yaml:"client_id"`
	GroupID         string   `yaml:"group_id"`
}

type RunnerConfig struct {
	WorkDir       string `yaml:"work_dir"`
	MaxConcurrent int    `yaml:"max_concurrent"`
	DockerSocket  string `yaml:"docker_socket"`
	HTTPPort      int    `yaml:"http_port"`
	ArtifactDir   string `yaml:"artifact_dir"`
	CacheDir      string `yaml:"cache_dir"`
}
