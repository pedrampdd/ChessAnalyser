package api

import (
	"net/http"
	"strconv"

	"github.com/pedrampdd/ChessAnalyser/internal/models"
	"github.com/pedrampdd/ChessAnalyser/internal/service"
	"github.com/pedrampdd/ChessAnalyser/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Handler represents the API handlers
type Handler struct {
	gameService     *service.GameAnalyzerService
	analysisService *service.AnalysisService
}

// NewHandler creates a new API handler
func NewHandler(gameService *service.GameAnalyzerService, analysisService *service.AnalysisService) *Handler {
	return &Handler{
		gameService:     gameService,
		analysisService: analysisService,
	}
}

// GetGame retrieves game information by ID
func (h *Handler) GetGame(c *gin.Context) {
	gameID := c.Param("gameId")

	gameInfo, err := h.gameService.GetGameByID(gameID)
	if err != nil {
		if _, ok := err.(*errors.GameNotFoundError); ok {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error:   err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	gameResponse := models.GameResponse{
		GameID:      gameInfo.GameID,
		URL:         gameInfo.URL,
		FEN:         gameInfo.FEN,
		PGN:         gameInfo.PGN,
		TimeControl: gameInfo.TimeControl,
		Rules:       gameInfo.Rules,
		WhitePlayer: gameInfo.WhitePlayer,
		BlackPlayer: gameInfo.BlackPlayer,
		Result:      gameInfo.Result,
		ResultCode:  gameInfo.ResultCode,
		TimeClass:   gameInfo.TimeClass,
		Rated:       gameInfo.Rated,
		StartTime:   gameInfo.StartTime,
		EndTime:     gameInfo.EndTime,
		Tournament:  gameInfo.Tournament,
		Match:       gameInfo.Match,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gameResponse,
	})
}

// GetPlayerGames retrieves player's games for a specific month
func (h *Handler) GetPlayerGames(c *gin.Context) {
	username := c.Param("username")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	if yearStr == "" || monthStr == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Year and month parameters are required",
		})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid year parameter",
		})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid month parameter",
		})
		return
	}

	gamesData, err := h.gameService.GetPlayerGames(username, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    gamesData,
	})
}

// GetPlayerProfile retrieves player profile information
func (h *Handler) GetPlayerProfile(c *gin.Context) {
	username := c.Param("username")

	profileData, err := h.gameService.GetPlayerProfile(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    profileData,
	})
}

// GetPlayerStats retrieves player's statistics
func (h *Handler) GetPlayerStats(c *gin.Context) {
	username := c.Param("username")

	statsData, err := h.gameService.GetPlayerStats(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    statsData,
	})
}

// AnalyzeGame analyzes a chess game using Stockfish engine
func (h *Handler) AnalyzeGame(c *gin.Context) {
	var request models.AnalysisRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.AnalysisResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	// Validate required fields
	if request.PGN == "" {
		c.JSON(http.StatusBadRequest, models.AnalysisResponse{
			Success: false,
			Error:   "PGN is required",
		})
		return
	}

	// Set default settings if not provided
	if request.Settings.Depth == 0 {
		request.Settings.Depth = 15
	}
	if request.Settings.TimeLimit == 0 {
		request.Settings.TimeLimit = 5000
	}
	if request.Settings.Threads == 0 {
		request.Settings.Threads = 4
	}
	if request.Settings.HashSize == 0 {
		request.Settings.HashSize = 128
	}

	// Perform analysis
	analysis, err := h.analysisService.AnalyzeGame(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.AnalysisResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AnalysisResponse{
		Success: true,
		Data:    analysis,
		Message: "Game analysis completed successfully",
	})
}

// AnalyzePosition analyzes a single chess position
func (h *Handler) AnalyzePosition(c *gin.Context) {
	fen := c.Query("fen")
	if fen == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "FEN parameter is required",
		})
		return
	}

	// Parse optional settings from query parameters
	settings := models.EngineSettings{
		Depth:     getIntQuery(c, "depth", 15),
		TimeLimit: getIntQuery(c, "time_limit", 5000),
		Threads:   getIntQuery(c, "threads", 4),
		HashSize:  getIntQuery(c, "hash_size", 128),
		MultiPV:   getIntQuery(c, "multipv", 1),
	}

	// Analyze position
	result, err := h.analysisService.AnalyzePosition(c.Request.Context(), fen, settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetEngineStatus returns the status of analysis engines
func (h *Handler) GetEngineStatus(c *gin.Context) {
	status := h.analysisService.GetEngineStatus()
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    status,
	})
}

// ClearAnalysisCache clears the analysis cache
func (h *Handler) ClearAnalysisCache(c *gin.Context) {
	h.analysisService.ClearCache()
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]string{
			"message": "Analysis cache cleared successfully",
		},
	})
}

// HealthCheck provides a health check endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]string{
			"status":  "healthy",
			"service": "chess-analyzer",
		},
	})
}

// getIntQuery gets an integer query parameter with a default value
func getIntQuery(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
