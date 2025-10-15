package api

import (
	service "github.com/pedrampdd/ChessAnalyser/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(gameService *service.GameAnalyzerService, analysisService *service.AnalysisService) *gin.Engine {
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
	handler := NewHandler(gameService, analysisService)

	// Health check endpoint
	r.GET("/health", handler.HealthCheck)

	// API routes
	api := r.Group("/api")
	{
		// Game routes
		api.GET("/game/:gameId", handler.GetGame)
		api.GET("/player/:username/games", handler.GetPlayerGames)
		api.GET("/player/:username/profile", handler.GetPlayerProfile)
		api.GET("/player/:username/stats", handler.GetPlayerStats)

		// Analysis routes
		api.POST("/analyze/game", handler.AnalyzeGame)
		api.GET("/analyze/position", handler.AnalyzePosition)
		api.GET("/analyze/status", handler.GetEngineStatus)
		api.DELETE("/analyze/cache", handler.ClearAnalysisCache)
	}

	return r
}
