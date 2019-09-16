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

package dipperin

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeNodeService struct {
	err error
}

func (s fakeNodeService) Start() error {
	time.Sleep(time.Millisecond * 100)
	return s.err
}

func (s fakeNodeService) Stop() {
	return
}

func TestNewCsNode(t *testing.T) {
	csNode := NewCsNode([]NodeService{fakeNodeService{}})
	err := csNode.Start()
	assert.NoError(t, err)

	csNode.Stop()

	csNode.AddService(fakeNodeService{err: errors.New("test error")})
	err = csNode.Start()
	assert.Error(t, err)
	csNode.Wait()
}

func TestCsNode_Start(t *testing.T) {
	chokeTimeout = time.Millisecond * 50
	csNode := NewCsNode([]NodeService{})
	err := csNode.Start()
	assert.NoError(t, err)
}
