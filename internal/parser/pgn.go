package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pedrampdd/ChessAnalyser/internal/models"
)

// PGNParser handles parsing of PGN (Portable Game Notation) files
type PGNParser struct {
	gameRegex *regexp.Regexp
	moveRegex *regexp.Regexp
}

// ParsedGame represents a parsed chess game from PGN
type ParsedGame struct {
	Headers   map[string]string `json:"headers"`
	Moves     []ParsedMove      `json:"moves"`
	Result    string            `json:"result"`
	PGN       string            `json:"pgn"`
	MoveCount int               `json:"move_count"`
	GamePhase string            `json:"game_phase"`
}

// ParsedMove represents a single move in a parsed game
type ParsedMove struct {
	MoveNumber int    `json:"move_number"`
	Move       string `json:"move"`
	Color      string `json:"color"` // "white" or "black"
	FEN        string `json:"fen"`
	Comment    string `json:"comment,omitempty"`
	NAG        string `json:"nag,omitempty"` // Numeric Annotation Glyph
}

// NewPGNParser creates a new PGN parser
func NewPGNParser() *PGNParser {
	return &PGNParser{
		gameRegex: regexp.MustCompile(`\[([A-Za-z]+)\s+"([^"]*)"\]`),
		moveRegex: regexp.MustCompile(`(\d+)\.\s*([^\s]+)\s+([^\s]+)?`),
	}
}

// ParsePGN parses a PGN string and returns a ParsedGame
func (p *PGNParser) ParsePGN(pgn string) (*ParsedGame, error) {
	if strings.TrimSpace(pgn) == "" {
		return nil, fmt.Errorf("empty PGN string")
	}

	// Split PGN into headers and moves
	parts := strings.Split(pgn, "\n\n")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid PGN format: missing moves section")
	}

	headers := p.parseHeaders(parts[0])
	moves, result, err := p.parseMoves(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse moves: %w", err)
	}

	game := &ParsedGame{
		Headers:   headers,
		Moves:     moves,
		Result:    result,
		PGN:       pgn,
		MoveCount: len(moves),
		GamePhase: p.determineGamePhase(len(moves)),
	}

	return game, nil
}

// parseHeaders extracts headers from the PGN header section
func (p *PGNParser) parseHeaders(headerSection string) map[string]string {
	headers := make(map[string]string)
	matches := p.gameRegex.FindAllStringSubmatch(headerSection, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			key := strings.ToLower(match[1])
			value := match[2]
			headers[key] = value
		}
	}

	return headers
}

// parseMoves extracts moves from the moves section
func (p *PGNParser) parseMoves(movesSection string) ([]ParsedMove, string, error) {
	var moves []ParsedMove
	var result string

	// Clean up the moves section
	movesSection = strings.TrimSpace(movesSection)

	// Extract result at the end
	if strings.HasSuffix(movesSection, " 1-0") || strings.HasSuffix(movesSection, " 0-1") ||
		strings.HasSuffix(movesSection, " 1/2-1/2") || strings.HasSuffix(movesSection, " *") {
		parts := strings.Fields(movesSection)
		if len(parts) > 0 {
			result = parts[len(parts)-1]
			movesSection = strings.TrimSuffix(movesSection, " "+result)
		}
	}

	// Parse individual moves
	lines := strings.Split(movesSection, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse moves in this line
		lineMoves, err := p.parseMoveLine(line)
		if err != nil {
			continue // Skip invalid lines
		}
		moves = append(moves, lineMoves...)
	}

	return moves, result, nil
}

// parseMoveLine parses a line containing chess moves
func (p *PGNParser) parseMoveLine(line string) ([]ParsedMove, error) {
	var moves []ParsedMove

	// Remove comments and annotations
	line = p.removeComments(line)

	// Split by move numbers
	parts := strings.Fields(line)
	var currentMoveNumber int
	var moveIndex int // Track moves within the current move number

	for _, part := range parts {
		// Check if this is a move number
		if strings.HasSuffix(part, ".") {
			if num, err := strconv.Atoi(strings.TrimSuffix(part, ".")); err == nil {
				currentMoveNumber = num
				moveIndex = 0 // Reset move index for new move number
			}
			continue
		}

		// Skip result indicators
		if part == "1-0" || part == "0-1" || part == "1/2-1/2" || part == "*" {
			continue
		}

		// This should be a move
		if currentMoveNumber > 0 {
			move := ParsedMove{
				MoveNumber: currentMoveNumber,
				Move:       part,
				Color:      p.determineMoveColor(currentMoveNumber, moveIndex),
			}
			moves = append(moves, move)
			moveIndex++
		}
	}

	return moves, nil
}

// removeComments removes comments and annotations from move text
func (p *PGNParser) removeComments(text string) string {
	// Remove {comments}
	commentRegex := regexp.MustCompile(`\{[^}]*\}`)
	text = commentRegex.ReplaceAllString(text, "")

	// Remove ;comments
	semicolonIndex := strings.Index(text, ";")
	if semicolonIndex != -1 {
		text = text[:semicolonIndex]
	}

	// Remove NAGs (Numeric Annotation Glyphs)
	nagRegex := regexp.MustCompile(`\$\d+`)
	text = nagRegex.ReplaceAllString(text, "")

	return strings.TrimSpace(text)
}

