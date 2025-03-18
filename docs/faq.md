# Frequently Asked Questions

## General Questions

### What is Tiger-Tail Microblog?

Tiger-Tail Microblog is a lightweight, high-performance microblogging platform built with Go. It allows users to create short posts, follow other users, and engage with content through likes and comments.

### Is Tiger-Tail Microblog open source?

Yes, Tiger-Tail Microblog is open source and available under the [LICENSE](../LICENSE.md) included in this repository.

### What are the system requirements?

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Redis (optional, for caching)
- 1GB RAM minimum (2GB+ recommended for production)
- 1 CPU core minimum (2+ recommended for production)

## Development Questions

### How do I contribute to the project?

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write tests for your changes
5. Submit a pull request

### How do I report a bug?

Please open an issue on the GitHub repository with the following information:
- Description of the bug
- Steps to reproduce
- Expected behavior
- Actual behavior
- System information (OS, Go version, etc.)

### How do I run the tests?

```
go test ./...
```

For integration tests that require a database:

```
go test ./... -tags=integration
```

## Deployment Questions

### How do I deploy to production?

See the [Getting Started](getting_started.md) guide for deployment instructions.

### How do I scale the application?

The application is designed to be horizontally scalable. You can run multiple instances behind a load balancer.

### How do I monitor the application?

The application exposes Prometheus metrics at the `/metrics` endpoint. You can use Prometheus and Grafana to monitor the application.

## API Questions

### Is there a rate limit on API requests?

Yes, the API has rate limiting to prevent abuse. The default limit is 100 requests per minute per IP address.

### How do I authenticate API requests?

See the [API Endpoints](api_endpoints.md) documentation for authentication details.

### Can I use the API in my application?

Yes, the API is designed to be used by third-party applications. See the [API Endpoints](api_endpoints.md) documentation for details.
