package commands

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/markdown"
)

// returns the current dir if args is empty, or the absolute path to the folder pointed to by args[0]
func parseProjectDir(args []string) (string, error) {
	currentdir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) == 0 {
		return currentdir, nil
	}

	projectdir := path.Join(currentdir, args[0])
	if !utils.FolderExists(projectdir) {
		return "", fmt.Errorf("%w: %s", ErrNotAFolder, projectdir)
	}

	return projectdir, nil
}

func formatInlineCode(text string) string {
	return color.New(color.Reset, color.Italic, color.FgYellow).Sprint(text)
}

func prettyPrintLinters(linters map[api.Category]api.Linter) {
	if len(linters) == 0 {
		color.Red("Oh no! Your mllint configuration has disabled ALL rules!")
		fmt.Println()
	}

	for cat, linter := range linters {
		prettyPrintCategory(cat)
		prettyPrintLinter(linter)
		fmt.Println()
	}
}

func prettyPrintCategory(cat api.Category) {
	bold := color.New(color.Bold)
	code := color.New(color.Bold, color.Italic, color.FgYellow)
	bold.Print(cat.Name, " (")
	code.Print(cat.Slug)
	bold.Println(")")
}

func prettyPrintLinter(linter api.Linter) {
	if linter == nil {
		return
	}

	rules := linter.Rules()
	if len(rules) == 0 {
		fmt.Println("None")
		return
	}

	builder := strings.Builder{}
	for _, rule := range rules {
		if !rule.Disabled {
			builder.WriteString(fmt.Sprintf("- `%s` — %s\n", rule.Slug, rule.Name))
		}
	}

	fmt.Print(markdown.Render(builder.String()))
}

func prettyPrintRule(rule api.Rule) {
	bold := color.New(color.Bold)
	code := color.New(color.Bold, color.Italic, color.FgYellow)
	bold.Print(rule.Name, " (")
	code.Print(rule.Slug)
	bold.Println(")")

	faint := color.New(color.Faint)
	faint.Println("Rule")
	faint.Printf("Weight: %.0f\n", rule.Weight)

	fmt.Println()
	fmt.Println(markdown.Render(rule.Details))
}
