package config

import (
	"os"
	"sync"
)

type Config struct {
	PORT string
}

var (
	config Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		config.PORT = os.Getenv("PORT")
	})

	return &config
}
