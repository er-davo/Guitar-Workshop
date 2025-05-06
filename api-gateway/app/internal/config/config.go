package config

import (
	"os"
	"sync"
)

type Config struct {
	PORT        string
	SupabaseURL string
	SupabaseKey string
}

var (
	config Config
	once   sync.Once
)

func Load() *Config {
	once.Do(func() {
		config.PORT = os.Getenv("PORT")
		config.SupabaseURL = os.Getenv("SUPABASE_URL")
		config.SupabaseKey = os.Getenv("ACCESS_KEY")
	})

	return &config
}
