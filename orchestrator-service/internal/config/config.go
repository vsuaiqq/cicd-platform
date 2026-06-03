package config

import "time"

type Config struct {
	Kafka          KafkaConfig    `yaml:"kafka"`
	Server         ServerConfig   `yaml:"server"`
	Postgres       PostgresConfig `yaml:"postgres"`
	Projects       ProjectsConfig `yaml:"projects"`
	Analytics      AnalyticsConfig `yaml:"analytics"`
	InternalAPIKey string         `yaml:"-"`
}

type KafkaConfig struct {
	Brokers         []string `yaml:"brokers"`
	InputTopic      string   `yaml:"input_topic"`
	JobsTopic       string   `yaml:"jobs_topic"`
	JobResultsTopic string   `yaml:"job_results_topic"`
	CancelJobsTopic string   `yaml:"cancel_jobs_topic"`
	RunEventsTopic  string   `yaml:"run_events_topic"`
	ClientID        string   `yaml:"client_id"`
	GroupID         string   `yaml:"group_id"`
	ResultsGroupID  string   `yaml:"results_group_id"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN             string        `yaml:"-"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout"`
}

type ProjectsConfig struct {
	GRPCAddress string        `yaml:"grpc_address"`
	Timeout     time.Duration `yaml:"timeout"`
}

type AnalyticsConfig struct {
	GRPCAddress string        `yaml:"grpc_address"`
	Timeout     time.Duration `yaml:"timeout"`
}
