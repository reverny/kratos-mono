.PHONY: init build run wire test clean

# Init dependencies
init:
	go mod tidy
	go mod download

# Generate wire code
wire:
	cd cmd/filemanagement && wire

# Build binary
build:
	go build -o ../../bin/filemanagement ./cmd/filemanagement

# Run service
run:
	../../bin/filemanagement -conf ./configs

# Run with live reload (requires air)
dev:
	air

# Test
test:
	go test -v ./...

# Clean
clean:
	rm -f ../../bin/filemanagement
	rm -rf ./uploads

# All: init + wire + build
all: init wire build
