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
	"github.com/dipperin/dipperin-core/third-party/log"
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
	Caller() ContractRef
	Self() ContractRef
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

func (service *resolverNeedExternalService) Transfer(toAddr common.Address, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	//todo gas used
	gas := uint64(0)
	/*gas := self.evm.callGasTemp
	if value.Sign() != 0 {
		gas += params.CallStipend
	}*/
	log.Info("Service#Transfer", "from", service.Self().Address(), "to", toAddr, "value", value, "gasLimit", gas)
	ret, returnGas, err := service.Call(service.Self(), toAddr, nil, gas, value)
	return ret, returnGas, err
}
