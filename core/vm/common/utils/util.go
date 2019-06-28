package utils

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"io/ioutil"
	"reflect"
	"strings"
)

var (
	errReturnInvalidRlpFormat    = errors.New("vm_utils: invalid rlp format")
	errReturnInsufficientParams  = errors.New("vm_utils: invalid input. ele must greater than 1")
	errReturnInvalidOutputLength = errors.New("vm_utils: invalid init function outputs length")
	errLengthInputAbiNotMatch    = errors.New("vm_utils: length of input and abi not match")
)

// RLP([code][abi][init params])
func ParseCallContractData(abi []byte, rlpInput []byte) (extraData []byte, err error) {
	// decode rlpInput
	inputPtr := new(interface{})
	//err = rlp.Decode(bytes.NewReader(rlpInput), &inputPtr)
	err = rlp.Decode(bytes.NewReader(rlpInput), &inputPtr)
	if err != nil {
		log.Error("ParseCallContractData#rlp.Decode", "err", err)
		return
	}
	inputRlpList := reflect.ValueOf(inputPtr).Elem().Interface()
	if _, ok := inputRlpList.([]interface{}); !ok {
		return nil, errReturnInvalidRlpFormat
	}

	inRlpList := inputRlpList.([]interface{})
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
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold(v.Name, funcName) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
			break
		}
	}

	var (
		paramStr string
		params []string
	)
	if v, ok := inRlpList[1].([]byte); ok {
		paramStr = string(v)
	}

	if paramStr != "" {
		params = strings.Split(paramStr, ",")
	}

	if len(args) != len(params) {
		log.Error("ParseCallContractData failed", "err", fmt.Sprintf("rlpInput:%v, abi:%v", len(params), len(args)))
		return nil, errLengthInputAbiNotMatch
	}

	rlpParams := []interface{}{
		funcName,
	}

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
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(rlpData), &ptr)
	if err != nil {
		return
	}

	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return nil, errReturnInvalidRlpFormat
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) <= 1 {
		return nil, errReturnInsufficientParams
	}

	// return if no init params
	if len(iRlpList) == 2 {
		log.Info("init function has no params", "len", len(iRlpList))
		return rlpData, nil
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
	for _, v := range abi.AbiArr {
		if strings.EqualFold(v.Name, "init") && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
			if len(v.Outputs) != 0 {
				return nil, errReturnInvalidOutputLength
			}
			break
		}
	}

	var (
		paramStr string
		params []string
	)
	if v, ok := iRlpList[2].([]byte); ok {
		paramStr = string(v)
	}

	if paramStr != "" {
		params = strings.Split(paramStr, ",")
	}

	if len(args) != len(params) {
		log.Error("ParseCreateContractData failed", "err", fmt.Sprintf("input:%v, abi:%v", len(params), len(args)))
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

// get Contract data
func GetExtraData(abiPath, wasmPath string, params []string) (err error, extraData []byte) {
	// GetContractExtraData
	abiBytes, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return
	}
	var wasmAbi WasmAbi
	err = wasmAbi.FromJson(abiBytes)
	if err != nil {
		return
	}
	var args []InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold("init", v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
		}
	}
	//params := []string{"dipp", "DIPP", "100000000"}
	wasmBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return
	}
	rlpParams := []interface{}{
		wasmBytes, abiBytes,
	}
	if len(params) != len(args) {
		return errors.New("params length and args length not equal"), nil
	}
	for i, v := range args {
		bts := params[i]
		re, err := StringConverter(bts, v.Type)
		if err != nil {
			return err, nil
		}
		rlpParams = append(rlpParams, re)
		//inputParams = append(inputParams, re)
	}
	data, err := rlp.EncodeToBytes(rlpParams)
	if err != nil {
		return err, nil
	}
	return err, data
}
