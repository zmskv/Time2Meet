package config

import (
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host    string
	Port    int
	Name    string
	User    string
	Pass    string
	SSLMode string
}

type HTTPConfig struct {
	Addr string
}

type Config struct {
	Database DatabaseConfig
	HTTP     HTTPConfig
}

func LoadFromEnv() (Config, error) {
	var cfg Config

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
	cfg.Database.Name = os.Getenv("DB_NAME")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Pass = os.Getenv("DB_PASSWORD")

	portStr := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return Config{}, fmt.Errorf("invalid DB_PORT: %q", portStr)
	}
	cfg.Database.Port = port

	cfg.HTTP.Addr = getEnv("HTTP_ADDR", ":8080")

	if cfg.Database.Name == "" {
		return Config{}, fmt.Errorf("DB_NAME is required")
	}
	if cfg.Database.User == "" {
		return Config{}, fmt.Errorf("DB_USER is required")
	}
	if cfg.Database.Pass == "" {
		return Config{}, fmt.Errorf("DB_PASSWORD is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
