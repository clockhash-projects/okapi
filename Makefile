.PHONY: build run test lint clean help

BINARY_NAME=okapi

# Default target: show help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the okapi binary"
	@echo "  run         Run the application locally"
	@echo "  test        Run tests with race detector and coverage"
	@echo "  lint        Run golangci-lint"
	@echo "  clean       Remove build artifacts"

build:
	go build -trimpath -ldflags="-s -w" -o $(BINARY_NAME) main.go

run:
	go run main.go

test:
	go test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -rf tmp/
