// Code generated by MockGen. DO NOT EDIT.
// Source: log/stream_factory_strategy.go

// Package log is a generated GoMock package.
package log

import (
	gomock "github.com/golang/mock/gomock"
	config "github.com/happyhippyhippo/servlet/config"
	reflect "reflect"
)

// MockStreamFactoryStrategy is a mock of StreamFactoryStrategy interface
type MockStreamFactoryStrategy struct {
	ctrl     *gomock.Controller
	recorder *MockStreamFactoryStrategyMockRecorder
}

// MockStreamFactoryStrategyMockRecorder is the mock recorder for MockStreamFactoryStrategy
type MockStreamFactoryStrategyMockRecorder struct {
	mock *MockStreamFactoryStrategy
}

// NewMockStreamFactoryStrategy creates a new mock instance
func NewMockStreamFactoryStrategy(ctrl *gomock.Controller) *MockStreamFactoryStrategy {
	mock := &MockStreamFactoryStrategy{ctrl: ctrl}
	mock.recorder = &MockStreamFactoryStrategyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStreamFactoryStrategy) EXPECT() *MockStreamFactoryStrategyMockRecorder {
	return m.recorder
}

// Accept mocks base method
func (m *MockStreamFactoryStrategy) Accept(stype string, args ...interface{}) bool {
	m.ctrl.T.Helper()
	varargs := []interface{}{stype}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Accept", varargs...)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Accept indicates an expected call of Accept
func (mr *MockStreamFactoryStrategyMockRecorder) Accept(stype interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{stype}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accept", reflect.TypeOf((*MockStreamFactoryStrategy)(nil).Accept), varargs...)
}

// AcceptConfig mocks base method
func (m *MockStreamFactoryStrategy) AcceptConfig(conf config.Partial) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptConfig", conf)
	ret0, _ := ret[0].(bool)
	return ret0
}

// AcceptConfig indicates an expected call of AcceptConfig
func (mr *MockStreamFactoryStrategyMockRecorder) AcceptConfig(conf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptConfig", reflect.TypeOf((*MockStreamFactoryStrategy)(nil).AcceptConfig), conf)
}

// Create mocks base method
func (m *MockStreamFactoryStrategy) Create(args ...interface{}) (Stream, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockStreamFactoryStrategyMockRecorder) Create(args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStreamFactoryStrategy)(nil).Create), args...)
}

// CreateConfig mocks base method
func (m *MockStreamFactoryStrategy) CreateConfig(conf config.Partial) (Stream, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateConfig", conf)
	ret0, _ := ret[0].(Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateConfig indicates an expected call of CreateConfig
func (mr *MockStreamFactoryStrategyMockRecorder) CreateConfig(conf interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateConfig", reflect.TypeOf((*MockStreamFactoryStrategy)(nil).CreateConfig), conf)
}
