package engine

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"chess-analyzer/internal/models"
)

// StockfishEngine represents a Stockfish chess engine instance
type StockfishEngine struct {
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	stdout      io.ReadCloser
	stderr      io.ReadCloser
	scanner     *bufio.Scanner
	mu          sync.RWMutex
	isReady     bool
	isAnalyzing bool
	settings    models.EngineSettings
	version     string
}

// EnginePool manages multiple Stockfish engine instances
type EnginePool struct {
	Engines    []*StockfishEngine
	Available  chan *StockfishEngine
	mu         sync.RWMutex
	maxEngines int
	settings   models.EngineSettings
}

// NewStockfishEngine creates a new Stockfish engine instance
func NewStockfishEngine(executablePath string, settings models.EngineSettings) (*StockfishEngine, error) {
	cmd := exec.Command(executablePath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Stockfish: %w", err)
	}

	engine := &StockfishEngine{
		cmd:      cmd,
		stdin:    stdin,
		stdout:   stdout,
		stderr:   stderr,
		scanner:  bufio.NewScanner(stdout),
		settings: settings,
	}

	// Initialize the engine
	if err := engine.initialize(); err != nil {
		engine.Close()
		return nil, fmt.Errorf("failed to initialize engine: %w", err)
	}

	return engine, nil
}

// initialize sets up the engine with UCI protocol
func (e *StockfishEngine) initialize() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Send UCI command
	if err := e.sendCommand("uci"); err != nil {
		return err
	}

	// Wait for uciok
	if err := e.waitForResponse("uciok"); err != nil {
		return err
	}

	// Set engine options
	if err := e.configureEngine(); err != nil {
		return err
	}

	// Set ready
	if err := e.sendCommand("isready"); err != nil {
		return err
	}

	if err := e.waitForResponse("readyok"); err != nil {
		return err
	}

	e.isReady = true
	return nil
}

// configureEngine sets engine parameters
func (e *StockfishEngine) configureEngine() error {
	commands := []string{
		fmt.Sprintf("setoption name Threads value %d", e.settings.Threads),
		fmt.Sprintf("setoption name Hash value %d", e.settings.HashSize),
		fmt.Sprintf("setoption name Skill Level value %d", e.settings.SkillLevel),
		fmt.Sprintf("setoption name Contempt value %d", e.settings.Contempt),
	}

	for _, cmd := range commands {
		if err := e.sendCommand(cmd); err != nil {
			return err
		}
	}

	return nil
}

// sendCommand sends a command to the engine
func (e *StockfishEngine) sendCommand(command string) error {
	_, err := fmt.Fprintf(e.stdin, "%s\n", command)
	return err
}

// waitForResponse waits for a specific response from the engine
func (e *StockfishEngine) waitForResponse(expected string) error {
	timeout := time.After(10 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for response: %s", expected)
		default:
			if e.scanner.Scan() {
				line := strings.TrimSpace(e.scanner.Text())
				if strings.Contains(line, expected) {
					return nil
				}
			} else {
				return fmt.Errorf("scanner error while waiting for: %s", expected)
			}
		}
	}
}

// AnalyzePosition analyzes a chess position
func (e *StockfishEngine) AnalyzePosition(ctx context.Context, fen string, settings models.EngineSettings) (*models.AnalysisResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.isReady {
		return nil, fmt.Errorf("engine is not ready")
	}

	e.isAnalyzing = true
	defer func() { e.isAnalyzing = false }()

	// Set position
	if err := e.sendCommand(fmt.Sprintf("position fen %s", fen)); err != nil {
		return nil, err
	}

	// Start analysis
	analysisCmd := fmt.Sprintf("go depth %d", settings.Depth)
	if settings.TimeLimit > 0 {
		analysisCmd = fmt.Sprintf("go movetime %d", settings.TimeLimit)
	}
	if settings.MultiPV > 1 {
		analysisCmd += fmt.Sprintf(" multipv %d", settings.MultiPV)
	}

	if err := e.sendCommand(analysisCmd); err != nil {
		return nil, err
	}

	// Parse analysis results
	result, err := e.parseAnalysisOutput(ctx, settings.MultiPV)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// parseAnalysisOutput parses the engine's analysis output
func (e *StockfishEngine) parseAnalysisOutput(ctx context.Context, multiPV int) (*models.AnalysisResult, error) {
	var result models.AnalysisResult
	var pvLines []string

	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("analysis timeout")
		default:
			if e.scanner.Scan() {
				line := strings.TrimSpace(e.scanner.Text())

				if strings.HasPrefix(line, "bestmove") {
					// Analysis complete
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						result.BestMove = parts[1]
					}
					result.PrincipalVariation = pvLines
					return &result, nil
				}

				// Parse info lines
				if strings.HasPrefix(line, "info") {
					if err := e.parseInfoLine(line, &result, &pvLines); err != nil {
						continue // Continue parsing even if one line fails
					}
				}
			} else {
				return nil, fmt.Errorf("scanner error during analysis")
			}
		}
	}
}

