package linters_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/api/mock_api"
	"github.com/bvobart/mllint/categories"
	"github.com/bvobart/mllint/config"
	"github.com/bvobart/mllint/linters"
	"github.com/bvobart/mllint/linters/common"
)

func TestDisableIgnore(t *testing.T) {
	linters.ByCategory = map[api.Category]api.Linter{}
	require.Equal(t, 0, linters.Disable("something"))
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
	mockLinters[1].EXPECT().Rules().Times(1).Return([]*api.Rule{{}})
	mockLinters[2].EXPECT().Rules().Times(1).Return([]*api.Rule{{}, {}})

	linters.ByCategory = map[api.Category]api.Linter{
		categories.VersionControl: mockLinters[0],
		categories.CodeQuality:    mockLinters[1],
		categories.DataQuality:    mockLinters[2],
		categories.Testing:        mockLinters[3],
	}

	disabled := []string{
		categories.CodeQuality.Slug, // category exists
		categories.DataQuality.Slug, // category exists
		categories.VersionControl.Slug + "/some/specific-rule",
		categories.DependencyMgmt.Slug, // category not in linters.ByCategory
	}
	require.Equal(t, 3, linters.DisableAll(disabled))

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
		{Name: "test1"}: testLinter1,
		{Name: "test2"}: testLinter2,
		{Name: "test3"}: testLinter3,
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
		{Name: "test1"}: testLinter1,
		{Name: "test2"}: testLinter2,
		{Name: "test3"}: testLinter3,
	}

	conf := config.Default()
	testErr := errors.New("test-error")
	testLinter1.EXPECT().Configure(conf).MaxTimes(1).Return(nil)
	testLinter2.EXPECT().Configure(conf).MaxTimes(1).Return(testErr)
	testLinter3.EXPECT().Configure(conf).MaxTimes(1).Return(nil)

	err := linters.ConfigureAll(conf)
	require.Error(t, err)
	require.ErrorIs(t, err, testErr)
}

func TestDisableNotExactlyMatchingCategory(t *testing.T) {
	linters.ByCategory = make(map[api.Category]api.Linter)
	require.Equal(t, 0, linters.Disable("version-control/"))
}

func TestDisableRuleNormalLinter(t *testing.T) {
	mockctl := gomock.NewController(t)
	mockLinter := mock_api.NewMockLinter(mockctl)
	mockRules := []*api.Rule{
		{Slug: "cat/test-rule-1"},
		{Slug: "cat/test-rule-2"},
		{Slug: "cat/test-rule-3"},
	}
	mockLinter.EXPECT().Rules().Times(1).Return(mockRules)

	require.Equal(t, 1, linters.DisableRule(mockLinter, "cat/test-rule-2"))
	require.False(t, mockRules[0].Disabled)
	require.True(t, mockRules[1].Disabled)
	require.False(t, mockRules[2].Disabled)
}

func TestDisableRuleCompositeLinter(t *testing.T) {
	mockctl := gomock.NewController(t)
	mockLinter1 := mock_api.NewMockLinter(mockctl)
	mockLinter2 := mock_api.NewMockLinter(mockctl)
	mockRules := []*api.Rule{
		{Slug: "mock-1/test-rule-1"},
		{Slug: "mock-1/test-rule-2"},
		{Slug: "mock-2/test-rule-1"},
		{Slug: "mock-2/test-rule-2"},
	}
	mockLinter1.EXPECT().Rules().Times(1).Return(mockRules[:2])
	mockLinter2.EXPECT().Rules().Times(1).Return(mockRules[2:])

	compLinter := common.NewCompositeLinter("Testing", mockLinter1, mockLinter2)
	require.Equal(t, 1, linters.DisableRule(compLinter, "mock-1/test-rule-2"))
	require.Equal(t, 1, linters.DisableRule(compLinter, "mock-2/test-rule-1"))

	require.False(t, mockRules[0].Disabled)
	require.True(t, mockRules[1].Disabled)
	require.True(t, mockRules[2].Disabled)
	require.False(t, mockRules[3].Disabled)
}

