package perplexity

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/types"
)

func createTestServer(t *testing.T, expectedModel string, verifyRequest func(*http.Request, *types.PerplexityRequest)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Decode request
		var req types.PerplexityRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request: %v", err)
		}

		// Verify model if specified
		if expectedModel != "" && req.Model != expectedModel {
			t.Errorf("Model mismatch: got %s, want %s", req.Model, expectedModel)
		}

		// Custom verification if provided
		if verifyRequest != nil {
			verifyRequest(r, &req)
		}

		// Send response
		resp := types.PerplexityResponse{
			ID:      "test-id",
			Model:   req.Model,
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Choices: []types.Choice{
				{
					Index:        0,
					FinishReason: "stop",
					Message: types.Message{
						Role:    "assistant",
						Content: "Test response for: " + req.Messages[0].Content,
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
}

func createTestConfig() *config.Config {
	return &config.Config{
		APIKey:          "test-api-key",
		DefaultModel:    types.ModelSonar,
		MaxTokens:       1024,
		Temperature:     0.2,
		TopP:            0.9,
		TopK:            0,
		Timeout:         30 * time.Second,
		ReturnCitations: true,
		ReturnImages:    false,
		ReturnRelated:   false,
	}
}

func TestSearch(t *testing.T) {
	cfg := createTestConfig()
	
	tests := []struct {
		name           string
		params         map[string]interface{}
		expectedModel  string
		verifyRequest  func(*http.Request, *types.PerplexityRequest)
		wantErr        bool
		errContains    string
	}{
		{
			name: "basic search",
			params: map[string]interface{}{
				"query": "test query",
			},
			expectedModel: types.ModelSonar,
		},
		{
			name: "search with custom model",
			params: map[string]interface{}{
				"query": "test query",
				"model": types.ModelSonarPro,
			},
			expectedModel: types.ModelSonarPro,
		},
		{
			name: "search with filters",
			params: map[string]interface{}{
				"query":                  "test query",
				"search_domain_filter":   []string{"example.com", "test.com"},
				"search_recency_filter":  types.RecencyWeek,
				"return_citations":       false,
				"return_images":          true,
			},
			expectedModel: types.ModelSonar,
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				if len(req.SearchDomainFilter) != 2 {
					t.Errorf("Expected 2 domain filters, got %d", len(req.SearchDomainFilter))
				}
				if req.SearchRecencyFilter != types.RecencyWeek {
					t.Errorf("Expected recency filter %s, got %s", types.RecencyWeek, req.SearchRecencyFilter)
				}
				if req.ReturnCitations != false {
					t.Errorf("Expected return_citations false, got %v", req.ReturnCitations)
				}
				if req.ReturnImages != true {
					t.Errorf("Expected return_images true, got %v", req.ReturnImages)
				}
			},
		},
		{
			name:        "missing query",
			params:      map[string]interface{}{},
			wantErr:     true,
			errContains: "query parameter is required",
		},
		{
			name: "empty query",
			params: map[string]interface{}{
				"query": "",
			},
			wantErr:     true,
			errContains: "query parameter is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				server := createTestServer(t, tt.expectedModel, tt.verifyRequest)
				defer server.Close()

				client := NewClient("test-api-key", 30*time.Second)
				client.baseURL = server.URL

				result, err := client.Search(context.Background(), tt.params, cfg)
				if err != nil {
					t.Fatalf("Search failed: %v", err)
				}

				if result == "" {
					t.Error("Expected non-empty result")
				}
			} else {
				client := NewClient("test-api-key", 30*time.Second)
				_, err := client.Search(context.Background(), tt.params, cfg)
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("Error mismatch: got %v, want to contain %s", err, tt.errContains)
				}
			}
		})
	}
}

