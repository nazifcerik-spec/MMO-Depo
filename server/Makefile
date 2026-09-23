.PHONY: help tidy build run dev test test-race lint fmt migrate-up migrate-down migrate-new clean

BINARY_API=bin/api
BINARY_WORKER=bin/worker

MIGRATE_DIR=migrations
MIGRATE_DB=$(DATABASE_URL)

help:
	@echo "Otter server — Makefile targets"
	@echo "  tidy          go mod tidy"
	@echo "  build         build api + worker binaries"
	@echo "  run           run api (foreground)"
	@echo "  dev           run api with auto-reload (needs air)"
	@echo "  test          run unit tests"
	@echo "  test-race     run tests with race detector"
	@echo "  lint          staticcheck (needs staticcheck installed)"
	@echo "  fmt           gofmt -w ."
	@echo "  migrate-up    apply pending migrations"
	@echo "  migrate-down  revert last migration"
	@echo "  migrate-new name=<slug>  create new migration files"
	@echo "  clean         remove build artifacts"

tidy:
	go mod tidy

build:
	@mkdir -p bin
	go build -ldflags="-s -w" -o $(BINARY_API) ./cmd/api
	go build -ldflags="-s -w" -o $(BINARY_WORKER) ./cmd/worker

run:
	go run ./cmd/api

dev:
	@command -v air >/dev/null 2>&1 || { echo "air not found — install with: go install github.com/air-verse/air@latest"; exit 1; }
	air -c .air.toml

test:
	go test ./...

test-race:
	go test -race -count=1 ./...

lint:
	@command -v staticcheck >/dev/null 2>&1 || { echo "staticcheck not found — go install honnef.co/go/tools/cmd/staticcheck@latest"; exit 1; }
	staticcheck ./...

fmt:
	gofmt -w .

# --- Migrations (requires golang-migrate CLI) ---
# Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	@[ -n "$(MIGRATE_DB)" ] || { echo "DATABASE_URL not set"; exit 1; }
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_DB)" up

migrate-down:
	@[ -n "$(MIGRATE_DB)" ] || { echo "DATABASE_URL not set"; exit 1; }
	migrate -path $(MIGRATE_DIR) -database "$(MIGRATE_DB)" down 1

migrate-new:
	@[ -n "$(name)" ] || { echo "usage: make migrate-new name=add_something"; exit 1; }
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

clean:
	rm -rf bin/ coverage.*