// parseInfoLine parses a single info line from Stockfish
func (e *StockfishEngine) parseInfoLine(line string, result *models.AnalysisResult, pvLines *[]string) error {
	// Extract depth
	if depth := extractInt(line, "depth"); depth > 0 {
		result.Depth = depth
	}

	// Extract nodes
	if nodes := extractInt64(line, "nodes"); nodes > 0 {
		result.Nodes = nodes
	}

	// Extract time
	if time := extractInt64(line, "time"); time > 0 {
		result.Time = time
	}

	// Extract evaluation
	if eval := extractFloat(line, "score cp"); eval != 0 {
		result.Evaluation = eval / 100.0 // Convert centipawns to pawns
	} else if mate := extractInt(line, "score mate"); mate != 0 {
		// Handle mate scores
		if mate > 0 {
			result.Evaluation = 1000.0 - float64(mate)
		} else {
			result.Evaluation = -1000.0 - float64(mate)
		}
	}

	// Extract principal variation
	if strings.Contains(line, "pv") {
		pv := extractPV(line)
		if len(pv) > 0 {
			*pvLines = pv
		}
	}

	return nil
}

// extractInt extracts an integer value from a string
func extractInt(line, key string) int {
	re := regexp.MustCompile(fmt.Sprintf(`%s\s+(\d+)`, key))
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			return val
		}
	}
	return 0
}

// extractInt64 extracts an int64 value from a string
func extractInt64(line, key string) int64 {
	re := regexp.MustCompile(fmt.Sprintf(`%s\s+(\d+)`, key))
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			return val
		}
	}
	return 0
}

// extractFloat extracts a float value from a string
func extractFloat(line, key string) float64 {
	re := regexp.MustCompile(fmt.Sprintf(`%s\s+(-?\d+)`, key))
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return val
		}
	}
	return 0
}

// extractPV extracts the principal variation from a line
func extractPV(line string) []string {
	parts := strings.Fields(line)
	var pv []string
	inPV := false

	for _, part := range parts {
		if part == "pv" {
			inPV = true
			continue
		}
		if inPV {
			pv = append(pv, part)
		}
	}

	return pv
}

// GetVersion returns the engine version
func (e *StockfishEngine) GetVersion() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.version
}

// IsReady returns whether the engine is ready
func (e *StockfishEngine) IsReady() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.isReady
}

// IsAnalyzing returns whether the engine is currently analyzing
func (e *StockfishEngine) IsAnalyzing() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.isAnalyzing
}

// Close shuts down the engine
func (e *StockfishEngine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.stdin != nil {
		e.stdin.Close()
	}
	if e.stdout != nil {
		e.stdout.Close()
	}
	if e.stderr != nil {
		e.stderr.Close()
	}

	if e.cmd != nil && e.cmd.Process != nil {
		return e.cmd.Process.Kill()
	}

	return nil
}

// NewEnginePool creates a new engine pool
func NewEnginePool(maxEngines int, executablePath string, settings models.EngineSettings) (*EnginePool, error) {
	pool := &EnginePool{
		Engines:    make([]*StockfishEngine, 0, maxEngines),
		Available:  make(chan *StockfishEngine, maxEngines),
		maxEngines: maxEngines,
		settings:   settings,
	}

	// Create initial engines
	for i := 0; i < maxEngines; i++ {
		engine, err := NewStockfishEngine(executablePath, settings)
		if err != nil {
			// Clean up any created engines
			pool.Close()
			return nil, fmt.Errorf("failed to create engine %d: %w", i, err)
		}
		pool.Engines = append(pool.Engines, engine)
		pool.Available <- engine
	}

	return pool, nil
}

// GetEngine gets an available engine from the pool
func (p *EnginePool) GetEngine() *StockfishEngine {
	return <-p.Available
}

// ReturnEngine returns an engine to the pool
func (p *EnginePool) ReturnEngine(engine *StockfishEngine) {
	p.Available <- engine
}

// Close shuts down all Engines in the pool
func (p *EnginePool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var errs []error
	for _, engine := range p.Engines {
		if err := engine.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	close(p.Available)

	if len(errs) > 0 {
		return fmt.Errorf("errors closing Engines: %v", errs)
	}

	return nil
}
