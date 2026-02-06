package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/liliang-cn/mcp-websearch-server/mcp"
)

func main() {
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	if *help {
		fmt.Println("MCP Web Search Server")
		fmt.Println("\nUsage: mcp-websearch-server [options]")
		fmt.Println("\nOptions:")
		fmt.Println("  --help    Show this help message")
		fmt.Println("\nDescription:")
		fmt.Println("  This server provides web search capabilities via the Model Context Protocol (MCP).")
		fmt.Println("  It runs in stdio mode, reading MCP protocol messages from stdin and writing responses to stdout.")
		fmt.Println("\nAvailable Tools:")
		fmt.Println("  - websearch_basic: Basic search returning titles, URLs and snippets from a single engine")
		fmt.Println("  - websearch_with_content: Search with intelligent page content extraction")
		fmt.Println("  - websearch_multi_engine: Comprehensive multi-engine search with content extraction")
		fmt.Println("  - websearch_ai_summary: Aggregated content optimized for AI analysis")
		fmt.Println("  - fetch_page_content: Directly extract content from any URL")
		fmt.Println("\nSearch Engines:")
		fmt.Println("  - DuckDuckGo (primary)")
		fmt.Println("  - Bing (fallback)")
		fmt.Println("  - Brave (fallback)")
		fmt.Println("\nIntegration with Claude Desktop:")
		fmt.Println("  Add to ~/Library/Application Support/Claude/claude_desktop_config.json:")
		fmt.Println(`  {
    "mcpServers": {
      "websearch": {
        "command": "/path/to/mcp-websearch-server"
      }
    }
  }`)
		os.Exit(0)
	}

	ctx := context.Background()

	server, err := mcp.NewServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
