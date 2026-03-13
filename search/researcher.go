package search

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/liliang-cn/mcp-websearch-server/extraction"
)

// ComprehensiveResearcher orchestrates a deep research process
type ComprehensiveResearcher struct {
	searcher MultiEngineSearcher
	reader   *extraction.DeepReader
}

func NewComprehensiveResearcher(searcher MultiEngineSearcher) *ComprehensiveResearcher {
	return &ComprehensiveResearcher{
		searcher: searcher,
		reader:   extraction.NewDeepReader(extraction.WithMaxLinks(5)), // Sample 5 sub-links per site
	}
}

// ResearchResult represents the final aggregated report
type ResearchResult struct {
	Topic       string
	Summary     string
	Sources     []*extraction.DeepReadResult
	ExtractedAt time.Time
}

func (cr *ComprehensiveResearcher) ConductResearch(ctx context.Context, topic string, maxSources int) (*ResearchResult, error) {
	// 1. Multi-engine search to get broad candidates
	results, err := cr.searcher.DeepSearch(ctx, topic, SearchOptions{
		MaxResults:     20,
		ExtractContent: false, // We'll do deep extraction later
		Engines:        []string{"duckduckgo", "bing", "brave"},
	})
	if err != nil {
		return nil, fmt.Errorf("initial search failed: %w", err)
	}

	// 2. Select top unique domains to ensure diversity
	selectedURLs := cr.pickDiverseURLs(results, maxSources)
	if len(selectedURLs) == 0 {
		return nil, fmt.Errorf("no suitable research sources found")
	}

	// 3. Deep read selected sources concurrently
	var wg sync.WaitGroup
	sourceResults := make([]*extraction.DeepReadResult, len(selectedURLs))
	mu := sync.Mutex{}

	for i, targetURL := range selectedURLs {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()

			// Use a shorter sub-timeout for individual site deep-reading
			subCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
			defer cancel()

			res, err := cr.reader.DeepRead(subCtx, u)
			if err == nil {
				mu.Lock()
				sourceResults[idx] = res
				mu.Unlock()
			}
		}(i, targetURL)
	}

	wg.Wait()

	// 4. Assemble final result
	validSources := []*extraction.DeepReadResult{}
	for _, s := range sourceResults {
		if s != nil {
			validSources = append(validSources, s)
		}
	}

	return &ResearchResult{
		Topic:       topic,
		Sources:     validSources,
		ExtractedAt: time.Now(),
	}, nil
}

func (cr *ComprehensiveResearcher) pickDiverseURLs(results []SearchResult, max int) []string {
	seenDomains := make(map[string]bool)
	var selected []string

	// Exclude non-research patterns
	exclude := []string{"twitter.com", "facebook.com", "youtube.com", "instagram.com", "linkedin.com", "reddit.com"}

	for _, r := range results {
		if len(selected) >= max {
			break
		}

		u, err := url.Parse(r.URL)
		if err != nil {
			continue
		}

		domain := strings.ToLower(u.Host)

		// Skip excluded domains
		skip := false
		for _, p := range exclude {
			if strings.Contains(domain, p) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Ensure diversity: only 1 URL per domain for breadth
		if !seenDomains[domain] {
			seenDomains[domain] = true
			selected = append(selected, r.URL)
		}
	}

	return selected
}

func (rr *ResearchResult) ToMarkdown() string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Comprehensive Research Report: %s\n\n", rr.Topic))

	sb.WriteString(fmt.Sprintf("*Report generated at: %s*\n\n", rr.ExtractedAt.Format(time.RFC1123)))

	sb.WriteString("## Executive Summary\n")

	sb.WriteString(fmt.Sprintf("This report synthesizes information from %d authoritative sources across multiple search engines. Each source has been deeply analyzed including its relevant sub-pages.\n\n", len(rr.Sources)))

	sb.WriteString("## Deep Dive by Source\n\n")

	for i, src := range rr.Sources {

		sb.WriteString(fmt.Sprintf("### Source %d: %s\n", i+1, src.MainTitle))

		sb.WriteString(fmt.Sprintf("**URL:** [%s](%s)\n\n", src.MainURL, src.MainURL))

		sb.WriteString("#### Core Content\n")

		sb.WriteString(src.MainContent)

		sb.WriteString("\n\n")

		if len(src.SubPages) > 0 {

			sb.WriteString("#### Supporting Knowledge (from sub-pages)\n")

			for _, sub := range src.SubPages {

				if sub.Error == "" && len(sub.Content) > 100 {

					sb.WriteString(fmt.Sprintf("- **%s**: %s\n", sub.LinkText, crTruncate(sub.Content, 300)))

				}

			}

			sb.WriteString("\n")

		}

		sb.WriteString("\n---\n\n")

	}

	sb.WriteString("## Bibliography & Source Links\n")

	for _, src := range rr.Sources {

		sb.WriteString(fmt.Sprintf("- [%s](%s)\n", src.MainTitle, src.MainURL))

	}

	return sb.String()

}

func crTruncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
