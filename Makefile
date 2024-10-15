# project name
PROJECT_NAME = eth-auth

.PHONY: build
build:
	@echo "=== Building $(PROJECT_NAME)..."
	@go build -o $(PROJECT_NAME) main.go

# Run the application
.PHONY: run
run:
	@echo "=== Running server..."
	@go mod tidy
	@go run main.go

.PHONY: test
test:
	@echo "=== Running tests with race detector"
	go test -vet=off -count=1 -race -timeout=30s ./...

.PHONY: clean
clean:
	@echo "=== Cleaning $(PROJECT_NAME)..."
	@rm -f $(PROJECT_NAME)