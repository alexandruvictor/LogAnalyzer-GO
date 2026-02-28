# Log Analyzer - Go

![Go](https://img.shields.io/badge/Language-Go-blue?logo=go&logoColor=white)
![Status](https://img.shields.io/badge/Status-Production-ready-brightgreen)

## Description

**Log Analyzer** is a high-performance tool written in **Go** designed to parse large log files, extract essential metrics, and highlight critical errors.  
This project demonstrates strong skills in **Go, parsing, data structures, concurrency, and CLI tooling**, making it ideal for real-world applications or as a showcase in your developer portfolio.

---

## Features

- Parses structured log files (`timestamp level status path latency`)  
- Collects statistics:
  - Total lines processed  
  - Errors per endpoint and status code  
  - Average and maximum latency  
- Automatically ignores malformed lines  
- Modular: separate `parser` and `stats` packages for scalability  
- Easy to extend: streaming, concurrency, JSON output  

---

## Project Structure

```text
log-analyzer/
├── go.mod             # Go module definition
├── main.go            # application entry point
├── parser/            # functions for parsing log lines
│   └── parser.go
├── stats/             # functions for collecting and calculating statistics
│   └── stats.go
├── logs/              # test log files
│   └── access.log
└── README.md