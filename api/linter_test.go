package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
)

type testLinter struct {
	name       string
	configured bool
}

func (l *testLinter) Name() string {
	return l.name
}

func (l *testLinter) Rules() []string {
	return nil
}

func (l *testLinter) Configure(_ *config.Config) error {
	l.configured = true
	return nil
}

func (l *testLinter) LintProject(projectdir string) ([]api.Issue, error) {
	return nil, nil
}

func TestFilterEnabled(t *testing.T) {
	linters := api.LinterList{
		&testLinter{name: "1"},
		&testLinter{name: "2"},
		&testLinter{name: "3"},
		&testLinter{name: "4"},
	}
	ruleconf := config.RuleConfig{Disabled: []string{"1", "3"}}
	actualEnabled := linters.FilterEnabled(ruleconf)

	enabled := api.LinterList{
		&testLinter{name: "2"},
		&testLinter{name: "4"},
	}
	require.Equal(t, enabled, actualEnabled)
}

func TestConfigure(t *testing.T) {
	linters := api.LinterList{
		&testLinter{name: "1"},
		&testLinter{name: "2"},
		&testLinter{name: "3"},
		&testLinter{name: "4"},
	}
	conf := config.Default()

	linters, err := linters.Configure(conf)
	require.NoError(t, err)

	for _, l := range linters {
		require.IsType(t, &testLinter{}, l)
		linter := l.(*testLinter)
		require.True(t, linter.configured)
	}
}
