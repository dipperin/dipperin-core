package vm

import (
	"dipperin-vm/common/utils"
	"github.com/dipperin/dipperin-core/common/vmcommon"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestWASMInterpreter_Run_map_string(t *testing.T) {
/*	var testPath = "./map-string"
	testVm := getTestVm()
	interpreter := testVm.Interpreter

	key := []byte("balance")
	value := vmcommon.Int32ToBytes(255)

	expect := make([]byte, 32)
	param := [][]byte{key, value}

	inputs := genInput(t, "setBalance", param)
	contract := getContract(t, contractAddr, testPath+"/map2.wasm", testPath+"/StringMap.cpp.abi.json", inputs)
	result, err := interpreter.Run(testVm, contract, false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")

	key1 := []byte("bbb")
	value1 := vmcommon.Int32ToBytes(222)
	param1 := [][]byte{key1, value1}

	inputs = genInput(t, "setBalance", param1)
	contract = getContract(t, contractAddr, testPath+"/map2.wasm", testPath+"/StringMap.cpp.abi.json", inputs)
	result, err = interpreter.Run(testVm, contract, false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")
	inputs = genInput(t, "getBalance", [][]byte{key})
	contract = getContract(t, contractAddr, testPath+"/map2.wasm", testPath+"/StringMap.cpp.abi.json", inputs)
	result, err = interpreter.Run(testVm, contract, false)
	expect = append(expect[:28], value...)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)*/
}

func TestWASMInterpreter_Run_event(t *testing.T) {
	//var testPath = "/home/qydev/go/src/github.com/PlatONnetwork/PlatON-CDT/build/bin/event"
	var testPath = "./event"
	testVm := getTestVm()
	interpreter := testVm.Interpreter

	name := []byte("RunEvent")
	num := vmcommon.Int64ToBytes(100)
	log.Info("the num is:", "num", num)

	expect := make([]byte, 32)
	param := [][]byte{name, num}

	inputs := genInput(t, "hello", param)
	contract := getContract(t, contractAddr, testPath+"/event.wasm", testPath+"/event.cpp.abi.json", inputs)
	result, err := interpreter.Run(testVm, contract, false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func TestParseInputFromAbiByInit(t *testing.T)  {
	abiBytes, err := ioutil.ReadFile("./event/token/token.cpp.abi.json")
	assert.NoError(t, err)
	var wasmAbi utils.WasmAbi
	err = wasmAbi.FromJson(abiBytes)
	//err = json.Unmarshal(abiBytes, &wasmAbi.AbiArr)
	assert.NoError(t, err)

	var args []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold("init", v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
		}
	}

	params := []string{"dipp", "DIPP", "100000000"}
	wasmBytes, err := ioutil.ReadFile("./event/token/token5.wasm")
	assert.NoError(t, err)

	rlpParams := []interface{}{
		wasmBytes, abiBytes,
	}

	inputParams := []interface{}{}
	for i, v := range args {
		bts := params[i]
		re, err := vmcommon.StringConverter(bts, v.Type)
		assert.NoError(t, err)
		rlpParams = append(rlpParams, re)
		inputParams = append(inputParams, re)
	}

	//data, err :=  rlp.EncodeToBytes(rlpParams)
	input, err := rlp.EncodeToBytes(inputParams)
	assert.NoError(t, err)

	//　life方法注入新建虚拟机
	//solver := resolver.NewResolver(vm, contract, in.state)
	//lifeVm, err := exec.NewVirtualMachine(wasmBytes, in.config, solver, nil)

	funcName, pms, returnType, err := parseInputFromAbiByInit(nil,input, abiBytes)
	log.Info("result", "funcName", funcName, "pms", pms, "returnType", returnType, "err", err)
}