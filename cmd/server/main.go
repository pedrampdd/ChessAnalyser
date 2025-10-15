package main

import (
	"log"

	"chess-analyzer/internal/api"
	"chess-analyzer/internal/service"
)

func main() {
	// Initialize the game analyzer service
	gameService := service.NewGameAnalyzerService()

	// Setup routes
	router := api.SetupRoutes(gameService)

	// Start the server
	log.Println("Starting Chess Analyzer API server on :8080")
	log.Println("Available endpoints:")
	log.Println("  GET /health - Health check")
	log.Println("  GET /api/game/{gameId} - Get game by ID")
	log.Println("  GET /api/player/{username}/games?year=YYYY&month=MM - Get player's games")
	log.Println("  GET /api/player/{username}/profile - Get player profile")
	log.Println("  GET /api/player/{username}/stats - Get player stats")

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
