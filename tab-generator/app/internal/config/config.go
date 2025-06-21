package config

import (
	"os"
	"sync"
)

type Config struct {
	PORT                string
	OnsetsAndFramesPort string
	OnsetsAndFramesHost string
}

var (
	config Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		config.PORT = os.Getenv("PORT")
		config.OnsetsAndFramesPort = os.Getenv("ONSETS_AND_FRAMES_PORT")
		config.OnsetsAndFramesHost = os.Getenv("ONSETS_AND_FRAMES_HOST")
	})

	return &config
}
