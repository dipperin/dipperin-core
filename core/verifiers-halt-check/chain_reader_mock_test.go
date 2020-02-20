// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/verifiers-halt-check (interfaces: NeedChainReaderFunction)

// Package verifiers_halt_check is a generated GoMock package.
package verifiers_halt_check

import (
	common "github.com/dipperin/dipperin-core/common"
	chain "github.com/dipperin/dipperin-core/core/chain"
	registerdb "github.com/dipperin/dipperin-core/core/chain/registerdb"
	model "github.com/dipperin/dipperin-core/core/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockNeedChainReaderFunction is a mock of NeedChainReaderFunction interface
type MockNeedChainReaderFunction struct {
	ctrl     *gomock.Controller
	recorder *MockNeedChainReaderFunctionMockRecorder
}

// MockNeedChainReaderFunctionMockRecorder is the mock recorder for MockNeedChainReaderFunction
type MockNeedChainReaderFunctionMockRecorder struct {
	mock *MockNeedChainReaderFunction
}

// NewMockNeedChainReaderFunction creates a new mock instance
func NewMockNeedChainReaderFunction(ctrl *gomock.Controller) *MockNeedChainReaderFunction {
	mock := &MockNeedChainReaderFunction{ctrl: ctrl}
	mock.recorder = &MockNeedChainReaderFunctionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNeedChainReaderFunction) EXPECT() *MockNeedChainReaderFunctionMockRecorder {
	return m.recorder
}

// BlockProcessor mocks base method
func (m *MockNeedChainReaderFunction) BlockProcessor(arg0 common.Hash) (*chain.BlockProcessor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockProcessor", arg0)
	ret0, _ := ret[0].(*chain.BlockProcessor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockProcessor indicates an expected call of BlockProcessor
func (mr *MockNeedChainReaderFunctionMockRecorder) BlockProcessor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockProcessor", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).BlockProcessor), arg0)
}

// BlockProcessorByNumber mocks base method
func (m *MockNeedChainReaderFunction) BlockProcessorByNumber(arg0 uint64) (*chain.BlockProcessor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockProcessorByNumber", arg0)
	ret0, _ := ret[0].(*chain.BlockProcessor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockProcessorByNumber indicates an expected call of BlockProcessorByNumber
func (mr *MockNeedChainReaderFunctionMockRecorder) BlockProcessorByNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockProcessorByNumber", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).BlockProcessorByNumber), arg0)
}

// BuildRegisterProcessor mocks base method
func (m *MockNeedChainReaderFunction) BuildRegisterProcessor(arg0 common.Hash) (*registerdb.RegisterDB, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildRegisterProcessor", arg0)
	ret0, _ := ret[0].(*registerdb.RegisterDB)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BuildRegisterProcessor indicates an expected call of BuildRegisterProcessor
func (mr *MockNeedChainReaderFunctionMockRecorder) BuildRegisterProcessor(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildRegisterProcessor", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).BuildRegisterProcessor), arg0)
}

// CurrentBlock mocks base method
func (m *MockNeedChainReaderFunction) CurrentBlock() model.AbstractBlock {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentBlock")
	ret0, _ := ret[0].(model.AbstractBlock)
	return ret0
}

// CurrentBlock indicates an expected call of CurrentBlock
func (mr *MockNeedChainReaderFunctionMockRecorder) CurrentBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentBlock", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).CurrentBlock))
}

// GetBlockByNumber mocks base method
func (m *MockNeedChainReaderFunction) GetBlockByNumber(arg0 uint64) model.AbstractBlock {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockByNumber", arg0)
	ret0, _ := ret[0].(model.AbstractBlock)
	return ret0
}

// GetBlockByNumber indicates an expected call of GetBlockByNumber
func (mr *MockNeedChainReaderFunctionMockRecorder) GetBlockByNumber(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByNumber", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).GetBlockByNumber), arg0)
}

// GetCurrVerifiers mocks base method
func (m *MockNeedChainReaderFunction) GetCurrVerifiers() []common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrVerifiers")
	ret0, _ := ret[0].([]common.Address)
	return ret0
}

// GetCurrVerifiers indicates an expected call of GetCurrVerifiers
func (mr *MockNeedChainReaderFunctionMockRecorder) GetCurrVerifiers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrVerifiers", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).GetCurrVerifiers))
}

// GetLastChangePoint mocks base method
func (m *MockNeedChainReaderFunction) GetLastChangePoint(arg0 model.AbstractBlock) *uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastChangePoint", arg0)
	ret0, _ := ret[0].(*uint64)
	return ret0
}

// GetLastChangePoint indicates an expected call of GetLastChangePoint
func (mr *MockNeedChainReaderFunctionMockRecorder) GetLastChangePoint(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastChangePoint", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).GetLastChangePoint), arg0)
}

// GetSeenCommit mocks base method
func (m *MockNeedChainReaderFunction) GetSeenCommit(arg0 uint64) []model.AbstractVerification {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeenCommit", arg0)
	ret0, _ := ret[0].([]model.AbstractVerification)
	return ret0
}

// GetSeenCommit indicates an expected call of GetSeenCommit
func (mr *MockNeedChainReaderFunctionMockRecorder) GetSeenCommit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeenCommit", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).GetSeenCommit), arg0)
}

// GetSlot mocks base method
func (m *MockNeedChainReaderFunction) GetSlot(arg0 model.AbstractBlock) *uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlot", arg0)
	ret0, _ := ret[0].(*uint64)
	return ret0
}

// GetSlot indicates an expected call of GetSlot
func (mr *MockNeedChainReaderFunctionMockRecorder) GetSlot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlot", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).GetSlot), arg0)
}

// GetVerifiers mocks base method
func (m *MockNeedChainReaderFunction) GetVerifiers(arg0 uint64) []common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVerifiers", arg0)
	ret0, _ := ret[0].([]common.Address)
	return ret0
}

// GetVerifiers indicates an expected call of GetVerifiers
func (mr *MockNeedChainReaderFunctionMockRecorder) GetVerifiers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVerifiers", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).GetVerifiers), arg0)
}

// IsChangePoint mocks base method
func (m *MockNeedChainReaderFunction) IsChangePoint(arg0 model.AbstractBlock, arg1 bool) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsChangePoint", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsChangePoint indicates an expected call of IsChangePoint
func (mr *MockNeedChainReaderFunctionMockRecorder) IsChangePoint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsChangePoint", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).IsChangePoint), arg0, arg1)
}

// SaveBlock mocks base method
func (m *MockNeedChainReaderFunction) SaveBlock(arg0 model.AbstractBlock, arg1 []model.AbstractVerification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBlock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBlock indicates an expected call of SaveBlock
func (mr *MockNeedChainReaderFunctionMockRecorder) SaveBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBlock", reflect.TypeOf((*MockNeedChainReaderFunction)(nil).SaveBlock), arg0, arg1)
}