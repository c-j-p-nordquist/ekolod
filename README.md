# Ekolod

Ekolod is a simple, efficient HTTP probe for monitoring web services.

## Quick Start

1. Clone the repository:
   ```
   git clone https://github.com/your-username/ekolod.git
   cd ekolod
   ```

2. Run the probe:
   ```
   go run cmd/ekolod/main.go
   ```

This will start probing https://example.com every 10 seconds and print the results.

## Features

- Simple HTTP probing
- Configurable targets
- JSON output

## TODO

- [ ] Support for multiple targets
- [ ] HTTP server to expose probe results
- [ ] More sophisticated checks (response time thresholds, content checks)
- [ ] In-memory storage for recent probe results
- [ ] Basic OpenTelemetry instrumentation