# Quick Start Guide - PGN Analysis

This guide will help you get started with the chess PGN analysis features using Stockfish engine.

## Prerequisites

1. **Go 1.21+** installed
2. **Stockfish engine** installed
3. **Chess Analyzer** project cloned

## Installation

### 1. Install Stockfish

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install stockfish
```

**macOS:**
```bash
brew install stockfish
```

**Manual Installation:**
```bash
# Download and extract Stockfish
wget https://github.com/official-stockfish/Stockfish/releases/latest/download/stockfish-ubuntu-x86-64-avx2.tar.gz
tar -xzf stockfish-ubuntu-x86-64-avx2.tar.gz
mv stockfish-ubuntu-x86-64-avx2 stockfish
```

### 2. Setup Project

```bash
git clone <your-repo>
cd chessAnalyser
go mod tidy
```

### 3. Configure Stockfish Path (if needed)

```bash
# If Stockfish is not in default location
export STOCKFISH_PATH="/path/to/your/stockfish"
```

### 4. Start Server

```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

## Quick Examples

### 1. Analyze a Simple Game

```bash
curl -X POST "http://localhost:8080/api/analyze/game" \
  -H "Content-Type: application/json" \
  -d '{
    "pgn": "[Event \"Quick Game\"]\n[Site \"Chess.com\"]\n[Date \"2023.01.01\"]\n[Round \"1\"]\n[White \"Player1\"]\n[Black \"Player2\"]\n[Result \"1-0\"]\n\n1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 1-0",
    "max_moves": 5
  }'
```

### 2. Analyze a Position

```bash
curl "http://localhost:8080/api/analyze/position?fen=rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR%20w%20KQkq%20-%200%201&depth=10"
```

### 3. Check Engine Status

```bash
curl "http://localhost:8080/api/analyze/status"
```

## JavaScript Examples

### Analyze Game with JavaScript

```javascript
async function analyzeGame(pgn) {
  const response = await fetch('http://localhost:8080/api/analyze/game', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      pgn: pgn,
      settings: {
        depth: 15,
        time_limit: 5000
      },
      max_moves: 20
    })
  });

  const data = await response.json();
  
  if (data.success) {
    console.log('Analysis completed!');
    console.log('Average accuracy:', data.data.accuracy.average_accuracy);
    console.log('Blunders:', data.data.accuracy.blunders);
    console.log('Mistakes:', data.data.accuracy.mistakes);
    
    // Display move analysis
    data.data.moves.forEach(move => {
      console.log(`Move ${move.move_number}: ${move.move} - Accuracy: ${move.accuracy.toFixed(1)}%`);
    });
  } else {
    console.error('Analysis failed:', data.error);
  }
}

// Example usage
const pgn = `[Event "Test Game"]
[Site "Chess.com"]
[Date "2023.01.01"]
[Round "1"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 1-0`;

analyzeGame(pgn);
```

### Analyze Position with JavaScript

```javascript
async function analyzePosition(fen) {
  const response = await fetch(`http://localhost:8080/api/analyze/position?fen=${encodeURIComponent(fen)}&depth=15`);
  const data = await response.json();
  
  if (data.success) {
    console.log('Best move:', data.data.best_move);
    console.log('Evaluation:', data.data.evaluation);
    console.log('Principal variation:', data.data.pv.join(' '));
  } else {
    console.error('Position analysis failed:', data.error);
  }
}

// Example usage
analyzePosition('rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1');
```

## Python Examples

### Analyze Game with Python

```python
import requests
import json

def analyze_game(pgn):
    url = "http://localhost:8080/api/analyze/game"
    data = {
        "pgn": pgn,
        "settings": {
            "depth": 15,
            "time_limit": 5000
        },
        "max_moves": 20
    }
    
    response = requests.post(url, json=data)
    result = response.json()
    
    if result["success"]:
        print("Analysis completed!")
        print(f"Average accuracy: {result['data']['accuracy']['average_accuracy']:.1f}%")
        print(f"Blunders: {result['data']['accuracy']['blunders']}")
        print(f"Mistakes: {result['data']['accuracy']['mistakes']}")
        
        # Display move analysis
        for move in result['data']['moves']:
            print(f"Move {move['move_number']}: {move['move']} - Accuracy: {move['accuracy']:.1f}%")
    else:
        print(f"Analysis failed: {result['error']}")

# Example usage
pgn = """[Event "Test Game"]
[Site "Chess.com"]
[Date "2023.01.01"]
[Round "1"]
[White "Player1"]
[Black "Player2"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 1-0"""

analyze_game(pgn)
```

### Analyze Position with Python

```python
import requests

def analyze_position(fen):
    url = f"http://localhost:8080/api/analyze/position?fen={fen}&depth=15"
    response = requests.get(url)
    result = response.json()
    
    if result["success"]:
        print(f"Best move: {result['data']['best_move']}")
        print(f"Evaluation: {result['data']['evaluation']}")
        print(f"Principal variation: {' '.join(result['data']['pv'])}")
    else:
        print(f"Position analysis failed: {result['error']}")

# Example usage
analyze_position("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
```

## Configuration Options

### Environment Variables

```bash
# Stockfish Configuration
export STOCKFISH_PATH="./stockfish/stockfish"
export STOCKFISH_MAX_ENGINES=4
export STOCKFISH_DEFAULT_DEPTH=15
export STOCKFISH_DEFAULT_TIME_LIMIT=5000

# Analysis Configuration
export ANALYSIS_MAX_CACHE_SIZE=1000
export ANALYSIS_ENABLE_CACHING=true
export ANALYSIS_CONCURRENT=true
```

### Analysis Settings

| Setting | Description | Default | Recommended |
|---------|-------------|---------|-------------|
| `depth` | Search depth | 15 | 15-20 for accuracy |
| `time_limit` | Time limit (ms) | 5000 | 3000-10000 |
| `threads` | CPU threads | 4 | Match CPU cores |
| `hash_size` | Hash table (MB) | 128 | 128-512 |
| `multipv` | Principal variations | 1 | 1-3 |

## Performance Tips

1. **For Speed**: Lower depth (10-12), shorter time limits (1000-3000ms)
2. **For Accuracy**: Higher depth (18-20), longer time limits (5000-10000ms)
3. **For Memory**: Smaller hash tables (64-128MB), fewer engines
4. **For Concurrency**: More engines (4-8), larger pool size

## Troubleshooting

### Common Issues

1. **"Stockfish not found"**
   - Check if Stockfish is installed: `which stockfish`
   - Set correct path: `export STOCKFISH_PATH="/path/to/stockfish"`

2. **"Analysis timeout"**
   - Increase time limits
   - Reduce analysis depth
   - Check system resources

3. **"Engine communication error"**
   - Verify Stockfish UCI compatibility
   - Check executable permissions
   - Restart the server

### Debug Commands

```bash
# Test Stockfish directly
echo "uci" | stockfish

# Check engine status
curl "http://localhost:8080/api/analyze/status"

# Clear cache if needed
curl -X DELETE "http://localhost:8080/api/analyze/cache"
```

## Next Steps

1. **Explore the API**: Try different analysis settings
2. **Build Applications**: Use the API in your chess applications
3. **Integrate with Databases**: Store analysis results
4. **Add Features**: Implement opening books, endgame tables
5. **Scale Up**: Use multiple servers for high-volume analysis

## Support

- Check the main README.md for complete documentation
- Review API_DOCUMENTATION.md for detailed endpoint information
- Check PGN_ANALYSIS.md for advanced usage patterns
- Open issues for bugs or feature requests
