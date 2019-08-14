package utils

import (
	"encoding/json"
)

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
		return errEmptyInput
	}
	err := json.Unmarshal(body, &abi.AbiArr)
	return err
}
