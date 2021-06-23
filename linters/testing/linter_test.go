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
	require.Equal(t, []*api.Rule{&testing.RuleHasTests, &testing.RuleTestsPass, &testing.RuleTestCoverage, &testing.RuleTestsFolder}, linter.Rules())

	suite := testutils.NewLinterTestSuite(linter, []testutils.LinterTest{
		{
			Name: "NoTestsNoFiles",
			Dir:  ".",
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 0, report.Scores[testing.RuleHasTests])

				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "No test report was provided")
				require.Contains(t, report.Details[testing.RuleTestsPass], "update the `testing.report` setting")

				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "No test coverage report was provided")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "update the `testing.coverage` setting")

				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
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
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
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
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
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
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassed/ZeroCoverage",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage = "coverage-0.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "achieved 0.0% line test coverage")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "80.0% is the target amount of test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassedHalfCoverage",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage = "coverage-50.xml"
					c.Testing.CoverageTarget = 100
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 50, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "achieved 50.0% line test coverage")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "100.0% is the target amount of test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassedHalfCoverageButTargetAchieved",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage = "coverage-50.xml"
					c.Testing.CoverageTarget = 50
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "achieved 50.0% line test coverage")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "which meets the target of 50.0% test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassedFullCoverage",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage = "coverage-100.xml"
					c.Testing.CoverageTarget = 100
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "Wow! Congratulations!")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "100% line test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllFailed",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-failed-all.xml"
					c.Testing.Coverage = ""
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "**None** of the 4 tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
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
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
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
			Name: "UnfindableCoverageReport",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Coverage = "non-existant-file.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "`non-existant-file.xml`")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "file could not be found")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "update the `testing.coverage` setting")
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
			Name: "MalformedCoverageReport",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Coverage = "coverage-malformed.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "was provided and found, but there was an error parsing the Cobertura XML contents")
			},
		},
		{
			Name: "EmptyTestReports",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-empty.xml"
					c.Testing.Coverage = "coverage-empty.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])

				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "No tests were run")

				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "no lines were covered")
			},
		},
		{
			Name: "ZeroCoverageTarget",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Coverage = "coverage-0.xml"
					c.Testing.CoverageTarget = 0
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestCoverage])
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
