# Chess PGN Analysis Package

A high-performance chess game analysis package using the Stockfish engine for your chess analyzer application.

## Features

- **PGN Parsing**: Complete PGN (Portable Game Notation) parsing with move extraction
- **Stockfish Integration**: Full UCI protocol support with engine pooling
- **Position Analysis**: Deep position analysis with multiple variations
- **Game Analysis**: Complete game analysis with accuracy metrics
- **Caching**: Intelligent caching system for performance optimization
- **Concurrent Processing**: Multi-engine support for parallel analysis
- **REST API**: Complete REST API integration

## Architecture

### Core Components

1. **Models** (`internal/models/analysis.go`)
   - `AnalysisResult`: Single position analysis result
   - `MoveAnalysis`: Individual move analysis with accuracy metrics
   - `GameAnalysis`: Complete game analysis with statistics
   - `EngineSettings`: Stockfish engine configuration

2. **Engine** (`internal/engine/stockfish.go`)
   - `StockfishEngine`: UCI protocol communication
   - `EnginePool`: Multi-engine management for concurrent analysis

3. **Parser** (`internal/parser/pgn.go`)
   - `PGNParser`: PGN parsing and validation
   - `ParsedGame`: Structured game representation

4. **Service** (`internal/service/analysis.go`)
   - `AnalysisService`: Main analysis orchestration
   - Caching and performance optimization

5. **API** (`internal/api/`)
   - REST endpoints for analysis operations
   - Integration with existing handlers

## Installation & Setup

### Prerequisites

1. **Stockfish Engine**: Download and install Stockfish
   ```bash
   # Ubuntu/Debian
   sudo apt-get install stockfish
   
   # macOS
   brew install stockfish
   
   # Or download from https://stockfishchess.org/download/
   ```

2. **Go Dependencies**: Already included in your project

### Configuration

Configure via environment variables:

```bash
# Stockfish Configuration
export STOCKFISH_PATH="/usr/bin/stockfish"
export STOCKFISH_MAX_ENGINES=4
export STOCKFISH_DEFAULT_DEPTH=15
export STOCKFISH_DEFAULT_TIME_LIMIT=5000
export STOCKFISH_DEFAULT_THREADS=4
export STOCKFISH_DEFAULT_HASH_SIZE=128
export STOCKFISH_DEFAULT_SKILL_LEVEL=20
export STOCKFISH_DEFAULT_CONTEMPT=0

# Analysis Configuration
export ANALYSIS_MAX_CACHE_SIZE=1000
export ANALYSIS_CACHE_EXPIRATION=60
export ANALYSIS_MAX_MOVES_PER_GAME=100
export ANALYSIS_ENABLE_CACHING=true
export ANALYSIS_CONCURRENT=true
```

## API Endpoints

### Game Analysis

**POST** `/api/analyze/game`

Analyze a complete chess game:

```json
{
  "pgn": "[Event \"Test Game\"]\n[Site \"Test Site\"]\n...",
  "settings": {
    "depth": 15,
    "time_limit": 5000,
    "threads": 4,
    "hash_size": 128,
    "multipv": 1
  },
  "include_moves": true,
  "max_moves": 50
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "game_id": "12345",
    "pgn": "...",
    "analysis_time": "2023-01-01T12:00:00Z",
    "engine_version": "Stockfish 15",
    "moves": [
      {
        "move": "e4",
        "move_number": 1,
        "evaluation": 0.2,
        "accuracy": 95.5,
        "blunder": false,
        "mistake": false,
        "inaccuracy": false,
        "best_move": "e4",
        "alternatives": []
      }
    ],
    "accuracy": {
      "white_accuracy": 92.3,
      "black_accuracy": 89.7,
      "average_accuracy": 91.0,
      "blunders": 2,
      "mistakes": 5,
      "inaccuracies": 8
    },
    "summary": {
      "total_moves": 40,
      "analysis_depth": 15,
      "total_time": 45000,
      "nodes_searched": 15000000,
      "game_phase": "middlegame",
      "complexity": "medium",
      "recommendations": [
        "Focus on tactical calculations to reduce blunders"
      ]
    }
  }
}
```

### Position Analysis

**GET** `/api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&depth=15&time_limit=5000`

Analyze a single chess position:

