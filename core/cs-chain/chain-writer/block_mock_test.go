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
// Source: github.com/caiqingfeng/dipperin-core/core/model (interfaces: AbstractBlock)

// Package chain_writer is a generated GoMock package.
package chain_writer

import (
	common "github.com/dipperin/dipperin-core/common"
	bloom "github.com/dipperin/dipperin-core/core/bloom"
	model "github.com/dipperin/dipperin-core/core/model"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
)

// MockAbstractBlock is a mock of AbstractBlock interface
type MockAbstractBlock struct {
	ctrl     *gomock.Controller
	recorder *MockAbstractBlockMockRecorder
}

// MockAbstractBlockMockRecorder is the mock recorder for MockAbstractBlock
type MockAbstractBlockMockRecorder struct {
	mock *MockAbstractBlock
}

// NewMockAbstractBlock creates a new mock instance
func NewMockAbstractBlock(ctrl *gomock.Controller) *MockAbstractBlock {
	mock := &MockAbstractBlock{ctrl: ctrl}
	mock.recorder = &MockAbstractBlockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAbstractBlock) EXPECT() *MockAbstractBlockMockRecorder {
	return m.recorder
}

// Body mocks base method
func (m *MockAbstractBlock) Body() model.AbstractBody {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Body")
	ret0, _ := ret[0].(model.AbstractBody)
	return ret0
}

// Body indicates an expected call of Body
func (mr *MockAbstractBlockMockRecorder) Body() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Body", reflect.TypeOf((*MockAbstractBlock)(nil).Body))
}

// CoinBase mocks base method
func (m *MockAbstractBlock) CoinBase() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CoinBase")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// CoinBase indicates an expected call of CoinBase
func (mr *MockAbstractBlockMockRecorder) CoinBase() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CoinBase", reflect.TypeOf((*MockAbstractBlock)(nil).CoinBase))
}

// CoinBaseAddress mocks base method
func (m *MockAbstractBlock) CoinBaseAddress() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CoinBaseAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// CoinBaseAddress indicates an expected call of CoinBaseAddress
func (mr *MockAbstractBlockMockRecorder) CoinBaseAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CoinBaseAddress", reflect.TypeOf((*MockAbstractBlock)(nil).CoinBaseAddress))
}

// Difficulty mocks base method
func (m *MockAbstractBlock) Difficulty() common.Difficulty {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Difficulty")
	ret0, _ := ret[0].(common.Difficulty)
	return ret0
}

// Difficulty indicates an expected call of Difficulty
func (mr *MockAbstractBlockMockRecorder) Difficulty() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Difficulty", reflect.TypeOf((*MockAbstractBlock)(nil).Difficulty))
}

// EncodeRlpToBytes mocks base method
func (m *MockAbstractBlock) EncodeRlpToBytes() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncodeRlpToBytes")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EncodeRlpToBytes indicates an expected call of EncodeRlpToBytes
func (mr *MockAbstractBlockMockRecorder) EncodeRlpToBytes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodeRlpToBytes", reflect.TypeOf((*MockAbstractBlock)(nil).EncodeRlpToBytes))
}

// FormatForRpc mocks base method
func (m *MockAbstractBlock) FormatForRpc() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FormatForRpc")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// FormatForRpc indicates an expected call of FormatForRpc
func (mr *MockAbstractBlockMockRecorder) FormatForRpc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FormatForRpc", reflect.TypeOf((*MockAbstractBlock)(nil).FormatForRpc))
}

// GetAbsTransactions mocks base method
func (m *MockAbstractBlock) GetAbsTransactions() []model.AbstractTransaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAbsTransactions")
	ret0, _ := ret[0].([]model.AbstractTransaction)
	return ret0
}

// GetAbsTransactions indicates an expected call of GetAbsTransactions
func (mr *MockAbstractBlockMockRecorder) GetAbsTransactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAbsTransactions", reflect.TypeOf((*MockAbstractBlock)(nil).GetAbsTransactions))
}

