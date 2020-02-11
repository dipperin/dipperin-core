// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/csbft/statemachine (interfaces: ChainReader,MsgSigner,MsgSender,Validator,Fetcher)

// Package components is a generated GoMock package.
package components

import (
	common "github.com/dipperin/dipperin-core/common"
	model "github.com/dipperin/dipperin-core/core/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockChainReader is a mock of ChainReader interface
type MockChainReader struct {
	ctrl     *gomock.Controller
	recorder *MockChainReaderMockRecorder
}

// MockChainReaderMockRecorder is the mock recorder for MockChainReader
type MockChainReaderMockRecorder struct {
	mock *MockChainReader
}

// NewMockChainReader creates a new mock instance
func NewMockChainReader(ctrl *gomock.Controller) *MockChainReader {
	mock := &MockChainReader{ctrl: ctrl}
	mock.recorder = &MockChainReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChainReader) EXPECT() *MockChainReaderMockRecorder {
	return m.recorder
}

// CurrentBlock mocks base method
func (m *MockChainReader) CurrentBlock() model.AbstractBlock {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CurrentBlock")
	ret0, _ := ret[0].(model.AbstractBlock)
	return ret0
}

// CurrentBlock indicates an expected call of CurrentBlock
func (mr *MockChainReaderMockRecorder) CurrentBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CurrentBlock", reflect.TypeOf((*MockChainReader)(nil).CurrentBlock))
}

// GetCurrVerifiers mocks base method
func (m *MockChainReader) GetCurrVerifiers() []common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrVerifiers")
	ret0, _ := ret[0].([]common.Address)
	return ret0
}

// GetCurrVerifiers indicates an expected call of GetCurrVerifiers
func (mr *MockChainReaderMockRecorder) GetCurrVerifiers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrVerifiers", reflect.TypeOf((*MockChainReader)(nil).GetCurrVerifiers))
}

// GetNextVerifiers mocks base method
func (m *MockChainReader) GetNextVerifiers() []common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextVerifiers")
	ret0, _ := ret[0].([]common.Address)
	return ret0
}

// GetNextVerifiers indicates an expected call of GetNextVerifiers
func (mr *MockChainReaderMockRecorder) GetNextVerifiers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextVerifiers", reflect.TypeOf((*MockChainReader)(nil).GetNextVerifiers))
}

// GetSeenCommit mocks base method
func (m *MockChainReader) GetSeenCommit(arg0 uint64) []model.AbstractVerification {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSeenCommit", arg0)
	ret0, _ := ret[0].([]model.AbstractVerification)
	return ret0
}

// GetSeenCommit indicates an expected call of GetSeenCommit
func (mr *MockChainReaderMockRecorder) GetSeenCommit(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSeenCommit", reflect.TypeOf((*MockChainReader)(nil).GetSeenCommit), arg0)
}

// IsChangePoint mocks base method
func (m *MockChainReader) IsChangePoint(arg0 model.AbstractBlock, arg1 bool) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsChangePoint", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsChangePoint indicates an expected call of IsChangePoint
func (mr *MockChainReaderMockRecorder) IsChangePoint(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsChangePoint", reflect.TypeOf((*MockChainReader)(nil).IsChangePoint), arg0, arg1)
}

// SaveBlock mocks base method
func (m *MockChainReader) SaveBlock(arg0 model.AbstractBlock, arg1 []model.AbstractVerification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBlock", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBlock indicates an expected call of SaveBlock
func (mr *MockChainReaderMockRecorder) SaveBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBlock", reflect.TypeOf((*MockChainReader)(nil).SaveBlock), arg0, arg1)
}

// MockMsgSigner is a mock of MsgSigner interface
type MockMsgSigner struct {
	ctrl     *gomock.Controller
	recorder *MockMsgSignerMockRecorder
}

// MockMsgSignerMockRecorder is the mock recorder for MockMsgSigner
type MockMsgSignerMockRecorder struct {
	mock *MockMsgSigner
}

// NewMockMsgSigner creates a new mock instance
func NewMockMsgSigner(ctrl *gomock.Controller) *MockMsgSigner {
	mock := &MockMsgSigner{ctrl: ctrl}
	mock.recorder = &MockMsgSignerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMsgSigner) EXPECT() *MockMsgSignerMockRecorder {
	return m.recorder
}

// GetAddress mocks base method
func (m *MockMsgSigner) GetAddress() common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAddress")
	ret0, _ := ret[0].(common.Address)
	return ret0
}

