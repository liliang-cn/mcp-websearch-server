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
		searcher:  search.NewHybridSearcher(),
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
		Name:        "websearch_basic",
		Description: "Basic web search returning titles, URLs and snippets from a single search engine",
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

		// Convert results to formatted text content
		var content string
		for i, result := range results {
			content += fmt.Sprintf("### Result %d\n", i+1)
			content += fmt.Sprintf("**Title:** %s\n", result.Title)
			content += fmt.Sprintf("**URL:** %s\n", result.URL)
			content += fmt.Sprintf("**Snippet:** %s\n", result.Snippet)
			content += fmt.Sprintf("**Engine:** %s\n\n", result.Engine)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})

	type searchWithContentArgs struct {
		Query          string `json:"query" jsonschema:"the search query to execute"`
		MaxResults     int    `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
		ExtractContent bool   `json:"extract_content,omitempty" jsonschema:"whether to extract full page content"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "websearch_with_content",
		Description: "Web search with intelligent content extraction from result pages",
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

		// Convert results to formatted text content with extracted content
		var content string
		for i, result := range results {
			content += fmt.Sprintf("### Result %d\n", i+1)
			content += fmt.Sprintf("**Title:** %s\n", result.Title)
			content += fmt.Sprintf("**URL:** %s\n", result.URL)
			content += fmt.Sprintf("**Snippet:** %s\n", result.Snippet)
			content += fmt.Sprintf("**Engine:** %s\n", result.Engine)
			
			if result.Content != "" {
				// Truncate content if too long
				extractedContent := result.Content
				if len(extractedContent) > 1000 {
					extractedContent = extractedContent[:1000] + "..."
				}
				content += fmt.Sprintf("\n**Extracted Content:**\n%s\n", extractedContent)
			}
			content += "\n---\n\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})

	type deepSearchArgs struct {
		Query      string   `json:"query" jsonschema:"the search query to execute"`
		MaxResults int      `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
		Engines    []string `json:"engines,omitempty" jsonschema:"search engines to use (bing, brave, duckduckgo)"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "websearch_multi_engine",
		Description: "Comprehensive search across multiple engines (Bing, Brave, DuckDuckGo) with content extraction",
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

		// Convert results to formatted text content with full extraction
		var content string
		content += fmt.Sprintf("## Deep Search Results (%d results from %d engines)\n\n", len(results), len(args.Engines))
		
		for i, result := range results {
			content += fmt.Sprintf("### Result %d\n", i+1)
			content += fmt.Sprintf("**Title:** %s\n", result.Title)
			content += fmt.Sprintf("**URL:** %s\n", result.URL)
			content += fmt.Sprintf("**Snippet:** %s\n", result.Snippet)
			content += fmt.Sprintf("**Engine:** %s\n", result.Engine)
			
			if result.Content != "" {
				// Truncate content if too long
				extractedContent := result.Content
				if len(extractedContent) > 1500 {
					extractedContent = extractedContent[:1500] + "..."
				}
				content += fmt.Sprintf("\n**Extracted Content:**\n%s\n", extractedContent)
			}
			content += "\n---\n\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})

	type searchAndAggregateArgs struct {
		Query      string `json:"query" jsonschema:"the search query to execute"`
		MaxResults int    `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "websearch_ai_summary",
		Description: "Search and return AI-ready aggregated content optimized for analysis and summarization",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchAndAggregateArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 {
			args.MaxResults = 3
		}

		// Check if we have a hybrid searcher with aggregation capability
		if hybridSearcher, ok := s.searcher.(*search.HybridMultiEngineSearcher); ok {
			aggregatedContent, err := hybridSearcher.SearchAndAggregate(ctx, args.Query, args.MaxResults)
			if err != nil {
				return nil, nil, fmt.Errorf("search and aggregate failed: %w", err)
			}

			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: aggregatedContent},
				},
			}, nil, nil
		}

		// Fallback to regular search with content extraction
		results, err := s.searcher.Search(ctx, args.Query, search.SearchOptions{
			MaxResults:     args.MaxResults,
			ExtractContent: true,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("search failed: %w", err)
		}

		// Format as aggregated content manually
		var content string
		content += fmt.Sprintf("# Search Results for: %s\n\n", args.Query)
		
		for i, result := range results {
			content += fmt.Sprintf("## %d. %s\n", i+1, result.Title)
			content += fmt.Sprintf("**Source:** %s\n", result.URL)
			content += fmt.Sprintf("**Engine:** %s\n\n", result.Engine)
			
			if result.Content != "" {
				extractedContent := result.Content
				if len(extractedContent) > 1500 {
					extractedContent = extractedContent[:1500] + "..."
				}
				content += extractedContent
			} else if result.Snippet != "" {
				content += result.Snippet
			}
			
			content += "\n\n---\n\n"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: content},
			},
		}, nil, nil
	})

	return nil
}
