package search

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

// Client handles Perplexity API communication
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