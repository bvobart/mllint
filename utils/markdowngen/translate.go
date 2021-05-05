package markdowngen

import (
	"fmt"
	"strings"
)

func List(items []interface{}) string {
	sb := strings.Builder{}
	for _, item := range items {
		sb.WriteString(fmt.Sprintln("-", item))
	}
	return sb.String()
}
