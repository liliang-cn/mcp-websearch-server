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
		fmt.Println("  - web_search: Basic search returning titles and URLs")
		fmt.Println("  - web_search_with_content: Search with page content extraction")
		fmt.Println("  - deep_web_search: Multi-engine comprehensive search")
		fmt.Println("\nSearch Engines:")
		fmt.Println("  - Bing (primary)")
		fmt.Println("  - Brave (fallback)")
		fmt.Println("  - DuckDuckGo (fallback)")
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
