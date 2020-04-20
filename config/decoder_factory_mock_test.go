// Code generated by MockGen. DO NOT EDIT.
// Source: config/decoder_factory.go

// Package config is a generated GoMock package.
package config

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDecoderFactory is a mock of DecoderFactory interface
type MockDecoderFactory struct {
	ctrl     *gomock.Controller
	recorder *MockDecoderFactoryMockRecorder
}

// MockDecoderFactoryMockRecorder is the mock recorder for MockDecoderFactory
type MockDecoderFactoryMockRecorder struct {
	mock *MockDecoderFactory
}

// NewMockDecoderFactory creates a new mock instance
func NewMockDecoderFactory(ctrl *gomock.Controller) *MockDecoderFactory {
	mock := &MockDecoderFactory{ctrl: ctrl}
	mock.recorder = &MockDecoderFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDecoderFactory) EXPECT() *MockDecoderFactoryMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockDecoderFactory) Register(strategy DecoderFactoryStrategy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", strategy)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockDecoderFactoryMockRecorder) Register(strategy interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockDecoderFactory)(nil).Register), strategy)
}

// Create mocks base method
func (m *MockDecoderFactory) Create(format string, args ...interface{}) (Decoder, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{format}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(Decoder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockDecoderFactoryMockRecorder) Create(format interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{format}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDecoderFactory)(nil).Create), varargs...)
}