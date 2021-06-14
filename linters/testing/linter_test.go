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
			Expect: func(t *stdtesting.T, report api.Report, err error) {
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
			Expect: func(t *stdtesting.T, report api.Report, err error) {
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
			Expect: func(t *stdtesting.T, report api.Report, err error) {
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
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassed",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage = "" // TODO
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllFailed",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-failed-all.xml"
					c.Testing.Coverage = "" // TODO
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "**None** of the 4 tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name: "FourTestsSixteenFiles/HalfPassed",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-half.xml"
					c.Testing.Coverage = "" // TODO
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 50, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "**2** out of **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				// require.Equal(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name: "UnfindableTestReport",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "non-existant-file.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "`non-existant-file.xml`")
				require.Contains(t, report.Details[testing.RuleTestsPass], "file could not be found")
				require.Contains(t, report.Details[testing.RuleTestsPass], "update the `testing.report` setting")
			},
		},
		{
			Name: "MalformedTestReport",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-malformed.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "was provided and found, but there was an error parsing the JUnit XML contents")
			},
		},
		{
			Name: "EmptyTestReport",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-empty.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "No tests were run")
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
