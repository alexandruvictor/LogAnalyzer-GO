package stats

import (
	"strings"
	"testing"
	"time"

	"log-analyzer/parser"
)

func makeEntry(level string, status int, path string, latency int) *parser.LogEntry {
	return &parser.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Status:    status,
		Path:      path,
		LatencyMs: latency,
	}
}

func TestStats_AddEntry(t *testing.T) {
	s := NewStats()

	s.AddEntry(makeEntry("INFO", 200, "/ping", 10))
	s.AddEntry(makeEntry("ERROR", 500, "/api", 300))
	s.AddEntry(makeEntry("ERROR", 500, "/api", 100))

	if s.TotalLines != 3 {
		t.Errorf("expected 3 total lines, got %d", s.TotalLines)
	}

	if len(s.TotalErrors) != 1 {
		t.Errorf("expected 1 error key, got %d", len(s.TotalErrors))
	}

	key := "500 /api"
	if s.TotalErrors[key] != 2 {
		t.Errorf("expected 2 errors for %s, got %d", key, s.TotalErrors[key])
	}

	if s.LatencySum != 410 {
		t.Errorf("expected latency sum 410, got %d", s.LatencySum)
	}

	if s.LatencyMax != 300 {
		t.Errorf("expected max latency 300, got %d", s.LatencyMax)
	}
}

func BenchmarkStats_AddEntry(b *testing.B) {
	s := NewStats()
	entry := makeEntry("ERROR", 500, "/api", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.AddEntry(entry)
	}
}
func TestStats_ToJSON(t *testing.T) {
	s := NewStats()
	s.AddEntry(makeEntry("ERROR", 500, "/api", 100))
	data, err := s.ToJSON()
	if err != nil {
		t.Fatalf("unexpected JSON error: %v", err)
	}
	if !strings.Contains(string(data), "\"total_lines\": 1") {
		t.Error("JSON output missing total_lines")
	}
}

func TestStats_ToCSV(t *testing.T) {
	s := NewStats()
	s.AddEntry(makeEntry("ERROR", 500, "/api", 100))
	var buf strings.Builder
	if err := s.ToCSV(&buf); err != nil {
		t.Fatalf("CSV error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "total_lines,1") {
		t.Error("CSV output missing total_lines")
	}
	if !strings.Contains(out, "error,500 /api:1") {
		t.Error("CSV output missing error line")
	}
}
