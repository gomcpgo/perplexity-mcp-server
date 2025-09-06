package perplexity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/prasanthmj/perplexity/pkg/cache"
	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/types"
)

const (
	baseURL = "https://api.perplexity.ai/chat/completions"
)

// Client represents a Perplexity API client
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Perplexity API client
func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

// callAPI makes a request to the Perplexity API
func (c *Client) callAPI(ctx context.Context, req *types.PerplexityRequest) (*types.PerplexityResponse, error) {
	// Marshal request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Handle errors
	if resp.StatusCode != http.StatusOK {
		var errResp types.ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, handleAPIError(resp.StatusCode, &errResp)
	}

	// Parse successful response
	var perplexityResp types.PerplexityResponse
	if err := json.Unmarshal(body, &perplexityResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &perplexityResp, nil
}

// handleAPIError converts API errors to meaningful error messages with helpful hints
func handleAPIError(statusCode int, errResp *types.ErrorResponse) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("authentication failed: invalid API key. Please check your PERPLEXITY_API_KEY environment variable")
	case http.StatusTooManyRequests:
		return fmt.Errorf("rate limit exceeded: %s. Try reducing request frequency or using 'sonar' model for lower rate limits", errResp.Error.Message)
	case http.StatusBadRequest:
		// Add model-specific hints
		if contains(errResp.Error.Message, "Invalid model") {
			return fmt.Errorf("bad request: %s. Use 'sonar' for quick searches or 'sonar-pro' for comprehensive searches", errResp.Error.Message)
		}
		return fmt.Errorf("bad request: %s. Check your query parameters and try simplifying the request", errResp.Error.Message)
	case http.StatusInternalServerError:
		return fmt.Errorf("server error: %s. The Perplexity API is experiencing issues, please try again later", errResp.Error.Message)
	default:
		return fmt.Errorf("API error (%s): %s", errResp.Error.Type, errResp.Error.Message)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstring(s, substr, 0)
}

func containsSubstring(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[start+i] != substr[i] {
			if start+1 < len(s) {
				return containsSubstring(s, substr, start+1)
			}
			return false
		}
	}
	return true
}

// buildRequest creates a PerplexityRequest from search parameters
func buildRequest(query string, params map[string]interface{}, defaultModel string, defaultMaxTokens int, defaultTemperature float64) *types.PerplexityRequest {
	req := &types.PerplexityRequest{
		Model: defaultModel,
		Messages: []types.Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		MaxTokens:       defaultMaxTokens,
		Temperature:     defaultTemperature,
		ReturnCitations: true, // Always return citations for LLM to potentially fetch more info
	}

	// Override with provided parameters
	if model, ok := params["model"].(string); ok && model != "" {
		req.Model = model
	}

	if domains, ok := params["search_domain_filter"].([]string); ok {
		req.SearchDomainFilter = domains
	}

	if excludeDomains, ok := params["search_exclude_domains"].([]string); ok {
		req.SearchExcludeDomains = excludeDomains
	}

	if recency, ok := params["search_recency_filter"].(string); ok && recency != "" {
		req.SearchRecencyFilter = recency
	}

	// Always return citations for LLM to fetch more info if needed
	// Even if user sets it to false, we override it
	req.ReturnCitations = true

	if images, ok := params["return_images"].(bool); ok {
		req.ReturnImages = images
	}

	if related, ok := params["return_related_questions"].(bool); ok {
		req.ReturnRelatedQuestions = related
	}

	if maxTokens, ok := params["max_tokens"].(float64); ok {
		req.MaxTokens = int(maxTokens)
	}

	if temperature, ok := params["temperature"].(float64); ok {
		req.Temperature = temperature
	}

	if topP, ok := params["top_p"].(float64); ok {
		req.TopP = topP
	}

	if topK, ok := params["top_k"].(float64); ok {
		req.TopK = int(topK)
	}

	if searchMode, ok := params["search_mode"].(string); ok && searchMode != "" {
		req.SearchMode = searchMode
	}

	if citationQuality, ok := params["citation_quality"].(string); ok && citationQuality != "" {
		req.CitationQuality = citationQuality
	}

	if dateStart, ok := params["date_range_start"].(string); ok && dateStart != "" {
		req.DateRangeStart = dateStart
	}

	if dateEnd, ok := params["date_range_end"].(string); ok && dateEnd != "" {
		req.DateRangeEnd = dateEnd
	}

	if location, ok := params["location"].(string); ok && location != "" {
		req.Location = location
	}

	if contextSize, ok := params["search_context_size"].(float64); ok {
		req.SearchContextSize = int(contextSize)
	}

	return req
}

