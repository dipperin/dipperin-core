// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/chain-communication (interfaces: P2PPeer)

// Package chain_communication is a generated GoMock package.
package chain_communication

import (
	p2p "github.com/dipperin/dipperin-core/third-party/p2p"
	enode "github.com/dipperin/dipperin-core/third-party/p2p/enode"
	gomock "github.com/golang/mock/gomock"
	net "net"
	reflect "reflect"
)

// MockP2PPeer is a mock of P2PPeer interface
type MockP2PPeer struct {
	ctrl     *gomock.Controller
	recorder *MockP2PPeerMockRecorder
}

// MockP2PPeerMockRecorder is the mock recorder for MockP2PPeer
type MockP2PPeerMockRecorder struct {
	mock *MockP2PPeer
}

// NewMockP2PPeer creates a new mock instance
func NewMockP2PPeer(ctrl *gomock.Controller) *MockP2PPeer {
	mock := &MockP2PPeer{ctrl: ctrl}
	mock.recorder = &MockP2PPeerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockP2PPeer) EXPECT() *MockP2PPeerMockRecorder {
	return m.recorder
}

// Disconnect mocks base method
func (m *MockP2PPeer) Disconnect(arg0 p2p.DiscReason) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Disconnect", arg0)
}

// Disconnect indicates an expected call of Disconnect
func (mr *MockP2PPeerMockRecorder) Disconnect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockP2PPeer)(nil).Disconnect), arg0)
}

// ID mocks base method
func (m *MockP2PPeer) ID() enode.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(enode.ID)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockP2PPeerMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockP2PPeer)(nil).ID))
}

// RemoteAddr mocks base method
func (m *MockP2PPeer) RemoteAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr
func (mr *MockP2PPeerMockRecorder) RemoteAddr() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockP2PPeer)(nil).RemoteAddr))
}
