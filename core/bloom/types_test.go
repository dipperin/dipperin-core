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


package iblt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyXOR(t *testing.T) {
	key1 := make(Data, 4)
	key2 := make(Data, 4)
	key1[0] = 15
	key2[1] = 1

	sum := make(Data, 4)
	sum.Xor(key1, key2)
	assert.Equal(t, byte(15), sum[0])
	assert.Equal(t, byte(1), sum[1])

	sum.Xor(sum, key1)
	assert.Equal(t, byte(0), sum[0])
	assert.Equal(t, byte(1), sum[1])
}

func TestKHash(t *testing.T) {
	k := uint(7)

	hashes := kHash([]byte{}, k)

	assert.EqualValues(t, k, len(hashes))
}

func TestData_Less(t *testing.T) {
	key1 := make(Data, 4)
	key2 := make(Data, 4)
	key1[2] = 7
	key1[3] = 14
	key2[1] = 5
	key2[2] = 4
	assert.EqualValues(t, key1.Less(key2), true)
}

func TestDistinct(t *testing.T) {
	s := make([]uint, 4)
	v := make([]uint, 4)
	n := uint(5)
	s[0] = 8
	s[3] = 1
	distinct(s, n)
	v[0] = 8
	v[2] = 1
	v[3] = 2
	assert.EqualValues(t, s, v)
}

func TestKHash2(t *testing.T) {
	k := uint(10)
	hashes := kHash([]byte{}, k)
	assert.NotEqual(t, k, len(hashes))
}

func TestData_SetBytes(t *testing.T) {
	h := make(Data, 4)

	h.SetBytes([]byte("test"))

	assert.Equal(t, h.Bytes(), []byte("test"))
}

func TestDataHash_SetBytes(t *testing.T) {
	h := make(DataHash, 4)

	h.SetBytes([]byte("test"))

	assert.Equal(t, h.Bytes(), []byte("test"))
}

func TestData_Big(t *testing.T) {
	h := make(Data, 4)

	h.SetBytes([]byte("test"))

	assert.Equal(t, h.Big().String(), "1952805748")
}

func TestDataHash_Big(t *testing.T) {
	h := make(DataHash, 4)

	h.SetBytes([]byte("test"))

	assert.Equal(t, h.Big().String(), "1952805748")
}

func TestDataHash_Lsh(t *testing.T) {
	h := make(DataHash, 4)

	h.SetBytes([]byte("test"))

	assert.Equal(t, h.Lsh(0), DataHash("test"))
}

