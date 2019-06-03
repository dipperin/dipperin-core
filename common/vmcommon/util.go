package vmcommon

import (
	"bytes"
	"errors"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"reflect"
	"strconv"
	"strings"
)

//  RLP([txType][code][abi][init params])
func ParseAndGetRlpData(rlpData []byte, input []byte) (extraData []byte, err error) {

	inputPtr := new(interface{})
	err = rlp.Decode(bytes.NewReader(input), &inputPtr)
	if err != nil {
		return
	}
	inputRlpList := reflect.ValueOf(inputPtr).Elem().Interface()
	if _, ok := inputRlpList.([]interface{}); !ok {
		return nil, errors.New("call contract: invalid input param")
	}
	inRlpList := inputRlpList.([]interface{})
	var funcName string
	if v, ok := inRlpList[0].([]byte); ok {
		funcName = string(v)
	}

	var paramStr string
	if v, ok := inRlpList[1].([]byte); ok {
		paramStr = string(v)
	}

	params := strings.Split(paramStr, ",")

	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(rlpData), &ptr)
	if err != nil {
		return
	}
	rlpList := reflect.ValueOf(ptr).Elem().Interface()

	if _, ok := rlpList.([]interface{}); !ok {
		return nil, errors.New("call contract: invalid rlp format")
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) <= 2 {
		return nil, errors.New("invalid input, ele must greater than 2")
	}
	var (
		abi []byte
	)

	if v, ok := iRlpList[2].([]byte); ok {
		abi = v
	}

	wasmAbi := new(utils.WasmAbi)
	err = wasmAbi.FromJson(abi)
	if err != nil {
		log.Error("ParseAndGetRlpData abi from json", "err", err)
		return nil, errors.New("call contract: invalid abi, encoded fail")
	}

	var args []utils.InputParam
	for _, v := range wasmAbi.AbiArr {
		if strings.EqualFold(funcName, v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
			break
		}
	}

	if len(args) != len(params) {
		return nil, errors.New("interpreter_life: length of input and abi not match")
	}

	rlpParams := []interface{}{
		funcName,
	}

	for i, v := range args {
		bts := params[i]
		switch v.Type {
		case "string":
			rlpParams = append(rlpParams, bts)
		case "int8":
			result, err := strconv.ParseInt(bts, 10, 8)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		case "int16":
			result, err := strconv.ParseInt(bts, 10, 16)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		case "int32", "int":
			result, err := strconv.ParseInt(bts, 10, 32)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		case "int64":

			result, err := strconv.ParseInt(bts, 10, 64)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		case "uint8":
			result, err := strconv.ParseUint(bts, 10, 8)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		case "uint32", "uint":
			result, err := strconv.ParseUint(bts, 10, 32)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		case "uint64":
			result, err := strconv.ParseUint(bts, 10, 64)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		case "bool":
			result, err := strconv.ParseBool(bts)
			if err != nil {
				return nil, errors.New("contract param type is wrong")
			}
			rlpParams = append(rlpParams, result)
		}
	}

	return rlp.EncodeToBytes(rlpParams)
}



