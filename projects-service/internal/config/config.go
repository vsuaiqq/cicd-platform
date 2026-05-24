package config

import "time"

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Projects ProjectsConfig `yaml:"projects"`
}

type ServerConfig struct {
	GRPCPort int `yaml:"grpc_port"`
	HTTPPort int `yaml:"http_port"`
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
	EncryptionKey  string `yaml:"-"`
	WebhookBaseURL string `yaml:"-"`
}
