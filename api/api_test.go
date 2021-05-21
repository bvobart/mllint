package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
)

func TestCategoryString(t *testing.T) {
	cat := categories.Testing
	require.Equal(t, cat.Name, cat.String())
}

func TestRuleEnableDisable(t *testing.T) {
	rule := api.Rule{
		Slug:     "test-rule",
		Name:     "Test Rule",
		Details:  "Some details about this rule",
		Weight:   1,
		Disabled: false,
	}
	rule.Disable()
	require.True(t, rule.Disabled)
	rule.Enable()
	require.False(t, rule.Disabled)
}

func TestNewReport(t *testing.T) {
	report := api.NewReport()
	require.NotNil(t, report.Scores)
	require.NotNil(t, report.Details)
}

func TestMergeReports(t *testing.T) {
	rule1 := api.Rule{Slug: "test-1"}
	rule2 := api.Rule{Slug: "test-2"}
	rule3 := api.Rule{Slug: "test-3"}

	report1 := api.NewReport()

	report2 := api.NewReport()
	report2.Scores[rule1] = 100
	report2.Details[rule1] = "something"
	report2.Scores[rule2] = 65
	report2.Details[rule2] = "something else"

	finalReport := api.MergeReports(report1, report2)
	require.Equal(t, report2, finalReport)

	finalReport = api.MergeReports(report2, report2, report2)
	require.Equal(t, report2, finalReport)

	report1.Scores[rule3] = 42
	report1.Details[rule3] = "something completely different"

	finalReport = api.MergeReports(report1, report2)
	expectedReport := api.NewReport()
	expectedReport.Scores[rule1] = 100
	expectedReport.Details[rule1] = "something"
	expectedReport.Scores[rule2] = 65
	expectedReport.Details[rule2] = "something else"
	expectedReport.Scores[rule3] = 42
	expectedReport.Details[rule3] = "something completely different"

	require.Equal(t, expectedReport, finalReport)
}
