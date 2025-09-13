package search

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type braveSearchEngine struct {
	client *http.Client
}

func NewBraveSearchEngine() SearchEngine {
	return &braveSearchEngine{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (b *braveSearchEngine) Name() string {
	return "brave"
}

func (b *braveSearchEngine) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://search.brave.com/search?q=%s", url.QueryEscape(query))

	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var results []SearchResult
	var nodes []*cdp.Node

	err := chromedp.Run(allocCtx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible(`#results`, chromedp.ByID),
		chromedp.Nodes(`#results .snippet`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search Brave: %w", err)
	}

	for i, node := range nodes {
		if i >= maxResults {
			break
		}

		var title, link, snippet string

		err := chromedp.Run(allocCtx,
			chromedp.Text(`.snippet-title`, &title, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue(`.result-header a`, "href", &link, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`.snippet-description`, &snippet, chromedp.ByQuery, chromedp.FromNode(node)),
		)

		if err == nil && link != "" {
			if !strings.HasPrefix(link, "http") {
				link = "https://" + link
			}

			results = append(results, SearchResult{
				Title:   strings.TrimSpace(title),
				URL:     link,
				Snippet: strings.TrimSpace(snippet),
				Engine:  b.Name(),
			})
		}
	}

	return results, nil
}