func TestDisableRulePartialSlug(t *testing.T) {
	mockctl := gomock.NewController(t)
	mockLinter := mock_api.NewMockLinter(mockctl)
	mockRules := []*api.Rule{
		{Slug: "cat/test-rule-1"},
		{Slug: "cat/test-rule-2"},
		{Slug: "cat/test-except-3"},
		{Slug: "cat/test-rule-4"},
	}
	mockLinter.EXPECT().Rules().Times(1).Return(mockRules)

	require.Equal(t, 3, linters.DisableRule(mockLinter, "cat/test-rule"))

	require.True(t, mockRules[0].Disabled)
	require.True(t, mockRules[1].Disabled)
	require.False(t, mockRules[2].Disabled)
	require.True(t, mockRules[3].Disabled)
}

func TestGetRule(t *testing.T) {
	rules1 := []*api.Rule{
		{Name: "Test Rule 1", Slug: "version-control/test-rule-1"},
		{Name: "Test Rule 2", Slug: "version-control/test-rule-2"},
	}
	rules2 := []*api.Rule{
		{Name: "Test Rule 3", Slug: "code-quality/test-rule-3"},
		{Name: "Test Rule 4", Slug: "code-quality/test-rule-4"},
	}

	mockctl := gomock.NewController(t)
	mockLinter1 := mock_api.NewMockLinter(mockctl)
	mockLinter2 := mock_api.NewMockLinter(mockctl)
	mockLinter1.EXPECT().Rules().Times(2).Return(rules1)
	mockLinter2.EXPECT().Rules().Times(1).Return(rules2)

	linters.ByCategory = map[api.Category]api.Linter{
		categories.VersionControl: mockLinter1,
		categories.CodeQuality:    mockLinter2,
	}

	t.Run("NoCategory", func(t *testing.T) {
		rule := linters.GetRule("basic-rulename")
		require.Nil(t, rule)
	})

	t.Run("UnimplementedCategory", func(t *testing.T) {
		rule := linters.GetRule("data-quality/some-rule")
		require.Nil(t, rule)
	})

	t.Run("CategoryKnownButNoRule", func(t *testing.T) {
		rule := linters.GetRule("version-control/some-rule")
		require.Nil(t, rule)
	})

	t.Run("ExistingRules", func(t *testing.T) {
		rule := linters.GetRule("version-control/test-rule-1")
		require.Equal(t, rules1[0], rule)

		rule = linters.GetRule("code-quality/test-rule-4")
		require.Equal(t, rules2[1], rule)
	})
}

func TestFindRules(t *testing.T) {
	rules1 := []*api.Rule{
		{Name: "Test Rule 1", Slug: "version-control/test-rule-1"},
		{Name: "Test Rule 2", Slug: "version-control/test-rule-2"},
	}
	rules2 := []*api.Rule{
		{Name: "Test Rule 3", Slug: "code-quality/test-rule-3"},
		{Name: "Test Rule 4", Slug: "code-quality/test-rule-4"},
		{Name: "Actual Rule 5", Slug: "code-quality/actual-rule-5"},
	}

	mockctl := gomock.NewController(t)
	mockLinter1 := mock_api.NewMockLinter(mockctl)
	mockLinter2 := mock_api.NewMockLinter(mockctl)
	mockLinter1.EXPECT().Rules().Times(2).Return(rules1)
	mockLinter2.EXPECT().Rules().Times(1).Return(rules2)

	linters.ByCategory = map[api.Category]api.Linter{
		categories.VersionControl: mockLinter1,
		categories.CodeQuality:    mockLinter2,
	}

	t.Run("NoCategory", func(t *testing.T) {
		rules := linters.FindRules("this-doesn't-work")
		require.Equal(t, []*api.Rule{}, rules)
	})

	t.Run("UnimplementedCategory", func(t *testing.T) {
		rules := linters.FindRules("data-quality")
		require.Equal(t, []*api.Rule{}, rules)
	})

	t.Run("EntireCategory", func(t *testing.T) {
		rules := linters.FindRules("version-control")
		require.Equal(t, rules1, rules)
	})

	t.Run("ExactMatch", func(t *testing.T) {
		rules := linters.FindRules("version-control/test-rule-1")
		require.Equal(t, []*api.Rule{rules1[0]}, rules)
	})

	t.Run("MatchMultiple", func(t *testing.T) {
		rules := linters.FindRules("code-quality/test")
		require.Equal(t, rules2[:2], rules)
	})
}
