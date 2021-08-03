package api_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
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

func TestNewCustomRule(t *testing.T) {
	cr := config.CustomRule{
		Name:    "Custom Test Rule",
		Slug:    "custom/test-rule",
		Details: "Tests whether parsing a custom rule from a YAML config works",
		Weight:  420,
		Run:     "python ./scripts/mllint-test-rule.py",
	}
	rule := api.NewCustomRule(cr)

	require.Equal(t, cr.Name, rule.Name)
	require.Equal(t, cr.Slug, rule.Slug)
	require.Equal(t, cr.Details, rule.Details)
	require.Equal(t, cr.Weight, rule.Weight)
	require.False(t, rule.Disabled)
}

func TestReportOverallScore(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		report := api.NewReport()
		require.Equal(t, 0.0, report.OverallScore())
	})

	t.Run("SameWeight", func(t *testing.T) {
		report := api.NewReport()
		rule1 := api.Rule{Slug: "test-1", Weight: 1}
		rule2 := api.Rule{Slug: "test-2", Weight: 1}
		rule3 := api.Rule{Slug: "test-3", Weight: 1}
		rule4 := api.Rule{Slug: "test-4", Weight: 1}

		report.Scores[rule1] = 100
		report.Scores[rule2] = 30
		report.Scores[rule3] = 70
		report.Scores[rule4] = 40
		require.Equal(t, 60.0, report.OverallScore())
	})

	t.Run("Weighted", func(t *testing.T) {
		report := api.NewReport()
		rule1 := api.Rule{Slug: "test-1", Weight: 1}
		rule2 := api.Rule{Slug: "test-2", Weight: 2}
		rule3 := api.Rule{Slug: "test-3", Weight: 3}
		rule4 := api.Rule{Slug: "test-4", Weight: 4}

		report.Scores[rule1] = 100
		report.Scores[rule2] = 50
		report.Scores[rule3] = 70
		report.Scores[rule4] = 100
		require.Equal(t, 81.0, report.OverallScore())
	})
}