// GetBlockTxsBloom mocks base method
func (m *MockAbstractBlock) GetBlockTxsBloom() *bloom.Bloom {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockTxsBloom")
	ret0, _ := ret[0].(*bloom.Bloom)
	return ret0
}

// GetBlockTxsBloom indicates an expected call of GetBlockTxsBloom
func (mr *MockAbstractBlockMockRecorder) GetBlockTxsBloom() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockTxsBloom", reflect.TypeOf((*MockAbstractBlock)(nil).GetBlockTxsBloom))
}

// GetBloom mocks base method
func (m *MockAbstractBlock) GetBloom() bloom.Bloom {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBloom")
	ret0, _ := ret[0].(bloom.Bloom)
	return ret0
}

// GetBloom indicates an expected call of GetBloom
func (mr *MockAbstractBlockMockRecorder) GetBloom() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBloom", reflect.TypeOf((*MockAbstractBlock)(nil).GetBloom))
}

// GetEiBloomBlockData mocks base method
func (m *MockAbstractBlock) GetEiBloomBlockData(arg0 *bloom.HybridEstimator) *model.BloomBlockData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEiBloomBlockData", arg0)
	ret0, _ := ret[0].(*model.BloomBlockData)
	return ret0
}

// GetEiBloomBlockData indicates an expected call of GetEiBloomBlockData
func (mr *MockAbstractBlockMockRecorder) GetEiBloomBlockData(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEiBloomBlockData", reflect.TypeOf((*MockAbstractBlock)(nil).GetEiBloomBlockData), arg0)
}

// GetInterLinkRoot mocks base method
func (m *MockAbstractBlock) GetInterLinkRoot() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInterLinkRoot")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// GetInterLinkRoot indicates an expected call of GetInterLinkRoot
func (mr *MockAbstractBlockMockRecorder) GetInterLinkRoot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInterLinkRoot", reflect.TypeOf((*MockAbstractBlock)(nil).GetInterLinkRoot))
}

// GetInterlinks mocks base method
func (m *MockAbstractBlock) GetInterlinks() model.InterLink {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInterlinks")
	ret0, _ := ret[0].(model.InterLink)
	return ret0
}

// GetInterlinks indicates an expected call of GetInterlinks
func (mr *MockAbstractBlockMockRecorder) GetInterlinks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInterlinks", reflect.TypeOf((*MockAbstractBlock)(nil).GetInterlinks))
}

// GetRegisterRoot mocks base method
func (m *MockAbstractBlock) GetRegisterRoot() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRegisterRoot")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// GetRegisterRoot indicates an expected call of GetRegisterRoot
func (mr *MockAbstractBlockMockRecorder) GetRegisterRoot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRegisterRoot", reflect.TypeOf((*MockAbstractBlock)(nil).GetRegisterRoot))
}

// GetTransactionFees mocks base method
func (m *MockAbstractBlock) GetTransactionFees() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionFees")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// GetTransactionFees indicates an expected call of GetTransactionFees
func (mr *MockAbstractBlockMockRecorder) GetTransactionFees() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionFees", reflect.TypeOf((*MockAbstractBlock)(nil).GetTransactionFees))
}

// GetTransactions mocks base method
func (m *MockAbstractBlock) GetTransactions() []*model.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactions")
	ret0, _ := ret[0].([]*model.Transaction)
	return ret0
}

// GetTransactions indicates an expected call of GetTransactions
func (mr *MockAbstractBlockMockRecorder) GetTransactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactions", reflect.TypeOf((*MockAbstractBlock)(nil).GetTransactions))
}

// GetVerifications mocks base method
func (m *MockAbstractBlock) GetVerifications() []model.AbstractVerification {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVerifications")
	ret0, _ := ret[0].([]model.AbstractVerification)
	return ret0
}

