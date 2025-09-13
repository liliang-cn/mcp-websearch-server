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

type duckDuckGoSearchEngine struct {
	client *http.Client
}

func NewDuckDuckGoSearchEngine() SearchEngine {
	return &duckDuckGoSearchEngine{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (d *duckDuckGoSearchEngine) Name() string {
	return "duckduckgo"
}

func (d *duckDuckGoSearchEngine) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://duckduckgo.com/?q=%s", url.QueryEscape(query))

	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var results []SearchResult
	var nodes []*cdp.Node

	err := chromedp.Run(allocCtx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible(`#links`, chromedp.ByID),
		chromedp.Nodes(`#links .result`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search DuckDuckGo: %w", err)
	}

	for i, node := range nodes {
		if i >= maxResults {
			break
		}

		var title, link, snippet string

		err := chromedp.Run(allocCtx,
			chromedp.Text(`.result__title`, &title, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue(`.result__title a`, "href", &link, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`.result__snippet`, &snippet, chromedp.ByQuery, chromedp.FromNode(node)),
		)

		if err == nil && link != "" {
			if strings.HasPrefix(link, "//") {
				link = "https:" + link
			} else if !strings.HasPrefix(link, "http") {
				link = "https://duckduckgo.com" + link
			}

			results = append(results, SearchResult{
				Title:   strings.TrimSpace(title),
				URL:     link,
				Snippet: strings.TrimSpace(snippet),
				Engine:  d.Name(),
			})
		}
	}

	return results, nil
}
