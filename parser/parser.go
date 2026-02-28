package parser

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// LogEntry reprezintă o linie de log
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Status    int
	Path      string
	LatencyMs int
}

// ReadLines citește fișierul linie cu linie
func ReadLines(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

// ParseLine parsează o linie și returnează un LogEntry
func ParseLine(line string) (*LogEntry, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return nil, errors.New("malformed line")
	}

	ts, err := time.Parse(time.RFC3339, fields[0])
	if err != nil {
		return nil, err
	}

	status, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, err
	}

	latencyStr := strings.TrimSuffix(fields[4], "ms")
	latency, err := strconv.Atoi(latencyStr)
	if err != nil {
		return nil, err
	}

	entry := &LogEntry{
		Timestamp: ts,
		Level:     fields[1],
		Status:    status,
		Path:      fields[3],
		LatencyMs: latency,
	}
	return entry, nil
}
