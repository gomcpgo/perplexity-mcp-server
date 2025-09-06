#!/bin/bash

# Source .env file if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

## To set your API key, create a .env file or export PERPLEXITY_API_KEY="your-api-key-here"

# Show usage if no command provided
function show_usage() {
    echo "Usage: ./run.sh [command] [options]"
    echo ""
    echo "Build and Server Commands:"
    echo "  build                          Build the perplexity MCP server binary"
    echo "  run                           Run the perplexity MCP server in MCP mode"
    echo "  test                          Run unit tests"
    echo ""
    echo "Terminal Mode Testing Commands:"
    echo "  search <query> [model]        Test general search"
    echo "  academic <query> [model]      Test academic search"
    echo "  financial <query> [model]     Test financial search"
    echo "  filtered <query> [model]      Test filtered search"
    echo "  list                          List previous cached queries"
    echo "  get <result_id>               Get cached result by unique ID"
    echo ""
    echo "Integration Testing:"
    echo "  integration-test              Run integration tests against real API"
    echo ""
    echo "Examples:"
    echo "  ./run.sh search 'latest AI news' sonar-pro"
    echo "  ./run.sh academic 'quantum computing' sonar-pro"
    echo "  ./run.sh financial 'AAPL earnings' sonar-pro"
    echo "  ./run.sh list"
    echo "  ./run.sh get ABC123XYZ0"
    echo ""
    exit 1
}

# Handle different commands
case "$1" in
    build)
        echo "Building perplexity MCP server..."
        go build -o perplexity ./cmd
        ;;
    
    test)
        echo "Running unit tests..."
        go test ./pkg/...
        ;;
    
    search)
        if [ -z "$2" ]; then
            echo "Usage: ./run.sh search <query> [model]"
            exit 1
        fi
        echo "Testing general search: '$2'"
        if [ -n "$3" ]; then
            go run ./cmd -search "$2" -model "$3"
        else
            go run ./cmd -search "$2"
        fi
        ;;
    
    academic)
        if [ -z "$2" ]; then
            echo "Usage: ./run.sh academic <query> [model]"
            exit 1
        fi
        echo "Testing academic search: '$2'"
        if [ -n "$3" ]; then
            go run ./cmd -academic "$2" -model "$3"
        else
            go run ./cmd -academic "$2"
        fi
        ;;
    
    financial)
        if [ -z "$2" ]; then
            echo "Usage: ./run.sh financial <query> [model]"
            exit 1
        fi
        echo "Testing financial search: '$2'"
        if [ -n "$3" ]; then
            go run ./cmd -financial "$2" -model "$3"
        else
            go run ./cmd -financial "$2"
        fi
        ;;
    
    filtered)
        if [ -z "$2" ]; then
            echo "Usage: ./run.sh filtered <query> [model]"
            exit 1
        fi
        echo "Testing filtered search: '$2'"
        if [ -n "$3" ]; then
            go run ./cmd -filtered "$2" -model "$3"
        else
            go run ./cmd -filtered "$2"
        fi
        ;;
    
    list)
        echo "Listing previous cached queries..."
        go run ./cmd -list
        ;;
    
    get)
        if [ -z "$2" ]; then
            echo "Usage: ./run.sh get <result_id>"
            exit 1
        fi
        echo "Getting cached result: '$2'"
        go run ./cmd -get "$2"
        ;;
    
    integration-test)
        echo "Running integration tests against real API..."
        go run ./cmd -test
        ;;
    
    run)
        echo "Running perplexity MCP server..."
        go run ./cmd
        ;;
    
    *)
        show_usage
        ;;
esac