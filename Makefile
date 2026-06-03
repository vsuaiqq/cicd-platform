.DEFAULT_GOAL := help

.PHONY: help setup bootstrap-vault vault-seed seed up down start restart status health wait logs ui observability build

TOKEN_DIR := docker/vault/tokens

help: ## Show available commands
	@echo "CI/CD platform — local dev"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*##' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick start:  make start"
	@echo "Frontend:     make ui   (http://localhost:5173)"
	@echo "API:          http://localhost:8080"

setup: ## Create .env and token dir if missing
	@test -f .env || cp .env.example .env
	@mkdir -p $(TOKEN_DIR)

bootstrap-vault: setup ## Ensure Vault is up, seeded, and service tokens are fresh
	docker compose up -d vault
	docker compose rm -sf vault-init >/dev/null 2>&1 || true
	docker compose run --rm vault-init

vault-seed: bootstrap-vault ## Seed Vault secrets from .env and refresh service tokens

seed: vault-seed ## Alias for vault-seed

up: bootstrap-vault ## Start all backend services (Docker)
	docker compose up -d --build

down: ## Stop all services
	docker compose down

start: up wait health ## Full startup: Vault seed → Docker → health check
	@echo ""
	@echo "Backend is up."
	@echo "  API gateway   http://localhost:8080"
	@echo "  Git webhooks  http://localhost:8083  (expose via ngrok for GitHub)"
	@echo "  Vault UI      http://localhost:8200"
	@echo ""
	@echo "Next: make ui"

restart: down up wait health ## Restart stack and verify health

status: ## Show container status
	docker compose ps -a

health: ## Ping HTTP health endpoints
	@fail=0; \
	for entry in \
	  "api-gateway|http://localhost:8080/readyz" \
	  "auth|http://localhost:8082/readyz" \
	  "orchestrator|http://localhost:8081/healthz" \
	  "projects|http://localhost:8084/healthz" \
	  "runner|http://localhost:8085/healthz" \
	  "ai|http://localhost:8086/healthz" \
	  "analytics|http://localhost:8087/healthz" \
	  "git-gateway|http://localhost:8083/healthz"; \
	do \
	  name=$${entry%%|*}; url=$${entry#*|}; \
	  if curl -sf --max-time 5 "$$url" >/dev/null 2>&1; then \
	    printf "  \033[32m✓\033[0m %-14s %s\n" "$$name" "$$url"; \
	  else \
	    printf "  \033[31m✗\033[0m %-14s %s\n" "$$name" "$$url"; \
	    fail=1; \
	  fi; \
	done; \
	if [ $$fail -ne 0 ]; then exit 1; fi

wait: ## Wait until api-gateway responds (used by start/restart)
	@echo "Waiting for services..."
	@for i in $$(seq 1 60); do \
		if curl -sf --max-time 3 http://localhost:8080/readyz >/dev/null 2>&1; then \
			echo "  api-gateway ready ($${i}0s)"; \
			exit 0; \
		fi; \
		sleep 10; \
	done; \
	echo "Timeout: api-gateway did not become ready in 10 minutes"; \
	exit 1

logs: ## Follow logs from all services
	docker compose logs -f --tail=100

ui: ## Start frontend dev server (Vite)
	cd client && npm run dev

observability: setup ## Start Grafana/Loki/Tempo/Prometheus stack
	docker compose --profile observability up -d

build: ## Compile all Go services locally
	@for d in shared auth-service projects-service orchestrator-service analytics-service ai-service api-gateway git-gateway-service runner-service; do \
		echo "==> $$d" && (cd $$d && go build ./...) || exit 1; \
	done
