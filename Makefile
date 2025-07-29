CC = go
BINARY_NAME = recall
MIN_GO_VERSION = 1.21

.PHONY: build install clean test run check-go

# Check if Go is installed and meets minimum version requirement
check-go:
	@command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed. Please install Go $(MIN_GO_VERSION) or later."; exit 1; }
	@echo "✓ Go $$(go version | cut -d' ' -f3) detected"

build: check-go
	go build -o $(BINARY_NAME) ./src
	@echo "✓ Build successful! Binary created: $(BINARY_NAME)"

install: check-go build
	sudo cp $(BINARY_NAME) /usr/local/bin/

clean:
	go clean
	rm -f $(BINARY_NAME)

test: check-go
	go test ./src/...

run: check-go
	go run ./src

# Development helpers
dev-build:
	go build -o $(BINARY_NAME) ./src

dev-test:
	./$(BINARY_NAME)
