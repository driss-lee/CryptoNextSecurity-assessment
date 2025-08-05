package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	StorageMaxSize   int
	SniffingInterval time.Duration
	ServerPort       string
	ShutdownTimeout  time.Duration
}

// Load loads configuration from .env file and environment variables
func Load() *Config {
	// Load appropriate .env file
	loadEnvFile(getEnvWithDefault("ENV", "development"))

	// Return config with environment variables (override .env file values)
	return &Config{
		StorageMaxSize:   getEnvIntWithDefault("STORAGE_MAX_SIZE", 1000),
		SniffingInterval: getEnvDurationWithDefault("SNIFFING_INTERVAL", 5*time.Second),
		ServerPort:       getEnvWithDefault("SERVER_PORT", "8080"),
		ShutdownTimeout:  getEnvDurationWithDefault("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
	}
}

// loadEnvFile loads the appropriate .env file based on environment
func loadEnvFile(env string) {
	// Determine file to load
	envFile := ".env.development"
	if env == "production" {
		envFile = ".env.production"
	}

	// Try to load the environment-specific file
	if err := godotenv.Load(envFile); err == nil {
		log.Printf("Loaded configuration from %s", envFile)
		return
	}

	// Fallback for production: try development file
	if env == "production" {
		if err := godotenv.Load(".env.development"); err == nil {
			log.Println("No .env.production found, using .env.development")
			return
		}
	}

	log.Println("No .env file found, using environment variables only")
}

// getEnvWithDefault returns environment variable value or default if not set
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntWithDefault returns environment variable as int or default if not set
func getEnvIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDurationWithDefault returns environment variable as duration or default if not set
func getEnvDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
