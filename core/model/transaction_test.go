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
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var txAmount = big.NewInt(10000)

func TestNewTransaction(t *testing.T) {
	result := NewTransaction(1, common.HexToAddress("123"), big.NewInt(100), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{123})
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
	//log.Info("the tx is:","txData",tx)
	result, err := tx.EncodeRlpToBytes()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	//log.Info("the tx rlp data is:","rlpData",hexutil.Encode(result))
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
	assert.Equal(t, big.NewInt(52000), tx.Cost())
	assert.Equal(t, big.NewInt(107), tx.EstimateFee())
	assert.Equal(t, []byte{}, tx.ExtraData())
	assert.Equal(t, uint64(0), tx.Nonce())
	assert.Equal(t, big.NewInt(0), tx.TimeLock())
	assert.Equal(t, &bobAddr, tx.To())
	assert.Equal(t, common.TxType(common.AddressTypeNormal), tx.GetType())
	assert.True(t, tx.IsEqual(*tx))
	assert.NotNil(t, tx.String())

	// read from cache
	assert.Equal(t, common.StorageSize(107), tx.Size())
	assert.Equal(t, common.StorageSize(107), tx.Size())

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

	assert.False(t, txs.Less(0, 1))
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

// Tests that transactions can be correctly sorted according to their price in
// decreasing order, but at the same time with increasing nonces when issued by
// the same account.
func TestTransactionPriceNonceSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 2)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	signer := MercurySigner{big.NewInt(1)}
	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address][]AbstractTransaction{}
	for start, key := range keys {
		addr := cs_crypto.GetNormalAddress(key.PublicKey)
		for i := 0; i < 2; i++ {
			tx := NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), g_testData.TestGasPrice, g_testData.TestGasLimit, nil)
			tx.SignTx(key, signer)
			tx.PaddingTxIndex(i)
			groups[addr] = append(groups[addr], tx)
		}
	}

	/*	log.Info("the group txs is:")
		for addr,txs:=range groups{
			log.Info("the addr is:","addr",addr.Hex())
			for _,tx := range txs{
				log.Info("the tx is:","txId",tx.CalTxId().Hex())
			}
		}*/

	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByFeeAndNonce(signer, groups)

	log.Info("the txset head is:")
	for _, tx := range txset.heads {
		log.Info("the head tx is:", "txId", tx.CalTxId().Hex())
	}

	txs := make([]AbstractTransaction, 0)
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != 2*2 {
		t.Errorf("expected %d transactions, found %d", 25*25, len(txs))
	}
	for i, txi := range txs {
		fromi, _ := txi.Sender(signer)

		// Make sure the nonce order is valid
		for j, txj := range txs[i+1:] {
			fromj, _ := txj.Sender(signer)

			if fromi == fromj && txi.Nonce() > txj.Nonce() {
				t.Errorf("invalid nonce ordering: tx #%d (A=%x N=%v) < tx #%d (A=%x N=%v)", i, fromi[:4], txi.Nonce(), i+j, fromj[:4], txj.Nonce())
			}
		}

		// If the next tx has different from account, the price must be lower than the current one
		if i+1 < len(txs) {
			next := txs[i+1]
			fromNext, _ := next.Sender(signer)
			if fromi != fromNext && txi.GetGasPrice().Cmp(next.GetGasPrice()) < 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", i, fromi[:4], txi.GetGasPrice(), i+1, fromNext[:4], next.GetGasPrice())
			}
		}
	}
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
