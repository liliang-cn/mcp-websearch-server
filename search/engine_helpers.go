package search

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNoResults     = errors.New("no search results")
	ErrSearchBlocked = errors.New("search request blocked")
)

func normalizeSearchOutcome(engineName string, results []SearchResult, err error) ([]SearchResult, error) {
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("%w: %s returned zero results", ErrNoResults, engineName)
	}

	return results, nil
}

func validateSearchResponse(engineName string, statusCode int, body string) error {
	if statusCode != http.StatusOK {
		return fmt.Errorf("%s returned status %d", engineName, statusCode)
	}

	if isBlockedSearchPage(body) {
		return fmt.Errorf("%w: %s returned a challenge page", ErrSearchBlocked, engineName)
	}

	return nil
}

func isBlockedSearchPage(body string) bool {
	lowerBody := strings.ToLower(body)

	blockMarkers := []string{
		"please solve the challenge below to continue",
		`id="turnstile-widget"`,
		`class="captcha"`,
		"cf-turnstile-wrapper",
		"verificationcomplete",
		"verificationfailed",
	}

	for _, marker := range blockMarkers {
		if strings.Contains(lowerBody, marker) {
			return true
		}
	}

	return false
}
