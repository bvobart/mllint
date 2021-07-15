package custom_test

import (
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters/custom"
	"github.com/bvobart/mllint/linters/testutils"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

func TestCustomLinter(t *testing.T) {
	linter := custom.NewLinter()
	require.Equal(t, "Custom", linter.Name())
	require.Equal(t, []*api.Rule(nil), linter.Rules())

	suite := testutils.NewLinterTestSuite(linter, []testutils.LinterTest{
		{
			Name: "SimpleEcho",
			Dir:  ".",
			Options: testutils.NewOptions().WithConfig(func() *config.Config {
				conf := config.Default()
				conf.Rules.Custom = append(conf.Rules.Custom, createCustomRule1())
				return conf
			}()),
			Expect: func(t *testing.T, report api.Report, err error) {
				rule1 := api.NewCustomRule(createCustomRule1())

				require.Equal(t, []*api.Rule{&rule1}, linter.Rules())
				require.NoError(t, err)
				require.EqualValues(t, 100, report.Scores[rule1])
				require.Equal(t, "All good!", report.Details[rule1])
			},
		},
		{
			Name: "EmptyJSON",
			Dir:  ".",
			Options: testutils.NewOptions().WithConfig(func() *config.Config {
				conf := config.Default()
				conf.Rules.Custom = append(conf.Rules.Custom, createCustomRule2())
				return conf
			}()),
			Expect: func(t *testing.T, report api.Report, err error) {
				rule2 := api.NewCustomRule(createCustomRule2())

				require.Equal(t, []*api.Rule{&rule2}, linter.Rules())
				require.NoError(t, err)
				require.EqualValues(t, 0, report.Scores[rule2])
				require.Equal(t, "", report.Details[rule2])
			},
		},
		{
			Name: "ErrorsInCommandOrExecution",
			Dir:  ".",
			Options: testutils.NewOptions().WithConfig(func() *config.Config {
				conf := config.Default()
				conf.Rules.Custom = append(conf.Rules.Custom, createErrorRule1())
				conf.Rules.Custom = append(conf.Rules.Custom, createErrorRule2())
				conf.Rules.Custom = append(conf.Rules.Custom, createErrorRule3())
				conf.Rules.Custom = append(conf.Rules.Custom, createErrorRule4())
				conf.Rules.Custom = append(conf.Rules.Custom, createErrorRule5())
				return conf
			}()),
			Expect: func(t *testing.T, report api.Report, err error) {
				customRule1 := createErrorRule1()
				customRule2 := createErrorRule2()
				customRule3 := createErrorRule3()
				customRule4 := createErrorRule4()
				customRule5 := createErrorRule5()
				rule1 := api.NewCustomRule(customRule1)
				rule2 := api.NewCustomRule(customRule2)
				rule3 := api.NewCustomRule(customRule3)
				rule4 := api.NewCustomRule(customRule4)
				rule5 := api.NewCustomRule(customRule5)

				require.Equal(t, []*api.Rule{&rule1, &rule2, &rule3, &rule4, &rule5}, linter.Rules())
				require.Error(t, err)
				require.IsType(t, &multierror.Error{}, err)
				require.Contains(t, err.Error(), "5 errors occurred:")
				require.Contains(t, err.Error(), "custom rule `custom/error-rule-1` has invalid run command")
				require.Contains(t, err.Error(), customRule1.Run)
				require.Contains(t, err.Error(), "custom rule `custom/error-rule-2` was run, but exited with an error: exec: \"some-executable-that-will-definitely-not-exist\": executable file not found in $PATH")
				require.Contains(t, err.Error(), "custom rule `custom/error-rule-3` was run, but exited with an error: exit status 1")
				require.Contains(t, err.Error(), "custom rule `custom/error-rule-4` was run, but exited with an error: exit status 1")
				require.Contains(t, err.Error(), "```\ndate: invalid option -- 'e'\nTry 'date --help' for more information.\n```")
				require.Contains(t, err.Error(), "custom rule `custom/error-rule-5` executed successfully, but the output was not a valid JSON / YAML object: yaml: did not find expected key. Output: `score 100, details: \"\" }`")
			},
		},
	})
	suite.DefaultOptions().WithConfig(config.Default())
	suite.RunAll(t)
}

func createCustomRule1() config.CustomRule {
	return config.CustomRule{
		Name:    "Custom Rule 1",
		Slug:    "custom/rule-1",
		Details: "This rule just echoes a correct JSON structure with a score of 100 and some details.",
		Weight:  420,
		Run:     `echo '{ score: 100, details: "All good!" }'`,
	}
}

func createCustomRule2() config.CustomRule {
	return config.CustomRule{
		Name:    "Custom Rule 2",
		Slug:    "custom/rule-2",
		Details: "This rule's command returns an empty JSON object",
		Weight:  420,
		Run:     `echo "{}"`,
	}
}

func createErrorRule1() config.CustomRule {
	return config.CustomRule{
		Name:    "Error Rule 1",
		Slug:    "custom/error-rule-1",
		Details: "This rule's command is invalid because a thief stole this closing quote from the argument to echo --> ' ",
		Weight:  420,
		Run:     `echo '{ score: 100, details: "All good!" }`,
	}
}

func createErrorRule2() config.CustomRule {
	return config.CustomRule{
		Name:    "Error Rule 2",
		Slug:    "custom/error-rule-2",
		Details: "This rule's command exits with an error because the executable cannot be found",
		Weight:  420,
		Run:     `some-executable-that-will-definitely-not-exist`,
	}
}

func createErrorRule3() config.CustomRule {
	return config.CustomRule{
		Name:    "Error Rule 3",
		Slug:    "custom/error-rule-3",
		Details: "This rule's command exits with a non-zero exit code as that is all that `false` does",
		Weight:  420,
		Run:     `false`,
	}
}

func createErrorRule4() config.CustomRule {
	return config.CustomRule{
		Name:    "Error Rule 4",
		Slug:    "custom/error-rule-4",
		Details: "This rule's command contains `&&` which is a shell feature, but os/exec just passes it and everything afterwards as an argument to date. To fix, use e.g. `bash -c 'your && shell || command here`",
		Weight:  420,
		Run:     `date && echo -e 'score: 100\ndetails: "sliated"'`,
	}
}

func createErrorRule5() config.CustomRule {
	return config.CustomRule{
		Name:    "Error Rule 5",
		Slug:    "custom/error-rule-5",
		Details: "This rule's command returns invalid JSON (missing colon after `score`)",
		Weight:  420,
		Run:     `echo 'score 100, details: "" }'`,
	}
}
