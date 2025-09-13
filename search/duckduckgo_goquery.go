package search

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type duckDuckGoGoQueryEngine struct {
	client *http.Client
}

func NewDuckDuckGoGoQueryEngine() SearchEngine {
	return &duckDuckGoGoQueryEngine{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (d *duckDuckGoGoQueryEngine) Name() string {
	return "duckduckgo"
}

func (d *duckDuckGoGoQueryEngine) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	// DuckDuckGo HTML version
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))
	
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	
	// Set headers to appear more like a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	
	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DuckDuckGo results: %w", err)
	}
	defer resp.Body.Close()
	
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}
	
	var results []SearchResult
	
	// For DuckDuckGo HTML version, results are in a simpler format
	doc.Find(".result, .web-result").Each(func(i int, s *goquery.Selection) {
		if i >= maxResults {
			return
		}
		
		// Extract title and link
		var title, link string
		
		// For HTML version
		titleElem := s.Find(".result__title a, h2 a").First()
		if titleElem.Length() == 0 {
			titleElem = s.Find("a.result__a").First()
		}
		
		title = strings.TrimSpace(titleElem.Text())
		link, _ = titleElem.Attr("href")
		
		// Extract snippet
		snippet := strings.TrimSpace(s.Find(".result__snippet").Text())
		if snippet == "" {
			snippet = strings.TrimSpace(s.Find(".snippet").Text())
		}
		if snippet == "" {
			snippet = strings.TrimSpace(s.Find("a.result__snippet").Text())
		}
		
		if link != "" && title != "" {
			// Clean up DuckDuckGo redirect URLs
			if strings.Contains(link, "duckduckgo.com/l/") {
				// Extract actual URL from redirect if possible
				if u, err := url.Parse(link); err == nil {
					if actualURL := u.Query().Get("uddg"); actualURL != "" {
						if decoded, err := url.QueryUnescape(actualURL); err == nil {
							link = decoded
						}
					}
				}
			}
			
			// Ensure proper URL format
			if strings.HasPrefix(link, "//") {
				link = "https:" + link
			} else if !strings.HasPrefix(link, "http") {
				if !strings.Contains(link, "duckduckgo.com") {
					link = "https://" + link
				}
			}
			
			results = append(results, SearchResult{
				Title:   title,
				URL:     link,
				Snippet: snippet,
				Engine:  d.Name(),
			})
		}
	})
	
	// Try alternative selectors for the no-JS version
	if len(results) == 0 {
		doc.Find(".links_main a.result__a").Each(func(i int, s *goquery.Selection) {
			if i >= maxResults {
				return
			}
			
			title := strings.TrimSpace(s.Text())
			link, _ := s.Attr("href")
			
			// Get snippet from next sibling
			snippet := strings.TrimSpace(s.Parent().Find(".result__snippet").Text())
			
			if link != "" && title != "" {
				results = append(results, SearchResult{
					Title:   title,
					URL:     link,
					Snippet: snippet,
					Engine:  d.Name(),
				})
			}
		})
	}
	
	return results, nil
}