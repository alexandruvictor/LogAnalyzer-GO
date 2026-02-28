package main

import (
	"fmt"
	"log"
	"log-analyzer/parser"
	"log-analyzer/stats"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: log-analyzer <logfile>")
		return
	}

	logFile := os.Args[1]
	lines, err := parser.ReadLines(logFile)
	if err != nil {
		log.Fatalf("Failed to read log file: %v", err)
	}

	s := stats.NewStats()
	for _, line := range lines {
		entry, err := parser.ParseLine(line)
		if err != nil {
			continue // skip malformed lines
		}
		s.AddEntry(entry)
	}

	s.PrintSummary()
}
