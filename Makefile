APP_NAME := bc-server
CMD_PATH := ./cmd/server
PORT ?= 8080
ROUND_DURATION ?= 0s

# Docker compose
DC := docker compose

# DB
PG_URL ?= postgres://bc:bc@localhost:5432/bc?sslmode=disable

# Redis
REDIS_ADDR ?= localhost:6379
MATCH_TTL ?= 24h

GO := go

.PHONY: help
help:
	@echo ""
	@echo "Dev:"
	@echo "  make run                 - run server locally"
	@echo ""
	@echo "Services (docker compose):"
	@echo "  make services-up         - start postgres+redis"
	@echo "  make services-down       - stop services"
	@echo "  make services-logs       - follow logs"
	@echo "  make services-reset      - delete volumes (DANGER)"
	@echo ""
	@echo "DB migrations (goose):"
	@echo "  make migrate-up          - apply migrations"
	@echo "  make migrate-down        - rollback last migration"
	@echo "  make migrate-status      - show status"
	@echo ""
	@echo "Tests:"
	@echo "  make test                - unit tests"
	@echo "  make test-race           - unit tests with race"
	@echo "  make test-integration    - integration tests (requires redis)"
	@echo ""

# -------------------------
# Run
# -------------------------
.PHONY: run
run:
	PORT=$(PORT) ROUND_DURATION=$(ROUND_DURATION) REDIS_ADDR=$(REDIS_ADDR) MATCH_TTL=$(MATCH_TTL) \
	$(GO) run $(CMD_PATH)

# -------------------------
# Services
# -------------------------
.PHONY: services-up
up:
	$(DC) up -d postgres redis

.PHONY: services-down
down:
	$(DC) down

.PHONY: services-logs
logs:
	$(DC) logs -f --tail=200

.PHONY: services-reset
reset:
	$(DC) down -v

# -------------------------
# Migrations (requires goose installed)
# -------------------------
.PHONY: migrate-up
migrate-up:
	goose -dir db/migrations postgres "$(PG_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir db/migrations postgres "$(PG_URL)" down

.PHONY: migrate-status
migrate-status:
	goose -dir db/migrations postgres "$(PG_URL)" status

# -------------------------
# Tests
# -------------------------
.PHONY: test
test:
	$(GO) test ./... -count=1

.PHONY: test-race
test-race:
	$(GO) test ./... -race -count=1

# Интеграционные тесты: стартуем сервисы, ждём Redis и гоняем тесты с build tag "integration"
.PHONY: test-integration
test-integration: services-up
	@echo "Waiting for Redis at $(REDIS_ADDR)..."
	@until docker exec bc_redis redis-cli ping >/dev/null 2>&1; do sleep 0.2; done
	REDIS_ADDR=$(REDIS_ADDR) MATCH_TTL=$(MATCH_TTL) \
	$(GO) test ./... -count=1 -tags=integration
