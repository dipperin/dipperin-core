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
// Source: github.com/caiqingfeng/dipperin-core/core/chain-communication (interfaces: PeerManager)

package chain_communication

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of PeerManager interface
type MockPeerManager struct {
	ctrl     *gomock.Controller
	recorder *_MockPeerManagerRecorder
}

// Recorder for MockPeerManager (not exported)
type _MockPeerManagerRecorder struct {
	mock *MockPeerManager
}

func NewMockPeerManager(ctrl *gomock.Controller) *MockPeerManager {
	mock := &MockPeerManager{ctrl: ctrl}
	mock.recorder = &_MockPeerManagerRecorder{mock}
	return mock
}

func (_m *MockPeerManager) EXPECT() *_MockPeerManagerRecorder {
	return _m.recorder
}

func (_m *MockPeerManager) BestPeer() PmAbstractPeer {
	ret := _m.ctrl.Call(_m, "BestPeer")
	ret0, _ := ret[0].(PmAbstractPeer)
	return ret0
}

func (_mr *_MockPeerManagerRecorder) BestPeer() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "BestPeer")
}

func (_m *MockPeerManager) GetPeer(_param0 string) PmAbstractPeer {
	ret := _m.ctrl.Call(_m, "GetPeer", _param0)
	ret0, _ := ret[0].(PmAbstractPeer)
	return ret0
}

func (_mr *_MockPeerManagerRecorder) GetPeer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetPeer", arg0)
}

func (_m *MockPeerManager) GetPeers() map[string]PmAbstractPeer {
	ret := _m.ctrl.Call(_m, "GetPeers")
	ret0, _ := ret[0].(map[string]PmAbstractPeer)
	return ret0
}

func (_mr *_MockPeerManagerRecorder) GetPeers() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetPeers")
}

func (_m *MockPeerManager) IsSync() bool {
	ret := _m.ctrl.Call(_m, "IsSync")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockPeerManagerRecorder) IsSync() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsSync")
}

func (_m *MockPeerManager) RemovePeer(_param0 string) {
	_m.ctrl.Call(_m, "RemovePeer", _param0)
}

func (_mr *_MockPeerManagerRecorder) RemovePeer(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemovePeer", arg0)
}
