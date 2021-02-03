// Code generated by MockGen. DO NOT EDIT.
// Source: provider.go

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	auth "github.com/aintsashqa/go-simple-blog/pkg/auth"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	reflect "reflect"
)

// MockAuthorizationProvider is a mock of AuthorizationProvider interface
type MockAuthorizationProvider struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationProviderMockRecorder
}

// MockAuthorizationProviderMockRecorder is the mock recorder for MockAuthorizationProvider
type MockAuthorizationProviderMockRecorder struct {
	mock *MockAuthorizationProvider
}

// NewMockAuthorizationProvider creates a new mock instance
func NewMockAuthorizationProvider(ctrl *gomock.Controller) *MockAuthorizationProvider {
	mock := &MockAuthorizationProvider{ctrl: ctrl}
	mock.recorder = &MockAuthorizationProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthorizationProvider) EXPECT() *MockAuthorizationProviderMockRecorder {
	return m.recorder
}

// NewToken mocks base method
func (m *MockAuthorizationProvider) NewToken(arg0 auth.TokenParams) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewToken", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewToken indicates an expected call of NewToken
func (mr *MockAuthorizationProviderMockRecorder) NewToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewToken", reflect.TypeOf((*MockAuthorizationProvider)(nil).NewToken), arg0)
}

// Parse mocks base method
func (m *MockAuthorizationProvider) Parse(arg0 string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parse", arg0)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Parse indicates an expected call of Parse
func (mr *MockAuthorizationProviderMockRecorder) Parse(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parse", reflect.TypeOf((*MockAuthorizationProvider)(nil).Parse), arg0)
}
