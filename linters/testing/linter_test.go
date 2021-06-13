package testing_test

import (
	"fmt"
	stdtesting "testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters/testing"
	"github.com/bvobart/mllint/linters/testutils"
	"github.com/bvobart/mllint/utils"
	"github.com/stretchr/testify/require"
)

func TestTestingLinter(t *stdtesting.T) {
	linter := testing.NewLinter()
	require.Equal(t, "Testing", linter.Name())

	suite := testutils.NewLinterTestSuite(linter, []testutils.LinterTest{
		{
			Name: "NoTestsNoFiles",
			Dir:  ".",
			Expect: func(report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 0, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name:    "NoTestsSixteenFiles",
			Dir:     ".",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16)),
			Expect: func(report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 0, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name:    "OneTestSixteenFiles",
			Dir:     ".",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(1))),
			Expect: func(report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 25, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name:    "FourTestsSixteenFiles",
			Dir:     ".",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))),
			Expect: func(report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
	})

	suite.DefaultOptions().WithConfig(config.Default())
	suite.RunAll(t)
}

func createPythonFilenames(n int) utils.Filenames {
	files := make(utils.Filenames, n)
	for i := 0; i < n; i++ {
		files[i] = fmt.Sprint("file", i, ".py")
	}
	return files
}

func createPythonTestFilenames(n int) utils.Filenames {
	files := make(utils.Filenames, n)
	for i := 0; i < n; i++ {
		files[i] = fmt.Sprint("file", i, "_test.py")
	}
	return files
}
