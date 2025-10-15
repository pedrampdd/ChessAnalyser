# Chess.com API Golang Client & Analyzer

A comprehensive Go-based REST API client and analyzer for Chess.com with Stockfish engine integration for advanced PGN analysis and position evaluation.

## Features

- **Chess.com API Client** - Complete Go client for Chess.com Published Data API
- **Game Retrieval** - Get chess games, player profiles, and statistics
- **PGN Analysis** - Advanced game analysis using Stockfish engine
- **Position Analysis** - Real-time position evaluation and move suggestions
- **Accuracy Metrics** - Blunder detection, mistake analysis, and accuracy scoring
- **Multi-engine Support** - Concurrent analysis with engine pooling
- **RESTful API** - Clean JSON API with comprehensive documentation
- **Caching** - Intelligent caching for improved performance
- **Docker Support** - Easy deployment with Docker and Docker Compose

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

## PGN Analysis Endpoints

### Analyze Chess Game
```
POST /api/analyze/game
```

Analyze a complete chess game using Stockfish engine.

**Request Body:**
```json
{
  "pgn": "[Event \"Test Game\"]\n[Site \"Test Site\"]\n[Date \"2023.01.01\"]\n[Round \"1\"]\n[White \"TestWhite\"]\n[Black \"TestBlack\"]\n[Result \"1-0\"]\n\n1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 1-0",
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

**Parameters:**
- `pgn` (required): PGN string to analyze
- `settings` (optional): Engine analysis settings
  - `depth`: Search depth (default: 15)
  - `time_limit`: Time limit in milliseconds (default: 5000)
  - `threads`: Number of threads (default: 4)
  - `hash_size`: Hash table size in MB (default: 128)
  - `multipv`: Number of principal variations (default: 1)
- `include_moves` (optional): Include move-by-move analysis (default: true)
- `max_moves` (optional): Maximum moves to analyze (default: 0 = all)

**Response:**
```json
{
  "success": true,
  "data": {
    "game_id": "12345",
    "pgn": "[Event \"Test Game\"]...",
    "analysis_time": "2023-01-01T12:00:00Z",
    "engine_version": "Stockfish 17.1",
    "engine_settings": {
      "depth": 15,
      "time_limit": 5000,
      "threads": 4,
      "hash_size": 128,
      "multipv": 1
    },
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
      },
      {
        "move": "e5",
        "move_number": 1,
        "evaluation": 0.1,
        "accuracy": 98.2,
        "blunder": false,
        "mistake": false,
        "inaccuracy": false,
        "best_move": "e5",
        "alternatives": []
      }
    ],
    "accuracy": {
      "white_accuracy": 92.3,
      "black_accuracy": 89.7,
      "average_accuracy": 91.0,
      "blunders": 2,
      "mistakes": 5,
      "inaccuracies": 8,
      "brilliant_moves": 1,
      "great_moves": 3,
      "best_moves": 15
    },
    "summary": {
      "total_moves": 40,
      "analysis_depth": 15,
      "total_time": 45000,
      "nodes_searched": 15000000,
      "game_phase": "middlegame",
      "complexity": "medium",
      "recommendations": [
        "Focus on tactical calculations to reduce blunders",
        "Study opening theory to improve early game play"
      ]
    }
  },
  "message": "Game analysis completed successfully"
}
```

### Analyze Chess Position
```
GET /api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&depth=15&time_limit=5000
```

Analyze a single chess position using Stockfish engine.

**Query Parameters:**
- `fen` (required): FEN position string
- `depth` (optional): Search depth (default: 15)
- `time_limit` (optional): Time limit in milliseconds (default: 5000)
- `threads` (optional): Number of threads (default: 4)
- `hash_size` (optional): Hash table size in MB (default: 128)
- `multipv` (optional): Number of principal variations (default: 1)

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
    "pv": ["e4", "e5", "Nf3", "Nc6", "Bb5"]
  }
}
```

### Get Engine Status
```
GET /api/analyze/status
```

Get the status of analysis engines in the pool.

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

### Clear Analysis Cache
```
DELETE /api/analyze/cache
```

