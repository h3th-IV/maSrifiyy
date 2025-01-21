BINARY=bin/maSrifiyy
MAIN=cmd/main.go

build:
	@echo "Building the application..."
	@go build -o $(BINARY) $(MAIN)

run: build
	@echo "Running the application..."
	@./$(BINARY)

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY)

help:
	@echo "Available targets:"
	@echo "  build   - Build the application"
	@echo "  run     - Build and run the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean build artifacts"
