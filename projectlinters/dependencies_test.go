package projectlinters_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/bvobart/mllint/api"
	"gitlab.com/bvobart/mllint/projectlinters"
	"gitlab.com/bvobart/mllint/utils/depsmgmt"
)

type test struct {
	Name     string
	Dir      string
	Expected []api.Issue
}

func TestUseDependencyManager(t *testing.T) {
	linter := projectlinters.UseDependencyManager{}

	tests := []test{
		{Name: "Correct/Pipenv", Dir: "test-resources/dependencies/correct-pipenv", Expected: nil},
		{Name: "Correct/Poetry", Dir: "test-resources/dependencies/correct-poetry", Expected: nil},
		{Name: "Correct/Poetry", Dir: "test-resources/dependencies/correct-poetry", Expected: nil},
		{Name: "Invalid/None", Dir: "test-resources/dependencies/none", Expected: []api.Issue{
			api.NewIssue(linter.Name(), api.SeverityError, projectlinters.MsgUseDependencyManager),
		}},
		{Name: "Invalid/RequirementsTxt", Dir: "test-resources/dependencies/requirementstxt", Expected: []api.Issue{
			api.NewIssue(linter.Name()+"/no-pip", api.SeverityWarning, projectlinters.MsgDontUsePip),
		}},
		{Name: "Invalid/Multiple", Dir: "test-resources/dependencies/multiple", Expected: []api.Issue{
			api.NewIssue(linter.Name()+"/single", api.SeverityError, fmt.Sprintf(projectlinters.MsgUseSingleDependencyManager, []depsmgmt.DependencyManagerType{depsmgmt.TypePoetry, depsmgmt.TypePipenv})),
		}},
		{Name: "Invalid/MultipleWithPip", Dir: "test-resources/dependencies/multiple-pip", Expected: []api.Issue{
			api.NewIssue(linter.Name()+"/single", api.SeverityError, fmt.Sprintf(projectlinters.MsgUseSingleDependencyManager, []depsmgmt.DependencyManagerType{depsmgmt.TypePipenv, depsmgmt.TypePip})),
			api.NewIssue(linter.Name()+"/no-pip", api.SeverityWarning, projectlinters.MsgDontUsePip),
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
