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

type bingSearchEngine struct {
	client *http.Client
}

func NewBingSearchEngine() SearchEngine {
	return &bingSearchEngine{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (b *bingSearchEngine) Name() string {
	return "bing"
}

func (b *bingSearchEngine) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://www.bing.com/search?q=%s", url.QueryEscape(query))

	allocCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var results []SearchResult
	var nodes []*cdp.Node

	err := chromedp.Run(allocCtx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible(`#b_results`, chromedp.ByID),
		chromedp.Nodes(`#b_results .b_algo`, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search Bing: %w", err)
	}

	for i, node := range nodes {
		if i >= maxResults {
			break
		}

		var title, link, snippet string

		err := chromedp.Run(allocCtx,
			chromedp.Text(`h2`, &title, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue(`h2 a`, "href", &link, nil, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.Text(`.b_caption p`, &snippet, chromedp.ByQuery, chromedp.FromNode(node)),
		)

		if err == nil && link != "" {
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
