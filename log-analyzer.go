// log-analyzer.go is the entry point for the Log Analyzer CLI application.
// It orchestrates reading the input file, parsing each log line, and
// delegating statistics aggregation to the stats package.  The code is kept
// simple for readability so that non‑technical reviewers (e.g. HR) can follow
// the overall flow.

package main

import (
	"fmt" // for user-facing messages
	"log" // for fatal errors that terminate execution
	"os"  // access to command-line arguments

	"log-analyzer/parser" // custom package for parsing log lines
	"log-analyzer/stats"  // custom package for collecting statistics
)

func main() {
	// Expect exactly one argument: the path to the log file.
	// If missing, show a usage message and exit gracefully.
	if len(os.Args) < 2 {
		fmt.Println("Usage: log-analyzer <logfile>")
		return
	}

	// Path of the log file provided by the user.
	logFilePath := os.Args[1]

	// Read all lines from the file. A failure here is a fatal error;
	// the program cannot continue without input.
	lines, err := parser.ReadLines(logFilePath)
	if err != nil {
		log.Fatalf("failed to read log file %s: %v", logFilePath, err)
	}

	// Initialize statistics collector. The stats package exposes a simple
	// API for incrementing counters and calculating latencies.
	statsCollector := stats.NewStats()

	// Process each line from the file. Malformed lines are ignored but do not
	// stop the program.
	for _, rawLine := range lines {
		entry, err := parser.ParseLine(rawLine)
		if err != nil {
			// Skip lines that cannot be parsed. These are treated as noise.
			continue
		}
		statsCollector.AddEntry(entry)
	}

	// After processing all lines, print a summary to stdout. The output
	// format is handled by the stats package.
	statsCollector.PrintSummary()
}
