# =========================
# Project config
# =========================
APP_NAME := bc-server
CMD_PATH := ./cmd/server
PORT ?= 8080
ROUND_DURATION ?= 0s
IMAGE_NAME := bc-mvp

GO := go
GOFLAGS :=

# =========================
# Helpers
# =========================
.PHONY: help
help:
	@echo ""
	@echo "Available targets:"
	@echo "  make run           - run server locally"
	@echo "  make test          - run all tests"
	@echo "  make test-race     - run tests with race detector"
	@echo "  make fmt           - gofmt all files"
	@echo "  make lint          - run golangci-lint (if installed)"
	@echo "  make build         - build local binary"
	@echo "  make docker-build  - build docker image"
	@echo "  make docker-run    - run docker container"
	@echo "  make clean         - remove build artifacts"
	@echo ""

# =========================
# Local development
# =========================
.PHONY: run
run:
	PORT=$(PORT) ROUND_DURATION=$(ROUND_DURATION) \
	$(GO) run $(CMD_PATH)

# =========================
# Build
# =========================
.PHONY: build
build:
	CGO_ENABLED=0 $(GO) build -o bin/$(APP_NAME) $(CMD_PATH)

# =========================
# Tests
# =========================
.PHONY: test
test:
	$(GO) test ./...

.PHONY: test-race
test-race:
	$(GO) test ./... -race

# =========================
# Formatting & lint
# =========================
.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: fmt-check
fmt-check:
	@fmt_out=$$(gofmt -l .); \
	if [ -n "$$fmt_out" ]; then \
		echo "gofmt required on:"; \
		echo "$$fmt_out"; \
		exit 1; \
	fi

.PHONY: lint
lint:
	golangci-lint run

# =========================
# Docker
# =========================
.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE_NAME) .

.PHONY: docker-run
docker-run:
	docker run --rm \
		-p $(PORT):$(PORT) \
		-e PORT=$(PORT) \
		-e ROUND_DURATION=$(ROUND_DURATION) \
		$(IMAGE_NAME)

# =========================
# Cleanup
# =========================
.PHONY: clean
clean:
	rm -rf bin
