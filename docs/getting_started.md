# Getting Started with Tiger-Tail Microblog

## Prerequisites

- Go 1.21 or higher
- Docker (for local development)
- Git

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/JoobyPM/tiger-tail-microblog.git
   cd tiger-tail-microblog
   ```

2. Install dependencies:
   ```
   go mod download
   ```

## Configuration

The application can be configured using environment variables or a configuration file.

### Environment Variables

- `TT_DB_HOST`: Database host
- `TT_DB_PORT`: Database port
- `TT_DB_USER`: Database username
- `TT_DB_PASSWORD`: Database password
- `TT_DB_NAME`: Database name
- `TT_SERVER_PORT`: Server port

### Configuration File

Create a `config.yaml` file in the root directory:

```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: password
  name: tigertail

server:
  port: 8080
```

## Running the Application

### Development Mode

```
go run cmd/tigertail/main.go
```

### Production Mode

```
go build -o tigertail cmd/tigertail/main.go
./tigertail
```

## Testing

```
go test ./...
```

## Docker Deployment

```
docker build -t tigertail .
docker run -p 8080:8080 tigertail
```

## Kubernetes Deployment

See the `charts/tigertail` directory for Helm charts.