Clear the analysis cache to free memory.

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "Analysis cache cleared successfully"
  }
}
```

## Installation and Setup

1. **Install Go** (version 1.21 or higher)

2. **Install Stockfish Engine:**
   
   **Ubuntu/Debian:**
   ```bash
   sudo apt update
   sudo apt install stockfish
   ```
   
   **macOS:**
   ```bash
   brew install stockfish
   ```
   
   **Windows:**
   - Download from https://stockfishchess.org/download/
   - Extract and add to PATH
   
   **Manual Installation:**
   ```bash
   # Download latest release
   wget https://github.com/official-stockfish/Stockfish/releases/latest/download/stockfish-ubuntu-x86-64-avx2.tar.gz
   tar -xzf stockfish-ubuntu-x86-64-avx2.tar.gz
   mv stockfish-ubuntu-x86-64-avx2 stockfish
   ```

3. **Clone and setup the project:**
   ```bash
   git clone <repository-url>
   cd chessAnalyser
   go mod tidy
   ```

4. **Configure Stockfish (Optional):**
   ```bash
   # Set custom Stockfish path if needed
   export STOCKFISH_PATH="/path/to/stockfish"
   
   # Configure analysis settings
   export STOCKFISH_MAX_ENGINES=4
   export STOCKFISH_DEFAULT_DEPTH=15
   export STOCKFISH_DEFAULT_TIME_LIMIT=5000
   ```

5. **Run the server:**
   ```bash
   go run cmd/server/main.go
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

# Analyze a chess game
curl -X POST "http://localhost:8080/api/analyze/game" \
  -H "Content-Type: application/json" \
  -d '{
    "pgn": "[Event \"Test Game\"]\n[Site \"Test Site\"]\n[Date \"2023.01.01\"]\n[Round \"1\"]\n[White \"TestWhite\"]\n[Black \"TestBlack\"]\n[Result \"1-0\"]\n\n1. e4 e5 2. Nf3 Nc6 1-0",
    "settings": {
      "depth": 15,
      "time_limit": 5000
    },
    "max_moves": 10
  }'

# Analyze a chess position
curl "http://localhost:8080/api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&depth=15"

# Get engine status
curl "http://localhost:8080/api/analyze/status"

# Clear analysis cache
curl -X DELETE "http://localhost:8080/api/analyze/cache"
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

// Analyze a chess game
const analysisResponse = await fetch('http://localhost:8080/api/analyze/game', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    pgn: '[Event "Test Game"]\n[Site "Test Site"]\n[Date "2023.01.01"]\n[Round "1"]\n[White "TestWhite"]\n[Black "TestBlack"]\n[Result "1-0"]\n\n1. e4 e5 2. Nf3 Nc6 1-0',
    settings: {
      depth: 15,
      time_limit: 5000
    },
    max_moves: 10
  })
});

const analysisData = await analysisResponse.json();
if (analysisData.success) {
  console.log('Analysis completed!');
  console.log('Average accuracy:', analysisData.data.accuracy.average_accuracy);
  console.log('Blunders:', analysisData.data.accuracy.blunders);
  console.log('Mistakes:', analysisData.data.accuracy.mistakes);
  console.log('Recommendations:', analysisData.data.summary.recommendations);
}

// Analyze a chess position
const positionResponse = await fetch('http://localhost:8080/api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&depth=15');
const positionData = await positionResponse.json();
if (positionData.success) {
  console.log('Best move:', positionData.data.best_move);
  console.log('Evaluation:', positionData.data.evaluation);
  console.log('Principal variation:', positionData.data.pv);
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

- [x] **PGN Analysis with Stockfish Engine** ✅
- [x] **Position Analysis and Move Evaluation** ✅
- [x] **Game Accuracy Metrics and Statistics** ✅
- [x] **Multi-engine Concurrent Analysis** ✅
- [ ] Redis Cache
- [ ] Direct game URL parsing
- [ ] Game search by multiple criteria
- [ ] Opening and endgame database integration
- [ ] Real-time game monitoring
- [ ] Database persistence for analysis results
- [ ] Authentication and rate limiting
- [ ] WebSocket support for live games
- [ ] Advanced FEN position calculation
- [ ] Tournament analysis features
- [ ] Opening book integration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.

## Keywords

`chess.com api golang` `chess api client` `golang chess` `chess.com golang` `chess api go` `stockfish golang` `chess analysis go` `chess.com client` `golang chess engine` `chess pgn analysis` `chess position analysis` `chess accuracy metrics` `chess blunder detection` `chess.com published data api` `golang chess library` `chess game analysis` `chess engine integration` `chess api wrapper` `golang chess client` `chess.com go client`
