// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/third_party/p2p (interfaces: MsgReadWriter)

// Package chaincommunication is a generated GoMock package.
package chaincommunication

import (
	p2p "github.com/dipperin/dipperin-core/third_party/p2p"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMsgReadWriter is a mock of MsgReadWriter interface
type MockMsgReadWriter struct {
	ctrl     *gomock.Controller
	recorder *MockMsgReadWriterMockRecorder
}

// MockMsgReadWriterMockRecorder is the mock recorder for MockMsgReadWriter
type MockMsgReadWriterMockRecorder struct {
	mock *MockMsgReadWriter
}

// NewMockMsgReadWriter creates a new mock instance
func NewMockMsgReadWriter(ctrl *gomock.Controller) *MockMsgReadWriter {
	mock := &MockMsgReadWriter{ctrl: ctrl}
	mock.recorder = &MockMsgReadWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMsgReadWriter) EXPECT() *MockMsgReadWriterMockRecorder {
	return m.recorder
}

// ReadMsg mocks base method
func (m *MockMsgReadWriter) ReadMsg() (p2p.Msg, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMsg")
	ret0, _ := ret[0].(p2p.Msg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMsg indicates an expected call of ReadMsg
func (mr *MockMsgReadWriterMockRecorder) ReadMsg() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMsg", reflect.TypeOf((*MockMsgReadWriter)(nil).ReadMsg))
}

// WriteMsg mocks base method
func (m *MockMsgReadWriter) WriteMsg(arg0 p2p.Msg) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMsg indicates an expected call of WriteMsg
func (mr *MockMsgReadWriterMockRecorder) WriteMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMsg", reflect.TypeOf((*MockMsgReadWriter)(nil).WriteMsg), arg0)
}
