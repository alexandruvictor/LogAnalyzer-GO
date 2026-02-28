// log-analyzer.go is the entry point for the Log Analyzer CLI application.
// It orchestrates reading the input file, parsing each log line, and
// delegating statistics aggregation to the stats package.  The code is kept
// simple for readability so that non‑technical reviewers (e.g. HR) can follow
// the overall flow.

package main

import (
	"fmt"     // for user-facing messages
	"log"     // for fatal errors that terminate execution
	"os"      // access to command-line arguments
	"runtime" // to discover CPU count
	"sync"    // synchronization primitives

	"log-analyzer/parser" // custom package for parsing log lines
	"log-analyzer/stats"  // custom package for collecting statistics
)

// workerCount returns a suitable number of goroutines to run in parallel.
// We choose the number of logical CPUs, but allow the caller to override
// via the GOMAXPROCS environment variable if desired.
func workerCount() int {
	return runtime.GOMAXPROCS(0)
}

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

	// Channel where parsed entries are sent for aggregation.
	entryChan := make(chan *parser.LogEntry, 1000)
	var wg sync.WaitGroup

	// Launch a goroutine that consumes entries and updates stats.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for e := range entryChan {
			statsCollector.AddEntry(e)
		}
	}()

	// Spawn a pool of workers to parse lines concurrently.
	workers := workerCount()
	var parseWG sync.WaitGroup
	parseWG.Add(workers)

	lineChan := make(chan string, 1000)

	// worker function
	for i := 0; i < workers; i++ {
		go func() {
			defer parseWG.Done()
			for raw := range lineChan {
				entry, err := parser.ParseLine(raw)
				if err != nil {
					continue // malformed, ignore
				}
				entryChan <- entry
			}
		}()
	}

	// Feed lines into the line channel
	go func() {
		for _, raw := range lines {
			lineChan <- raw
		}
		close(lineChan)
	}()

	// wait for parsers to finish, then close entry channel so stats goroutine
	// can exit
	parseWG.Wait()
	close(entryChan)

	// wait for stats aggregation to complete
	wg.Wait()

	// After processing all lines, print a summary to stdout. The output
	// format is handled by the stats package.
	statsCollector.PrintSummary()
}
