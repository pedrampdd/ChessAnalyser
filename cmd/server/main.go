package main

import (
	"log"

	"chess-analyzer/internal/api"
	"chess-analyzer/internal/config"
	"chess-analyzer/internal/models"
	service "chess-analyzer/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the game analyzer service
	gameService := service.NewGameAnalyzerService()

	// Initialize the analysis service
	defaultSettings := models.EngineSettings{
		Depth:      cfg.Stockfish.DefaultDepth,
		TimeLimit:  cfg.Stockfish.DefaultTimeLimit,
		Threads:    cfg.Stockfish.DefaultThreads,
		HashSize:   cfg.Stockfish.DefaultHashSize,
		SkillLevel: cfg.Stockfish.DefaultSkillLevel,
		Contempt:   cfg.Stockfish.DefaultContempt,
		MultiPV:    1,
	}

	analysisService, err := service.NewAnalysisService(
		cfg.Stockfish.ExecutablePath,
		cfg.Stockfish.MaxEngines,
		defaultSettings,
	)
	if err != nil {
		log.Fatal("Failed to initialize analysis service:", err)
	}
	defer analysisService.Close()

	// Setup routes
	router := api.SetupRoutes(gameService, analysisService)

	// Start the server
	log.Printf("Starting Chess Analyzer API server on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Println("Available endpoints:")
	log.Println("  GET /health - Health check")
	log.Println("  GET /api/game/{gameId} - Get game by ID")
	log.Println("  GET /api/player/{username}/games?year=YYYY&month=MM - Get player's games")
	log.Println("  GET /api/player/{username}/profile - Get player profile")
	log.Println("  GET /api/player/{username}/stats - Get player stats")
	log.Println("  POST /api/analyze/game - Analyze a chess game")
	log.Println("  GET /api/analyze/position?fen=FEN - Analyze a chess position")
	log.Println("  GET /api/analyze/status - Get engine status")
	log.Println("  DELETE /api/analyze/cache - Clear analysis cache")

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
