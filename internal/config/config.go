package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	ChessAPI  ChessAPIConfig
	Stockfish StockfishConfig
	Analysis  AnalysisConfig
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

// StockfishConfig holds Stockfish engine configuration
type StockfishConfig struct {
	ExecutablePath    string
	MaxEngines        int
	DefaultDepth      int
	DefaultTimeLimit  int
	DefaultThreads    int
	DefaultHashSize   int
	DefaultSkillLevel int
	DefaultContempt   int
}

// AnalysisConfig holds analysis service configuration
type AnalysisConfig struct {
	MaxCacheSize       int
	CacheExpiration    int // in minutes
	MaxMovesPerGame    int
	EnableCaching      bool
	ConcurrentAnalysis bool
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
		Stockfish: StockfishConfig{
			ExecutablePath:    getEnv("STOCKFISH_PATH", "./stockfish/stockfish"),
			MaxEngines:        getEnvAsInt("STOCKFISH_MAX_ENGINES", 4),
			DefaultDepth:      getEnvAsInt("STOCKFISH_DEFAULT_DEPTH", 15),
			DefaultTimeLimit:  getEnvAsInt("STOCKFISH_DEFAULT_TIME_LIMIT", 5000), // 5 seconds
			DefaultThreads:    getEnvAsInt("STOCKFISH_DEFAULT_THREADS", 4),
			DefaultHashSize:   getEnvAsInt("STOCKFISH_DEFAULT_HASH_SIZE", 128), // 128 MB
			DefaultSkillLevel: getEnvAsInt("STOCKFISH_DEFAULT_SKILL_LEVEL", 20),
			DefaultContempt:   getEnvAsInt("STOCKFISH_DEFAULT_CONTEMPT", 0),
		},
		Analysis: AnalysisConfig{
			MaxCacheSize:       getEnvAsInt("ANALYSIS_MAX_CACHE_SIZE", 1000),
			CacheExpiration:    getEnvAsInt("ANALYSIS_CACHE_EXPIRATION", 60), // 60 minutes
			MaxMovesPerGame:    getEnvAsInt("ANALYSIS_MAX_MOVES_PER_GAME", 100),
			EnableCaching:      getEnvAsBool("ANALYSIS_ENABLE_CACHING", true),
			ConcurrentAnalysis: getEnvAsBool("ANALYSIS_CONCURRENT", true),
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

// getEnvAsBool gets an environment variable as boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
