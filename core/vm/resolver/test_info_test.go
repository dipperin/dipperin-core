// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package resolver

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/third-party/life/mem-manage"
	"math/big"
)

var (
	aliceAddr    = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	bobAddr      = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
	contractAddr = common.HexToAddress("0x0014B5Df12F50295469Fe33951403b8f4E63231Ef488")
)

var TEST_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: mem_manage.DefaultPageSize,
}

type fakeContractRef struct {
	address common.Address
}

func (ref *fakeContractRef) Address() common.Address {
	return ref.address
}

type fakeStateDBService struct {
	logList  []*model.Log
	stateMap map[common.Address]map[string][]byte
}

func NewFakeStateDBService() *fakeStateDBService {
	return &fakeStateDBService{
		logList:  make([]*model.Log, 0),
		stateMap: make(map[common.Address]map[string][]byte, 0),
	}
}

func (state *fakeStateDBService) GetBalance(addr common.Address) *big.Int {
	return big.NewInt(10000)
}

func (state *fakeStateDBService) AddLog(addedLog *model.Log) {
	state.logList = append(state.logList, addedLog)
}

func (state *fakeStateDBService) SetState(addr common.Address, key []byte, value []byte) {
	if state.stateMap[addr] == nil {
		state.stateMap[addr] = make(map[string][]byte, 0)
	}
	state.stateMap[addr][string(key)] = value
}

func (state *fakeStateDBService) GetState(addr common.Address, key []byte) (data []byte) {
	return state.stateMap[addr][string(key)]
}

func (state *fakeStateDBService) GetNonce(common.Address) (uint64, error) {
	return uint64(0), nil
}

type fakeContractService struct {
}

func (contract *fakeContractService) Caller() ContractRef {
	return &fakeContractRef{aliceAddr}
}

func (contract *fakeContractService) Self() ContractRef {
	return &fakeContractRef{contractAddr}
}

func (contract *fakeContractService) Address() common.Address {
	return contractAddr
}

func (contract *fakeContractService) CallValue() *big.Int {
	return g_testData.TestValue
}

func (contract *fakeContractService) GetGas() uint64 {
	return g_testData.TestGasLimit
}

type fakeVmContextService struct {
}

func (context *fakeVmContextService) GetGasPrice() *big.Int {
	return g_testData.TestGasPrice
}

func (context *fakeVmContextService) GetGasLimit() uint64 {
	return g_testData.TestGasLimit
}

func (context *fakeVmContextService) GetBlockHash(num uint64) common.Hash {
	return common.HexToHash("blockHash")
}

func (context *fakeVmContextService) GetBlockNumber() *big.Int {
	return big.NewInt(1)
}

func (context *fakeVmContextService) GetTime() *big.Int {
	return big.NewInt(10)
}

func (context *fakeVmContextService) GetCoinBase() common.Address {
	return aliceAddr
}

func (context *fakeVmContextService) GetOrigin() common.Address {
	return bobAddr
}

func (context *fakeVmContextService) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	return nil, gas, nil
}

func (context *fakeVmContextService) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	return nil, gas, nil
}

func (context *fakeVmContextService) GetTxHash() common.Hash {
	return common.HexToHash("txHash")
}
