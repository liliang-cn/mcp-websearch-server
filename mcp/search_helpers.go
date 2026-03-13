package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/liliang-cn/mcp-websearch-server/extraction"
	"github.com/liliang-cn/mcp-websearch-server/search"
)

func shouldEnhanceBasicQuery(query string) bool {
	lowerQuery := strings.ToLower(query)

	keywords := []string{
		"today", "current", "now", "live", "latest", "price", "quote", "stock", "weather", "temperature", "exchange rate", "gold",
		"今天", "今日", "当前", "现在", "实时", "最新", "价格", "多少钱", "金价", "汇率", "股价", "行情", "天气", "温度",
	}

	for _, keyword := range keywords {
		if strings.Contains(lowerQuery, keyword) {
			return true
		}
	}

	return false
}

func resolveExtractContent(extractContent *bool) bool {
	if extractContent == nil {
		return true
	}

	return *extractContent
}

func extractQuickAnswerContext(ctx context.Context, results []search.SearchResult) string {
	extractor := extraction.NewHybridExtractor()
	maxSources := min(2, len(results))

	for i := 0; i < maxSources; i++ {
		extractCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
		summary, err := extractor.ExtractSummary(extractCtx, results[i].URL, 1200)
		cancel()
		if err != nil {
			continue
		}

		summary = trimPreview(summary, 800)
		if summary == "" {
			continue
		}

		return fmt.Sprintf("### Quick Answer Context\n**Source:** %s\n\n%s\n\n", results[i].URL, summary)
	}

	return ""
}

func trimPreview(content string, maxLen int) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	runes := []rune(content)
	if len(runes) <= maxLen {
		return content
	}

	truncated := string(runes[:maxLen])
	if idx := strings.LastIndexAny(truncated, "。\n.!?"); idx > len(truncated)/2 {
		return strings.TrimSpace(truncated[:idx+1])
	}

	return strings.TrimSpace(truncated) + "..."
}
