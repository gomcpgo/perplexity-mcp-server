package search

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prasanthmj/perplexity/pkg/cache"
	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/types"
)

// Searcher handles search operations with caching
type Searcher struct {
	client *Client
	config *config.Config
}

// NewSearcher creates a new searcher instance
func NewSearcher(cfg *config.Config) (*Searcher, error) {
	client := NewClient(cfg.APIKey, cfg.Timeout)

	return &Searcher{
		client: client,
		config: cfg,
	}, nil
}

// validateSearchParams validates search parameters for security and correctness
func validateSearchParams(params *SearchParams) error {
	// Validate query length (max 10000 characters to prevent abuse)
	if len(params.Query) == 0 {
		return fmt.Errorf("query cannot be empty")
	}
	if len(params.Query) > 10000 {
		return fmt.Errorf("query too long: maximum 10000 characters")
	}

	// Validate domain filters (max 100 domains to prevent abuse)
	if len(params.SearchDomainFilter) > 100 {
		return fmt.Errorf("too many domain filters: maximum 100 domains")
	}
	if len(params.SearchExcludeDomains) > 100 {
		return fmt.Errorf("too many exclude domains: maximum 100 domains")
	}

	// Validate temperature if specified
	if params.Temperature != nil && (*params.Temperature < 0 || *params.Temperature > 2) {
		return fmt.Errorf("temperature must be between 0 and 2")
	}

	// Validate max tokens if specified
	if params.MaxTokens != nil && (*params.MaxTokens < 1 || *params.MaxTokens > 100000) {
		return fmt.Errorf("max_tokens must be between 1 and 100000")
	}

	return nil
}

// Search performs a general web search
func (s *Searcher) Search(ctx context.Context, params *SearchParams) (string, error) {
	// Validate parameters
	if err := validateSearchParams(params); err != nil {
		return "", fmt.Errorf("invalid search parameters: %w", err)
	}

	// Build request with default model for general search
	req := s.buildRequest(params, s.config.DefaultModel)

	// Apply config defaults if not specified in params
	if params.ReturnImages == nil {
		req.ReturnImages = s.config.ReturnImages
	}
	if params.ReturnRelatedQuestions == nil {
		req.ReturnRelatedQuestions = s.config.ReturnRelated
	}

	// Make API call
	resp, err := s.client.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return s.formatResponseWithCache(resp, params), nil
}

// AcademicSearch performs an academic-focused search using native search_mode: "academic"
func (s *Searcher) AcademicSearch(ctx context.Context, params *SearchParams) (string, error) {
	// Validate parameters
	if err := validateSearchParams(params); err != nil {
		return "", fmt.Errorf("invalid search parameters: %w", err)
	}

	// Use sonar-pro model for academic search if not specified
	if params.Model == "" {
		params.Model = types.ModelSonarPro
	}

	// Build request
	req := s.buildRequest(params, s.config.DefaultModel)

	// Use native academic search mode (new API feature from June 2025)
	req.SearchMode = "academic"

	// Higher context size for academic content (if not specified)
	if req.SearchContextSize == 0 {
		req.SearchContextSize = 10
	}

	// Handle subject area if provided
	if params.SubjectArea != "" {
		req.Messages[0].Content = fmt.Sprintf("[Subject: %s] %s", params.SubjectArea, params.Query)
	}

	// Make API call
	resp, err := s.client.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return s.formatResponseWithCache(resp, params), nil
}

// FinancialSearch performs a financial/SEC filing focused search with native SEC domain support
func (s *Searcher) FinancialSearch(ctx context.Context, params *SearchParams) (string, error) {
	// Validate parameters
	if err := validateSearchParams(params); err != nil {
		return "", fmt.Errorf("invalid search parameters: %w", err)
	}

	// Use sonar-pro model for financial search if not specified
	if params.Model == "" {
		params.Model = types.ModelSonarPro
	}

	// Build request
	req := s.buildRequest(params, s.config.DefaultModel)

	// Use native SEC domain filter if report type is specified (new API feature from July 2025)
	if params.ReportType != "" || params.Ticker != "" {
		req.SearchDomain = "sec"
	}

	// Handle financial-specific parameters
	var contextAdditions []string
	if params.Ticker != "" {
		contextAdditions = append(contextAdditions, fmt.Sprintf("Ticker: %s", params.Ticker))
	}
	if params.CompanyName != "" {
		contextAdditions = append(contextAdditions, fmt.Sprintf("Company: %s", params.CompanyName))
	}
	if params.ReportType != "" {
		contextAdditions = append(contextAdditions, fmt.Sprintf("Report Type: %s", params.ReportType))
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
		req.Messages[0].Content = fmt.Sprintf("[%s] %s", contextStr, params.Query)
	}

	// Make API call
	resp, err := s.client.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return s.formatResponseWithCache(resp, params), nil
}

