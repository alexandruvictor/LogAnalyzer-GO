package stats

import (
	"fmt"
	"log-analyzer/parser"
)

// Stats ține toate metricile
type Stats struct {
	TotalLines  int
	TotalErrors map[string]int
	LatencySum  int
	LatencyMax  int
}

// NewStats creează un obiect Stats
func NewStats() *Stats {
	return &Stats{
		TotalErrors: make(map[string]int),
	}
}

// AddEntry adaugă o linie în stats
func (s *Stats) AddEntry(entry *parser.LogEntry) {
	s.TotalLines++
	if entry.Level == "ERROR" {
		s.TotalErrors[fmt.Sprintf("%d %s", entry.Status, entry.Path)]++
	}

	s.LatencySum += entry.LatencyMs
	if entry.LatencyMs > s.LatencyMax {
		s.LatencyMax = entry.LatencyMs
	}
}

// PrintSummary afișează rezultatele
func (s *Stats) PrintSummary() {
	fmt.Printf("Total lines processed: %d\n", s.TotalLines)
	fmt.Println("\nErrors:")
	for k, v := range s.TotalErrors {
		fmt.Printf("%s -> %d\n", k, v)
	}
	if s.TotalLines > 0 {
		fmt.Printf("\nLatency:\navg: %dms\nmax: %dms\n", s.LatencySum/s.TotalLines, s.LatencyMax)
	}
}
