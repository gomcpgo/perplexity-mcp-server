package handler

import (
	"context"
	"fmt"

	"github.com/gomcpgo/mcp/pkg/protocol"
	"github.com/prasanthmj/perplexity/pkg/config"
	"github.com/prasanthmj/perplexity/pkg/search"
)

// Handler handles MCP protocol operations
type Handler struct {
	searcher *search.Searcher
	config   *config.Config
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, debugMode bool) (*Handler, error) {
	searcher, err := search.NewSearcher(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create searcher: %w", err)
	}

	return &Handler{
		searcher: searcher,
		config:   cfg,
	}, nil
}

// CallTool handles MCP tool calls
func (h *Handler) CallTool(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResponse, error) {
	var result string
	var err error

	switch req.Name {
	case "perplexity_search":
		result, err = h.handlePerplexitySearch(ctx, req.Arguments)
	case "perplexity_academic_search":
		result, err = h.handleAcademicSearch(ctx, req.Arguments)
	case "perplexity_financial_search":
		result, err = h.handleFinancialSearch(ctx, req.Arguments)
	case "perplexity_filtered_search":
		result, err = h.handleFilteredSearch(ctx, req.Arguments)
	case "list_previous":
		result, err = h.handleListPrevious(ctx, req.Arguments)
	case "get_previous_result":
		result, err = h.handleGetPreviousResult(ctx, req.Arguments)
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