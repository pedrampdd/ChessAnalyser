package service_test

import (
	"context"
	"testing"
	"time"

	"chess-analyzer/internal/models"
	"chess-analyzer/internal/service"
)

func TestAnalysisService_AnalyzeGame(t *testing.T) {
	// Mock test - in real implementation, you'd need Stockfish binary
	// t.Skip("Skipping integration test - requires Stockfish binary")

	service, err := service.NewAnalysisService("../../stockfish/stockfish", 1, models.EngineSettings{
		Depth:     10,
		TimeLimit: 1000,
		Threads:   1,
		HashSize:  64,
	})
	if err != nil {
		t.Fatalf("Failed to create analysis service: %v", err)
	}
	defer service.Close()

	testPGN := `[Event "Test Game"]
[Site "Test Site"]
[Date "2023.01.01"]
[Round "1"]
[White "TestWhite"]
[Black "TestBlack"]
[Result "1-0"]

1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 4. Ba4 Nf6 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7 11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15. Nb1 h6 16. Bh4 c5 17. dxe5 Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21. Nc4 Nxc4 22. Bxc4 Nb6 23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7 27. Qe3 Qg5 28. Qxg5 hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33. f3 Bc8 34. Kf2 Bf5 35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5 40. Rd6 Kc5 41. Ra6 Nf2 42. g4 Bd3 43. Re6 1-0`

	request := &models.AnalysisRequest{
		PGN:          testPGN,
		Settings:     models.EngineSettings{Depth: 10, TimeLimit: 1000},
		IncludeMoves: true,
		MaxMoves:     10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	analysis, err := service.AnalyzeGame(ctx, request)
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	if analysis == nil {
		t.Fatal("Analysis result is nil")
	}

	if analysis.PGN != testPGN {
		t.Errorf("Expected PGN to match input, got: %s", analysis.PGN)
	}

	if len(analysis.Moves) == 0 {
		t.Error("Expected moves to be analyzed")
	}

	if analysis.Accuracy.AverageAccuracy == 0 {
		t.Error("Expected accuracy to be calculated")
	}
}

func TestAnalysisService_AnalyzePosition(t *testing.T) {
	t.Skip("Skipping integration test - requires Stockfish binary")

	service, err := service.NewAnalysisService("../../stockfish/stockfish", 1, models.EngineSettings{
		Depth:     10,
		TimeLimit: 1000,
		Threads:   1,
		HashSize:  64,
	})
	if err != nil {
		t.Fatalf("Failed to create analysis service: %v", err)
	}
	defer service.Close()

	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	settings := models.EngineSettings{Depth: 10, TimeLimit: 1000}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := service.AnalyzePosition(ctx, fen, settings)
	if err != nil {
		t.Fatalf("Position analysis failed: %v", err)
	}

	if result == nil {
		t.Fatal("Analysis result is nil")
	}

	if result.BestMove == "" {
		t.Error("Expected best move to be found")
	}

	if result.Depth == 0 {
		t.Error("Expected depth to be set")
	}
}

func TestAnalysisService_GetEngineStatus(t *testing.T) {
	t.Skip("Skipping integration test - requires Stockfish binary")

	service, err := service.NewAnalysisService("../../stockfish/stockfish", 2, models.EngineSettings{
		Depth:     10,
		TimeLimit: 1000,
		Threads:   1,
		HashSize:  64,
	})
	if err != nil {
		t.Fatalf("Failed to create analysis service: %v", err)
	}
	defer service.Close()

	status := service.GetEngineStatus()
	if status == nil {
		t.Fatal("Engine status is nil")
	}

	if totalEngines, ok := status["total_engines"].(int); !ok || totalEngines != 2 {
		t.Errorf("Expected total_engines to be 2, got: %v", status["total_engines"])
	}
}

func TestAnalysisService_ClearCache(t *testing.T) {
	service, err := service.NewAnalysisService("../../stockfish/stockfish", 1, models.EngineSettings{
		Depth:     10,
		TimeLimit: 1000,
		Threads:   1,
		HashSize:  64,
	})
	if err != nil {
		t.Fatalf("Failed to create analysis service: %v", err)
	}
	defer service.Close()

	// Clear cache should not panic
	service.ClearCache()

	status := service.GetEngineStatus()
	if cacheSize, ok := status["cache_size"].(int); !ok || cacheSize != 0 {
		t.Errorf("Expected cache_size to be 0 after clear, got: %v", status["cache_size"])
	}
}
