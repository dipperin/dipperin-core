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
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestService_GetInProcHandler(t *testing.T) {
	expect := &rpc.Server{}

	s := Service{inprocHandler: expect}
	assert.Equal(t, expect, s.GetInProcHandler())
}

func TestService_AddApis(t *testing.T) {
	expect1 := 1
	expect2 := rpc.API{Namespace: "test-api", Version: "0.0.1"}

	s := Service{}
	s.AddApis([]rpc.API{expect2})
	if assert.Equal(t, expect1, len(s.apis)) {
		api := s.apis[0]
		assert.Equal(t, expect2.Namespace, api.Namespace)
		assert.Equal(t, expect2.Version, api.Version)
	}
}

func TestService_Start(t *testing.T) {
	s := Service{
		httpEndpoint: "127.0.0.1:32001",
		wsEndpoint:   "127.0.0.1:32002",
	}
	if assert.NoError(t, s.Start()) {
		time.Sleep(500 * time.Millisecond)
		s.Stop()
	}
}
