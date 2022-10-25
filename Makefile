PROJECT := "avito-test"

.PHONY: up-app
up-app: ## Run project in docker environment
	docker-compose -p $(PROJECT) up -d --no-deps --build server

.PHONY: up-db
up-db:
	docker-compose -p $(PROJECT) up -d balancedb

.PHONY: down
down: ## Stop project in docker environment
	docker-compose -p $(PROJECT) --env-file .env down

.PHONY: logs
logs: ## View project logs from the docker container
	docker-compose -p $(PROJECT) logs -f

.PHONY: lint
lint: ## Run all the linters
	golangci-lint run -v --color=always --timeout 4m ./...

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help
