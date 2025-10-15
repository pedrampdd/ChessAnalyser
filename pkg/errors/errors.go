package errors

import "fmt"

// GameNotFoundError represents an error when a game is not found
type GameNotFoundError struct {
	GameID string
	Err    error
}

func (e *GameNotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("game with ID %s not found: %v", e.GameID, e.Err)
	}
	return fmt.Sprintf("game with ID %s not found", e.GameID)
}

func (e *GameNotFoundError) Unwrap() error {
	return e.Err
}

// APIError represents an error with the Chess.com API
type APIError struct {
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("API error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("API error: %s", e.Message)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for field %s: %s", e.Field, e.Message)
}

// NewGameNotFoundError creates a new GameNotFoundError
func NewGameNotFoundError(gameID string, err error) *GameNotFoundError {
	return &GameNotFoundError{
		GameID: gameID,
		Err:    err,
	}
}

// NewAPIError creates a new APIError
func NewAPIError(message string, err error) *APIError {
	return &APIError{
		Message: message,
		Err:     err,
	}
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
