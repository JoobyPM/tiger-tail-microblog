#!/bin/bash
set -e

# Change to the project root directory
cd "$(dirname "$0")/../.."

# Run the tests
go test -v ./test/e2e/...
