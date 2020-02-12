package vm

import (
	"encoding/json"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/vm/base"
	"github.com/dipperin/dipperin-core/core/vm/base/utils"
	"github.com/dipperin/dipperin-core/tests/factory/vminfo"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	eventContractName = "event"
	callFuncName      = "returnString"
)

func Test_findParams(t *testing.T) {
	ctrl, _, _ := GetBaseVmInfo(t)
	defer ctrl.Finish()

	code, abi := vminfo.GetTestData(eventContractName)
	lifeVm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, nil, nil)
	assert.NoError(t, err)
	input, err := rlp.EncodeToBytes([]interface{}{"winner"})
	assert.NoError(t, err)
	t.Log("input info", input, string(input))

	type result struct {
		err        error
		params     []int64
		returnType string
	}

	testCases := []struct {
		name   string
		given  func() result
		expect result
	}{
		{
			name: "abiErr",
			given: func() result {
				_, _, err = findParams(lifeVm, []byte{}, callFuncName, []interface{}{[]byte{}})
				return result{err, []int64{}, ""}
			},
			expect: result{gerror.ErrInvalidAbi, []int64{}, ""},
		},
		{
			name: "errInputParam",
			given: func() result {
				_, _, err = findParams(lifeVm, abi, callFuncName, []interface{}{})
				return result{err, []int64{}, ""}
			},
			expect: result{gerror.ErrInputAbiNotMatch, []int64{}, ""},
		},
		{
			name: "CallFindParamsRight",
			given: func() result {
				params, returnType, err := findParams(lifeVm, abi, callFuncName, []interface{}{input})
				return result{err, params, returnType}
			},
			expect: result{err, []int64{131072}, "string"},
		},
	}

	for _, tc := range testCases {
		result := tc.given()
		t.Log("result", result)
		if result.err != nil {
			assert.Equal(t, result.err.Error(), tc.expect.err.Error())
		} else {
			assert.NoError(t, result.err)
			assert.Equal(t, result.params, tc.expect.params)
			assert.Equal(t, result.returnType, tc.expect.returnType)
		}
	}
}

func Test_ParseCreateExtraData(t *testing.T) {
	code, abi := vminfo.GetTestData(eventContractName)

	type result struct {
		code    []byte
		abi     []byte
		rlpInit []byte
		err     error
	}

	testCases := []struct {
		name   string
		given  func() []byte
		expect result
	}{
		{
			name: "ErrEmptyInput",
			given: func() []byte {
				return []byte{}
			},
			expect: result{nil, nil, nil, gerror.ErrEmptyInput},
		},
		{
			name: "ErrInvalidRlpFormat",
			given: func() []byte {
				rlpDataErr := "errData"
				return []byte(rlpDataErr)
			},
			expect: result{nil, nil, nil, gerror.ErrInvalidRlpFormat},
		},
		{
			name: "ErrInsufficientParams",
			given: func() []byte {
				rlpData, err := rlp.EncodeToBytes([]interface{}{code})
				assert.NoError(t, err)
				return []byte(rlpData)
			},
			expect: result{nil, nil, nil, gerror.ErrInsufficientParams},
		},
		{
			name: "ParseCreateExtraDataRight",
			given: func() []byte {
				rlpData, err := rlp.EncodeToBytes([]interface{}{code, abi})
				assert.NoError(t, err)
				return []byte(rlpData)
			},
			expect: result{code, abi, nil, nil},
		},
		{
			name: "ParseCreateExtraDataRightMoreParams",
			given: func() []byte {
				rlpData, err := rlp.EncodeToBytes([]interface{}{code, abi,"params"})
				assert.NoError(t, err)
				return []byte(rlpData)
			},
			expect: result{code, abi, []byte("params"), nil},
		},
	}

	for _, tc := range testCases {
		input := tc.given()

		code, abi, rlpInit, err := ParseCreateExtraData(input)

		if err != nil {
			assert.Equal(t, tc.expect.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expect.code, code)
			assert.Equal(t, tc.expect.abi, abi)
			if tc.expect.rlpInit != nil {
				result, err := rlp.EncodeToBytes([]interface{}{string(tc.expect.rlpInit)})
				assert.NoError(t, err)
				assert.Equal(t, result, rlpInit)
			}
		}
	}

}

