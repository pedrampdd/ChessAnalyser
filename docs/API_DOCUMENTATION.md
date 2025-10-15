# Chess Analyzer API Documentation

## Overview

The Chess Analyzer API provides comprehensive chess game analysis using the Stockfish engine. It offers both basic game retrieval from Chess.com and advanced PGN analysis capabilities.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, no authentication is required. All endpoints are publicly accessible.

## Response Format

All API responses follow a consistent format:

```json
{
  "success": boolean,
  "data": object | null,
  "error": string | null,
  "message": string | null
}
```

## Endpoints

### Game Retrieval Endpoints

#### Get Game by ID
- **URL:** `GET /api/game/{gameId}`
- **Description:** Retrieve game information by game ID
- **Parameters:**
  - `gameId` (path): Game identifier in format `username/YYYY/MM`

#### Get Player Games
- **URL:** `GET /api/player/{username}/games`
- **Description:** Get player's games for a specific month
- **Parameters:**
  - `username` (path): Player username
  - `year` (query): Year (required)
  - `month` (query): Month 1-12 (required)

#### Get Player Profile
- **URL:** `GET /api/player/{username}/profile`
- **Description:** Get player profile information
- **Parameters:**
  - `username` (path): Player username

#### Get Player Statistics
- **URL:** `GET /api/player/{username}/stats`
- **Description:** Get player's chess statistics
- **Parameters:**
  - `username` (path): Player username

### Analysis Endpoints

#### Analyze Chess Game
- **URL:** `POST /api/analyze/game`
- **Description:** Analyze a complete chess game using Stockfish engine
- **Content-Type:** `application/json`

**Request Body:**
```json
{
  "pgn": "string (required)",
  "settings": {
    "depth": "integer (default: 15)",
    "time_limit": "integer (default: 5000)",
    "threads": "integer (default: 4)",
    "hash_size": "integer (default: 128)",
    "multipv": "integer (default: 1)",
    "skill_level": "integer (default: 20)",
    "contempt": "integer (default: 0)"
  },
  "include_moves": "boolean (default: true)",
  "max_moves": "integer (default: 0 = all)"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "game_id": "string",
    "pgn": "string",
    "analysis_time": "ISO 8601 timestamp",
    "engine_version": "string",
    "engine_settings": {
      "depth": "integer",
      "time_limit": "integer",
      "threads": "integer",
      "hash_size": "integer",
      "multipv": "integer",
      "skill_level": "integer",
      "contempt": "integer"
    },
    "moves": [
      {
        "move": "string",
        "move_number": "integer",
        "evaluation": "float",
        "accuracy": "float",
        "blunder": "boolean",
        "mistake": "boolean",
        "inaccuracy": "boolean",
        "best_move": "string",
        "alternatives": [
          {
            "move": "string",
            "evaluation": "float",
            "depth": "integer"
          }
        ]
      }
    ],
    "accuracy": {
      "white_accuracy": "float",
      "black_accuracy": "float",
      "average_accuracy": "float",
      "blunders": "integer",
      "mistakes": "integer",
      "inaccuracies": "integer",
      "brilliant_moves": "integer",
      "great_moves": "integer",
      "best_moves": "integer"
    },
    "summary": {
      "total_moves": "integer",
      "analysis_depth": "integer",
      "total_time": "integer",
      "nodes_searched": "integer",
      "game_phase": "string",
      "complexity": "string",
      "recommendations": ["string"]
    }
  },
  "message": "string"
}
```

#### Analyze Chess Position
- **URL:** `GET /api/analyze/position`
- **Description:** Analyze a single chess position using Stockfish engine
- **Parameters:**
  - `fen` (query, required): FEN position string
  - `depth` (query, optional): Search depth (default: 15)
  - `time_limit` (query, optional): Time limit in milliseconds (default: 5000)
  - `threads` (query, optional): Number of threads (default: 4)
  - `hash_size` (query, optional): Hash table size in MB (default: 128)
  - `multipv` (query, optional): Number of principal variations (default: 1)

**Response:**
```json
{
  "success": true,
  "data": {
    "position": "string",
    "move_number": "integer",
    "best_move": "string",
    "evaluation": "float",
    "depth": "integer",
    "nodes": "integer",
    "time": "integer",
    "pv": ["string"]
  }
}
```

#### Get Engine Status
- **URL:** `GET /api/analyze/status`
- **Description:** Get the status of analysis engines in the pool

**Response:**
```json
{
  "success": true,
  "data": {
    "total_engines": "integer",
    "available_engines": "integer",
    "cache_size": "integer",
    "max_cache_size": "integer"
  }
}
```

