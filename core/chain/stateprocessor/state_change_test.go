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

package stateprocessor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

func TestStateChangeList_DecodeRLP(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)

	// create accounts
	processor.NewAccountState(charlieAddr)
	processor.AddBalance(charlieAddr, big.NewInt(500))
	processor.DeleteAccountState(charlieAddr)
	processor.NewAccountState(aliceAddr)
	processor.NewAccountState(bobAddr)

	// add nonce, stake, balance
	for j := 0; j < 100; j++ {
		processor.AddNonce(bobAddr, uint64(j))
	}

	// snapShot
	snapShot := processor.Snapshot()
	for i := 0; i < 100; i++ {
		processor.AddStake(aliceAddr, big.NewInt(int64(i)))
		processor.AddBalance(aliceAddr, big.NewInt(int64(i)))
	}

	// set info
	for i := 0; i < 5; i++ {
		processor.SetCode(aliceAddr, []byte{123})
		processor.SetAbi(aliceAddr, []byte{123})
		processor.SetDataRoot(aliceAddr, common.HexToHash("dataRoot"))
		processor.SetData(aliceAddr, "key", []byte("value"))
		processor.SetTimeLock(aliceAddr, big.NewInt(500))
		processor.SetHashLock(aliceAddr, common.HexToHash("hashLock"))
		processor.SetVerifyNum(aliceAddr, uint64(100))
		processor.SetCommitNum(aliceAddr, uint64(100))
		processor.SetPerformance(aliceAddr, uint64(100))
		processor.SetLastElect(aliceAddr, uint64(100))
		processor.PutContract(aliceAddr, reflect.ValueOf(&erc20{}))
	}

	// encode state change list
	sclSent := processor.stateChangeList.digest()
	enc, err := rlp.EncodeToBytes(sclSent)
	assert.NoError(t, err)

	// decode state change list
	var sclGet StateChangeList
	err2 := rlp.DecodeBytes(enc, &sclGet)
	assert.NoError(t, err2)

	processor2, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor2.stateChangeList = &sclGet
	processor2.stateChangeList.recover(processor2)

	// assert accounts
	_, charlieErr := processor2.GetAccountState(charlieAddr)
	assert.Equal(t, charlieErr, gerror.ErrAccountNotExist)
	aliceBalance, err := processor2.GetBalance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, aliceBalance, big.NewInt(4950))
	aliceStake, err := processor2.GetBalance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, aliceStake, big.NewInt(4950))
	bobNonce, err := processor2.GetNonce(bobAddr)
	assert.NoError(t, err)
	assert.Equal(t, bobNonce, uint64(4950))

	// revert to snap shot
	processor.RevertToSnapshot(snapShot)

	// assert accounts
	assert.Equal(t, false, processor.IsEmptyAccount(aliceAddr))
	assert.Equal(t, false, processor.IsEmptyAccount(bobAddr))
	assert.Equal(t, true, processor.IsEmptyAccount(charlieAddr))
	aliceBalance, err = processor.GetBalance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, aliceBalance, big.NewInt(0))
	aliceStake, err = processor.GetBalance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, aliceStake, big.NewInt(0))
	bobNonce, err = processor.GetNonce(bobAddr)
	assert.NoError(t, err)
	assert.Equal(t, bobNonce, uint64(4950))
}

func TestStateChangeList_AccountState(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)

	// create accounts
	snapShot := processor.Snapshot()
	processor.NewAccountState(aliceAddr)
	processor.AddBalance(aliceAddr, big.NewInt(500))
	processor.NewAccountState(aliceAddr)
	assert.Panics(t, func() {
		processor.stateChangeList.digest()
	})

	processor.RevertToSnapshot(snapShot)
	assert.Equal(t, true, processor.IsEmptyAccount(aliceAddr))
	snapShot = processor.Snapshot()

	// delete accounts
	processor.NewAccountState(aliceAddr)
	processor.AddBalance(aliceAddr, big.NewInt(500))
	processor.DeleteAccountState(aliceAddr)
	processor.DeleteAccountState(aliceAddr)
	assert.Panics(t, func() {
		processor.stateChangeList.digest()
	})

	processor.RevertToSnapshot(snapShot)
	assert.Equal(t, true, processor.IsEmptyAccount(aliceAddr))
}

func TestStateChangeList_Logs(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)

	// create accounts
	snapShot := processor.Snapshot()
	txHash := common.HexToHash("txHash")
	log := &model.Log{TxHash: txHash}
	processor.NewAccountState(aliceAddr)
	for i := 0; i < 50; i++ {
		processor.AddLog(log)
	}
	snapShot = processor.Snapshot()
	assert.Equal(t, 50, len(processor.GetLogs(txHash)))
	for i := 0; i < 50; i++ {
		processor.AddLog(log)
	}
	assert.Equal(t, 100, len(processor.GetLogs(txHash)))
	processor.RevertToSnapshot(snapShot)
	assert.Equal(t, 50, len(processor.GetLogs(txHash)))

	// encode state change list
	enc, err := rlp.EncodeToBytes(processor.stateChangeList)
	assert.NoError(t, err)

	// decode state change list
	var sclGet StateChangeList
	err2 := rlp.DecodeBytes(enc, &sclGet)
	assert.NoError(t, err2)

	processor2, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor2.stateChangeList = &sclGet
	processor2.stateChangeList.recover(processor2)
	assert.Equal(t, 50, len(processor2.GetLogs(txHash)))
}


