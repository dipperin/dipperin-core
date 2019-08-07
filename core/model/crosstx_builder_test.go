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
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestGetLockAddress(t *testing.T) {
	buf := cs_crypto.GetLockAddress(aliceAddr, bobAddr)
	fact1 := buf[2:]
	fact2 := buf[:2]
	enc, _ := rlpHash([]interface{}{aliceAddr, bobAddr})
	exp1 := enc[12:]
	var exp2 = []byte{0, 1}
	assert.Equal(t, exp1, fact1)
	assert.Equal(t, exp2, fact2)
}

func TestCreateRawLockTx(t *testing.T) {
	key1, _ := CreateKey()
	fs := NewSigner(big.NewInt(1))

	hashKey := []byte("123")
	hashLock := cs_crypto.Keccak256Hash(hashKey)
	tx := CreateRawLockTx(1, hashLock, big.NewInt(34564), big.NewInt(10000), big.NewInt(1), model.TxGas, aliceAddr, bobAddr)
	tx.SignTx(key1, fs)

	testAlice, _ := fs.GetSender(tx)
	//todo change this to a digest method to parse the extradata
	testBob := common.BytesToAddress(tx.data.ExtraData)
	assert.Equal(t, cs_crypto.GetLockAddress(testAlice, testBob), *tx.data.Recipient)
}

func TestCreateRawRefundTx(t *testing.T) {
	key1, _ := CreateKey()
	fs := NewSigner(big.NewInt(1))
	tx := CreateRawRefundTx(1, big.NewInt(10000), big.NewInt(1), model.TxGas, aliceAddr, bobAddr)
	tx.SignTx(key1, fs)

	testAlice, _ := fs.GetSender(tx)
	testBob := common.BytesToAddress(tx.data.ExtraData)
	assert.Equal(t, cs_crypto.GetLockAddress(testAlice, testBob), *tx.data.Recipient)
}

func TestCreateRawClaimTx(t *testing.T) {
	_, key2 := CreateKey()
	fs := NewSigner(big.NewInt(1))

	hashKey := []byte("123")
	hashLock, _ := rlpHash(hashKey)
	tx := CreateRawClaimTx(1, hashKey, big.NewInt(10000), big.NewInt(1), model.TxGas, aliceAddr, bobAddr)
	tx.SignTx(key2, fs)

	testBob, _ := fs.GetSender(tx)
	testAlice := common.BytesToAddress(tx.data.ExtraData)
	assert.Equal(t, cs_crypto.GetLockAddress(testAlice, testBob), *tx.data.Recipient)
	try, _ := rlpHash(tx.wit.HashKey)
	assert.Equal(t, try, hashLock)
}