#### Clear Analysis Cache
- **URL:** `DELETE /api/analyze/cache`
- **Description:** Clear the analysis cache to free memory

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "string"
  }
}
```

### Utility Endpoints

#### Health Check
- **URL:** `GET /health`
- **Description:** Check if the service is running

**Response:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "service": "chess-analyzer"
  }
}
```

## Error Codes

| HTTP Status | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid parameters |
| 404 | Not Found - Game or player not found |
| 500 | Internal Server Error - Server error |

## Rate Limiting

- No rate limiting is currently implemented
- Respects Chess.com API rate limits for data retrieval
- Analysis requests are processed concurrently using engine pooling

## Configuration

The API can be configured using environment variables:

### Server Configuration
- `SERVER_PORT`: Server port (default: 8080)
- `SERVER_HOST`: Server host (default: 0.0.0.0)

### Chess.com API Configuration
- `CHESS_API_BASE_URL`: Chess.com API base URL (default: https://api.chess.com/pub)
- `CHESS_API_USER_AGENT`: User agent string (default: ChessAnalyzer/1.0)
- `CHESS_API_TIMEOUT`: Request timeout in seconds (default: 30)

### Stockfish Configuration
- `STOCKFISH_PATH`: Path to Stockfish executable (default: ./stockfish/stockfish)
- `STOCKFISH_MAX_ENGINES`: Maximum number of engines in pool (default: 4)
- `STOCKFISH_DEFAULT_DEPTH`: Default search depth (default: 15)
- `STOCKFISH_DEFAULT_TIME_LIMIT`: Default time limit in milliseconds (default: 5000)
- `STOCKFISH_DEFAULT_THREADS`: Default number of threads (default: 4)
- `STOCKFISH_DEFAULT_HASH_SIZE`: Default hash table size in MB (default: 128)
- `STOCKFISH_DEFAULT_SKILL_LEVEL`: Default skill level (default: 20)
- `STOCKFISH_DEFAULT_CONTEMPT`: Default contempt factor (default: 0)

### Analysis Configuration
- `ANALYSIS_MAX_CACHE_SIZE`: Maximum cache size (default: 1000)
- `ANALYSIS_CACHE_EXPIRATION`: Cache expiration in minutes (default: 60)
- `ANALYSIS_MAX_MOVES_PER_GAME`: Maximum moves per game (default: 100)
- `ANALYSIS_ENABLE_CACHING`: Enable caching (default: true)
- `ANALYSIS_CONCURRENT`: Enable concurrent analysis (default: true)

## Examples

### Analyze a Game with Custom Settings

```bash
curl -X POST "http://localhost:8080/api/analyze/game" \
  -H "Content-Type: application/json" \
  -d '{
    "pgn": "[Event \"World Championship\"]\n[Site \"Chess.com\"]\n[Date \"2023.01.01\"]\n[Round \"1\"]\n[White \"Magnus Carlsen\"]\n[Black \"Hikaru Nakamura\"]\n[Result \"1-0\"]\n\n1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15. Nb1 h6 16. Bh4 c5 17. dxe5 Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21. Nc4 Nxc4 22. Bxc4 Nb6 23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7 27. Qe3 Qg5 28. Qxg5 hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33. f3 Bc8 34. Kf2 Bf5 35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5 40. Rd6 Kc5 41. Ra6 Nf2 42. g4 Bd3 43. Re6 1-0",
    "settings": {
      "depth": 20,
      "time_limit": 10000,
      "threads": 8,
      "hash_size": 512
    },
    "max_moves": 50
  }'
```

### Analyze Multiple Positions

```bash
# Analyze starting position
curl "http://localhost:8080/api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&depth=15"

# Analyze after 1.e4
curl "http://localhost:8080/api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR%20b%20KQkq%20e3%200%201&depth=15"

# Analyze after 1.e4 e5
curl "http://localhost:8080/api/analyze/position?fen=rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR%20w%20KQkq%20e5%200%202&depth=15"
```

## Performance Tips

1. **Engine Pool Size**: Adjust `STOCKFISH_MAX_ENGINES` based on your CPU cores
2. **Analysis Depth**: Higher depth = more accurate but slower analysis
3. **Time Limits**: Use time limits for consistent response times
4. **Caching**: Enable caching for repeated analyses
5. **Concurrent Requests**: The API supports multiple simultaneous analysis requests

## Troubleshooting

### Common Issues

1. **Stockfish not found**: Ensure Stockfish is installed and `STOCKFISH_PATH` is correct
2. **Analysis timeout**: Increase time limits or reduce analysis depth
3. **Memory issues**: Reduce hash table size or engine pool size
4. **Engine communication errors**: Check Stockfish UCI protocol compatibility

### Debug Information

Use the `/api/analyze/status` endpoint to monitor engine pool status and cache usage.
