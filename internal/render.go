package internal

import (
	"bytes"
	"log"
	"math/rand"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	mdParser      goldmark.Markdown
	htmlSanitizer *bluemonday.Policy
)

func init() {
	// Goldmarkパーサーの初期化
	mdParser = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	// HTMLサニタイザーの初期化
	htmlSanitizer = bluemonday.UGCPolicy()
}

type CursorInfo struct {
	UserID   string `json:"userId"`
	Position int    `json:"position"`
	Color    string `json:"color"`
}

var seed = time.Now().UnixNano()

func generateUserID() string {
	r := rand.New(rand.NewSource(seed))
	return "user_" + string(rune(r.Intn(26)+65)) + string(rune(r.Intn(26)+65)) + string(rune(r.Intn(26)+65))
}

func generateColor() string {
	colors := []string{"#FF6B6B", "#4ECDC4", "#45B7D1", "#FFA07A", "#98D8C8", "#F7DC6F", "#BB8FCE", "#85C1E9"}
	r := rand.New(rand.NewSource(seed))
	return colors[r.Intn(len(colors))]

}

// MarkdownをHTMLにレンダリングしてサニタイズする
func renderMarkdown(source string) string {
	var buf bytes.Buffer
	if err := mdParser.Convert([]byte(source), &buf); err != nil {
		log.Printf("Markdown rendering error: %v", err)
		return ""
	}
	return htmlSanitizer.Sanitize(buf.String())
}
