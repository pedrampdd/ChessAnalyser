# Chess.com Golang Game Analyzer API

A Go-based REST API for analyzing chess games using the Chess.com Published Data API.

## Features

- Retrieve chess game information by game ID
- Get player profiles and statistics
- Fetch player's games by month
- Caching for improved performance
- RESTful API design with JSON responses

## API Endpoints

### Get Game by ID
```
GET /api/game/{gameId}
```

**Supported Game ID formats:**
- `username/YYYY/MM` - Player's games for a specific month (e.g., `hikaru/2024/01`)
- Direct Chess.com game URLs (planned feature)
- Direct game identifiers (planned feature)

**Response:**
```json
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
    "end_time": "2024-01-01T12:15:00Z",
    "tournament": "",
    "match": ""
  }
}
```

### Get Player's Games
```
GET /api/player/{username}/games?year=YYYY&month=MM
```

**Parameters:**
- `year` (required): Year (e.g., 2024)
- `month` (required): Month (1-12)

**Response:**
```json
{
  "success": true,
  "data": {
    "games": [
      {
        "url": "https://www.chess.com/game/live/123456789",
        "pgn": "1. e4 e5 2. Nf3 Nc6...",
        "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
        "time_control": "600+0",
        "rules": "chess",
        "white": {
          "username": "hikaru",
          "rating": 2800,
          "result": "win"
        },
        "black": {
          "username": "magnus",
          "rating": 2850,
          "result": "lose"
        },
        "time_class": "blitz",
        "rated": true,
        "start_time": 1640995200,
        "end_time": 1640996100
      }
    ]
  }
}
```

### Get Player Profile
```
GET /api/player/{username}/profile
```

**Response:**
```json
{
  "success": true,
  "data": {
    "username": "hikaru",
    "player_id": 123456,
    "title": "GM",
    "status": "premium",
    "name": "Hikaru Nakamura",
    "avatar": "https://images.chesscomfiles.com/uploads/v1/user/123456.1234567890.1234567890.1234567890.jpeg",
    "location": "United States",
    "country": "US",
    "joined": 1234567890,
    "last_online": 1640995200,
    "followers": 500000,
    "is_streamer": true,
    "twitch_url": "https://twitch.tv/gmhikaru",
    "url": "https://www.chess.com/member/hikaru"
  }
}
```

### Get Player Statistics
```
GET /api/player/{username}/stats
```

**Response:**
```json
{
  "success": true,
  "data": {
    "chess_daily": {
      "last": {
        "rating": 2800,
        "date": 1640995200
      },
      "best": {
        "rating": 2850,
        "date": 1640908800
      },
      "record": {
        "win": 1000,
        "loss": 200,
        "draw": 50
      }
    },
    "chess_rapid": {
      "last": {
        "rating": 2750,
        "date": 1640995200
      },
      "best": {
        "rating": 2800,
        "date": 1640908800
      },
      "record": {
        "win": 500,
        "loss": 100,
        "draw": 25
      }
    }
  }
}
```

## Installation and Setup

1. **Install Go** (version 1.21 or higher)

2. **Clone and setup the project:**
   ```bash
   git clone <repository-url>
   cd chessAnalyser
   go mod tidy
   ```

3. **Run the server:**
   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

## Usage Examples

### Using curl

```bash
# Get a game by ID (player/month format)
curl "http://localhost:8080/api/game/hikaru/2024/01"

# Get player's games for January 2024
curl "http://localhost:8080/api/player/hikaru/games?year=2024&month=1"

# Get player profile
curl "http://localhost:8080/api/player/hikaru/profile"

# Get player statistics
curl "http://localhost:8080/api/player/hikaru/stats"
```

### Using JavaScript/Fetch

```javascript
// Get game information
const response = await fetch('http://localhost:8080/api/game/hikaru/2024/01');
const data = await response.json();

if (data.success) {
  console.log('Game:', data.data);
  console.log('Players:', data.data.white_player.username, 'vs', data.data.black_player.username);
  console.log('Result:', data.data.result);
  console.log('PGN:', data.data.pgn);
} else {
  console.error('Error:', data.error);
}
```

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `200 OK` - Successful request
- `400 Bad Request` - Invalid parameters
- `404 Not Found` - Game or player not found
- `429 Too Many Requests` - Rate limit exceeded (from Chess.com API)
- `500 Internal Server Error` - Server error

Error response format:
```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

## Rate Limiting

The API respects Chess.com's rate limiting policies:
- Serial requests are unlimited
- Parallel requests may be rate limited
- Use appropriate User-Agent headers
- Implement retry logic with exponential backoff

## Data Caching

The API implements in-memory caching for:
- Game information by game ID
- Player profiles and statistics

Cache invalidation occurs when:
- New data is requested
- Cache reaches memory limits (planned feature)

## Future Enhancements

- [ ] Direct game URL parsing
- [ ] Game search by multiple criteria
- [ ] Move-by-move analysis
- [ ] Opening and endgame database integration
- [ ] Real-time game monitoring
- [ ] Database persistence
- [ ] Authentication and rate limiting
- [ ] WebSocket support for live games

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
