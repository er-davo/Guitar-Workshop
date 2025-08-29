package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App App `yaml:"app"`

	DatabaseURL string

	Tabgen   GrpcClient `yaml:"tabgen"`
	AudioSep GrpcClient `yaml:"audiosep"`

	Supabase Supabase
}

type App struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type GrpcClient struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type Supabase struct {
	URL string `yaml:"url"`
	Key string `yaml:"key"`
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
	cfg.Supabase.URL = os.Getenv("SUPABASE_URL")
	cfg.Supabase.Key = os.Getenv("ACCESS_KEY")

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")

	return cfg, nil
}