**Response:**
```json
{
  "success": true,
  "data": {
    "position": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
    "move_number": 0,
    "best_move": "e4",
    "evaluation": 0.2,
    "depth": 15,
    "nodes": 1500000,
    "time": 5000,
    "pv": ["e4", "e5", "Nf3", "Nc6"]
  }
}
```

### Engine Status

**GET** `/api/analyze/status`

Get engine pool status:

**Response:**
```json
{
  "success": true,
  "data": {
    "total_engines": 4,
    "available_engines": 3,
    "cache_size": 150,
    "max_cache_size": 1000
  }
}
```

### Cache Management

**DELETE** `/api/analyze/cache`

Clear analysis cache:

**Response:**
```json
{
  "success": true,
  "message": "Analysis cache cleared successfully"
}
```

## Usage Examples

### Basic Game Analysis

```go
package main

import (
    "context"
    "log"
    
    "github.com/pedrampdd/ChessAnalyser/internal/models"
    "github.com/pedrampdd/ChessAnalyser/internal/service"
)

func main() {
    // Initialize analysis service
    analysisService, err := service.NewAnalysisService(
        "/usr/bin/stockfish",
        4, // max engines
        models.EngineSettings{
            Depth:     15,
            TimeLimit: 5000,
            Threads:   4,
            HashSize:  128,
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    defer analysisService.Close()

    // Analyze a game
    request := &models.AnalysisRequest{
        PGN: `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 1-0`,
        Settings: models.EngineSettings{
            Depth:     15,
            TimeLimit: 5000,
        },
        IncludeMoves: true,
        MaxMoves:     20,
    }

    ctx := context.Background()
    analysis, err := analysisService.AnalyzeGame(ctx, request)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Analysis completed: %+v", analysis.Accuracy)
}
```

### Position Analysis

```go
// Analyze a single position
fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
settings := models.EngineSettings{
    Depth:     15,
    TimeLimit: 5000,
}

result, err := analysisService.AnalyzePosition(ctx, fen, settings)
if err != nil {
    log.Fatal(err)
}

log.Printf("Best move: %s, Evaluation: %.2f", result.BestMove, result.Evaluation)
```

## Performance Optimization

### Engine Pooling
- Multiple Stockfish engines for concurrent analysis
- Automatic engine management and load balancing
- Configurable pool size based on system resources

### Caching
- Intelligent caching of analysis results
- Configurable cache size and expiration
- Cache key based on PGN and analysis settings

### Concurrent Processing
- Parallel analysis of multiple positions
- Non-blocking API endpoints
- Context-based cancellation support

## Best Practices

### 1. Resource Management
- Set appropriate engine pool size based on CPU cores
- Monitor memory usage with hash table settings
- Use time limits to prevent long-running analyses

### 2. Analysis Settings
- **Depth**: Higher depth = more accurate but slower (15-20 recommended)
- **Time Limit**: Balance between speed and accuracy (3-10 seconds)
- **Threads**: Match your CPU core count
- **Hash Size**: 128-512 MB depending on available RAM

### 3. Error Handling
- Always check for analysis errors
- Implement timeout handling for long analyses
- Use context cancellation for user-initiated stops

### 4. Caching Strategy
- Enable caching for repeated analyses
- Set appropriate cache expiration times
- Monitor cache hit rates for optimization

## Testing

Run the test suite:

```bash
# Unit tests
go test ./internal/parser/
go test ./internal/service/

# Integration tests (requires Stockfish)
go test -tags=integration ./internal/service/
```

## Troubleshooting

### Common Issues

1. **Stockfish not found**
   - Ensure Stockfish is installed and in PATH
   - Check `STOCKFISH_PATH` environment variable

2. **Analysis timeout**
   - Increase time limits for complex positions
   - Reduce analysis depth for faster results

3. **Memory issues**
   - Reduce hash table size
   - Decrease engine pool size
   - Monitor system memory usage

4. **Engine communication errors**
   - Check Stockfish UCI protocol compatibility
   - Verify engine executable permissions

### Performance Tuning

1. **For Speed**: Lower depth, shorter time limits
2. **For Accuracy**: Higher depth, longer time limits
3. **For Memory**: Smaller hash tables, fewer engines
4. **For Concurrency**: More engines, larger pool size

## Contributing

When extending the analysis package:

1. Follow the existing architecture patterns
2. Add comprehensive tests for new features
3. Update documentation for API changes
4. Consider performance implications
5. Maintain backward compatibility

## License

This package is part of your chess analyzer application and follows the same license terms.
