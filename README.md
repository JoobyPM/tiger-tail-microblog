![logo](docs/images/icon.svg)
# TigerTail Microblog

TigerTail is a minimal microblog service, built in Go following our [Tiger Style](docs/tiger_style.md) principles, featuring:
- **PostgreSQL** persistence with stub support for testing
- **Redis** caching with stub support for testing
- **Basic Auth** for admin post creation
- **/livez** and **/readyz** for Kubernetes health checks
- **Docker Compose** for local dev
- **Helm** chart for production deployments (in `charts/tiger-tail-chart`)

## Quick Start

1. **Clone the repo**:
   ```bash
   git clone https://github.com/JoobyPM/tiger-tail-microblog.git
   cd tiger-tail-microblog
   ```

2. **Configure environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env file if needed
   ```

3. **Run locally** with Docker Compose:
   ```bash
   docker-compose up --build
   ```

4. **Test endpoints**:
   - **GET** `/api/posts` - fetch posts from Redis (or DB if cache is empty).
   - **POST** `/api/posts` - create new post (requires Basic Auth with env-based credentials).
   - **GET** `/livez` - liveness probe.
   - **GET** `/readyz` - readiness probe.

See the documentation for more details:
- [Architecture](docs/architecture.md) - System design and data flow
- [Getting Started](docs/getting_started.md) - Setup, running, and deployment
- [API Endpoints](docs/api_endpoints.md) - API documentation
- [FAQ](docs/faq.md) - Common questions and troubleshooting
- [Tiger Style](docs/tiger_style.md) - Our coding principles


## Development

### Running Tests

The project includes both unit tests and integration/e2e tests:

```bash
# Run unit tests
make test-unit

# Run integration/e2e tests (requires Docker)
make test-e2e

# Run all tests
make test
```

The integration tests use Docker Compose to spin up the application with PostgreSQL and Redis, then run tests against the running services.

### Docker Build and Deployment

The project includes Makefile targets for building and publishing Docker images:

1. **Configure Docker environment variables**:
   ```bash
   # Edit .env file to set Docker Hub username and version
   # Required variables: HUB_USERNAME, REPO_NAME, VERSION
   ```

2. **Build Docker image**:
   ```bash
   make docker-build
   ```

3. **Push to Docker Hub**:
   ```bash
   make docker-push
   ```

4. **Run Docker container locally**:
   ```bash
   make docker-run
   ```

5. **Build and push multi-architecture image** (amd64 and arm64):
   ```bash
   make docker-buildx
   ```

These commands will use the values from your `.env` file for Docker Hub username, repository name, and version tag.

### Helm Chart Packaging and Publishing

The project includes Makefile targets for packaging and publishing Helm charts to Docker Hub using the OCI format:

1. **Update Helm dependencies**:
   ```bash
   make helm-deps
   ```

2. **Package Helm chart**:
   ```bash
   make helm-package
   ```

3. **Push Helm chart to Docker Hub OCI registry**:
   ```bash
   make helm-push
   ```

4. **Install Helm chart locally**:
   ```bash
   make helm-install
   ```

These commands will use the chart version from `charts/tiger-tail-chart/Chart.yaml` and the Docker Hub username from your `.env` file.

## Contributing

We welcome pull requests! Please open an issue first to discuss any feature or bug fix.


## License

[MIT](./LICENSE.md)
