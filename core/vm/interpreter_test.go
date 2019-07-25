package vm

import (
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
)

func TestWASMInterpreter_Run(t *testing.T) {
	WASMPath := g_testData.GetWasmPath("event")
	AbiPath := g_testData.GetAbiPath("event")
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

	name = []byte("contract")
	inputs = genInput(t, "returnInt", [][]byte{name})
	contract.Input = inputs
	result, err = interpreter.Run(testVm, contract, false)
	resp = utils.Align32BytesConverter(result, "int64")
	assert.Equal(t, int64(50), resp)
	assert.NoError(t, err)

	name = []byte("contract")
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
	assert.Error(t, err)

	WASMPath := g_testData.GetWasmPath("event")
	AbiPath := g_testData.GetAbiPath("event")
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)
	contract.Code = code
	contract.ABI = abi
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Error(t, err)

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
	assert.Error(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{"returnString"})
	assert.NoError(t, err)
	contract.Input = input
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Error(t, err)

	name := []byte("string")
	param := [][]byte{name}
	input = genInput(t, "returnString", param)
	contract.Input = input
	result, err = interpreter.Run(testVm, contract, false)
	assert.Nil(t, result)
	assert.Error(t, err)

	contract.Input = []byte{123}
	result, err = interpreter.Run(testVm, contract, true)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestParseInputForFuncName(t *testing.T) {
	funcName, err := ParseInputForFuncName(nil)
	assert.Equal(t, "", funcName)
	assert.Equal(t, errEmptyInput, err)

	funcName, err = ParseInputForFuncName([]byte{1,2,3})
	assert.Equal(t, "", funcName)
	assert.Equal(t, errReturnInvalidRlpFormat, err)

	input, err := rlp.EncodeToBytes([]interface{}{})
	funcName, err = ParseInputForFuncName(input)
	assert.Equal(t, "", funcName)
	assert.Equal(t, errReturnInsufficientParams, err)

	input, err = rlp.EncodeToBytes([]interface{}{"funcName"})
	funcName, err = ParseInputForFuncName(input)
	assert.Equal(t, "funcName", funcName)
	assert.NoError(t, err)
}