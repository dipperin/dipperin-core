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
	"math/big"
	"testing"
)

func TestMakeDefaultBlockDecoder(t *testing.T) {
	bd := MakeDefaultBlockDecoder()
	assert.NotNil(t, bd)
}

func Test_defaultBlockDecoder_DecodeRlpBlockFromHeaderAndBodyBytes(t *testing.T) {
	bd := MakeDefaultBlockDecoder()
	h := NewHeader(1, 100, common.HexToHash("001010001010"), common.HexToHash("10111101011"), common.HexToDiff("1111ffff"), big.NewInt(10100), common.HexToAddress("111111000000"), common.BlockNonceFromInt(100))
	header, err := rlp.EncodeToBytes(h)
	assert.NoError(t, err)
	body, err := rlp.EncodeToBytes(Body{})
	assert.NoError(t, err)

	_, err = bd.DecodeRlpBlockFromHeaderAndBodyBytes(header, body)
	assert.NoError(t, err)
	_, err = bd.DecodeRlpBlockFromHeaderAndBodyBytes([]byte{123}, body)
	assert.Equal(t, "rlp: expected input list for model.Header", err.Error())
	_, err = bd.DecodeRlpBlockFromHeaderAndBodyBytes(header, []byte{123})
	assert.Equal(t, "rlp: expected input list for model.PBFTBody", err.Error())
}

func Test_defaultBlockDecoder_DecodeRlpHeaderFromBytes(t *testing.T) {
	bd := MakeDefaultBlockDecoder()
	h := NewHeader(1, 100, common.HexToHash("001010001010"), common.HexToHash("10111101011"), common.HexToDiff("1111ffff"), big.NewInt(10100), common.HexToAddress("111111000000"), common.BlockNonceFromInt(100))
	header, err := rlp.EncodeToBytes(h)
	assert.NoError(t, err)
	_, err2 := bd.DecodeRlpHeaderFromBytes(header)
	assert.NoError(t, err2)
	_, err = bd.DecodeRlpHeaderFromBytes([]byte{123})
	assert.Equal(t, "rlp: expected input list for model.Header", err.Error())
}

func Test_defaultBlockDecoder_DecodeRlpBodyFromBytes(t *testing.T) {
	bd := MakeDefaultBlockDecoder()
	body, err := rlp.EncodeToBytes(Body{})
	assert.NoError(t, err)
	_, err = bd.DecodeRlpBodyFromBytes(body)
	assert.NoError(t, err)
	_, err = bd.DecodeRlpBodyFromBytes([]byte{123})
	assert.Equal(t, "rlp: expected input list for model.PBFTBody", err.Error())
}

func Test_defaultBlockDecoder_DecodeRlpBlockFromBytes(t *testing.T) {
	bd := MakeDefaultBlockDecoder()
	b := CreateBlock(1, common.HexToHash("123"), 10)
	block, err := rlp.EncodeToBytes(b)
	assert.NoError(t, err)
	_, err = bd.DecodeRlpBlockFromBytes(block)
	assert.NoError(t, err)
	_, err = bd.DecodeRlpBlockFromBytes([]byte{123})
	assert.Equal(t, "rlp: expected input list for model.blockForRlp", err.Error())
}

func Test_defaultBlockDecoder_DecodeRlpTransactionFromBytes(t *testing.T) {
	bd := MakeDefaultBlockDecoder()
	testTx := CreateSignedTx(0, big.NewInt(10000))
	tx, err := rlp.EncodeToBytes(testTx)
	assert.NoError(t, err)
	_, err = bd.DecodeRlpTransactionFromBytes(tx)
	assert.NoError(t, err)
	_, err = bd.DecodeRlpTransactionFromBytes([]byte{123})
	assert.Equal(t, "rlp: expected input list for model.TransactionRLP", err.Error())
}
