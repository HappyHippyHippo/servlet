// Code generated by MockGen. DO NOT EDIT.
package log

import (
	"bytes"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockWriter is a mock of io.Writer interface
type MockWriter struct {
	ctrl     *gomock.Controller
	recorder *MockWriterRecorder
	writer   *bytes.Buffer
}

// MockWriterRecorder is the mock recorder for MockWriter
type MockWriterRecorder struct {
	mock *MockWriter
}

// NewMockWriter creates a new mock instance
func NewMockWriter(ctrl *gomock.Controller) *MockWriter {
	mock := &MockWriter{ctrl: ctrl, writer: bytes.NewBufferString("")}
	mock.recorder = &MockWriterRecorder{mock}
	return mock
}

// Content returns the content writen into the mocked writer
func (m *MockWriter) Content() string {
	return m.writer.String()
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWriter) EXPECT() *MockWriterRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockWriter) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockWriterRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWriter)(nil).Close))
}

// Write mocks base method
func (m *MockWriter) Write(p []byte) (n int, err error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)

	m.writer.Write(p)

	return ret0, ret1
}

// Write indicates an expected call of Write
func (mr *MockWriterRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockWriter)(nil).Write), p)
}