// GetVerifications indicates an expected call of GetVerifications
func (mr *MockAbstractBlockMockRecorder) GetVerifications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVerifications", reflect.TypeOf((*MockAbstractBlock)(nil).GetVerifications))
}

// Hash mocks base method
func (m *MockAbstractBlock) Hash() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// Hash indicates an expected call of Hash
func (mr *MockAbstractBlockMockRecorder) Hash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockAbstractBlock)(nil).Hash))
}

// Header mocks base method
func (m *MockAbstractBlock) Header() model.AbstractHeader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(model.AbstractHeader)
	return ret0
}

// Header indicates an expected call of Header
func (mr *MockAbstractBlockMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockAbstractBlock)(nil).Header))
}

// IsSpecial mocks base method
func (m *MockAbstractBlock) IsSpecial() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsSpecial")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsSpecial indicates an expected call of IsSpecial
func (mr *MockAbstractBlockMockRecorder) IsSpecial() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsSpecial", reflect.TypeOf((*MockAbstractBlock)(nil).IsSpecial))
}

// Nonce mocks base method
func (m *MockAbstractBlock) Nonce() common.BlockNonce {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Nonce")
	ret0, _ := ret[0].(common.BlockNonce)
	return ret0
}

// Nonce indicates an expected call of Nonce
func (mr *MockAbstractBlockMockRecorder) Nonce() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Nonce", reflect.TypeOf((*MockAbstractBlock)(nil).Nonce))
}

// Number mocks base method
func (m *MockAbstractBlock) Number() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Number")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Number indicates an expected call of Number
func (mr *MockAbstractBlockMockRecorder) Number() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Number", reflect.TypeOf((*MockAbstractBlock)(nil).Number))
}

// PreHash mocks base method
func (m *MockAbstractBlock) PreHash() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreHash")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// PreHash indicates an expected call of PreHash
func (mr *MockAbstractBlockMockRecorder) PreHash() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreHash", reflect.TypeOf((*MockAbstractBlock)(nil).PreHash))
}

// RefreshHashCache mocks base method
func (m *MockAbstractBlock) RefreshHashCache() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshHashCache")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// RefreshHashCache indicates an expected call of RefreshHashCache
func (mr *MockAbstractBlockMockRecorder) RefreshHashCache() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshHashCache", reflect.TypeOf((*MockAbstractBlock)(nil).RefreshHashCache))
}

// Seed mocks base method
func (m *MockAbstractBlock) Seed() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seed")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// Seed indicates an expected call of Seed
func (mr *MockAbstractBlockMockRecorder) Seed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seed", reflect.TypeOf((*MockAbstractBlock)(nil).Seed))
}

// SetInterLinkRoot mocks base method
func (m *MockAbstractBlock) SetInterLinkRoot(arg0 common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetInterLinkRoot", arg0)
}

// SetInterLinkRoot indicates an expected call of SetInterLinkRoot
func (mr *MockAbstractBlockMockRecorder) SetInterLinkRoot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInterLinkRoot", reflect.TypeOf((*MockAbstractBlock)(nil).SetInterLinkRoot), arg0)
}

// SetInterLinks mocks base method
func (m *MockAbstractBlock) SetInterLinks(arg0 model.InterLink) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetInterLinks", arg0)
}

// SetInterLinks indicates an expected call of SetInterLinks
func (mr *MockAbstractBlockMockRecorder) SetInterLinks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInterLinks", reflect.TypeOf((*MockAbstractBlock)(nil).SetInterLinks), arg0)
}

// SetNonce mocks base method
func (m *MockAbstractBlock) SetNonce(arg0 common.BlockNonce) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNonce", arg0)
}

// SetNonce indicates an expected call of SetNonce
func (mr *MockAbstractBlockMockRecorder) SetNonce(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNonce", reflect.TypeOf((*MockAbstractBlock)(nil).SetNonce), arg0)
}

// SetRegisterRoot mocks base method
func (m *MockAbstractBlock) SetRegisterRoot(arg0 common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRegisterRoot", arg0)
}

