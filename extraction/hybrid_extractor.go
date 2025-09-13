package extraction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// HybridExtractor uses chromedp for intelligent content extraction
type HybridExtractor struct {
	timeout time.Duration
}

func NewHybridExtractor() *HybridExtractor {
	return &HybridExtractor{
		timeout: 30 * time.Second,
	}
}

// ExtractContent extracts the main content from a webpage
func (e *HybridExtractor) ExtractContent(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var title string
	var paragraphs []string
	var articleContent string

	err := chromedp.Run(allocCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Title(&title),
		// Try to get article content first
		chromedp.Evaluate(`
			(() => {
				// Remove script and style elements first
				document.querySelectorAll('script, style, noscript').forEach(el => el.remove());
				
				// Try to find main article content
				const articleSelectors = [
					'article', 
					'main article',
					'[role="main"]',
					'.article-content',
					'.post-content', 
					'.entry-content',
					'.content-body',
					'#article-body',
					'.story-body'
				];
				
				for (const selector of articleSelectors) {
					const elem = document.querySelector(selector);
					if (elem && elem.innerText && elem.innerText.length > 200) {
						return elem.innerText;
					}
				}
				
				// Fallback: get all paragraphs
				return null;
			})()
		`, &articleContent),
		// If no article content, get paragraphs
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('p'))
				.map(p => p.innerText.trim())
				.filter(text => text.length > 50) // Filter short paragraphs
				.slice(0, 20) // Limit to first 20 paragraphs
		`, &paragraphs),
	)

	if err != nil {
		return "", fmt.Errorf("failed to extract content from %s: %w", url, err)
	}

	// Build the final content
	var content strings.Builder
	
	if title != "" {
		content.WriteString(fmt.Sprintf("# %s\n\n", title))
	}

	// Use article content if found
	if articleContent != "" && len(articleContent) > 200 {
		content.WriteString(cleanText(articleContent))
	} else if len(paragraphs) > 0 {
		// Otherwise use paragraphs
		for _, p := range paragraphs {
			if p != "" {
				content.WriteString(p)
				content.WriteString("\n\n")
			}
		}
	}

	result := content.String()
	if result == "" || (title != "" && result == fmt.Sprintf("# %s\n\n", title)) {
		return "", fmt.Errorf("no content extracted from %s", url)
	}

	return result, nil
}

// ExtractSummary extracts a summary-friendly version of the content
func (e *HybridExtractor) ExtractSummary(ctx context.Context, url string, maxLength int) (string, error) {
	content, err := e.ExtractContent(ctx, url)
	if err != nil {
		return "", err
	}

	// Truncate if necessary
	if len(content) > maxLength {
		// Try to cut at a sentence boundary
		truncated := content[:maxLength]
		lastPeriod := strings.LastIndex(truncated, ". ")
		if lastPeriod > maxLength/2 {
			content = truncated[:lastPeriod+1]
		} else {
			content = truncated + "..."
		}
	}

	return content, nil
}

// ExtractMultiple extracts content from multiple URLs concurrently
func (e *HybridExtractor) ExtractMultiple(ctx context.Context, urls []string) map[string]string {
	results := make(map[string]string)
	resultChan := make(chan struct {
		url     string
		content string
	}, len(urls))

	// Create a shared browser context for efficiency
	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	for _, url := range urls {
		go func(u string) {
			content, err := e.extractWithContext(allocCtx, u)
			if err != nil {
				content = fmt.Sprintf("Error extracting %s: %v", u, err)
			}
			resultChan <- struct {
				url     string
				content string
			}{url: u, content: content}
		}(url)
	}

	// Collect results
	for i := 0; i < len(urls); i++ {
		result := <-resultChan
		results[result.url] = result.content
	}

	return results
}

func (e *HybridExtractor) extractWithContext(ctx context.Context, url string) (string, error) {
	var title string
	var paragraphs []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Title(&title),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('p'))
				.map(p => p.innerText.trim())
				.filter(text => text.length > 30)
				.slice(0, 10)
		`, &paragraphs),
	)

	if err != nil {
		return "", err
	}

	content := fmt.Sprintf("## %s\n\n", title)
	for _, p := range paragraphs {
		content += p + "\n\n"
	}

	return content, nil
}

// AggregateContent combines multiple contents into a single string for summarization
func AggregateContent(contents map[string]string) string {
	var aggregated strings.Builder
	
	aggregated.WriteString("# Aggregated Content from Multiple Sources\n\n")
	
	for url, content := range contents {
		aggregated.WriteString(fmt.Sprintf("## Source: %s\n\n", url))
		aggregated.WriteString(content)
		aggregated.WriteString("\n\n---\n\n")
	}
	
	return aggregated.String()
}