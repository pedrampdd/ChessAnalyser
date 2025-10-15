package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ChessComAPI represents the Chess.com API client
type ChessComAPI struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// NewChessComAPI creates a new Chess.com API client
func NewChessComAPI() *ChessComAPI {
	return &ChessComAPI{
		BaseURL: "https://api.chess.com/pub",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent: "ChessAnalyzer/1.0",
	}
}

// GetPlayerProfile retrieves player profile information
func (api *ChessComAPI) GetPlayerProfile(username string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/player/%s", api.BaseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPlayerGames retrieves player's games for a specific month
func (api *ChessComAPI) GetPlayerGames(username string, year, month int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/player/%s/games/%d/%02d", api.BaseURL, username, year, month)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPlayerStats retrieves player's statistics
func (api *ChessComAPI) GetPlayerStats(username string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/player/%s/stats", api.BaseURL, username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (api *ChessComAPI) GetGameByID(gameID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/game/live/%s", api.BaseURL, gameID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
