package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	ChessAPI ChessAPIConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// ChessAPIConfig holds Chess.com API configuration
type ChessAPIConfig struct {
	BaseURL   string
	UserAgent string
	Timeout   int
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		ChessAPI: ChessAPIConfig{
			BaseURL:   getEnv("CHESS_API_BASE_URL", "https://api.chess.com/pub"),
			UserAgent: getEnv("CHESS_API_USER_AGENT", "ChessAnalyzer/1.0"),
			Timeout:   getEnvAsInt("CHESS_API_TIMEOUT", 30),
		},
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
