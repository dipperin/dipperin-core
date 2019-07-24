package vm

import (
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWASMInterpreter_Run_map_string(t *testing.T) {
	/*	WASMPath := g_testData.GetWASMPath("map-string")
		AbiPath := g_testData.GetAbiPath("map-string")*/
	/*	WASMPath := "/home/qydev/go/src/dipperin-c/dipc/testcontract/mapString/mapString.wasm"
		AbiPath := "/home/qydev/go/src/dipperin-c/dipc/testcontract/mapString/mapString.cpp.abi.json"
		testVm := getTestVm()
		interpreter := testVm.Interpreter

		key := []byte("balance")
		value := utils.Int32ToBytes(255)

		expect := make([]byte, 32)
		param := [][]byte{key, value}

		inputs := genInput(t, "setBalance", param)
		contract := getContract(t, contractAddr, WASMPath, AbiPath, inputs)
		result, err := interpreter.Run(testVm, contract, false)
		assert.Equal(t, expect, result)
		assert.NoError(t, err)

		fmt.Println("-----------------------------------------")*/
	/*
		key1 := []byte("bbb")
		value1 := utils.Int32ToBytes(222)
		param1 := [][]byte{key1, value1}

		inputs = genInput(t, "setBalance", param1)
		contract = getContract(t, contractAddr, WASMPath, AbiPath, inputs)
		result, err = interpreter.Run(testVm, contract, false)
		assert.Equal(t, expect, result)
		assert.NoError(t, err)

		fmt.Println("-----------------------------------------")
		inputs = genInput(t, "getBalance", [][]byte{key})
		contract = getContract(t, contractAddr, WASMPath, AbiPath, inputs)
		result, err = interpreter.Run(testVm, contract, false)
		expect = append(expect[:28], value...)
		assert.Equal(t, expect, result)
		assert.NoError(t, err)*/
}

func TestWASMInterpreter_Run_event(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event",g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event",g_testData.CoreVmTestData)
	testVm := getTestVm()
	interpreter := testVm.Interpreter

	name := []byte("RunEvent")
	num := utils.Int64ToBytes(100)
	log.Info("the num is:", "num", num)

	expect := make([]byte, 32)
	param := [][]byte{name, num}

	inputs := genInput(t, "hello", param)
	contract := getContract(t, contractAddr, WASMPath, AbiPath, inputs)
	result, err := interpreter.Run(testVm, contract, false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_DIPCLibContract(t *testing.T) {
	t.Skip()
	testVm := getTestVm()
	interpreter := testVm.Interpreter
	inputs := genInput(t, g_testData.ContractTestPar.CallFuncName,[][]byte{})
	log.Info("the wasmPath is:","wasmPath",g_testData.ContractTestPar.WASMPath)
	log.Info("the abiPath is:","abiPath",g_testData.ContractTestPar.AbiPath)
	contract := getContract(t, contractAddr, g_testData.ContractTestPar.WASMPath, g_testData.ContractTestPar.AbiPath, inputs)
	_, err := interpreter.Run(testVm, contract, false)
	assert.NoError(t,err)
}
