package extraction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type ChromedpExtractor struct {
	timeout time.Duration
}

func NewChromedpExtractor() *ChromedpExtractor {
	return &ChromedpExtractor{
		timeout: 30 * time.Second,
	}
}

func (e *ChromedpExtractor) ExtractContent(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var content string
	var title string
	var bodyText string

	err := chromedp.Run(allocCtx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Title(&title),
		chromedp.Evaluate(`
			(function() {
				// Remove script and style elements
				var scripts = document.querySelectorAll('script, style, noscript');
				scripts.forEach(function(el) { el.remove(); });
				
				// Try to find main content areas
				var mainContent = document.querySelector('main, article, .content, #content, .post, .entry-content');
				if (mainContent) {
					return mainContent.innerText;
				}
				
				// Fallback to body text
				return document.body.innerText;
			})()
		`, &bodyText),
	)

	if err != nil {
		return "", fmt.Errorf("failed to extract content from %s: %w", url, err)
	}

	bodyText = cleanText(bodyText)

	if title != "" {
		content = fmt.Sprintf("# %s\n\n%s", title, bodyText)
	} else {
		content = bodyText
	}

	return content, nil
}

func cleanText(text string) string {
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	lastWasEmpty := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanedLines = append(cleanedLines, line)
			lastWasEmpty = false
		} else if !lastWasEmpty && len(cleanedLines) > 0 {
			cleanedLines = append(cleanedLines, "")
			lastWasEmpty = true
		}
	}

	result := strings.Join(cleanedLines, "\n")
	result = strings.TrimSpace(result)

	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}

	return result
}
