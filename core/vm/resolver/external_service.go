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
	"encoding/hex"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	model2 "github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/base/utils"
	"go.uber.org/zap"
	"math/big"
	"strings"
)

type ContractRef interface {
	Address() common.Address
}

//go:generate mockgen -destination=./vm_context_service_mock.go -package=resolver github.com/dipperin/dipperin-core/core/vm/resolver VmContextService
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
	TransferValue(caller ContractRef,  toAddr common.Address, value *big.Int) error
}

//go:generate mockgen -destination=./contract_service_mock.go -package=resolver github.com/dipperin/dipperin-core/core/vm/resolver ContractService
type ContractService interface {
	Caller() ContractRef
	Self() ContractRef
	Address() common.Address
	CallValue() *big.Int
	GetGas() uint64
}

//go:generate mockgen -destination=./statedb_service_mock.go -package=resolver github.com/dipperin/dipperin-core/core/vm/resolver StateDBService
type StateDBService interface {
	GetBalance(addr common.Address) *big.Int
	AddLog(addedLog *model2.Log)
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
	log.DLogger.Info("Service#Transfer", zap.Any("from", service.Self().Address()), zap.Any("to", toAddr), zap.Any("value", value), zap.Any("gasLimit", gas))
	err = service.TransferValue(service.Self(), toAddr, value)
	//ret, returnGas, err := service.Call(service.Self(), toAddr, nil, gas, value)
	return []byte{}, gas, err
}

func (service *resolverNeedExternalService) ResolverCall(addr, param []byte) ([]byte, error) {
	funcName, err := utils.ParseInputForFuncName(param)
	if err != nil {
		log.DLogger.Error("ResolverCall#ParseInputForFuncName failed", zap.Error(err))
		return nil, err
	}

	// check funcName
	if strings.EqualFold(funcName, "init") {
		log.DLogger.Error("ResolverCall can't call init function")
		return nil, gerror.ErrFunctionInitCanNotCalled
	}

	contractAddr := common.HexToAddress(hex.EncodeToString(addr))
	log.DLogger.Info("Call ResolverCall", zap.Any("caller", service.Self().Address()), zap.Any("contractAddr", contractAddr), zap.Uint64("gas", service.GetGas()), zap.Any("value", service.CallValue()), zap.Uint8s("inputs", param))
	ret, _, err := service.Call(service.ContractService, contractAddr, param, service.GetGas(), service.CallValue())
	return ret, err
}

func (service *resolverNeedExternalService) ResolverDelegateCall(addr, param []byte) ([]byte, error) {
	funcName, err := utils.ParseInputForFuncName(param)
	if err != nil {
		log.DLogger.Error("ResolverDelegateCall#ParseInputForFuncName failed", zap.Error(err))
		return nil, err
	}

	// check funcName
	if strings.EqualFold(funcName, "init") {
		log.DLogger.Error("ResolverDelegateCall can't call init function")
		return nil, gerror.ErrFunctionInitCanNotCalled
	}

	contractAddr := common.HexToAddress(hex.EncodeToString(addr))
	log.DLogger.Info("Call ResolverDelegateCall", zap.Any("caller", service.Self().Address()), zap.Any("contractAddr", contractAddr), zap.Uint64("gas", service.GetGas()), zap.Uint8s("inputs", param))
	ret, _, err := service.DelegateCall(service.ContractService, contractAddr, param, service.GetGas())
	return ret, err
}
