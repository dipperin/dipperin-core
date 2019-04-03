// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.


package rpc_interface

import (
	"fmt"
	"runtime"
)

type debugAPI interface {
	Metrics(raw bool) (map[string]interface{}, error)
}

type DipperinDebugApi struct {
	service debugAPI
}

func (api *DipperinDebugApi) Metrics(raw bool) (map[string]interface{}, error) {
	return api.service.Metrics(raw)
}

func (api *DipperinDebugApi) PrintGos() {
	buf := make([]byte, 5 * 1024 * 1024)
	buf = buf[:runtime.Stack(buf, true)]
	fmt.Println(string(buf))
}