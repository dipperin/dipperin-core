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

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestInvBloom_RLP(t *testing.T) {
	bloom := NewInvBloom(defaultInvBloomConfig)

	KV := generateRandomKV(10000)

	bloom.InsertMap(KV)

	bytes, err := rlp.EncodeToBytes(bloom)

	var rehydrate InvBloom

	if err != nil {
		fmt.Println(err)
	}

	err = rlp.DecodeBytes(bytes, &rehydrate)

	if err != nil {
		fmt.Println(err)
	}

	assert.EqualValues(t, bloom, &rehydrate)
}

func TestBloom_RLP(t *testing.T) {
	b := NewBloom(defaultConfig)
	b.Digest([]byte{1, 2, 3})

	bytes, err := rlp.EncodeToBytes(b)
	assert.NoError(t, err)

	var rehydrate Bloom

	err = rlp.DecodeBytes(bytes, &rehydrate)
	assert.NoError(t, err)

	assert.True(t, rehydrate.IsEqual(b))
}

func TestHybridEstimator_RLP(t *testing.T) {
	e := NewHybridEstimator(NewHybridEstimatorConfig())

	whole := generateRandomKV(200)

	for _, v := range whole {
		e.Encode(v)
	}

	bytes, err := rlp.EncodeToBytes(e)

	assert.NoError(t, err)

	var rehydrate HybridEstimator

	err = rlp.DecodeBytes(bytes, &rehydrate)
	assert.NoError(t, err)

	assert.EqualValues(t, e, &rehydrate)
}

func TestStrataEstimator_DecodeRLP(t *testing.T) {
	e := NewEstimator(NewEstimatorConfig(6))

	bytes, err := rlp.EncodeToBytes(e)

	assert.NoError(t, err)

	var rehydrate StrataEstimator

	err = rlp.DecodeBytes(bytes, &rehydrate)
	assert.NoError(t, err)

	assert.EqualValues(t, e, &rehydrate)
}

func TestHashPool_DecodeRLP(t *testing.T) {
	e := NewHashPool(NewHashPoolConfig(4, 4))

	whole := generateRandomKV(10)

	for _, v := range whole {
		e.Encode(v)
	}

	bytes, err := rlp.EncodeToBytes(e)
	assert.NoError(t, err)

	var rehydrate HashPool

	err = rlp.DecodeBytes(bytes, &rehydrate)
	assert.NoError(t, err)

	assert.EqualValues(t, e, &rehydrate)
}

func TestRLPByteSlice(t *testing.T) {
	b := make(SortDataHash, 4)

	b[0] = []byte{1, 2, 3}
	b[1] = []byte{1, 2, 3}
	b[2] = []byte{1, 2, 3}
	b[3] = []byte{1, 2, 3}
	bytes, err := rlp.EncodeToBytes(b)
	assert.NoError(t, err)

	var a SortDataHash

	err = rlp.DecodeBytes(bytes, &a)
	assert.NoError(t, err)

	assert.EqualValues(t, b, a)
}

func TestGraphene_grapheneRLP(t *testing.T) {
	bloom := NewGraphene(defaultInvBloomConfig, defaultConfig)
	bloomRlp := bloom.grapheneRLP()

	assert.Equal(t, bloomRlp.Bloom, bloom.bloom.BloomRLP())
	assert.Equal(t, bloomRlp.InvBloom, bloom.invBloom.invBloomRLP())
}

func TestGrapheneRLP_graphene(t *testing.T) {
	bloom := NewGraphene(defaultInvBloomConfig, defaultConfig)
	bloomRlp := bloom.grapheneRLP()
	bloomRlp.graphene(bloom)

	assert.Equal(t, bloomRlp.Bloom, bloom.bloom.BloomRLP())
	assert.Equal(t, bloomRlp.InvBloom, bloom.invBloom.invBloomRLP())
}

func TestGraphene_EncodeRLP(t *testing.T) {

}
