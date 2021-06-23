package testing

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/bvobart/gocover-cobertura/cobertura"
	"github.com/joshdk/go-junit"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/utils"
	"github.com/bvobart/mllint/utils/markdowngen"
)

var ErrCoverageTargetTooHigh = errors.New("coverage target higher than 100%")
var ErrCoverageTargetTooLow = errors.New("coverage target lower than 0%")

func NewLinter() api.ConfigurableLinter {
	return &TestingLinter{}
}

type TestingLinter struct {
	Config    config.TestingConfig
	TestFiles utils.Filenames
}

func (l *TestingLinter) Name() string {
	return categories.Testing.Name
}

func (l *TestingLinter) Configure(conf *config.Config) error {
	l.Config = conf.Testing
	if l.Config.CoverageTarget > 100 {
		return fmt.Errorf("%w: %.1f", ErrCoverageTargetTooHigh, l.Config.CoverageTarget)
	} else if l.Config.CoverageTarget < 0 {
		return fmt.Errorf("%w: %.1f", ErrCoverageTargetTooLow, l.Config.CoverageTarget)
	}
	return nil
}

func (l *TestingLinter) Rules() []*api.Rule {
	return []*api.Rule{&RuleHasTests, &RuleTestsPass, &RuleTestCoverage, &RuleTestsFolder}
}

func (l *TestingLinter) LintProject(project api.Project) (api.Report, error) {
	report := api.NewReport()

	l.TestFiles = project.PythonFiles.Filter(func(filename string) bool {
		return strings.HasSuffix(filename, "_test.py") || strings.HasPrefix(path.Base(filename), "test_")
	})

	l.ScoreRuleHasTests(&report, project)
	l.ScoreRuleTestsFolder(&report, project)
	l.ScoreRuleTestsPass(&report, project)
	l.ScoreRuleTestCoverage(&report, project)

	// TODO: determine possible config options:
	// - target amount of tests per file

	return report, nil
}

//---------------------------------------------------------------------------------------

func (l *TestingLinter) ScoreRuleHasTests(report *api.Report, project api.Project) {
	if len(project.PythonFiles) == 0 {
		report.Scores[RuleHasTests] = 0
		return
	}

	// Possible TODO: have a config option for the target amount of tests per file?

	// there should be at 1 test file per 4 non-test Python files.
	report.Scores[RuleHasTests] = 100 * (float64(len(l.TestFiles)) * 4 / float64(len(project.PythonFiles)-len(l.TestFiles)))
}

//---------------------------------------------------------------------------------------

func (l *TestingLinter) ScoreRuleTestsFolder(report *api.Report, project api.Project) {
	if len(project.PythonFiles) == 0 {
		report.Scores[RuleTestsFolder] = 0
		return
	}

	if len(l.TestFiles) == 0 {
		if utils.FolderExists(path.Join(project.Dir, "tests")) {
			report.Scores[RuleTestsFolder] = 100
			report.Details[RuleTestsFolder] = "While no tests were detected in your project, it's good that your project already has a `tests` folder!"
		} else {
			report.Scores[RuleTestsFolder] = 0
			report.Details[RuleTestsFolder] = "Tip for when you start implementing tests: create a folder called `tests` at the root of your project and place all your Python test files in there, as per common convention."
		}
		return
	}

	notInTestsFolder := utils.Filenames{}
	for _, testFile := range l.TestFiles {
		if !isInTestsFolder(project.Dir, testFile) {
			notInTestsFolder = append(notInTestsFolder, testFile)
		}
	}

	// score is percentage of test files that _are_ in the tests folder.
	report.Scores[RuleTestsFolder] = 100 * (1 - float64(len(notInTestsFolder))/float64(len(l.TestFiles)))
	if len(notInTestsFolder) > 0 {
		report.Details[RuleTestsFolder] = "The following test files have been detected that are **not** in the `tests` folder at the root of your project:\n\n" +
			markdowngen.ListFiles(notInTestsFolder)
	}
}

