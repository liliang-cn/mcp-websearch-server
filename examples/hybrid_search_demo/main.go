package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/liliang-cn/mcp-websearch-server/search"
)

func main() {
	fmt.Println("=== Hybrid Search Test (goquery + chromedp) ===\n")
	
	// Create hybrid searcher
	searcher := search.NewHybridSearcher()
	ctx := context.Background()
	
	// Test 1: Basic search with DuckDuckGo
	fmt.Println("📰 Test 1: Search for Trump news and extract content")
	fmt.Println(strings.Repeat("=", 50))
	
	results, err := searcher.Search(ctx, "Trump latest news", search.SearchOptions{
		MaxResults:     3,
		ExtractContent: true,
		Timeout:        45 * time.Second,
	})
	
	if err != nil {
		log.Printf("Search error: %v\n", err)
	} else {
		for i, result := range results {
			fmt.Printf("\n📌 Result %d:\n", i+1)
			fmt.Printf("Title: %s\n", result.Title)
			fmt.Printf("URL: %s\n", result.URL)
			fmt.Printf("Engine: %s\n", result.Engine)
			
			if result.Content != "" {
				// Show first 300 chars of extracted content
				content := result.Content
				if len(content) > 300 {
					content = content[:300] + "..."
				}
				fmt.Printf("Extracted Content:\n%s\n", content)
			} else {
				fmt.Printf("Snippet: %s\n", result.Snippet)
			}
			fmt.Println()
		}
	}
	
	// Test 2: Aggregated search (ready for AI summarization)
	fmt.Println("\n\n🤖 Test 2: Aggregated Search (Ready for AI Summary)")
	fmt.Println(strings.Repeat("=", 50))
	
	if hybridSearcher, ok := searcher.(*search.HybridMultiEngineSearcher); ok {
		aggregated, err := hybridSearcher.SearchAndAggregate(ctx, "iPhone 17 features", 3)
		if err != nil {
			log.Printf("Aggregation error: %v\n", err)
		} else {
			// Show first 1000 chars of aggregated content
			if len(aggregated) > 1000 {
				aggregated = aggregated[:1000] + "\n\n[... truncated for display ...]"
			}
			fmt.Println("Aggregated Content (for AI processing):")
			fmt.Println(aggregated)
		}
	}
	
	// Test 3: Deep search across multiple engines
	fmt.Println("\n\n🌐 Test 3: Deep Search (Multiple Engines)")
	fmt.Println(strings.Repeat("=", 50))
	
	deepResults, err := searcher.DeepSearch(ctx, "China economy 2025", search.SearchOptions{
		MaxResults:     6,
		ExtractContent: true,
		Engines:        []string{"duckduckgo", "bing", "brave"},
		Timeout:        60 * time.Second,
	})
	
	if err != nil {
		log.Printf("Deep search error: %v\n", err)
	} else {
		fmt.Printf("Found %d results from multiple engines:\n\n", len(deepResults))
		
		// Group by engine
		byEngine := make(map[string]int)
		for _, r := range deepResults {
			byEngine[r.Engine]++
			fmt.Printf("• [%s] %s\n", r.Engine, r.Title)
		}
		
		fmt.Println("\nResults by engine:")
		for engine, count := range byEngine {
			fmt.Printf("  %s: %d results\n", engine, count)
		}
	}
	
	fmt.Println("\n✅ All tests completed!")
}