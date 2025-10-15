# Chess Analyzer API

A comprehensive chess analysis API with Stockfish engine integration for PGN analysis and position evaluation.

## Quick Start

### Prerequisites
- Go 1.21+
- Stockfish engine installed

### Using Docker (Recommended)
```bash
# Build and run with Docker Compose
docker-compose up --build

# Or build and run manually
docker build -t chess-analyzer .
docker run -p 8080:8080 chess-analyzer
```

### Using Go directly
```bash
# Install Stockfish (Ubuntu/Debian)
sudo apt install stockfish

# Install Stockfish (macOS)
brew install stockfish

# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go

# Run tests
go test -v
```

## API Usage Examples

### Get Game Information
```bash
# Get a game by player/month format
curl "http://localhost:8080/api/game/hikaru/2024/01"

# Example response:
{
  "success": true,
  "data": {
    "game_id": "hikaru_1640995200",
    "url": "https://www.chess.com/game/live/123456789",
    "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
    "pgn": "1. e4 e5 2. Nf3 Nc6...",
    "time_control": "600+0",
    "rules": "chess",
    "white_player": {
      "username": "hikaru",
      "title": "GM",
      "country": "US"
    },
    "black_player": {
      "username": "magnus",
      "title": "GM", 
      "country": "NO"
    },
    "result": "1-0",
    "result_code": "win",
    "time_class": "blitz",
    "rated": true,
    "start_time": "2024-01-01T12:00:00Z",
    "end_time": "2024-01-01T12:15:00Z"
  }
}
```

### Get Player's Games
```bash
# Get player's games for January 2024
curl "http://localhost:8080/api/player/hikaru/games?year=2024&month=1"
```

### Get Player Profile
```bash
# Get player profile information
curl "http://localhost:8080/api/player/hikaru/profile"
```

### Get Player Statistics
```bash
# Get player's rating statistics
curl "http://localhost:8080/api/player/hikaru/stats"
```

## PGN Analysis Examples

### Analyze Chess Game
```bash
# Analyze a complete game
curl -X POST "http://localhost:8080/api/analyze/game" \
  -H "Content-Type: application/json" \
  -d '{
    "pgn": "[Event \"Test Game\"]\n[Site \"Chess.com\"]\n[Date \"2023.01.01\"]\n[Round \"1\"]\n[White \"Player1\"]\n[Black \"Player2\"]\n[Result \"1-0\"]\n\n1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 1-0",
    "settings": {
      "depth": 15,
      "time_limit": 5000
    },
    "max_moves": 10
  }'

# Example response:
{
  "success": true,
  "data": {
    "game_id": "12345",
    "pgn": "[Event \"Test Game\"]...",
    "analysis_time": "2023-01-01T12:00:00Z",
    "engine_version": "Stockfish 17.1",
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
      "total_moves": 20,
      "analysis_depth": 15,
      "total_time": 45000,
      "nodes_searched": 15000000,
      "game_phase": "opening",
      "complexity": "medium",
      "recommendations": [
        "Focus on tactical calculations to reduce blunders"
      ]
    }
  },
  "message": "Game analysis completed successfully"
}
```

### Analyze Chess Position
```bash
# Analyze a single position
curl "http://localhost:8080/api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&depth=15"

# Example response:
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
    "pv": ["e4", "e5", "Nf3", "Nc6", "Bb5"]
  }
}
```

### Engine Management
```bash
# Check engine status
curl "http://localhost:8080/api/analyze/status"

# Clear analysis cache
curl -X DELETE "http://localhost:8080/api/analyze/cache"
```

## Supported Game ID Formats

1. **Player/Month Format**: `username/YYYY/MM`
   - Example: `hikaru/2024/01`
   - Returns the first game from that player's monthly archive

2. **Direct URLs** (planned): `https://www.chess.com/game/live/123456789`

3. **Direct Game IDs** (planned): `123456789`

## Error Handling

The API returns appropriate HTTP status codes:
- `200` - Success
- `400` - Bad Request (invalid parameters)
- `404` - Not Found (game/player doesn't exist)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

## Rate Limiting

- Serial requests are unlimited
- Parallel requests may be rate limited by Chess.com
- Implement retry logic with exponential backoff for production use

## Development

### Project Structure
```
chessAnalyser/
├── cmd/
│   └── server/
│       └── main.go              # Main server application
├── internal/
│   ├── api/
│   │   ├── handlers.go          # API request handlers
│   │   └── routes.go            # Route definitions
│   ├── client/
│   │   └── chesscom.go         # Chess.com API client
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── engine/
│   │   └── stockfish.go        # Stockfish engine integration
│   ├── models/
│   │   ├── game.go             # Game data models
│   │   └── analysis.go         # Analysis data models
│   ├── parser/
│   │   └── pgn.go              # PGN parsing functionality
│   └── service/
│       ├── game.go             # Game service logic
│       ├── analysis.go         # Analysis service logic
│       ├── game_test.go        # Game service tests
│       └── analysis_test.go    # Analysis service tests
├── pkg/
│   └── errors/
│       ├── errors.go           # Custom error types
│       └── errors_test.go      # Error tests
├── docs/
│   ├── API_DOCUMENTATION.md    # Complete API documentation
│   ├── PGN_ANALYSIS.md         # PGN analysis guide
│   └── QUICK_START.md          # Quick start guide
├── stockfish/                  # Stockfish engine binary
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── Dockerfile                  # Docker configuration
├── docker-compose.yml         # Docker Compose configuration
├── README.md                  # Main documentation
└── USAGE.md                   # Usage examples
```

### Adding New Features

1. **New Endpoints**: Add routes in `setupRoutes()` function
2. **New Data Types**: Define structs in `main.go`
3. **Chess.com API Integration**: Add methods to `ChessComAPI` struct
4. **Caching**: Extend the `gameCache` map or implement persistent storage

### Testing

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run specific test
go test -v -run TestParseGameID
```

## Production Considerations

1. **Environment Variables**: Add configuration for API keys, database URLs, etc.
2. **Logging**: Implement structured logging with levels
3. **Metrics**: Add Prometheus metrics for monitoring
4. **Database**: Implement persistent storage for caching
5. **Authentication**: Add API key or JWT authentication
6. **Rate Limiting**: Implement client-side rate limiting
7. **Health Checks**: Add comprehensive health check endpoints
8. **Graceful Shutdown**: Implement proper shutdown handling
