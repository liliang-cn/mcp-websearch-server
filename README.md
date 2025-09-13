# MCP Web Search Server

A Model Context Protocol (MCP) server that provides multi-engine web search capabilities with content extraction.

## Features

- 🔍 **Multi-Engine Search**: Prioritizes Bing → Brave → DuckDuckGo for optimal reliability
- 📄 **Content Extraction**: Fetches and extracts full page content from search results
- 🚀 **Concurrent Processing**: Extracts content from multiple pages simultaneously
- 🔄 **Smart Fallback**: Automatically switches to alternative search engines on failure
- 🛠️ **MCP Protocol**: Full compliance with Model Context Protocol specification

## Installation

### Via `go install`

```bash
go install github.com/liliang-cn/mcp-websearch-server@latest
```

### From Source

```bash
git clone https://github.com/liliang-cn/mcp-websearch-server
cd mcp-websearch-server
go build -o mcp-websearch-server
```

## Usage

### Standalone

```bash
# Show help
mcp-websearch-server --help

# Run the server (stdio mode)
mcp-websearch-server
```

### Integration with Claude Desktop

Add to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`  
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "websearch": {
      "command": "mcp-websearch-server"
    }
  }
}
```

If installed via `go install`, make sure `~/go/bin` is in your PATH.

## Available Tools

### 🔍 `web_search`
Basic web search returning titles and URLs.

**Parameters:**
- `query` (string, required): The search query
- `max_results` (int, optional): Maximum results to return (default: 10)

### 📄 `web_search_with_content`
Search with automatic content extraction from result pages.

**Parameters:**
- `query` (string, required): The search query
- `max_results` (int, optional): Maximum results to return (default: 5)
- `extract_content` (bool, optional): Extract full page content (default: true)

### 🚀 `deep_web_search`
Comprehensive search across multiple engines with content extraction.

**Parameters:**
- `query` (string, required): The search query
- `max_results` (int, optional): Maximum results to return (default: 3)
- `engines` (array, optional): Search engines to use ["bing", "brave", "duckduckgo"] (default: all)

## Architecture

```
mcp-websearch-server/
├── main.go                 # Entry point with CLI flags
├── mcp/                    # MCP protocol implementation
│   └── server.go          # MCP server and tool registration
├── search/                 # Search engine implementations
│   ├── interface.go       # Common interfaces
│   ├── multi_engine.go    # Multi-engine orchestration
│   ├── bing.go           # Bing search
│   ├── brave.go          # Brave search
│   └── duckduckgo.go     # DuckDuckGo search
├── extraction/            # Content extraction
│   └── chromedp.go       # Browser-based extraction
└── utils/                 # Utilities
    └── retry.go          # Retry logic with backoff
```

## Development

### Prerequisites

- Go 1.21 or higher
- Chrome/Chromium browser (for content extraction)

### Building

```bash
# Build the server
go build -o mcp-websearch-server

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run
```

### Testing

The project includes comprehensive unit tests with 60%+ coverage:

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## How It Works

1. **Search Request**: Receives search query via MCP protocol
2. **Engine Selection**: Chooses primary engine (Bing) or fallback
3. **Search Execution**: Performs search using browser automation
4. **Content Extraction**: Optionally extracts full page content
5. **Response**: Returns structured results via MCP protocol

## Search Engine Priority

The server prioritizes search engines in this order:
1. **Bing** - Primary engine
2. **Brave** - First fallback
3. **DuckDuckGo** - Second fallback

If one engine fails, the server automatically tries the next available engine.

## Error Handling

- Implements retry logic with exponential backoff
- Graceful fallback to alternative search engines
- Structured error messages via MCP protocol
- Timeout handling for long-running operations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Built with [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- Uses [chromedp](https://github.com/chromedp/chromedp) for browser automation
- Implements [Model Context Protocol](https://modelcontextprotocol.io) specification