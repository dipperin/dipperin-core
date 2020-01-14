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

package txpool

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

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