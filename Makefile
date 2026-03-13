.PHONY: help setup api-setup ui-setup api-dev ui-dev api-build ui-build ui-lint ui-format docker-up docker-down

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Setup
api-setup: ## Set up API (copy env, download Go modules)
	cd api && cp -n .env.example .env || true
	cd api && go mod download

ui-setup: ## Set up UI (copy env, install npm packages)
	cd ui && cp -n .env.example .env || true
	cd ui && npm install

setup: api-setup ui-setup ## Set up both projects

# Development
api-dev: ## Start API with hot reload (requires air)
	cd api && air

ui-dev: ## Start UI dev server on port 3001
	cd ui && npm run dev

# Build
api-build: ## Build API binary
	cd api && go build -o ./tmp/main .

ui-build: ## Build UI for production
	cd ui && npm run build

# Lint / Format
ui-lint: ## Lint UI code
	cd ui && npm run lint

ui-format: ## Format UI code
	cd ui && npm run format

# Docker
docker-up: ## Start all services via Docker Compose
	docker compose up --build

docker-down: ## Stop all services
	docker compose down
