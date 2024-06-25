# Ekolod

Ekolod is a simple, efficient HTTP probe for monitoring web services, with built-in support for OpenTelemetry metrics and Prometheus export.

## Features

- Configurable HTTP probing of multiple targets
- OpenTelemetry metrics integration
- Prometheus metrics endpoint
- YAML-based configuration
- Docker support for easy deployment
- Hot reload of configuration

## Roadmap

The following features are planned for future releases:

- [ ] Dynamic Target Management: API endpoints to add, remove, or modify probe targets at runtime
- [ ] Advanced Checks: More sophisticated checks beyond HTTP status codes (e.g., response body content matching, header checks)
- [ ] Alerting: Integration with alerting systems to notify users when probes fail or meet certain conditions
- [ ] Web UI: Simple web interface to view probe status and metrics
- [ ] Authentication: Basic authentication for metrics and management endpoints
- [ ] Structured Logging: Improved logging for easier debugging and monitoring
- [ ] Testing: Comprehensive unit and integration tests
- [x] Docker Support: Dockerfile for containerization
- [ ] Helm Chart: For deploying the application in Kubernetes environments
- [x] Hot Reload: Ability to reload configuration without restarting the application

## Quick Start

1. Clone the repository:
   ```
   git clone https://github.com/c-j-p-nordquist/ekolod.git
   cd ekolod
   ```

2. Create a `probe/config.yaml` file:
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

3. Build and run the Docker container:
   ```
   docker-compose up --build
   ```

4. Access the endpoints:
   - Prometheus metrics: `http://localhost:8080/metrics`
   - Health check: `http://localhost:8080/health`
   - Reload configuration: `curl -X POST http://localhost:8080/reload`

## Configuration

Ekolod uses a YAML configuration file (`probe/config.yaml`) to define probe targets and server settings.

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

1. Ensure you have Docker and Docker Compose installed.
2. Clone the repository.
3. Make your changes.
4. Build and run the Docker container:
   ```
   docker-compose up --build
   ```
5. To reload the configuration after changes:
   ```
   curl -X POST http://localhost:8080/reload
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.