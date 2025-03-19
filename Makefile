# Tiger-Tail Microblog Makefile

SHELL := bash

################################################
# STEP 1: Fail if .env is missing
################################################
ifeq (,$(wildcard .env))
  $(error ".env file is missing! Please create a .env with HUB_USERNAME, REPO_NAME, VERSION, etc.")
endif

################################################
# STEP 2: Include .env, then export the vars
################################################
include .env
export $(shell sed 's/=.*//' .env)

################################################
# STEP 3: Fallback defaults (if not in .env)
################################################
HUB_USERNAME ?= 
REPO_NAME ?= tigertail
VERSION ?= latest

# For multi-arch builds
PLATFORMS=linux/amd64,linux/arm64

# Helm chart info
CHART_PATH=./charts/tigertail
CHART_NAME=tiger-tail
# Extract only the main chart version, not dependency versions
CHART_VERSION=$(shell grep '^version:' $(CHART_PATH)/Chart.yaml | head -1 | awk '{print $$2}')
CHART_PACKAGE=$(CHART_NAME)-$(CHART_VERSION).tgz
CHART_REGISTRY=oci://registry-1.docker.io/$(HUB_USERNAME)

################################################
# STEP 4: A special 'check-env' target that
#         fails if any required variable is empty
################################################
.PHONY: check-env
check-env:
ifeq ($(strip $(HUB_USERNAME)),)
	$(error "HUB_USERNAME is not set! Please set it in .env")
endif
ifeq ($(strip $(REPO_NAME)),)
	$(error "REPO_NAME is not set! Please set it in .env")
endif
ifeq ($(strip $(VERSION)),)
	$(error "VERSION is not set! Please set it in .env")
endif

.PHONY: all build test test-unit test-e2e clean docker-build docker-push docker-run docker-buildx

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

################################################
# STEP 5: Docker build targets, each depends
#         on 'check-env' to ensure vars are set
################################################

# Build Docker image
docker-build: check-env
	docker build -f cmd/tigertail/Dockerfile \
		-t $(HUB_USERNAME)/$(REPO_NAME):$(VERSION) .

# Push Docker image to registry
docker-push: check-env
	docker push $(HUB_USERNAME)/$(REPO_NAME):$(VERSION)

# Run Docker container
docker-run: check-env
	docker run -d -p $(SERVER_PORT):8080 --cpus=2 --name $(REPO_NAME) \
		$(HUB_USERNAME)/$(REPO_NAME):$(VERSION)

# Build and push multi-architecture Docker image
docker-buildx: check-env
	docker buildx build \
		--platform $(PLATFORMS) \
		-f cmd/tigertail/Dockerfile \
		-t $(HUB_USERNAME)/$(REPO_NAME):$(VERSION) \
		--push .

################################################
# STEP 6: Helm chart packaging and publishing
################################################

# Package the Helm chart
helm-package:
	@echo "Packaging Helm chart $(CHART_NAME) version $(CHART_VERSION)..."
	helm package $(CHART_PATH) -d ./charts

# Push the Helm chart to Docker Hub OCI registry
helm-push: check-env helm-package
	@echo "Pushing Helm chart $(CHART_PACKAGE) to $(CHART_REGISTRY)..."
	helm push ./charts/$(CHART_PACKAGE) $(CHART_REGISTRY)

# Update Helm dependencies
helm-deps:
	helm dependency update $(CHART_PATH)

# Install the Helm chart locally
helm-install: check-env helm-package
	helm install $(CHART_NAME) ./charts/$(CHART_PACKAGE)
