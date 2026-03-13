package search

import (
	"context"
	"errors"
	"testing"
)

func TestNormalizeSearchOutcome_EmptyResults(t *testing.T) {
	_, err := normalizeSearchOutcome("bing", nil, nil)
	if !errors.Is(err, ErrNoResults) {
		t.Fatalf("expected ErrNoResults, got %v", err)
	}
}

func TestValidateSearchResponse_BlockedChallenge(t *testing.T) {
	err := validateSearchResponse("bing", 200, `<div class="captcha"><div id="turnstile-widget"></div><div>Please solve the challenge below to continue</div></div>`)
	if !errors.Is(err, ErrSearchBlocked) {
		t.Fatalf("expected ErrSearchBlocked, got %v", err)
	}
}

func TestMultiEngineSearcher_SearchFallsBackOnEmptyResults(t *testing.T) {
	searcher := &multiEngineSearcher{
		engines: map[string]SearchEngine{
			"bing": &mockSearchEngine{name: "bing"},
			"brave": &mockSearchEngine{
				name: "brave",
				results: []SearchResult{
					{Title: "Fallback Result", URL: "http://fallback.example", Engine: "brave"},
				},
			},
		},
	}

	results, err := searcher.Search(context.Background(), "gold price", SearchOptions{MaxResults: 1})
	if err != nil {
		t.Fatalf("expected fallback search to succeed, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Engine != "brave" {
		t.Fatalf("expected brave fallback result, got %s", results[0].Engine)
	}
}

func TestHybridSearcher_SearchFallsBackOnEmptyResults(t *testing.T) {
	searcher := &HybridMultiEngineSearcher{
		engines: map[string]SearchEngine{
			"duckduckgo": &mockSearchEngine{name: "duckduckgo"},
			"bing": &mockSearchEngine{
				name: "bing",
				results: []SearchResult{
					{Title: "Fallback Result", URL: "http://fallback.example", Engine: "bing"},
				},
			},
		},
	}

	results, err := searcher.Search(context.Background(), "gold price", SearchOptions{MaxResults: 1})
	if err != nil {
		t.Fatalf("expected fallback search to succeed, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if results[0].Engine != "bing" {
		t.Fatalf("expected bing fallback result, got %s", results[0].Engine)
	}
}
