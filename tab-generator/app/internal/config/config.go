package config

import (
	"os"
	"time"

	"github.com/stretchr/testify/assert/yaml"
)

type Config struct {
	App      App        `yaml:"app"`
	Analyzer GrpcClient `yaml:"analyzer"`
}

type App struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type GrpcClient struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

func Load(yamlConfigFilePath string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(yamlConfigFilePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
