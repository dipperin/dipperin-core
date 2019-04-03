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

// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/caiqingfeng/dipperin-core/core/chain-communication (interfaces: P2PServer)

package chain_communication

import (
	enode "github.com/dipperin/dipperin-core/third-party/p2p/enode"
	gomock "github.com/golang/mock/gomock"
)

// Mock of P2PServer interface
type MockP2PServer struct {
	ctrl     *gomock.Controller
	recorder *_MockP2PServerRecorder
}

// Recorder for MockP2PServer (not exported)
type _MockP2PServerRecorder struct {
	mock *MockP2PServer
}

func NewMockP2PServer(ctrl *gomock.Controller) *MockP2PServer {
	mock := &MockP2PServer{ctrl: ctrl}
	mock.recorder = &_MockP2PServerRecorder{mock}
	return mock
}

func (_m *MockP2PServer) EXPECT() *_MockP2PServerRecorder {
	return _m.recorder
}

func (_m *MockP2PServer) AddPeer(_param0 *enode.Node) {
	_m.ctrl.Call(_m, "AddPeer", _param0)
}

func (_mr *_MockP2PServerRecorder) AddPeer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddPeer", arg0)
}

func (_m *MockP2PServer) RemovePeer(_param0 *enode.Node) {
	_m.ctrl.Call(_m, "RemovePeer", _param0)
}

func (_mr *_MockP2PServerRecorder) RemovePeer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemovePeer", arg0)
}

func (_m *MockP2PServer) Self() *enode.Node {
	ret := _m.ctrl.Call(_m, "Self")
	ret0, _ := ret[0].(*enode.Node)
	return ret0
}

func (_mr *_MockP2PServerRecorder) Self() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Self")
}
