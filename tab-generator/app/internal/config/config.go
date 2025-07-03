package config

import (
	"os"
	"sync"
)

type Config struct {
	PORT         string
	AnalyzerPort string
	AnalyzerHost string
}

var (
	config Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		config.PORT = os.Getenv("PORT")
		config.AnalyzerPort = os.Getenv("ANALYZER_PORT")
		config.AnalyzerHost = os.Getenv("ANALYZER_HOST")
	})

	return &config
}
