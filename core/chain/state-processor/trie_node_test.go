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


package state_processor

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewNodeSet(t *testing.T) {
	nodeSet := NewNodeSet()
	key := []byte{1}
	value := []byte{2}

	err := nodeSet.Put(key, value)
	assert.NoError(t, err)

	err = nodeSet.Put(key, value)
	assert.NoError(t, err)

	result, err := nodeSet.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	result, err = nodeSet.Get([]byte{3})
	assert.Error(t, err)
	assert.Nil(t, result)


	has, err := nodeSet.Has(key)
	assert.NoError(t, err)
	assert.True(t, has)

	num := nodeSet.KeyCount()
	assert.Equal(t, 1, num)

	size := nodeSet.DataSize()
	assert.Equal(t, 1, size)

	nodeList := nodeSet.NodeList()
	assert.NotNil(t, nodeList)

	nodeSet.Store(&nodeList)
}

func TestNodeList(t *testing.T) {
	nodeList := NodeList{}
	key := []byte{1}
	value := []byte{2}

	err := nodeList.Put(key, value)
	assert.NoError(t, err)

	num := nodeList.DataSize()
	assert.Equal(t, 1, num)

	nodeList.Store(NewNodeSet())

	nodeSet := nodeList.NodeSet()
	assert.NotNil(t, nodeSet)
}