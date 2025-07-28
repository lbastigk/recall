CC = go
BINARY_NAME = recall

.PHONY: build install clean test run

build:
	go build -o $(BINARY_NAME) ./src

install: build
	sudo cp $(BINARY_NAME) /usr/local/bin/

clean:
	go clean
	rm -f $(BINARY_NAME)

test:
	go test ./src/...

run:
	go run ./src

# Development helpers
dev-build:
	go build -o $(BINARY_NAME) ./src

dev-test:
	./$(BINARY_NAME)
