package config

import "time"

type Config struct {
	Server ServerConfig `yaml:"server"`
	LLM    LLMConfig    `yaml:"llm"`
}

type ServerConfig struct {
	HTTPPort int `yaml:"http_port"`
	GRPCPort int `yaml:"grpc_port"`
}

type LLMConfig struct {
	BaseURL string `yaml:"-"`
	APIKey  string `yaml:"-"`
	Model   string `yaml:"-"`

	MaxTokens int `yaml:"max_tokens"`

	Timeout time.Duration `yaml:"timeout"`
}
