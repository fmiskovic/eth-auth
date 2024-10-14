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

.PHONY: clean
clean:
	@echo "=== Cleaning $(PROJECT_NAME)..."
	@rm -f $(PROJECT_NAME)