package parser

import (
	"testing"
)

func TestPGNParser_ParsePGN(t *testing.T) {
	parser := NewPGNParser()

	testPGN := `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 1-0`

	game, err := parser.ParsePGN(testPGN)
	if err != nil {
		t.Fatalf("Failed to parse PGN: %v", err)
	}

	if game == nil {
		t.Fatal("Parsed game is nil")
	}

	// Test headers
	expectedHeaders := map[string]string{
		"event":  "Test Game",
		"site":   "Test Site",
		"date":   "2023.01.01",
		"round":  "1",
		"white":  "TestWhite",
		"black":  "TestBlack",
		"result": "1-0",
	}

	for key, expectedValue := range expectedHeaders {
		if actualValue, exists := game.Headers[key]; !exists || actualValue != expectedValue {
			t.Errorf("Header %s: expected %s, got %s", key, expectedValue, actualValue)
		}
	}

	// Test moves
	if len(game.Moves) == 0 {
		t.Error("Expected moves to be parsed")
	}

	// Test first move
	firstMove := game.Moves[0]
	if firstMove.MoveNumber != 1 || firstMove.Move != "e4" || firstMove.Color != "white" {
		t.Errorf("First move parsing failed: %+v", firstMove)
	}

	// Test result
	if game.Result != "1-0" {
		t.Errorf("Expected result 1-0, got %s", game.Result)
	}

	// Test game phase
	if game.GamePhase != "opening" {
		t.Errorf("Expected game phase 'opening', got %s", game.GamePhase)
	}
}

func TestPGNParser_ValidatePGN(t *testing.T) {
	parser := NewPGNParser()

	validPGN := `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]

1. e4 e5 1-0`

	err := parser.ValidatePGN(validPGN)
	if err != nil {
		t.Errorf("Valid PGN should not return error: %v", err)
	}

	// Test empty PGN
	err = parser.ValidatePGN("")
	if err == nil {
		t.Error("Empty PGN should return error")
	}

	// Test missing headers
	invalidPGN := `[Event "Test Game"]
1. e4 e5 1-0`
	err = parser.ValidatePGN(invalidPGN)
	if err == nil {
		t.Error("PGN with missing headers should return error")
	}

	// Test missing moves
	invalidPGN2 := `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]`
	err = parser.ValidatePGN(invalidPGN2)
	if err == nil {
		t.Error("PGN with missing moves should return error")
	}
}

func TestPGNParser_IsValidMove(t *testing.T) {
	parser := NewPGNParser()

	validMoves := []string{
		"e4", "Nf3", "O-O", "O-O-O", "Qxd5", "Bxf7+", "Nxe4", "Qh5#",
	}

	invalidMoves := []string{
		"invalid", "e9", "Nf10", "O-O-O-O", "x", "",
	}

	for _, move := range validMoves {
		if !parser.IsValidMove(move) {
			t.Errorf("Move %s should be valid", move)
		}
	}

	for _, move := range invalidMoves {
		if parser.IsValidMove(move) {
			t.Errorf("Move %s should be invalid", move)
		}
	}
}

func TestPGNParser_GetMoveAtPosition(t *testing.T) {
	parser := NewPGNParser()

	testPGN := `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 1-0`

	game, err := parser.ParsePGN(testPGN)
	if err != nil {
		t.Fatalf("Failed to parse PGN: %v", err)
	}

	// Test white move
	whiteMove, err := parser.GetMoveAtPosition(game, 1, "white")
	if err != nil {
		t.Errorf("Failed to get white move: %v", err)
	}
	if whiteMove.Move != "e4" {
		t.Errorf("Expected white move 'e4', got %s", whiteMove.Move)
	}

	// Test black move
	blackMove, err := parser.GetMoveAtPosition(game, 1, "black")
	if err != nil {
		t.Errorf("Failed to get black move: %v", err)
	}
	if blackMove.Move != "e5" {
		t.Errorf("Expected black move 'e5', got %s", blackMove.Move)
	}

	// Test non-existent move
	_, err = parser.GetMoveAtPosition(game, 10, "white")
	if err == nil {
		t.Error("Expected error for non-existent move")
	}
}

func TestPGNParser_GetGameLength(t *testing.T) {
	parser := NewPGNParser()

	testPGN := `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 1-0`

	game, err := parser.ParsePGN(testPGN)
	if err != nil {
		t.Fatalf("Failed to parse PGN: %v", err)
	}

	length := parser.GetGameLength(game)
	expectedLength := 6 // 3 moves each for white and black
	if length != expectedLength {
		t.Errorf("Expected game length %d, got %d", expectedLength, length)
	}
}

func TestPGNParser_ConvertToGameInfo(t *testing.T) {
	parser := NewPGNParser()

	testPGN := `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]

1. e4 e5 1-0`

	parsedGame, err := parser.ParsePGN(testPGN)
	if err != nil {
		t.Fatalf("Failed to parse PGN: %v", err)
	}

	gameInfo := parser.ConvertToGameInfo(parsedGame)
	if gameInfo == nil {
		t.Fatal("Converted GameInfo is nil")
	}

	if gameInfo.PGN != testPGN {
		t.Errorf("Expected PGN to match, got: %s", gameInfo.PGN)
	}

	if gameInfo.WhitePlayer.Username != "TestWhite" {
		t.Errorf("Expected white player 'TestWhite', got %s", gameInfo.WhitePlayer.Username)
	}

	if gameInfo.BlackPlayer.Username != "TestBlack" {
		t.Errorf("Expected black player 'TestBlack', got %s", gameInfo.BlackPlayer.Username)
	}

	if len(gameInfo.Moves) == 0 {
		t.Error("Expected moves to be converted")
	}
}
