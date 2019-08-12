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
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/third-party/rpc"
)

func MakeDipperinVenusApi(service *service.VenusFullChainService) *DipperinVenusApi {
	return &DipperinVenusApi{service: service}
}

func MakeDipperinDebugApi(service debugAPI) *DipperinDebugApi {
	return &DipperinDebugApi{service: service}
}

func MakeDipperinP2PApi(service P2PAPI) *DipperinP2PApi {
	return &DipperinP2PApi{service: service}
}

type nodeConf interface {
	IpcEndpoint() string
	HttpEndpoint() string
	WsEndpoint() string
}

func MakeRpcService(conf nodeConf, apis []rpc.API, allowHosts []string) *Service {
	return &Service{
		ipcEndpoint:  conf.IpcEndpoint(),
		httpEndpoint: conf.HttpEndpoint(),
		wsEndpoint:   conf.WsEndpoint(),
		apis:         apis,
		allowHosts:   allowHosts,
	}
}
