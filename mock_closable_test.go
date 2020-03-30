// Code generated by MockGen. DO NOT EDIT.
// Source: closable.go

// Package servlet is a generated GoMock package.
package servlet

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClosable is a mock of Closable interface
type MockClosable struct {
	ctrl     *gomock.Controller
	recorder *MockClosableMockRecorder
}

// MockClosableMockRecorder is the mock recorder for MockClosable
type MockClosableMockRecorder struct {
	mock *MockClosable
}

// NewMockClosable creates a new mock instance
func NewMockClosable(ctrl *gomock.Controller) *MockClosable {
	mock := &MockClosable{ctrl: ctrl}
	mock.recorder = &MockClosableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClosable) EXPECT() *MockClosableMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockClosable) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockClosableMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClosable)(nil).Close))
}
