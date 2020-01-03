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
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestWASMInterpreter_Run_DIPCLibContract(t *testing.T) {
	t.Skip()
	testVm := getTestVm()
	interpreter := testVm.Interpreter
	inputs := genInput(t, g_testData.ContractTestPar.CallFuncName, [][]byte{})
	log.DLogger.Info("the wasmPath is:", zap.String("wasmPath", g_testData.ContractTestPar.WASMPath))
	log.DLogger.Info("the abiPath is:", zap.String("abiPath", g_testData.ContractTestPar.AbiPath))
	contract := getContract(g_testData.ContractTestPar.WASMPath, g_testData.ContractTestPar.AbiPath, inputs)
	_, err := interpreter.Run(testVm, contract, false)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	testVm := getTestVm()
	interpreter := testVm.Interpreter

	// create contract
	code, _ := g_testData.GetCodeAbi(WASMPath, AbiPath)
	contract := getContract(WASMPath, AbiPath, nil)
	result, err := interpreter.Run(testVm, contract, true)
	assert.Equal(t, code, result)
	assert.NoError(t, err)

	// call contract
	name := []byte("contract")
	inputs := genInput(t, "returnString", [][]byte{name})
	contract.Input = inputs
	result, err = interpreter.Run(testVm, contract, false)
	resp := utils.Align32BytesConverter(result, "string")
	assert.Equal(t, "contract", resp)
	assert.NoError(t, err)

	inputs = genInput(t, "returnInt", [][]byte{name})
	contract.Input = inputs
	result, err = interpreter.Run(testVm, contract, false)
	resp = utils.Align32BytesConverter(result, "int64")
	assert.Equal(t, int64(50), resp)
	assert.NoError(t, err)

	inputs = genInput(t, "returnUint", [][]byte{name})
	contract.Input = inputs
	result, err = interpreter.Run(testVm, contract, false)
	resp = utils.Align32BytesConverter(result, "uint64")
	assert.Equal(t, uint64(50), resp)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_Error(t *testing.T) {
	testVm := getTestVm()
	interpreter := testVm.Interpreter
	assert.Equal(t, true, interpreter.CanRun(nil))

	contract := &Contract{}
	result, err := interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.NoError(t, err)

	contract.Code = []byte{123}
	contract.ABI = []byte{123}
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Equal(t, "unexpected EOF", err.Error())

	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)
	contract.Code = code
	contract.ABI = abi
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Equal(t, errEmptyInput, err)

	input, err := rlp.EncodeToBytes([]interface{}{})
	assert.NoError(t, err)
	contract.Input = input
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.NoError(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{"func"})
	assert.NoError(t, err)
	contract.Input = input
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Equal(t, errFuncNameNotFound, err)

	input, err = rlp.EncodeToBytes([]interface{}{"returnString"})
	assert.NoError(t, err)
	contract.Input = input
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Equal(t, errInputAbiNotMatch, err)

	name := []byte("string")
	param := [][]byte{name}
	input = genInput(t, "returnString", param)
	contract.Input = input
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Equal(t, "out of gas  cost:1 GasUsed:0 GasLimit:0", err.Error())

	contract.Input = []byte{1, 2, 3}
	result, err = interpreter.Run(testVm, contract, true)
	assert.Nil(t, result)
	assert.Equal(t, errInvalidRlpFormat, err)
}

func TestParseInitFunctionByABI(t *testing.T) {
	lifeVm := &exec.VirtualMachine{}
	num := utils.Uint64ToBytes(100)
	input := genInput(t, "", [][]byte{num})

	_, _, err := ParseInitFunctionByABI(lifeVm, nil, nil)
	assert.NoError(t, err)

	_, _, err = ParseInitFunctionByABI(lifeVm, input, nil)
	assert.Equal(t, errEmptyABI, err)

	_, _, err = ParseInitFunctionByABI(lifeVm, []byte{1, 2, 3}, []byte{1, 2, 3})
	assert.Equal(t, errInvalidRlpFormat, err)

	_, _, err = ParseInitFunctionByABI(lifeVm, input, []byte(abi1))
	assert.Equal(t, errInputAbiNotMatch, err)

	_, _, err = ParseInitFunctionByABI(lifeVm, input, []byte(abi2))
	assert.Equal(t, errInvalidReturnType, err)

	_, _, err = ParseInitFunctionByABI(lifeVm, input, []byte(abi3))
	assert.NoError(t, err)
}

