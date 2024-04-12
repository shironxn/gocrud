BIN := main
BIN_DIR := ./bin
DOCKER_COMPOSE := docker-compose
DOCKER := docker
GOPATH := $(GOPATH)

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: help
help: ## Display usage information
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

.PHONY: tidy
tidy: ## Run go mod tidy
	@echo "Running go mod tidy..."
	@go mod tidy

.PHONY: clean
clean: ## Clean the project
	@echo "Cleaning the project..."
	@rm -rf $(BIN_DIR)

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY:test
test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

.PHONY: dev
dev: ## Run development environment
	@$(GOPATH)/bin/air

.PHONY: build
build: ## Build the project
	@echo "Building the project..."
	@go build -o $(BIN_DIR)/$(BIN) ./cmd

.PHONY: run
run: tidy build ## Run the project
	@echo "Running the project..."
	@$(BIN_DIR)/$(BIN)

.PHONY: docker-up
docker-up: ## Start Docker Compose services
	@echo "Starting Docker Compose services..."
	@$(DOCKER_COMPOSE) up -d

.PHONY: docker-down
docker-down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@$(DOCKER_COMPOSE) down
