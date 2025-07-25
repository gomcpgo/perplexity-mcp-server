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
				Description: "Perform a general web search using Perplexity AI",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The search query"
						},
						"model": {
							"type": "string",
							"description": "Model to use (sonar, sonar-pro, sonar-reasoning, sonar-finance)",
							"enum": ["sonar", "sonar-pro", "sonar-reasoning", "sonar-finance"]
						},
						"search_domain_filter": {
							"type": "array",
							"items": {"type": "string"},
							"description": "List of domains to include in search"
						},
						"search_exclude_domains": {
							"type": "array",
							"items": {"type": "string"},
							"description": "List of domains to exclude from search"
						},
						"search_recency_filter": {
							"type": "string",
							"description": "Time-based filter (hour, day, week, month, year)",
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
				Description: "Search academic papers and scholarly content",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The academic search query"
						},
						"subject_area": {
							"type": "string",
							"description": "Academic subject area (e.g., Physics, Computer Science)"
						},
						"model": {
							"type": "string",
							"description": "Model to use (defaults to sonar-reasoning)",
							"enum": ["sonar", "sonar-pro", "sonar-reasoning"]
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
				Description: "Search financial data, SEC filings, and market information",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The financial search query"
						},
						"ticker": {
							"type": "string",
							"description": "Stock ticker symbol (e.g., AAPL)"
						},
						"company_name": {
							"type": "string",
							"description": "Company name"
						},
						"report_type": {
							"type": "string",
							"description": "Type of financial report (e.g., 10-K, 10-Q, 8-K)"
						},
						"model": {
							"type": "string",
							"description": "Model to use (defaults to sonar-finance)",
							"enum": ["sonar-finance", "sonar-reasoning", "sonar-pro"]
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
				Description: "Advanced search with comprehensive filtering options",
				InputSchema: json.RawMessage(`{
					"type": "object",
					"properties": {
						"query": {
							"type": "string",
							"description": "The search query"
						},
						"model": {
							"type": "string",
							"description": "Model to use (defaults to sonar-pro)",
							"enum": ["sonar", "sonar-pro", "sonar-reasoning"]
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