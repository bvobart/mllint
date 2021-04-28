package dependencymgmt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters/dependencymgmt"
)

var linter = dependencymgmt.NewLinter()

func TestName(t *testing.T) {
	require.Equal(t, categories.DependencyMgmt.Name, linter.Name())
}

func TestRules(t *testing.T) {
	require.Equal(t, []*api.Rule{&dependencymgmt.RuleUse, &dependencymgmt.RuleSingle}, linter.Rules())
}

type linterTest struct {
	Name   string
	Dir    string
	Expect func(report api.Report, err error)
}

func TestLintProject(t *testing.T) {
	perfectScore := func(report api.Report, err error) {
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
		require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleSingle])
	}

	tests := []linterTest{
		{Name: "Correct/Pipenv", Dir: "test-resources/correct-pipenv", Expect: perfectScore},
		{Name: "Correct/Poetry", Dir: "test-resources/correct-poetry", Expect: perfectScore},
		{Name: "Invalid/None", Dir: "test-resources/none", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
		}},
		{Name: "Invalid/RequirementsTxt", Dir: "test-resources/requirementstxt", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 20, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleUse], dependencymgmt.DetailsNoRequirementsTxt)
		}},
		{Name: "Invalid/SetupPy", Dir: "test-resources/setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 30, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleUse], dependencymgmt.DetailsNoSetupPy)
		}},
		{Name: "Invalid/Multiple/Pipenv+SetupPy", Dir: "test-resources/multiple/pipenv+setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsPipenvSetupPy)
		}},
		{Name: "Invalid/Multiple/RequirementsTxt+SetupPy", Dir: "test-resources/multiple/requirementstxt+setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 20, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsRequirementsTxtSetupPy)
		}},
		{Name: "Invalid/Multiple/Pipenv+RequirementsTxt", Dir: "test-resources/multiple/pipenv+requirementstxt", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsRequirementsTxtPipenv)
		}},
		{Name: "Invalid/Multiple/Poetry+RequirementsTxt", Dir: "test-resources/multiple/poetry+requirementstxt", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsRequirementsTxtPoetry)
		}},
		{Name: "Invalid/Multiple/Poetry+SetupPy", Dir: "test-resources/multiple/poetry+setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsPoetrySetupPy)
		}},
		{Name: "Invalid/Multiple/Poetry+Pipenv", Dir: "test-resources/multiple/poetry+pipenv", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
		}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			project := api.Project{Dir: test.Dir}
			report, err := linter.LintProject(project)
			test.Expect(report, err)
		})
	}
}
