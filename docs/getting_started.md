# Getting Started with Tiger-Tail Microblog

This guide will help you set up, run, test, and deploy the Tiger-Tail Microblog application, following our [Tiger Style](tiger_style.md) principles.

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git
- Helm (for Kubernetes deployment)
- kubectl (for Kubernetes deployment)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/JoobyPM/tiger-tail-microblog.git
cd tiger-tail-microblog
```

### 2. Configure Environment Variables

```bash
cp .env.example .env
# Edit .env file with your configuration
```

### 3. Run with Docker Compose

The fastest way to get started is with Docker Compose, which sets up PostgreSQL, Redis, and the application:

```bash
docker-compose up --build
```

This will start the application at http://localhost:8080.

### 4. Test the API

Once the application is running, you can test the API endpoints:

```bash
# Get all posts
curl http://localhost:8080/posts

# Create a new post (requires authentication)
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(echo -n username:password | base64)" \
  -d '{"content":"Hello, Tiger-Tail!"}'
```

## Development Setup

### 1. Install Dependencies

```bash
go mod download
```

### 2. Run PostgreSQL and Redis

You can use Docker Compose to run just the dependencies:

```bash
docker-compose up postgres redis -d
```

### 3. Run the Application Locally

```bash
go run cmd/tigertail/main.go
```

## Configuration

The application can be configured using environment variables or a .env file.

### Environment Variables

| Variable       | Description              | Default   |
|----------------|--------------------------|-----------|
| DB_HOST        | PostgreSQL host          | postgres  |
| DB_PORT        | PostgreSQL port          | 5432      |
| DB_USER        | PostgreSQL username      | postgres  |
| DB_PASSWORD    | PostgreSQL password      | postgres  |
| DB_NAME        | PostgreSQL database name | tigertail |
| DB_SSLMODE     | PostgreSQL SSL mode      | disable   |
| REDIS_HOST     | Redis host               | redis     |
| REDIS_PORT     | Redis port               | 6379      |
| REDIS_PASSWORD | Redis password           |           |
| REDIS_DB       | Redis database           | 0         |
| SERVER_PORT    | Server port              | 8080      |

## Running Tests

Tiger-Tail includes comprehensive test suites following our Tiger Style principles.

### Unit Tests

```bash
make test-unit
# or
go test -v ./internal/... ./cmd/...
```

### Integration/E2E Tests

```bash
make test-e2e
# or
./test/e2e/run_tests.sh
```

### All Tests

```bash
make test
```

## Building and Deploying

### Building Docker Images

```bash
# Configure Docker Hub username in .env
HUB_USERNAME=yourusername
REPO_NAME=tigertail
VERSION=0.1.0

# Build the Docker image
make docker-build

# Push to Docker Hub
make docker-push

# Build and push multi-architecture image
make docker-buildx
```

### Deploying with Helm

Tiger-Tail includes Helm charts for Kubernetes deployment:

```bash
# Update Helm dependencies
make helm-deps

# Package the Helm chart
make helm-package

# Install locally (for testing)
make helm-install

# Push to Docker Hub OCI registry
make helm-push
```

To deploy from the OCI registry:

```bash
helm install tigertail oci://registry-1.docker.io/yourusername/tiger-tail --version 0.1.0
```

## Next Steps

- Check out the [Architecture](architecture.md) document to understand the system design
- Explore the [API Endpoints](api_endpoints.md) documentation
- Read the [FAQ](faq.md) for common questions and troubleshooting
- Review our [Tiger Style](tiger_style.md) guide for coding principles
