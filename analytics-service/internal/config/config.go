package config

import "time"

type Config struct {
	Kafka      KafkaConfig      `yaml:"kafka"`
	Server     ServerConfig     `yaml:"server"`
	ClickHouse ClickHouseConfig `yaml:"clickhouse"`
}

type KafkaConfig struct {
	Brokers         []string `yaml:"brokers"`
	RunEventsTopic  string   `yaml:"run_events_topic"`
	JobResultsTopic string   `yaml:"job_results_topic"`
	ClientID        string   `yaml:"client_id"`
	RunsGroupID     string   `yaml:"runs_group_id"`
	JobsGroupID     string   `yaml:"jobs_group_id"`
}

type ServerConfig struct {
	GRPCPort int `yaml:"grpc_port"`
	HTTPPort int `yaml:"http_port"`
}

type ClickHouseConfig struct {
	Addr            string        `yaml:"addr"`
	Database        string        `yaml:"database"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"-"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	DialTimeout     time.Duration `yaml:"dial_timeout"`
}