// SetRegisterRoot indicates an expected call of SetRegisterRoot
func (mr *MockAbstractBlockMockRecorder) SetRegisterRoot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRegisterRoot", reflect.TypeOf((*MockAbstractBlock)(nil).SetRegisterRoot), arg0)
}

// SetStateRoot mocks base method
func (m *MockAbstractBlock) SetStateRoot(arg0 common.Hash) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStateRoot", arg0)
}

// SetStateRoot indicates an expected call of SetStateRoot
func (mr *MockAbstractBlockMockRecorder) SetStateRoot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStateRoot", reflect.TypeOf((*MockAbstractBlock)(nil).SetStateRoot), arg0)
}

// SetVerifications mocks base method
func (m *MockAbstractBlock) SetVerifications(arg0 []model.AbstractVerification) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetVerifications", arg0)
}

// SetVerifications indicates an expected call of SetVerifications
func (mr *MockAbstractBlockMockRecorder) SetVerifications(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetVerifications", reflect.TypeOf((*MockAbstractBlock)(nil).SetVerifications), arg0)
}

// StateRoot mocks base method
func (m *MockAbstractBlock) StateRoot() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StateRoot")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// StateRoot indicates an expected call of StateRoot
func (mr *MockAbstractBlockMockRecorder) StateRoot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StateRoot", reflect.TypeOf((*MockAbstractBlock)(nil).StateRoot))
}

// Timestamp mocks base method
func (m *MockAbstractBlock) Timestamp() *big.Int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Timestamp")
	ret0, _ := ret[0].(*big.Int)
	return ret0
}

// Timestamp indicates an expected call of Timestamp
func (mr *MockAbstractBlockMockRecorder) Timestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Timestamp", reflect.TypeOf((*MockAbstractBlock)(nil).Timestamp))
}

// TxCount mocks base method
func (m *MockAbstractBlock) TxCount() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxCount")
	ret0, _ := ret[0].(int)
	return ret0
}

// TxCount indicates an expected call of TxCount
func (mr *MockAbstractBlockMockRecorder) TxCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxCount", reflect.TypeOf((*MockAbstractBlock)(nil).TxCount))
}

// TxIterator mocks base method
func (m *MockAbstractBlock) TxIterator(arg0 func(int, model.AbstractTransaction) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxIterator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// TxIterator indicates an expected call of TxIterator
func (mr *MockAbstractBlockMockRecorder) TxIterator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxIterator", reflect.TypeOf((*MockAbstractBlock)(nil).TxIterator), arg0)
}

// TxRoot mocks base method
func (m *MockAbstractBlock) TxRoot() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxRoot")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// TxRoot indicates an expected call of TxRoot
func (mr *MockAbstractBlockMockRecorder) TxRoot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxRoot", reflect.TypeOf((*MockAbstractBlock)(nil).TxRoot))
}

// VerificationRoot mocks base method
func (m *MockAbstractBlock) VerificationRoot() common.Hash {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerificationRoot")
	ret0, _ := ret[0].(common.Hash)
	return ret0
}

// VerificationRoot indicates an expected call of VerificationRoot
func (mr *MockAbstractBlockMockRecorder) VerificationRoot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerificationRoot", reflect.TypeOf((*MockAbstractBlock)(nil).VerificationRoot))
}

// VersIterator mocks base method
func (m *MockAbstractBlock) VersIterator(arg0 func(int, model.AbstractVerification, model.AbstractBlock) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VersIterator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// VersIterator indicates an expected call of VersIterator
func (mr *MockAbstractBlockMockRecorder) VersIterator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VersIterator", reflect.TypeOf((*MockAbstractBlock)(nil).VersIterator), arg0)
}

// Version mocks base method
func (m *MockAbstractBlock) Version() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Version")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// Version indicates an expected call of Version
func (mr *MockAbstractBlockMockRecorder) Version() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Version", reflect.TypeOf((*MockAbstractBlock)(nil).Version))
}
