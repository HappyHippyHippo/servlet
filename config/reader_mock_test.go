// Code generated by MockGen. DO NOT EDIT.
package config

import (
	"bytes"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockReader is a mock of io.Reader interface
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderRecorder
	reader   *bytes.Buffer
}

// MockReaderRecorder is the mock recorder for MockReader
type MockReaderRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance
func NewMockReader(ctrl *gomock.Controller, content string) *MockReader {
	mock := &MockReader{ctrl: ctrl, reader: bytes.NewBufferString(content)}
	mock.recorder = &MockReaderRecorder{mock}
	return mock
}

// Content returns the content of the mocked reader
func (m *MockReader) Content() string {
	return m.reader.String()
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReader) EXPECT() *MockReaderRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockReader) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockReaderRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReader)(nil).Close))
}

// Read mocks base method
func (m *MockReader) Read(b []byte) (n int, err error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", b)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)

	m.reader.Read(b)

	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockReaderRecorder) Read(b interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReader)(nil).Read), b)
}