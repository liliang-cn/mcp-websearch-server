package search

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type braveGoQueryEngine struct {
	client *http.Client
}

func NewBraveGoQueryEngine() SearchEngine {
	return &braveGoQueryEngine{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (b *braveGoQueryEngine) Name() string {
	return "brave"
}

func (b *braveGoQueryEngine) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://search.brave.com/search?q=%s", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers to appear more like a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Brave results: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Brave response: %w", err)
	}

	if err := validateSearchResponse(b.Name(), resp.StatusCode, string(body)); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []SearchResult

	// Try multiple selectors for Brave results
	doc.Find(".snippet, .result-card, article[data-type='web']").Each(func(i int, s *goquery.Selection) {
		if i >= maxResults {
			return
		}

		title, link, snippet := extractBraveResult(s)

		if link != "" && title != "" {
			// Ensure link has protocol
			if strings.HasPrefix(link, "//") {
				link = "https:" + link
			} else if strings.HasPrefix(link, "/") {
				link = "https://search.brave.com" + link
			} else if !strings.HasPrefix(link, "http") {
				link = "https://" + link
			}

			results = append(results, SearchResult{
				Title:   title,
				URL:     link,
				Snippet: snippet,
				Engine:  b.Name(),
			})
		}
	})

	// If no results with primary selectors, try backup approach
	if len(results) == 0 {
		doc.Find("#results a[href]").Each(func(i int, s *goquery.Selection) {
			if i >= maxResults {
				return
			}

			title := strings.TrimSpace(s.Text())
			link, _ := s.Attr("href")

			// Skip navigation/internal links
			if link != "" && title != "" && strings.Contains(link, "http") {
				if !strings.HasPrefix(link, "http") {
					link = "https://" + link
				}

				results = append(results, SearchResult{
					Title:   title,
					URL:     link,
					Snippet: "",
					Engine:  b.Name(),
				})
			}
		})
	}

	return results, nil
}

func extractBraveResult(s *goquery.Selection) (string, string, string) {
	linkElem := s.Find("a.l1[href]").First()
	if linkElem.Length() == 0 {
		linkElem = s.Find("a[data-testid='result-title']").First()
	}
	if linkElem.Length() == 0 {
		linkElem = s.Find("h3 a").First()
	}
	if linkElem.Length() == 0 {
		linkElem = s.Find("a[href]").First()
	}

	titleElem := s.Find(".search-snippet-title, .title, .snippet-title").First()
	title := strings.TrimSpace(titleElem.Text())
	if title == "" {
		title = strings.TrimSpace(linkElem.Find(".search-snippet-title, .title, .snippet-title").First().Text())
	}
	if title == "" {
		title = strings.TrimSpace(linkElem.Text())
	}

	link, _ := linkElem.Attr("href")
	if link == "" {
		link, _ = s.Find("a[href]").First().Attr("href")
	}

	snippet := strings.TrimSpace(s.Find(".generic-snippet .content").First().Text())
	if snippet == "" {
		snippet = strings.TrimSpace(s.Find(".snippet-description").Text())
	}
	if snippet == "" {
		snippet = strings.TrimSpace(s.Find("[data-testid='result-description']").Text())
	}
	if snippet == "" {
		snippet = strings.TrimSpace(s.Find(".desc").Text())
	}
	if snippet == "" {
		snippet = strings.TrimSpace(s.Find(".content").First().Text())
	}
	if snippet == "" {
		snippet = strings.TrimSpace(s.Find("p").First().Text())
	}

	return title, link, snippet
}
