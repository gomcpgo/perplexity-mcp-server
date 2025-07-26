package perplexity

import (
	"context"
	"fmt"

	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/types"
)

// Search performs a general web search using Perplexity API
func (c *Client) Search(ctx context.Context, params map[string]interface{}, cfg *config.Config) (string, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return "", fmt.Errorf("query parameter is required")
	}

	// Build request with default model for general search
	req := buildRequest(query, params, cfg.DefaultModel, cfg.MaxTokens, cfg.Temperature)

	// Apply config defaults if not specified in params
	// Citations are always returned (set to true in buildRequest)
	if _, ok := params["return_images"]; !ok {
		req.ReturnImages = cfg.ReturnImages
	}
	if _, ok := params["return_related_questions"]; !ok {
		req.ReturnRelatedQuestions = cfg.ReturnRelated
	}

	// Make API call
	resp, err := c.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return formatResponse(resp), nil
}

// AcademicSearch performs an academic-focused search
func (c *Client) AcademicSearch(ctx context.Context, params map[string]interface{}, cfg *config.Config) (string, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return "", fmt.Errorf("query parameter is required")
	}

	// Use sonar-pro model for academic search if not specified (better depth for scholarly content)
	if _, ok := params["model"]; !ok {
		params["model"] = types.ModelSonarPro
	}

	// Set academic search mode
	params["search_mode"] = "academic"

	// Citations are always returned (set to true in buildRequest)

	// Set higher context size for academic content
	if _, ok := params["search_context_size"]; !ok {
		params["search_context_size"] = float64(10)
	}

	// Build request
	req := buildRequest(query, params, cfg.DefaultModel, cfg.MaxTokens, cfg.Temperature)

	// Handle subject area if provided
	if subjectArea, ok := params["subject_area"].(string); ok && subjectArea != "" {
		// Prepend subject area to the query for better context
		req.Messages[0].Content = fmt.Sprintf("[Subject: %s] %s", subjectArea, query)
	}

	// Make API call
	resp, err := c.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return formatResponse(resp), nil
}

// FinancialSearch performs a financial/SEC filing focused search
func (c *Client) FinancialSearch(ctx context.Context, params map[string]interface{}, cfg *config.Config) (string, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return "", fmt.Errorf("query parameter is required")
	}

	// Use sonar-pro model for financial search if not specified (comprehensive data needed)
	if _, ok := params["model"]; !ok {
		params["model"] = types.ModelSonarPro
	}

	// Citations are always returned (set to true in buildRequest)

	// Build request
	req := buildRequest(query, params, cfg.DefaultModel, cfg.MaxTokens, cfg.Temperature)

	// Handle financial-specific parameters
	var contextAdditions []string

	if ticker, ok := params["ticker"].(string); ok && ticker != "" {
		contextAdditions = append(contextAdditions, fmt.Sprintf("Ticker: %s", ticker))
	}

	if companyName, ok := params["company_name"].(string); ok && companyName != "" {
		contextAdditions = append(contextAdditions, fmt.Sprintf("Company: %s", companyName))
	}

	if reportType, ok := params["report_type"].(string); ok && reportType != "" {
		contextAdditions = append(contextAdditions, fmt.Sprintf("Report Type: %s", reportType))
	}

	// Add financial context to query
	if len(contextAdditions) > 0 {
		contextStr := ""
		for i, addition := range contextAdditions {
			if i > 0 {
				contextStr += ", "
			}
			contextStr += addition
		}
		req.Messages[0].Content = fmt.Sprintf("[%s] %s", contextStr, query)
	}

	// Make API call
	resp, err := c.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return formatResponse(resp), nil
}

// FilteredSearch performs an advanced search with comprehensive filtering options
func (c *Client) FilteredSearch(ctx context.Context, params map[string]interface{}, cfg *config.Config) (string, error) {
	query, ok := params["query"].(string)
	if !ok || query == "" {
		return "", fmt.Errorf("query parameter is required")
	}

	// Use sonar-pro model for filtered search if not specified
	if _, ok := params["model"]; !ok {
		params["model"] = types.ModelSonarPro
	}

	// Build request
	req := buildRequest(query, params, cfg.DefaultModel, cfg.MaxTokens, cfg.Temperature)

	// Handle advanced filtering parameters
	var filterContext []string

	if contentType, ok := params["content_type"].(string); ok && contentType != "" {
		filterContext = append(filterContext, fmt.Sprintf("Content Type: %s", contentType))
	}

	if fileType, ok := params["file_type"].(string); ok && fileType != "" {
		filterContext = append(filterContext, fmt.Sprintf("File Type: %s", fileType))
	}

	if language, ok := params["language"].(string); ok && language != "" {
		filterContext = append(filterContext, fmt.Sprintf("Language: %s", language))
	}

	if country, ok := params["country"].(string); ok && country != "" {
		filterContext = append(filterContext, fmt.Sprintf("Country: %s", country))
		// Also set location parameter if not already set
		if req.Location == "" {
			req.Location = country
		}
	}

	// Add filter context to query if any filters are specified
	if len(filterContext) > 0 {
		contextStr := ""
		for i, filter := range filterContext {
			if i > 0 {
				contextStr += ", "
			}
			contextStr += filter
		}
		req.Messages[0].Content = fmt.Sprintf("[Filters: %s] %s", contextStr, query)
	}

	// Handle custom filters
	if customFilters, ok := params["custom_filters"].(map[string]interface{}); ok {
		customContext := ""
		for key, value := range customFilters {
			if customContext != "" {
				customContext += ", "
			}
			customContext += fmt.Sprintf("%s: %v", key, value)
		}
		if customContext != "" {
			req.Messages[0].Content = fmt.Sprintf("[Custom Filters: %s] %s", customContext, req.Messages[0].Content)
		}
	}

	// Make API call
	resp, err := c.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return formatResponse(resp), nil
}