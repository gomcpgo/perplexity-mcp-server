package handler

import (
	"context"
	"fmt"

	"github.com/prasanthmj/perplexity/pkg/search"
)

// handlePerplexitySearch handles general web search
func (h *Handler) handlePerplexitySearch(ctx context.Context, args map[string]interface{}) (string, error) {
	params, err := h.extractSearchParams(args, "general")
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	return h.searcher.Search(ctx, params)
}

// handleAcademicSearch handles academic search
func (h *Handler) handleAcademicSearch(ctx context.Context, args map[string]interface{}) (string, error) {
	params, err := h.extractSearchParams(args, "academic")
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	// Add academic-specific parameter
	if subjectArea, ok := args["subject_area"].(string); ok && subjectArea != "" {
		params.SubjectArea = subjectArea
	}

	return h.searcher.AcademicSearch(ctx, params)
}

// handleFinancialSearch handles financial search
func (h *Handler) handleFinancialSearch(ctx context.Context, args map[string]interface{}) (string, error) {
	params, err := h.extractSearchParams(args, "financial")
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	// Add financial-specific parameters
	if ticker, ok := args["ticker"].(string); ok && ticker != "" {
		params.Ticker = ticker
	}
	if companyName, ok := args["company_name"].(string); ok && companyName != "" {
		params.CompanyName = companyName
	}
	if reportType, ok := args["report_type"].(string); ok && reportType != "" {
		params.ReportType = reportType
	}

	return h.searcher.FinancialSearch(ctx, params)
}

// handleFilteredSearch handles filtered search
func (h *Handler) handleFilteredSearch(ctx context.Context, args map[string]interface{}) (string, error) {
	params, err := h.extractSearchParams(args, "filtered")
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	// Add filtering-specific parameters
	if contentType, ok := args["content_type"].(string); ok && contentType != "" {
		params.ContentType = contentType
	}
	if fileType, ok := args["file_type"].(string); ok && fileType != "" {
		params.FileType = fileType
	}
	if language, ok := args["language"].(string); ok && language != "" {
		params.Language = language
	}
	if country, ok := args["country"].(string); ok && country != "" {
		params.Country = country
	}
	if customFilters, ok := args["custom_filters"].(map[string]interface{}); ok {
		params.CustomFilters = customFilters
	}

	return h.searcher.FilteredSearch(ctx, params)
}

// handleListPrevious handles listing previous queries
func (h *Handler) handleListPrevious(ctx context.Context, args map[string]interface{}) (string, error) {
	return h.searcher.ListPrevious(ctx)
}

// handleGetPreviousResult handles getting previous results
func (h *Handler) handleGetPreviousResult(ctx context.Context, args map[string]interface{}) (string, error) {
	uniqueID, ok := args["unique_id"].(string)
	if !ok || uniqueID == "" {
		return "", fmt.Errorf("unique_id parameter is required")
	}

	return h.searcher.GetPreviousResult(ctx, uniqueID)
}

// extractSearchParams extracts common search parameters from map[string]interface{}
func (h *Handler) extractSearchParams(args map[string]interface{}, searchType string) (*search.SearchParams, error) {
	// Required parameter
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query parameter is required")
	}

	params := &search.SearchParams{
		Query:      query,
		SearchType: searchType,
	}

	// Optional parameters with type checking
	if model, ok := args["model"].(string); ok && model != "" {
		params.Model = model
	}

	if domains, ok := args["search_domain_filter"].([]interface{}); ok {
		params.SearchDomainFilter = convertToStringSlice(domains)
	}

	if excludeDomains, ok := args["search_exclude_domains"].([]interface{}); ok {
		params.SearchExcludeDomains = convertToStringSlice(excludeDomains)
	}

	if recency, ok := args["search_recency_filter"].(string); ok && recency != "" {
		params.SearchRecencyFilter = recency
	}

	if images, ok := args["return_images"].(bool); ok {
		params.ReturnImages = &images
	}

	if related, ok := args["return_related_questions"].(bool); ok {
		params.ReturnRelatedQuestions = &related
	}

	if maxTokens, ok := args["max_tokens"].(float64); ok {
		maxTokensInt := int(maxTokens)
		params.MaxTokens = &maxTokensInt
	}

	if temperature, ok := args["temperature"].(float64); ok {
		params.Temperature = &temperature
	}

	if dateStart, ok := args["date_range_start"].(string); ok && dateStart != "" {
		params.DateRangeStart = dateStart
	}

	if dateEnd, ok := args["date_range_end"].(string); ok && dateEnd != "" {
		params.DateRangeEnd = dateEnd
	}

	if location, ok := args["location"].(string); ok && location != "" {
		params.Location = location
	}

	return params, nil
}

// convertToStringSlice safely converts []interface{} to []string
func convertToStringSlice(interfaces []interface{}) []string {
	result := make([]string, 0, len(interfaces))
	for _, item := range interfaces {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}