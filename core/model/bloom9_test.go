// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package model

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateBloom(t *testing.T) {
	topics := []common.Hash{common.HexToHash("topic")}
	l := &Log{Topics: topics}
	receipt1 := NewReceipt([]byte{}, false, uint64(100), []*Log{l})
	bloom := CreateBloom(Receipts{receipt1})
	assert.Equal(t, true, BloomLookup(bloom, topics[0]))
	assert.Equal(t, false, BloomLookup(bloom, common.HexToHash("bloom")))

	bloom.Add(big.NewInt(100))
	assert.Equal(t, true, bloom.Test(big.NewInt(100)))
	assert.Equal(t, true, bloom.TestBytes(big.NewInt(100).Bytes()))
	assert.Equal(t, bloom.Hex(), hexutil.Encode(bloom.Bytes()))
	assert.Panics(t, func() {
		d := make([]byte, BloomByteLength+1)
		bloom.SetBytes(d)
	})
}

func TestBloom_MarshalText(t *testing.T) {
	topics := []common.Hash{common.HexToHash("topic")}
	l := &Log{Topics: topics}
	receipt1 := NewReceipt([]byte{}, false, uint64(100), []*Log{l})
	bloom := CreateBloom(Receipts{receipt1})

	// MarshalJSON
	enc, err1 := bloom.MarshalText()
	assert.NoError(t, err1)

	// UnmarshalJSON
	b1get := &Bloom{}
	err2 := b1get.UnmarshalText(enc)
	assert.NoError(t, err2)
	assert.EqualValues(t, &bloom, b1get)
}

func bloom9t(b []byte) *big.Int {
	b = crypto.Keccak256(b)
	fmt.Println("crypto : ", common.Bytes2Hex(b))

	r := new(big.Int)
	fmt.Println("=====r", r)

	for i := 0; i < 6; i += 2 {
		t := big.NewInt(1)
		b := (uint(b[i+1]) + (uint(b[i]) << 8)) & 2047
		fmt.Println(i, "=====b", b)
		t = t.Lsh(t, b)
		fmt.Println(i, "=====t", t)
		r.Or(r, t)
		fmt.Println(i, "=====r", r)

	}

	return r
}

func TestBloom9(t *testing.T) {
	fmt.Println(bloom9t([]byte("t")))
}
