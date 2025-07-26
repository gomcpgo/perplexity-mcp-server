package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gomcpgo/mcp/pkg/handler"
	"github.com/gomcpgo/mcp/pkg/protocol"
	"github.com/gomcpgo/mcp/pkg/server"
	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/perplexity"
	"github.com/prasanthmj/perplexity/test"
)

type PerplexityMCPServer struct {
	client *perplexity.Client
	config *config.Config
}

func NewPerplexityMCPServer(cfg *config.Config) *PerplexityMCPServer {
	return &PerplexityMCPServer{
		client: perplexity.NewClient(cfg.APIKey, cfg.Timeout),
		config: cfg,
	}
}

func (s *PerplexityMCPServer) ListTools(ctx context.Context) (*protocol.ListToolsResponse, error) {
	return &protocol.ListToolsResponse{
		Tools: []protocol.Tool{
			{
				Name:        "perplexity_search",
				Description: "General web search with real-time information. Best for: current events, general knowledge, quick facts, web content. Use 'sonar' model for quick searches, 'sonar-pro' for comprehensive results.",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The search query. Be specific and clear for best results."
						},
						"model": {
							"type": "string",
							"description": "Choose 'sonar' for quick factual searches (faster, cheaper) or 'sonar-pro' for comprehensive searches (better depth, more thorough)",
							"enum": ["sonar", "sonar-pro"],
							"default": "sonar"
						},
						"search_domain_filter": {
							"type": "array",
							"items": {"type": "string"},
							"description": "Limit search to specific domains (e.g., ['wikipedia.org', 'nature.com'])"
						},
						"search_exclude_domains": {
							"type": "array",
							"items": {"type": "string"},
							"description": "Exclude specific domains from results (e.g., ['reddit.com', 'quora.com'])"
						},
						"search_recency_filter": {
							"type": "string",
							"description": "Filter by recency: 'hour' for breaking news, 'day' for today's updates, 'week' for recent events, 'month' for recent trends, 'year' for current year",
							"enum": ["hour", "day", "week", "month", "year"]
						},
						"return_citations": {
							"type": "boolean",
							"description": "Include citations in response"
						},
						"return_images": {
							"type": "boolean",
							"description": "Include images in response"
						},
						"return_related_questions": {
							"type": "boolean",
							"description": "Include related questions"
						},
						"max_tokens": {
							"type": "number",
							"description": "Maximum tokens in response"
						},
						"temperature": {
							"type": "number",
							"description": "Response randomness (0-2)"
						},
						"date_range_start": {
							"type": "string",
							"description": "Start date for filtering (YYYY-MM-DD)"
						},
						"date_range_end": {
							"type": "string",
							"description": "End date for filtering (YYYY-MM-DD)"
						},
						"location": {
							"type": "string",
							"description": "Location for geo-specific search"
						}
					},
					"required": ["query"]
				}`),
			},
			{
				Name:        "perplexity_academic_search",
				Description: "Search academic papers, research articles, and scholarly content. Automatically filters to academic sources (arxiv.org, pubmed, journals). Best for: research papers, scientific studies, academic citations.",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The academic search query. Include key terms, authors, or specific topics."
						},
						"subject_area": {
							"type": "string",
							"description": "Optional: Specify academic field to narrow results (e.g., 'Physics', 'Computer Science', 'Medicine')"
						},
						"model": {
							"type": "string",
							"description": "Defaults to 'sonar-pro' for comprehensive academic results. Use 'sonar' only for quick lookups.",
							"enum": ["sonar", "sonar-pro"],
							"default": "sonar-pro"
						},
						"search_domain_filter": {
							"type": "array",
							"items": {"type": "string"},
							"description": "List of academic domains to include"
						},
						"search_recency_filter": {
							"type": "string",
							"description": "Time-based filter",
							"enum": ["hour", "day", "week", "month", "year"]
						},
						"return_citations": {
							"type": "boolean",
							"description": "Include citations (default: true)"
						},
						"max_tokens": {
							"type": "number",
							"description": "Maximum tokens in response"
						},
						"temperature": {
							"type": "number",
							"description": "Response randomness (0-2)"
						}
					},
					"required": ["query"]
				}`),
			},
			{
				Name:        "perplexity_financial_search",
				Description: "Search financial data, SEC filings, earnings reports, and market information. Optimized for financial domains and recent data. Best for: stock analysis, earnings, SEC filings, market trends.",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The financial search query. Include company names, tickers, or specific financial metrics."
						},
						"ticker": {
							"type": "string",
							"description": "Optional: Stock ticker symbol (e.g., 'AAPL', 'MSFT') to focus search"
						},
						"company_name": {
							"type": "string",
							"description": "Optional: Company name to ensure accurate results"
						},
						"report_type": {
							"type": "string",
							"description": "Optional: SEC report type (e.g., '10-K' for annual, '10-Q' for quarterly, '8-K' for current)"
						},
						"model": {
							"type": "string",
							"description": "Defaults to 'sonar-pro' for comprehensive financial data. Use 'sonar' for quick stock quotes.",
							"enum": ["sonar", "sonar-pro"],
							"default": "sonar-pro"
						},
						"search_recency_filter": {
							"type": "string",
							"description": "Time-based filter",
							"enum": ["hour", "day", "week", "month", "year"]
						},
						"date_range_start": {
							"type": "string",
							"description": "Start date for reports (YYYY-MM-DD)"
						},
						"date_range_end": {
							"type": "string",
							"description": "End date for reports (YYYY-MM-DD)"
						},
						"return_citations": {
							"type": "boolean",
							"description": "Include citations (default: true)"
						},
						"max_tokens": {
							"type": "number",
							"description": "Maximum tokens in response"
						}
					},
					"required": ["query"]
				}`),
			},
			{
				Name:        "perplexity_filtered_search",
				Description: "Advanced search with multiple filters. Best for: specific requirements, domain-specific searches, content type filtering, location-based searches. Use when other specialized searches don't fit your needs.",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The search query"
						},
						"model": {
							"type": "string",
							"description": "Choose based on needs: 'sonar' for quick filtered searches, 'sonar-pro' for comprehensive filtered results",
							"enum": ["sonar", "sonar-pro"],
							"default": "sonar-pro"
						},
						"search_domain_filter": {
							"type": "array",
							"items": {"type": "string"},
							"description": "List of domains to include"
						},
						"search_exclude_domains": {
							"type": "array",
							"items": {"type": "string"},
							"description": "List of domains to exclude"
						},
						"search_recency_filter": {
							"type": "string",
							"description": "Time-based filter",
							"enum": ["hour", "day", "week", "month", "year"]
						},
						"content_type": {
							"type": "string",
							"description": "Type of content (news, academic, blog, etc.)"
						},
						"file_type": {
							"type": "string",
							"description": "File type filter (pdf, doc, html, etc.)"
						},
						"language": {
							"type": "string",
							"description": "Language filter"
						},
						"country": {
							"type": "string",
							"description": "Country for geo-specific search"
						},
						"date_range_start": {
							"type": "string",
							"description": "Start date (YYYY-MM-DD)"
						},
						"date_range_end": {
							"type": "string",
							"description": "End date (YYYY-MM-DD)"
						},
						"return_citations": {
							"type": "boolean",
							"description": "Include citations"
						},
						"return_images": {
							"type": "boolean",
							"description": "Include images"
						},
						"return_related_questions": {
							"type": "boolean",
							"description": "Include related questions"
						},
						"max_tokens": {
							"type": "number",
							"description": "Maximum tokens in response"
						},
						"temperature": {
							"type": "number",
							"description": "Response randomness (0-2)"
						},
						"custom_filters": {
							"type": "object",
							"description": "Additional custom filters as key-value pairs"
						}
					},
					"required": ["query"]
				}`),
			},
		},
	}, nil
}

