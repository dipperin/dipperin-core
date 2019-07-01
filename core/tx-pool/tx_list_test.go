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

package tx_pool

import (
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
)

var testPriv1 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
var testPriv2 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
var testPriv3 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033"
var tx1hash = common.HexToHash("0x528131488f97c6314b2fa0dff404f1037067e787b65cb244d79c7ecea007c0d5")
var tx2hash = common.HexToHash("0x0aedd7a6779339cc44fe1e51cdf42b4bf3a557d52e646390e6d6bf6d489a5de3")

func createKey() (*ecdsa.PrivateKey, *ecdsa.PrivateKey, *ecdsa.PrivateKey) {
	key1, err1 := crypto.HexToECDSA(testPriv1)
	key2, err2 := crypto.HexToECDSA(testPriv2)
	key3, err3 := crypto.HexToECDSA(testPriv3)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, nil, nil
	}
	return key1, key2, key3
}

func createKeyBatch(num int) (keys []*ecdsa.PrivateKey) {
	keyBase := []byte(testPriv1)
	baseLen := len(keyBase)

	pat := []byte{}
	patLen := len([]byte(strconv.Itoa(num)))
	for i := 0; i < patLen; i++ {
		pat = append(pat, '0')
	}

	keyBase = append(keyBase[:baseLen-patLen], pat...)

	for i := 0; i < num; i++ {
		s := []byte(strconv.Itoa(i))
		slice := append(keyBase[:baseLen-len(s)], s...)
		//fmt.Println(i, "=", string(slice))
		key, err := crypto.HexToECDSA(string(slice))
		if err != nil {
			return nil
		}
		keys = append(keys, key)
	}

	return
}

func createTxList(n int) []*model.Transaction {
	keyAlice, keyBob, _ := createKey()
	ms := model.NewMercurySigner(big.NewInt(1))

	bob := cs_crypto.GetNormalAddress(keyBob.PublicKey)

	var res []*model.Transaction
	for i := 0; i < n; i++ {
		temptx := model.NewTransaction(uint64(i+1), bob, big.NewInt(int64(i)), big.NewInt(0).Mul(big.NewInt(int64(i)), g_testData.TestGasPrice), g_testData.TestGasLimit, []byte{})
		temptx.SignTx(keyAlice, ms)
		res = append(res, temptx)
	}
	return res
}

func createTxListWithFee(n int) []*model.Transaction {
	keyAlice, keyBob, _ := createKey()
	ms := model.NewMercurySigner(big.NewInt(1))

	bob := cs_crypto.GetNormalAddress(keyBob.PublicKey)

	var res []*model.Transaction
	for i := 0; i < n; i++ {
		temptx := model.NewTransaction(uint64(i+1), bob, big.NewInt(int64(i)), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})

		temptx.SignTx(keyAlice, ms)
		res = append(res, temptx)
	}
	return res
}

type MockPool struct {
	all     *txLookup
	feeList *txFeeList
}

func TestTxList_Put(t *testing.T) {
	txMap := newTxSortedMap()

	txs := createTxList(2)

	txMap.Put(txs[1])
}

func TestTxLookup_Remove(t *testing.T) {
	p := &MockPool{
		all: newTxLookup(),
	}

	p.feeList = newTxFeeList(p.all)

	txs := createTxList(100)
	for _, tx := range txs {
		p.all.Add(tx)
	}

	for i := 0; i < 100; i++ {
		p.all.Remove(txs[i].CalTxId())
		p.feeList.Removed()
	}

	p.feeList.Removed()

}

