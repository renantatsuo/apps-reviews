.PHONY: help init dev dev-server dev-web

help:
	@echo "Available commands:"
	@echo "  init       - Install all dependencies"
	@echo "  dev        - Start both server and web in development mode"
	@echo "  dev-server - Start Go server in development mode"
	@echo "  dev-web    - Start React development server"
	

# Variables
GO_CMD=go
NPM_CMD=npm
SERVER_DIR=server
WEB_DIR=web

# Install dependencies
init: deps-web

deps-web:
	@echo "Installing npm dependencies..."
	cd $(WEB_DIR) && $(NPM_CMD) install

# Development targets
dev:
	@echo "Starting development servers..."
	@echo "Note: This will start both servers. Use 'make dev-server' and 'make dev-web' in separate terminals for better control."
	@make dev-server &
	@make dev-web

dev-server:
	@echo "Starting Go server in development mode..."
	cd $(SERVER_DIR) && $(GO_CMD) run ./cmd/server

dev-web:
	@echo "Starting React development server..."
	cd $(WEB_DIR) && $(NPM_CMD) run dev
