package testing

import (
	"fmt"
	"path"
	"strings"

	"github.com/joshdk/go-junit"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/utils"
)

func NewLinter() api.ConfigurableLinter {
	return &TestingLinter{}
}

type TestingLinter struct {
	Config config.TestingConfig
}

func (l *TestingLinter) Name() string {
	return categories.Testing.Name
}

func (l *TestingLinter) Configure(conf *config.Config) error {
	l.Config = conf.Testing
	return nil
}

func (l *TestingLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleHasTests, &RuleTestsPass, &RuleTestCoverage, &RuleTestsFolder}
}

func (l *TestingLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()

	l.ScoreRuleHasTests(&report, project)
	l.ScoreRuleTestsPass(&report, project)

	// TODO: implement the linting for RuleTestCoverage, which checks whether there is a Cobertura XML coverage report and analyses it for test coverage.
	// TODO: check whether all test files are in tests folder.
	// TODO: determine possible config options:
	// - target amount of tests per file
	// - target test coverage

	return report, nil
}

func (l *TestingLinter) ScoreRuleHasTests(report *api.Report, project api.Project) {
	if len(project.PythonFiles) == 0 {
		report.Scores[RuleHasTests] = 0
		return
	}

	testFiles := project.PythonFiles.Filter(func(filename string) bool {
		return strings.HasSuffix(filename, "_test.py") || strings.HasPrefix(path.Base(filename), "test_")
	})

	// Possible TODO: have a config option for the target amount of tests per file?

	// there should be at 1 test file per 4 non-test Python files.
	report.Scores[RuleHasTests] = 100 * (float64(len(testFiles)) * 4 / float64(len(project.PythonFiles)-len(testFiles)))
}

func (l *TestingLinter) ScoreRuleTestsPass(report *api.Report, project api.Project) {
	if l.Config.Report == "" {
		report.Scores[RuleTestsPass] = 0
		report.Details[RuleTestsPass] = "No test report was provided. Please update the `testing.report` setting in your project's `mllint` configuration to specify the path to your project's test report.\n\n" + howToMakeJUnitXML
		return
	}

	junitReportPath := path.Join(project.Dir, l.Config.Report)
	if !utils.FileExists(junitReportPath) {
		report.Scores[RuleTestsPass] = 0
		report.Details[RuleTestsPass] = fmt.Sprintf("A test report was provided, namely `%s`, but this file could not be found. Please update the `testing.report` setting in your project's `mllint` configuration to fix the path to your project's test report. Remember that this path must be relative to the root of your project directory.", l.Config.Report)
		return
	}

	suites, err := junit.IngestFile(junitReportPath)
	if err != nil {
		report.Scores[RuleTestsPass] = 0
		report.Details[RuleTestsPass] = fmt.Sprintf(`A test report file `+"`%s`"+` was provided and found, but there was an error parsing the JUnit XML contents:

%s

Please make sure your test report file is a valid JUnit XML file. %s`, l.Config.Report, "```\n"+err.Error()+"\n```", howToMakeJUnitXML)
		return
	}

	passedTests := 0
	totalTests := 0
	for _, suite := range suites {
		totalTests += suite.Totals.Tests
		passedTests += suite.Totals.Passed
	}

	if totalTests == 0 {
		report.Scores[RuleTestsPass] = 0
		report.Details[RuleTestsPass] = fmt.Sprintf(`No tests were run, according to the provided test report file `+"`%s`"+`. Don't be shy, implement some tests!`, l.Config.Report)
		return
	}

	score := 100 * float64(passedTests) / float64(totalTests)
	report.Scores[RuleTestsPass] = score
	if passedTests == totalTests {
		report.Details[RuleTestsPass] = fmt.Sprintf("Congratulations, all **%d** tests in your project passed!", totalTests)
	} else if passedTests == 0 {
		report.Details[RuleTestsPass] = fmt.Sprintf("Oh no! What a shame... **None** of the %d tests in your project passed! There must be something terribly wrong.", totalTests)
	} else if score < 0.25 {
		report.Details[RuleTestsPass] = fmt.Sprintf("Oh no! Only **%d** out of **%d** tests in your project passed... That's less than a quarter of all your project's tests...", passedTests, totalTests)
	} else if score > 0.75 {
		report.Details[RuleTestsPass] = fmt.Sprintf("Hmm, only **%d** out of **%d** tests in your project passed... That's over three quarter of all tests in your project, but it's not enough: _all tests must pass_. Good luck fixing the broken tests!", passedTests, totalTests)
	} else {
		report.Details[RuleTestsPass] = fmt.Sprintf("Oh my, only **%d** out of **%d** tests in your project passed... You can do better, right? Good luck fixing those tests!", passedTests, totalTests)
	}
}

const howToMakeJUnitXML = "When using `pytest` to run your project's tests, use the `--junitxml=<filename>` option to generate such a test report, e.g.: `pytest --junitxml=tests-report.xml`"
