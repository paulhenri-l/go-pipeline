// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/paulhenri-l/go-pipeline/contracts (interfaces: Sink)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSink is a mock of Sink interface.
type MockSink struct {
	ctrl     *gomock.Controller
	recorder *MockSinkMockRecorder
}

// MockSinkMockRecorder is the mock recorder for MockSink.
type MockSinkMockRecorder struct {
	mock *MockSink
}

// NewMockSink creates a new mock instance.
func NewMockSink(ctrl *gomock.Controller) *MockSink {
	mock := &MockSink{ctrl: ctrl}
	mock.recorder = &MockSinkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSink) EXPECT() *MockSinkMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockSink) Start(arg0 context.Context, arg1 <-chan interface{}) <-chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0, arg1)
	ret0, _ := ret[0].(<-chan struct{})
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockSinkMockRecorder) Start(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockSink)(nil).Start), arg0, arg1)
}
