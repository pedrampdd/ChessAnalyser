package chessanalyser

import (
	"github.com/pedrampdd/ChessAnalyser/internal/client"
	"github.com/pedrampdd/ChessAnalyser/internal/models"
	"github.com/pedrampdd/ChessAnalyser/internal/service"
)

// NewChessComClient creates a new Chess.com API client
func NewChessComClient() *client.ChessComAPI {
	return client.NewChessComAPI()
}

// NewGameAnalyzer creates a new game analyzer service
func NewGameAnalyzer() *service.GameAnalyzerService {
	return service.NewGameAnalyzerService()
}

// NewAnalysisService creates a new analysis service with Stockfish integration
func NewAnalysisService(stockfishPath string, maxEngines int, settings models.EngineSettings) (*service.AnalysisService, error) {
	return service.NewAnalysisService(stockfishPath, maxEngines, settings)
}

// EngineSettings represents Stockfish engine configuration
type EngineSettings = models.EngineSettings

// GameInfo represents chess game information
type GameInfo = models.GameInfo

// AnalysisResult represents game analysis results
type AnalysisResult = models.GameAnalysis
