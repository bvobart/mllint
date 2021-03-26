package api_test

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/config"
)

func TestIssuesFilterEnabled(t *testing.T) {
	issues := api.IssueList{
		api.NewIssue("testLinter", "0", api.SeverityError, "test message"),
		api.NewIssue("testLinter", "1", api.SeverityError, "test message"),
		api.NewIssue("testLinter", "2", api.SeverityError, "test message"),
		api.NewIssue("testLinter", "3", api.SeverityError, "test message"),
	}
	ruleconf := config.RuleConfig{Disabled: []string{"testLinter/0", "testLinter/2"}}
	actualEnabled := issues.FilterEnabled(ruleconf)

	enabled := api.IssueList{issues[1], issues[3]}
	require.Equal(t, enabled, actualEnabled)
}

func TestIssuesFilterEnabledDifferentLinters(t *testing.T) {
	issues := api.IssueList{
		api.NewIssue("testLinter", "0", api.SeverityError, "test message"),
		api.NewIssue("testLinter", "1", api.SeverityError, "test message"),
		api.NewIssue("fakeLinter", "2", api.SeverityError, "test message"),
		api.NewIssue("fakeLinter", "3", api.SeverityError, "test message"),
		api.NewIssue("singleLinter", "", api.SeverityError, "test message"),
	}
	ruleconf := config.RuleConfig{Disabled: []string{"testLinter/0", "fakeLinter/2", "singleLinter"}}
	actualEnabled := issues.FilterEnabled(ruleconf)

	enabled := api.IssueList{issues[1], issues[3]}
	require.Equal(t, enabled, actualEnabled)
}

func TestSeverityString(t *testing.T) {
	require.Equal(t, color.RedString(string(api.SeverityError)), api.SeverityError.String())
	require.Equal(t, color.YellowString(string(api.SeverityWarning)), api.SeverityWarning.String())
	require.Equal(t, color.BlueString(string(api.SeverityInfo)), api.SeverityInfo.String())
	require.Equal(t, "test", api.Severity("test").String())
}
