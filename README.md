# Ekolod

Ekolod is a simple, efficient HTTP probe for monitoring web services, with built-in support for OpenTelemetry metrics and Prometheus export.

## Features

- Configurable HTTP probing of multiple targets
- OpenTelemetry metrics integration
- Prometheus metrics endpoint
- YAML-based configuration
- Easy to deploy and use

## Roadmap

The following features are planned for future releases:

- [ ] Dynamic Target Management: API endpoints to add, remove, or modify probe targets at runtime
- [ ] Advanced Checks: More sophisticated checks beyond HTTP status codes (e.g., response body content matching, header checks)
- [ ] Alerting: Integration with alerting systems to notify users when probes fail or meet certain conditions
- [ ] Web UI: Simple web interface to view probe status and metrics
- [ ] Authentication: Basic authentication for metrics and management endpoints
- [ ] Structured Logging: Improved logging for easier debugging and monitoring
- [ ] Testing: Comprehensive unit and integration tests
- [ ] Docker Support: Dockerfile for containerization
- [ ] Helm Chart: For deploying the application in Kubernetes environments
- [ ] Hot Reload: Ability to reload configuration without restarting the application

## Quick Start

1. Clone the repository:
   ```
   git clone https://github.com/c-j-p-nordquist/ekolod.git
   cd ekolod
   ```

2. Create a `config.yaml` file in the project root:
   ```yaml
   targets:
     - name: "Example"
       url: "https://example.com"
       interval: "10s"
       timeout: "5s"
     - name: "Google"
       url: "https://www.google.com"
       interval: "15s"
       timeout: "5s"

   server:
     port: 8080
   ```

3. Run the application:
   ```
   go run cmd/ekolod/main.go
   ```

4. Access the metrics:
   - Prometheus metrics: `http://localhost:8080/metrics`
   - Health check: `http://localhost:8080/health`

## Configuration

Ekolod uses a YAML configuration file (`config.yaml`) to define probe targets and server settings. The configuration file should be placed in the same directory as the executable.

### Configuration Options

- `targets`: A list of targets to probe
  - `name`: A unique name for the target
  - `url`: The URL to probe
  - `interval`: The interval between probes (e.g., "10s", "1m")
  - `timeout`: The timeout for each probe (e.g., "5s")
- `server`:
  - `port`: The port on which to serve metrics and health check endpoints

## Development

To set up the development environment:

1. Ensure you have Go 1.20 or later installed.
2. Clone the repository.
3. Run `go mod tidy` to install dependencies.
4. Make your changes.
5. Run tests with `go test ./...`.
6. Build the application with `go build -o ekolod cmd/ekolod/main.go`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.