package mcp

import (
	"context"
	"fmt"

	"github.com/liliang-cn/mcp-websearch-server/extraction"
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
	// ... (basicSearchArgs omitted for brevity, but I will write the full file)
	// I'll use replace for specific parts to be safer, but since I have the content, 
	// I'll just rewrite the file with all tools correctly.
	return s.doRegisterTools()
}

// I'll split the registration to keep it clean
func (s *Server) doRegisterTools() error {
	// websearch_basic
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
		results, err := s.searcher.Search(ctx, args.Query, search.SearchOptions{MaxResults: args.MaxResults})
		if err != nil {
			return nil, nil, err
		}
		var content string
		for i, result := range results {
			content += fmt.Sprintf("### Result %d\n**Title:** %s\n**URL:** %s\n**Snippet:** %s\n\n", i+1, result.Title, result.URL, result.Snippet)
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil, nil
	})

	// websearch_with_content
	type searchWithContentArgs struct {
		Query          string `json:"query" jsonschema:"the search query to execute"`
		MaxResults     int    `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
		ExtractContent bool   `json:"extract_content,omitempty" jsonschema:"whether to extract full page content"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "websearch_with_content",
		Description: "Web search with intelligent content extraction from result pages",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchWithContentArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 { args.MaxResults = 5 }
		results, err := s.searcher.Search(ctx, args.Query, search.SearchOptions{MaxResults: args.MaxResults, ExtractContent: true})
		if err != nil { return nil, nil, err }
		var content string
		for i, result := range results {
			content += fmt.Sprintf("### Result %d\n**Title:** %s\n**URL:** %s\n", i+1, result.Title, result.URL)
			if result.Content != "" {
				ext := result.Content
				if len(ext) > 1500 { ext = ext[:1500] + "..." }
				content += fmt.Sprintf("\n**Content:**\n%s\n", ext)
			}
			content += "\n---\n\n"
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil, nil
	})

	// websearch_multi_engine
	type deepSearchArgs struct {
		Query      string   `json:"query" jsonschema:"the search query to execute"`
		MaxResults int      `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
		Engines    []string `json:"engines,omitempty" jsonschema:"search engines to use"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "websearch_multi_engine",
		Description: "Comprehensive search across multiple engines with content extraction",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args deepSearchArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 { args.MaxResults = 10 }
		results, err := s.searcher.DeepSearch(ctx, args.Query, search.SearchOptions{MaxResults: args.MaxResults, Engines: args.Engines, ExtractContent: true})
		if err != nil { return nil, nil, err }
		var content string
		for i, result := range results {
			content += fmt.Sprintf("### Result %d\n**Title:** %s\n**URL:** %s\n", i+1, result.Title, result.URL)
			if result.Content != "" {
				ext := result.Content
				if len(ext) > 1500 { ext = ext[:1500] + "..." }
				content += fmt.Sprintf("\n**Content:**\n%s\n", ext)
			}
			content += "\n---\n\n"
		}
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil, nil
	})

	// websearch_ai_summary
	type searchAndAggregateArgs struct {
		Query      string `json:"query" jsonschema:"the search query to execute"`
		MaxResults int    `json:"max_results,omitempty" jsonschema:"maximum number of results to return"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "websearch_ai_summary",
		Description: "Search and return AI-ready aggregated content optimized for analysis and summarization",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args searchAndAggregateArgs) (*mcp.CallToolResult, any, error) {
		if args.MaxResults == 0 { args.MaxResults = 5 }
		if hs, ok := s.searcher.(*search.HybridMultiEngineSearcher); ok {
			aggregated, err := hs.SearchAndAggregate(ctx, args.Query, args.MaxResults)
			if err != nil { return nil, nil, err }
			return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: aggregated}}}, nil, nil
		}
		return nil, nil, fmt.Errorf("aggregation not supported")
	})

	// fetch_page_content
	type fetchPageContentArgs struct {
		URL string `json:"url" jsonschema:"the URL of the page to fetch content from"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "fetch_page_content",
		Description: "Directly fetch and extract the main content from a specific URL using Readability and Markdown conversion",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args fetchPageContentArgs) (*mcp.CallToolResult, any, error) {
		if args.URL == "" { return nil, nil, fmt.Errorf("URL is required") }
		content, err := extraction.NewHybridExtractor().ExtractContent(ctx, args.URL)
		if err != nil { return nil, nil, err }
		return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}}, nil, nil
	})

	// take_screenshot
	type takeScreenshotArgs struct {
		URL      string `json:"url" jsonschema:"the URL of the page to screenshot"`
		FullPage bool   `json:"full_page,omitempty" jsonschema:"whether to take a full page screenshot"`
	}

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "take_screenshot",
		Description: "Capture a screenshot of a webpage",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args takeScreenshotArgs) (*mcp.CallToolResult, any, error) {
		if args.URL == "" { return nil, nil, fmt.Errorf("URL is required") }
		imgData, err := extraction.NewChromedpExtractor().CaptureScreenshot(ctx, args.URL, args.FullPage)
		if err != nil { return nil, nil, err }
		
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.ImageContent{
					Data:     imgData,
					MIMEType: "image/png",
				},
				&mcp.TextContent{Text: fmt.Sprintf("Successfully captured screenshot of %s (%d bytes).", args.URL, len(imgData))},
			},
		}, nil, nil
	})

	return nil
}