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
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestDeriveShaByHash_DeriveSha(t *testing.T) {
	txs1 := CreateSignedTxList(300)
	try := DeriveSha(Transactions(txs1))
	ms := NewMercurySigner(big.NewInt(1))

	tree := new(trie.Trie)
	for i := 0; i < 300; i++ {
		enc, _ := rlp.EncodeToBytes(txs1[i])
		key := txs1[i].CalTxId().Bytes()
		tree.Update(key, enc)
	}
	get := tree.Hash()
	assert.Equal(t, try, get)

	key1, _ := crypto.HexToECDSA(alicePriv)
	tryTx1 := NewTransaction(uint64(1), bobAddr, big.NewInt(1000), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	tryTx1.SignTx(key1, ms)
	enc1 := tree.Get(tryTx1.CalTxId().Bytes())
	getTx1 := new(Transaction)
	if len(enc1) == 0 {
		t.Error("i am wrong ,dont punch me.")
	} else {

		rlp.DecodeBytes(enc1, getTx1)
	}

	assert.Equal(t, tryTx1.CalTxId(), getTx1.CalTxId())

	tryTx2 := NewTransaction(uint64(1000), bobAddr, big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	enc2 := tree.Get(tryTx2.CalTxId().Bytes())
	assert.Equal(t, len(enc2), 0)
}

func TestVerifications_GetKey(t *testing.T) {
	v := Verifications{}
	result := v.GetKey(0)
	assert.NotNil(t, result)
}

func TestVerifications_GetRlp(t *testing.T) {
	vote := CreateSignedVote(1, 0, common.Hash{}, VoteMessage)
	votes := Verifications{vote}
	result := votes.GetRlp(0)
	assert.NotNil(t, result)
}

func TestVerifications_Len(t *testing.T) {
	v := Verifications{}
	result := v.Len()
	assert.NotNil(t, result)
}

func TestAbsTransactions_GetKey(t *testing.T) {
	tx := AbsTransactions{CreateSignedTx(0, big.NewInt(10000))}
	result := tx.GetKey(0)
	assert.NotNil(t, result)
}

func TestAbsTransactions_GetRlp(t *testing.T) {
	tx := AbsTransactions{CreateSignedTx(0, big.NewInt(10000))}
	result := tx.GetRlp(0)
	assert.NotNil(t, result)
}

func TestAbsTransactions_Len(t *testing.T) {
	tx := AbsTransactions{CreateSignedTx(0, big.NewInt(10000))}
	result := tx.Len()
	assert.NotNil(t, result)
}
