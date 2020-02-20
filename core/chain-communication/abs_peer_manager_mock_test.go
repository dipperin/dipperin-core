// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/chain-communication (interfaces: AbsPeerManager)

// Package chain_communication is a generated GoMock package.
package chain_communication

import (
	//chain_communication "github.com/dipperin/dipperin-core/core/chain-communication"
	enode "github.com/dipperin/dipperin-core/third-party/p2p/enode"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAbsPeerManager is a mock of AbsPeerManager interface
type MockAbsPeerManager struct {
	ctrl     *gomock.Controller
	recorder *MockAbsPeerManagerMockRecorder
}

// MockAbsPeerManagerMockRecorder is the mock recorder for MockAbsPeerManager
type MockAbsPeerManagerMockRecorder struct {
	mock *MockAbsPeerManager
}

// NewMockAbsPeerManager creates a new mock instance
func NewMockAbsPeerManager(ctrl *gomock.Controller) *MockAbsPeerManager {
	mock := &MockAbsPeerManager{ctrl: ctrl}
	mock.recorder = &MockAbsPeerManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAbsPeerManager) EXPECT() *MockAbsPeerManagerMockRecorder {
	return m.recorder
}

// BestPeer mocks base method
func (m *MockAbsPeerManager) BestPeer() PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BestPeer")
	ret0, _ := ret[0].(PmAbstractPeer)
	return ret0
}

// BestPeer indicates an expected call of BestPeer
func (mr *MockAbsPeerManagerMockRecorder) BestPeer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BestPeer", reflect.TypeOf((*MockAbsPeerManager)(nil).BestPeer))
}

// ConnectPeer mocks base method
func (m *MockAbsPeerManager) ConnectPeer(arg0 *enode.Node) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ConnectPeer", arg0)
}

// ConnectPeer indicates an expected call of ConnectPeer
func (mr *MockAbsPeerManagerMockRecorder) ConnectPeer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectPeer", reflect.TypeOf((*MockAbsPeerManager)(nil).ConnectPeer), arg0)
}

// CurrentVerifierPeersSet mocks base method
func (m *MockAbsPeerManager) CurrentVerifierPeersSet() AbstractPeerSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentVerifierPeersSet")
	ret0, _ := ret[0].(AbstractPeerSet)
	return ret0
}

// CurrentVerifierPeersSet indicates an expected call of CurrentVerifierPeersSet
func (mr *MockAbsPeerManagerMockRecorder) CurrentVerifierPeersSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentVerifierPeersSet", reflect.TypeOf((*MockAbsPeerManager)(nil).CurrentVerifierPeersSet))
}

// GetSelfNode mocks base method
func (m *MockAbsPeerManager) GetSelfNode() *enode.Node {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSelfNode")
	ret0, _ := ret[0].(*enode.Node)
	return ret0
}

// GetSelfNode indicates an expected call of GetSelfNode
func (mr *MockAbsPeerManagerMockRecorder) GetSelfNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSelfNode", reflect.TypeOf((*MockAbsPeerManager)(nil).GetSelfNode))
}

// GetVerifierBootNode mocks base method
func (m *MockAbsPeerManager) GetVerifierBootNode() map[string]PmAbstractPeer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVerifierBootNode")
	ret0, _ := ret[0].(map[string]PmAbstractPeer)
	return ret0
}

// GetVerifierBootNode indicates an expected call of GetVerifierBootNode
func (mr *MockAbsPeerManagerMockRecorder) GetVerifierBootNode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVerifierBootNode", reflect.TypeOf((*MockAbsPeerManager)(nil).GetVerifierBootNode))
}

// HaveEnoughVerifiers mocks base method
func (m *MockAbsPeerManager) HaveEnoughVerifiers(arg0 bool) (uint, uint) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HaveEnoughVerifiers", arg0)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(uint)
	return ret0, ret1
}

// HaveEnoughVerifiers indicates an expected call of HaveEnoughVerifiers
func (mr *MockAbsPeerManagerMockRecorder) HaveEnoughVerifiers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HaveEnoughVerifiers", reflect.TypeOf((*MockAbsPeerManager)(nil).HaveEnoughVerifiers), arg0)
}

// NextVerifierPeersSet mocks base method
func (m *MockAbsPeerManager) NextVerifierPeersSet() AbstractPeerSet {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NextVerifierPeersSet")
	ret0, _ := ret[0].(AbstractPeerSet)
	return ret0
}

// NextVerifierPeersSet indicates an expected call of NextVerifierPeersSet
func (mr *MockAbsPeerManagerMockRecorder) NextVerifierPeersSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NextVerifierPeersSet", reflect.TypeOf((*MockAbsPeerManager)(nil).NextVerifierPeersSet))
}

// SelfIsCurrentVerifier mocks base method
func (m *MockAbsPeerManager) SelfIsCurrentVerifier() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelfIsCurrentVerifier")
	ret0, _ := ret[0].(bool)
	return ret0
}

// SelfIsCurrentVerifier indicates an expected call of SelfIsCurrentVerifier
func (mr *MockAbsPeerManagerMockRecorder) SelfIsCurrentVerifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelfIsCurrentVerifier", reflect.TypeOf((*MockAbsPeerManager)(nil).SelfIsCurrentVerifier))
}

// SelfIsNextVerifier mocks base method
func (m *MockAbsPeerManager) SelfIsNextVerifier() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelfIsNextVerifier")
	ret0, _ := ret[0].(bool)
	return ret0
}

// SelfIsNextVerifier indicates an expected call of SelfIsNextVerifier
func (mr *MockAbsPeerManagerMockRecorder) SelfIsNextVerifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelfIsNextVerifier", reflect.TypeOf((*MockAbsPeerManager)(nil).SelfIsNextVerifier))
}