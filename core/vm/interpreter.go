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

package vm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/core/vm/common"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/core/vm/resolver"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"math/big"
	"reflect"
	"runtime"
	"strings"
)

var (
	errEmptyInput         = errors.New("interpreter_life: empty input")
	errEmptyABI           = errors.New("interpreter_life: empty abi")
	errInvalidRlpFormat   = errors.New("interpreter_life: invalid rlp format")
	errInsufficientParams = errors.New("interpreter_life: invalid input params")
	errInvalidAbi         = errors.New("interpreter_life: invalid abi, from json fail")
	errInputAbiNotMatch   = errors.New("interpreter_life: length of input and abi not match")
	errInvalidReturnType  = errors.New("interpreter_life: return type not void")
	errFuncNameNotFound   = errors.New("interpreter_life: function name not found")
)

const (
	CALL_CONTRACT_FLAG = 9
)

type Interpreter interface {
	// Run loops and evaluates the contract's code with the given input data and returns
	// the return byte-slice and an error if one occurred.
	Run(vm *VM, contract *Contract, create bool) ([]byte, error)
	// CanRun tells if the contract, passed as an argument, can be
	CanRun([]byte) bool
}

type WASMInterpreter struct {
	state    common.StateDB
	context  *common.Context
	config   exec.VMConfig
	resolver exec.ImportResolver
}

// NewWASMInterpreter returns a new instance of the Interpreter
func NewWASMInterpreter(state common.StateDB, context common.Context, vmConfig exec.VMConfig) *WASMInterpreter {
	return &WASMInterpreter{
		state,
		&context,
		vmConfig,
		&resolver.Resolver{},
	}
}

func (in *WASMInterpreter) Run(vm *VM, contract *Contract, create bool) (ret []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			//fmt.Println(stack())
			stackInfo := stack()
			log.DLogger.Info("WASMInterpreter panic err", zap.String("stackInfo", stackInfo))
			ret, err = nil, fmt.Errorf("VM execute fail: %v", er)
		}
	}()
	vm.depth++
	defer func() {
		vm.depth--
		if vm.depth == 0 {
			log.DLogger.Info("VM depth = 0")
		}
	}()

	if len(contract.Code) == 0 || len(contract.ABI) == 0 {
		log.DLogger.Debug("code or ABI Length is 0", zap.Int("code", len(contract.Code)), zap.Int("abi", len(contract.ABI)))
		return nil, nil
	}

	//　life方法注入新建虚拟机
	solver := resolver.NewResolver(vm, contract, in.state)
	lifeVm, err := exec.NewVirtualMachine(contract.Code, in.config, solver, nil)
	if err != nil {
		log.DLogger.Info("NewVirtualMachine failed", zap.Error(err))
		return nil, err
	}
	lifeVm.GasLimit = contract.Gas
	defer func() {
		lifeVm.Stop()
	}()

	var (
		funcName   string
		params     []int64
		returnType string
	)

	if create {
		// init function.
		funcName = "init"
		params, returnType, err = ParseInitFunctionByABI(lifeVm, contract.Input, contract.ABI)
		if err != nil {
			log.DLogger.Error("ParseInitFunctionByABI failed", zap.Error(err))
			return nil, err
		}
	} else {
		// parse input
		funcName, params, returnType, err = ParseCallExtraDataByABI(lifeVm, contract.Input, contract.ABI)
		if err != nil {
			if err == errInsufficientParams { // transfer to contract address.
				return nil, nil
			}
			log.DLogger.Error("ParseCallExtraDataByABI failed", zap.Error(err))
			return nil, err
		}
	}
	log.DLogger.Info("WASMInterpreter Run", zap.String("funcName", funcName), zap.Int64s("params", params), zap.String("return", returnType), zap.Error(err))

	//　获取entryID
	entryID, ok := lifeVm.GetFunctionExport(funcName)
	if !ok {
		return nil, errFuncNameNotFound
	}

	res, err := lifeVm.Run(entryID, params...)
	log.DLogger.Info("Run lifeVm", zap.Uint64("gasUsed", lifeVm.GasUsed), zap.Uint64("gasLimit", lifeVm.GasLimit))
	if err != nil {
		log.DLogger.Error("throw exception:", zap.Error(err))
		return nil, err
	}

	contract.Gas = contract.Gas - lifeVm.GasUsed
	if create {
		return contract.Code, nil
	}

	switch returnType {
	case "void", "bool", "int8", "int16", "int32", "int64":
		bigRes := new(big.Int)
		bigRes.SetInt64(res)
		finalRes := utils.Align32Bytes(math.U256(bigRes).Bytes())
		return finalRes, nil
	case "uint8", "uint16", "uint32", "uint64":
		finalRes := utils.Align32Bytes(utils.Uint64ToBytes(uint64(res)))
		return finalRes, nil
	case "string":
		returnBytes := make([]byte, 0)
		copyData := lifeVm.Memory.Memory[res:]
		for _, v := range copyData {
			if v == 0 {
				break
			}
			returnBytes = append(returnBytes, v)
		}
		var dataRealSize = len(returnBytes)
		if (dataRealSize % 32) != 0 {
			dataRealSize = dataRealSize + (32 - (dataRealSize % 32))
		}
		dataByt := make([]byte, dataRealSize)
		copy(dataByt[:], returnBytes)
		return dataByt, nil
	}
	return nil, nil
}