func isInTestsFolder(projectdir, testFile string) bool {
	// files passed into a linter through the project are generally absolute paths.
	if path.IsAbs(testFile) {
		return strings.HasPrefix(testFile, path.Join(projectdir, "tests"))
	}

	// if the path is not absolute, it is assumed to be relative to the project root.
	return strings.HasPrefix(testFile, "tests")
}

//---------------------------------------------------------------------------------------

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

//---------------------------------------------------------------------------------------

func (l *TestingLinter) ScoreRuleTestCoverage(report *api.Report, project api.Project) {
	if l.Config.Coverage == "" {
		report.Scores[RuleTestCoverage] = 0
		report.Details[RuleTestCoverage] = "No test coverage report was provided. Please update the `testing.coverage` setting in your project's `mllint` configuration to specify the path to your project's test coverage report.\n\n" + howToMakeCoverageXML
		return
	}

	covReportFile, err := utils.OpenFile(project.Dir, l.Config.Coverage)
	if err != nil {
		report.Scores[RuleTestCoverage] = 0
		report.Details[RuleTestCoverage] = fmt.Sprintf("A test coverage report was provided, namely `%s`, but this file could not be found or opened (%s). Please update the `testing.coverage` setting in your project's `mllint` configuration to fix the path to your project's test report. Remember that this path must be relative to the root of your project directory.", l.Config.Coverage, err.Error())
		return
	}

	var covReport cobertura.Coverage
	covReportData, err := io.ReadAll(covReportFile)
	if err == nil {
		err = xml.Unmarshal(covReportData, &covReport)
	}
	if err != nil {
		report.Scores[RuleTestCoverage] = 0
		report.Details[RuleTestCoverage] = fmt.Sprintf(`A test report file `+"`%s`"+` was provided and found, but there was an error parsing the Cobertura XML contents:

%s

Please make sure your test report file is a valid Cobertura-compatible XML file. %s`, l.Config.Report, "```\n"+err.Error()+"\n```", howToMakeCoverageXML)
		return
	}

	totalLines := covReport.NumLines()
	hitLines := covReport.NumLinesWithHits()
	hitRate := 100 * float64(hitLines) / float64(totalLines) // percentage of lines covered.
	score := 100 * hitRate / l.Config.CoverageTarget         // percentage of coverage target achieved.
	if totalLines == 0 {
		score = 0
	}
	if l.Config.CoverageTarget == 0 {
		score = 100
	}
	report.Scores[RuleTestCoverage] = score

	if totalLines != 0 && hitLines == totalLines {
		report.Details[RuleTestCoverage] = "Wow! Congratulations! You've achieved full 100% line test coverage! Great job!"
	} else if hitRate < l.Config.CoverageTarget {
		report.Details[RuleTestCoverage] = fmt.Sprintf("Your project's tests achieved %.1f%% line test coverage, but %.1f%% is the target amount of test coverage to beat. You'll need to further improve your tests.", hitRate, l.Config.CoverageTarget)
	} else if hitRate >= l.Config.CoverageTarget {
		report.Details[RuleTestCoverage] = fmt.Sprintf("Congratulations, your project's tests have achieved %.1f%% line test coverage, which meets the target of %.1f%% test coverage!", hitRate, l.Config.CoverageTarget)
	} else if totalLines == 0 {
		report.Details[RuleTestCoverage] = "It seems your test coverage report is empty, no lines were covered."
	}
}

//---------------------------------------------------------------------------------------

const howToMakeJUnitXML = "When using `pytest` to run your project's tests, use the `--junitxml=<filename>` option to generate such a test report, e.g.: `pytest --junitxml=tests-report.xml`"
const howToMakeCoverageXML = "Generating a test coverage report with `pytest` can be done by adding and installing `pytest-cov` as a development dependency of your project. Then use the following command to run your tests and generate both a test report as well as a coverage report:" + `
` + "```" + `
pytest --junitxml=tests-report.xml --cov=path_to_package_under_test --cov-report=xml
` + "```\n"
