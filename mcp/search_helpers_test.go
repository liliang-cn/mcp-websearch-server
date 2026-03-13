package mcp

import "testing"

func TestShouldEnhanceBasicQuery(t *testing.T) {
	tests := []struct {
		query string
		want  bool
	}{
		{query: "hi", want: false},
		{query: "现在的金价是多少", want: true},
		{query: "current gold price", want: true},
		{query: "golang interfaces", want: false},
	}

	for _, tt := range tests {
		if got := shouldEnhanceBasicQuery(tt.query); got != tt.want {
			t.Fatalf("shouldEnhanceBasicQuery(%q) = %v, want %v", tt.query, got, tt.want)
		}
	}
}

func TestResolveExtractContent(t *testing.T) {
	if !resolveExtractContent(nil) {
		t.Fatal("nil extract_content should default to true")
	}

	falseValue := false
	if resolveExtractContent(&falseValue) {
		t.Fatal("false extract_content should stay false")
	}

	trueValue := true
	if !resolveExtractContent(&trueValue) {
		t.Fatal("true extract_content should stay true")
	}
}

func TestTrimPreview(t *testing.T) {
	content := "Sentence one. Sentence two should be kept. Sentence three should be removed."
	trimmed := trimPreview(content, 44)
	if trimmed != "Sentence one. Sentence two should be kept." {
		t.Fatalf("unexpected trimmed content: %q", trimmed)
	}
}
