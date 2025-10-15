// Package chessanalyser provides a comprehensive Go client for Chess.com API with Stockfish engine integration.
//
// This package offers:
//   - Complete Chess.com Published Data API client
//   - Advanced PGN analysis using Stockfish engine
//   - Position evaluation and move suggestions
//   - Game accuracy metrics and blunder detection
//   - Multi-engine concurrent analysis
//   - RESTful API with comprehensive documentation
//
// Example usage:
//
//	import "github.com/pedrampdd/ChessAnalyser/internal/client"
//	import "github.com/pedrampdd/ChessAnalyser/internal/service"
//
//	// Initialize Chess.com API client
//	chessClient := client.NewChessComAPI()
//
//	// Get player profile
//	profile, err := chessClient.GetPlayerProfile("hikaru")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Initialize game analyzer
//	gameAnalyzer := service.NewGameAnalyzerService()
//
//	// Get a game
//	game, err := gameAnalyzer.GetGameByID("hikaru/2024/01")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Game: %+v\n", game)
//
// For more examples and documentation, visit:
// https://github.com/pedrampdd/ChessAnalyser
//
// Keywords: chess.com api golang, chess analysis, stockfish golang, chess api client, pgn analysis
package main
