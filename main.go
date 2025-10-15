package main

import (
	"log"

	"chess-analyzer/internal/api"
	"chess-analyzer/internal/config"
	"chess-analyzer/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the game analyzer service
	gameService := service.NewGameAnalyzerService()

	// Setup routes
	router := api.SetupRoutes(gameService)

	// Start the server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting Chess Analyzer API server on %s", serverAddr)
	log.Println("Available endpoints:")
	log.Println("  GET /health - Health check")
	log.Println("  GET /api/game/{gameId} - Get game by ID")
	log.Println("  GET /api/player/{username}/games?year=YYYY&month=MM - Get player's games")
	log.Println("  GET /api/player/{username}/profile - Get player profile")
	log.Println("  GET /api/player/{username}/stats - Get player stats")

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
