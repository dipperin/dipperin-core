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

package utils

import (
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/tests/factory/vminfo"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	abiStr = `[{
        "name": "init",
        "inputs": [
            {
                "name": "inputName",
                "type": "string"
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
        "name": "test",
        "inputs": [
            {
                "name": "inputName",
                "type": "string"
            }
        ],
        "outputs": [],
        "constant": "false",
        "type": "function"
    }]`
)

func TestConvertInputs(t *testing.T) {

	_, abi := vminfo.GetTestData("token")

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

	type result struct {
		data []byte
		err error
	}

	testCases := []struct{
		name string
		given func() []byte
		expect result
	}{
		{
			name:"ErrEmptyInput",
			given: func() []byte {
				return nil
			},
			expect:result{[]byte(nil),gerror.ErrEmptyInput},
		},
		{
			name:"ErrInvalidRlpFormat",
			given: func() []byte {
				return []byte{123}
			},
			expect:result{[]byte(nil),gerror.ErrInvalidRlpFormat},
		},
		{
			name:"ErrLengthInputAbiNotMatch",
			given: func() []byte {
				data, err := rlp.EncodeToBytes([]interface{}{"DIPP", "WU"})
				assert.NoError(t, err)
				return data
			},
			expect:result{[]byte(nil),gerror.ErrLengthInputAbiNotMatch},
		},
		{
			name:"ConvertInputsRight",
			given: func() []byte {
				data, err := rlp.EncodeToBytes([]interface{}{"DIPP", "WU", uint64(100)})
				assert.NoError(t, err)
				return  data
				//return ConvertInputs(data, args)
			},
			expect:result{[]byte("DIPP,WU,100,"),nil},
		},
	}

	for _, tc := range testCases{
		t.Log("test case name ", tc.name)
		data:= tc.given()
		result, err := ConvertInputs(data, args)
		assert.Equal(t, tc.expect.data, result)
		assert.Equal(t, tc.expect.err, err)
	}

}

func TestParseCallContractData(t *testing.T) {

	_, abi := vminfo.GetTestData("token")
	type result struct {
		data []byte
		err error
	}
	expectData, err := rlp.EncodeToBytes([]interface{}{"init", []byte("DIPP"), []byte("WU"), Uint64ToBytes(100)})
	assert.NoError(t, err)
	startInput, err := rlp.EncodeToBytes([]interface{}{"getBalance", "owner"})
	assert.NoError(t, err)
	right2Input, err := rlp.EncodeToBytes([]interface{}{"init", "DIPP"})
	assert.NoError(t, err)

	testCases := []struct{
		name string
		given func() ([]byte, []byte)
		expect result
	}{
		{
			name: "ErrEmptyInput",
			given: func() ([]byte, []byte) {
				return abi,nil
			},
			expect: result{[]byte(nil), gerror.ErrEmptyInput},
		},
		{
			name: "ErrInvalidRlpFormat",
			given: func() ([]byte, []byte) {
				return abi,[]byte{123}
			},
			expect: result{[]byte(nil), gerror.ErrInvalidRlpFormat},
		},
		{
			name: "ErrInsufficientParams",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{})
				assert.NoError(t, err)
				return abi,input
			},
			expect: result{[]byte(nil), gerror.ErrInsufficientParams},
		},
		{
			name: "ErrLengthInputAbiNotMatch",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{"init"})
				assert.NoError(t, err)
				return abi,input
			},
			expect: result{[]byte(nil), gerror.ErrLengthInputAbiNotMatch},
		},
		{
			name: "ParseCallContractDataRight",
			given: func() ([]byte, []byte) {
				return abi,startInput
			},
			expect: result{startInput, nil},
		},
		{
			name: "ErrFuncNameNotFound",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{"test", ""})
				assert.NoError(t, err)
				return abi,input
			},
			expect: result{[]byte(nil), gerror.ErrFuncNameNotFound},
		},
		{
			name: "ErrLengthInputAbiNotMatch",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{"init", "DIPP,dipc"})
				assert.NoError(t, err)
				return []byte(abiStr),input
			},
			expect: result{[]byte(nil), gerror.ErrLengthInputAbiNotMatch},
		},
		{
			name: "ParseCallContractDataRight2",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{"init", "DIPP"})
				assert.NoError(t, err)
				return []byte(abiStr),input
			},
			expect: result{right2Input, nil},
		},
		{
			name: "ParseCallContractDataRight3",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{"init", "DIPP,WU,100"})
				assert.NoError(t, err)
				return abi,input
			},
			expect: result{expectData, nil},
		},
	}

	for _, tc := range testCases{
		t.Log("test case name ", tc.name)
		abiData, data := tc.given()
		result, err := ParseCallContractData(abiData,data)
		assert.Equal(t, tc.expect.data, result)
		assert.Equal(t, tc.expect.err, err)
	}

}

func TestParseCreateContractData(t *testing.T) {

	type result struct {
		data []byte
		err error
	}

	code, abi := vminfo.GetTestData("token")
	expectData, err := rlp.EncodeToBytes([]interface{}{code, abi, []byte("DIPP"), []byte("WU"), Uint64ToBytes(100)})
	assert.NoError(t, err)
	expectData2, err := rlp.EncodeToBytes([]interface{}{code, []byte(abi2)})
	assert.NoError(t, err)

	testCases := []struct{
		name string
		given func() []byte
		expect result
	}{
		{
			name:"ErrEmptyInput",
			given: func() []byte {
				return nil
			},
			expect:result{[]byte(nil), gerror.ErrEmptyInput},
		},
		{
			name:"ErrInvalidRlpFormat",
			given: func() []byte {
				return []byte{123}
			},
			expect:result{[]byte(nil), gerror.ErrInvalidRlpFormat},
		},
		{
			name:"ErrInsufficientParams",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{code})
				assert.NoError(t, err)
				return input
			},
			expect:result{[]byte(nil), gerror.ErrInsufficientParams},
		},
		{
			name:"ErrLengthInputAbiNotMatch",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{code, abi})
				assert.NoError(t, err)
				return input
			},
			expect:result{[]byte(nil), gerror.ErrLengthInputAbiNotMatch},
		},
		{
			name:"ErrInvalidOutputLength",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{code, []byte(abi1)})
				assert.NoError(t, err)
				return input
			},
			expect:result{[]byte(nil), gerror.ErrInvalidOutputLength},
		},
		{
			name:"ParseCreateContractDataRight",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{code, []byte(abi2)})
				assert.NoError(t, err)
				return input
			},
			expect:result{expectData2, nil},
		},
		{
			name:"ErrLengthInputAbiNotMatch",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{code, abi, "DIPP"})
				assert.NoError(t, err)
				return input
			},
			expect:result{[]byte(nil), gerror.ErrLengthInputAbiNotMatch},
		},
		{
			name:"ErrFuncNameNotFound",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{code, []byte(abi3), "DIPP"})
				assert.NoError(t, err)
				return input
			},
			expect:result{[]byte(nil), gerror.ErrFuncNameNotFound},
		},
		{
			name:"ParseCreateContractData",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{code, abi, "DIPP,WU,100"})
				assert.NoError(t, err)
				//return  ParseCreateContractData(input)
				return input
			},
			expect:result{expectData, nil},
		},
	}

	for _, tc := range testCases{
		t.Log("test case name ", tc.name)
		data := tc.given()
		result, err := ParseCreateContractData(data)
		assert.Equal(t, tc.expect.data, result)
		assert.Equal(t, tc.expect.err, err)
	}
}

func TestParseInputForFuncName(t *testing.T) {
	type  result struct {
		funcName string
		err error
	}
	testCases := []struct{
		name string
		given func()  []byte
		expect result
	}{
		{
			name:"ErrEmptyInput",
			given: func() []byte {
				return nil
			},
			expect:result{"", gerror.ErrEmptyInput},
		},
		{
			name:"ErrEmptyInput",
			given: func() []byte {
				return 	[]byte{}

			},
			expect:result{"", gerror.ErrEmptyInput},
		},
		{
			name:"ErrInvalidRlpFormat",
			given: func() []byte {
				return []byte{1, 2, 3}

			},
			expect:result{"", gerror.ErrInvalidRlpFormat},
		},
		{
			name:"ErrInsufficientParams",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{})
				assert.NoError(t, err)
				return input
			},
			expect:result{"", gerror.ErrInsufficientParams},
		},
		{
			name:"ParseInputForFuncNameRight",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{""})
				assert.NoError(t, err)
				return input
			},
			expect:result{"", nil},
		},
		{
			name:"ParseInputForFuncNameRight",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{"funcName"})
				assert.NoError(t, err)
				return input
			},
			expect:result{"funcName", nil},
		},
		{
			name:"ParseInputForFuncNameRight",
			given: func() []byte {
				input, err := rlp.EncodeToBytes([]interface{}{"funcName",""})
				assert.NoError(t, err)
				return input
				//return ParseInputForFuncName(input)
			},
			expect:result{"funcName", nil},
		},
	}


	for _, tc := range testCases {
		input := tc.given()
		funcName, err := ParseInputForFuncName(input)
		assert.Equal(t, tc.expect.err, err)
		assert.Equal(t, tc.expect.funcName, funcName)
	}
}
