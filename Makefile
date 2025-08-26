.PHONY: help init dev dev-server dev-web

help:
	@echo "Available commands:"
	@echo "  init          - Install all dependencies"
	@echo "  dev           - Start both server and web in development mode"
	@echo "  dev-server    - Start Go server in development mode"
	@echo "  dev-scheduler - Start scheduler in development mode"
	@echo "  dev-consumer  - Start consumer in development mode"
	@echo "  dev-web       - Start React development server"
	@echo "  migrate-up    - Run migrations up"
	@echo "  migrate-down  - Run migrations down"

# Variables
GO_CMD=go
NPM_CMD=npm
SERVER_DIR=server
WEB_DIR=web

# Install dependencies
init: deps-web migrate-up

deps-web:
	@echo "Installing npm dependencies..."
	cd $(WEB_DIR) && $(NPM_CMD) install

migrate-up:
	@echo "Migrating up..."
	cd $(SERVER_DIR) && goose up

migrate-down:
	@echo "Migrating down..."
	GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=file:$(SERVER_DIR)/data/database.db goose -dir $(SERVER_DIR)/migrations $(GOOSE_ARGS)

# Development targets
dev:
	@echo "Starting development servers..."
	@echo "Note: This will start both servers. Use 'make dev-server' and 'make dev-web' in separate terminals for better control."
	@make dev-server &
	@make dev-web &
	@make dev-scheduler &
	@make dev-consumer

dev-server:
	@echo "Starting Go server in development mode..."
	cd $(SERVER_DIR) && $(GO_CMD) run ./cmd/server

dev-scheduler:
	@echo "Starting scheduler in development mode..."
	cd $(SERVER_DIR) && $(GO_CMD) run ./cmd/scheduler

dev-consumer:
	@echo "Starting consumer in development mode..."
	cd $(SERVER_DIR) && $(GO_CMD) run ./cmd/consumer

dev-web:
	@echo "Starting React development server..."
	cd $(WEB_DIR) && $(NPM_CMD) run dev
