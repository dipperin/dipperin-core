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
	"time"

	"github.com/dipperin/dipperin-core/third-party/rpc"
)

func TestService(t *testing.T) {
	s := &Service{}
	s.AddApis([]rpc.API{
		{Namespace: "test", Version: "0.0.1", Service: &fakeAPI{}, Public: true},
	})
	assert.Error(t, s.Start())

	s.apis = []rpc.API{
		{Namespace: "test", Version: "0.0.1", Service: &FakeAPI{}, Public: true},
	}
	assert.NoError(t, s.Start())

	s.httpEndpoint = ":123"
	//assert.Error(t, s.Start())

	_ = s.Start()

	s.httpEndpoint = ""
	s.wsEndpoint = ":123"
	assert.Error(t, s.Start())

	s.httpEndpoint = ":15214"
	s.wsEndpoint = ":15213"
	assert.NoError(t, s.Start())
	time.Sleep(time.Millisecond)
	s.Stop()
}

type FakeAPI struct{}

func (f *FakeAPI) GetNum() uint64 {
	return 1
}

type fakeAPI struct{}

func (f *fakeAPI) GetNum() uint64 {
	return 1
}
