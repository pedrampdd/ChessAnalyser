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
//	import "github.com/pedrampdd/ChessAnalyser/pkg/client"
//	import "github.com/pedrampdd/ChessAnalyser/pkg/analysis"
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
//	// Initialize analysis service
//	analysisService, err := analysis.NewService("./stockfish/stockfish", 4)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer analysisService.Close()
//
//	// Analyze a PGN
//	result, err := analysisService.AnalyzePGN(pgnString)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Game accuracy: %.1f%%\n", result.Accuracy.AverageAccuracy)
//
// For more examples and documentation, visit:
// https://github.com/pedrampdd/ChessAnalyser
//
// Keywords: chess.com api golang, chess analysis, stockfish golang, chess api client, pgn analysis
package chessanalyser
