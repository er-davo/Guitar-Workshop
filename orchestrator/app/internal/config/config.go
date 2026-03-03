package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App App `mapstructure:"app"`

	Kafka Kafka `mapstructure:"kafka"`

	Retry Retry `mapstructure:"retry"`

	DatabaseURL string `mapstructure:"database_url"`
}

type App struct {
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	LogLevel        string        `mapstructure:"log_level"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topics  struct {
		AudioSeparationCompleted string `mapstructure:"audio_separation_completed"`
		TabGenerationRequested   string `mapstructure:"tab_generation_requested"`
		TabGenerationStart       string `mapstructure:"tab_generation_start"`
	} `mapstructure:"topics"`
	GroupID string `mapstructure:"group_id"`
}

type Retry struct {
	Backoff     string        `mapstructure:"backoff"`
	MaxAttempts int           `mapstructure:"max_attempts"`
	Base        time.Duration `mapstructure:"base"`
	Factor      float64       `mapstructure:"factor"`
	Max         time.Duration `mapstructure:"max"`
	Jitter      float64       `mapstructure:"jitter"`
}

// Load reads configuration from file or environment variables.
// Config file is optional; environment variables override file values.
func Load(configFilePath string) (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configFilePath != "" {
		v.SetConfigFile(configFilePath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
