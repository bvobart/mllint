package template_test

import (
	"errors"
	stdtesting "testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters/template"
	"github.com/bvobart/mllint/linters/testutils"
	"github.com/bvobart/mllint/utils"
	"github.com/stretchr/testify/require"
)

// This test serves as an example test for a linter.
// When implementing a new linter, copy and edit this test to use the linter you're implementing,
// and create a suite of integration tests that cover the functionality of the linter.
//
// You can use options on the suite or on the test in order to e.g. specify which Python files to set on the project,
// supply a specific Config that will be passed to the linter if it implements api.Configurable,
// or provide / auto-detect the dependency managers in the project.
func TestTemplateLinter(t *stdtesting.T) {
	linter := template.NewLinter()
	require.Equal(t, "Linter Template", linter.Name())
	require.Equal(t, []*api.Rule{&template.RuleSomething}, linter.Rules())

	suite := testutils.NewLinterTestSuite(linter, []testutils.LinterTest{
		{
			Name:    "ExampleTest",
			Dir:     ".",
			Options: testutils.NewOptions().WithConfig(config.Default()),
			Expect: func(report api.Report, err error) {
				require.Error(t, err, errors.New("not implemented"))
				require.EqualValues(t, 80, report.Scores[template.RuleSomething])
			},
		},
	})
	// use DefaultOptions to edit the options that will be applied to every test (unless overridden by test options)
	suite.DefaultOptions().UsePythonFiles(utils.Filenames{"example.py"})
	suite.RunAll(t)
}
