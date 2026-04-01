.PHONY: build test coverage lint clean fmt doc help

build:
	go build -o bin/erlc-go ./...

test:
	go test -v -race -timeout 30s ./...

coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

doc:
	godoc -http=:6060

clean:
	rm -rf bin/ coverage.out coverage.html
	go clean ./...

help:
	@echo "build      - Build the project"
	@echo "test       - Run tests"
	@echo "coverage   - Generate coverage report"
	@echo "lint       - Run linter"
	@echo "fmt        - Format code"
	@echo "doc        - Serve documentation"
	@echo "clean      - Clean build artifacts"
