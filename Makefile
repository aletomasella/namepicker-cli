

# BUILD COMMAND
all: build

build:
	@echo "Building..."
	@go build -o namepicker cmd/namepicker/main.go

run:
	@echo "Running..."
	@go run cmd/namepicker/main.go

