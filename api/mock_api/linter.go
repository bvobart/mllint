// Code generated by MockGen. DO NOT EDIT.
// Source: api/linter.go

// Package mock_api is a generated GoMock package.
package mock_api

import (
	api "github.com/bvobart/mllint/api"
	config "github.com/bvobart/mllint/config"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLinter is a mock of Linter interface
type MockLinter struct {
	ctrl     *gomock.Controller
	recorder *MockLinterMockRecorder
}

// MockLinterMockRecorder is the mock recorder for MockLinter
type MockLinterMockRecorder struct {
	mock *MockLinter
}

// NewMockLinter creates a new mock instance
func NewMockLinter(ctrl *gomock.Controller) *MockLinter {
	mock := &MockLinter{ctrl: ctrl}
	mock.recorder = &MockLinterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLinter) EXPECT() *MockLinterMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockLinter) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockLinterMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockLinter)(nil).Name))
}

// Rules mocks base method
func (m *MockLinter) Rules() []*api.Rule {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rules")
	ret0, _ := ret[0].([]*api.Rule)
	return ret0
}

// Rules indicates an expected call of Rules
func (mr *MockLinterMockRecorder) Rules() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rules", reflect.TypeOf((*MockLinter)(nil).Rules))
}

// LintProject mocks base method
func (m *MockLinter) LintProject(project api.Project) (api.Report, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LintProject", project)
	ret0, _ := ret[0].(api.Report)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LintProject indicates an expected call of LintProject
func (mr *MockLinterMockRecorder) LintProject(project interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LintProject", reflect.TypeOf((*MockLinter)(nil).LintProject), project)
}

// MockConfigurable is a mock of Configurable interface
type MockConfigurable struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurableMockRecorder
}

// MockConfigurableMockRecorder is the mock recorder for MockConfigurable
type MockConfigurableMockRecorder struct {
	mock *MockConfigurable
}

// NewMockConfigurable creates a new mock instance
func NewMockConfigurable(ctrl *gomock.Controller) *MockConfigurable {
	mock := &MockConfigurable{ctrl: ctrl}
	mock.recorder = &MockConfigurableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigurable) EXPECT() *MockConfigurableMockRecorder {
	return m.recorder
}

// Configure mocks base method
func (m *MockConfigurable) Configure(conf *config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configure", conf)
	ret0, _ := ret[0].(error)
	return ret0
}

// Configure indicates an expected call of Configure
func (mr *MockConfigurableMockRecorder) Configure(conf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockConfigurable)(nil).Configure), conf)
}

// MockConfigurableLinter is a mock of ConfigurableLinter interface
type MockConfigurableLinter struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurableLinterMockRecorder
}

// MockConfigurableLinterMockRecorder is the mock recorder for MockConfigurableLinter
type MockConfigurableLinterMockRecorder struct {
	mock *MockConfigurableLinter
}

// NewMockConfigurableLinter creates a new mock instance
func NewMockConfigurableLinter(ctrl *gomock.Controller) *MockConfigurableLinter {
	mock := &MockConfigurableLinter{ctrl: ctrl}
	mock.recorder = &MockConfigurableLinterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConfigurableLinter) EXPECT() *MockConfigurableLinterMockRecorder {
	return m.recorder
}

// Name mocks base method
func (m *MockConfigurableLinter) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockConfigurableLinterMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockConfigurableLinter)(nil).Name))
}

// Rules mocks base method
func (m *MockConfigurableLinter) Rules() []*api.Rule {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rules")
	ret0, _ := ret[0].([]*api.Rule)
	return ret0
}

// Rules indicates an expected call of Rules
func (mr *MockConfigurableLinterMockRecorder) Rules() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rules", reflect.TypeOf((*MockConfigurableLinter)(nil).Rules))
}

// LintProject mocks base method
func (m *MockConfigurableLinter) LintProject(project api.Project) (api.Report, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LintProject", project)
	ret0, _ := ret[0].(api.Report)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LintProject indicates an expected call of LintProject
func (mr *MockConfigurableLinterMockRecorder) LintProject(project interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LintProject", reflect.TypeOf((*MockConfigurableLinter)(nil).LintProject), project)
}

// Configure mocks base method
func (m *MockConfigurableLinter) Configure(conf *config.Config) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configure", conf)
	ret0, _ := ret[0].(error)
	return ret0
}

// Configure indicates an expected call of Configure
func (mr *MockConfigurableLinterMockRecorder) Configure(conf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockConfigurableLinter)(nil).Configure), conf)
}
