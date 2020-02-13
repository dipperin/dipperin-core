// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/accounts/accountsbase (interfaces: Wallet)

// Package softwallet is a generated GoMock package.
package softwallet

import (
	ecdsa "crypto/ecdsa"
	common "github.com/dipperin/dipperin-core/common"
	accountsbase "github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	model "github.com/dipperin/dipperin-core/core/model"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

// MockWallet is a mock of Wallet interface
type MockWallet struct {
	ctrl     *gomock.Controller
	recorder *MockWalletMockRecorder
}

// MockWalletMockRecorder is the mock recorder for MockWallet
type MockWalletMockRecorder struct {
	mock *MockWallet
}

// NewMockWallet creates a new mock instance
func NewMockWallet(ctrl *gomock.Controller) *MockWallet {
	mock := &MockWallet{ctrl: ctrl}
	mock.recorder = &MockWalletMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWallet) EXPECT() *MockWalletMockRecorder {
	return m.recorder
}

// Accounts mocks base method
func (m *MockWallet) Accounts() ([]accountsbase.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accounts")
	ret0, _ := ret[0].([]accountsbase.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accounts indicates an expected call of Accounts
func (mr *MockWalletMockRecorder) Accounts() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accounts", reflect.TypeOf((*MockWallet)(nil).Accounts))
}

// Close mocks base method
func (m *MockWallet) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockWalletMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWallet)(nil).Close))
}

// Contains mocks base method
func (m *MockWallet) Contains(arg0 accountsbase.Account) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Contains", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Contains indicates an expected call of Contains
func (mr *MockWalletMockRecorder) Contains(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Contains", reflect.TypeOf((*MockWallet)(nil).Contains), arg0)
}

// Derive mocks base method
func (m *MockWallet) Derive(arg0 accountsbase.DerivationPath, arg1 bool) (accountsbase.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Derive", arg0, arg1)
	ret0, _ := ret[0].(accountsbase.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Derive indicates an expected call of Derive
func (mr *MockWalletMockRecorder) Derive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Derive", reflect.TypeOf((*MockWallet)(nil).Derive), arg0, arg1)
}

// Establish mocks base method
func (m *MockWallet) Establish(arg0, arg1, arg2, arg3 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Establish", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Establish indicates an expected call of Establish
func (mr *MockWalletMockRecorder) Establish(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Establish", reflect.TypeOf((*MockWallet)(nil).Establish), arg0, arg1, arg2, arg3)
}

// Evaluate mocks base method
func (m *MockWallet) Evaluate(arg0 accountsbase.Account, arg1 []byte) ([32]byte, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Evaluate", arg0, arg1)
	ret0, _ := ret[0].([32]byte)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Evaluate indicates an expected call of Evaluate
func (mr *MockWalletMockRecorder) Evaluate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Evaluate", reflect.TypeOf((*MockWallet)(nil).Evaluate), arg0, arg1)
}

// GetAddressNonce mocks base method
func (m *MockWallet) GetAddressNonce(arg0 common.Address) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddressNonce", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAddressNonce indicates an expected call of GetAddressNonce
func (mr *MockWalletMockRecorder) GetAddressNonce(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddressNonce", reflect.TypeOf((*MockWallet)(nil).GetAddressNonce), arg0)
}

// GetPKFromAddress mocks base method
func (m *MockWallet) GetPKFromAddress(arg0 accountsbase.Account) (*ecdsa.PublicKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPKFromAddress", arg0)
	ret0, _ := ret[0].(*ecdsa.PublicKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPKFromAddress indicates an expected call of GetPKFromAddress
func (mr *MockWalletMockRecorder) GetPKFromAddress(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPKFromAddress", reflect.TypeOf((*MockWallet)(nil).GetPKFromAddress), arg0)
}

// GetSKFromAddress mocks base method
func (m *MockWallet) GetSKFromAddress(arg0 common.Address) (*ecdsa.PrivateKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSKFromAddress", arg0)
	ret0, _ := ret[0].(*ecdsa.PrivateKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSKFromAddress indicates an expected call of GetSKFromAddress
func (mr *MockWalletMockRecorder) GetSKFromAddress(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSKFromAddress", reflect.TypeOf((*MockWallet)(nil).GetSKFromAddress), arg0)
}

// GetWalletIdentifier mocks base method
func (m *MockWallet) GetWalletIdentifier() (accountsbase.WalletIdentifier, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletIdentifier")
	ret0, _ := ret[0].(accountsbase.WalletIdentifier)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletIdentifier indicates an expected call of GetWalletIdentifier
func (mr *MockWalletMockRecorder) GetWalletIdentifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletIdentifier", reflect.TypeOf((*MockWallet)(nil).GetWalletIdentifier))
}

// Open mocks base method
func (m *MockWallet) Open(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open
func (mr *MockWalletMockRecorder) Open(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockWallet)(nil).Open), arg0, arg1, arg2)
}

// PaddingAddressNonce mocks base method
func (m *MockWallet) PaddingAddressNonce(arg0 accountsbase.AddressInfoReader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaddingAddressNonce", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// PaddingAddressNonce indicates an expected call of PaddingAddressNonce
func (mr *MockWalletMockRecorder) PaddingAddressNonce(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaddingAddressNonce", reflect.TypeOf((*MockWallet)(nil).PaddingAddressNonce), arg0)
}

// RestoreWallet mocks base method
func (m *MockWallet) RestoreWallet(arg0, arg1, arg2, arg3, arg4 string, arg5 accountsbase.AddressInfoReader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestoreWallet", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreWallet indicates an expected call of RestoreWallet
func (mr *MockWalletMockRecorder) RestoreWallet(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreWallet", reflect.TypeOf((*MockWallet)(nil).RestoreWallet), arg0, arg1, arg2, arg3, arg4, arg5)
}

// SelfDerive mocks base method
func (m *MockWallet) SelfDerive(arg0 accountsbase.DerivationPath) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelfDerive", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelfDerive indicates an expected call of SelfDerive
func (mr *MockWalletMockRecorder) SelfDerive(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelfDerive", reflect.TypeOf((*MockWallet)(nil).SelfDerive), arg0)
}

// SetAddressNonce mocks base method
func (m *MockWallet) SetAddressNonce(arg0 common.Address, arg1 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAddressNonce", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAddressNonce indicates an expected call of SetAddressNonce
func (mr *MockWalletMockRecorder) SetAddressNonce(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAddressNonce", reflect.TypeOf((*MockWallet)(nil).SetAddressNonce), arg0, arg1)
}

// SignHash mocks base method
func (m *MockWallet) SignHash(arg0 accountsbase.Account, arg1 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignHash", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignHash indicates an expected call of SignHash
func (mr *MockWalletMockRecorder) SignHash(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignHash", reflect.TypeOf((*MockWallet)(nil).SignHash), arg0, arg1)
}

// SignTx mocks base method
func (m *MockWallet) SignTx(arg0 accountsbase.Account, arg1 *model.Transaction, arg2 *big.Int) (*model.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignTx", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignTx indicates an expected call of SignTx
func (mr *MockWalletMockRecorder) SignTx(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignTx", reflect.TypeOf((*MockWallet)(nil).SignTx), arg0, arg1, arg2)
}

// Status mocks base method
func (m *MockWallet) Status() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status
func (mr *MockWalletMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockWallet)(nil).Status))
}
