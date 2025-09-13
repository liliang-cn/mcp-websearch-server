package mcp

import (
	"context"
	"fmt"

	"github.com/liliang-cn/mcp-websearch-server/search"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	mcpServer *mcp.Server
	searcher  search.MultiEngineSearcher
}

func NewServer() (*Server, error) {
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-websearch-server",
			Version: "1.0.0",
		},
		nil,
	)

	s := &Server{
		mcpServer: mcpServer,
		searcher:  search.NewMultiEngineSearcher(),
	}

	if err := s.registerTools(); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

func (s *Server) registerTools() error {
	type basicSearchArgs struct {
		Query      string `json:"query" jsonschema:"the search query to execute"`
		MaxResults int    `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "web_search",
		Description: "Perform a basic web search returning titles and URLs",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args basicSearchArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 {
			args.MaxResults = 10
		}

		results, err := s.searcher.Search(ctx, args.Query, search.SearchOptions{
			MaxResults:     args.MaxResults,
			ExtractContent: false,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("search failed: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Found %d results", len(results))},
			},
		}, results, nil
	})

	type searchWithContentArgs struct {
		Query          string `json:"query" jsonschema:"the search query to execute"`
		MaxResults     int    `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
		ExtractContent bool   `json:"extract_content,omitempty" jsonschema:"whether to extract full page content"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "web_search_with_content",
		Description: "Search the web and extract content from result pages",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchWithContentArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 {
			args.MaxResults = 5
		}
		if !args.ExtractContent {
			args.ExtractContent = true
		}

		results, err := s.searcher.Search(ctx, args.Query, search.SearchOptions{
			MaxResults:     args.MaxResults,
			ExtractContent: args.ExtractContent,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("search with content failed: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Found %d results with content", len(results))},
			},
		}, results, nil
	})

	type deepSearchArgs struct {
		Query      string   `json:"query" jsonschema:"the search query to execute"`
		MaxResults int      `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
		Engines    []string `json:"engines,omitempty" jsonschema:"search engines to use (bing, brave, duckduckgo)"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "deep_web_search",
		Description: "Comprehensive search with full page analysis across multiple engines",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deepSearchArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 {
			args.MaxResults = 3
		}
		if len(args.Engines) == 0 {
			args.Engines = []string{"bing", "brave", "duckduckgo"}
		}

		results, err := s.searcher.DeepSearch(ctx, args.Query, search.SearchOptions{
			MaxResults:     args.MaxResults,
			ExtractContent: true,
			Engines:        args.Engines,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("deep search failed: %w", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Deep search found %d results from %d engines", len(results), len(args.Engines))},
			},
		}, results, nil
	})

	return nil
}
