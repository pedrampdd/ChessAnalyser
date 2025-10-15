package models

import "time"

// Player represents a chess player
type Player struct {
	Username string `json:"username"`
	PlayerID *int   `json:"player_id,omitempty"`
	URL      string `json:"url,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Country  string `json:"country,omitempty"`
	Title    string `json:"title,omitempty"`
}

// GameMove represents a single move in a chess game
type GameMove struct {
	MoveNumber    int    `json:"move_number"`
	WhiteMove     string `json:"white_move,omitempty"`
	BlackMove     string `json:"black_move,omitempty"`
	FEN           string `json:"fen,omitempty"`
	TimeRemaining *int   `json:"time_remaining,omitempty"`
}

// GameInfo represents complete game information
type GameInfo struct {
	GameID      string     `json:"game_id"`
	URL         string     `json:"url"`
	FEN         string     `json:"fen"`
	PGN         string     `json:"pgn"`
	TimeControl string     `json:"time_control"`
	Rules       string     `json:"rules"`
	WhitePlayer Player     `json:"white_player"`
	BlackPlayer Player     `json:"black_player"`
	Result      string     `json:"result"`
	ResultCode  string     `json:"result_code"`
	TimeClass   string     `json:"time_class"`
	Rated       bool       `json:"rated"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Moves       []GameMove `json:"moves,omitempty"`
	Tournament  string     `json:"tournament,omitempty"`
	Match       string     `json:"match,omitempty"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GameResponse represents the response structure for game data
type GameResponse struct {
	GameID      string     `json:"game_id"`
	URL         string     `json:"url"`
	FEN         string     `json:"fen"`
	PGN         string     `json:"pgn"`
	TimeControl string     `json:"time_control"`
	Rules       string     `json:"rules"`
	WhitePlayer Player     `json:"white_player"`
	BlackPlayer Player     `json:"black_player"`
	Result      string     `json:"result"`
	ResultCode  string     `json:"result_code"`
	TimeClass   string     `json:"time_class"`
	Rated       bool       `json:"rated"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Tournament  string     `json:"tournament,omitempty"`
	Match       string     `json:"match,omitempty"`
}