func TestTxList_Add(t *testing.T) {
	keyAlice, keyBob, _ := createKey()
	ms := model.NewMercurySigner(big.NewInt(1))

	l := newTxList(true)

	bob := cs_crypto.GetNormalAddress(keyBob.PublicKey)
	tx1 := model.NewTransaction(1, bob, big.NewInt(1), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	tx1.SignTx(keyAlice, ms)

	tx2 := model.NewTransaction(2, bob, big.NewInt(2), big.NewInt(2), g_testData.TestGasLimit, []byte{})
	tx2.SignTx(keyAlice, ms)

	tx3 := model.NewTransaction(2, bob, big.NewInt(3), big.NewInt(3), g_testData.TestGasLimit, []byte{})
	tx3.SignTx(keyAlice, ms)

	tx4 := model.NewTransaction(1, bob, big.NewInt(3), big.NewInt(0), g_testData.TestGasLimit, []byte{})
	tx4.SignTx(keyAlice, ms)

	assert.True(t, l.Empty())

	// successfully add
	ok, replace := l.Add(tx1, 1)
	assert.False(t, l.Overlaps(tx2))
	assert.True(t, ok)
	assert.Nil(t, replace)

	assert.False(t, l.Empty())

	// successfully add
	ok, replace = l.Add(tx2, 2)
	assert.True(t, l.Overlaps(tx2))
	assert.True(t, ok)
	assert.Nil(t, replace)

	// overlaps due to the same nonce
	assert.True(t, l.Overlaps(tx3))
	assert.Equal(t, 2, l.Len())

	// fail due to old tx feeBump could be higher than
	// current tx fee
	ok, replace = l.Add(tx3, 200)
	assert.False(t, ok)
	assert.Nil(t, replace)

	// lowers down the feeBump, in real application, fee bump
	// is a configurable param consistent through program
	ok, replace = l.Add(tx3, 1)
	assert.True(t, ok)
	assert.EqualValues(t, replace, tx2)

	// fail due to old tx fee is higher than current fee
	ok, replace = l.Add(tx4, 0)
	assert.False(t, ok)
	assert.Nil(t, replace)

	assert.Equal(t, 2, l.Len())
}

func TestTxList_Cap(t *testing.T) {
	l := newTxList(true)

	txs := createTxList(10)

	for _, tx := range txs {
		l.Add(tx, 1)
	}

	// cap threshold is higher than current length
	// returns empty slice and original len stay unchanged
	capedTxs := l.Cap(20)

	assert.Len(t, capedTxs, 0)
	assert.Equal(t, 10, l.Len())

	// cap threshold is lower than current length
	// returns (len - cap) txs that are with higher nonce
	capedTxs = l.Cap(5)

	assert.Len(t, capedTxs, 5)
	assert.Equal(t, 5, l.Len())
}

func TestTxList_Ready(t *testing.T) {
	l := newTxList(true)

	txs := createTxList(10)

	for _, tx := range txs {
		l.Add(tx, 1)
	}

	// if there is no tx nonce that are less or equal than
	// the start param, it returns nil
	readyTx := l.Ready(0)
	assert.Len(t, readyTx, 0)

	// start param equals to the lowest nonce, retrieves all
	// the consecutive txs
	readyTx = l.Ready(1)

	assert.Len(t, readyTx, 10)
	assert.Equal(t, 0, l.Len())

	for _, tx := range txs {
		l.Add(tx, 1)
	}

	// start param is greater than the lowest nonce, retrieves all
	// the consecutive starting from the lowest nonce,
	// this is something should never happen, but be handle this
	// here rather than failing.
	readyTx = l.Ready(2)
}

func TestTxList_Filter(t *testing.T) {
	l := newTxList(true)

	txs := createTxList(10)
	for _, tx := range txs {
		l.Add(tx, 10)
	}

	// filter removes and returns all the txs that the cost is higher but not equal
	// to the given threshold
	removed, invalids := l.Filter(big.NewInt(0).Mul(big.NewInt(5),big.NewInt(int64(g_testData.TestGasLimit))))

	assert.Len(t, removed, 5)

	for _, tx := range removed {
		assert.True(t, tx.Nonce() > 5)
	}

	assert.Len(t, invalids, 0)
}

func TestTxList_FilterNonce(t *testing.T) {
	l := newTxList(true)

	txs := createTxList(10)
	for _, tx := range txs {
		l.Add(tx, 10)
	}

	threshold := uint64(5)
	removed := l.FilterNonce(threshold)

	assert.Len(t, removed, 4)
	for _, tx := range removed {
		assert.True(t, tx.Nonce() < threshold)
	}
}

func TestTxList_Remove(t *testing.T) {
	l := newTxList(true)

	txs := createTxList(10)
	for i := 0; i < 5; i++ {
		l.Add(txs[i], 10)
	}

	e, rem := l.Remove(txs[5])
	assert.Nil(t, rem)
	assert.Equal(t, e, false)

	//strict mode
	l.strict = true
	e, rem = l.Remove(txs[3])
	assert.Equal(t, e, true)
	assert.NotNil(t, rem)
}

func TestTxList_FilterNonce_WithCache(t *testing.T) {
	l := newTxList(true)

	txs := createTxList(10)
	for _, tx := range txs {
		l.Add(tx, 10)
	}

	l.txs.Sort()
	threshold := uint64(5)
	removed := l.FilterNonce(threshold)

	assert.Len(t, removed, 4)
	for _, tx := range removed {
		assert.True(t, tx.Nonce() < threshold)
	}
}

func TestTxList_Flatten(t *testing.T) {
	l := newTxList(true)

	txs := createTxList(10)
	for _, tx := range txs {
		l.Add(tx, 10)
	}

	flatten := l.Flatten()

	for i, tx := range txs {
		assert.EqualValues(t, tx.CalTxId(), flatten[i].CalTxId())
	}
}

func TestTxList_Pop(t *testing.T) {
	txs := createTxList(1)
	ph := (priceHeap)([]model.AbstractTransaction{txs[0]})
	tx := ph.Pop().(model.AbstractTransaction)

	tfl := newTxFeeList(&txLookup{all: make(map[common.Hash]model.AbstractTransaction)})
	tfl.Put(tx)
	tfl.Cap(big.NewInt(1), &accountSet{
		accounts: map[common.Address]struct{}{},
		signer:   model.NewMercurySigner(big.NewInt(1)),
	})

	tfl = newTxFeeList(&txLookup{all: map[common.Hash]model.AbstractTransaction{
		tx.CalTxId(): tx,
	}})
	tfl.Put(tx)
	tfl.Cap(big.NewInt(1), &accountSet{
		accounts: map[common.Address]struct{}{},
		signer:   model.NewMercurySigner(big.NewInt(1)),
	})

	tfl = newTxFeeList(&txLookup{all: map[common.Hash]model.AbstractTransaction{}})
	tfl.Put(tx)
	tfl.UnderPriced(tx, &accountSet{
		accounts: map[common.Address]struct{}{},
		signer:   model.NewMercurySigner(big.NewInt(1)),
	})

	tfl = newTxFeeList(&txLookup{all: map[common.Hash]model.AbstractTransaction{}})
	tfl.Put(tx)
	tfl.Discard(1, &accountSet{
		accounts: map[common.Address]struct{}{},
		signer:   model.NewMercurySigner(big.NewInt(1)),
	})

	tfl = newTxFeeList(&txLookup{all: map[common.Hash]model.AbstractTransaction{
		tx.CalTxId(): tx,
	}})
	tfl.Put(tx)
	tfl.Discard(1, &accountSet{
		accounts: map[common.Address]struct{}{},
		signer:   model.NewMercurySigner(big.NewInt(1)),
	})
}
