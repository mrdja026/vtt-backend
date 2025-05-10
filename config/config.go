package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DBPath      string
	JWTSecret   string
	Environment string
	Port        string
	SRDAPIBaseURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env file if it exists
	godotenv.Load()

	config := &Config{
		Environment: getEnv("ENV", "development"),
		Port:        getEnv("PORT", "8000"),
		SRDAPIBaseURL: getEnv("SRD_API_BASE_URL", "https://www.dnd5eapi.co/api"),
	}

	// JWT Secret is required
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// In development, use a default secret
		if config.Environment == "development" {
			jwtSecret = "dev-secret-do-not-use-in-production"
		} else {
			return nil, errors.New("JWT_SECRET environment variable is required")
		}
	}
	config.JWTSecret = jwtSecret

	// Database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Use a default path
		dataDir := "./data"
		
		// Create data directory if it doesn't exist
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, errors.New("Failed to create data directory")
		}
		
		dbPath = filepath.Join(dataDir, "dnd_combat.db")
	}
	config.DBPath = dbPath

	return config, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