func (s *PerplexityMCPServer) CallTool(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResponse, error) {
	var result string
	var err error

	switch req.Name {
	case "perplexity_search":
		result, err = s.client.Search(ctx, req.Arguments, s.config)
	case "perplexity_academic_search":
		result, err = s.client.AcademicSearch(ctx, req.Arguments, s.config)
	case "perplexity_financial_search":
		result, err = s.client.FinancialSearch(ctx, req.Arguments, s.config)
	case "perplexity_filtered_search":
		result, err = s.client.FilteredSearch(ctx, req.Arguments, s.config)
	default:
		return nil, fmt.Errorf("unknown tool: %s", req.Name)
	}

	if err != nil {
		return nil, err
	}

	return &protocol.CallToolResponse{
		Content: []protocol.ToolContent{
			{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

func main() {
	// Parse command line flags
	testMode := flag.Bool("test", false, "Run integration tests")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Run tests if requested
	if *testMode {
		test.RunIntegrationTests()
		os.Exit(0)
	}

	// Create and run MCP server
	perplexityServer := NewPerplexityMCPServer(cfg)

	registry := handler.NewHandlerRegistry()
	registry.RegisterToolHandler(perplexityServer)

	srv := server.New(server.Options{
		Name:     "perplexity",
		Version:  "2.0.0",
		Registry: registry,
	})

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}