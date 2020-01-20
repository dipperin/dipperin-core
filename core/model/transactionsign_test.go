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
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestMercurySigner_Sender(t *testing.T) {
	tx1, tx2 := createTestTx()

	fs := NewSigner(big.NewInt(1))
	tryAddr, err := fs.GetSender(tx1)
	assert.NoError(t, err)
	assert.Equal(t, tryAddr, aliceAddr)

	fs = NewSigner(big.NewInt(3))
	tryAddr, err = fs.GetSender(tx2)
	assert.NoError(t, err)
	assert.Equal(t, tryAddr, bobAddr)
	assert.Equal(t, DipperinSigner{chainId: new(big.Int)}, NewSigner(nil))
}

func TestMercurySigner_Equal(t *testing.T) {
	fs1 := NewSigner(big.NewInt(1))
	fs2 := NewSigner(big.NewInt(4))
	fs3 := NewSigner(big.NewInt(1))
	assert.Equal(t, fs1.Equal(fs2), false)
	assert.Equal(t, fs1.Equal(fs3), true)
}

func TestMercurySigner_GetSignHash(t *testing.T) {
	tx1, tx2 := createTestTx()
	fs1 := NewSigner(big.NewInt(1))
	fs2 := NewSigner(big.NewInt(3))
	tryHash1, err1 := rlpHash([]interface{}{tx1.data, fs1.chainId})
	tryHash2, err2 := rlpHash([]interface{}{tx2.data, fs2.chainId})
	assert.NoError(t, err1)
	assert.NoError(t, err2)

	getHash1, err3 := fs1.GetSignHash(tx1)
	assert.NoError(t, err3)
	getHash2, err4 := fs2.GetSignHash(tx2)
	assert.NoError(t, err4)
	assert.Equal(t, tryHash1, getHash1)
	assert.Equal(t, tryHash2, getHash2)

	tryHash3, err6 := rlpHash([]interface{}{tx1.data, fs2.chainId})
	assert.NoError(t, err6)
	getHash3, err5 := fs2.GetSignHash(tx1)
	assert.NoError(t, err5)
	assert.Equal(t, tryHash3, getHash3)
	assert.NotEqual(t, getHash1, getHash3)
}

func TestMercurySigner_SignatureValues(t *testing.T) {
	tx1, _ := createTestTx()
	key1, _ := CreateKey()
	fs1 := NewSigner(big.NewInt(1))
	getR := tx1.wit.R
	getS := tx1.wit.S
	getV := tx1.wit.V
	sigHash, err1 := fs1.GetSignHash(tx1)
	assert.NoError(t, err1)

	sig, err2 := crypto.Sign(sigHash[:], key1)
	assert.NoError(t, err2)
	r, s, v, err3 := fs1.SignatureValues(tx1, sig)
	assert.NoError(t, err3)
	assert.Equal(t, r, getR)
	assert.Equal(t, s, getS)
	assert.Equal(t, v, getV)
	assert.Panics(t, func() {
		fs1.SignatureValues(tx1, []byte{123})
	})
}

func TestDeriveChainId(t *testing.T) {
	tx1, tx2 := createTestTx()
	assert.Equal(t, deriveChainId(tx1.wit.V), big.NewInt(1))
	assert.Equal(t, deriveChainId(tx2.wit.V), big.NewInt(3))
}

func TestTransaction_SignTx(t *testing.T) {
	key1, _ := CreateKey()
	fs1 := NewSigner(big.NewInt(1))
	tx := NewTransaction(10, bobAddr, big.NewInt(10000), TestGasPrice, TestGasLimit, []byte{})

	result, err := tx.SignTx(key1, fs1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestNewMercurySigner(t *testing.T) {
	fs1 := NewSigner(big.NewInt(1))
	assert.NotNil(t, fs1)
}

func TestMercurySigner_GetSender(t *testing.T) {
	tx := CreateSignedTx(0, big.NewInt(10000))
	ms := NewSigner(big.NewInt(1))

	sender, err := ms.GetSender(tx)
	assert.Equal(t, aliceAddr, sender)
	assert.NoError(t, err)
}

func TestMercurySigner_GetSenderPublicKey(t *testing.T) {
	tx := CreateSignedTx(0, big.NewInt(10000))
	ms := NewSigner(big.NewInt(1))

	key, _ := CreateKey()
	pubKey, err2 := ms.GetSenderPublicKey(tx)

	assert.Equal(t, &key.PublicKey, pubKey)
	assert.NoError(t, err2)
}

func TestMakeSigner(t *testing.T) {
	config := chainconfig.GetChainConfig()
	result := DipperinSigner{chainId: config.ChainId}
	signer := MakeSigner(config, 10)
	assert.Equal(t, result, signer)
}
