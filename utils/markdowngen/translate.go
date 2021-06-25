package markdowngen

import (
	"fmt"
	"strings"

	"github.com/bvobart/mllint/utils"
)

func List(items []interface{}) string {
	sb := strings.Builder{}
	for _, item := range items {
		sb.WriteString(fmt.Sprintln("-", item))
	}
	return sb.String()
}

func ListFiles(items utils.Filenames) string {
	sb := strings.Builder{}
	for _, item := range items {
		sb.WriteString(fmt.Sprintln("-", item))
	}
	return sb.String()
}

func CodeBlock(code string) string {
	return "```\n" + code + "\n```\n"
}
