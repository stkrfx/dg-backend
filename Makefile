# Load environment variables from .env file if it exists
ifneq ("$(wildcard .env)","")
    include .env
    export
endif

# Main entry point path
MAIN_PACKAGE_PATH := ./cmd/server/main.go
BINARY_NAME := dg-backend

.PHONY: help run build tidy clean

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: Run the application with environment variables
run: tidy
	@echo "Starting the backend..."
	go run $(MAIN_PACKAGE_PATH)

## build: Build the production binary
build: tidy
	@echo "Building binary..."
	go build -o tmp/bin/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

## tidy: Clean up go.mod and download dependencies
tidy:
	@echo "Tidying up modules..."
	go mod tidy

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -rf tmp/bin/$(BINARY_NAME)

## db-status: Quick check if DATABASE_URL is set (Internal use)
db-status:
	@if [ -z "$(DATABASE_URL)" ]; then echo "ERROR: DATABASE_URL is not set. Check your .env file."; exit 1; else echo "DATABASE_URL is set."; fi
