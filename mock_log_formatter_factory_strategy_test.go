// Code generated by MockGen. DO NOT EDIT.
// Source: log_formatter_factory_strategy.go

// Package servlet is a generated GoMock package.
package servlet

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLogFormatterFactoryStrategy is a mock of LogFormatterFactoryStrategy interface
type MockLogFormatterFactoryStrategy struct {
	ctrl     *gomock.Controller
	recorder *MockLogFormatterFactoryStrategyMockRecorder
}

// MockLogFormatterFactoryStrategyMockRecorder is the mock recorder for MockLogFormatterFactoryStrategy
type MockLogFormatterFactoryStrategyMockRecorder struct {
	mock *MockLogFormatterFactoryStrategy
}

// NewMockLogFormatterFactoryStrategy creates a new mock instance
func NewMockLogFormatterFactoryStrategy(ctrl *gomock.Controller) *MockLogFormatterFactoryStrategy {
	mock := &MockLogFormatterFactoryStrategy{ctrl: ctrl}
	mock.recorder = &MockLogFormatterFactoryStrategyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogFormatterFactoryStrategy) EXPECT() *MockLogFormatterFactoryStrategyMockRecorder {
	return m.recorder
}

// Accept mocks base method
func (m *MockLogFormatterFactoryStrategy) Accept(format string, args ...interface{}) bool {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Accept", varargs...)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept
func (mr *MockLogFormatterFactoryStrategyMockRecorder) Accept(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockLogFormatterFactoryStrategy)(nil).Accept), varargs...)
}

// Create mocks base method
func (m *MockLogFormatterFactoryStrategy) Create(args ...interface{}) (LogFormatter, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(LogFormatter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockLogFormatterFactoryStrategyMockRecorder) Create(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLogFormatterFactoryStrategy)(nil).Create), args...)
}