package projectlinters_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/config"
	"gitlab.com/bvobart/mllint/projectlinters"
)

type testLinter struct {
	name string
}

func (l testLinter) Name() string {
	return l.name
}

func (l testLinter) Rules() []string {
	return nil
}

func (l testLinter) LintProject(projectdir string) ([]api.Issue, error) {
	return nil, nil
}

func TestFilterEnabled(t *testing.T) {
	linters := []api.Linter{
		testLinter{name: "1"},
		testLinter{name: "2"},
		testLinter{name: "3"},
		testLinter{name: "4"},
	}
	ruleconf := config.RuleConfig{Disabled: []string{"1", "3"}}
	actualEnabled := projectlinters.FilterEnabled(linters, ruleconf)

	enabled := []api.Linter{
		testLinter{name: "2"},
		testLinter{name: "4"},
	}
	require.Equal(t, enabled, actualEnabled)
}
