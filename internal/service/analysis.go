package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pedrampdd/ChessAnalyser/internal/engine"
	"github.com/pedrampdd/ChessAnalyser/internal/models"
	"github.com/pedrampdd/ChessAnalyser/internal/parser"
	"github.com/pedrampdd/ChessAnalyser/pkg/errors"
)

// AnalysisService provides chess game analysis using Stockfish engine
type AnalysisService struct {
	enginePool      *engine.EnginePool
	pgnParser       *parser.PGNParser
	cache           map[string]*models.GameAnalysis
	cacheMutex      sync.RWMutex
	defaultSettings models.EngineSettings
	maxCacheSize    int
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(executablePath string, maxEngines int, defaultSettings models.EngineSettings) (*AnalysisService, error) {
	enginePool, err := engine.NewEnginePool(maxEngines, executablePath, defaultSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to create engine pool: %w", err)
	}

	return &AnalysisService{
		enginePool:      enginePool,
		pgnParser:       parser.NewPGNParser(),
		cache:           make(map[string]*models.GameAnalysis),
		defaultSettings: defaultSettings,
		maxCacheSize:    1000, // Maximum cached analyses
	}, nil
}

// AnalyzeGame analyzes a complete chess game
func (s *AnalysisService) AnalyzeGame(ctx context.Context, request *models.AnalysisRequest) (*models.GameAnalysis, error) {
	// Check cache first
	cacheKey := s.generateCacheKey(request)
	if cached := s.getFromCache(cacheKey); cached != nil {
		return cached, nil
	}

	// Validate PGN
	if err := s.pgnParser.ValidatePGN(request.PGN); err != nil {
		return nil, errors.NewValidationError("pgn", err.Error())
	}

	// Parse PGN
	parsedGame, err := s.pgnParser.ParsePGN(request.PGN)
	if err != nil {
		return nil, errors.NewValidationError("pgn", fmt.Sprintf("failed to parse PGN: %v", err))
	}

	// Extract positions
	if err := s.pgnParser.ExtractPositions(parsedGame); err != nil {
		return nil, errors.NewAPIError("failed to extract positions", err)
	}

	// Perform analysis
	analysis, err := s.performGameAnalysis(ctx, parsedGame, request.Settings, request.MaxMoves)
	if err != nil {
		return nil, errors.NewAPIError("analysis failed", err)
	}

	// Cache the result
	s.addToCache(cacheKey, analysis)

	return analysis, nil
}

// performGameAnalysis performs the actual game analysis
func (s *AnalysisService) performGameAnalysis(ctx context.Context, game *parser.ParsedGame, settings models.EngineSettings, maxMoves int) (*models.GameAnalysis, error) {
	startTime := time.Now()

	// Get engine from pool
	stockfishEngine := s.enginePool.GetEngine()
	defer s.enginePool.ReturnEngine(stockfishEngine)

	// Initialize analysis result
	analysis := &models.GameAnalysis{
		GameID:         game.Headers["gameid"],
		PGN:            game.PGN,
		AnalysisTime:   startTime,
		EngineVersion:  stockfishEngine.GetVersion(),
		EngineSettings: settings,
		Moves:          make([]models.MoveAnalysis, 0, len(game.Moves)),
		Accuracy:       models.GameAccuracy{},
		Summary:        models.AnalysisSummary{},
	}

	// Determine how many moves to analyze
	movesToAnalyze := len(game.Moves)
	if maxMoves > 0 && maxMoves < movesToAnalyze {
		movesToAnalyze = maxMoves
	}

	// Analyze each move
	var totalNodes int64
	var totalTime int64
	var whiteBlunders, blackBlunders int
	var whiteMistakes, blackMistakes int
	var whiteInaccuracies, blackInaccuracies int
	var whiteBestMoves, blackBestMoves int

	for i := 0; i < movesToAnalyze; i++ {
		move := game.Moves[i]

		// Analyze the position after this move
		result, err := stockfishEngine.AnalyzePosition(ctx, move.FEN, settings)
		if err != nil {
			// Continue with next move if analysis fails
			continue
		}

		// Create move analysis
		moveAnalysis := s.createMoveAnalysis(move, result, i+1)
		analysis.Moves = append(analysis.Moves, moveAnalysis)

		// Update statistics
		totalNodes += result.Nodes
		totalTime += result.Time

		// Count move quality
		if move.Color == "white" {
			if moveAnalysis.Blunder {
				whiteBlunders++
			} else if moveAnalysis.Mistake {
				whiteMistakes++
			} else if moveAnalysis.Inaccuracy {
				whiteInaccuracies++
			} else if moveAnalysis.Accuracy >= 95 {
				whiteBestMoves++
			}
		} else {
			if moveAnalysis.Blunder {
				blackBlunders++
			} else if moveAnalysis.Mistake {
				blackMistakes++
			} else if moveAnalysis.Inaccuracy {
				blackInaccuracies++
			} else if moveAnalysis.Accuracy >= 95 {
				blackBestMoves++
			}
		}
	}

	// Calculate final statistics
	s.calculateGameStatistics(analysis, totalNodes, totalTime,
		whiteBlunders, blackBlunders, whiteMistakes, blackMistakes,
		whiteInaccuracies, blackInaccuracies, whiteBestMoves, blackBestMoves)

	return analysis, nil
}

// createMoveAnalysis creates a MoveAnalysis from a ParsedMove and AnalysisResult
func (s *AnalysisService) createMoveAnalysis(move parser.ParsedMove, result *models.AnalysisResult, moveNumber int) models.MoveAnalysis {
	// Calculate move accuracy based on evaluation
	accuracy := s.calculateMoveAccuracy(result.Evaluation)

	// Determine move quality
	blunder := accuracy < 50
	mistake := accuracy >= 50 && accuracy < 80
	inaccuracy := accuracy >= 80 && accuracy < 90

	// Get alternative moves (simplified for now)
	alternatives := make([]models.MoveAlternative, 0)
	if len(result.PrincipalVariation) > 1 {
		alt := models.MoveAlternative{
			Move:       result.PrincipalVariation[0],
			Evaluation: result.Evaluation,
			Depth:      result.Depth,
		}
		alternatives = append(alternatives, alt)
	}

	return models.MoveAnalysis{
		Move:         move.Move,
		MoveNumber:   moveNumber,
		Evaluation:   result.Evaluation,
		Accuracy:     accuracy,
		Blunder:      blunder,
		Mistake:      mistake,
		Inaccuracy:   inaccuracy,
		BestMove:     result.BestMove,
		Alternatives: alternatives,
	}
}

// calculateMoveAccuracy calculates the accuracy percentage for a move
func (s *AnalysisService) calculateMoveAccuracy(evaluation float64) float64 {
	// This is a simplified accuracy calculation
	// In practice, you'd compare against the best move evaluation
	if evaluation >= 0 {
		return 100.0 - (evaluation * 10) // Penalize positive evaluations less
	} else {
		return 100.0 + (evaluation * 15) // Penalize negative evaluations more
	}
}

// calculateGameStatistics calculates overall game statistics
func (s *AnalysisService) calculateGameStatistics(analysis *models.GameAnalysis, totalNodes, totalTime int64,
	whiteBlunders, blackBlunders, whiteMistakes, blackMistakes, whiteInaccuracies, blackInaccuracies, whiteBestMoves, blackBestMoves int) {

	totalMoves := len(analysis.Moves)
	if totalMoves == 0 {
		return
	}

	// Calculate accuracies
	whiteMoves := 0
	blackMoves := 0
	var whiteAccuracySum, blackAccuracySum float64

	for _, move := range analysis.Moves {
		if move.MoveNumber%2 == 1 { // White moves
			whiteMoves++
			whiteAccuracySum += move.Accuracy
		} else { // Black moves
			blackMoves++
			blackAccuracySum += move.Accuracy
		}
	}

	analysis.Accuracy.WhiteAccuracy = whiteAccuracySum / float64(whiteMoves)
	analysis.Accuracy.BlackAccuracy = blackAccuracySum / float64(blackMoves)
	analysis.Accuracy.AverageAccuracy = (whiteAccuracySum + blackAccuracySum) / float64(totalMoves)
	analysis.Accuracy.Blunders = whiteBlunders + blackBlunders
	analysis.Accuracy.Mistakes = whiteMistakes + blackMistakes
	analysis.Accuracy.Inaccuracies = whiteInaccuracies + blackInaccuracies
	analysis.Accuracy.BestMoves = whiteBestMoves + blackBestMoves

	// Calculate summary
	analysis.Summary.TotalMoves = totalMoves
	analysis.Summary.TotalTime = totalTime
	analysis.Summary.NodesSearched = totalNodes
	analysis.Summary.GamePhase = s.determineGamePhase(totalMoves)
	analysis.Summary.Complexity = s.determineComplexity(analysis.Accuracy.AverageAccuracy)
	analysis.Summary.Recommendations = s.generateRecommendations(analysis)
}

// determineGamePhase determines the game phase based on move count
func (s *AnalysisService) determineGamePhase(moveCount int) string {
	if moveCount <= 20 {
		return "opening"
	} else if moveCount <= 40 {
		return "middlegame"
	} else {
		return "endgame"
	}
}

// determineComplexity determines game complexity based on accuracy
func (s *AnalysisService) determineComplexity(accuracy float64) string {
	if accuracy >= 90 {
		return "low"
	} else if accuracy >= 75 {
		return "medium"
	} else {
		return "high"
	}
}

// generateRecommendations generates analysis recommendations
func (s *AnalysisService) generateRecommendations(analysis *models.GameAnalysis) []string {
	var recommendations []string

	if analysis.Accuracy.Blunders > 5 {
		recommendations = append(recommendations, "Consider spending more time on tactical calculations to reduce blunders")
	}

	if analysis.Accuracy.Mistakes > 10 {
		recommendations = append(recommendations, "Focus on positional understanding to minimize mistakes")
	}

	if analysis.Accuracy.AverageAccuracy < 80 {
		recommendations = append(recommendations, "Overall game accuracy could be improved with more careful move selection")
	}

	if analysis.Summary.GamePhase == "opening" && analysis.Accuracy.AverageAccuracy < 85 {
		recommendations = append(recommendations, "Study opening theory to improve early game play")
	}

	return recommendations
}

// generateCacheKey generates a cache key for the analysis request
func (s *AnalysisService) generateCacheKey(request *models.AnalysisRequest) string {
	return fmt.Sprintf("%s_%d_%d_%d",
		request.PGN,
		request.Settings.Depth,
		request.Settings.TimeLimit,
		request.MaxMoves)
}

// getFromCache retrieves analysis from cache
func (s *AnalysisService) getFromCache(key string) *models.GameAnalysis {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()
	return s.cache[key]
}

// addToCache adds analysis to cache
func (s *AnalysisService) addToCache(key string, analysis *models.GameAnalysis) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	// Simple cache eviction if cache is full
	if len(s.cache) >= s.maxCacheSize {
		// Remove oldest entry (simplified)
		for k := range s.cache {
			delete(s.cache, k)
			break
		}
	}

	s.cache[key] = analysis
}

// AnalyzePosition analyzes a single chess position
func (s *AnalysisService) AnalyzePosition(ctx context.Context, fen string, settings models.EngineSettings) (*models.AnalysisResult, error) {
	stockfishEngine := s.enginePool.GetEngine()
	defer s.enginePool.ReturnEngine(stockfishEngine)

	return stockfishEngine.AnalyzePosition(ctx, fen, settings)
}

// GetEngineStatus returns the status of engines in the pool
func (s *AnalysisService) GetEngineStatus() map[string]interface{} {
	return map[string]interface{}{
		"total_engines":     len(s.enginePool.Engines),
		"available_engines": len(s.enginePool.Available),
		"cache_size":        len(s.cache),
		"max_cache_size":    s.maxCacheSize,
	}
}

// ClearCache clears the analysis cache
func (s *AnalysisService) ClearCache() {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.cache = make(map[string]*models.GameAnalysis)
}

// Close shuts down the analysis service
func (s *AnalysisService) Close() error {
	return s.enginePool.Close()
}
