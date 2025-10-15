# Chess Analyzer API

## Quick Start

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
# Install dependencies
go mod tidy

# Run the server
go run main.go

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
├── main.go              # Main application and API handlers
├── main_test.go         # Unit tests
├── go.mod               # Go module definition
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
└── README.md            # Detailed documentation
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
