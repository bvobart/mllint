package markdown

import (
	"regexp"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
)

// Renders Markdown formatted text using go-term-markdown.
func Render(text string) string {
	text = strings.ReplaceAll(text, "❌", color.RedString("❌"))
	text = prettifyInlineCode(text)
	return string(markdown.Render(text, 120, 0))
}

func prettifyInlineCode(text string) string {
	regex := regexp.MustCompile("`([^`\n\r]+)`")
	repl := color.Set(color.Italic, color.FgYellow).Sprint("$1")
	color.Unset()
	return regex.ReplaceAllString(text, repl)
}
