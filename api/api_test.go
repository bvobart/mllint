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
