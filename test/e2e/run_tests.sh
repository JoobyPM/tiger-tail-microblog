#!/bin/bash
set -e

# Change to the e2e test directory
cd "$(dirname "$0")"

# Create a .env file for testing if it doesn't exist
if [ ! -f .env ]; then
  cp ../../.env.example .env
  # Override with test-specific values if needed
  echo "DB_NAME=tigertail_test" >> .env
fi

# Ensure we clean up on exit
function cleanup {
  echo "Cleaning up..."
  # Show logs before cleaning up
  echo "Container logs:"
  docker-compose -f docker-compose.test.yml logs
  docker-compose -f docker-compose.test.yml down -v
}
trap cleanup EXIT

# Run go mod vendor to fix vendoring issues
echo "Running go mod vendor to fix vendoring issues..."
go mod vendor

# Start the services
echo "Starting services..."
docker-compose -f docker-compose.test.yml up -d --build

# Check container status
echo "Checking container status..."
docker-compose -f docker-compose.test.yml ps

# Wait a bit for the service to start
echo "Waiting for tigertail service to start..."
sleep 10

# Try to access the livez endpoint directly
echo "Trying to access the livez endpoint directly..."
curl -v http://localhost:8080/livez || echo "Failed to access livez endpoint"

# Wait for services to be ready
echo "Waiting for services to be ready..."
docker-compose -f docker-compose.test.yml run --rm test

# Show logs if tests failed
if [ $? -ne 0 ]; then
  echo "Tests failed. Showing logs..."
  docker-compose -f docker-compose.test.yml logs
  exit 1
fi

echo "Tests completed successfully!"
