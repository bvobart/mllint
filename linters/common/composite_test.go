package common_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters/common"
	"github.com/stretchr/testify/require"
)

var testRule1 = api.Rule{
	Slug:    "test-rule-1",
	Name:    "Test should have first rule",
	Details: "",
}

var testRule2 = api.Rule{
	Slug:    "test-rule-2",
	Name:    "Test should have second rule",
	Details: "",
}

var testRule3 = api.Rule{
	Slug:    "test-rule-3",
	Name:    "Test should have third rule",
	Details: "",
}

var testRule4 = api.Rule{
	Slug:    "test-rule-4",
	Name:    "Test should have fourth rule",
	Details: "",
}

type testLinter struct {
	name         string
	rules        []api.Rule
	configured   bool
	configureErr error
	report       api.Report
	lintErr      error
}

func (l *testLinter) Name() string {
	return l.name
}

func (l *testLinter) Rules() []api.Rule {
	return l.rules
}

func (l *testLinter) Configure(_ *config.Config) error {
	if l.configureErr != nil {
		return l.configureErr
	}
	l.configured = true
	return nil
}

func (l *testLinter) LintProject(projectdir string) (api.Report, error) {
	return l.report, l.lintErr
}

func TestCompositeLinterNameRules(t *testing.T) {
	name := "TestCategory"
	rulePrefix := name + ": "
	linter := common.NewCompositeLinter(name,
		&testLinter{name: "Linter1", rules: []api.Rule{testRule1, testRule2}},
		&testLinter{name: "Linter2", rules: []api.Rule{testRule3, testRule4}},
	)

	require.Equal(t, name, linter.Name())

	rules := linter.Rules()
	for _, rule := range rules {
		require.True(t, strings.HasPrefix(rule.Name, rulePrefix))
	}
	require.Equal(t, testRule1.Name, strings.TrimPrefix(rules[0].Name, rulePrefix))
	require.Equal(t, testRule2.Name, strings.TrimPrefix(rules[1].Name, rulePrefix))
	require.Equal(t, testRule3.Name, strings.TrimPrefix(rules[2].Name, rulePrefix))
	require.Equal(t, testRule4.Name, strings.TrimPrefix(rules[3].Name, rulePrefix))
}

func TestCompositeLinterConfigure(t *testing.T) {
	name := "TestCategory"
	linter1 := &testLinter{name: "Linter1", rules: []api.Rule{testRule1, testRule2}}
	linter2 := &testLinter{name: "Linter2", rules: []api.Rule{testRule3, testRule4}}
	compLinter := common.NewCompositeLinter(name, linter1, linter2)

	conf := config.Default()
	require.NoError(t, compLinter.Configure(conf))
	require.True(t, linter1.configured)
	require.True(t, linter2.configured)
}

func TestCompositeLinterConfigureErr(t *testing.T) {
	name := "TestCategory"
	configureErr := errors.New("test error")
	linter1 := &testLinter{name: "Linter1", rules: []api.Rule{testRule1, testRule2}}
	linter2 := &testLinter{name: "Linter2", rules: []api.Rule{testRule3, testRule4}, configureErr: configureErr}
	compLinter := common.NewCompositeLinter(name, linter1, linter2)

	conf := config.Default()
	err := compLinter.Configure(conf)

	require.Error(t, err)
	require.ErrorIs(t, err, configureErr)
	require.True(t, strings.Contains(err.Error(), linter2.Name()))

	require.True(t, linter1.configured)
	require.False(t, linter2.configured)
}

func TestCompositeLinterLintProject(t *testing.T) {
	// Given: 2 linters and their expected reports, combined into a CompositeLinter
	report1 := api.Report{
		Scores: map[api.Rule]float64{
			testRule1: 100,
			testRule2: 65,
		},
		Details: map[api.Rule]string{
			testRule2: "rule2 details",
		},
	}
	report2 := api.Report{
		Scores: map[api.Rule]float64{
			testRule4: 42,
		},
		Details: map[api.Rule]string{
			testRule4: "rule4 details",
		},
	}
	linter1 := &testLinter{name: "linter1", rules: []api.Rule{testRule1, testRule2}, report: report1}
	linter2 := &testLinter{name: "linter2", rules: []api.Rule{testRule3, testRule4}, report: report2}

	name := "TestCategory"
	compLinter := common.NewCompositeLinter(name, linter1, linter2)

	// When: compLinter.LintProject is called
	projectdir := "test"
	report, err := compLinter.LintProject(projectdir)
	require.NoError(t, err)

	// Then: expect that the report contains the scores and details from the expected reports above,
	// for a subset of the composite linter's Rules()
	scores := []float64{}
	details := []string{}
	for _, rule := range compLinter.Rules() {
		score, found := report.Scores[rule]
		if found {
			scores = append(scores, score)
		}

		detail, found := report.Details[rule]
		if found {
			details = append(details, detail)
		}
	}

	require.Equal(t, []float64{100, 65, 42}, scores)
	require.Equal(t, []string{"rule2 details", "rule4 details"}, details)
}

func TestCompositeLinterLintProjectErr(t *testing.T) {
	name := "TestCategory"
	lintErr := errors.New("test error")
	linter1 := &testLinter{name: "linter1", rules: []api.Rule{testRule1, testRule2}}
	linter2 := &testLinter{name: "linter2", rules: []api.Rule{testRule3, testRule4}, lintErr: lintErr}
	compLinter := common.NewCompositeLinter(name, linter1, linter2)

	projectdir := "test"
	_, err := compLinter.LintProject(projectdir)
	require.Error(t, err)
	require.ErrorIs(t, err, lintErr)
	require.True(t, strings.Contains(err.Error(), linter2.Name()))
}
