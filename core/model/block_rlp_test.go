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
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlock_RLP(t *testing.T) {
	SetBlockRlpHandler(&PBFTBlockRlpHandler{})
	block := CreateBlock(1, common.HexToHash("123"), 10)
	b, err := rlp.EncodeToBytes(block)
	assert.NoError(t, err)

	var dBlock Block
	err = rlp.DecodeBytes(b, &dBlock)
	assert.NoError(t, err)
	assert.Equal(t, dBlock.Hash(), block.Hash())
	assert.Equal(t, dBlock.header.Nonce, block.header.Nonce)
	assert.Equal(t, dBlock.header.Number, block.header.Number)
	assert.Equal(t, dBlock.header.PreHash, block.header.PreHash)
	assert.Equal(t, dBlock.header.TransactionRoot, block.header.TransactionRoot)
	assert.Equal(t, dBlock.body.Txs[0].CalTxId(), block.body.Txs[0].CalTxId())

	// test body rlp
	body, err := rlp.EncodeToBytes(block.body)
	assert.NoError(t, err)

	var dBody Body
	err = rlp.DecodeBytes(body, &dBody)
	assert.NoError(t, err)
	assert.Equal(t, block.body.Txs[0].CalTxId(), dBody.Txs[0].CalTxId())
}

func TestSetBlockRlpHandler(t *testing.T) {
	SetBlockRlpHandler(blockRlpHandler)
}
