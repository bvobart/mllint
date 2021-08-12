// Code generated by MockGen. DO NOT EDIT.
// Source: api/depmanagers.go

// Package mock_api is a generated GoMock package.
package mock_api

import (
	api "github.com/bvobart/mllint/api"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDependencyManagerType is a mock of DependencyManagerType interface
type MockDependencyManagerType struct {
	ctrl     *gomock.Controller
	recorder *MockDependencyManagerTypeMockRecorder
}

// MockDependencyManagerTypeMockRecorder is the mock recorder for MockDependencyManagerType
type MockDependencyManagerTypeMockRecorder struct {
	mock *MockDependencyManagerType
}

// NewMockDependencyManagerType creates a new mock instance
func NewMockDependencyManagerType(ctrl *gomock.Controller) *MockDependencyManagerType {
	mock := &MockDependencyManagerType{ctrl: ctrl}
	mock.recorder = &MockDependencyManagerTypeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDependencyManagerType) EXPECT() *MockDependencyManagerTypeMockRecorder {
	return m.recorder
}

// String mocks base method
func (m *MockDependencyManagerType) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockDependencyManagerTypeMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockDependencyManagerType)(nil).String))
}

// Detect mocks base method
func (m *MockDependencyManagerType) Detect(project api.Project) (api.DependencyManager, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Detect", project)
	ret0, _ := ret[0].(api.DependencyManager)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Detect indicates an expected call of Detect
func (mr *MockDependencyManagerTypeMockRecorder) Detect(project interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Detect", reflect.TypeOf((*MockDependencyManagerType)(nil).Detect), project)
}

// MockDependencyManager is a mock of DependencyManager interface
type MockDependencyManager struct {
	ctrl     *gomock.Controller
	recorder *MockDependencyManagerMockRecorder
}

// MockDependencyManagerMockRecorder is the mock recorder for MockDependencyManager
type MockDependencyManagerMockRecorder struct {
	mock *MockDependencyManager
}

// NewMockDependencyManager creates a new mock instance
func NewMockDependencyManager(ctrl *gomock.Controller) *MockDependencyManager {
	mock := &MockDependencyManager{ctrl: ctrl}
	mock.recorder = &MockDependencyManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDependencyManager) EXPECT() *MockDependencyManagerMockRecorder {
	return m.recorder
}

// Dependencies mocks base method
func (m *MockDependencyManager) Dependencies() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dependencies")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Dependencies indicates an expected call of Dependencies
func (mr *MockDependencyManagerMockRecorder) Dependencies() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dependencies", reflect.TypeOf((*MockDependencyManager)(nil).Dependencies))
}

// HasDependency mocks base method
func (m *MockDependencyManager) HasDependency(dependency string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasDependency", dependency)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasDependency indicates an expected call of HasDependency
func (mr *MockDependencyManagerMockRecorder) HasDependency(dependency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasDependency", reflect.TypeOf((*MockDependencyManager)(nil).HasDependency), dependency)
}

// HasDevDependency mocks base method
func (m *MockDependencyManager) HasDevDependency(dependency string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasDevDependency", dependency)
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasDevDependency indicates an expected call of HasDevDependency
func (mr *MockDependencyManagerMockRecorder) HasDevDependency(dependency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasDevDependency", reflect.TypeOf((*MockDependencyManager)(nil).HasDevDependency), dependency)
}

// Type mocks base method
func (m *MockDependencyManager) Type() api.DependencyManagerType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(api.DependencyManagerType)
	return ret0
}

// Type indicates an expected call of Type
func (mr *MockDependencyManagerMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockDependencyManager)(nil).Type))
}