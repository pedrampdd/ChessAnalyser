package service

import (
	"testing"
	"time"
)

func TestParseGameID(t *testing.T) {
	service := NewGameAnalyzerService()

	tests := []struct {
		name    string
		gameID  string
		wantErr bool
	}{
		{
			name:    "Valid player/month format",
			gameID:  "hikaru/2024/01",
			wantErr: false,
		},
		{
			name:    "Invalid format",
			gameID:  "invalid-game-id",
			wantErr: true,
		},
		{
			name:    "URL format (not implemented)",
			gameID:  "https://www.chess.com/game/live/123456789",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.parseGameID(tt.gameID)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGameID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseGameData(t *testing.T) {
	service := NewGameAnalyzerService()

	// Mock game data
	gameData := map[string]any{
		"url":          "https://www.chess.com/game/live/123456789",
		"fen":          "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"pgn":          "1. e4 e5 2. Nf3 Nc6",
		"time_control": "600+0",
		"rules":        "chess",
		"white": map[string]any{
			"username":  "hikaru",
			"player_id": float64(123456),
			"title":     "GM",
			"country":   "US",
		},
		"black": map[string]any{
			"username":  "magnus",
			"player_id": float64(789012),
			"title":     "GM",
			"country":   "NO",
		},
		"result":      "1-0",
		"result_code": "win",
		"time_class":  "blitz",
		"rated":       true,
		"start_time":  float64(1640995200),
		"end_time":    float64(1640996100),
	}

	gameInfo, err := service.parseGameData(gameData)
	if err != nil {
		t.Fatalf("parseGameData() error = %v", err)
	}

	// Test basic fields
	if gameInfo.URL != gameData["url"] {
		t.Errorf("URL = %v, want %v", gameInfo.URL, gameData["url"])
	}

	if gameInfo.FEN != gameData["fen"] {
		t.Errorf("FEN = %v, want %v", gameInfo.FEN, gameData["fen"])
	}

	if gameInfo.WhitePlayer.Username != "hikaru" {
		t.Errorf("WhitePlayer.Username = %v, want hikaru", gameInfo.WhitePlayer.Username)
	}

	if gameInfo.BlackPlayer.Username != "magnus" {
		t.Errorf("BlackPlayer.Username = %v, want magnus", gameInfo.BlackPlayer.Username)
	}

	if gameInfo.Rated != true {
		t.Errorf("Rated = %v, want true", gameInfo.Rated)
	}

	// Test timestamp parsing
	expectedStartTime := time.Unix(1640995200, 0)
	if !gameInfo.StartTime.Equal(expectedStartTime) {
		t.Errorf("StartTime = %v, want %v", gameInfo.StartTime, expectedStartTime)
	}
}

func TestHelperFunctions(t *testing.T) {
	data := map[string]interface{}{
		"string_val": "test",
		"float_val":  float64(123.45),
		"bool_val":   true,
		"int_val":    float64(42),
	}

	// Test getStringValue
	if got := getStringValue(data, "string_val"); got != "test" {
		t.Errorf("getStringValue() = %v, want test", got)
	}

	if got := getStringValue(data, "nonexistent"); got != "" {
		t.Errorf("getStringValue() = %v, want empty string", got)
	}

	// Test getFloatValue
	if got := getFloatValue(data, "float_val"); got != 123.45 {
		t.Errorf("getFloatValue() = %v, want 123.45", got)
	}

	if got := getFloatValue(data, "nonexistent"); got != 0 {
		t.Errorf("getFloatValue() = %v, want 0", got)
	}

	// Test getBoolValue
	if got := getBoolValue(data, "bool_val"); got != true {
		t.Errorf("getBoolValue() = %v, want true", got)
	}

	if got := getBoolValue(data, "nonexistent"); got != false {
		t.Errorf("getBoolValue() = %v, want false", got)
	}
}
