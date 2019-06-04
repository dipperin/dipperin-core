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
	"bytes"
	"github.com/dipperin/dipperin-core/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var txAmount = big.NewInt(10000)

func TestNewTransaction(t *testing.T) {
	result := NewTransaction(1, common.HexToAddress("123"), big.NewInt(100), big.NewInt(10), []byte{123})
	assert.NotNil(t, result)
}

func TestNewContractCreation(t *testing.T) {
	result := NewContractCreation(1, big.NewInt(100), big.NewInt(10), uint64(21000), []byte{123})
	assert.NotNil(t, result)
}

func TestTransaction_EncodeRLP(t *testing.T) {
	tx := CreateSignedTx(0, big.NewInt(10000))
	buffer := new(bytes.Buffer)
	assert.NoError(t, tx.EncodeRLP(buffer))

	stream := rlp.NewStream(buffer, 0)
	assert.NoError(t, tx.DecodeRLP(stream))
}

func TestTransaction_EncodeRlpToBytes(t *testing.T) {
	tx := CreateSignedTx(0, big.NewInt(10000))
	result, err := tx.EncodeRlpToBytes()
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestTransaction_RLP(t *testing.T) {
	txSent, _ := createTestTx()
	enc, err := rlp.EncodeToBytes(txSent)
	assert.NoError(t, err)
	var txGet = &Transaction{}
	rlp.DecodeBytes(enc, txGet)
	assert.Equal(t, txSent.CalTxId(), txGet.CalTxId())
	assert.Equal(t, txSent.wit.R.Cmp(txGet.wit.R), 0)
	assert.Equal(t, txSent.wit.S.Cmp(txGet.wit.S), 0)
	assert.Equal(t, txSent.wit.V.Cmp(txGet.wit.V), 0)
	assert.True(t, bytes.Equal(txSent.wit.HashKey, txGet.wit.HashKey))
}

func TestTransaction_SenderPublicKey(t *testing.T) {
	tx := CreateSignedTx(0, big.NewInt(10000))
	result, err := tx.SenderPublicKey(nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	signer := NewMercurySigner(big.NewInt(100))
	result, err = tx.SenderPublicKey(signer)
	assert.Equal(t, ErrInvalidSig, err)
	assert.NotNil(t, result)

	signer = NewMercurySigner(big.NewInt(1))
	result, err = tx.SenderPublicKey(signer)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestTransaction(t *testing.T) {
	tx := CreateSignedTx(0, txAmount)
	assert.Equal(t, txAmount, tx.Amount())
	assert.Equal(t, big.NewInt(1), tx.ChainId())
	assert.Equal(t, big.NewInt(220000), tx.Cost())
	assert.Equal(t, big.NewInt(111), tx.EstimateFee())
	assert.Equal(t, []byte{}, tx.ExtraData())
	assert.Equal(t, big.NewInt(210000), tx.Fee())
	assert.Equal(t, uint64(0), tx.Nonce())
	assert.Equal(t, big.NewInt(0), tx.TimeLock())
	assert.Equal(t, &bobAddr, tx.To())
	assert.Equal(t, common.TxType(common.AddressTypeNormal), tx.GetType())
	assert.True(t, tx.IsEqual(*tx))
	assert.NotNil(t, tx.String())

	// read from cache
	assert.Equal(t, common.StorageSize(111), tx.Size())
	assert.Equal(t, common.StorageSize(111), tx.Size())

	signer := NewMercurySigner(big.NewInt(1))
	assert.Equal(t, signer, tx.GetSigner())

	var empty []byte
	var hashLock *common.Hash
	assert.Equal(t, empty, tx.HashKey())
	assert.Equal(t, hashLock, tx.HashLock())

	result, _, _ := tx.RawSignatureValues()
	assert.NotNil(t, result)

	sender, err := tx.Sender(signer)
	assert.NoError(t, err)
	assert.Equal(t, aliceAddr, sender)

	tx = &Transaction{}
	assert.Nil(t, tx.To())
}

func TestTransactionBy_Sort(t *testing.T) {
	block1 := CreateBlock(0, common.HexToHash("123"), 0)
	block2 := CreateBlock(1, common.HexToHash("123"), 0)

	bs := Blocks{block1, block2}
	BlockBy(func(b1, b2 *Block) bool { return true }).Sort(bs)
	TransactionBy(func(b1, b2 *Block) bool { return true }).Sort(bs)
}

func TestTransactions(t *testing.T) {
	tx1 := CreateSignedTx(0, txAmount)
	tx2 := CreateSignedTx(1, txAmount)
	txs := Transactions{tx1, tx2}

	assert.True(t, txs.Less(0, 1))
	assert.NotNil(t, txs.GetRlp(0))
	assert.NotNil(t, txs.String())

	txs.Swap(0, 1)
	assert.Equal(t, 2, txs.Len())
	assert.Equal(t, tx2.CalTxId().Bytes(), txs.GetKey(0))
	assert.Equal(t, tx1.CalTxId().Bytes(), txs.GetKey(1))
}

func TestTransactionsByFeeAndNonce(t *testing.T) {
	txs := make(map[common.Address][]AbstractTransaction, 1)
	tx1 := CreateSignedTx(0, txAmount)
	tx2 := CreateSignedTx(1, txAmount)

	txs[aliceAddr] = []AbstractTransaction{tx1, tx2}
	signer := NewMercurySigner(big.NewInt(1))

	tx := NewTransactionsByFeeAndNonce(signer, txs)
	assert.Equal(t, tx1.CalTxId(), tx.Peek().CalTxId())

	tx.Shift()
	assert.Equal(t, tx2.CalTxId(), tx.Peek().CalTxId())
	assert.Equal(t, 1, tx.heads.Len())

	tx.Pop()
	assert.Equal(t, nil, tx.Peek())
	assert.Equal(t, 0, tx.heads.Len())
}

func TestTxByFee(t *testing.T) {
	tx1 := CreateSignedTx(0, big.NewInt(100))
	tx2 := CreateSignedTx(1, big.NewInt(100))
	txs := TxByFee{tx1}

	txs.Push(tx2)
	assert.Equal(t, 2, txs.Len())
	assert.False(t, txs.Less(0, 1))

	txs.Swap(0, 1)
	assert.Equal(t, tx2.CalTxId(), txs[0].CalTxId())
	assert.Equal(t, tx1.CalTxId(), txs[1].CalTxId())

	txs.Pop()
	assert.Equal(t, 1, txs.Len())
	assert.Equal(t, tx2.CalTxId(), txs[0].CalTxId())
}

func TestTxSorter(t *testing.T) {
	block1 := CreateBlock(1, common.HexToHash("123"), 0)
	block2 := CreateBlock(2, block1.Hash(), 0)
	sorter := txSorter{
		blocks: []*Block{block1, block2},
		by: func(b1, b2 *Block) bool {
			return true
		},
	}
	assert.Equal(t, 2, sorter.Len())
	assert.True(t, sorter.Less(0, 1))

	sorter.Swap(0, 1)
	assert.Equal(t, block2.Number(), sorter.blocks[0].Number())
}

func TestTxDifference(t *testing.T) {
	tx1, tx2 := createTestTx()
	result := TxDifference([]AbstractTransaction{tx1}, []AbstractTransaction{tx2})
	assert.NotNil(t, result)
}
