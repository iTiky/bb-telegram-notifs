# Version
GIT_TAG ?= $(shell git describe --tags $(git rev-list --tags --max-count=1))
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DATE ?= $(shell date +%FT%T%z)

DOCKER_IMAGE="bbtt"
DOCKER_TAG="latest"

# Output names
BINARY_NAME := bbtt # BitBucket to Telegram

.PHONY: deps
deps:
	@echo "Downloading go.mod dependencies"
	go mod download

.PHONY: build
build: deps
	@echo "Building binary ($(BINARY_NAME)): $(GIT_TAG).$(GIT_COMMIT).$(BUILD_DATE)"
	@go build -o $(BINARY_NAME) ./cmd/main.go

.PHONY: docker-build
docker-build:
	@echo "Building docker image ($(DOCKER_IMAGE):$(DOCKER_TAG)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