func TestAcademicSearch(t *testing.T) {
	cfg := createTestConfig()
	
	tests := []struct {
		name          string
		params        map[string]interface{}
		verifyRequest func(*http.Request, *types.PerplexityRequest)
	}{
		{
			name: "basic academic search",
			params: map[string]interface{}{
				"query": "quantum computing research",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				if req.Model != types.ModelSonarReasoning {
					t.Errorf("Expected model %s, got %s", types.ModelSonarReasoning, req.Model)
				}
				if req.SearchMode != "academic" {
					t.Errorf("Expected search_mode academic, got %s", req.SearchMode)
				}
				if req.ReturnCitations != true {
					t.Errorf("Expected return_citations true, got %v", req.ReturnCitations)
				}
				if req.SearchContextSize != 10 {
					t.Errorf("Expected search_context_size 10, got %d", req.SearchContextSize)
				}
			},
		},
		{
			name: "academic search with subject area",
			params: map[string]interface{}{
				"query":        "neural networks",
				"subject_area": "Computer Science",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				expectedContent := "[Subject: Computer Science] neural networks"
				if req.Messages[0].Content != expectedContent {
					t.Errorf("Expected content %s, got %s", expectedContent, req.Messages[0].Content)
				}
			},
		},
		{
			name: "academic search with custom model",
			params: map[string]interface{}{
				"query": "physics research",
				"model": types.ModelSonarPro,
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				if req.Model != types.ModelSonarPro {
					t.Errorf("Expected model %s, got %s", types.ModelSonarPro, req.Model)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createTestServer(t, "", tt.verifyRequest)
			defer server.Close()

			client := NewClient("test-api-key", 30*time.Second)
			client.baseURL = server.URL

			result, err := client.AcademicSearch(context.Background(), tt.params, cfg)
			if err != nil {
				t.Fatalf("AcademicSearch failed: %v", err)
			}

			if result == "" {
				t.Error("Expected non-empty result")
			}
		})
	}
}

func TestFinancialSearch(t *testing.T) {
	cfg := createTestConfig()
	
	tests := []struct {
		name          string
		params        map[string]interface{}
		verifyRequest func(*http.Request, *types.PerplexityRequest)
	}{
		{
			name: "basic financial search",
			params: map[string]interface{}{
				"query": "Apple earnings report",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				if req.Model != types.ModelSonarReasoningPro {
					t.Errorf("Expected model %s, got %s", types.ModelSonarReasoningPro, req.Model)
				}
				if req.ReturnCitations != true {
					t.Errorf("Expected return_citations true, got %v", req.ReturnCitations)
				}
			},
		},
		{
			name: "financial search with ticker",
			params: map[string]interface{}{
				"query":  "earnings report",
				"ticker": "AAPL",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				expectedContent := "[Ticker: AAPL] earnings report"
				if req.Messages[0].Content != expectedContent {
					t.Errorf("Expected content %s, got %s", expectedContent, req.Messages[0].Content)
				}
			},
		},
		{
			name: "financial search with all parameters",
			params: map[string]interface{}{
				"query":        "quarterly results",
				"ticker":       "MSFT",
				"company_name": "Microsoft Corporation",
				"report_type":  "10-Q",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				expectedContent := "[Ticker: MSFT, Company: Microsoft Corporation, Report Type: 10-Q] quarterly results"
				if req.Messages[0].Content != expectedContent {
					t.Errorf("Expected content %s, got %s", expectedContent, req.Messages[0].Content)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createTestServer(t, "", tt.verifyRequest)
			defer server.Close()

			client := NewClient("test-api-key", 30*time.Second)
			client.baseURL = server.URL

			result, err := client.FinancialSearch(context.Background(), tt.params, cfg)
			if err != nil {
				t.Fatalf("FinancialSearch failed: %v", err)
			}

			if result == "" {
				t.Error("Expected non-empty result")
			}
		})
	}
}

func TestFilteredSearch(t *testing.T) {
	cfg := createTestConfig()
	
	tests := []struct {
		name          string
		params        map[string]interface{}
		verifyRequest func(*http.Request, *types.PerplexityRequest)
	}{
		{
			name: "basic filtered search",
			params: map[string]interface{}{
				"query": "technology news",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				if req.Model != types.ModelSonarPro {
					t.Errorf("Expected model %s, got %s", types.ModelSonarPro, req.Model)
				}
			},
		},
		{
			name: "filtered search with content type",
			params: map[string]interface{}{
				"query":        "AI research",
				"content_type": "academic papers",
				"file_type":    "pdf",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				expectedContent := "[Filters: Content Type: academic papers, File Type: pdf] AI research"
				if req.Messages[0].Content != expectedContent {
					t.Errorf("Expected content %s, got %s", expectedContent, req.Messages[0].Content)
				}
			},
		},
		{
			name: "filtered search with country",
			params: map[string]interface{}{
				"query":   "startup ecosystem",
				"country": "Germany",
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				if req.Location != "Germany" {
					t.Errorf("Expected location Germany, got %s", req.Location)
				}
				expectedContent := "[Filters: Country: Germany] startup ecosystem"
				if req.Messages[0].Content != expectedContent {
					t.Errorf("Expected content %s, got %s", expectedContent, req.Messages[0].Content)
				}
			},
		},
		{
			name: "filtered search with custom filters",
			params: map[string]interface{}{
				"query": "electric vehicles",
				"custom_filters": map[string]interface{}{
					"industry": "automotive",
					"year":     2024,
				},
			},
			verifyRequest: func(r *http.Request, req *types.PerplexityRequest) {
				// The exact content depends on map iteration order, so just check it contains the expected parts
				content := req.Messages[0].Content
				if !contains(content, "[Custom Filters:") {
					t.Errorf("Expected content to contain custom filters prefix")
				}
				if !contains(content, "electric vehicles") {
					t.Errorf("Expected content to contain original query")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := createTestServer(t, "", tt.verifyRequest)
			defer server.Close()

			client := NewClient("test-api-key", 30*time.Second)
			client.baseURL = server.URL

			result, err := client.FilteredSearch(context.Background(), tt.params, cfg)
			if err != nil {
				t.Fatalf("FilteredSearch failed: %v", err)
			}

			if result == "" {
				t.Error("Expected non-empty result")
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsAt(s, substr, 0)
}

func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[start+i] != substr[i] {
			return containsAt(s, substr, start+1)
		}
	}
	return true
}