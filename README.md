# Log Analyzer (Go)

![Go](https://img.shields.io/badge/Language-Go-blue?logo=go&logoColor=white)
![Status](https://img.shields.io/badge/status-production%20ready-brightgreen)

## Overview

**Log Analyzer** is a simple, high-performance CLI tool written in Go. It reads structured log files, gathers error counts and latency statistics, and displays results in human-readable or JSON format. The project is organized with modular packages for parsing and stats, making it easy to extend or reuse in other Go applications.

## Features

- Parse logs in the format: timestamp level status path latency
- Compute:
  - Total lines processed (ignoring malformed entries)
  - Error counts by status code and endpoint
  - Average and maximum latency
- Optional JSON output (--json)
- Clean separation between parser and stats packages
- Lightweight and fast; suitable for large files
- Easily extended for streaming, concurrency, or new metrics

## Project Structure

```text
log-analyzer/
├── go.mod             # Go module definition
├── log-analyzer.go    # application entry point
├── parser/            # parsing utilities
│   └── parser.go
├── stats/             # statistics aggregation
│   └── stats.go
├── logs/              # example log files
│   └── access.log
└── README.md
```

## Requirements

- Go 1.22 or newer
- Compatible with Windows, macOS, and Linux
- Terminal/CLI access

## Getting Started

1. Clone the repository:
   ```sh
   git clone https://your-repo-url.git
   cd log-analyzer
   ```

2. Build the application:
   ```sh
   go build -o log-analyzer log-analyzer.go
   ```

3. Run against a log file:
   ```sh
   ./log-analyzer logs/access.log
   ```

### Sample Log (`logs/access.log`)

```
2026-02-01T10:15:01Z INFO 200 /api/login 120ms
2026-02-01T10:15:01Z ERROR 500 /api/login 450ms
2026-02-01T10:15:02Z INFO 200 /api/products 80ms
2026-02-01T10:15:03Z ERROR 504 /api/products 900ms
2026-02-01T10:15:05Z WARN 301 /api/redirect 30ms
2026-02-01T10:15:06Z INFO 200 /api/cart 50ms
2026-02-01T10:15:07Z ERROR 403 /api/cart 300ms
2026-02-01T10:15:08Z INFO 200 /api/checkout 100ms
2026-02-01T10:15:09Z ERROR 500 /api/checkout 700ms
MALFORMED LINE HERE
```

### Example Output

**Console:**
```
Total lines processed: 10

Errors:
500 /api/login -> 1
504 /api/products -> 1
403 /api/cart -> 1
500 /api/checkout -> 1

Latency:
avg: 387ms
max: 900ms
```

**JSON Output (`--json`):**
```json
{
  "total_lines": 10,
  "errors": {
    "500 /api/login": 1,
    "504 /api/products": 1,
    "403 /api/cart": 1,
    "500 /api/checkout": 1
  },
  "latency": {
    "avg": 387,
    "max": 900
  }
}
```

## Extending the Tool

- Add streaming support for very large log files.
- Introduce concurrency for parsing multiple files.
- Output to different formats (CSV, Prometheus metrics, etc.)

## License

MIT License
