# ----- 1. Build stage: Go build environment -----
FROM golang:1.21-alpine AS build

WORKDIR /app

# Copy module files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application with static linking
ENV CGO_ENABLED=0
RUN go build -o tigertail ./cmd/tigertail

# ----- 2. Final stage: Minimal runtime image -----
FROM alpine:3.18

# Add CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=build /app/tigertail /app/tigertail

# Set environment variables
ENV PORT=8080
ENV GIN_MODE=release

# Expose the application port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/tigertail"]
