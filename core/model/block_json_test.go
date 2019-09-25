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

package model

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeader_JSON(t *testing.T) {
	b1 := CreateBlock(0, common.Hash{}, 1)
	h1 := b1.header

	// MarshalJSON
	enc, err1 := h1.MarshalJSON()
	assert.NoError(t, err1)

	// UnmarshalJSON
	h1get := &Header{}
	err2 := h1get.UnmarshalJSON(enc)
	assert.NoError(t, err2)
	assert.EqualValues(t, h1, h1get)
	assert.Equal(t, h1.Hash(), h1get.Hash())
}

func TestBlock_JSON(t *testing.T) {
	b1 := CreateBlock(0, common.Hash{}, 1)

	// MarshalJSON
	enc, err1 := b1.MarshalJSON()
	assert.NoError(t, err1)

	// UnmarshalJSON
	b1get := &Block{}
	err2 := b1get.UnmarshalJSON(enc)
	assert.NoError(t, err2)
	assert.Equal(t, b1.Hash(), b1get.Hash())
}

func TestBody_UnmarshalJSON(t *testing.T) {
	b := CreateBlock(0, common.Hash{}, 1)
	enc, err := b.MarshalJSON()
	assert.NoError(t, err)
	err = b.body.UnmarshalJSON(enc)
	assert.NoError(t, err)
}

func TestPBFTBlockJsonHandler_DecodeBody(t *testing.T) {
	pb := PBFTBlockJsonHandler{}
	b := CreateBlock(0, common.Hash{}, 1)
	enc, err := b.MarshalJSON()
	assert.NoError(t, err)
	err = pb.DecodeBody(b.body, enc)
	assert.NoError(t, err)
}

func TestSetBlockJsonHandler(t *testing.T) {
	bjh := blockJsonHandler
	SetBlockJsonHandler(bjh)
}
