// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dipperin/dipperin-core/core/economymodel (interfaces: EconomyNeedService)

// Package economymodel is a generated GoMock package.
package economymodel

import (
	common "github.com/dipperin/dipperin-core/common"
	model "github.com/dipperin/dipperin-core/core/model"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEconomyNeedService is a mock of EconomyNeedService interface
type MockEconomyNeedService struct {
	ctrl     *gomock.Controller
	recorder *MockEconomyNeedServiceMockRecorder
}

// MockEconomyNeedServiceMockRecorder is the mock recorder for MockEconomyNeedService
type MockEconomyNeedServiceMockRecorder struct {
	mock *MockEconomyNeedService
}

// NewMockEconomyNeedService creates a new mock instance
func NewMockEconomyNeedService(ctrl *gomock.Controller) *MockEconomyNeedService {
	mock := &MockEconomyNeedService{ctrl: ctrl}
	mock.recorder = &MockEconomyNeedServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEconomyNeedService) EXPECT() *MockEconomyNeedServiceMockRecorder {
	return m.recorder
}

// GetSlot mocks base method
func (m *MockEconomyNeedService) GetSlot(arg0 model.AbstractBlock) *uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlot", arg0)
	ret0, _ := ret[0].(*uint64)
	return ret0
}

// GetSlot indicates an expected call of GetSlot
func (mr *MockEconomyNeedServiceMockRecorder) GetSlot(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlot", reflect.TypeOf((*MockEconomyNeedService)(nil).GetSlot), arg0)
}

// GetVerifiers mocks base method
func (m *MockEconomyNeedService) GetVerifiers(arg0 uint64) []common.Address {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVerifiers", arg0)
	ret0, _ := ret[0].([]common.Address)
	return ret0
}

// GetVerifiers indicates an expected call of GetVerifiers
func (mr *MockEconomyNeedServiceMockRecorder) GetVerifiers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVerifiers", reflect.TypeOf((*MockEconomyNeedService)(nil).GetVerifiers), arg0)
}