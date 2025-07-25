package test

import (
	"context"
	"fmt"
	"log"

	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/perplexity"
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

	// Create client
	client := perplexity.NewClient(cfg.APIKey, cfg.Timeout)
	ctx := context.Background()

	// Test cases
	tests := []struct {
		name   string
		testFn func(context.Context, *perplexity.Client, *config.Config) error
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

		err := test.testFn(ctx, client, cfg)
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

func testGeneralSearch(ctx context.Context, client *perplexity.Client, cfg *config.Config) error {
	params := map[string]interface{}{
		"query": "What is the capital of France?",
	}

	result, err := client.Search(ctx, params, cfg)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testAcademicSearch(ctx context.Context, client *perplexity.Client, cfg *config.Config) error {
	params := map[string]interface{}{
		"query":        "quantum computing applications",
		"subject_area": "Physics",
	}

	result, err := client.AcademicSearch(ctx, params, cfg)
	if err != nil {
		return fmt.Errorf("academic search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testFinancialSearch(ctx context.Context, client *perplexity.Client, cfg *config.Config) error {
	params := map[string]interface{}{
		"query":       "latest earnings report",
		"ticker":      "AAPL",
		"report_type": "10-K",
	}

	result, err := client.FinancialSearch(ctx, params, cfg)
	if err != nil {
		return fmt.Errorf("financial search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testFilteredSearch(ctx context.Context, client *perplexity.Client, cfg *config.Config) error {
	params := map[string]interface{}{
		"query":        "artificial intelligence",
		"content_type": "news",
		"language":     "English",
		"country":      "United States",
	}

	result, err := client.FilteredSearch(ctx, params, cfg)
	if err != nil {
		return fmt.Errorf("filtered search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testSearchWithParameters(ctx context.Context, client *perplexity.Client, cfg *config.Config) error {
	params := map[string]interface{}{
		"query":                    "climate change",
		"search_recency_filter":    "week",
		"return_citations":         true,
		"return_related_questions": true,
		"max_tokens":               float64(512),
		"temperature":              0.5,
	}

	result, err := client.Search(ctx, params, cfg)
	if err != nil {
		return fmt.Errorf("search with parameters failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	// Check if citations are included
	if !contains(result, "Citations:") && params["return_citations"].(bool) {
		fmt.Println("Warning: Citations requested but not found in response")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testDomainFiltering(ctx context.Context, client *perplexity.Client, cfg *config.Config) error {
	params := map[string]interface{}{
		"query":                 "machine learning",
		"search_domain_filter":  []string{"arxiv.org", "nature.com"},
		"search_exclude_domains": []string{"wikipedia.org"},
	}

	result, err := client.Search(ctx, params, cfg)
	if err != nil {
		return fmt.Errorf("domain filtering search failed: %w", err)
	}

	if result == "" {
		return fmt.Errorf("empty result")
	}

	fmt.Printf("Result preview: %.100s...\n", result)
	return nil
}

func testErrorHandling(ctx context.Context, client *perplexity.Client, cfg *config.Config) error {
	// Test with empty query
	params := map[string]interface{}{
		"query": "",
	}

	_, err := client.Search(ctx, params, cfg)
	if err == nil {
		return fmt.Errorf("expected error for empty query, got nil")
	}

	fmt.Printf("Expected error received: %v\n", err)

	// Test with invalid model (this will be caught by API)
	params = map[string]interface{}{
		"query": "test",
		"model": "invalid-model-name",
	}

	_, err = client.Search(ctx, params, cfg)
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