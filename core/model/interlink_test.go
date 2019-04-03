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
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

func TestNewInterLink(t *testing.T) {
	block := CreateBlock(1, common.HexToHash("123"), 1)

	// curNum = 1
	mInter := NewInterLink([]common.Hash{common.HexToHash("123")}, block)
	assert.Equal(t, InterLink{block.PreHash(), block.PreHash()}, mInter)

	// curNum = 2
	block = CreateBlock(2, common.HexToHash("123"), 1)
	mInter = NewInterLink(block.GetInterlinks(), block)
	assert.NotNil(t, DeriveSha(mInter))
}

func TestInterLink_RLP(t *testing.T) {
	inter := common.HexToHash("1111222233334444")
	enc, err := rlp.EncodeToBytes([]common.Hash{inter})
	fmt.Println(enc, err)
}

func Test1(t *testing.T) {
	inter := common.HexToHash("1111222233334444")

	ii := InterLink{inter, inter}

	v, err := rlp.EncodeToBytes(ii)
	assert.NoError(t, err)

	assert.Equal(t, true, len(v) > 0)

	var tmpIi InterLink
	assert.NoError(t, rlp.DecodeBytes(v, &tmpIi))

	assert.Equal(t, 2, len(tmpIi))
}

func Test2(t *testing.T) {
	voteA := CreateSignedVote(1, 2, common.HexToHash("0x123456"), VoteMessage)
	block := NewBlock(&Header{Bloom: iblt.NewBloom(DefaultBlockBloomConfig)}, nil, []AbstractVerification{
		voteA,
	})

	inter := common.HexToHash("1111222233334444")
	ii := InterLink{inter, inter}

	block.SetInterLinks(ii)

	b, err := rlp.EncodeToBytes(block)
	assert.NoError(t, err)
	assert.Equal(t, true, len(b) > 0)

	var dBlock Block
	err = rlp.DecodeBytes(b, &dBlock)
	assert.NoError(t, err)

	fmt.Println(dBlock.body.Inters)

	body := dBlock.body
	bb, err := rlp.EncodeToBytes(body)
	assert.NoError(t, err)

	var dBody Body
	err = rlp.DecodeBytes(bb, &dBody)
	assert.NoError(t, err)

	fmt.Println(dBody)
}
