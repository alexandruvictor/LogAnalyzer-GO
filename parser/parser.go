package parser

import (
	"errors"  // used for returning parse errors
	"os"      // file reading
	"strconv" // string conversions
	"strings" // string manipulation utilities
	"time"    // timestamp parsing
)

// LogEntry represents a single parsed log line. Fields correspond to the
// expected log format:
//
//	<timestamp> <level> <status> <path> <latency>ms
//
// where timestamp is RFC3339, status is an HTTP status code, and latency is
// an integer measured in milliseconds.
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Status    int
	Path      string
	LatencyMs int
}

// ReadLines reads the entire file identified by filename and returns its
// contents split into lines. The caller is responsible for handling empty
// or malformed lines as necessary.
func ReadLines(filename string) ([]string, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Split on newline; trailing newline will produce a final empty string.
	lines := strings.Split(string(raw), "\n")
	return lines, nil
}

// ParseLine attempts to parse a single line of log text. If the line does not
// conform to the expected structure, an error is returned. On success a
// pointer to a populated LogEntry is returned.
func ParseLine(line string) (*LogEntry, error) {
	parts := strings.Fields(line)
	if len(parts) < 5 {
		return nil, errors.New("malformed line: insufficient fields")
	}

	// Parse timestamp (RFC3339 format).
	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, err
	}

	// Convert status code to integer.
	statusCode, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	// Latency value ends with "ms"; strip the suffix before conversion.
	latencyStr := strings.TrimSuffix(parts[4], "ms")
	latencyMs, err := strconv.Atoi(latencyStr)
	if err != nil {
		return nil, err
	}

	entry := &LogEntry{
		Timestamp: timestamp,
		Level:     parts[1],
		Status:    statusCode,
		Path:      parts[3],
		LatencyMs: latencyMs,
	}
	return entry, nil
}
