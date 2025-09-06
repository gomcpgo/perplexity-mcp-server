package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gomcpgo/mcp/pkg/handler"
	"github.com/gomcpgo/mcp/pkg/protocol"
	"github.com/gomcpgo/mcp/pkg/server"
	"github.com/prasanthmj/perplexity/pkg/config"
	mcpHandler "github.com/prasanthmj/perplexity/pkg/handler"
	"github.com/prasanthmj/perplexity/pkg/search"
	"github.com/prasanthmj/perplexity/test"
)

func main() {
	// Parse command line flags
	var (
		testMode        = flag.Bool("test", false, "Run integration tests")
		searchQuery     = flag.String("search", "", "Test general search: ./perplexity -search 'query'")
		academicQuery   = flag.String("academic", "", "Test academic search: ./perplexity -academic 'query'")
		financialQuery  = flag.String("financial", "", "Test financial search: ./perplexity -financial 'query'")
		filteredQuery   = flag.String("filtered", "", "Test filtered search: ./perplexity -filtered 'query'")
		listPrevious    = flag.Bool("list", false, "List previous cached queries")
		getResult       = flag.String("get", "", "Get cached result by ID: ./perplexity -get 'ABC123XYZ0'")
		model           = flag.String("model", "", "Model to use (sonar, sonar-pro)")
		debugMode       = flag.Bool("debug", false, "Enable debug mode")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Terminal mode operations for testing
	if *searchQuery != "" || *academicQuery != "" || *financialQuery != "" || *filteredQuery != "" || *listPrevious || *getResult != "" {
		err := runTerminalMode(cfg, *searchQuery, *academicQuery, *financialQuery, *filteredQuery, *listPrevious, *getResult, *model, *debugMode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Run integration tests if requested
	if *testMode {
		test.RunIntegrationTests()
		os.Exit(0)
	}

	// MCP Server mode (default)
	err = runMCPServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// runTerminalMode executes terminal mode for CLI testing
func runTerminalMode(cfg *config.Config, searchQuery, academicQuery, financialQuery, filteredQuery string, listPrevious bool, getResult, model string, debugMode bool) error {
	ctx := context.Background()

	// Create searcher for direct testing
	searcher, err := search.NewSearcher(cfg)
	if err != nil {
		return fmt.Errorf("failed to create searcher: %w", err)
	}

	// Handle list previous queries
	if listPrevious {
		result, err := searcher.ListPrevious(ctx)
		if err != nil {
			return fmt.Errorf("failed to list previous queries: %w", err)
		}
		fmt.Println(result)
		return nil
	}

	// Handle get previous result
	if getResult != "" {
		result, err := searcher.GetPreviousResult(ctx, getResult)
		if err != nil {
			return fmt.Errorf("failed to get previous result: %w", err)
		}
		fmt.Println(result)
		return nil
	}

	// Create search parameters
	var params *search.SearchParams
	var searchType string

	if searchQuery != "" {
		searchType = "general"
		params = &search.SearchParams{
			Query:      searchQuery,
			SearchType: searchType,
			Model:      model,
		}
	} else if academicQuery != "" {
		searchType = "academic"
		params = &search.SearchParams{
			Query:      academicQuery,
			SearchType: searchType,
			Model:      model,
		}
	} else if financialQuery != "" {
		searchType = "financial"
		params = &search.SearchParams{
			Query:      financialQuery,
			SearchType: searchType,
			Model:      model,
		}
	} else if filteredQuery != "" {
		searchType = "filtered"
		params = &search.SearchParams{
			Query:      filteredQuery,
			SearchType: searchType,
			Model:      model,
		}
	}

	if params == nil {
		return fmt.Errorf("no query provided")
	}

	// Execute search based on type
	var result string
	switch searchType {
	case "general":
		result, err = searcher.Search(ctx, params)
	case "academic":
		result, err = searcher.AcademicSearch(ctx, params)
	case "financial":
		result, err = searcher.FinancialSearch(ctx, params)
	case "filtered":
		result, err = searcher.FilteredSearch(ctx, params)
	}

	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	fmt.Println(result)
	return nil
}

// runMCPServer starts the MCP server
func runMCPServer(cfg *config.Config) error {
	// Create handler
	h, err := mcpHandler.NewHandler(cfg, false)
	if err != nil {
		return fmt.Errorf("failed to create handler: %w", err)
	}

	// Create MCP server
	registry := handler.NewHandlerRegistry()
	registry.RegisterToolHandler(h)

	srv := server.New(server.Options{
		Name:     "perplexity",
		Version:  "2.1.0",
		Registry: registry,
	})

	return srv.Run()
}

// PerplexityMCPServer wraps the handler to implement the required interfaces
type PerplexityMCPServer struct {
	handler *mcpHandler.Handler
}

// NewPerplexityMCPServer creates a new MCP server wrapper (for backward compatibility)
func NewPerplexityMCPServer(cfg *config.Config) (*PerplexityMCPServer, error) {
	h, err := mcpHandler.NewHandler(cfg, false)
	if err != nil {
		return nil, err
	}

	return &PerplexityMCPServer{
		handler: h,
	}, nil
}

// ListTools implements the ListTools interface
func (s *PerplexityMCPServer) ListTools(ctx context.Context) (*protocol.ListToolsResponse, error) {
	return s.handler.ListTools(ctx)
}

// CallTool implements the CallTool interface
func (s *PerplexityMCPServer) CallTool(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResponse, error) {
	return s.handler.CallTool(ctx, req)
}