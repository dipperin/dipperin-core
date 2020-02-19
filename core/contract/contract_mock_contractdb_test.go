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
// Source: github.com/dipperin/dipperin-core/core/contract (interfaces: ContractDB)

// Package contract is a generated GoMock package.
package contract

import (
	common "github.com/dipperin/dipperin-core/common"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockContractDB is a mock of ContractDB interface
type MockContractDB struct {
	ctrl     *gomock.Controller
	recorder *MockContractDBMockRecorder
}

// MockContractDBMockRecorder is the mock recorder for MockContractDB
type MockContractDBMockRecorder struct {
	mock *MockContractDB
}

// NewMockContractDB creates a new mock instance
func NewMockContractDB(ctrl *gomock.Controller) *MockContractDB {
	mock := &MockContractDB{ctrl: ctrl}
	mock.recorder = &MockContractDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContractDB) EXPECT() *MockContractDBMockRecorder {
	return m.recorder
}

// ContractExist mocks base method
func (m *MockContractDB) ContractExist(arg0 common.Address) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ContractExist", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ContractExist indicates an expected call of ContractExist
func (mr *MockContractDBMockRecorder) ContractExist(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ContractExist", reflect.TypeOf((*MockContractDB)(nil).ContractExist), arg0)
}

// GetContract mocks base method
func (m *MockContractDB) GetContract(arg0 common.Address, arg1 reflect.Type) (reflect.Value, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContract", arg0, arg1)
	ret0, _ := ret[0].(reflect.Value)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContract indicates an expected call of GetContract
func (mr *MockContractDBMockRecorder) GetContract(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContract", reflect.TypeOf((*MockContractDB)(nil).GetContract), arg0, arg1)
}

// PutContract mocks base method
func (m *MockContractDB) PutContract(arg0 common.Address, arg1 reflect.Value) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutContract", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutContract indicates an expected call of PutContract
func (mr *MockContractDBMockRecorder) PutContract(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutContract", reflect.TypeOf((*MockContractDB)(nil).PutContract), arg0, arg1)
}
