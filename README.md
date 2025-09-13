# MCP Web Search Server

A Model Context Protocol (MCP) server that provides multi-engine web search capabilities with intelligent content extraction using a hybrid approach.

## Features

- 🔍 **Hybrid Search Engine**: Fast goquery-based search results + intelligent chromedp content extraction
- 🌐 **Multi-Engine Support**: Bing, Brave, and DuckDuckGo with smart fallback mechanisms
- 📄 **Intelligent Content Extraction**: Advanced article parsing with multiple content selectors
- 🚀 **Concurrent Processing**: Parallel content extraction with rate limiting
- 🤖 **AI-Ready Summaries**: Aggregated content optimized for AI analysis and summarization
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

### 🔍 `websearch_basic`
Basic web search returning titles, URLs and snippets from a single search engine using the hybrid approach.

**Parameters:**
- `query` (string, required): The search query
- `max_results` (int, optional): Maximum results to return (default: 10)

### 📄 `websearch_with_content`
Web search with intelligent content extraction from result pages using chromedp.

**Parameters:**
- `query` (string, required): The search query
- `max_results` (int, optional): Maximum results to return (default: 5)
- `extract_content` (bool, optional): Extract full page content (default: true)

### 🚀 `websearch_multi_engine`
Comprehensive search across multiple engines (Bing, Brave, DuckDuckGo) with content extraction.

**Parameters:**
- `query` (string, required): The search query
- `max_results` (int, optional): Maximum results to return (default: 3)
- `engines` (array, optional): Search engines to use ["bing", "brave", "duckduckgo"] (default: all)

### 🤖 `websearch_ai_summary`
Search and return AI-ready aggregated content optimized for analysis and summarization.

**Parameters:**
- `query` (string, required): The search query
- `max_results` (int, optional): Maximum results to return (default: 3)

**Returns:** Formatted markdown content with proper structure for AI processing.

## Architecture

```
mcp-websearch-server/
├── main.go                     # Entry point with CLI flags
├── mcp/                        # MCP protocol implementation
│   └── server.go              # MCP server and tool registration
├── search/                     # Search engine implementations
│   ├── interface.go           # Common interfaces
│   ├── hybrid_searcher.go     # Hybrid multi-engine searcher
│   ├── multi_engine.go        # Basic multi-engine orchestration
│   ├── bing_goquery.go        # Fast Bing search with goquery
│   ├── brave_goquery.go       # Fast Brave search with goquery
│   ├── duckduckgo_goquery.go  # Fast DuckDuckGo search with goquery
│   ├── bing.go               # Original Bing search (chromedp)
│   ├── brave.go              # Original Brave search (chromedp)
│   └── duckduckgo.go         # Original DuckDuckGo search (chromedp)
├── extraction/                 # Content extraction
│   ├── hybrid_extractor.go   # Intelligent chromedp-based extraction
│   └── chromedp.go           # Basic browser-based extraction
├── examples/                   # Demo applications
│   ├── basic_search_demo/     # Basic search functionality demo
│   ├── hybrid_search_demo/    # Hybrid search with content extraction
│   └── mcp_tools_demo/        # MCP server tools demonstration
└── utils/                     # Utilities
    └── retry.go              # Retry logic with backoff
```

## Hybrid Approach

The server uses a sophisticated hybrid approach for optimal performance:

### 1. Fast Search Results (goquery)
- **Bing**: Scrapes `www.bing.com/search` with proper CSS selectors
- **Brave**: Scrapes `search.brave.com/search` for results
- **DuckDuckGo**: Scrapes `duckduckgo.com` with lite interface
- **Benefits**: Fast response times, reliable result parsing

### 2. Intelligent Content Extraction (chromedp)
- **Article Detection**: Uses advanced selectors to find main content
- **Content Cleaning**: Removes scripts, styles, and navigation elements
- **Fallback Strategy**: Falls back to paragraph extraction if article content not found
- **Benefits**: High-quality content extraction, JavaScript handling

### 3. AI-Ready Aggregation
- **Structured Output**: Properly formatted markdown for AI processing
- **Content Summarization**: Truncates content intelligently at sentence boundaries
- **Multi-Source**: Combines content from multiple search engines
- **Benefits**: Optimized for AI analysis and summarization

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

### Example Applications

```bash
# Test basic search functionality
go run ./examples/basic_search_demo/main.go

# Test hybrid search with content extraction
go run ./examples/hybrid_search_demo/main.go

# Test MCP server tools
go run ./examples/mcp_tools_demo/main.go
```

## How It Works

1. **Search Request**: Receives search query via MCP protocol
2. **Engine Selection**: Uses goquery-based engines for fast results
3. **Search Execution**: Performs HTTP-based search with proper headers
4. **Content Extraction**: Uses chromedp for intelligent content extraction
5. **Aggregation**: Combines and formats content for AI analysis
6. **Response**: Returns structured results via MCP protocol

## Search Engine Priority

The hybrid searcher prioritizes engines in this order:
1. **DuckDuckGo** - Primary engine (privacy-focused)
2. **Bing** - First fallback (comprehensive results)
3. **Brave** - Second fallback (independent search)

If one engine fails, the server automatically tries the next available engine.

## Error Handling

- Implements retry logic with exponential backoff
- Graceful fallback to alternative search engines
- Structured error messages via MCP protocol
- Timeout handling for long-running operations
- Rate limiting for content extraction

## Performance

- **Search Speed**: ~200-500ms per search using goquery
- **Content Extraction**: ~2-5s per page using chromedp
- **Concurrent Extraction**: Limited to 2-3 simultaneous browser instances
- **Memory Usage**: Optimized with proper context cleanup

## Dependencies

- **MCP Go SDK**: Model Context Protocol implementation
- **chromedp**: Browser automation for content extraction
- **goquery**: Fast HTML parsing and scraping
- **Standard Library**: HTTP client, context, sync primitives

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Built with [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- Uses [chromedp](https://github.com/chromedp/chromedp) for browser automation
- Uses [goquery](https://github.com/PuerkitoBio/goquery) for HTML parsing
- Implements [Model Context Protocol](https://modelcontextprotocol.io) specification