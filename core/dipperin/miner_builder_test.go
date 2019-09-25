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
	"github.com/stretchr/testify/assert"
	"testing"
)

const url = "enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@127.0.0.1:10003"

func TestNewMinerNode(t *testing.T) {
	node, err := NewMinerNode(url, "", 1, aliceAddr.String())
	assert.Error(t, err)
	assert.Nil(t, node)

	node, err = NewMinerNode("", "coinbase", 1, aliceAddr.String())
	assert.Error(t, err)
	assert.Nil(t, node)

	node, err = NewMinerNode(url, "coinbase", 1, aliceAddr.String())
	assert.NoError(t, err)
	assert.NotNil(t, node)
}
