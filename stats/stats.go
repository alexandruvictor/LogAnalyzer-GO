package stats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

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
// serialized to JSON or CSV if the user requests it (see README).
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

		// simple ASCII bar chart scaled down by factor 10
		scale := 10
		avgBar := strings.Repeat("=", avgLatency/scale)
		maxBar := strings.Repeat("=", s.LatencyMax/scale)
		fmt.Println("\nLatency graph (each '=' is", scale, "ms):")
		fmt.Printf("avg [%s]\n", avgBar)
		fmt.Printf("max [%s]\n", maxBar)
	}
}

// ToJSON returns the statistics serialized as JSON. The structure matches the
// example shown in the README.
func (s *Stats) ToJSON() ([]byte, error) {
	type out struct {
		TotalLines int            `json:"total_lines"`
		Errors     map[string]int `json:"errors"`
		Latency    struct {
			Avg int `json:"avg"`
			Max int `json:"max"`
		} `json:"latency"`
	}
	result := out{
		TotalLines: s.TotalLines,
		Errors:     s.TotalErrors,
	}
	if s.TotalLines > 0 {
		result.Latency.Avg = s.LatencySum / s.TotalLines
		result.Latency.Max = s.LatencyMax
	}
	return json.MarshalIndent(result, "", "  ")
}

// ToCSV writes the statistics in CSV format to the provided writer. The
// output begins with headers and includes a row per error key followed by
// latency lines.
func (s *Stats) ToCSV(w io.Writer) error {
	writer := csv.NewWriter(w)
	// headers
	if err := writer.Write([]string{"metric", "value"}); err != nil {
		return err
	}
	// total lines
	if err := writer.Write([]string{"total_lines", fmt.Sprint(s.TotalLines)}); err != nil {
		return err
	}
	// errors
	for k, v := range s.TotalErrors {
		if err := writer.Write([]string{"error", fmt.Sprintf("%s:%d", k, v)}); err != nil {
			return err
		}
	}
	if s.TotalLines > 0 {
		if err := writer.Write([]string{"latency_avg", fmt.Sprint(s.LatencySum / s.TotalLines)}); err != nil {
			return err
		}
		if err := writer.Write([]string{"latency_max", fmt.Sprint(s.LatencyMax)}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
