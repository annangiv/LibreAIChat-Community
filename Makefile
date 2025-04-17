.PHONY: all build watch clean css css-dev css-build

# Source .env file before any command
include .env
export

# Default target
all: build

# Build everything
build: css-build
	go build -o LibreAIChat main.go

# Development build with watch mode
watch:
	go run main.go

# Clean build artifacts
clean:
	go mod tidy
	rm -f static/css/output.css

# CSS specific commands
css-dev:
	tailwindcss -i static/css/input.css -o static/css/output.css

css-build:
	tailwindcss -i static/css/input.css -o static/css/output.css --minify

# Install dependencies (run once)
setup:
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.4/tailwindcss-linux-x64
	chmod +x tailwindcss-linux-x64
	sudo mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss