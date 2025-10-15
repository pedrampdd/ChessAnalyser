package api

import (
	"chess-analyzer/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(gameService *service.GameAnalyzerService) *gin.Engine {
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Initialize handlers
	handler := NewHandler(gameService)

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	// API routes
	api := r.Group("/api")
	{
		api.GET("/game/:gameId", handler.GetGame)
		api.GET("/player/:username/games", handler.GetPlayerGames)
		api.GET("/player/:username/profile", handler.GetPlayerProfile)
		api.GET("/player/:username/stats", handler.GetPlayerStats)
	}

	return r
}
