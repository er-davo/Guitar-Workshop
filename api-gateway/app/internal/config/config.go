package config

import (
	"os"
	"sync"
)

type Config struct {
	PORT string

	SupabaseURL string
	SupabaseKey string

	DatabaseURL string

	TabgenPort string
	TabgenHost string

	AudioProcPort string
	AudioProcHost string

	AudioSeparatorPort string
	AudioSeparatorHost string
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

		config.DatabaseURL = os.Getenv("DATABASE_URL")

		config.TabgenPort = os.Getenv("TABGEN_PORT")
		config.TabgenHost = os.Getenv("TABGEN_HOST")

		config.AudioProcPort = os.Getenv("AUDIO_PROC_PORT")
		config.AudioProcHost = os.Getenv("AUDIO_PROC_HOST")

		config.AudioSeparatorPort = os.Getenv("AUDIO_SEPARATOR_PORT")
		config.AudioSeparatorHost = os.Getenv("AUDIO_SEPARATOR_HOST")
	})

	return &config
}
