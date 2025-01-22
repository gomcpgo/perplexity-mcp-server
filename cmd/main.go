package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gomcpgo/mcp/pkg/handler"
	"github.com/gomcpgo/mcp/pkg/protocol"
	"github.com/gomcpgo/mcp/pkg/server"
)

const (
	perplexityAPIURL = "https://api.perplexity.ai/chat/completions"
	//defaultModel     = "llama-3.1-sonar-small-128k-online"

	defaultModel = "sonar-pro"
)

type PerplexityServer struct {
	apiKey string
}

type PerplexityRequest struct {
	Model    string              `json:"model"`
	Messages []PerplexityMessage `json:"messages"`
}

type PerplexityMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PerplexityResponse struct {
	ID        string   `json:"id"`
	Model     string   `json:"model"`
	Citations []string `json:"citations"`
	Choices   []Choice `json:"choices"`
}

type Choice struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

func NewPerplexityServer() (*PerplexityServer, error) {
	apiKey := os.Getenv("PERPLEXITY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PERPLEXITY_API_KEY environment variable is required")
	}
	return &PerplexityServer{apiKey: apiKey}, nil
}

func (s *PerplexityServer) ListTools(ctx context.Context) (*protocol.ListToolsResponse, error) {
	return &protocol.ListToolsResponse{
		Tools: []protocol.Tool{
			{
				Name:        "research",
				Description: "Search the internet and provide up-to-date information about a topic using Perplexity.ai's Sonar Pro model",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The research query or question"
						}
					},
					"required": ["query"]
				}`),
			},
		},
	}, nil
}

func (s *PerplexityServer) CallTool(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResponse, error) {
	if req.Name != "research" {
		return nil, fmt.Errorf("unknown tool: %s", req.Name)
	}

	query, ok := req.Arguments["query"].(string)
	if !ok {
		return nil, fmt.Errorf("query must be a string")
	}

	// Create Perplexity request
	perplexityReq := PerplexityRequest{
		Model: defaultModel,
		Messages: []PerplexityMessage{
			{
				Role:    "system",
				Content: "You are a helpful research assistant. Provide accurate and well-sourced information.",
			},
			{
				Role:    "user",
				Content: query,
			},
		},
	}

	// Marshal request to JSON
	reqBody, err := json.Marshal(perplexityReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", perplexityAPIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var perplexityResp PerplexityResponse
	if err := json.Unmarshal(body, &perplexityResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(perplexityResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned")
	}

	// Format response with citations
	content := perplexityResp.Choices[0].Message.Content
	if len(perplexityResp.Citations) > 0 {
		content += "\n\nSources:\n"
		for _, citation := range perplexityResp.Citations {
			content += fmt.Sprintf("- %s\n", citation)
		}
	}

	return &protocol.CallToolResponse{
		Content: []protocol.ToolContent{
			{
				Type: "text",
				Text: content,
			},
		},
	}, nil
}

func main() {
	perplexityServer, err := NewPerplexityServer()
	if err != nil {
		log.Fatal(err)
	}

	registry := handler.NewHandlerRegistry()
	registry.RegisterToolHandler(perplexityServer)

	srv := server.New(server.Options{
		Name:     "perplexity",
		Version:  "1.0.0",
		Registry: registry,
	})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