// FilteredSearch performs an advanced search with comprehensive filtering options
func (s *Searcher) FilteredSearch(ctx context.Context, params *SearchParams) (string, error) {
	// Validate parameters
	if err := validateSearchParams(params); err != nil {
		return "", fmt.Errorf("invalid search parameters: %w", err)
	}

	// Use sonar-pro model for filtered search if not specified
	if params.Model == "" {
		params.Model = types.ModelSonarPro
	}

	// Build request
	req := s.buildRequest(params, s.config.DefaultModel)

	// Handle advanced filtering parameters
	var filterContext []string
	if params.ContentType != "" {
		filterContext = append(filterContext, fmt.Sprintf("Content Type: %s", params.ContentType))
	}
	if params.FileType != "" {
		filterContext = append(filterContext, fmt.Sprintf("File Type: %s", params.FileType))
	}
	if params.Language != "" {
		filterContext = append(filterContext, fmt.Sprintf("Language: %s", params.Language))
	}
	if params.Country != "" {
		filterContext = append(filterContext, fmt.Sprintf("Country: %s", params.Country))
		// Also set location parameter if not already set
		if req.Location == "" {
			req.Location = params.Country
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
		req.Messages[0].Content = fmt.Sprintf("[Filters: %s] %s", contextStr, params.Query)
	}

	// Handle custom filters
	if params.CustomFilters != nil && len(params.CustomFilters) > 0 {
		customContext := ""
		for key, value := range params.CustomFilters {
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
	resp, err := s.client.callAPI(ctx, req)
	if err != nil {
		return "", err
	}

	return s.formatResponseWithCache(resp, params), nil
}

// ListPrevious lists previous cached queries
func (s *Searcher) ListPrevious(ctx context.Context) (string, error) {
	if !cache.IsCachingEnabled(s.config.ResultsRootFolder) {
		return "[]", fmt.Errorf("results caching is not enabled. Set PERPLEXITY_RESULTS_ROOT_FOLDER environment variable to enable caching")
	}
	
	queries, err := cache.ListPreviousQueries(s.config.ResultsRootFolder)
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
func (s *Searcher) GetPreviousResult(ctx context.Context, uniqueID string) (string, error) {
	if !cache.IsCachingEnabled(s.config.ResultsRootFolder) {
		return "", fmt.Errorf("results caching is not enabled. Set PERPLEXITY_RESULTS_ROOT_FOLDER environment variable to enable caching")
	}
	
	result, err := cache.GetPreviousResult(s.config.ResultsRootFolder, uniqueID)
	if err != nil {
		return "", fmt.Errorf("failed to get previous result: %w", err)
	}
	
	return result, nil
}

// buildRequest creates a PerplexityRequest from search parameters
func (s *Searcher) buildRequest(params *SearchParams, defaultModel string) *types.PerplexityRequest {
	req := &types.PerplexityRequest{
		Model: defaultModel,
		Messages: []types.Message{
			{
				Role:    "user",
				Content: params.Query,
			},
		},
		MaxTokens:       s.config.MaxTokens,
		Temperature:     s.config.Temperature,
		ReturnCitations: true, // Always return citations for LLM to potentially fetch more info
	}

	// Override with provided parameters
	if params.Model != "" {
		req.Model = params.Model
	}

	if len(params.SearchDomainFilter) > 0 {
		req.SearchDomainFilter = params.SearchDomainFilter
	}

	if len(params.SearchExcludeDomains) > 0 {
		req.SearchExcludeDomains = params.SearchExcludeDomains
	}

	if params.SearchRecencyFilter != "" {
		req.SearchRecencyFilter = params.SearchRecencyFilter
	}

	if params.ReturnImages != nil {
		req.ReturnImages = *params.ReturnImages
	}

	if params.ReturnRelatedQuestions != nil {
		req.ReturnRelatedQuestions = *params.ReturnRelatedQuestions
	}

	if params.MaxTokens != nil {
		req.MaxTokens = *params.MaxTokens
	}

	if params.Temperature != nil {
		req.Temperature = *params.Temperature
	}

	if params.DateRangeStart != "" {
		req.DateRangeStart = params.DateRangeStart
	}

	if params.DateRangeEnd != "" {
		req.DateRangeEnd = params.DateRangeEnd
	}

	if params.Location != "" {
		req.Location = params.Location
	}

	return req
}

// formatResponse formats the API response for MCP
func (s *Searcher) formatResponse(resp *types.PerplexityResponse) string {
	if len(resp.Choices) == 0 {
		return "No response from Perplexity API"
	}

	content := resp.Choices[0].Message.Content

	// Prefer search_results (new API format) over citations (deprecated)
	if len(resp.SearchResults) > 0 {
		// Extract URLs for quick reference
		content += "\n\n## Source URLs\n"
		for i, result := range resp.SearchResults {
			content += fmt.Sprintf("%d. %s\n", i+1, result.URL)
		}

		// Include detailed search results with metadata
		content += "\n\n## Detailed Sources\n"
		for i, result := range resp.SearchResults {
			content += fmt.Sprintf("\n%d. **%s**\n", i+1, result.Title)
			content += fmt.Sprintf("   URL: %s\n", result.URL)
			if result.PublicationDate != "" {
				content += fmt.Sprintf("   Published: %s\n", result.PublicationDate)
			}
			if result.Snippet != "" {
				content += fmt.Sprintf("   Snippet: %s\n", result.Snippet)
			}
		}
	} else if len(resp.Citations) > 0 {
		// Fallback to citations if search_results not available (legacy API response)
		content += "\n\n## Source URLs\n"
		for i, url := range resp.Citations {
			content += fmt.Sprintf("%d. %s\n", i+1, url)
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
func (s *Searcher) formatResponseWithCache(resp *types.PerplexityResponse, params *SearchParams) string {
	content := s.formatResponse(resp)
	
	// Save to cache if caching is enabled
	if cache.IsCachingEnabled(s.config.ResultsRootFolder) {
		model := s.config.DefaultModel
		if params.Model != "" {
			model = params.Model
		}
		
		// Convert params to map for cache storage
		paramsMap := s.convertParamsToMap(params)
		
		uniqueID, err := cache.SaveResult(s.config.ResultsRootFolder, params.Query, params.SearchType, model, content, paramsMap)
		if err == nil && uniqueID != "" {
			// Return artifact-compatible JSON when caching is enabled
			return s.formatAsArtifactData(uniqueID, content, params, model)
		}
		// Silently ignore cache errors - don't break the search functionality
	}
	
	return content
}

// formatAsArtifactData formats the response as artifact-compatible JSON
func (s *Searcher) formatAsArtifactData(uniqueID, content string, params *SearchParams, model string) string {
	// Get current timestamp
	timestamp := time.Now().Format(time.RFC3339)
	
	// Build file paths
	resultFile := fmt.Sprintf("%s/%s/result.md", s.config.ResultsRootFolder, uniqueID)
	metadataFile := fmt.Sprintf("%s/%s/metadata.yaml", s.config.ResultsRootFolder, uniqueID)
	
	// Create artifact-compatible data structure
	artifactData := map[string]interface{}{
		"unique_id":   uniqueID,
		"query":       params.Query,
		"search_type": params.SearchType,
		"model":       model,
		"timestamp":   timestamp,
		"status":      "completed",
		"paths": map[string]interface{}{
			"result_file":   resultFile,
			"metadata_file": metadataFile,
		},
		"parameters": s.convertParamsToMap(params),
	}
	
	// Marshal to JSON
	jsonBytes, err := json.MarshalIndent(artifactData, "", "  ")
	if err != nil {
		// Fall back to text response if JSON marshaling fails
		return content + fmt.Sprintf("\n\n**Result ID:** %s", uniqueID)
	}
	
	return string(jsonBytes)
}

// convertParamsToMap converts SearchParams to map[string]interface{} for cache storage
func (s *Searcher) convertParamsToMap(params *SearchParams) map[string]interface{} {
	result := make(map[string]interface{})
	
	result["query"] = params.Query
	result["search_type"] = params.SearchType
	
	if params.Model != "" {
		result["model"] = params.Model
	}
	if len(params.SearchDomainFilter) > 0 {
		result["search_domain_filter"] = params.SearchDomainFilter
	}
	if len(params.SearchExcludeDomains) > 0 {
		result["search_exclude_domains"] = params.SearchExcludeDomains
	}
	if params.SearchRecencyFilter != "" {
		result["search_recency_filter"] = params.SearchRecencyFilter
	}
	if params.ReturnImages != nil {
		result["return_images"] = *params.ReturnImages
	}
	if params.ReturnRelatedQuestions != nil {
		result["return_related_questions"] = *params.ReturnRelatedQuestions
	}
	if params.MaxTokens != nil {
		result["max_tokens"] = *params.MaxTokens
	}
	if params.Temperature != nil {
		result["temperature"] = *params.Temperature
	}
	if params.DateRangeStart != "" {
		result["date_range_start"] = params.DateRangeStart
	}
	if params.DateRangeEnd != "" {
		result["date_range_end"] = params.DateRangeEnd
	}
	if params.Location != "" {
		result["location"] = params.Location
	}
	
	// Add type-specific parameters
	if params.SubjectArea != "" {
		result["subject_area"] = params.SubjectArea
	}
	if params.Ticker != "" {
		result["ticker"] = params.Ticker
	}
	if params.CompanyName != "" {
		result["company_name"] = params.CompanyName
	}
	if params.ReportType != "" {
		result["report_type"] = params.ReportType
	}
	if params.ContentType != "" {
		result["content_type"] = params.ContentType
	}
	if params.FileType != "" {
		result["file_type"] = params.FileType
	}
	if params.Language != "" {
		result["language"] = params.Language
	}
	if params.Country != "" {
		result["country"] = params.Country
	}
	if params.CustomFilters != nil {
		result["custom_filters"] = params.CustomFilters
	}
	
	return result
}