// formatResponse formats the API response for MCP
func formatResponse(resp *types.PerplexityResponse) string {
	if len(resp.Choices) == 0 {
		return "No response from Perplexity API"
	}

	content := resp.Choices[0].Message.Content

	// Always append source URLs if available (for LLM to fetch if needed)
	if len(resp.Citations) > 0 {
		content += "\n\n## Source URLs\n"
		for i, url := range resp.Citations {
			content += fmt.Sprintf("%d. %s\n", i+1, url)
		}
	}

	// Include detailed search results if available
	if len(resp.SearchResults) > 0 {
		content += "\n\n## Detailed Sources\n"
		for i, result := range resp.SearchResults {
			content += fmt.Sprintf("\n%d. **%s**\n", i+1, result.Title)
			content += fmt.Sprintf("   URL: %s\n", result.URL)
			if result.Snippet != "" {
				content += fmt.Sprintf("   Snippet: %s\n", result.Snippet)
			}
		}
	}

	// Append related questions if available
	if len(resp.RelatedQuestions) > 0 {
		content += "\n\n## Related Questions\n"
		for _, question := range resp.RelatedQuestions {
			content += fmt.Sprintf("- %s\n", question)
		}
	}

	return content
}

// formatResponseWithCache formats the API response and handles caching
func formatResponseWithCache(resp *types.PerplexityResponse, query, searchType string, params map[string]interface{}, cfg *config.Config) string {
	content := formatResponse(resp)
	
	// Save to cache if caching is enabled
	if cache.IsCachingEnabled(cfg.ResultsRootFolder) {
		model := cfg.DefaultModel
		if paramModel, ok := params["model"].(string); ok && paramModel != "" {
			model = paramModel
		}
		
		uniqueID, err := cache.SaveResult(cfg.ResultsRootFolder, query, searchType, model, content, params)
		if err == nil && uniqueID != "" {
			content += fmt.Sprintf("\n\n**Result ID:** %s", uniqueID)
		}
		// Silently ignore cache errors - don't break the search functionality
	}
	
	return content
}

// ListPrevious lists previous cached queries
func (c *Client) ListPrevious(ctx context.Context, cfg *config.Config) (string, error) {
	if !cache.IsCachingEnabled(cfg.ResultsRootFolder) {
		return "[]", fmt.Errorf("results caching is not enabled. Set PERPLEXITY_RESULTS_ROOT_FOLDER environment variable to enable caching")
	}
	
	queries, err := cache.ListPreviousQueries(cfg.ResultsRootFolder)
	if err != nil {
		return "", fmt.Errorf("failed to list previous queries: %w", err)
	}
	
	if len(queries) == 0 {
		return "[]", fmt.Errorf("no previous queries found. The results folder may be empty or not configured properly")
	}
	
	// Convert to JSON
	jsonBytes, err := json.MarshalIndent(queries, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format query list: %w", err)
	}
	
	return string(jsonBytes), nil
}

// GetPreviousResult retrieves a cached result by unique ID
func (c *Client) GetPreviousResult(ctx context.Context, params map[string]interface{}, cfg *config.Config) (string, error) {
	if !cache.IsCachingEnabled(cfg.ResultsRootFolder) {
		return "", fmt.Errorf("results caching is not enabled. Set PERPLEXITY_RESULTS_ROOT_FOLDER environment variable to enable caching")
	}
	
	uniqueID, ok := params["unique_id"].(string)
	if !ok || uniqueID == "" {
		return "", fmt.Errorf("unique_id parameter is required")
	}
	
	result, err := cache.GetPreviousResult(cfg.ResultsRootFolder, uniqueID)
	if err != nil {
		return "", fmt.Errorf("failed to get previous result: %w", err)
	}
	
	return result, nil
}