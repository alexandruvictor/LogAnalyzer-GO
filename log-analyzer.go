// log-analyzer.go is the entry point for the Log Analyzer CLI application.
// It orchestrates reading the input file, parsing each log line, and
// delegating statistics aggregation to the stats package.  The code is kept
// simple for readability so that non‑technical reviewers (e.g. HR) can follow
// the overall flow.

package main

import (
	"flag"    // command-line flags
	"fmt"     // for user-facing messages
	"log"     // fallback logging
	"os"      // access to command-line arguments and file I/O
	"runtime" // to discover CPU count
	"sync"    // synchronization primitives

	"log-analyzer/parser" // custom package for parsing log lines
	"log-analyzer/stats"  // custom package for collecting statistics

	"github.com/sirupsen/logrus" // structured logging
)

var (
	jsonFlag = flag.Bool("json", false, "output results in JSON format to stdout")
	csvFile  = flag.String("csv", "", "path to write CSV report")

	errorLogger *logrus.Logger
)

// We choose the number of logical CPUs, but allow the caller to override
// via the GOMAXPROCS environment variable if desired.
func workerCount() int {
	return runtime.GOMAXPROCS(0)
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: log-analyzer [--json] [--csv=file] <logfile>")
		return
	}

	// collect all lines from input paths (files or directories)
	var allLines []string
	for _, path := range flag.Args() {
		info, err := os.Stat(path)
		if err != nil {
			log.Fatalf("invalid path %s: %v", path, err)
		}
		if info.IsDir() {
			entries, err := os.ReadDir(path)
			if err != nil {
				log.Fatalf("cannot read directory %s: %v", path, err)
			}
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				filePath := path + string(os.PathSeparator) + e.Name()
				lines, err := parser.ReadLines(filePath)
				if err != nil {
					logrus.WithField("file", filePath).Warnf("skipping unreadable file: %v", err)
					continue
				}
				allLines = append(allLines, lines...)
			}
		} else {
			lines, err := parser.ReadLines(path)
			if err != nil {
				log.Fatalf("failed to read log file %s: %v", path, err)
			}
			allLines = append(allLines, lines...)
		}
	}
	// use allLines for processing
	lines := allLines

	// Initialize statistics collector. The stats package exposes a simple
	// API for incrementing counters and calculating latencies.
	statsCollector := stats.NewStats()

	// setup structured logger for parse errors
	errorLogger = logrus.New()
	errLogFile, err := os.OpenFile("errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Fatalf("failed to open error log file: %v", err)
	}
	defer errLogFile.Close()
	errorLogger.Out = errLogFile
	errorLogger.SetLevel(logrus.ErrorLevel)
	errorLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

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
					// record failure in separate error log
					if errorLogger != nil {
						errorLogger.WithField("line", raw).WithError(err).Error("parse failed")
					}
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

	// optional exports
	if *jsonFlag {
		data, err := statsCollector.ToJSON()
		if err != nil {
			log.Fatalf("failed to generate JSON: %v", err)
		}
		fmt.Println(string(data))
	}

	if *csvFile != "" {
		f, err := os.Create(*csvFile)
		if err != nil {
			log.Fatalf("unable to create csv file: %v", err)
		}
		defer f.Close()
		if err := statsCollector.ToCSV(f); err != nil {
			log.Fatalf("csv export failed: %v", err)
		}
	}
}
