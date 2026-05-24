package config

import "time"

type Config struct {
	Kafka    KafkaConfig    `yaml:"kafka"`
	Server   ServerConfig   `yaml:"server"`
	Projects ProjectsConfig `yaml:"projects"`
}

type KafkaConfig struct {
	Brokers     []string `yaml:"brokers"`
	OutputTopic string   `yaml:"output_topic"`
	ClientID    string   `yaml:"client_id"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type ProjectsConfig struct {
	GRPCAddress string        `yaml:"grpc_address"`
	Timeout     time.Duration `yaml:"timeout"`
}