// determineMoveColor determines if a move is white or black
func (p *PGNParser) determineMoveColor(moveNumber, position int) string {
	// White moves are at even positions (0, 2, 4...)
	// Black moves are at odd positions (1, 3, 5...)
	if position%2 == 0 {
		return "white"
	}
	return "black"
}

// determineGamePhase determines the phase of the game based on move count
func (p *PGNParser) determineGamePhase(moveCount int) string {
	if moveCount <= 20 {
		return "opening"
	} else if moveCount <= 40 {
		return "middlegame"
	} else {
		return "endgame"
	}
}

// ExtractPositions extracts FEN positions for each move
func (p *PGNParser) ExtractPositions(game *ParsedGame) error {
	// For now, generate basic FEN positions
	// In a real implementation, you'd use a chess library to generate proper FEN strings
	for i := range game.Moves {
		// Generate a simple FEN based on move number
		// This is a placeholder - real implementation would parse moves and update position
		game.Moves[i].FEN = fmt.Sprintf("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - %d %d", i, (i/2)+1)
	}
	return nil
}

// ConvertToGameInfo converts a ParsedGame to GameInfo
func (p *PGNParser) ConvertToGameInfo(parsedGame *ParsedGame) *models.GameInfo {
	gameInfo := &models.GameInfo{
		PGN:    parsedGame.PGN,
		Result: parsedGame.Result,
		Moves:  make([]models.GameMove, len(parsedGame.Moves)),
	}

	// Convert headers
	if event, ok := parsedGame.Headers["event"]; ok {
		gameInfo.Tournament = event
	}

	if site, ok := parsedGame.Headers["site"]; ok {
		gameInfo.URL = site
	}

	if date, ok := parsedGame.Headers["date"]; ok {
		if t, err := time.Parse("2006.01.02", date); err == nil {
			gameInfo.StartTime = t
		}
	}

	if timeControl, ok := parsedGame.Headers["timecontrol"]; ok {
		gameInfo.TimeControl = timeControl
	}

	if rules, ok := parsedGame.Headers["rules"]; ok {
		gameInfo.Rules = rules
	}

	// Convert players
	if white, ok := parsedGame.Headers["white"]; ok {
		gameInfo.WhitePlayer = models.Player{Username: white}
	}

	if black, ok := parsedGame.Headers["black"]; ok {
		gameInfo.BlackPlayer = models.Player{Username: black}
	}

	// Convert moves
	for i, move := range parsedGame.Moves {
		gameMove := models.GameMove{
			MoveNumber: move.MoveNumber,
			FEN:        move.FEN,
		}

		if move.Color == "white" {
			gameMove.WhiteMove = move.Move
		} else {
			gameMove.BlackMove = move.Move
		}

		gameInfo.Moves[i] = gameMove
	}

	return gameInfo
}

// ValidatePGN validates if a PGN string is well-formed
func (p *PGNParser) ValidatePGN(pgn string) error {
	if strings.TrimSpace(pgn) == "" {
		return fmt.Errorf("empty PGN")
	}

	// Check for required headers
	headers := p.parseHeaders(strings.Split(pgn, "\n\n")[0])
	requiredHeaders := []string{"event", "site", "date", "round", "white", "black", "result"}

	for _, header := range requiredHeaders {
		if _, exists := headers[header]; !exists {
			return fmt.Errorf("missing required header: %s", header)
		}
	}

	// Check for moves section
	parts := strings.Split(pgn, "\n\n")
	if len(parts) < 2 {
		return fmt.Errorf("missing moves section")
	}

	// Basic move validation
	movesSection := parts[1]
	if strings.TrimSpace(movesSection) == "" {
		return fmt.Errorf("empty moves section")
	}

	return nil
}

// GetMoveAtPosition returns the move at a specific position number
func (p *PGNParser) GetMoveAtPosition(game *ParsedGame, moveNumber int, color string) (*ParsedMove, error) {
	for _, move := range game.Moves {
		if move.MoveNumber == moveNumber && move.Color == color {
			return &move, nil
		}
	}
	return nil, fmt.Errorf("move not found at position %d for %s", moveNumber, color)
}

// GetGameLength returns the total number of moves in the game
func (p *PGNParser) GetGameLength(game *ParsedGame) int {
	return len(game.Moves)
}

// IsValidMove checks if a move string is valid algebraic notation
func (p *PGNParser) IsValidMove(move string) bool {
	// Basic validation - this could be enhanced with more sophisticated checks
	moveRegex := regexp.MustCompile(`^[KQRBN]?[a-h]?[1-8]?x?[a-h][1-8](?:=[QRBN])?[+#]?$|^O-O(-O)?[+#]?$`)
	return moveRegex.MatchString(move)
}