func TestParseCallExtraDataByABI(t *testing.T) {
	lifeVm := &exec.VirtualMachine{}
	_, _, _, err := ParseCallExtraDataByABI(lifeVm, nil, nil)
	assert.Equal(t, errEmptyInput, err)

	_, _, _, err = ParseCallExtraDataByABI(lifeVm, []byte{1, 2, 3}, nil)
	assert.Equal(t, errEmptyABI, err)

	_, _, _, err = ParseCallExtraDataByABI(lifeVm, []byte{1, 2, 3}, []byte{1, 2, 3})
	assert.Equal(t, errInvalidRlpFormat, err)

	input := genInput(t, "", [][]byte{})
	_, _, _, err = ParseCallExtraDataByABI(lifeVm, input, []byte(abi2))
	assert.Equal(t, errInsufficientParams, err)

	input = genInput(t, "test", [][]byte{})
	_, _, _, err = ParseCallExtraDataByABI(lifeVm, input, []byte(abi2))
	assert.NoError(t, err)

	_, _, _, err = ParseCallExtraDataByABI(lifeVm, input, []byte{123})
	assert.Equal(t, errInvalidAbi, err)

	input = genInput(t, "init", [][]byte{})
	_, _, _, err = ParseCallExtraDataByABI(lifeVm, input, []byte(abi2))
	assert.Equal(t, errInputAbiNotMatch, err)

	num := utils.Uint64ToBytes(100)
	input = genInput(t, "init", [][]byte{num})
	funcName, params, returnType, err := ParseCallExtraDataByABI(lifeVm, input, []byte(abi2))
	assert.Equal(t, "init", funcName)
	assert.Equal(t, 1, len(params))
	assert.Equal(t, "string", returnType)
	assert.NoError(t, err)

	funcName, params, returnType, err = ParseCallExtraDataByABI(lifeVm, input, []byte(abi3))
	assert.Equal(t, "init", funcName)
	assert.Equal(t, 1, len(params))
	assert.Equal(t, "void", returnType)
	assert.NoError(t, err)

	input = genInput(t, "init", [][]byte{})
	funcName, params, returnType, err = ParseCallExtraDataByABI(lifeVm, input, []byte(abi1))
	assert.Equal(t, "init", funcName)
	assert.Equal(t, 0, len(params))
	assert.Equal(t, "void", returnType)
	assert.NoError(t, err)

	funcName, params, returnType, err = ParseCallExtraDataByABI(lifeVm, input, []byte(abi4))
	assert.Equal(t, "init", funcName)
	assert.Equal(t, 0, len(params))
	assert.Equal(t, "string", returnType)
	assert.NoError(t, err)

	num1 := utils.Uint64ToBytes(10)
	num2 := utils.Uint32ToBytes(10)
	num3 := utils.Uint16ToBytes(10)
	num4 := []byte{255}
	num5 := []byte{1}
	input = genInput(t, "init", [][]byte{num1, num2, num3, num4, num5})
	funcName, params, returnType, err = ParseCallExtraDataByABI(lifeVm, input, []byte(abi5))
	assert.Equal(t, "init", funcName)
	assert.Equal(t, 5, len(params))
	assert.Equal(t, "void", returnType)
	assert.NoError(t, err)
}

func TestParseCreateExtraData(t *testing.T) {
	_, _, _, err := ParseCreateExtraData(nil)
	assert.Equal(t, errEmptyInput, err)

	_, _, _, err = ParseCreateExtraData([]byte{123})
	assert.Equal(t, errInvalidRlpFormat, err)

	_, _, _, err = ParseCreateExtraData([]byte{1, 2, 3})
	assert.Equal(t, errInvalidRlpFormat, err)

	input, err := rlp.EncodeToBytes([]interface{}{})
	_, _, _, err = ParseCreateExtraData(input)
	assert.Equal(t, errInsufficientParams, err)

	input, err = rlp.EncodeToBytes([]interface{}{""})
	_, _, _, err = ParseCreateExtraData(input)
	assert.Equal(t, errInsufficientParams, err)

	input, err = rlp.EncodeToBytes([]interface{}{"code", "abi"})
	_, _, _, err = ParseCreateExtraData(input)
	assert.NoError(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{"code", "abi", "init"})
	_, _, _, err = ParseCreateExtraData(input)
	assert.NoError(t, err)
}
