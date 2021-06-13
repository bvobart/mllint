package dependencymgmt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/linters/dependencymgmt"
	"github.com/bvobart/mllint/linters/testutils"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/bvobart/mllint/utils/markdowngen"
)

func TestName(t *testing.T) {
	linter := dependencymgmt.NewLinter()
	require.Equal(t, categories.DependencyMgmt.Name, linter.Name())
}

func TestRules(t *testing.T) {
	linter := dependencymgmt.NewLinter()
	require.Equal(t, []*api.Rule{&dependencymgmt.RuleUse, &dependencymgmt.RuleSingle, &dependencymgmt.RuleUseDev}, linter.Rules())
}

func TestLintProject(t *testing.T) {
	perfectScore := func(report api.Report, err error) {
		require.NoError(t, err)
		require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
		require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleSingle])
	}

	tests := []testutils.LinterTest{
		{Name: "Correct/Pipenv", Dir: "test-resources/correct-pipenv", Expect: perfectScore},
		{Name: "Correct/Poetry", Dir: "test-resources/correct-poetry", Expect: perfectScore},
		{Name: "Invalid/None", Dir: "test-resources/none", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
		}},
		{Name: "Invalid/RequirementsTxt", Dir: "test-resources/requirementstxt", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 20, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleUse], dependencymgmt.DetailsNoRequirementsTxt)
		}},
		{Name: "Invalid/SetupPy", Dir: "test-resources/setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 30, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleUse], dependencymgmt.DetailsNoSetupPy)
		}},
		{Name: "Invalid/Multiple/Pipenv+SetupPy", Dir: "test-resources/multiple/pipenv+setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsPipenvSetupPy)
		}},
		{Name: "Invalid/Multiple/RequirementsTxt+SetupPy", Dir: "test-resources/multiple/requirementstxt+setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 20, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsRequirementsTxtSetupPy)
		}},
		{Name: "Invalid/Multiple/Pipenv+RequirementsTxt", Dir: "test-resources/multiple/pipenv+requirementstxt", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsRequirementsTxtPipenv)
		}},
		{Name: "Invalid/Multiple/Poetry+RequirementsTxt", Dir: "test-resources/multiple/poetry+requirementstxt", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsRequirementsTxtPoetry)
		}},
		{Name: "Invalid/Multiple/Poetry+SetupPy", Dir: "test-resources/multiple/poetry+setuppy", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
			require.Contains(t, report.Details[dependencymgmt.RuleSingle], dependencymgmt.DetailsPoetrySetupPy)
		}},
		{Name: "Invalid/Multiple/Poetry+Pipenv", Dir: "test-resources/multiple/poetry+pipenv", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUse])
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleSingle])
		}},
		{Name: "Invalid/DevDeps/Poetry", Dir: "test-resources/dev-dependencies/poetry", Expect: func(report api.Report, err error) {
			require.NoError(t, err)
		}},
	}

	linter := dependencymgmt.NewLinter()
	suite := testutils.NewLinterTestSuite(linter, tests)
	suite.DefaultOptions().DetectDepManagers()
	suite.RunAll(t)
}

type ruleUseDevTest struct {
	Name        string
	Dir         string
	Expect      func(report api.Report)
	ManagerType api.DependencyManagerType
}

func TestRuleUseDev(t *testing.T) {
	tests := []ruleUseDevTest{
		{Name: "Poetry/Correct", Dir: "test-resources/dev-dependencies/poetry/correct", ManagerType: depmanagers.TypePoetry, Expect: func(report api.Report) {
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUseDev])
		}},
		{Name: "Poetry/Invalid", Dir: "test-resources/dev-dependencies/poetry/invalid", ManagerType: depmanagers.TypePoetry, Expect: func(report api.Report) {
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleUseDev])
			require.Contains(t, report.Details[dependencymgmt.RuleUseDev], markdowngen.List([]interface{}{"pytest", "tox"}))
			require.Contains(t, report.Details[dependencymgmt.RuleUseDev], "`pyproject.toml`")
			require.Contains(t, report.Details[dependencymgmt.RuleUseDev], "`poetry lock`")
		}},
		{Name: "Pipenv/Correct", Dir: "test-resources/dev-dependencies/pipenv/correct", ManagerType: depmanagers.TypePipenv, Expect: func(report api.Report) {
			require.EqualValues(t, 100, report.Scores[dependencymgmt.RuleUseDev])
		}},
		{Name: "Pipenv/Invalid", Dir: "test-resources/dev-dependencies/pipenv/invalid", ManagerType: depmanagers.TypePipenv, Expect: func(report api.Report) {
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleUseDev])
			require.Contains(t, report.Details[dependencymgmt.RuleUseDev], markdowngen.List([]interface{}{"dvc", "isort", "mypy"}))
		}},
		{Name: "RequirementsTxt", Dir: "test-resources/dev-dependencies/requirementstxt", ManagerType: depmanagers.TypeRequirementsTxt, Expect: func(report api.Report) {
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleUseDev])
			require.Equal(t, "Your project's main dependency manager is a `requirements.txt` file, which doesn't distinguish between regular dependencies and development dependencies.", report.Details[dependencymgmt.RuleUseDev])
		}},
		{Name: "SetupPy", Dir: "test-resources/dev-dependencies/setuppy", ManagerType: depmanagers.TypeSetupPy, Expect: func(report api.Report) {
			require.EqualValues(t, 0, report.Scores[dependencymgmt.RuleUseDev])
			require.Equal(t, "Your project's main dependency manager is a `setup.py` file, which doesn't distinguish between regular dependencies and development dependencies.", report.Details[dependencymgmt.RuleUseDev])
		}},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			project := api.Project{Dir: test.Dir}

			manager, err := test.ManagerType.Detect(project)
			require.NoError(t, err)

			report := api.NewReport()
			linter := dependencymgmt.NewLinter().(*dependencymgmt.DependenciesLinter)
			linter.ScoreRuleUseDev(&report, manager)
			test.Expect(report)
		})
	}
}
