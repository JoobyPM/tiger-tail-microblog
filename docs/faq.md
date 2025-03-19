# Tiger-Tail Microblog FAQ

This document addresses frequently asked questions about the Tiger-Tail Microblog application, following our [Tiger Style](tiger_style.md) principles of safety, performance, and developer experience.

## General Questions

### What is Tiger-Tail Microblog?

Tiger-Tail is a lightweight, high-performance microblogging platform built with Go. It features PostgreSQL persistence, Redis caching, and a clean, RESTful API. The application follows our Tiger Style principles, emphasizing safety, performance, and developer experience.

### Is Tiger-Tail Microblog open source?

Yes, Tiger-Tail Microblog is open source and available under the [MIT License](../LICENSE.md) included in this repository.

### What are the system requirements?

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Redis 6 or higher
- Docker and Docker Compose (for local development)
- Kubernetes and Helm (for production deployment)
- 1GB RAM minimum (2GB+ recommended for production)
- 1 CPU core minimum (2+ recommended for production)

## Environment Variables

Tiger-Tail uses environment variables for configuration. Here's a comprehensive list:

### Database Configuration

| Variable    | Description              | Default   | Required |
|-------------|--------------------------|-----------|----------|
| DB_HOST     | PostgreSQL host          | postgres  | Yes      |
| DB_PORT     | PostgreSQL port          | 5432      | Yes      |
| DB_USER     | PostgreSQL username      | postgres  | Yes      |
| DB_PASSWORD | PostgreSQL password      | postgres  | Yes      |
| DB_NAME     | PostgreSQL database name | tigertail | Yes      |
| DB_SSLMODE  | PostgreSQL SSL mode      | disable   | No       |

### Redis Configuration

| Variable       | Description    | Default | Required |
|----------------|----------------|---------|----------|
| REDIS_HOST     | Redis host     | redis   | Yes      |
| REDIS_PORT     | Redis port     | 6379    | Yes      |
| REDIS_PASSWORD | Redis password |         | No       |
| REDIS_DB       | Redis database | 0       | No       |

### Server Configuration

| Variable    | Description                              | Default     | Required |
|-------------|------------------------------------------|-------------|----------|
| SERVER_PORT | Server port                              | 8080        | No       |
| LOG_LEVEL   | Logging level (debug, info, warn, error) | info        | No       |
| ENVIRONMENT | Environment (development, production)    | development | No       |

### Docker Hub Configuration

| Variable     | Description         | Default   | Required                    |
|--------------|---------------------|-----------|-----------------------------|
| HUB_USERNAME | Docker Hub username |           | Yes (for Docker operations) |
| REPO_NAME    | Repository name     | tigertail | No                          |
| VERSION      | Version tag         | latest    | No                          |

## Common Issues and Solutions

### Connection Issues

**Issue**: Cannot connect to PostgreSQL or Redis.

**Solution**: 
1. Ensure the services are running: `docker-compose ps`
2. Check the environment variables in your `.env` file
3. Verify network connectivity: `telnet postgres 5432` or `telnet redis 6379`
4. Check service logs: `docker-compose logs postgres` or `docker-compose logs redis`

### Authentication Issues

**Issue**: Getting 401 Unauthorized when creating posts.

**Solution**:
1. Ensure you're using Basic Authentication with the correct credentials
2. Check that the Authorization header is properly formatted: `Authorization: Basic <base64-encoded-credentials>`
3. Verify that the credentials match those in your environment variables

### Performance Issues

**Issue**: API responses are slow.

**Solution**:
1. Ensure Redis is properly configured and connected for caching
2. Check database indexes: `EXPLAIN ANALYZE` your queries
3. Monitor resource usage: CPU, memory, and disk I/O
4. Consider scaling horizontally by adding more instances

### Docker Build Issues

**Issue**: Docker build fails.

**Solution**:
1. Ensure Docker is installed and running
2. Check that all required files are present
3. Verify that the Dockerfile is correctly formatted
4. Try building with verbose output: `docker build -t tigertail . --progress=plain`

## Best Practices

### Local Development

1. Use Docker Compose for local development to ensure consistent environments
2. Set `LOG_LEVEL=debug` for more detailed logs during development
3. Run tests before committing changes: `make test`
4. Use `go fmt` and `go vet` to maintain code quality

### Production Deployment

1. Always use HTTPS in production
2. Set up proper monitoring with Prometheus and Grafana
3. Configure appropriate resource limits in Kubernetes
4. Use a proper secrets management solution for sensitive environment variables
5. Implement a CI/CD pipeline for automated testing and deployment

### Performance Optimization

1. Ensure Redis caching is properly configured
2. Use connection pooling for database connections
3. Implement appropriate indexes in PostgreSQL
4. Consider using a CDN for static assets
5. Enable compression for HTTP responses

## Cross-References

- For setup instructions, see the [Getting Started](getting_started.md) guide
- For API details, see the [API Endpoints](api_endpoints.md) documentation
- For architecture information, see the [Architecture](architecture.md) document
- For coding principles, see the [Tiger Style](tiger_style.md) guide
