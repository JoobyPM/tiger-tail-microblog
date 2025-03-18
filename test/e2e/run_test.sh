#!/bin/bash
set -e

# Change to the project root directory
cd "$(dirname "$0")/../.."

# Run the e2e tests
go test -v ./test/e2e/...