// CanRun tells if the contract, passed as an argument, can be run
// by the current interpreter
func (in *WASMInterpreter) CanRun(code []byte) bool {
	return true
}

// input = RLP([params])
// returnType must void
func ParseInitFunctionByABI(vm *exec.VirtualMachine, input []byte, abi []byte) (params []int64, returnType string, err error) {
	if input == nil || len(input) <= 1 {
		log.DLogger.Info("InitFunc has no input")
		return
	}

	if abi == nil || len(abi) == 0 {
		err = errEmptyABI
		return
	}

	// rlp decode
	ptr := new(interface{})
	rlp.Decode(bytes.NewReader(input), &ptr)
	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		err = errInvalidRlpFormat
		return
	}

	params, returnType, err = findParams(vm, abi, "init", rlpList.([]interface{}))
	if err != nil {
		return
	}

	if returnType != "void" {
		err = errInvalidReturnType
		return
	}
	return
}

// input = RLP([funcName][params])
// get returnType[0] if more than 1 return
func ParseCallExtraDataByABI(vm *exec.VirtualMachine, input []byte, abi []byte) (funcName string, params []int64, returnType string, err error) {
	if input == nil || len(input) == 0 {
		err = errEmptyInput
		return
	}

	if abi == nil || len(abi) == 0 {
		err = errEmptyABI
		return
	}

	// rlp decode
	ptr := new(interface{})
	rlp.Decode(bytes.NewReader(input), &ptr)
	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		err = errInvalidRlpFormat
		return
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) < 1 {
		err = errInsufficientParams
		return
	}

	if v, ok := iRlpList[0].([]byte); ok {
		funcName = string(v)
	}

	var inputList []interface{}
	for _, value := range iRlpList[1:] {
		inputList = append(inputList, value)
	}
	params, returnType, err = findParams(vm, abi, funcName, inputList)
	return
}

func findParams(vm *exec.VirtualMachine, abi []byte, funcName string, inputList []interface{}) (params []int64, returnType string, err error) {
	wasmAbi := new(utils.WasmAbi)
	err = wasmAbi.FromJson(abi)
	if err != nil {
		log.DLogger.Error("findParams#FromJson failed", zap.Error(err))
		err = errInvalidAbi
		return
	}

	var abiParam []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold(v.Name, funcName) && strings.EqualFold(v.Type, "function") {
			abiParam = v.Inputs
			log.DLogger.Info("findParams", zap.Int("len outputs", len(v.Outputs)), zap.Int("abiParam", len(abiParam)), zap.Int("inputList", len(inputList)))
			if len(v.Outputs) != 0 {
				returnType = v.Outputs[0].Type
			} else {
				returnType = "void"
			}
			break
		}
	}

	if len(abiParam) != len(inputList) {
		log.DLogger.Error("findParams failed", zap.Error(errInputAbiNotMatch), zap.Int("abiLen", len(abiParam)), zap.Int("inputLen", len(inputList)))
		err = errInputAbiNotMatch
		return
	}

	// uint64 uint32  uint16 uint8 int64 int32  int16 int8 string void bool
	for i, v := range abiParam {
		input := inputList[i].([]byte)
		switch v.Type {
		case "string":
			pos := resolver.MallocString(vm, string(input))
			params = append(params, pos)
		case "int8", "uint8":
			params = append(params, int64(input[0]))
		case "int16", "uint16":
			params = append(params, int64(binary.BigEndian.Uint16(input)))
		case "int32", "int", "uint32", "uint":
			params = append(params, int64(binary.BigEndian.Uint32(input)))
		case "int64", "uint64":
			params = append(params, int64(binary.BigEndian.Uint64(input)))
		case "bool":
			params = append(params, int64(input[0]))
		}
	}
	return
}

// rlpData=RLP([code][abi][init params])
func ParseCreateExtraData(rlpData []byte) (code, abi, rlpInit []byte, err error) {
	if rlpData == nil || len(rlpData) == 0 {
		err = errEmptyInput
		return
	}

	ptr := new(interface{})
	rlp.Decode(bytes.NewReader(rlpData), &ptr)
	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		err = errInvalidRlpFormat
		return
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) < 2 {
		err = errInsufficientParams
		return
	}

	// parse code and abi
	if v, ok := iRlpList[0].([]byte); ok {
		code = v
	}

	if v, ok := iRlpList[1].([]byte); ok {
		abi = v
	}

	if len(iRlpList) == 2 {
		return
	}

	// encode init
	var init []interface{}
	for _, value := range iRlpList[2:] {
		init = append(init, value)
	}

	rlpInit, err = rlp.EncodeToBytes(init)
	return
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}
