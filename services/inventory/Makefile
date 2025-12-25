.PHONY: all build run test clean wire

SERVICE_NAME := inventory
CMD_DIR := ./cmd/$(SERVICE_NAME)
BIN_DIR := ../../bin
CONFIG_DIR := ./configs

all: wire build ## Generate wire and build

build: ## Build the service
	@echo "Building $(SERVICE_NAME)..."
	go build -o $(BIN_DIR)/$(SERVICE_NAME) $(CMD_DIR)

run: ## Run the service
	@echo "Running $(SERVICE_NAME)..."
	go run $(CMD_DIR) -conf $(CONFIG_DIR)

test: ## Run tests
	go test -v ./...

wire: ## Generate wire
	@echo "Generating wire..."
	cd $(CMD_DIR) && wire

clean: ## Clean build artifacts
	rm -f $(BIN_DIR)/$(SERVICE_NAME)

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
