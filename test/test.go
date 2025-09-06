package test

import (
	"context"
	"fmt"
	"log"

	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/search"
)

// RunIntegrationTests runs integration tests against the real Perplexity API
func RunIntegrationTests() {
	fmt.Println("Running Perplexity MCP Server Integration Tests")
	fmt.Println("=" + repeatString("=", 50))

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create searcher
	searcher, err := search.NewSearcher(cfg)
	if err != nil {
		log.Fatalf("Failed to create searcher: %v", err)
	}
	ctx := context.Background()

	// Test cases
	tests := []struct {
		name   string
		testFn func(context.Context, *search.Searcher, *config.Config) error
	}{
		{"General Search", testGeneralSearch},
		{"Academic Search", testAcademicSearch},
		{"Financial Search", testFinancialSearch},
		{"Filtered Search", testFilteredSearch},
		{"Search with Parameters", testSearchWithParameters},
		{"Domain Filtering", testDomainFiltering},
		{"Error Handling", testErrorHandling},
	}

	// Run tests
	passed := 0
	failed := 0

	for _, test := range tests {
		fmt.Printf("\nRunning: %s\n", test.name)
		fmt.Println(repeatString("-", 30))

		err := test.testFn(ctx, searcher, cfg)
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
			failed++
		} else {
			fmt.Println("✅ PASSED")
			passed++
		}
	}

	// Summary
	fmt.Printf("\n%s\n", repeatString("=", 50))
	fmt.Printf("Test Summary: %d passed, %d failed\n", passed, failed)
	if failed > 0 {
		log.Fatal("Some tests failed")
	}
}

func testGeneralSearch(ctx context.Context, searcher *search.Searcher, cfg *config.Config) error {
	params := &search.SearchParams{
		Query:      "What is the capital of France?",
		SearchType: "general",
	}

	result, err := searcher.Search(ctx, params)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testAcademicSearch(ctx context.Context, searcher *search.Searcher, cfg *config.Config) error {
	params := &search.SearchParams{
		Query:       "quantum computing applications",
		SearchType:  "academic",
		SubjectArea: "Physics",
	}

	result, err := searcher.AcademicSearch(ctx, params)
	if err != nil {
		return fmt.Errorf("academic search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testFinancialSearch(ctx context.Context, searcher *search.Searcher, cfg *config.Config) error {
	params := &search.SearchParams{
		Query:       "latest earnings report",
		SearchType:  "financial",
		Ticker:      "AAPL",
		ReportType:  "10-K",
	}

	result, err := searcher.FinancialSearch(ctx, params)
	if err != nil {
		return fmt.Errorf("financial search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testFilteredSearch(ctx context.Context, searcher *search.Searcher, cfg *config.Config) error {
	params := &search.SearchParams{
		Query:       "artificial intelligence",
		SearchType:  "filtered",
		ContentType: "news",
		Language:    "English",
		Country:     "United States",
	}

	result, err := searcher.FilteredSearch(ctx, params)
	if err != nil {
		return fmt.Errorf("filtered search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testSearchWithParameters(ctx context.Context, searcher *search.Searcher, cfg *config.Config) error {
	maxTokens := 512
	temperature := 0.5
	returnRelated := true
	
	params := &search.SearchParams{
		Query:                    "climate change",
		SearchType:               "general",
		SearchRecencyFilter:      "week",
		ReturnRelatedQuestions:   &returnRelated,
		MaxTokens:                &maxTokens,
		Temperature:              &temperature,
	}

	result, err := searcher.Search(ctx, params)
	if err != nil {
		return fmt.Errorf("search with parameters failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	// Check if citations are included (they should always be)
	if !contains(result, "Source URLs") {
		fmt.Println("Warning: Source URLs not found in response")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testDomainFiltering(ctx context.Context, searcher *search.Searcher, cfg *config.Config) error {
	params := &search.SearchParams{
		Query:                "machine learning",
		SearchType:           "general",
		SearchDomainFilter:   []string{"arxiv.org", "nature.com"},
		SearchExcludeDomains: []string{"wikipedia.org"},
	}

	result, err := searcher.Search(ctx, params)
	if err != nil {
		return fmt.Errorf("domain filtering search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testErrorHandling(ctx context.Context, searcher *search.Searcher, cfg *config.Config) error {
	// Test with empty query
	params := &search.SearchParams{
		Query:      "",
		SearchType: "general",
	}

	_, err := searcher.Search(ctx, params)
	if err == nil {
		return fmt.Errorf("expected error for empty query, got nil")
	}

	fmt.Printf("Expected error received: %v\n", err)

	// Test with invalid model (this will be caught by API)
	params = &search.SearchParams{
		Query:      "test",
		SearchType: "general",
		Model:      "invalid-model-name",
	}

	_, err = searcher.Search(ctx, params)
	if err == nil {
		// Some APIs might not validate model immediately
		fmt.Println("Warning: Invalid model was not rejected")
	} else {
		fmt.Printf("Model validation error: %v\n", err)
	}

	return nil
}

func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}