// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/chain/stateprocessor (interfaces: StateStorage)

// Package cschain is a generated GoMock package.
package cschain

import (
	common "github.com/dipperin/dipperin-core/common"
	stateprocessor "github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	trie "github.com/dipperin/dipperin-core/third_party/trie"
	ethdb "github.com/ethereum/go-ethereum/ethdb"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStateStorage is a mock of StateStorage interface
type MockStateStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStateStorageMockRecorder
}

// MockStateStorageMockRecorder is the mock recorder for MockStateStorage
type MockStateStorageMockRecorder struct {
	mock *MockStateStorage
}

// NewMockStateStorage creates a new mock instance
func NewMockStateStorage(ctrl *gomock.Controller) *MockStateStorage {
	mock := &MockStateStorage{ctrl: ctrl}
	mock.recorder = &MockStateStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStateStorage) EXPECT() *MockStateStorageMockRecorder {
	return m.recorder
}

// CopyTrie mocks base method
func (m *MockStateStorage) CopyTrie(arg0 stateprocessor.StateTrie) stateprocessor.StateTrie {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyTrie", arg0)
	ret0, _ := ret[0].(stateprocessor.StateTrie)
	return ret0
}

// CopyTrie indicates an expected call of CopyTrie
func (mr *MockStateStorageMockRecorder) CopyTrie(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyTrie", reflect.TypeOf((*MockStateStorage)(nil).CopyTrie), arg0)
}

// DiskDB mocks base method
func (m *MockStateStorage) DiskDB() ethdb.Database {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiskDB")
	ret0, _ := ret[0].(ethdb.Database)
	return ret0
}

// DiskDB indicates an expected call of DiskDB
func (mr *MockStateStorageMockRecorder) DiskDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiskDB", reflect.TypeOf((*MockStateStorage)(nil).DiskDB))
}

// OpenStorageTrie mocks base method
func (m *MockStateStorage) OpenStorageTrie(arg0, arg1 common.Hash) (stateprocessor.StateTrie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenStorageTrie", arg0, arg1)
	ret0, _ := ret[0].(stateprocessor.StateTrie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenStorageTrie indicates an expected call of OpenStorageTrie
func (mr *MockStateStorageMockRecorder) OpenStorageTrie(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenStorageTrie", reflect.TypeOf((*MockStateStorage)(nil).OpenStorageTrie), arg0, arg1)
}

// OpenTrie mocks base method
func (m *MockStateStorage) OpenTrie(arg0 common.Hash) (stateprocessor.StateTrie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenTrie", arg0)
	ret0, _ := ret[0].(stateprocessor.StateTrie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenTrie indicates an expected call of OpenTrie
func (mr *MockStateStorageMockRecorder) OpenTrie(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenTrie", reflect.TypeOf((*MockStateStorage)(nil).OpenTrie), arg0)
}

// TrieDB mocks base method
func (m *MockStateStorage) TrieDB() *trie.Database {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TrieDB")
	ret0, _ := ret[0].(*trie.Database)
	return ret0
}

// TrieDB indicates an expected call of TrieDB
func (mr *MockStateStorageMockRecorder) TrieDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TrieDB", reflect.TypeOf((*MockStateStorage)(nil).TrieDB))
}