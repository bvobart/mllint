package linters_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/api/categories"
	"github.com/bvobart/mllint/api/mock_api"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/linters/common"
)

func TestDisableIgnore(t *testing.T) {
	linters.ByCategory = map[api.Category]api.Linter{}
	linters.Disable("something")
}

func TestDisableAllCategories(t *testing.T) {
	mockctl := gomock.NewController(t)
	mockLinters := [4]*mock_api.MockLinter{
		mock_api.NewMockLinter(mockctl),
		mock_api.NewMockLinter(mockctl),
		mock_api.NewMockLinter(mockctl),
		mock_api.NewMockLinter(mockctl),
	}
	mockLinters[0].EXPECT().Rules().Times(1).Return([]*api.Rule{})

	linters.ByCategory = map[api.Category]api.Linter{
		categories.VersionControl: mockLinters[0],
		categories.CodeQuality:    mockLinters[1],
		categories.DataQuality:    mockLinters[2],
		categories.Testing:        mockLinters[3],
	}

	disabled := []string{
		categories.CodeQuality.Slug(), // category exists
		categories.DataQuality.Slug(), // category exists
		categories.VersionControl.Slug() + "/some/specific-rule",
		categories.DependencyMgmt.Slug(), // category not in linters.ByCategory
	}
	linters.DisableAll(disabled)

	enabled := map[api.Category]api.Linter{
		categories.VersionControl: mockLinters[0],
		categories.Testing:        mockLinters[3],
	}
	require.Equal(t, enabled, linters.ByCategory)
}

func TestConfigureAll(t *testing.T) {
	mockctl := gomock.NewController(t)
	testLinter1 := mock_api.NewMockConfigurableLinter(mockctl)
	testLinter2 := mock_api.NewMockConfigurableLinter(mockctl)
	testLinter3 := mock_api.NewMockConfigurableLinter(mockctl)

	linters.ByCategory = map[api.Category]api.Linter{
		"test1": testLinter1,
		"test2": testLinter2,
		"test3": testLinter3,
	}

	conf := config.Default()
	testLinter1.EXPECT().Configure(conf).Times(1).Return(nil)
	testLinter2.EXPECT().Configure(conf).Times(1).Return(nil)
	testLinter3.EXPECT().Configure(conf).Times(1).Return(nil)

	require.NoError(t, linters.ConfigureAll(conf))
}

func TestConfigureAllError(t *testing.T) {
	mockctl := gomock.NewController(t)
	testLinter1 := mock_api.NewMockConfigurableLinter(mockctl)
	testLinter2 := mock_api.NewMockConfigurableLinter(mockctl)
	testLinter3 := mock_api.NewMockConfigurableLinter(mockctl)

	linters.ByCategory = map[api.Category]api.Linter{
		"test1": testLinter1,
		"test2": testLinter2,
		"test3": testLinter3,
	}

	conf := config.Default()
	testErr := errors.New("test-error")
	testLinter1.EXPECT().Configure(conf).Times(1).Return(nil)
	testLinter2.EXPECT().Configure(conf).Times(1).Return(testErr)
	testLinter3.EXPECT().Configure(conf).Times(0).Return(nil)

	err := linters.ConfigureAll(conf)
	require.Error(t, err)
	require.ErrorIs(t, err, testErr)
}

func TestDisableRuleNormalLinter(t *testing.T) {
	mockctl := gomock.NewController(t)
	mockLinter := mock_api.NewMockLinter(mockctl)
	mockRules := []*api.Rule{
		{Slug: "test-rule-1"},
		{Slug: "test-rule-2"},
		{Slug: "test-rule-3"},
	}
	mockLinter.EXPECT().Rules().Times(1).Return(mockRules)

	linters.DisableRule(mockLinter, "test-rule-2")
	require.False(t, mockRules[0].Disabled)
	require.True(t, mockRules[1].Disabled)
	require.False(t, mockRules[2].Disabled)
}

func TestDisableRuleCompositeLinter(t *testing.T) {
	mockctl := gomock.NewController(t)
	mockLinter1 := mock_api.NewMockLinter(mockctl)
	mockLinter2 := mock_api.NewMockLinter(mockctl)
	mockRules := []*api.Rule{
		{Slug: "test-rule-1"},
		{Slug: "test-rule-2"},
		{Slug: "test-rule-1"},
		{Slug: "test-rule-2"},
	}
	mockLinter1.EXPECT().Name().Times(1).Return("Mock 1")
	mockLinter2.EXPECT().Name().Times(1).Return("Mock 2")
	mockLinter1.EXPECT().Rules().Times(2).Return(mockRules[:2])
	mockLinter2.EXPECT().Rules().Times(2).Return(mockRules[2:])

	compLinter := common.NewCompositeLinter("Testing", mockLinter1, mockLinter2)
	linters.DisableRule(compLinter, "mock-1/test-rule-2")
	linters.DisableRule(compLinter, "mock-2/test-rule-1")

	require.False(t, mockRules[0].Disabled)
	require.True(t, mockRules[1].Disabled)
	require.True(t, mockRules[2].Disabled)
	require.False(t, mockRules[3].Disabled)
}
