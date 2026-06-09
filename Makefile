.PHONY: help dev dev-down dev-logs build test lint migrate-up migrate-down migrate-create seed clean infra-init infra-plan infra-apply infra-destroy fmt vet

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Start local development stack (Docker Compose)
	@if command -v docker-compose >/dev/null 2>&1; then 		docker-compose up -d; 	else 		docker compose up -d; 	fi
	@echo "🚀 API: http://localhost:8080 | Frontend: http://localhost:3000"

dev-down: ## Stop local development stack
	@if command -v docker-compose >/dev/null 2>&1; then 		docker-compose down; 	else 		docker compose down; 	fi

dev-logs: ## Tail all container logs
	@if command -v docker-compose >/dev/null 2>&1; then 		docker-compose logs -f; 	else 		docker compose logs -f; 	fi

build: ## Build production Docker images
	docker build -t taskmaster-api:latest ./backend
	docker build -t taskmaster-frontend:latest ./frontend

test: ## Run backend unit tests
	cd backend && go test ./tests/unit/... -v -race -cover

test-integration: ## Run integration tests (requires Docker)
	cd backend && go test ./tests/integration/... -v

lint: ## Run Go linter
	cd backend && golangci-lint run

migrate-up: ## Run database migrations up
	cd backend && migrate -path migrations -database "postgres://taskmaster:taskmaster@localhost:5432/taskmaster?sslmode=disable" up

migrate-down: ## Rollback one migration
	cd backend && migrate -path migrations -database "postgres://taskmaster:taskmaster@localhost:5432/taskmaster?sslmode=disable" down 1

migrate-create: ## Create new migration (usage: make migrate-create name=add_users)
	cd backend && migrate create -ext sql -dir migrations $(name)

seed: ## Seed demo data into local database
	@bash scripts/seed.sh 2>/dev/null || cmd /c scripts\seed.bat 2>nul || echo "Seed script not found"

clean: ## Remove all containers, volumes, and images
	@if command -v docker-compose >/dev/null 2>&1; then 		docker-compose down -v --rmi all; 	else 		docker compose down -v --rmi all; 	fi

infra-init: ## Initialize Terraform
	cd infrastructure/environments/dev && terraform init

infra-plan: ## Plan Terraform changes
	cd infrastructure/environments/dev && terraform plan

infra-apply: ## Apply Terraform changes
	cd infrastructure/environments/dev && terraform apply

infra-destroy: ## Destroy Terraform infrastructure
	cd infrastructure/environments/dev && terraform destroy

fmt: ## Format Go code
	cd backend && gofmt -w .

vet: ## Run go vet
	cd backend && go vet ./...
