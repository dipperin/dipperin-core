package vm

import (
	"bytes"
	"fmt"
	common2 "github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func getTestContract(wasmFile, abiFile string, contractAddr common2.Address, t *testing.T) *Contract {

	fileCode, err := ioutil.ReadFile(wasmFile)
	assert.NoError(t, err)

	fileABI, err := ioutil.ReadFile(abiFile)
	assert.NoError(t, err)

	contract := &Contract{
		self: &Caller{contractAddr},
		Code: fileCode,
		ABI:  fileABI,
	}
	return contract
}

func TestWASMInterpreter_Run_small(t *testing.T) {
	contract := getTestContract("./small/small.wasm", "./small/small.cpp.abi.json", contractAddr, t)

	interpreter := NewWASMInterpreter(fakeStateDB{}, Context{}, DEFAULT_VM_CONFIG)

	param := [][]byte{[]byte("world")}
	result, err := interpreter.Run(contract, genInput(t, "hello", param))
	assert.Equal(t, make([]byte, 32), result)
	assert.NoError(t, err)

	result, err = interpreter.Run(contract, nil)
	assert.Equal(t, contract.Code, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_parseInputFromAbi(t *testing.T) {
	contract := getTestContract("./small/small.wasm", "./small/small.cpp.abi.json", contractAddr, t)

	param := [][]byte{[]byte("world")}
	input := genInput(t, "hello", param)
	vm, err := exec.NewVirtualMachine(contract.Code, DEFAULT_VM_CONFIG, nil, nil)
	txType, funcName, params, returnType, err := parseInputFromAbi(vm, input, contract.ABI)

	assert.NoError(t, err)
	assert.Equal(t, 1, txType)
	assert.Equal(t, "hello", funcName)
	assert.Equal(t, "world", string(vm.Memory.Memory[params[0]:params[0]+5]))
	assert.Equal(t, "void", returnType)
}

func TestWASMInterpreter_Run_testcontract(t *testing.T) {
	contract := getTestContract("./testcontract/testcontract.wasm", "./testcontract/testcontract.cpp.abi.json", contractAddr, t)

	interpreter := NewWASMInterpreter(fakeStateDB{}, Context{}, DEFAULT_VM_CONFIG)
	param := [][]byte{utils.Int32ToBytes(123), utils.Int32ToBytes(456)}
	result, err := interpreter.Run(contract, genInput(t, "saveKeyValue", param))
	assert.Equal(t, make([]byte, 32), result)
	assert.NoError(t, err)

	param = [][]byte{utils.Int32ToBytes(123)}
	result, err = interpreter.Run(contract, genInput(t, "getKeyValue", param))
	assert.Equal(t, make([]byte, 32), result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_example3(t *testing.T) {
	contract := getTestContract("./example3/example3.wasm", "./example3/example3.cpp.abi.json", contractAddr, t)
	interpreter := NewWASMInterpreter(NewStorage(), Context{}, DEFAULT_VM_CONFIG)

	expect := make([]byte, 32)
	param := [][]byte{utils.Int32ToBytes(6666)}

	inputs := genInput(t, "setBalance", param)
	result, err := interpreter.Run(contract, inputs)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	expect = append(expect[:28], utils.Int32ToBytes(6666)...)
	result, err = interpreter.Run(contract, genInput(t, "getBalance", [][]byte{}))
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_map_int(t *testing.T) {
	contract := getTestContract("./map-int/sMap.wasm", "./map-int/sMap.cpp.abi.json", contractAddr, t)

	interpreter := NewWASMInterpreter(NewStorage(), Context{}, DEFAULT_VM_CONFIG)

	key := utils.Int32ToBytes(333)
	value := utils.Int32ToBytes(444)

	expect := make([]byte, 32)
	param := [][]byte{key, value}

	inputs := genInput(t, "setBalance", param)
	result, err := interpreter.Run(contract, inputs)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")

	key1 := utils.Int32ToBytes(111)
	value1 := utils.Int32ToBytes(222)
	param1 := [][]byte{key1, value1}
	inputs = genInput(t, "setBalance", param1)
	result, err = interpreter.Run(contract, inputs)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")
	result, err = interpreter.Run(contract, genInput(t, "getBalance", [][]byte{key}))
	expect = append(expect[:28], value...)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_map_string(t *testing.T) {
	contract := getTestContract("./map-string/map2.wasm", "./map-string/StringMap.cpp.abi.json", contractAddr, t)
	interpreter := NewWASMInterpreter(NewStorage(), Context{}, DEFAULT_VM_CONFIG)

	key := []byte("aaa")
	value := utils.Int32ToBytes(111)

	expect := make([]byte, 32)
	param := [][]byte{key, value}

	inputs := genInput(t, "setBalance", param)
	result, err := interpreter.Run(contract, inputs)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")

	key1 := []byte("bbb")
	value1 := utils.Int32ToBytes(222)
	param1 := [][]byte{key1, value1}

	inputs = genInput(t, "setBalance", param1)
	result, err = interpreter.Run(contract, inputs)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")
	result, err = interpreter.Run(contract, genInput(t, "getBalance", [][]byte{key}))
	expect = append(expect[:28], value...)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func genInput(t *testing.T, funcName string, param [][]byte) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	// tx type
	input = append(input, utils.Int64ToBytes(1))
	// func name
	input = append(input, []byte(funcName))
	// func parameter
	for _, v := range param {
		input = append(input, v)
	}

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	assert.NoError(t, err)
	return buffer.Bytes()
}
