package testing_test

import (
	"fmt"
	"path"
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

				require.EqualValues(t, 0, report.Scores[testing.RuleTestsFolder])
			},
		},
		{
			Name:    "NoTestsSixteenFiles",
			Dir:     ".",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16)),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 0, report.Scores[testing.RuleHasTests])
				require.Equal(t, "There are **0** test files in your project, but `mllint` was expecting at least **1**.", report.Details[testing.RuleHasTests])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsFolder])
				require.Equal(t, "Tip for when you start implementing tests: create a folder called `tests` at the root of your project and place all your Python test files in there, as per common convention.", report.Details[testing.RuleTestsFolder])
			},
		},
		{
			Name:    "OneTestFifteenFiles", // 1 + 15 = 16, makes it easier to calculate percentages with.
			Dir:     ".",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(15).Concat(createPythonTestFilenames(1))),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 31.25, report.Scores[testing.RuleHasTests])
				require.Contains(t, report.Details[testing.RuleHasTests], "There is **1** test file in your project")
				require.Contains(t, report.Details[testing.RuleHasTests], "minimum of **1** test file required")
				require.Contains(t, report.Details[testing.RuleHasTests], "equates to **6.25%** of Python files in your project being tests")
				require.Contains(t, report.Details[testing.RuleHasTests], "`mllint` expects that **20%** of your project's Python files are tests")
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsFolder])
			},
		},
		{
			Name:    "FourTestsSixteenFiles",
			Dir:     ".",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.Contains(t, report.Details[testing.RuleHasTests], "Your project contains **4** test files")
				require.Contains(t, report.Details[testing.RuleHasTests], "meets the minimum of **1** test files required")
				require.Contains(t, report.Details[testing.RuleHasTests], "equates to **20%** of Python files in your project being tests")
				require.Contains(t, report.Details[testing.RuleHasTests], "meets the target ratio of **20%**")
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsPass])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestsFolder])
			},
		},
		{
			Name: "FourTestsSixteenFiles/InTestsFolder",
			Dir:  ".",
			Options: testutils.NewOptions().UsePythonFiles(
				createPythonFilenames(16).Concat(createPythonTestFilenames(4).Prefix("tests")),
			),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsFolder])
			},
		},
		{
			Name: "FourTestsSixteenFiles/HalfInTestsFolder",
			Dir:  ".",
			Options: testutils.NewOptions().UsePythonFiles(
				createPythonFilenames(16).
					Concat(createPythonTestFilenames(2).Prefix("tests")).
					Concat(createPythonTestFilenames(2)),
			),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 50, report.Scores[testing.RuleTestsFolder])
				require.Contains(t, report.Details[testing.RuleTestsFolder], "The following test files have been detected that are **not** in the `tests` folder at the root of your project")
				require.Contains(t, report.Details[testing.RuleTestsFolder], "- file0_test.py")
				require.Contains(t, report.Details[testing.RuleTestsFolder], "- file1_test.py")
			},
		},
		{
			Name: "FourTestsSixteenFiles/InTestsFolderAbsolute",
			Dir:  utils.AbsolutePath("."),
			Options: testutils.NewOptions().UsePythonFiles(
				createPythonFilenames(16).Concat(createPythonTestFilenames(4).Prefix(path.Join(utils.AbsolutePath("."), "tests"))),
			),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsFolder])
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassed/ZeroCoverage",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage.Report = "coverage-0.xml"
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 0, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "achieved **0.0%** line test coverage")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "**80.0%** is the target amount of test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassedHalfCoverage",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage.Report = "coverage-50.xml"
					c.Testing.Coverage.Targets.Line = 100
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 50, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "achieved **50.0%** line test coverage")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "**100.0%** is the target amount of test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassedHalfCoverageButTargetAchieved",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage.Report = "coverage-50.xml"
					c.Testing.Coverage.Targets.Line = 50
					return c
				}()),
			Expect: func(t *stdtesting.T, report api.Report, err error) {
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[testing.RuleHasTests])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestsPass])
				require.Contains(t, report.Details[testing.RuleTestsPass], "all **4** tests in your project passed")
				// require.Equal(t, 0, report.Scores[testing.RuleTestsFolder])
				require.EqualValues(t, 100, report.Scores[testing.RuleTestCoverage])
				require.Contains(t, report.Details[testing.RuleTestCoverage], "achieved **50.0%** line test coverage")
				require.Contains(t, report.Details[testing.RuleTestCoverage], "which meets the target of **50.0%** test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllPassedFullCoverage",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-passed-all.xml"
					c.Testing.Coverage.Report = "coverage-100.xml"
					c.Testing.Coverage.Targets.Line = 100
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
				require.Contains(t, report.Details[testing.RuleTestCoverage], "**100%** line test coverage")
			},
		},
		{
			Name: "FourTestsSixteenFiles/AllFailed",
			Dir:  "test-resources",
			Options: testutils.NewOptions().UsePythonFiles(createPythonFilenames(16).Concat(createPythonTestFilenames(4))).
				WithConfig(func() *config.Config {
					c := config.Default()
					c.Testing.Report = "junit-failed-all.xml"
					c.Testing.Coverage.Report = ""
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
					c.Testing.Coverage.Report = "" // TODO
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
					c.Testing.Coverage.Report = "non-existant-file.xml"
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
					c.Testing.Coverage.Report = "coverage-malformed.xml"
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
					c.Testing.Coverage.Report = "coverage-empty.xml"
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
					c.Testing.Coverage.Report = "coverage-0.xml"
					c.Testing.Coverage.Targets.Line = 0
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

func TestTestingLinterConfigure(t *stdtesting.T) {
	linter := testing.NewLinter()
	conf := config.Default()
	require.NoError(t, linter.Configure(conf))

	conf.Testing.Coverage.Targets.Line = -1
	require.ErrorIs(t, linter.Configure(conf), testing.ErrCoverageTargetTooLow)
	conf.Testing.Coverage.Targets.Line = 200
	require.ErrorIs(t, linter.Configure(conf), testing.ErrCoverageTargetTooHigh)
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
