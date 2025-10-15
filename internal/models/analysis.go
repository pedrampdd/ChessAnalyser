package models

import "time"

// AnalysisResult represents the result of a chess position analysis
type AnalysisResult struct {
	Position           string   `json:"position"`    // FEN position
	MoveNumber         int      `json:"move_number"` // Move number in the game
	BestMove           string   `json:"best_move"`   // Best move found by engine
	Evaluation         float64  `json:"evaluation"`  // Centipawn evaluation
	Depth              int      `json:"depth"`       // Search depth reached
	Nodes              int64    `json:"nodes"`       // Number of nodes searched
	Time               int64    `json:"time"`        // Analysis time in milliseconds
	PrincipalVariation []string `json:"pv"`          // Principal variation (best line)
	MultiPV            int      `json:"multipv"`     // Multi-PV line number
}

// MoveAnalysis represents analysis for a specific move
type MoveAnalysis struct {
	Move         string            `json:"move"`         // Move in algebraic notation
	MoveNumber   int               `json:"move_number"`  // Move number
	Evaluation   float64           `json:"evaluation"`   // Position evaluation after move
	Accuracy     float64           `json:"accuracy"`     // Move accuracy percentage
	Blunder      bool              `json:"blunder"`      // True if move is a blunder
	Mistake      bool              `json:"mistake"`      // True if move is a mistake
	Inaccuracy   bool              `json:"inaccuracy"`   // True if move is an inaccuracy
	BestMove     string            `json:"best_move"`    // Best move in this position
	Alternatives []MoveAlternative `json:"alternatives"` // Alternative moves
}

// MoveAlternative represents an alternative move suggestion
type MoveAlternative struct {
	Move       string  `json:"move"`       // Alternative move
	Evaluation float64 `json:"evaluation"` // Evaluation of this move
	Depth      int     `json:"depth"`      // Search depth
}

// GameAnalysis represents complete analysis of a chess game
type GameAnalysis struct {
	GameID         string          `json:"game_id"`         // Original game ID
	PGN            string          `json:"pgn"`             // Original PGN
	AnalysisTime   time.Time       `json:"analysis_time"`   // When analysis was performed
	EngineVersion  string          `json:"engine_version"`  // Stockfish version used
	EngineSettings EngineSettings  `json:"engine_settings"` // Analysis settings
	Moves          []MoveAnalysis  `json:"moves"`           // Analysis for each move
	GameEvaluation float64         `json:"game_evaluation"` // Overall game evaluation
	Accuracy       GameAccuracy    `json:"accuracy"`        // Overall accuracy metrics
	Summary        AnalysisSummary `json:"summary"`         // Analysis summary
}

// EngineSettings represents Stockfish engine configuration
type EngineSettings struct {
	Depth      int `json:"depth"`       // Search depth
	TimeLimit  int `json:"time_limit"`  // Time limit in milliseconds
	MultiPV    int `json:"multipv"`     // Number of principal variations
	Threads    int `json:"threads"`     // Number of threads
	HashSize   int `json:"hash_size"`   // Hash table size in MB
	SkillLevel int `json:"skill_level"` // Skill level (0-20)
	Contempt   int `json:"contempt"`    // Contempt factor
}

// GameAccuracy represents accuracy metrics for the entire game
type GameAccuracy struct {
	WhiteAccuracy   float64 `json:"white_accuracy"`   // White player accuracy
	BlackAccuracy   float64 `json:"black_accuracy"`   // Black player accuracy
	AverageAccuracy float64 `json:"average_accuracy"` // Average accuracy
	Blunders        int     `json:"blunders"`         // Number of blunders
	Mistakes        int     `json:"mistakes"`         // Number of mistakes
	Inaccuracies    int     `json:"inaccuracies"`     // Number of inaccuracies
	BrilliantMoves  int     `json:"brilliant_moves"`  // Number of brilliant moves
	GreatMoves      int     `json:"great_moves"`      // Number of great moves
	BestMoves       int     `json:"best_moves"`       // Number of best moves
}

// AnalysisSummary provides a high-level summary of the analysis
type AnalysisSummary struct {
	TotalMoves      int      `json:"total_moves"`     // Total number of moves analyzed
	AnalysisDepth   int      `json:"analysis_depth"`  // Average analysis depth
	TotalTime       int64    `json:"total_time"`      // Total analysis time in ms
	NodesSearched   int64    `json:"nodes_searched"`  // Total nodes searched
	GamePhase       string   `json:"game_phase"`      // Opening/Middlegame/Endgame
	Complexity      string   `json:"complexity"`      // Low/Medium/High complexity
	Recommendations []string `json:"recommendations"` // Analysis recommendations
}

// AnalysisRequest represents a request for game analysis
type AnalysisRequest struct {
	GameID       string         `json:"game_id"`       // Game identifier
	PGN          string         `json:"pgn"`           // PGN to analyze
	Settings     EngineSettings `json:"settings"`      // Analysis settings
	IncludeMoves bool           `json:"include_moves"` // Include move-by-move analysis
	MaxMoves     int            `json:"max_moves"`     // Maximum moves to analyze (0 = all)
}

// AnalysisResponse represents the response for an analysis request
type AnalysisResponse struct {
	Success bool          `json:"success"`
	Data    *GameAnalysis `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
	Message string        `json:"message,omitempty"`
}