// GetAddress indicates an expected call of GetAddress
func (mr *MockMsgSignerMockRecorder) GetAddress() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAddress", reflect.TypeOf((*MockMsgSigner)(nil).GetAddress))
}

// SignHash mocks base method
func (m *MockMsgSigner) SignHash(arg0 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignHash", arg0)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignHash indicates an expected call of SignHash
func (mr *MockMsgSignerMockRecorder) SignHash(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignHash", reflect.TypeOf((*MockMsgSigner)(nil).SignHash), arg0)
}

// MockMsgSender is a mock of MsgSender interface
type MockMsgSender struct {
	ctrl     *gomock.Controller
	recorder *MockMsgSenderMockRecorder
}

// MockMsgSenderMockRecorder is the mock recorder for MockMsgSender
type MockMsgSenderMockRecorder struct {
	mock *MockMsgSender
}

// NewMockMsgSender creates a new mock instance
func NewMockMsgSender(ctrl *gomock.Controller) *MockMsgSender {
	mock := &MockMsgSender{ctrl: ctrl}
	mock.recorder = &MockMsgSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMsgSender) EXPECT() *MockMsgSenderMockRecorder {
	return m.recorder
}

// BroadcastEiBlock mocks base method
func (m *MockMsgSender) BroadcastEiBlock(arg0 model.AbstractBlock) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BroadcastEiBlock", arg0)
}

// BroadcastEiBlock indicates an expected call of BroadcastEiBlock
func (mr *MockMsgSenderMockRecorder) BroadcastEiBlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastEiBlock", reflect.TypeOf((*MockMsgSender)(nil).BroadcastEiBlock), arg0)
}

// BroadcastMsg mocks base method
func (m *MockMsgSender) BroadcastMsg(arg0 uint64, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BroadcastMsg", arg0, arg1)
}

// BroadcastMsg indicates an expected call of BroadcastMsg
func (mr *MockMsgSenderMockRecorder) BroadcastMsg(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastMsg", reflect.TypeOf((*MockMsgSender)(nil).BroadcastMsg), arg0, arg1)
}

// SendReqRoundMsg mocks base method
func (m *MockMsgSender) SendReqRoundMsg(arg0 uint64, arg1 []common.Address, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendReqRoundMsg", arg0, arg1, arg2)
}

// SendReqRoundMsg indicates an expected call of SendReqRoundMsg
func (mr *MockMsgSenderMockRecorder) SendReqRoundMsg(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendReqRoundMsg", reflect.TypeOf((*MockMsgSender)(nil).SendReqRoundMsg), arg0, arg1, arg2)
}

// MockValidator is a mock of Validator interface
type MockValidator struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorMockRecorder
}

// MockValidatorMockRecorder is the mock recorder for MockValidator
type MockValidatorMockRecorder struct {
	mock *MockValidator
}

// NewMockValidator creates a new mock instance
func NewMockValidator(ctrl *gomock.Controller) *MockValidator {
	mock := &MockValidator{ctrl: ctrl}
	mock.recorder = &MockValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockValidator) EXPECT() *MockValidatorMockRecorder {
	return m.recorder
}

// FullValid mocks base method
func (m *MockValidator) FullValid(arg0 model.AbstractBlock) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullValid", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// FullValid indicates an expected call of FullValid
func (mr *MockValidatorMockRecorder) FullValid(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullValid", reflect.TypeOf((*MockValidator)(nil).FullValid), arg0)
}

// MockFetcher is a mock of Fetcher interface
type MockFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockFetcherMockRecorder
}

// MockFetcherMockRecorder is the mock recorder for MockFetcher
type MockFetcherMockRecorder struct {
	mock *MockFetcher
}

// NewMockFetcher creates a new mock instance
func NewMockFetcher(ctrl *gomock.Controller) *MockFetcher {
	mock := &MockFetcher{ctrl: ctrl}
	mock.recorder = &MockFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFetcher) EXPECT() *MockFetcherMockRecorder {
	return m.recorder
}

// FetchBlock mocks base method
func (m *MockFetcher) FetchBlock(arg0 common.Address, arg1 common.Hash) model.AbstractBlock {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchBlock", arg0, arg1)
	ret0, _ := ret[0].(model.AbstractBlock)
	return ret0
}

// FetchBlock indicates an expected call of FetchBlock
func (mr *MockFetcherMockRecorder) FetchBlock(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchBlock", reflect.TypeOf((*MockFetcher)(nil).FetchBlock), arg0, arg1)
}
