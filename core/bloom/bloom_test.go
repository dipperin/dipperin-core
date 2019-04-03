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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultConfig = NewBloomConfig(8, 4)

func TestBloom_Digest(t *testing.T) {
	b := NewBloom(defaultConfig)
	b.Digest([]byte{1, 2, 3})
	b.Digest([]byte{4, 5, 6})
	//bloomStats(b, t)
}

func TestBloom_Config(t *testing.T) {
	c := DeriveBloomConfig(32000)
	b := NewBloom(c)
	assert.EqualValues(t, c.BloomByteLength, len(b.bloom))
}

func TestBloom_LookUp(t *testing.T) {
	b := NewBloom(defaultConfig)

	b.Digest([]byte{1, 2, 3})
	b.Digest([]byte{4, 5, 6})

	assert.True(t, b.LookUp([]byte{1, 2, 3}))
	assert.False(t, b.LookUp([]byte{1, 4, 5}))
}

// bloomStats counts how many zeros and ones in a bloom.
func bloomStats(b *Bloom, t *testing.T) {
	count := uint(0)

	for i := uint(0); i < b.config.BloomByteLength; i++ {
		sub := b.bloom[i]
		for sub != 0 {
			sub = sub & (sub - 1)
			count++
		}
	}

	fmt.Printf("Bloom Counts:\n0: %d, 1: %d\n", b.config.BloomByteLength*8-count, count)
	assert.True(t, count < b.config.BloomBits)
}

func TestBloom_LookUp2(t *testing.T) {
	m := generateRandomKV(28000)

	bloom := NewGraphene(defaultInvBloomConfig, defaultConfig)
	bloom.InsertMap(m)

	cnt := 0
	testSamples := 5000

	mm := generateRandomKV(testSamples)

	//bloomStats(bloom.bloom,t)
	for _, i := range mm {
		if bloom.LookUp(i) {
			cnt++
		}
	}

	assert.True(t, cnt < int(testSamples*20/1000))

	//bloomStats(bloom.bloom, t)
}

func TestBloom_SetAt(t *testing.T) {
	b := NewBloom(defaultConfig)

	idx := uint32(1)
	b.SetAt(idx)
	assert.True(t, b.LookAt(idx))

	idx = 7
	b.SetAt(idx)
	assert.True(t, b.LookAt(idx))

	idx = 18
	b.SetAt(idx)
	assert.True(t, b.LookAt(idx))

	assert.False(t, b.LookAt(19))
	assert.False(t, b.LookAt(257))
}

func TestBloom_SetBytes(t *testing.T) {
	b := NewBloom(defaultConfig)

	b.SetBytes([]byte("test"))

	assert.Equal(t, b.Hex(), "0x0000000000000000000000000000000000000000000000000000000074657374")

	b2 := NewBloom(defaultConfig)

	b2.SetBytes([]byte("test111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"))

	assert.Equal(t, b2.Hex(), "0x3131313131313131313131313131313131313131313131313131313131313131")
}

func TestBloom_IsEqual(t *testing.T) {
	b := NewBloom(defaultConfig)

	b.SetBytes([]byte("test"))

	b2 := NewBloom(defaultConfig)

	b2.SetBytes([]byte("test"))

	assert.True(t, b.IsEqual(b2))
}

func TestBloom_Or(t *testing.T) {
	b := NewBloom(defaultConfig)

	b.SetBytes([]byte("test"))

	b2 := NewBloom(defaultConfig)

	b2.SetBytes([]byte("test"))

	b.Or(b2, b2)

	assert.Equal(t, b.Hex(), "0x0000000000000000000000000000000000000000000000000000000074657374")
}

func TestBloom_Big(t *testing.T) {
	b := NewBloom(defaultConfig)

	b.SetBytes([]byte("test"))

	assert.Equal(t, b.Big().String(), "1952805748")
}
