// Code generated by MockGen. DO NOT EDIT.
// Source: log_middleware.go

// Package servlet is a generated GoMock package.
package servlet

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockLogMiddlewareRequestReader is a mock of LogMiddlewareRequestReader interface
type MockLogMiddlewareRequestReader struct {
	ctrl     *gomock.Controller
	recorder *MockLogMiddlewareRequestReaderMockRecorder
}

// MockLogMiddlewareRequestReaderMockRecorder is the mock recorder for MockLogMiddlewareRequestReader
type MockLogMiddlewareRequestReaderMockRecorder struct {
	mock *MockLogMiddlewareRequestReader
}

// NewMockLogMiddlewareRequestReader creates a new mock instance
func NewMockLogMiddlewareRequestReader(ctrl *gomock.Controller) *MockLogMiddlewareRequestReader {
	mock := &MockLogMiddlewareRequestReader{ctrl: ctrl}
	mock.recorder = &MockLogMiddlewareRequestReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogMiddlewareRequestReader) EXPECT() *MockLogMiddlewareRequestReaderMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockLogMiddlewareRequestReader) Get(context GinContext) map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", context)
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockLogMiddlewareRequestReaderMockRecorder) Get(context interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLogMiddlewareRequestReader)(nil).Get), context)
}

// MockLogMiddlewareResponseReader is a mock of LogMiddlewareResponseReader interface
type MockLogMiddlewareResponseReader struct {
	ctrl     *gomock.Controller
	recorder *MockLogMiddlewareResponseReaderMockRecorder
}

// MockLogMiddlewareResponseReaderMockRecorder is the mock recorder for MockLogMiddlewareResponseReader
type MockLogMiddlewareResponseReaderMockRecorder struct {
	mock *MockLogMiddlewareResponseReader
}

// NewMockLogMiddlewareResponseReader creates a new mock instance
func NewMockLogMiddlewareResponseReader(ctrl *gomock.Controller) *MockLogMiddlewareResponseReader {
	mock := &MockLogMiddlewareResponseReader{ctrl: ctrl}
	mock.recorder = &MockLogMiddlewareResponseReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogMiddlewareResponseReader) EXPECT() *MockLogMiddlewareResponseReaderMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockLogMiddlewareResponseReader) Get(context GinContext) map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", context)
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockLogMiddlewareResponseReaderMockRecorder) Get(context interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLogMiddlewareResponseReader)(nil).Get), context)
}
