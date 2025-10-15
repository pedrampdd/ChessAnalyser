package errors

import (
	"testing"
)

func TestGameNotFoundError(t *testing.T) {
	err := NewGameNotFoundError("test-game-id", nil)

	if err.GameID != "test-game-id" {
		t.Errorf("GameID = %v, want test-game-id", err.GameID)
	}

	expectedMsg := "game with ID test-game-id not found"
	if err.Error() != expectedMsg {
		t.Errorf("Error() = %v, want %v", err.Error(), expectedMsg)
	}
}

func TestAPIError(t *testing.T) {
	err := NewAPIError("test message", nil)

	if err.Message != "test message" {
		t.Errorf("Message = %v, want test message", err.Message)
	}

	expectedMsg := "API error: test message"
	if err.Error() != expectedMsg {
		t.Errorf("Error() = %v, want %v", err.Error(), expectedMsg)
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("field", "message")

	if err.Field != "field" {
		t.Errorf("Field = %v, want field", err.Field)
	}

	if err.Message != "message" {
		t.Errorf("Message = %v, want message", err.Message)
	}

	expectedMsg := "validation error for field field: message"
	if err.Error() != expectedMsg {
		t.Errorf("Error() = %v, want %v", err.Error(), expectedMsg)
	}
}
