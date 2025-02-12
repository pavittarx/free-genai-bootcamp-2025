package config

import (
	"os"
	"strconv"
)

// Config holds the configuration for the application
type Config struct {
	DatabaseURL string
	ServerPort  int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Default values
	config := &Config{
		DatabaseURL: "lang_portal.db",
		ServerPort:  8080,
	}

	// Override with environment variables if set
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DatabaseURL = dbURL
	}

	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		config.ServerPort = port
	}

	return config, nil
}