func Test_ParseCallExtraDataByABI(t *testing.T) {
	ctrl, _, _ := GetBaseVmInfo(t)
	defer ctrl.Finish()

	code, abi := vminfo.GetTestData(eventContractName)
	lifeVm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, nil, nil)
	assert.NoError(t, err)
	//input, err := rlp.EncodeToBytes([]interface{}{"winner"})
	//assert.NoError(t, err)

	type result struct {
		funcName   string
		params     []int64
		returnType string
		err        error
	}

	testCases := []struct {
		name   string
		given  func() ([]byte, []byte)
		expect result
	}{
		{
			name: "ErrEmptyInput",
			given: func() ([]byte, []byte) {
				return []byte{},[]byte{}
			},
			expect: result{"", []int64{}, "", gerror.ErrEmptyInput},
		},
		{
			name: "ErrEmptyABI",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes("result")
				assert.NoError(t, err)
				return input,[]byte{}
			},
			expect: result{"", []int64{}, "", gerror.ErrEmptyABI},
		},
		{
			name: "ErrInvalidRlpFormat",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes("result")
				assert.NoError(t, err)
				return input,abi
			},
			expect: result{"", []int64{}, "", gerror.ErrInvalidRlpFormat},
		},
		{
			name: "ErrInsufficientParams",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{})
				assert.NoError(t, err)
				return input,abi
			},
			expect: result{callFuncName, []int64{131072}, "string", gerror.ErrInsufficientParams},
		},
		{
			name: "ParseCallExtraDataByABIRight",
			given: func() ([]byte, []byte) {
				input, err := rlp.EncodeToBytes([]interface{}{callFuncName, "winner"})
				assert.NoError(t, err)
				return input,abi
			},
			expect: result{callFuncName, []int64{131072}, "string", nil},
		},
	}

	for _, tc := range testCases {

		input,abi := tc.given()
		funcName, params, returnType, err := ParseCallExtraDataByABI(lifeVm, input, abi)
		if err != nil {
			assert.Equal(t, tc.expect.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expect.returnType, returnType)
			assert.Equal(t, tc.expect.params, params)
			assert.Equal(t, tc.expect.funcName, funcName)
		}
	}
}

func Test_ParseInitFunctionByABI(t *testing.T) {
	ctrl, _, _ := GetBaseVmInfo(t)
	defer ctrl.Finish()

	type result struct {
		params     []int64
		returnType string
		err        error
	}

	testCases := []struct {
		name   string
		given  func() (*exec.VirtualMachine, []byte, []byte)
		expect result
	}{
		{
			name: "ErrEmptyABI",
			given: func() (*exec.VirtualMachine, []byte, []byte){
				input, err := rlp.EncodeToBytes("result")
				assert.NoError(t, err)
				code, _ := vminfo.GetTestData(eventContractName)
				lifeVm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, nil, nil)
				assert.NoError(t, err)
				return lifeVm, input, []byte{}
			},
			expect: result{[]int64{}, "", gerror.ErrEmptyABI},
		},
		{
			name: "ErrInvalidReturnType",
			given: func() (*exec.VirtualMachine, []byte, []byte){
				input, err := rlp.EncodeToBytes([]interface{}{})
				assert.NoError(t, err)
				code, abi := vminfo.GetTestData(eventContractName)

				wasmAbi := new(utils.WasmAbi)
				err = wasmAbi.FromJson(abi)
				assert.NoError(t, err)
				wasmAbi.AbiArr[0].Outputs = append(wasmAbi.AbiArr[0].Outputs, utils.OutputsParam{Name:"return", Type:"string"} )
				abi, err = json.Marshal(wasmAbi.AbiArr)
				assert.NoError(t, err)

				lifeVm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, nil, nil)
				assert.NoError(t, err)
				return lifeVm, input, abi
			},
			expect: result{[]int64(nil), "", gerror.ErrInvalidReturnType},
		},
		{
			name: "ErrInvalidRlpFormat",
			given: func() (*exec.VirtualMachine, []byte, []byte){
				input, err := rlp.EncodeToBytes("result")
				assert.NoError(t, err)
				code, abi := vminfo.GetTestData(eventContractName)
				lifeVm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, nil, nil)
				assert.NoError(t, err)
				return lifeVm, input, abi
			},
			expect: result{[]int64{}, "", gerror.ErrInvalidRlpFormat},
		},
		{
			name: "ParseInitFunctionByABIRight",
			given: func() (*exec.VirtualMachine, []byte, []byte) {
				code, abi := vminfo.GetTestData("token-payable")
				lifeVm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, nil, nil)
				assert.NoError(t, err)
				input, err := rlp.EncodeToBytes([]interface{}{"dipc", "dipc", utils.Uint64ToBytes(1000)})
				assert.NoError(t, err)
				return lifeVm,input,abi
			},
			expect: result{[]int64{131072, 131080, 1000}, "void", nil},
		},
	}

	for _, tc := range testCases {
		lifeVm, input, abi := tc.given()
		params, returnType, err := ParseInitFunctionByABI(lifeVm, input, abi)
		if err != nil {
			assert.Equal(t, tc.expect.err.Error(), err.Error())
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expect.returnType, returnType)
			assert.Equal(t, tc.expect.params, params)
		}
	}
}
