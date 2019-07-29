package resolver

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"math/big"
)

type ContractRef interface {
	Address() common.Address
}

//go:generate mockgen -destination=/home/qydev/go/src/github.com/dipperin/dipperin-core/core/vm/resolver/vm_context_service_mock_test.go -package=resolver github.com/dipperin/dipperin-core/core/vm/resolver VmContextService
type VmContextService interface {
	GetGasPrice() *big.Int
	GetGasLimit() uint64
	GetBlockHash(num uint64) common.Hash
	GetBlockNumber() *big.Int
	GetTime() *big.Int
	GetCoinBase() common.Address
	GetOrigin() common.Address
	Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error)
	//GetCallGasTemp() uint64
	DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error)
	GetTxHash() common.Hash
}

//go:generate mockgen -destination=/home/qydev/go/src/github.com/dipperin/dipperin-core/core/vm/resolver/contract_service_mock_test.go -package=resolver github.com/dipperin/dipperin-core/core/vm/resolver ContractService
type ContractService interface {
	Caller() common.Address
	Address() common.Address
	CallValue() *big.Int
	GetGas() uint64
}

//go:generate mockgen -destination=/home/qydev/go/src/github.com/dipperin/dipperin-core/core/vm/resolver/statedb_service_mock_test.go -package=resolver github.com/dipperin/dipperin-core/core/vm/resolver StateDBService
type StateDBService interface {
	GetBalance(addr common.Address) *big.Int
	AddLog(addedLog *model.Log)
	SetState(addr common.Address, key []byte, value []byte)
	GetState(addr common.Address, key []byte) (data []byte)
	GetNonce(common.Address) (uint64, error)
}

type resolverNeedExternalService struct {
	ContractService
	VmContextService
	StateDBService
}
