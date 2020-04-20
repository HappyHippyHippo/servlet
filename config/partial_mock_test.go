// Code generated by MockGen. DO NOT EDIT.
// Source: config/partial.go

// Package config is a generated GoMock package.
package config

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPartial is a mock of Partial interface
type MockPartial struct {
	ctrl     *gomock.Controller
	recorder *MockPartialMockRecorder
}

// MockPartialMockRecorder is the mock recorder for MockPartial
type MockPartialMockRecorder struct {
	mock *MockPartial
}

// NewMockPartial creates a new mock instance
func NewMockPartial(ctrl *gomock.Controller) *MockPartial {
	mock := &MockPartial{ctrl: ctrl}
	mock.recorder = &MockPartialMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPartial) EXPECT() *MockPartialMockRecorder {
	return m.recorder
}

// Has mocks base method
func (m *MockPartial) Has(path string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Has", path)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Has indicates an expected call of Has
func (mr *MockPartialMockRecorder) Has(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Has", reflect.TypeOf((*MockPartial)(nil).Has), path)
}

// Get mocks base method
func (m *MockPartial) Get(path string, def ...interface{}) interface{} {
	m.ctrl.T.Helper()
	varargs := []interface{}{path}
	for _, a := range def {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockPartialMockRecorder) Get(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPartial)(nil).Get), varargs...)
}

// Int mocks base method
func (m *MockPartial) Int(path string, def ...int) int {
	m.ctrl.T.Helper()
	varargs := []interface{}{path}
	for _, a := range def {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Int", varargs...)
	ret0, _ := ret[0].(int)
	return ret0
}

// Int indicates an expected call of Int
func (mr *MockPartialMockRecorder) Int(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Int", reflect.TypeOf((*MockPartial)(nil).Int), varargs...)
}

// String mocks base method
func (m *MockPartial) String(path string, def ...string) string {
	m.ctrl.T.Helper()
	varargs := []interface{}{path}
	for _, a := range def {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "String", varargs...)
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockPartialMockRecorder) String(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockPartial)(nil).String), varargs...)
}

// Config mocks base method
func (m *MockPartial) Config(path string, def ...Partial) Partial {
	m.ctrl.T.Helper()
	varargs := []interface{}{path}
	for _, a := range def {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Config", varargs...)
	ret0, _ := ret[0].(Partial)
	return ret0
}

// Config indicates an expected call of Config
func (mr *MockPartialMockRecorder) Config(path interface{}, def ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{path}, def...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Config", reflect.TypeOf((*MockPartial)(nil).Config), varargs...)
}