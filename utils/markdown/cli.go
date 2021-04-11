package markdown

import markdown "github.com/MichaelMure/go-term-markdown"

// Renders Markdown formatted text using go-term-markdown.
func Render(text string) string {
	return string(markdown.Render(text, 120, 0))
}
