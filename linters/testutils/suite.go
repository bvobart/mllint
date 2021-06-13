package testutils

import (
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/setools/cqlinters"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/bvobart/mllint/utils"
	"github.com/stretchr/testify/require"
)

type LinterTest struct {
	Name    string
	Dir     string
	Expect  func(report api.Report, err error)
	Options *LinterTestOptions
}

type LinterTestSuite struct {
	linter      api.Linter
	tests       []LinterTest
	defaultOpts *LinterTestOptions
}

// NewLinterTestSuite initialises a test suite for a specific linter, with a list of tests
// that will be executed in parallel when suite.RunAll(t) is called.
func NewLinterTestSuite(linter api.Linter, tests []LinterTest) *LinterTestSuite {
	return &LinterTestSuite{linter, tests, NewOptions()}
}

// DefaultOptions returns a pointer to the options object that will be applied with every test (unless overridden by test options)
func (suite *LinterTestSuite) DefaultOptions() *LinterTestOptions {
	return suite.defaultOpts
}

// RunAll runs all tests in the suite in parallel.
// For each test, it creates a project with the Dir and options specified in the LinterTest.
// Then, it runs the linter's LintProject with that project and calls the test's Expect function.
func (suite *LinterTestSuite) RunAll(t *testing.T) {
	for _, tt := range suite.tests {
		test := tt
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			project := api.Project{Dir: test.Dir}
			suite.applyOptions(t, test.Options, &project)

			report, err := suite.linter.LintProject(project)
			test.Expect(report, err)
		})
	}
}

//---------------------------------------------------------------------------------------

// applies the default and test's options to the project that will be passed to LintProject.
func (suite *LinterTestSuite) applyOptions(t *testing.T, testOptions *LinterTestOptions, project *api.Project) {
	suite.applyPythonFilesOptions(t, testOptions, project)
	suite.applyDepManagerOptions(t, testOptions, project)
	suite.applyCQLinterOptions(t, testOptions, project)
	suite.applyConfigOption(t, testOptions)
}

//---------------------------------------------------------------------------------------

func (suite *LinterTestSuite) applyPythonFilesOptions(t *testing.T, testOptions *LinterTestOptions, project *api.Project) {
	if testOptions != nil && len(testOptions.usePythonFiles) > 0 {
		project.PythonFiles = testOptions.usePythonFiles
		return
	}

	if suite.defaultOpts != nil && len(suite.defaultOpts.usePythonFiles) > 0 {
		project.PythonFiles = suite.defaultOpts.usePythonFiles
		return
	}

	if suite.defaultOpts != nil && suite.defaultOpts.detectPythonFiles || testOptions != nil && testOptions.detectPythonFiles {
		pyfiles, err := utils.FindPythonFilesIn(project.Dir)
		require.NoError(t, err, "failed to parse Python files in test project")
		project.PythonFiles = pyfiles.Prefix(project.Dir)
	}
}

func (suite *LinterTestSuite) applyDepManagerOptions(t *testing.T, testOptions *LinterTestOptions, project *api.Project) {
	if testOptions != nil && len(testOptions.useDepManagers) > 0 {
		project.DepManagers = testOptions.useDepManagers
		return
	}

	if suite.defaultOpts != nil && len(suite.defaultOpts.useDepManagers) > 0 {
		project.DepManagers = suite.defaultOpts.useDepManagers
		return
	}

	if suite.defaultOpts != nil && suite.defaultOpts.detectDepManagers || testOptions != nil && testOptions.detectDepManagers {
		project.DepManagers = depmanagers.Detect(*project)
	}
}

func (suite *LinterTestSuite) applyCQLinterOptions(t *testing.T, testOptions *LinterTestOptions, project *api.Project) {
	if testOptions != nil && len(testOptions.useCQLinters) > 0 {
		project.CQLinters = testOptions.useCQLinters
		return
	}

	if suite.defaultOpts != nil && len(suite.defaultOpts.useCQLinters) > 0 {
		project.CQLinters = suite.defaultOpts.useCQLinters
		return
	}

	if suite.defaultOpts != nil && suite.defaultOpts.detectCQLinters || testOptions != nil && testOptions.detectCQLinters {
		project.CQLinters = cqlinters.Detect(*project)
	}
}

func (suite *LinterTestSuite) applyConfigOption(t *testing.T, testOptions *LinterTestOptions) {
	if configurable, ok := suite.linter.(api.ConfigurableLinter); ok {
		if testOptions != nil && testOptions.conf != nil {
			require.NoError(t, configurable.Configure(testOptions.conf), "error configuring test linter with configuration from test options")
			return
		}

		if suite.defaultOpts != nil && suite.defaultOpts.conf != nil {
			require.NoError(t, configurable.Configure(suite.defaultOpts.conf), "error configuring test linter with configuration from default options")
			return
		}
	}
}
