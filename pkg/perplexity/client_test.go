package perplexity

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prasanthmj/perplexity/pkg/types"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	timeout := 30 * time.Second

	client := NewClient(apiKey, timeout)

	if client.apiKey != apiKey {
		t.Errorf("API key mismatch: got %s, want %s", client.apiKey, apiKey)
	}
	if client.httpClient.Timeout != timeout {
		t.Errorf("Timeout mismatch: got %v, want %v", client.httpClient.Timeout, timeout)
	}
	if client.baseURL != baseURL {
		t.Errorf("Base URL mismatch: got %s, want %s", client.baseURL, baseURL)
	}
}

func TestCallAPISuccess(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Invalid authorization header: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Invalid content type: %s", r.Header.Get("Content-Type"))
		}

		// Send response
		resp := types.PerplexityResponse{
			ID:      "test-id",
			Model:   types.ModelSonar,
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Choices: []types.Choice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: types.Message{
						Role:    "assistant",
						Content: "Test response",
					},
				},
			},
			Usage: types.Usage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
			Citations: []string{"https://example.com"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-api-key", 30*time.Second)
	client.baseURL = server.URL

	// Make request
	req := &types.PerplexityRequest{
		Model: types.ModelSonar,
		Messages: []types.Message{
			{Role: "user", Content: "Test query"},
		},
	}

	resp, err := client.callAPI(context.Background(), req)
	if err != nil {
		t.Fatalf("callAPI failed: %v", err)
	}

	// Verify response
	if resp.ID != "test-id" {
		t.Errorf("ID mismatch: got %s, want test-id", resp.ID)
	}
	if len(resp.Choices) != 1 {
		t.Errorf("Choices count mismatch: got %d, want 1", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "Test response" {
		t.Errorf("Content mismatch: got %s, want Test response", resp.Choices[0].Message.Content)
	}
	if len(resp.Citations) != 1 {
		t.Errorf("Citations count mismatch: got %d, want 1", len(resp.Citations))
	}
}

func TestCallAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   types.ErrorResponse
		wantErr    string
	}{
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			response: types.ErrorResponse{
				Error: struct {
					Type    string `json:"type"`
					Message string `json:"message"`
					Code    string `json:"code,omitempty"`
				}{
					Type:    "authentication_error",
					Message: "Invalid API key",
				},
			},
			wantErr: "authentication failed: invalid API key",
		},
		{
			name:       "rate limit",
			statusCode: http.StatusTooManyRequests,
			response: types.ErrorResponse{
				Error: struct {
					Type    string `json:"type"`
					Message string `json:"message"`
					Code    string `json:"code,omitempty"`
				}{
					Type:    "rate_limit_error",
					Message: "Rate limit exceeded",
				},
			},
			wantErr: "rate limit exceeded: Rate limit exceeded",
		},
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			response: types.ErrorResponse{
				Error: struct {
					Type    string `json:"type"`
					Message string `json:"message"`
					Code    string `json:"code,omitempty"`
				}{
					Type:    "invalid_request_error",
					Message: "Invalid model specified",
				},
			},
			wantErr: "bad request: Invalid model specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			// Create client with test server URL
			client := NewClient("test-api-key", 30*time.Second)
			client.baseURL = server.URL

			// Make request
			req := &types.PerplexityRequest{
				Model: types.ModelSonar,
				Messages: []types.Message{
					{Role: "user", Content: "Test query"},
				},
			}

			_, err := client.callAPI(context.Background(), req)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("Error mismatch: got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	query := "Test query"
	params := map[string]interface{}{
		"model":                    types.ModelSonarPro,
		"search_domain_filter":     []string{"example.com"},
		"search_recency_filter":    types.RecencyWeek,
		"return_citations":         false,
		"return_images":            true,
		"max_tokens":               float64(512),
		"temperature":              0.5,
		"search_mode":              "academic",
		"date_range_start":         "2024-01-01",
		"location":                 "United States",
	}

	req := buildRequest(query, params, types.DefaultModel, types.DefaultMaxTokens, types.DefaultTemperature)

	if req.Model != types.ModelSonarPro {
		t.Errorf("Model mismatch: got %s, want %s", req.Model, types.ModelSonarPro)
	}
	if len(req.Messages) != 1 || req.Messages[0].Content != query {
		t.Errorf("Messages mismatch")
	}
	if len(req.SearchDomainFilter) != 1 || req.SearchDomainFilter[0] != "example.com" {
		t.Errorf("SearchDomainFilter mismatch")
	}
	if req.SearchRecencyFilter != types.RecencyWeek {
		t.Errorf("SearchRecencyFilter mismatch: got %s, want %s", req.SearchRecencyFilter, types.RecencyWeek)
	}
	if req.ReturnCitations != false {
		t.Errorf("ReturnCitations mismatch: got %v, want false", req.ReturnCitations)
	}
	if req.ReturnImages != true {
		t.Errorf("ReturnImages mismatch: got %v, want true", req.ReturnImages)
	}
	if req.MaxTokens != 512 {
		t.Errorf("MaxTokens mismatch: got %d, want 512", req.MaxTokens)
	}
	if req.Temperature != 0.5 {
		t.Errorf("Temperature mismatch: got %f, want 0.5", req.Temperature)
	}
	if req.SearchMode != "academic" {
		t.Errorf("SearchMode mismatch: got %s, want academic", req.SearchMode)
	}
	if req.DateRangeStart != "2024-01-01" {
		t.Errorf("DateRangeStart mismatch: got %s, want 2024-01-01", req.DateRangeStart)
	}
	if req.Location != "United States" {
		t.Errorf("Location mismatch: got %s, want United States", req.Location)
	}
}

func TestFormatResponse(t *testing.T) {
	tests := []struct {
		name     string
		response *types.PerplexityResponse
		want     string
	}{
		{
			name: "basic response",
			response: &types.PerplexityResponse{
				Choices: []types.Choice{
					{
						Message: types.Message{
							Content: "Test content",
						},
					},
				},
			},
			want: "Test content",
		},
		{
			name: "response with citations",
			response: &types.PerplexityResponse{
				Choices: []types.Choice{
					{
						Message: types.Message{
							Content: "Test content",
						},
					},
				},
				Citations: []string{"https://example.com", "https://test.com"},
			},
			want: "Test content\n\nCitations:\n1. https://example.com\n2. https://test.com\n",
		},
		{
			name: "response with related questions",
			response: &types.PerplexityResponse{
				Choices: []types.Choice{
					{
						Message: types.Message{
							Content: "Test content",
						},
					},
				},
				RelatedQuestions: []string{"Question 1?", "Question 2?"},
			},
			want: "Test content\n\nRelated Questions:\n- Question 1?\n- Question 2?\n",
		},
		{
			name: "empty response",
			response: &types.PerplexityResponse{
				Choices: []types.Choice{},
			},
			want: "No response from Perplexity API",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatResponse(tt.response)
			if got != tt.want {
				t.Errorf("formatResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}