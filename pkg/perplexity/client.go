package perplexity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

// handleAPIError converts API errors to meaningful error messages
func handleAPIError(statusCode int, errResp *types.ErrorResponse) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("authentication failed: invalid API key")
	case http.StatusTooManyRequests:
		return fmt.Errorf("rate limit exceeded: %s", errResp.Error.Message)
	case http.StatusBadRequest:
		return fmt.Errorf("bad request: %s", errResp.Error.Message)
	case http.StatusInternalServerError:
		return fmt.Errorf("server error: %s", errResp.Error.Message)
	default:
		return fmt.Errorf("API error (%s): %s", errResp.Error.Type, errResp.Error.Message)
	}
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
		ReturnCitations: types.DefaultReturnCitations,
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

	if citations, ok := params["return_citations"].(bool); ok {
		req.ReturnCitations = citations
	}

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

	// Append citations if available
	if len(resp.Citations) > 0 {
		content += "\n\nCitations:\n"
		for i, citation := range resp.Citations {
			content += fmt.Sprintf("%d. %s\n", i+1, citation)
		}
	}

	// Append related questions if available
	if len(resp.RelatedQuestions) > 0 {
		content += "\n\nRelated Questions:\n"
		for _, question := range resp.RelatedQuestions {
			content += fmt.Sprintf("- %s\n", question)
		}
	}

	return content
}