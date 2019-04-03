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
// Source: github.com/caiqingfeng/dipperin-core/core/chain-communication (interfaces: AbstractPbftProtocolManager)

// Package rpc_interface is a generated GoMock package.
package rpc_interface

import (
	common "github.com/dipperin/dipperin-core/common"
	chain_communication "github.com/dipperin/dipperin-core/core/chain-communication"
	enode "github.com/dipperin/dipperin-core/third-party/p2p/enode"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAbstractPbftProtocolManager is a mock of AbstractPbftProtocolManager interface
type MockAbstractPbftProtocolManager struct {
	ctrl     *gomock.Controller
	recorder *MockAbstractPbftProtocolManagerMockRecorder
}

// MockAbstractPbftProtocolManagerMockRecorder is the mock recorder for MockAbstractPbftProtocolManager
type MockAbstractPbftProtocolManagerMockRecorder struct {
	mock *MockAbstractPbftProtocolManager
}

// NewMockAbstractPbftProtocolManager creates a new mock instance
func NewMockAbstractPbftProtocolManager(ctrl *gomock.Controller) *MockAbstractPbftProtocolManager {
	mock := &MockAbstractPbftProtocolManager{ctrl: ctrl}
	mock.recorder = &MockAbstractPbftProtocolManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAbstractPbftProtocolManager) EXPECT() *MockAbstractPbftProtocolManagerMockRecorder {
	return m.recorder
}

// BestPeer mocks base method
func (m *MockAbstractPbftProtocolManager) BestPeer() chain_communication.PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BestPeer")
	ret0, _ := ret[0].(chain_communication.PmAbstractPeer)
	return ret0
}

// BestPeer indicates an expected call of BestPeer
func (mr *MockAbstractPbftProtocolManagerMockRecorder) BestPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestPeer", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).BestPeer))
}

// GetCurrentConnectPeers mocks base method
func (m *MockAbstractPbftProtocolManager) GetCurrentConnectPeers() map[string]common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentConnectPeers")
	ret0, _ := ret[0].(map[string]common.Address)
	return ret0
}

// GetCurrentConnectPeers indicates an expected call of GetCurrentConnectPeers
func (mr *MockAbstractPbftProtocolManagerMockRecorder) GetCurrentConnectPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentConnectPeers", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).GetCurrentConnectPeers))
}

// GetNextVerifierPeers mocks base method
func (m *MockAbstractPbftProtocolManager) GetNextVerifierPeers() map[string]chain_communication.PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextVerifierPeers")
	ret0, _ := ret[0].(map[string]chain_communication.PmAbstractPeer)
	return ret0
}

// GetNextVerifierPeers indicates an expected call of GetNextVerifierPeers
func (mr *MockAbstractPbftProtocolManagerMockRecorder) GetNextVerifierPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextVerifierPeers", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).GetNextVerifierPeers))
}

// GetPeer mocks base method
func (m *MockAbstractPbftProtocolManager) GetPeer(arg0 string) chain_communication.PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeer", arg0)
	ret0, _ := ret[0].(chain_communication.PmAbstractPeer)
	return ret0
}

// GetPeer indicates an expected call of GetPeer
func (mr *MockAbstractPbftProtocolManagerMockRecorder) GetPeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeer", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).GetPeer), arg0)
}

// GetPeers mocks base method
func (m *MockAbstractPbftProtocolManager) GetPeers() map[string]chain_communication.PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeers")
	ret0, _ := ret[0].(map[string]chain_communication.PmAbstractPeer)
	return ret0
}

// GetPeers indicates an expected call of GetPeers
func (mr *MockAbstractPbftProtocolManagerMockRecorder) GetPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeers", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).GetPeers))
}

// GetSelfNode mocks base method
func (m *MockAbstractPbftProtocolManager) GetSelfNode() *enode.Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelfNode")
	ret0, _ := ret[0].(*enode.Node)
	return ret0
}

// GetSelfNode indicates an expected call of GetSelfNode
func (mr *MockAbstractPbftProtocolManagerMockRecorder) GetSelfNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelfNode", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).GetSelfNode))
}

// GetVerifierBootNode mocks base method
func (m *MockAbstractPbftProtocolManager) GetVerifierBootNode() map[string]chain_communication.PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVerifierBootNode")
	ret0, _ := ret[0].(map[string]chain_communication.PmAbstractPeer)
	return ret0
}

// GetVerifierBootNode indicates an expected call of GetVerifierBootNode
func (mr *MockAbstractPbftProtocolManagerMockRecorder) GetVerifierBootNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVerifierBootNode", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).GetVerifierBootNode))
}

// IsSync mocks base method
func (m *MockAbstractPbftProtocolManager) IsSync() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSync")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSync indicates an expected call of IsSync
func (mr *MockAbstractPbftProtocolManagerMockRecorder) IsSync() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSync", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).IsSync))
}

// MatchCurrentVerifiersToNext mocks base method
func (m *MockAbstractPbftProtocolManager) MatchCurrentVerifiersToNext() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MatchCurrentVerifiersToNext")
}

// MatchCurrentVerifiersToNext indicates an expected call of MatchCurrentVerifiersToNext
func (mr *MockAbstractPbftProtocolManagerMockRecorder) MatchCurrentVerifiersToNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MatchCurrentVerifiersToNext", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).MatchCurrentVerifiersToNext))
}

// RemovePeer mocks base method
func (m *MockAbstractPbftProtocolManager) RemovePeer(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemovePeer", arg0)
}

// RemovePeer indicates an expected call of RemovePeer
func (mr *MockAbstractPbftProtocolManagerMockRecorder) RemovePeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePeer", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).RemovePeer), arg0)
}

// SelfIsBootNode mocks base method
func (m *MockAbstractPbftProtocolManager) SelfIsBootNode() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelfIsBootNode")
	ret0, _ := ret[0].(bool)
	return ret0
}

// SelfIsBootNode indicates an expected call of SelfIsBootNode
func (mr *MockAbstractPbftProtocolManagerMockRecorder) SelfIsBootNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelfIsBootNode", reflect.TypeOf((*MockAbstractPbftProtocolManager)(nil).SelfIsBootNode))
}
