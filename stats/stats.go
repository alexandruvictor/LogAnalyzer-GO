package stats

import (
	"fmt"

	"log-analyzer/parser" // used to interpret LogEntry values
)

// Stats aggregates metrics produced by scanning a sequence of log entries.
// It is intentionally simple: totals, error counts keyed by status+path, and
// latency statistics.
//
// All fields are exported so that future packages (e.g. a JSON printer) can
// access them directly.
type Stats struct {
	TotalLines  int            // number of lines (entries) processed
	TotalErrors map[string]int // "<status> <path>" -> count
	LatencySum  int            // sum of all latencies (ms)
	LatencyMax  int            // highest single latency seen (ms)
}

// NewStats returns an initialized Stats instance ready for use.
func NewStats() *Stats {
	return &Stats{
		TotalErrors: make(map[string]int),
	}
}

// AddEntry incorporates one parsed log entry into the running totals.
//
// If the entry's level is "ERROR" we increment the corresponding counter.
// Latency values are added to the sum and compared with the previous max.
func (s *Stats) AddEntry(entry *parser.LogEntry) {
	s.TotalLines++

	// Count error occurrences by status code and path.  Using a formatted key
	// keeps the representation compact and easy to print.
	if entry.Level == "ERROR" {
		key := fmt.Sprintf("%d %s", entry.Status, entry.Path)
		s.TotalErrors[key]++
	}

	s.LatencySum += entry.LatencyMs
	if entry.LatencyMs > s.LatencyMax {
		s.LatencyMax = entry.LatencyMs
	}
}

// PrintSummary writes a human-readable report to stdout.  This is the default
// display used by the CLI.  In practice the same Stats object can also be
// serialized to JSON if the user requests it (see README).
func (s *Stats) PrintSummary() {
	fmt.Printf("Total lines processed: %d\n", s.TotalLines)

	if len(s.TotalErrors) > 0 {
		fmt.Println("\nErrors:")
		for k, v := range s.TotalErrors {
			fmt.Printf("%s -> %d\n", k, v)
		}
	}

	if s.TotalLines > 0 {
		avgLatency := s.LatencySum / s.TotalLines
		fmt.Printf("\nLatency:\navg: %dms\nmax: %dms\n", avgLatency, s.LatencyMax)
	}
}
