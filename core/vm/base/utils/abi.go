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
	"encoding/json"
	"github.com/dipperin/dipperin-core/common/gerror"
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
		return gerror.ErrEmptyInput
	}
	err := json.Unmarshal(body, &abi.AbiArr)
	return err
}
