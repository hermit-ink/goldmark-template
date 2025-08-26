.PHONY: all test build clean fmt lint vet bench tools help

all: test

test:
	gotestsum --format=testname -- ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	go build -v ./...

clean:
	go clean
	rm -f coverage.out coverage.html

fmt:
	gofumpt -w .
	go mod tidy

lint:
	golangci-lint run

vet:
	go vet ./...

bench:
	go test -bench=. -benchmem ./...

tools:
	go install mvdan.cc/gofumpt@latest
	go install gotest.tools/gotestsum@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin
	npm install -g @mermaid-js/mermaid-cli

help:
	@echo "Available targets:"
	@echo "  test           - Run all tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  build          - Build the package"
	@echo "  clean          - Remove build artifacts and coverage files"
	@echo "  fmt            - Format code and tidy modules"
	@echo "  lint           - Run golangci-lint (requires installation)"
	@echo "  vet            - Run go vet"
	@echo "  bench          - Run benchmarks"
	@echo "  tools          - Install development tools"
	@echo "  help           - Show this help message"