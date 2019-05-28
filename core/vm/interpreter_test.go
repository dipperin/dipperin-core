package vm

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"io/ioutil"
	"bytes"
	"github.com/dipperin/dipperin-core/common"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
)

func testGetHash(blockNumber uint64) common.Hash{
	return common.Hash{}
}

func getTestVm() *VM{
	return NewVM(Context{BlockNumber:big.NewInt(1),GetHash:testGetHash}, fakeStateDB{},DEFAULT_VM_CONFIG)
}

func TestWASMInterpreter_Run_small(t *testing.T) {
	fileCode, err := ioutil.ReadFile("./small/small.wasm")
	assert.NoError(t, err)

	var testPath = "./small"
	contract := getContract(t, contractAddr, testPath + "/small.wasm", testPath + "/small.cpp.abi.json")

	testVm := getTestVm()
	interpreter :=testVm.Interpreter

	param := [][]byte{[]byte("hello")}
	fmt.Println(param)
	result, err := interpreter.Run(testVm,contract, genInput(t, "hello", param),false)
	assert.Equal(t, make([]byte, 32), result)
	assert.NoError(t, err)

	result, err = interpreter.Run(testVm,contract, nil,false)
	assert.Equal(t, fileCode, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_parseInputFromAbi(t *testing.T) {
	fileCode, err := ioutil.ReadFile("./small/small.wasm")
	assert.NoError(t, err)

	abi, err := ioutil.ReadFile("./small/small.cpp.abi.json")
	assert.NoError(t, err)

	param := [][]byte{[]byte("world")}
	input := genInput(t, "hello", param)
	vm, err := exec.NewVirtualMachine(fileCode, DEFAULT_VM_CONFIG, nil, nil)
	txType, funcName, params, returnType, err := parseInputFromAbi(vm, input, abi)

	assert.NoError(t, err)
	assert.Equal(t, 1, txType)
	assert.Equal(t, "hello", funcName)
	assert.Equal(t, "world", string(vm.Memory.Memory[params[0]:params[0]+5]))
	assert.Equal(t, "void", returnType)
}

func TestWASMInterpreter_Run_testcontract(t *testing.T) {
	//var testPath = "/home/qydev/go/src/github.com/PlatONnetwork/PlatON-CDT/build/bin/testcontract"
	var testPath = "./testcontract"
	contract := getContract(t, contractAddr, testPath + "/testcontract.wasm", testPath + "/testcontract.cpp.abi.json")

	testVm := getTestVm()
	interpreter :=testVm.Interpreter

	param := [][]byte{utils.Int32ToBytes(123), utils.Int32ToBytes(456)}
	result, err := interpreter.Run(testVm,contract, genInput(t, "saveKeyValue", param),false)
	assert.Equal(t, make([]byte, 32), result)
	assert.NoError(t, err)

	param = [][]byte{utils.Int32ToBytes(123)}
	result, err = interpreter.Run(testVm,contract, genInput(t, "getKeyValue", param),false)
	assert.Equal(t, make([]byte, 32), result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_example3(t *testing.T) {
	//var testPath = "/home/qydev/go/src/github.com/PlatONnetwork/PlatON-CDT/build/bin/example3"
	var testPath = "./example3"
	contract := getContract(t, contractAddr, testPath + "/example3.wasm", testPath + "/example3.cpp.abi.json")

	testVm := getTestVm()
	interpreter :=testVm.Interpreter

	expect := make([]byte, 32)
	param := [][]byte{utils.Int32ToBytes(6666)}

	inputs := genInput(t, "setBalance", param)
	result, err := interpreter.Run(testVm,contract, inputs,false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	expect = append(expect[:28], utils.Int32ToBytes(6666)...)
	result, err = interpreter.Run(testVm,contract, genInput(t, "getBalance", [][]byte{}),false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_map_int(t *testing.T) {
	//var testPath = "/home/qydev/go/src/github.com/PlatONnetwork/PlatON-CDT/build/bin/sMap"
	var testPath = "./map-int"
	contract := getContract(t, contractAddr, testPath + "/sMap.wasm", testPath + "/sMap.cpp.abi.json")

	testVm := getTestVm()
	interpreter :=testVm.Interpreter

	key := utils.Int32ToBytes(2147483647)
	value := utils.Int64ToBytes(9223372036854775807)
	expect := make([]byte, 32)
	param := [][]byte{key, value}

	inputs := genInput(t, "setBalance", param)
	result, err := interpreter.Run(testVm,contract, inputs,false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	/*	fmt.Println("-----------------------------------------")

		key1 := utils.Int32ToBytes(111)
		value1 := utils.Int32ToBytes(222)
		param1 := [][]byte{key1, value1}
		inputs = genInput(t, "setBalance", param1)

		result, err = Interpreter.Run(contract, inputs)
		assert.Equal(t, expect, result)
		assert.NoError(t, err)*/

	fmt.Println("-----------------------------------------")
	result, err = interpreter.Run(testVm,contract, genInput(t, "getBalance", [][]byte{key}),false)
	expect = append(expect[:24], value...)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_map_string(t *testing.T) {
	var testPath = "./map-string"
	contract := getContract(t, contractAddr, testPath + "/map2.wasm", testPath + "/StringMap.cpp.abi.json")

	testVm := getTestVm()
	interpreter :=testVm.Interpreter

	key := []byte("balance")
	value := utils.Int32ToBytes(255)

	expect := make([]byte, 32)
	param := [][]byte{key, value}

	inputs := genInput(t, "setBalance", param)
	result, err := interpreter.Run(testVm,contract, inputs,false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")

	key1 := []byte("bbb")
	value1 := utils.Int32ToBytes(222)
	param1 := [][]byte{key1, value1}

	inputs = genInput(t, "setBalance", param1)
	result, err = interpreter.Run(testVm,contract, inputs,false)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)

	fmt.Println("-----------------------------------------")
	result, err = interpreter.Run(testVm,contract, genInput(t, "getBalance", [][]byte{key}),false)
	expect = append(expect[:28], value...)
	assert.Equal(t, expect, result)
	assert.NoError(t, err)
}

func TestWASMInterpreter_Run_event(t *testing.T) {
	//var testPath = "/home/qydev/go/src/github.com/PlatONnetwork/PlatON-CDT/build/bin/event"
	var testPath = "./event"
	contract := getContract(t, contractAddr, testPath + "/event.wasm", testPath + "/event.cpp.abi.json")

	testVm := getTestVm()
	interpreter :=testVm.Interpreter


	name := []byte("0000")
	num := utils.Int64ToBytes(2)
	log.Info("the num is:","num",num)

	/*	name := []byte("logName")
		test, err := rlp.EncodeToBytes(append(name, []byte{133, 101, 118, 101, 110, 116, 127}...))
		fmt.Println(test)*/

	expect := make([]byte, 32)
	param := [][]byte{name, num}

	inputs := genInput(t, "hello", param)
	result, err := interpreter.Run(testVm,contract, inputs,false)
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
	for _, v := range (param) {
		input = append(input, v)
	}

	buffer := new(bytes.Buffer)
	err := rlp.Encode(buffer, input)
	assert.NoError(t, err)
	return buffer.Bytes()
}

func getContract(t *testing.T, addr common.Address, code, abi string) *Contract {
	fileCode, err := ioutil.ReadFile(code)
	assert.NoError(t, err)

	fileABI, err := ioutil.ReadFile(abi)
	assert.NoError(t, err)

	return &Contract{
		self: fakeContractRef{addr: addr},
		Code: fileCode,
		ABI:  fileABI,
	}
}
