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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeDipperinVenusApi(t *testing.T) {
	assert.NotNil(t, MakeDipperinVenusApi(nil))
}

func TestMakeDipperinDebugApi(t *testing.T) {
	assert.NotNil(t, MakeDipperinDebugApi(nil))
}

func TestMakeDipperinP2PApi(t *testing.T) {
	assert.NotNil(t, MakeDipperinP2PApi(nil))
}

func TestMakeDipperExternalApi(t *testing.T) {
	assert.NotNil(t, MakeDipperExternalApi(nil))
}

type nodeConfigMock struct{}

func (n nodeConfigMock) IpcEndpoint() string {
	return "IpcEndpoint"
}

func (n nodeConfigMock) HttpEndpoint() string {
	return "HttpEndpoint"
}

func (n nodeConfigMock) WsEndpoint() string {
	return "WsEndpoint"
}

func TestMakeRpcService(t *testing.T) {
	cnf := nodeConfigMock{}
	s := MakeRpcService(cnf, nil, nil)
	assert.NotNil(t, s)
	assert.Equal(t, cnf.IpcEndpoint(), s.ipcEndpoint)
	assert.Equal(t, cnf.HttpEndpoint(), s.httpEndpoint)
	assert.Equal(t, cnf.WsEndpoint(), s.wsEndpoint)
}
