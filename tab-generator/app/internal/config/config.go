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

	Analyzer GrpcClient `mapstructure:"analyzer"`

	DatabaseURL string `mapstructure:"database_url"`
}

type GrpcClient struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type App struct {
	Port            string        `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	LogLevel        string        `mapstructure:"log_level"`
	MaxWorkers      int64         `mapstructure:"max_workers"`
	Service         Service       `mapstructure:"service"`
}

type Service struct {
	MaxMLRequests int64 `mapstructure:"max_ml_requests"`
	Config        struct {
		TaskTimeout    time.Duration `mapstructure:"task_timeout"`
		DBTimeout      time.Duration `mapstructure:"db_timeout"`
		StorageTimeout time.Duration `mapstructure:"storage_timeout"`
		MLTimeout      time.Duration `mapstructure:"ml_timeout"`
	} `mapstructure:"config"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	Topics  struct {
		TabGenerationStart string `mapstructure:"tab_generation_start"`
	} `mapstructure:"topics"`
	GroupID string `mapstructure:"group_id"`
}

type Storage struct {
	Endpoint    string `mapstructure:"endpoint"`
	AccessKey   string `mapstructure:"access_key"`
	SecretKey   string `mapstructure:"secret_key"`
	UseSSL      bool   `mapstructure:"use_ssl"`
	Region      string `mapstructure:"region"`
	AudioBucket string `mapstructure:"audio_bucket"`
	TabBucket   string `mapstructure:"tab_bucket"`
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
