package resolver

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"math/big"
)

var (
	aliceAddr    = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	bobAddr      = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
	contractAddr = common.HexToAddress("0x0014B5Df12F50295469Fe33951403b8f4E63231Ef488")
)

var TEST_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: exec.DefaultPageSize,
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

func (contract *fakeContractService) Caller() common.Address {
	return bobAddr
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
	panic("implement me")
}

func (context *fakeVmContextService) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	panic("implement me")
}

func (context *fakeVmContextService) GetTxHash() common.Hash {
	return common.HexToHash("txHash")
}
