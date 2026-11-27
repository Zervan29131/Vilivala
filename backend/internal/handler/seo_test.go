package handler

import (
	"strings"
	"testing"
)

func TestInjectMetaHTMLArticle(t *testing.T) {
	tpl := `<!doctype html><html><head><title>Old</title><meta name="description" content="old" /></head><body></body></html>`
	out := InjectMetaHTML(tpl, PageMeta{
		Title:       `Go & Vue <Test>`,
		Description: `Desc "quoted"`,
		URL:         "https://example.com/posts/go-vue",
		Type:        "article",
		PublishedAt: "2026-11-27T12:00:00Z",
		UpdatedAt:   "2026-11-27T13:00:00Z",
		Author:      "Zervan",
	})

	checks := []string{
		`<title>Go &amp; Vue &lt;Test&gt;</title>`,
		`property="og:type" content="article"`,
		`rel="canonical" href="https://example.com/posts/go-vue"`,
		`application/ld+json`,
		`"@type":"BlogPosting"`,
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestExtractSnippet(t *testing.T) {
	md := "# Title\n\n```go\nfmt.Println(\"hidden\")\n```\n\n这是一段介绍个人博客架构的文字，包含关键的 MySQL 全文检索说明，方便后续查找。"
	s := ExtractSnippet(md, "", "MySQL", 30)
	if !strings.Contains(s, "MySQL") {
		t.Fatalf("expected snippet to contain query term, got %q", s)
	}
	if strings.Contains(s, "fmt.Println") {
		t.Fatalf("expected code block to be stripped from snippet, got %q", s)
	}
}
