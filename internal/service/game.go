package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pedrampdd/ChessAnalyser/internal/client"
	"github.com/pedrampdd/ChessAnalyser/internal/models"
	"github.com/pedrampdd/ChessAnalyser/pkg/errors"
)

// GameAnalyzerService represents the main service for game analysis
type GameAnalyzerService struct {
	chessAPI  *client.ChessComAPI
	gameCache map[string]*models.GameInfo
}

// NewGameAnalyzerService creates a new game analyzer service instance
func NewGameAnalyzerService() *GameAnalyzerService {
	return &GameAnalyzerService{
		chessAPI:  client.NewChessComAPI(),
		gameCache: make(map[string]*models.GameInfo),
	}
}

// GetGameByID retrieves game information by game ID
func (s *GameAnalyzerService) GetGameByID(gameID string) (*models.GameInfo, error) {
	// Check cache first
	if gameInfo, exists := s.gameCache[gameID]; exists {
		return gameInfo, nil
	}

	// Parse game ID and retrieve game information
	gameInfo, err := s.parseGameID(gameID)
	if err != nil {
		return nil, errors.NewGameNotFoundError(gameID, err)
	}

	// Cache the result
	s.gameCache[gameID] = gameInfo
	return gameInfo, nil
}

// GetPlayerGames retrieves player's games for a specific month
func (s *GameAnalyzerService) GetPlayerGames(username string, year, month int) (*models.GameInfo, error) {

	gameData, err := s.chessAPI.GetPlayerGames(username, year, month)
	if err != nil {
		return nil, errors.NewAPIError("failed to retrieve games", err)
	}

	gameInfo, err := s.parseGameData(gameData["games"].([]any)[0].(map[string]any))
	if err != nil {
		return nil, errors.NewAPIError("failed to parse games", err)
	}

	return gameInfo, nil
}

// GetPlayerProfile retrieves player profile information
func (s *GameAnalyzerService) GetPlayerProfile(username string) (map[string]any, error) {
	return s.chessAPI.GetPlayerProfile(username)
}

// GetPlayerStats retrieves player's statistics
func (s *GameAnalyzerService) GetPlayerStats(username string) (map[string]any, error) {
	return s.chessAPI.GetPlayerStats(username)
}

// parseGameID handles different game ID formats
func (s *GameAnalyzerService) parseGameID(gameID string) (*models.GameInfo, error) {
	if strings.HasPrefix(gameID, "http") {
		return s.getGameFromURL(gameID)
	} else if strings.Contains(gameID, "/") {
		parts := strings.Split(gameID, "/")
		if len(parts) >= 3 {
			username := parts[0]
			year, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, errors.NewValidationError("year", fmt.Sprintf("invalid year in game ID: %s", parts[1]))
			}
			month, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, errors.NewValidationError("month", fmt.Sprintf("invalid month in game ID: %s", parts[2]))
			}
			return s.getGameFromPlayerMonth(username, year, month)
		}
	}

	// Try to search for game by ID
	return s.searchGameByID(gameID)
}

// getGameFromURL extracts game information from Chess.com game URL
func (s *GameAnalyzerService) getGameFromURL(url string) (*models.GameInfo, error) {
	// This would need to parse the Chess.com URL structure
	// For now, return an error indicating this feature is not implemented
	return nil, errors.NewAPIError("URL parsing not yet implemented", nil)
}

// getGameFromPlayerMonth gets games from player's monthly archive
func (s *GameAnalyzerService) getGameFromPlayerMonth(username string, year, month int) (*models.GameInfo, error) {
	gamesData, err := s.chessAPI.GetPlayerGames(username, year, month)
	if err != nil {
		return nil, errors.NewAPIError("failed to retrieve games", err)
	}

	// Parse games and return the first one (or implement specific game selection)
	if games, ok := gamesData["games"].([]any); ok && len(games) > 0 {
		gameData := games[0].(map[string]any)
		return s.parseGameData(gameData)
	}

	return nil, errors.NewGameNotFoundError(fmt.Sprintf("%s/%d/%02d", username, year, month), nil)
}

// searchGameByID searches for game by ID across different methods
func (s *GameAnalyzerService) searchGameByID(gameID string) (*models.GameInfo, error) {
	// This would implement a more sophisticated search
	// For now, return an error
	return nil, errors.NewValidationError("gameID", fmt.Sprintf("game ID format not recognized: %s", gameID))
}

// parseGameData parses raw game data from Chess.com API into GameInfo struct
func (s *GameAnalyzerService) parseGameData(gameData map[string]any) (*models.GameInfo, error) {
	// Extract player information
	whiteData, _ := gameData["white"].(map[string]any)
	blackData, _ := gameData["black"].(map[string]any)

	whitePlayer := models.Player{
		Username: getStringValue(whiteData, "username"),
		URL:      getStringValue(whiteData, "url"),
		Avatar:   getStringValue(whiteData, "avatar"),
		Country:  getStringValue(whiteData, "country"),
		Title:    getStringValue(whiteData, "title"),
	}

	if playerID, ok := whiteData["player_id"].(float64); ok {
		id := int(playerID)
		whitePlayer.PlayerID = &id
	}

	blackPlayer := models.Player{
		Username: getStringValue(blackData, "username"),
		URL:      getStringValue(blackData, "url"),
		Avatar:   getStringValue(blackData, "avatar"),
		Country:  getStringValue(blackData, "country"),
		Title:    getStringValue(blackData, "title"),
	}

	if playerID, ok := blackData["player_id"].(float64); ok {
		id := int(playerID)
		blackPlayer.PlayerID = &id
	}

	// Parse timestamps
	startTime := time.Unix(int64(getFloatValue(gameData, "start_time")), 0)
	var endTime *time.Time
	if endTimeVal := getFloatValue(gameData, "end_time"); endTimeVal > 0 {
		et := time.Unix(int64(endTimeVal), 0)
		endTime = &et
	}

	// Create GameInfo object
	gameInfo := &models.GameInfo{
		URL:         getStringValue(gameData, "url"),
		FEN:         getStringValue(gameData, "fen"),
		PGN:         getStringValue(gameData, "pgn"),
		TimeControl: getStringValue(gameData, "time_control"),
		Rules:       getStringValue(gameData, "rules"),
		WhitePlayer: whitePlayer,
		BlackPlayer: blackPlayer,
		Result:      getStringValue(gameData, "result"),
		ResultCode:  getStringValue(gameData, "result_code"),
		TimeClass:   getStringValue(gameData, "time_class"),
		Rated:       getBoolValue(gameData, "rated"),
		StartTime:   startTime,
		EndTime:     endTime,
		Tournament:  getStringValue(gameData, "tournament"),
		Match:       getStringValue(gameData, "match"),
	}

	return gameInfo, nil
}

// Helper functions for type conversion
func getStringValue(data map[string]any, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getFloatValue(data map[string]any, key string) float64 {
	if val, ok := data[key].(float64); ok {
		return val
	}
	return 0
}

func getBoolValue(data map[string]any, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}
