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
// Source: github.com/caiqingfeng/dipperin-core/core/chain-communication (interfaces: PbftSigner)

// Package builder is a generated GoMock package.
package builder

import (
	ecdsa "crypto/ecdsa"
	common "github.com/dipperin/dipperin-core/common"
	accounts "github.com/dipperin/dipperin-core/core/accounts"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockPbftSigner is a mock of PbftSigner interface
type MockPbftSigner struct {
	ctrl     *gomock.Controller
	recorder *MockPbftSignerMockRecorder
}

// MockPbftSignerMockRecorder is the mock recorder for MockPbftSigner
type MockPbftSignerMockRecorder struct {
	mock *MockPbftSigner
}

// NewMockPbftSigner creates a new mock instance
func NewMockPbftSigner(ctrl *gomock.Controller) *MockPbftSigner {
	mock := &MockPbftSigner{ctrl: ctrl}
	mock.recorder = &MockPbftSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPbftSigner) EXPECT() *MockPbftSignerMockRecorder {
	return m.recorder
}

// Evaluate mocks base method
func (m *MockPbftSigner) Evaluate(arg0 accounts.Account, arg1 []byte) ([32]byte, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Evaluate", arg0, arg1)
	ret0, _ := ret[0].([32]byte)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Evaluate indicates an expected call of Evaluate
func (mr *MockPbftSignerMockRecorder) Evaluate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Evaluate", reflect.TypeOf((*MockPbftSigner)(nil).Evaluate), arg0, arg1)
}

// GetAddress mocks base method
func (m *MockPbftSigner) GetAddress() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetAddress indicates an expected call of GetAddress
func (mr *MockPbftSignerMockRecorder) GetAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockPbftSigner)(nil).GetAddress))
}

// PublicKey mocks base method
func (m *MockPbftSigner) PublicKey() *ecdsa.PublicKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublicKey")
	ret0, _ := ret[0].(*ecdsa.PublicKey)
	return ret0
}

// PublicKey indicates an expected call of PublicKey
func (mr *MockPbftSignerMockRecorder) PublicKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicKey", reflect.TypeOf((*MockPbftSigner)(nil).PublicKey))
}

// SetBaseAddress mocks base method
func (m *MockPbftSigner) SetBaseAddress(arg0 common.Address) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBaseAddress", arg0)
}

// SetBaseAddress indicates an expected call of SetBaseAddress
func (mr *MockPbftSignerMockRecorder) SetBaseAddress(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBaseAddress", reflect.TypeOf((*MockPbftSigner)(nil).SetBaseAddress), arg0)
}

// SignHash mocks base method
func (m *MockPbftSigner) SignHash(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignHash", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignHash indicates an expected call of SignHash
func (mr *MockPbftSignerMockRecorder) SignHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignHash", reflect.TypeOf((*MockPbftSigner)(nil).SignHash), arg0)
}

// ValidSign mocks base method
func (m *MockPbftSigner) ValidSign(arg0, arg1, arg2 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidSign", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidSign indicates an expected call of ValidSign
func (mr *MockPbftSignerMockRecorder) ValidSign(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidSign", reflect.TypeOf((*MockPbftSigner)(nil).ValidSign), arg0, arg1, arg2)
}
