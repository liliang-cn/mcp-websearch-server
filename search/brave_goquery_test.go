package search

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestExtractBraveResult(t *testing.T) {
	html := `
	<div class="snippet" data-type="web">
		<div class="result-wrapper">
			<div class="result-content">
				<a href="https://cngoldprice.com" class="l1">
					<div class="site-name-wrapper">
						<div class="desktop-small-semibold">Cngoldprice</div>
					</div>
					<div class="title search-snippet-title">水贝金价网-今日金价_国际金价</div>
				</a>
				<div class="generic-snippet">
					<div class="content">实时国内上金所金价、国际金价和品牌金价信息。</div>
				</div>
			</div>
		</div>
	</div>`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse test html: %v", err)
	}

	title, link, snippet := extractBraveResult(doc.Find(".snippet").First())
	if title != "水贝金价网-今日金价_国际金价" {
		t.Fatalf("unexpected title: %q", title)
	}
	if link != "https://cngoldprice.com" {
		t.Fatalf("unexpected link: %q", link)
	}
	if snippet != "实时国内上金所金价、国际金价和品牌金价信息。" {
		t.Fatalf("unexpected snippet: %q", snippet)
	}
}
