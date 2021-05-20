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

const name = "TestCategory"

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
	rules        []*api.Rule
	configured   bool
	configureErr error
	report       api.Report
	lintErr      error
}

func (l *testLinter) Name() string {
	return l.name
}

func (l *testLinter) Rules() []*api.Rule {
	return l.rules
}

func (l *testLinter) Configure(_ *config.Config) error {
	if l.configureErr != nil {
		return l.configureErr
	}
	l.configured = true
	return nil
}

func (l *testLinter) LintProject(project api.Project) (api.Report, error) {
	return l.report, l.lintErr
}

func TestCompositeLinterName(t *testing.T) {
	linter := common.NewCompositeLinter(name,
		&testLinter{name: "Linter1", rules: []*api.Rule{&testRule1, &testRule2}},
		&testLinter{name: "Linter2", rules: []*api.Rule{&testRule3, &testRule4}},
	)
	require.Equal(t, name, linter.Name())
}

func TestCompositeLinterRules(t *testing.T) {
	linter := common.NewCompositeLinter(name,
		&testLinter{name: "Linter 1", rules: []*api.Rule{&testRule1, &testRule2}},
		&testLinter{name: "Linter 2", rules: []*api.Rule{&testRule3, &testRule4}},
	)

	rules := linter.Rules()
	require.Equal(t, []*api.Rule{&testRule1, &testRule2, &testRule3, &testRule4}, rules)
}

func TestCompositeLinterSameRulePointers(t *testing.T) {
	linter := common.NewCompositeLinter(name,
		&testLinter{name: "Linter 1", rules: []*api.Rule{&testRule1, &testRule2}},
		&testLinter{name: "Linter 2", rules: []*api.Rule{&testRule3, &testRule4}},
	)

	rules1 := linter.Rules()
	rules2 := linter.Rules()
	require.Equal(t, len(rules1), len(rules2))

	for i, rule1 := range rules1 {
		rule2 := rules2[i]
		require.True(t, rule1 == rule2) // checks pointer equality
	}
}

func TestCompositeLinterConfigure(t *testing.T) {
	linter1 := &testLinter{name: "Linter1", rules: []*api.Rule{&testRule1, &testRule2}}
	linter2 := &testLinter{name: "Linter2", rules: []*api.Rule{&testRule3, &testRule4}}
	compLinter := common.NewCompositeLinter(name, linter1, linter2)

	conf := config.Default()
	require.NoError(t, compLinter.Configure(conf))
	require.True(t, linter1.configured)
	require.True(t, linter2.configured)
}

func TestCompositeLinterConfigureErr(t *testing.T) {
	configureErr := errors.New("test error")
	linter1 := &testLinter{name: "Linter1", rules: []*api.Rule{&testRule1, &testRule2}}
	linter2 := &testLinter{name: "Linter2", rules: []*api.Rule{&testRule3, &testRule4}, configureErr: configureErr}
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
	linter1 := &testLinter{name: "linter1", rules: []*api.Rule{&testRule1, &testRule2}, report: report1}
	linter2 := &testLinter{name: "linter2", rules: []*api.Rule{&testRule3, &testRule4}, report: report2}

	compLinter := common.NewCompositeLinter(name, linter1, linter2)
	compLinter.SetRunner(nil)

	// When: compLinter.LintProject is called
	project := api.Project{Dir: "test"}
	report, err := compLinter.LintProject(project)
	require.NoError(t, err)

	// Then: expect that the report contains the scores and details from the expected reports above,
	// for the subset of the composite linter's Rules() that the linter reported on.
	scores := []float64{}
	details := []string{}
	for _, rule := range compLinter.Rules() {
		score, found := report.Scores[*rule]
		if found {
			scores = append(scores, score)
		}

		detail, found := report.Details[*rule]
		if found {
			details = append(details, detail)
		}
	}

	require.Equal(t, []float64{100, 65, 42}, scores)
	require.Equal(t, []string{"rule2 details", "rule4 details"}, details)
}

func TestCompositeLinterLintProjectErr(t *testing.T) {
	lintErr := errors.New("test error")
	linter1 := &testLinter{name: "linter1", rules: []*api.Rule{&testRule1, &testRule2}}
	linter2 := &testLinter{name: "linter2", rules: []*api.Rule{&testRule3, &testRule4}, lintErr: lintErr}
	compLinter := common.NewCompositeLinter(name, linter1, linter2)

	project := api.Project{Dir: "test"}
	_, err := compLinter.LintProject(project)
	require.Error(t, err)
	require.ErrorIs(t, err, lintErr)
	require.True(t, strings.Contains(err.Error(), linter2.Name()))
}
