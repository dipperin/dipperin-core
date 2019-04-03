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
// Source: github.com/caiqingfeng/dipperin-core/core/chain-communication (interfaces: PmAbstractPeer)

package chain_communication

import (
	common "github.com/dipperin/dipperin-core/common"
	p2p "github.com/dipperin/dipperin-core/third-party/p2p"
	gomock "github.com/golang/mock/gomock"
	net "net"
)

// Mock of PmAbstractPeer interface
type MockPmAbstractPeer struct {
	ctrl     *gomock.Controller
	recorder *_MockPmAbstractPeerRecorder
}

// Recorder for MockPmAbstractPeer (not exported)
type _MockPmAbstractPeerRecorder struct {
	mock *MockPmAbstractPeer
}

func NewMockPmAbstractPeer(ctrl *gomock.Controller) *MockPmAbstractPeer {
	mock := &MockPmAbstractPeer{ctrl: ctrl}
	mock.recorder = &_MockPmAbstractPeerRecorder{mock}
	return mock
}

func (_m *MockPmAbstractPeer) EXPECT() *_MockPmAbstractPeerRecorder {
	return _m.recorder
}

func (_m *MockPmAbstractPeer) DisconnectPeer() {
	_m.ctrl.Call(_m, "DisconnectPeer")
}

func (_mr *_MockPmAbstractPeerRecorder) DisconnectPeer() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DisconnectPeer")
}

func (_m *MockPmAbstractPeer) GetCsPeerInfo() *p2p.CsPeerInfo {
	ret := _m.ctrl.Call(_m, "GetCsPeerInfo")
	ret0, _ := ret[0].(*p2p.CsPeerInfo)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) GetCsPeerInfo() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetCsPeerInfo")
}

func (_m *MockPmAbstractPeer) GetHead() (common.Hash, uint64) {
	ret := _m.ctrl.Call(_m, "GetHead")
	ret0, _ := ret[0].(common.Hash)
	ret1, _ := ret[1].(uint64)
	return ret0, ret1
}

func (_mr *_MockPmAbstractPeerRecorder) GetHead() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetHead")
}

func (_m *MockPmAbstractPeer) GetPeerRawUrl() string {
	ret := _m.ctrl.Call(_m, "GetPeerRawUrl")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) GetPeerRawUrl() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetPeerRawUrl")
}

func (_m *MockPmAbstractPeer) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

func (_m *MockPmAbstractPeer) IsRunning() bool {
	ret := _m.ctrl.Call(_m, "IsRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) IsRunning() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsRunning")
}

func (_m *MockPmAbstractPeer) NodeName() string {
	ret := _m.ctrl.Call(_m, "NodeName")
	ret0, _ := ret[0].(string)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) NodeName() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NodeName")
}

func (_m *MockPmAbstractPeer) NodeType() uint64 {
	ret := _m.ctrl.Call(_m, "NodeType")
	ret0, _ := ret[0].(uint64)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) NodeType() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "NodeType")
}

func (_m *MockPmAbstractPeer) ReadMsg() (p2p.Msg, error) {
	ret := _m.ctrl.Call(_m, "ReadMsg")
	ret0, _ := ret[0].(p2p.Msg)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockPmAbstractPeerRecorder) ReadMsg() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ReadMsg")
}

func (_m *MockPmAbstractPeer) RemoteAddress() net.Addr {
	ret := _m.ctrl.Call(_m, "RemoteAddress")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) RemoteAddress() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteAddress")
}

func (_m *MockPmAbstractPeer) RemoteVerifierAddress() common.Address {
	ret := _m.ctrl.Call(_m, "RemoteVerifierAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) RemoteVerifierAddress() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteVerifierAddress")
}

func (_m *MockPmAbstractPeer) SendMsg(_param0 uint64, _param1 interface{}) error {
	ret := _m.ctrl.Call(_m, "SendMsg", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPmAbstractPeerRecorder) SendMsg(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SendMsg", arg0, arg1)
}

func (_m *MockPmAbstractPeer) SetHead(_param0 common.Hash, _param1 uint64) {
	_m.ctrl.Call(_m, "SetHead", _param0, _param1)
}

func (_mr *_MockPmAbstractPeerRecorder) SetHead(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetHead", arg0, arg1)
}

func (_m *MockPmAbstractPeer) SetNodeName(_param0 string) {
	_m.ctrl.Call(_m, "SetNodeName", _param0)
}

func (_mr *_MockPmAbstractPeerRecorder) SetNodeName(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetNodeName", arg0)
}

func (_m *MockPmAbstractPeer) SetNodeType(_param0 uint64) {
	_m.ctrl.Call(_m, "SetNodeType", _param0)
}

func (_mr *_MockPmAbstractPeerRecorder) SetNodeType(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetNodeType", arg0)
}

func (_m *MockPmAbstractPeer) SetNotRunning() {
	_m.ctrl.Call(_m, "SetNotRunning")
}

func (_mr *_MockPmAbstractPeerRecorder) SetNotRunning() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetNotRunning")
}

func (_m *MockPmAbstractPeer) SetPeerRawUrl(_param0 string) {
	_m.ctrl.Call(_m, "SetPeerRawUrl", _param0)
}

func (_mr *_MockPmAbstractPeerRecorder) SetPeerRawUrl(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetPeerRawUrl", arg0)
}

func (_m *MockPmAbstractPeer) SetRemoteVerifierAddress(_param0 common.Address) {
	_m.ctrl.Call(_m, "SetRemoteVerifierAddress", _param0)
}

func (_mr *_MockPmAbstractPeerRecorder) SetRemoteVerifierAddress(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetRemoteVerifierAddress", arg0)
}
