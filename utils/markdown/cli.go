package markdown

import (
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
)

// Renders Markdown formatted text using go-term-markdown.
func Render(text string) string {
	text = strings.ReplaceAll(text, "❌", color.RedString("❌"))
	return string(markdown.Render(text, 120, 0))
}
