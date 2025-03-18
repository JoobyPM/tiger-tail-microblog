# Tiger-Tail Microblog Makefile

.PHONY: all build test test-unit test-e2e clean

# Default target
all: build

# Build the application
build:
	go build -o server.out ./cmd/tigertail

# Run the application
run: build
	./server.out

# Run all tests
test: test-unit test-e2e

# Run unit tests
test-unit:
	go test -v ./internal/... ./cmd/...

# Run integration/e2e tests
test-e2e:
	./test/e2e/run_tests.sh

# Clean build artifacts
clean:
	rm -f server.out
	go clean
