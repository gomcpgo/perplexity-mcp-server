#!/bin/bash


## export PERPLEXITY_API_KEY="your-api-key-here"

# Show usage if no command provided
function show_usage() {
    echo "Usage: ./run.sh [command]"
    echo "Commands:"
    echo "  build    Build the perplexity MCP server binary"
    echo "  run      Run the perplexity MCP server"
    exit 1
}

# Handle different commands
case "$1" in
  build)
    echo "Building perplexity MCP server..."
    go build -o bin/perplexity-server ./cmd
    ;;
  run)
    echo "Running perplexity MCP server..."
    go run ./cmd/main.go
    ;;
  *)
    show_usage
    ;;
esac