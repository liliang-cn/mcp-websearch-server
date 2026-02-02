package main

import (
	"context"
	"fmt"
	"time"

	"github.com/liliang-cn/mcp-websearch-server/search"
)

func main() {
	fmt.Println("=== MCP Web Search Server Test ===")
	
	searcher := search.NewMultiEngineSearcher()
	ctx := context.Background()
	
	queries := []string{"Trump", "China", "iPhone 17"}
	
	for _, query := range queries {
		fmt.Printf("Testing search for: '%s'\n", query)
		fmt.Println("----------------------------------------")
		
		// Test with timeout
		searchCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		
		results, err := searcher.Search(searchCtx, query, search.SearchOptions{
			MaxResults:     3,
			ExtractContent: false,
			Timeout:        20 * time.Second,
		})
		
		cancel()
		
		if err != nil {
			fmt.Printf("❌ Error searching for '%s': %v\n", query, err)
		} else {
			fmt.Printf("✅ Found %d results for '%s':\n", len(results), query)
			for i, result := range results {
				fmt.Printf("\n  Result %d:\n", i+1)
				fmt.Printf("    Title: %s\n", result.Title)
				fmt.Printf("    URL: %s\n", result.URL)
				fmt.Printf("    Engine: %s\n", result.Engine)
				if result.Snippet != "" {
					snippet := result.Snippet
					if len(snippet) > 100 {
						snippet = snippet[:100] + "..."
					}
					fmt.Printf("    Snippet: %s\n", snippet)
				}
			}
		}
		
		fmt.Println("\n========================================")
		
		// Small delay between searches
		time.Sleep(2 * time.Second)
	}
	
	fmt.Println("Test completed!")
}