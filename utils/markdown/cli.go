package markdown

import (
	"regexp"
	"strings"

	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/fatih/color"
	"github.com/nathan-fiscaletti/consolesize-go"
)

// Renders Markdown formatted text using go-term-markdown.
func Render(text string) string {
	text = strings.ReplaceAll(text, "❌", color.RedString("❌"))
	text = prettifyInlineCode(text)
	return string(markdown.Render(text, getRenderWidth(), 0))
}

func prettifyInlineCode(text string) string {
	regex := regexp.MustCompile("`([^`\n\r]+)`")
	repl := color.Set(color.Italic, color.FgYellow).Sprint("$1")
	color.Unset()
	return regex.ReplaceAllString(text, repl)
}

func getRenderWidth() int {
	cols, _ := consolesize.GetConsoleSize()
	if cols < 40 { // anything under 40 colums makes the output pretty much unreadable anyway
		return 120
	}
	return cols
}
