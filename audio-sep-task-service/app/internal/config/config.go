package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App App `mapstructure:"app"`

	Kafka   Kafka   `mapstructure:"kafka"`
	Storage Storage `mapstructure:"storage"`

	Retry Retry `mapstructure:"retry"`

	DatabaseURL string `mapstructure:"database_url"`
}

type App struct {
	Port            string        `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	MigrationDir    string        `mapstructure:"migration_dir"`
	LogLevel        string        `mapstructure:"log_level"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topics  struct {
		AudioSeparation string `mapstructure:"audio_separation"`
	} `mapstructure:"topics"`
}

type Storage struct {
	Endpoint       string        `mapstructure:"endpoint"`
	AccessKey      string        `mapstructure:"access_key"`
	SecretKey      string        `mapstructure:"secret_key"`
	UseSSL         bool          `mapstructure:"use_ssl"`
	Region         string        `mapstructure:"region"`
	AudioBucket    string        `mapstructure:"audio_bucket"`
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
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
