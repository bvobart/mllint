package projectlinters_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/projectlinters"
	"gitlab.com/bvobart/mllint/utils/depsmgmt"
)

type linterTest struct {
	Name     string
	Dir      string
	Expected []api.Issue
}

func TestUseDependencyManager(t *testing.T) {
	linter := projectlinters.UseDependencyManager{}

	tests := []linterTest{
		{Name: "Correct/Pipenv", Dir: "test-resources/dependencies/correct-pipenv", Expected: nil},
		{Name: "Correct/Poetry", Dir: "test-resources/dependencies/correct-poetry", Expected: nil},
		{Name: "Invalid/None", Dir: "test-resources/dependencies/none", Expected: []api.Issue{
			api.NewIssue(linter.Name(), "", api.SeverityError, projectlinters.MsgUseDependencyManager),
		}},
		{Name: "Invalid/RequirementsTxt", Dir: "test-resources/dependencies/requirementstxt", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleNoRequirementsTxt, api.SeverityWarning, projectlinters.MsgNoRequirementsTxt),
		}},
		{Name: "Invalid/SetupPy", Dir: "test-resources/dependencies/setuppy", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleNoSetupPy, api.SeverityWarning, projectlinters.MsgNoSetupPy),
		}},

		{Name: "Invalid/Multiple/Pipenv+SetupPy", Dir: "test-resources/dependencies/multiple/pipenv+setuppy", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleDontCombinePipenvSetupPy, api.SeverityInfo, projectlinters.MsgDontCombinePipenvSetupPy),
		}},
		{Name: "Invalid/Multiple/RequirementsTxt+SetupPy", Dir: "test-resources/dependencies/multiple/requirementstxt+setuppy", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleDontCombineRequirementsTxtSetupPy, api.SeverityWarning, projectlinters.MsgDontCombineRequirementsTxtSetupPy),
		}},
		{Name: "Invalid/Multiple/Pipenv+RequirementsTxt", Dir: "test-resources/dependencies/multiple/poetry+requirementstxt", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleDontCombineRequirementsTxtPoetryPipenv, api.SeverityWarning, projectlinters.MsgDontCombineRequirementsTxtPoetryPipenv),
		}},
		{Name: "Invalid/Multiple/Poetry+RequirementsTxt", Dir: "test-resources/dependencies/multiple/poetry+requirementstxt", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleDontCombineRequirementsTxtPoetryPipenv, api.SeverityWarning, projectlinters.MsgDontCombineRequirementsTxtPoetryPipenv),
		}},
		{Name: "Invalid/Multiple/Poetry+SetupPy", Dir: "test-resources/dependencies/multiple/poetry+setuppy", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleDontCombinePoetrySetupPy, api.SeverityInfo, projectlinters.MsgDontCombinePoetrySetupPy),
		}},
		{Name: "Invalid/Multiple/Poetry+Pipenv", Dir: "test-resources/dependencies/multiple/poetry+pipenv", Expected: []api.Issue{
			api.NewIssue(linter.Name(), projectlinters.RuleSingle, api.SeverityWarning, fmt.Sprintf(projectlinters.MsgUseSingleDependencyManager, []depsmgmt.DependencyManagerType{depsmgmt.TypePoetry, depsmgmt.TypePipenv})),
		}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			issues, err := linter.LintProject(test.Dir)
			require.NoError(t, err)
			require.Equal(t, test.Expected, issues)
		})
	}
}
