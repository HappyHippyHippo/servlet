// Code generated by MockGen. DO NOT EDIT.
// Source: container.go

// Package servlet is a generated GoMock package.
package servlet

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockContainer is a mock of Container interface
type MockContainer struct {
	ctrl     *gomock.Controller
	recorder *MockContainerMockRecorder
}

// MockContainerMockRecorder is the mock recorder for MockContainer
type MockContainerMockRecorder struct {
	mock *MockContainer
}

// NewMockContainer creates a new mock instance
func NewMockContainer(ctrl *gomock.Controller) *MockContainer {
	mock := &MockContainer{ctrl: ctrl}
	mock.recorder = &MockContainerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContainer) EXPECT() *MockContainerMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockContainer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockContainerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockContainer)(nil).Close))
}

// Has mocks base method
func (m *MockContainer) Has(id string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", id)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Has indicates an expected call of Has
func (mr *MockContainerMockRecorder) Has(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockContainer)(nil).Has), id)
}

// Add mocks base method
func (m *MockContainer) Add(id string, factory Factory) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", id, factory)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockContainerMockRecorder) Add(id, factory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockContainer)(nil).Add), id, factory)
}

// Remove mocks base method
func (m *MockContainer) Remove(id string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Remove", id)
}

// Remove indicates an expected call of Remove
func (mr *MockContainerMockRecorder) Remove(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockContainer)(nil).Remove), id)
}

// Get mocks base method
func (m *MockContainer) Get(id string) interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", id)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockContainerMockRecorder) Get(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockContainer)(nil).Get), id)
}
