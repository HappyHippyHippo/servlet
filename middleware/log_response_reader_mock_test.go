// Code generated by MockGen. DO NOT EDIT.
// Source: middleware/log_response_reader.go

// Package middleware is a generated GoMock package.
package middleware

import (
	gomock "github.com/golang/mock/gomock"
	servlet "github.com/happyhippyhippo/servlet"
	reflect "reflect"
)

// MockLogResponseReader is a mock of LogResponseReader interface
type MockLogResponseReader struct {
	ctrl     *gomock.Controller
	recorder *MockLogResponseReaderMockRecorder
}

// MockLogResponseReaderMockRecorder is the mock recorder for MockLogResponseReader
type MockLogResponseReaderMockRecorder struct {
	mock *MockLogResponseReader
}

// NewMockLogResponseReader creates a new mock instance
func NewMockLogResponseReader(ctrl *gomock.Controller) *MockLogResponseReader {
	mock := &MockLogResponseReader{ctrl: ctrl}
	mock.recorder = &MockLogResponseReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogResponseReader) EXPECT() *MockLogResponseReaderMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockLogResponseReader) Get(context servlet.Context) map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", context)
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockLogResponseReaderMockRecorder) Get(context interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLogResponseReader)(nil).Get), context)
}