package resolver

import (
	"encoding/hex"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/common/params"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"math/big"
)

type ContractRef interface {
	Address() common.Address
}

type VmContextService interface {
	GetGasPrice() int64
	GetGasLimit() uint64
	BlockHash(num uint64) common.Hash
	GetBlockNumber() *big.Int
	GetTime() *big.Int
	GetCoinBase() common.Address
	GetOrigin() common.Address
	Call(caller ContractRef, addr common.Address, input []byte,gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error)
	GetCallGasTemp() uint64
	DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error)
}

type ContractService interface {
	Caller() common.Address
	Address() common.Address
	CallValue() *big.Int
	GetGas() uint64
}

type StateDBService interface {
	GetBalance(addr common.Address) *big.Int
	AddLog(addedLog *model.Log)
	//AddLog(address common.Address, topics []common.Hash, data []byte, bn uint64)
	SetState(addr common.Address,key []byte, value []byte)
	GetState(addr common.Address,key []byte) (data []byte)
	GetNonce(common.Address) uint64
	TxHash() common.Hash
	TxIdx() uint32
}

type resolverNeedExternalService struct {
	ContractService
	VmContextService
	StateDBService
}

func (service *resolverNeedExternalService) GetCallerNonce() int64 {
	addr := service.Caller()

	return int64(service.StateDBService.GetNonce(addr))
}

func (service *resolverNeedExternalService) ReSolverSetState(key []byte, value []byte)  {
	service.StateDBService.SetState(service.Address(), key, value)
}

func (service *resolverNeedExternalService) ReSolverGetState(key []byte) []byte {
	return service.StateDBService.GetState(service.Address(), key)
}

func (service *resolverNeedExternalService) Transfer(toAddr common.Address, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	gas := service.GetCallGasTemp()
	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	fmt.Println("Transfer to:", toAddr.String())
	fmt.Println("Transfer caller:", service.Address().Hex())
	ret, returnGas, err := service.VmContextService.Call(service.ContractService, toAddr, nil, gas, value)
	return ret, returnGas, err
}

func (service *resolverNeedExternalService) ResolverCall(addr, param []byte) ([]byte, error) {

	ret, _, err := service.VmContextService.Call(service.ContractService, common.HexToAddress(hex.EncodeToString(addr)), param, service.ContractService.GetGas(), service.ContractService.CallValue())
	return ret, err
}

func (service *resolverNeedExternalService) ResolverDelegateCall(addr, param []byte) ([]byte, error) {

	ret, _, err := service.VmContextService.DelegateCall(service.ContractService, common.HexToAddress(hex.EncodeToString(addr)), param,service.ContractService.GetGas())
	return ret, err
}

