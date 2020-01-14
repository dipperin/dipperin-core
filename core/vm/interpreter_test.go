package vm

import (
	"github.com/dipperin/dipperin-core/core/vm/common"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/util"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
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

	code, abi := test_util.GetTestData(eventContractName)
	lifeVm, err := exec.NewVirtualMachine(code, common.DEFAULT_VM_CONFIG, nil, nil)
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
			expect: result{errInvalidAbi, []int64{}, ""},
		},
		{
			name: "errInputParam",
			given: func() result {
				_, _, err = findParams(lifeVm, abi, callFuncName, []interface{}{})
				return result{err, []int64{}, ""}
			},
			expect: result{errInputAbiNotMatch, []int64{}, ""},
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
	code, abi := test_util.GetTestData(eventContractName)

	type result struct {
		code    []byte
		abi     []byte
		rlpInit []byte
		err     error
	}

	testCases := []struct {
		name   string
		given  func() result
		expect result
	}{
		{
			name: "errEmptyInput",
			given: func() result {
				rlpDataErr := "errData"
				code, abi, rlpInit, err := ParseCreateExtraData([]byte(rlpDataErr))
				return result{
					code:    code,
					abi:     abi,
					rlpInit: rlpInit,
					err:     err,
				}
			},
			expect: result{nil, nil, nil, errInvalidRlpFormat},
		},
		{
			name: "errInsufficientParams",
			given: func() result {
				rlpData, err := rlp.EncodeToBytes([]interface{}{code})
				code, abi, rlpInit, err := ParseCreateExtraData([]byte(rlpData))
				return result{
					code:    code,
					abi:     abi,
					rlpInit: rlpInit,
					err:     err,
				}
			},
			expect: result{nil, nil, nil, errInsufficientParams},
		},
		{
			name: "RightParseCreateExtraData",
			given: func() result {
				rlpData, err := rlp.EncodeToBytes([]interface{}{code, abi})
				assert.NoError(t, err)
				code, abi, rlpInit, err := ParseCreateExtraData([]byte(rlpData))
				return result{
					code:    code,
					abi:     abi,
					rlpInit: rlpInit,
					err:     err,
				}
			},
			expect: result{code, abi, nil, nil},
		},
	}

	for _, tc := range testCases {
		result := tc.given()
		if result.err != nil {
			assert.Equal(t, tc.expect.err.Error(), result.err.Error())
		} else {
			assert.NoError(t, result.err)
			assert.Equal(t, tc.expect.code, result.code)
			assert.Equal(t, tc.expect.abi, result.abi)
			assert.Equal(t, tc.expect.rlpInit, result.rlpInit)
		}
	}

}

func Test_ParseCallExtraDataByABI(t *testing.T) {
	ctrl, _, _ := GetBaseVmInfo(t)
	defer ctrl.Finish()

	code, abi := test_util.GetTestData(eventContractName)
	lifeVm, err := exec.NewVirtualMachine(code, common.DEFAULT_VM_CONFIG, nil, nil)
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
		given  func() *result
		expect result
	}{
		{
			name: "errInvalidRlpFormat",
			given: func() *result {
				input, err := rlp.EncodeToBytes("result")
				assert.NoError(t, err)
				funcName, params, returnType, err := ParseCallExtraDataByABI(lifeVm, input, abi)
				return &result{
					funcName:   funcName,
					params:     params,
					returnType: returnType,
					err:        err,
				}
			},
			expect: result{"", []int64{}, "", errInvalidRlpFormat},
		},
		{
			name: "ParseCallExtraDataByABIRight",
			given: func() *result {
				input, err := rlp.EncodeToBytes([]interface{}{callFuncName, "winner"})
				assert.NoError(t, err)
				funcName, params, returnType, err := ParseCallExtraDataByABI(lifeVm, input, abi)
				return &result{
					funcName:   funcName,
					params:     params,
					returnType: returnType,
					err:        err,
				}
			},
			expect: result{callFuncName, []int64{131072}, "string", errInvalidRlpFormat},
		},
	}

	for _, tc := range testCases {
		res := tc.given()
		if res.err != nil {
			assert.Equal(t, tc.expect.err.Error(), res.err.Error())
		} else {
			assert.NoError(t, res.err)
			assert.Equal(t, tc.expect.returnType, res.returnType)
			assert.Equal(t, tc.expect.params, res.params)
			assert.Equal(t, tc.expect.funcName, res.funcName)
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
		given  func() *result
		expect result
	}{
		{
			name: "errInvalidRlpFormat",
			given: func() *result {
				input, err := rlp.EncodeToBytes("result")
				assert.NoError(t, err)
				code, abi := test_util.GetTestData(eventContractName)
				lifeVm, err := exec.NewVirtualMachine(code, common.DEFAULT_VM_CONFIG, nil, nil)
				assert.NoError(t, err)
				params, returnType, err := ParseInitFunctionByABI(lifeVm, input, abi)
				return &result{
					params:     params,
					returnType: returnType,
					err:        err,
				}
			},
			expect: result{[]int64{}, "", errInvalidRlpFormat},
		},
		{
			name: "ParseInitFunctionByABIRight",
			given: func() *result {
				code, abi := test_util.GetTestData("token-payable")
				lifeVm, err := exec.NewVirtualMachine(code, common.DEFAULT_VM_CONFIG, nil, nil)
				assert.NoError(t, err)
				input, err := rlp.EncodeToBytes([]interface{}{"dipc", "dipc", utils.Uint64ToBytes(1000)})
				assert.NoError(t, err)
				params, returnType, err := ParseInitFunctionByABI(lifeVm, input, abi)
				return &result{
					params:     params,
					returnType: returnType,
					err:        err,
				}
			},
			expect: result{[]int64{131072, 131080, 1000}, "void", nil},
		},
	}

	for _, tc := range testCases {
		res := tc.given()
		if res.err != nil {
			assert.Equal(t, tc.expect.err.Error(), res.err.Error())
		} else {
			assert.NoError(t, res.err)
			assert.Equal(t, tc.expect.returnType, res.returnType)
			assert.Equal(t, tc.expect.params, res.params)
		}
	}
}
