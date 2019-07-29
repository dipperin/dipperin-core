package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"reflect"
	"strings"
)

var (
	errEmptyInput             = errors.New("vm_utils: empty input")
	errInvalidRlpFormat       = errors.New("vm_utils: invalid rlp format")
	errInsufficientParams     = errors.New("vm_utils: invalid input params")
	errInvalidOutputLength    = errors.New("vm_utils: invalid init function outputs length")
	errLengthInputAbiNotMatch = errors.New("vm_utils: length of input and abi not match")
	errFuncNameNotFound       = errors.New("vm_utils: function name not found")
)

// RLP([funName][params])
func ParseCallContractData(abi []byte, rlpInput []byte) (extraData []byte, err error) {
	if rlpInput == nil || len(rlpInput) == 0 {
		return nil, errEmptyInput
	}

	// decode rlpInput
	inputPtr := new(interface{})
	rlp.Decode(bytes.NewReader(rlpInput), &inputPtr)
	inputRlpList := reflect.ValueOf(inputPtr).Elem().Interface()
	if _, ok := inputRlpList.([]interface{}); !ok {
		return nil, errInvalidRlpFormat
	}

	inRlpList := inputRlpList.([]interface{})
	if len(inRlpList) < 1 || len(inRlpList) > 2 {
		err = errInsufficientParams
		return
	}

	var funcName string
	if v, ok := inRlpList[0].([]byte); ok {
		funcName = string(v)
	}

	wasmAbi := new(WasmAbi)
	err = wasmAbi.FromJson(abi)
	if err != nil {
		log.Error("ParseCallContractData#wasmAbi.FromJson", "err", err)
		return nil, err
	}

	var args []InputParam
	found := false
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold(v.Name, funcName) && strings.EqualFold(v.Type, "function") {
			found = true
			args = v.Inputs
			break
		}
	}

	if !found {
		log.Error("ParseCallContractData failed", "err", errFuncNameNotFound, "funcName", funcName)
		err = errFuncNameNotFound
		return
	}

	var (
		paramStr string
		params   []string
	)

	// if function has params or not
	if len(args) == 0 && len(inRlpList) == 1 {
		return rlpInput, nil
	} else {
		if len(inRlpList) == 1 {
			log.Error("ParseCallContractData failed", "err", fmt.Sprintf("rlpInput:%v, abi:%v", len(params), len(args)))
			return nil, errLengthInputAbiNotMatch
		}

		if v, ok := inRlpList[1].([]byte); ok {
			paramStr = string(v)
		}

		if paramStr != "" {
			params = strings.Split(paramStr, ",")
		}
	}

	if len(args) != len(params) {
		log.Error("ParseCallContractData failed", "err", fmt.Sprintf("rlpInput:%v, abi:%v", len(params), len(args)))
		return nil, errLengthInputAbiNotMatch
	}

	rlpParams := []interface{}{funcName}
	for i, v := range args {
		bts := params[i]
		result, innerErr := StringConverter(bts, v.Type)
		if innerErr != nil {
			return nil, innerErr
		}
		rlpParams = append(rlpParams, result)
	}
	return rlp.EncodeToBytes(rlpParams)
}

// RLP([code][abi][init params])
func ParseCreateContractData(rlpData []byte) (extraData []byte, err error) {
	if rlpData == nil || len(rlpData) == 0 {
		return nil, errEmptyInput
	}

	ptr := new(interface{})
	rlp.Decode(bytes.NewReader(rlpData), &ptr)
	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return nil, errInvalidRlpFormat
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) < 2 {
		return nil, errInsufficientParams
	}

	var wasmBytes []byte
	if v, ok := iRlpList[0].([]byte); ok {
		wasmBytes = v
	}

	var abiBytes []byte
	if v, ok := iRlpList[1].([]byte); ok {
		abiBytes = v
	}

	var abi WasmAbi
	err = abi.FromJson(abiBytes)
	if err != nil {
		return nil, err
	}

	var args []InputParam
	found := false
	for _, v := range abi.AbiArr {
		if strings.EqualFold(v.Name, "init") && strings.EqualFold(v.Type, "function") {
			found = true
			args = v.Inputs
			if len(v.Outputs) != 0 {
				return nil, errInvalidOutputLength
			}
			break
		}
	}

	if !found {
		log.Error("ParseCreateContractData failed", "err", errFuncNameNotFound, "funcName", "init")
		err = errFuncNameNotFound
		return
	}

	var (
		paramStr string
		params   []string
	)
	// if function has params or not
	if len(args) == 0 && len(iRlpList) == 2 {
		return rlpData, nil
	} else {
		if len(iRlpList) == 2 {
			log.Error("ParseCallContractData failed", "err", fmt.Sprintf("rlpInput:%v, abi:%v", len(params), len(args)))
			return nil, errLengthInputAbiNotMatch
		}

		if v, ok := iRlpList[2].([]byte); ok {
			paramStr = string(v)
		}

		if paramStr != "" {
			params = strings.Split(paramStr, ",")
		}
	}

	if len(args) != len(params) {
		log.Error("ParseCallContractData failed", "err", fmt.Sprintf("rlpInput:%v, abi:%v", len(params), len(args)))
		return nil, errLengthInputAbiNotMatch
	}

	rlpParams := []interface{}{
		wasmBytes, abiBytes,
	}

	for i, v := range args {
		bts := params[i]
		re, innerErr := StringConverter(bts, v.Type)
		if innerErr != nil {
			return re, innerErr
		}
		rlpParams = append(rlpParams, re)
	}
	return rlp.EncodeToBytes(rlpParams)
}

func ConvertInputs(src []byte, abiInput []InputParam) ([]byte, error) {
	if src == nil || len(src) == 0 {
		log.Error("ConvertInputs failed", "err", errEmptyInput)
		return nil, errEmptyInput
	}

	ptr := new(interface{})
	rlp.Decode(bytes.NewReader(src), &ptr)
	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return nil, errInvalidRlpFormat
	}

	inputList := rlpList.([]interface{})
	if len(inputList) != len(abiInput) {
		log.Error("ConvertInputs failed", "length", fmt.Sprintf("input:%v, abi:%v", len(inputList), len(abiInput)))
		return nil, errLengthInputAbiNotMatch
	}

	var data []byte
	for i, v := range abiInput {
		input := inputList[i].([]byte)
		convert := Align32BytesConverter(input, v.Type)
		result := fmt.Sprintf("%v,", convert)
		data = append(data, []byte(result)...)
	}
	return data, nil
}
