.PHONY: all build test clean run

BINARY_NAME=my-go-app
GO_FILES=$(shell find . -type f -name "*.go" | grep -v "/vendor/")

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) .

test:
	@echo "Running tests..."
	go test ./...

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
