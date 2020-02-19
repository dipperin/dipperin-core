// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package common contains various helper functions.
package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToHex(t *testing.T) {
	assert.Equal(t, []byte{0x1, 0x23}, FromHex("0x123"))
	assert.Equal(t, "0x0", ToHex([]byte{}))

	assert.Nil(t, CopyBytes(nil))
	assert.Equal(t, []byte{1}, CopyBytes([]byte{1}))

	assert.True(t, hasHexPrefix("0x123"))
	assert.False(t, isHexCharacter('n'))

	assert.False(t, isHex("0x1"))
	assert.False(t, isHex("0x1q"))
	assert.False(t, isHex("0x1f"))
	assert.True(t, isHex("1f"))

	assert.Equal(t, []byte{0}, Hex2BytesFixed("0x123", 1))
	assert.Equal(t, []byte{0x12}, Hex2BytesFixed("123", 1))
	assert.Equal(t, []byte{0x35}, Hex2BytesFixed("12356", 1))

	assert.Equal(t, []byte{0x1, 0x2, 0x3}, RightPadBytes([]byte{1, 2, 3}, 1))
	assert.Equal(t, []byte{0x1, 0x2, 0x3, 0, 0}, RightPadBytes([]byte{1, 2, 3}, 5))

	assert.Equal(t, []byte{0x1, 0x2, 0x3}, LeftPadBytes([]byte{1, 2, 3}, 1))
	assert.Equal(t, []byte{0x0, 0x0, 0x1, 0x2, 0x3}, LeftPadBytes([]byte{1, 2, 3}, 5))
}
