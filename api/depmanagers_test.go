package api_test

import (
	"testing"

	"github.com/bvobart/mllint/api"
	"github.com/bvobart/mllint/api/mock_api"
	"github.com/bvobart/mllint/setools/depmanagers"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// Main returns the first dependency manager in the list, under the assumption that that is the main / primary dependency manager used in the project.
func TestMain(t *testing.T) {
	managers := api.DependencyManagerList{}
	require.Nil(t, managers.Main())

	ctrl := gomock.NewController(t)
	mock := mock_api.NewMockDependencyManager(ctrl)
	mock2 := mock_api.NewMockDependencyManager(ctrl)
	managers = append(managers, mock, mock2)

	require.Equal(t, mock, managers.Main())
}

func TestContains(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock_api.NewMockDependencyManager(ctrl)
	mock2 := mock_api.NewMockDependencyManager(ctrl)
	mock3 := mock_api.NewMockDependencyManager(ctrl)
	managers := api.DependencyManagerList{mock, mock2}

	require.False(t, managers.Contains(nil))
	require.True(t, managers.Contains(mock))
	require.True(t, managers.Contains(mock2))
	require.False(t, managers.Contains(mock3))
}

func TestContainsType(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock_api.NewMockDependencyManager(ctrl)
	mock2 := mock_api.NewMockDependencyManager(ctrl)
	mock.EXPECT().Type().Times(3).Return(depmanagers.TypePoetry)
	mock2.EXPECT().Type().Times(2).Return(depmanagers.TypeRequirementsTxt)
	managers := api.DependencyManagerList{mock, mock2}

	require.False(t, managers.ContainsType(nil))
	require.True(t, managers.ContainsType(depmanagers.TypePoetry))
	require.True(t, managers.ContainsType(depmanagers.TypeRequirementsTxt))
	require.False(t, managers.ContainsType(depmanagers.TypeSetupPy))
}

func TestContainsAllTypes(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock_api.NewMockDependencyManager(ctrl)
	mock2 := mock_api.NewMockDependencyManager(ctrl)
	mock.EXPECT().Type().Times(5).Return(depmanagers.TypePoetry)
	mock2.EXPECT().Type().Times(3).Return(depmanagers.TypeRequirementsTxt)
	managers := api.DependencyManagerList{mock, mock2}

	require.True(t, managers.ContainsAllTypes(depmanagers.TypePoetry, depmanagers.TypeRequirementsTxt))
	require.False(t, managers.ContainsAllTypes(depmanagers.TypePoetry, depmanagers.TypeRequirementsTxt, depmanagers.TypeSetupPy))
}
