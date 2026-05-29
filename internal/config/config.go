package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
	Port         string
}

func Load() Config {
	return Config{
		DatabaseURL:  envOrDefault("DATABASE_URL", "postgres://doheem_dev_user:simple_pswd@localhost:5432/doheem_dev_db?sslmode=disable"),
		JWTSecret:    envOrDefault("JWT_SECRET", "doheem-dev-secret-change-in-production"),
		JWTExpiresIn: envDurationOrDefault("JWT_EXPIRES_IN", 24*time.Hour),
		Port:         envOrDefault("PORT", "8080"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return fallback
}
