.PHONY: all api build test clean lint help wire run-inventory

# Variables
SERVICES := inventory
ROOT_DIR := $(shell pwd)
API_PROTO_FILES := $(shell find api -name "*.proto")

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: api build ## Generate API and build all services

# ===== API Generation =====
.PHONY: api-gen conf-gen
api: ## Generate API code from proto files
	@echo "Generating API code..."
	buf generate

conf-gen: ## Generate config proto for services
	@echo "Generating config proto..."
	@for service in $(SERVICES); do \
		if [ -f "services/$$service/internal/conf/conf.proto" ]; then \
			cd services/$$service && \
			protoc --proto_path=. \
				--proto_path=../../third_party \
				--go_out=paths=source_relative:. \
				internal/conf/conf.proto && \
			cd $(ROOT_DIR); \
		fi; \
	done

# ===== Build =====
.PHONY: build build-inventory
build: ## Build all services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd services/$$service && go build -o ../../bin/$$service ./cmd/$$service && cd $(ROOT_DIR); \
	done

build-inventory: ## Build inventory service
	@echo "Building inventory service..."
	cd services/inventory && go build -o ../../bin/inventory ./cmd/inventory

# ===== Wire =====
.PHONY: wire wire-inventory
wire: ## Generate wire code for all services
	@echo "Generating wire code..."
	@for service in $(SERVICES); do \
		echo "Wire $$service..."; \
		cd services/$$service/cmd/$$service && wire && cd $(ROOT_DIR); \
	done

wire-inventory: ## Generate wire code for inventory service
	@echo "Generating wire for inventory service..."
	cd services/inventory/cmd/inventory && wire

# ===== Run =====
.PHONY: run-% stop-%
run-%: ## Run a specific service (e.g., make run-inventory)
	@echo "Running $* service..."
	@cd services/$* && go run ./cmd/$* -conf ./configs

stop-%: ## Stop a specific service (e.g., make stop-inventory)
	@echo "Stopping $* service..."
	@pkill -f "$*.*-conf" || true

# ===== Test =====
.PHONY: test test-coverage
test: ## Run tests for all services
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ===== Lint =====
.PHONY: lint buf-lint
lint: ## Run linters
	@echo "Running golangci-lint..."
	golangci-lint run ./...

buf-lint: ## Run buf lint on proto files
	@echo "Running buf lint..."
	buf lint

# ===== Clean =====
.PHONY: clean clean-gen clean-bin
clean: clean-gen clean-bin ## Clean generated files and binaries

clean-gen: ## Clean generated code
	@echo "Cleaning generated code..."
	rm -rf gen/

clean-bin: ## Clean binaries
	@echo "Cleaning binaries..."
	rm -rf bin/

# ===== Dependencies =====
.PHONY: deps deps-install
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	@for service in $(SERVICES); do \
		cd services/$$service && go mod download && cd $(ROOT_DIR); \
	done

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy
	@for service in $(SERVICES); do \
		cd services/$$service && go mod tidy && cd $(ROOT_DIR); \
	done

# ===== Docker =====
.PHONY: docker-build docker-build-inventory
docker-build: ## Build docker images for all services
	@for service in $(SERVICES); do \
		echo "Building docker image for $$service..."; \
		docker build -t kratos-mono/$$service:latest -f services/$$service/Dockerfile .; \
	done

docker-build-inventory: ## Build docker image for inventory service
	docker build -t kratos-mono/inventory:latest -f services/inventory/Dockerfile .

# ===== Init =====
.PHONY: init install-tools
init: install-tools deps ## Initialize project (install tools and dependencies)

install-tools: ## Install required tools
	@echo "Installing tools..."
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@echo "Tools installed successfully!"
