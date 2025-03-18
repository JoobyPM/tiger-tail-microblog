

# TigerTail Microblog

TigerTail is a minimal microblog service, built in Go, featuring:
- **PostgreSQL** persistence
- **Redis** caching
- **Basic Auth** for admin post creation
- **/livez** and **/readyz** for Kubernetes health checks
- **Docker Compose** for local dev
- **Helm** chart for production deployments (in `charts/tigertail`)

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
   - **GET** `/posts` - fetch posts from Redis (or DB if cache is empty).
   - **POST** `/posts` - create new post (requires Basic Auth with env-based credentials).
   - **GET** `/livez` - liveness probe.
   - **GET** `/readyz` - readiness probe.

See the [docs](./docs) folder for more details on architecture, usage, and FAQs.


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

## Contributing

We welcome pull requests! Please open an issue first to discuss any feature or bug fix.


## License

[MIT](./LICENSE.md)
