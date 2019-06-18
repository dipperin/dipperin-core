package utils

import (
	"encoding/json"
	"fmt"
	"bytes"
	"reflect"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/dipperin/dipperin-core/common/g-error"
	"errors"
)

var errReturnInputAbiNotMatch = errors.New("ABI: length of input and abi not match")

type WasmAbi struct {
	AbiArr []AbiStruct `json:"abiArr"`
}

type AbiStruct struct {
	Name     string         `json:"name"`
	Inputs   []InputParam   `json:"inputs"`
	Outputs  []OutputsParam `json:"outputs"`
	Constant string         `json:"constant"`
	Type     string         `json:"type"`
}

type InputParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type OutputsParam struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (abi *WasmAbi) FromJson(body []byte) error {
	if body == nil {
		return fmt.Errorf("invalid param. %v", body)
	}
	err := json.Unmarshal(body, &abi.AbiArr)
	return err
}

func ConvertInputs(src []byte, inputParam []InputParam) ([]byte, error) {
	ptr := new(interface{})
	err := rlp.Decode(bytes.NewReader(src), &ptr)
	if err != nil {
		return nil, err
	}

	rlpList := reflect.ValueOf(ptr).Elem().Interface()
	if _, ok := rlpList.([]interface{}); !ok {
		return nil, g_error.ErrReturnInvalidRlpFormat
	}

	inputList := rlpList.([]interface{})
	if len(inputList) != len(inputParam) {
		return nil, errReturnInputAbiNotMatch
	}

	var data []byte
	for i, v := range inputParam {
		input := inputList[i].([]byte)
		convert := BytesConverter(input, v.Type)
		result := fmt.Sprintf("%v,", convert)
		data = append(data, []byte(result)...)
	}
	return data, nil
}
