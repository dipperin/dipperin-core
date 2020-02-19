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
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var path = "./transaction.out"

// generate n transactions list file on the specified path
func createTxListFile(n int, path string) {
	txs := createTxList(n)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Create(path)
	} else {
		os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, tx := range txs {
		err := rlp.Encode(f, tx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// convert tx list to map of address, todo: this should be optimized in the future
func txListToAddrMap(txs []model.AbstractTransaction) map[common.Address][]model.AbstractTransaction {
	_, keyBob, _ := createKey()
	bob := cs_crypto.GetNormalAddress(keyBob.PublicKey)
	m := make(map[common.Address][]model.AbstractTransaction)

	for _, tx := range txs {
		m[bob] = append(m[bob], tx)
	}

	return m
}

type transactions []model.AbstractTransaction

func TestTxJournalLoad(t *testing.T) {
	jNoExist := newTxJournal("not_exist")
	err := jNoExist.load(nil)
	assert.NoError(t, err)

	//jNoOpen := newTxJournal("/dev/block")
	//err = jNoOpen.load(nil)
	//assert.NotNil(t, err)

	mj := newTxJournal("transaction.out")
	err = mj.load(func(txs []model.AbstractTransaction) []error {
		return []error{errors.New("load test error")}
	})
	assert.NoError(t, err)
}

//func TestCreateTxListFile(t *testing.T) {
//	createTxListFile(1, path)
//}

// this is used for console inspection
func readFromTxListFile(t *testing.T, path string) {
	j := newTxJournal(path)
	defer j.close()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(path, "file not exists")
		return
	}
	fmt.Println("reading from", path)
	err := j.load(func(txs []model.AbstractTransaction) []error {
		var res []error
		for _, tx := range txs {
			fmt.Println(tx.Nonce())
			assert.NotNil(t, tx)
			res = append(res, nil)
		}
		return res
	})

	assert.NoError(t, err)
}

//func TestReadFromFile(t *testing.T) {
//	readFromTxListFile(t, path)
//}

func TestTxJournalRotate(t *testing.T) {
	j := newTxJournal(path)
	defer j.close()

	txs := createTxList(10)

	var abstractTxs transactions
	for _, tx := range txs {
		abstractTxs = append(abstractTxs, tx)
	}

	m := txListToAddrMap(abstractTxs)
	err := j.rotate(m)

	assert.NoError(t, err)
}

type fakePool struct {
	txs []model.AbstractTransaction
}

func newFakePool() *fakePool {
	return &fakePool{
		txs: make([]model.AbstractTransaction, 1),
	}
}

func (p *fakePool) Add(txs []model.AbstractTransaction) []error {
	p.txs = []model.AbstractTransaction{}

	res := make([]error, len(txs))

	for i, txs := range txs {
		p.txs = append(p.txs, txs)
		res[i] = nil
	}

	return res
}

func TestTxJournal_InsertFail(t *testing.T) {
	createTxListFile(2, path)

	p := newFakePool()

	j := newTxJournal(path)
	defer j.close()

	err := j.load(p.Add)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(p.txs))

	txs := createTxList(4)
	tx := txs[3]

	err = j.insert(tx)

	// no active journal
	assert.Error(t, err)
}

func TestTxJournal_Insert(t *testing.T) {
	createTxListFile(2, path)

	j := newTxJournal(path)
	defer j.close()

	p := newFakePool()

	err := j.load(p.Add)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(p.txs))

	// rotate and assign a writer
	txs := createTxList(3)

	// convert to abstraction transaction slice, must do this because interfaces are awesome!
	var abstractTxs transactions
	for _, tx := range txs {
		abstractTxs = append(abstractTxs, tx)
	}

	m := txListToAddrMap(abstractTxs)

	err = j.rotate(m)
	assert.NoError(t, err)

	txs = createTxList(4)
	tx := txs[3]

	// test how insert action differs form previous
	err = j.insert(tx)
	assert.NoError(t, err)

	// successfully inserted
	err = j.load(p.Add)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(p.txs))
}
