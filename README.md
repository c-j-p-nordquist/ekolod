# Ekolod

Ekolod is a simple, efficient HTTP probe for monitoring web services, with built-in support for OpenTelemetry metrics and Prometheus export. It now includes a SvelteKit-based UI for easy visualization and management.

## Features

- Configurable HTTP probing of multiple targets
- SvelteKit-based UI for visualization and management
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

## Quick Start with Docker

1. Clone the repository:
   ```
   git clone https://github.com/c-j-p-nordquist/ekolod.git
   cd ekolod
   ```

2. Create a `probe/config.yaml` file with your desired configuration:
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

3. Build and run the Docker containers:
   ```
   docker-compose up --build
   ```

4. Access the UI at `http://localhost:5173`
5. Access the metrics endpoint at `http://localhost:8080/metrics`

## Development

To set up the development environment:

1. Ensure you have Docker and Docker Compose installed.
2. Clone the repository.
3. Run `docker-compose up --build` to start both the probe and UI in development mode.
4. The probe API will be available at `http://localhost:8080`
5. The UI development server will be available at `http://localhost:5173`

## Project Structure

- `cmd/probe/`: Contains the main application for the probe
- `pkg/`: Contains the core packages for the probe
- `probe/`: Contains probe-specific files including Dockerfile and config
- `ui/`: Contains the SvelteKit-based UI application

## Configuration

Ekolod uses a YAML configuration file (`probe/config.yaml`) to define probe targets and server settings.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
