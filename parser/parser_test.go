package parser

import (
	"os"
	"testing"
	"time"
)

// sample valid line used across tests
const validLine = "2026-02-01T10:15:01Z INFO 200 /api/login 120ms"

func TestParseLine_Valid(t *testing.T) {
	entry, err := ParseLine(validLine)
	if err != nil {
		t.Fatalf("expected valid entry, got error: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("expected level INFO, got %s", entry.Level)
	}
	if entry.Status != 200 {
		t.Errorf("expected status 200, got %d", entry.Status)
	}
	if entry.Path != "/api/login" {
		t.Errorf("expected path /api/login, got %s", entry.Path)
	}
	if entry.LatencyMs != 120 {
		t.Errorf("expected latency 120, got %d", entry.LatencyMs)
	}
	// check timestamp parsed correctly
	wantTime := time.Date(2026, 02, 01, 10, 15, 1, 0, time.UTC)
	if !entry.Timestamp.Equal(wantTime) {
		t.Errorf("unexpected timestamp: %v", entry.Timestamp)
	}
}

func TestParseLine_Malformed(t *testing.T) {
	tests := []string{
		"",
		"just some text",
		"2026-02-01T10:15:01Z INFO /api/login 120ms",     // missing status
		"2026-02-01T10:15:01Z INFO 200 /api/login notms", // bad latency
	}
	for _, line := range tests {
		if _, err := ParseLine(line); err == nil {
			t.Errorf("expected error for malformed line %q", line)
		}
	}
}

func TestReadLines(t *testing.T) {
	data := "line1\nline2\n"
	// simulate by writing to a temp file
	f, err := os.CreateTemp("", "logtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(data); err != nil {
		t.Fatal(err)
	}
	f.Close()

	lines, err := ReadLines(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 { // trailing newline yields empty
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

// bench for ParseLine with many lines
func BenchmarkParseLine(b *testing.B) {
	// generate a slice containing 10000 copies of validLine
	lines := make([]string, 10000)
	for i := range lines {
		lines[i] = validLine
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, l := range lines {
			if _, err := ParseLine(l); err != nil {
				b.Fatalf("unexpected parse error: %v", err)
			}
		}
	}
}
