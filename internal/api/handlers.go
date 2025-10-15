package api

import (
	"net/http"
	"strconv"

	"chess-analyzer/internal/models"
	"chess-analyzer/internal/service"
	"chess-analyzer/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Handler represents the API handlers
type Handler struct {
	gameService *service.GameAnalyzerService
}

// NewHandler creates a new API handler
func NewHandler(gameService *service.GameAnalyzerService) *Handler {
	return &Handler{
		gameService: gameService,
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
