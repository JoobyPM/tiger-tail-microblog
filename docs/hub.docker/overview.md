# TigerTail Microblog

A minimal microblog service built in Go with PostgreSQL and Redis support.

## Features

- 🗄️ PostgreSQL persistence
- ⚡ Redis caching
- 🔒 Basic Auth for admin post creation
- 🏥 Kubernetes health checks (`/livez` and `/readyz`)
- 🐳 Docker Compose for local development
- 📦 Helm chart for production deployments

## Quick Start

```bash
# Pull the image
docker pull joobypm/tiger-tail:latest

# Run with Docker Compose
docker-compose up --build
```

## API Endpoints

- `GET /api/posts` - Fetch posts (cached in Redis)
- `POST /api/posts` - Create new post (requires Basic Auth)
- `GET /livez` - Liveness probe
- `GET /readyz` - Readiness probe

## Environment Variables

Required environment variables:
- `POSTGRES_HOST` - PostgreSQL host
- `POSTGRES_PORT` - PostgreSQL port
- `POSTGRES_USER` - PostgreSQL user
- `POSTGRES_PASSWORD` - PostgreSQL password
- `POSTGRES_DB` - PostgreSQL database name
- `REDIS_HOST` - Redis host
- `REDIS_PORT` - Redis port
- `ADMIN_USERNAME` - Basic Auth username
- `ADMIN_PASSWORD` - Basic Auth password

## Documentation

For detailed documentation, architecture overview, and FAQs, visit our [GitHub repository](https://github.com/JoobyPM/tiger-tail-microblog).

## License

MIT License 