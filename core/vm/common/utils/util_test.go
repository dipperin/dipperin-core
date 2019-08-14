package utils

import (
	"fmt"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	abiStr = `[{
        "name": "test",
        "inputs": [
            {
                "name": "inputName",
                "type": "unknown"
            }
        ],
        "outputs": [],
        "constant": "false",
        "type": "function"
    }]`

	abi1 = `[{
        "name": "init",
        "inputs": [],
        "outputs": [
            {
                "name": "inputName",
                "type": "unknown"
            }
        ],
        "constant": "false",
        "type": "function"
    }]`

	abi2 = `[{
        "name": "init",
        "inputs": [],
        "outputs": [],
        "constant": "false",
        "type": "function"
    }]`

	abi3 = `[{
        "name": "init",
        "inputs": [
            {
                "name": "inputName",
                "type": "unknown"
            }
        ],
        "outputs": [],
        "constant": "false",
        "type": "function"
    }]`
)

func TestConvertInputs(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	_, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

	wasmAbi := new(WasmAbi)
	err := wasmAbi.FromJson(abi)
	assert.NoError(t, err)

	var args []InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold(v.Name, "init") && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
			break
		}
	}

	_, err = ConvertInputs(nil, args)
	assert.Equal(t, errEmptyInput, err)

	_, err = ConvertInputs([]byte{123}, args)
	assert.Equal(t, errInvalidRlpFormat, err)

	data, err := rlp.EncodeToBytes([]interface{}{"DIPP", "WU"})
	assert.NoError(t, err)
	_, err = ConvertInputs(data, args)
	assert.Equal(t, errLengthInputAbiNotMatch, err)

	data, err = rlp.EncodeToBytes([]interface{}{"DIPP", "WU", uint64(100)})
	assert.NoError(t, err)
	convertData, err := ConvertInputs(data, args)
	assert.Equal(t, "DIPP,WU,100,", string(convertData))
	assert.NoError(t, err)
}

func TestParseCallContractData(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	_, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

	expectData, err := rlp.EncodeToBytes([]interface{}{"init", []byte("DIPP"), []byte("WU"), Uint64ToBytes(100)})
	assert.NoError(t, err)

	data, err := ParseCallContractData(abi, nil)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errEmptyInput, err)

	data, err = ParseCallContractData(abi, []byte{123})
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errInvalidRlpFormat, err)

	input, err := rlp.EncodeToBytes([]interface{}{})
	assert.NoError(t, err)
	data, err = ParseCallContractData(abi, input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errInsufficientParams, err)

	input, err = rlp.EncodeToBytes([]interface{}{"init"})
	assert.NoError(t, err)
	data, err = ParseCallContractData(abi, input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errLengthInputAbiNotMatch, err)

	input, err = rlp.EncodeToBytes([]interface{}{"start"})
	assert.NoError(t, err)
	data, err = ParseCallContractData(abi, input)
	assert.Equal(t, input, data)
	assert.NoError(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{"test", ""})
	assert.NoError(t, err)
	data, err = ParseCallContractData(abi, input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errFuncNameNotFound, err)

	data, err = ParseCallContractData([]byte{123}, input)
	assert.Equal(t, []byte(nil), data)
	assert.Error(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{"init", "DIPP"})
	assert.NoError(t, err)
	data, err = ParseCallContractData(abi, input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errLengthInputAbiNotMatch, err)

	input, err = rlp.EncodeToBytes([]interface{}{"test", "DIPP"})
	assert.NoError(t, err)
	data, err = ParseCallContractData([]byte(abiStr), input)
	assert.Equal(t, []byte(nil), data)
	assert.Error(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{"init", "DIPP,WU,100"})
	assert.NoError(t, err)
	data, err = ParseCallContractData(abi, input)
	assert.Equal(t, expectData, data)
	assert.NoError(t, err)
}

func TestParseCreateContractData(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)
	fmt.Println("code", WASMPath)
	fmt.Println("abi", AbiPath)

	data, err := ParseCreateContractData(nil)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errEmptyInput, err)

	data, err = ParseCreateContractData([]byte{123})
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errInvalidRlpFormat, err)

	input, err := rlp.EncodeToBytes([]interface{}{code})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errInsufficientParams, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, []byte{123}})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, []byte(nil), data)
	assert.Error(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, abi})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errLengthInputAbiNotMatch, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, []byte(abi1)})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errInvalidOutputLength, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, []byte(abi2)})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, input, data)
	assert.NoError(t, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, abi, "DIPP"})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errLengthInputAbiNotMatch, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, []byte(abiStr), "DIPP"})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, []byte(nil), data)
	assert.Equal(t, errFuncNameNotFound, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, []byte(abi3), "DIPP"})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, []byte(nil), data)
	assert.Error(t, errFuncNameNotFound, err)

	input, err = rlp.EncodeToBytes([]interface{}{code, abi, "DIPP,WU,100"})
	assert.NoError(t, err)
	expectData, err := rlp.EncodeToBytes([]interface{}{code, abi, []byte("DIPP"), []byte("WU"), Uint64ToBytes(100)})
	assert.NoError(t, err)
	data, err = ParseCreateContractData(input)
	assert.Equal(t, expectData, data)
	assert.NoError(t, err)
}
