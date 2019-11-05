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

package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

/*
Test verifier basic process
Include three function relate to verifier's stake
Stake: move money from balance to stake
UnStake: move all money from stake to balance
MoveStakeToAddress: move somebody's stake to another person's balance.
*/
func TestAccountStateProcessor_Stake(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)
	aliceOriginal, _ := processor.GetBalance(aliceAddr)
	assert.EqualValues(t, big.NewInt(9e6), aliceOriginal)

	//Valid
	//Alice stake 30
	err = processor.Stake(aliceAddr, big.NewInt(30))
	assert.NoError(t, err)
	aliceStake, err := processor.GetStake(aliceAddr)
	assert.NoError(t, err)
	aliceBalance, err := processor.GetBalance(aliceAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(30), aliceStake)
	assert.EqualValues(t, big.NewInt(8999970), aliceBalance)

	//Invalid
	err = processor.Stake(aliceAddr, big.NewInt(1e7)) //Not enough money
	assert.EqualValues(t, g_error.ErrBalanceNotEnough, err)

	err = processor.Stake(common.HexToAddress("test"), big.NewInt(20)) //Account not exit
	assert.EqualValues(t, g_error.ErrBalanceNotEnough, err)
}
func TestAccountStateProcessor_UnStake(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)
	err = processor.Stake(aliceAddr, big.NewInt(800))
	assert.NoError(t, err)

	aliceStake, _ := processor.GetStake(aliceAddr)
	assert.EqualValues(t, big.NewInt(800), aliceStake)
	aliceBalance, _ := processor.GetBalance(aliceAddr)
	assert.EqualValues(t, big.NewInt(8999200), aliceBalance)

	//Valid
	//Alice un stake
	err = processor.UnStake(common.HexToAddress("123"))
	assert.Error(t, err)
	err = processor.UnStake(aliceAddr)
	assert.NoError(t, err)

	aliceNewStake, _ := processor.GetStake(aliceAddr)
	aliceNewBalance, _ := processor.GetBalance(aliceAddr)
	assert.EqualValues(t, big.NewInt(0), aliceNewStake)
	assert.EqualValues(t, big.NewInt(9e6), aliceNewBalance)

	//Invalid
	bobStake, _ := processor.GetStake(bobAddr)
	assert.EqualValues(t, big.NewInt(0), bobStake)
	err = processor.UnStake(bobAddr) //Bob has no stake at all
	assert.EqualValues(t, g_error.ErrStakeNotEnough, err)
}
func TestAccountStateProcessor_MoveStakeToAddress(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))
	processor.Stake(aliceAddr, big.NewInt(800))
	aliceOriginalStake, _ := processor.GetStake(aliceAddr)
	aliceOriginalBalance, _ := processor.GetBalance(aliceAddr)
	bobOriginalStake, _ := processor.GetStake(bobAddr)
	bobOriginalBalance, _ := processor.GetBalance(bobAddr)

	assert.EqualValues(t, big.NewInt(800), aliceOriginalStake)
	assert.EqualValues(t, big.NewInt(8999200), aliceOriginalBalance)
	assert.EqualValues(t, big.NewInt(0), bobOriginalStake)
	assert.EqualValues(t, big.NewInt(0), bobOriginalBalance)

	//Valid
	// Move stake from alice's stake to bob's balance
	err := processor.MoveStakeToAddress(aliceAddr, bobAddr)
	aliceStake, _ := processor.GetStake(aliceAddr)
	aliceBalance, _ := processor.GetBalance(aliceAddr)
	bobStake, _ := processor.GetStake(bobAddr)
	bobBalance, _ := processor.GetBalance(bobAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0), aliceStake)
	assert.EqualValues(t, big.NewInt(8999200), aliceBalance)
	assert.EqualValues(t, big.NewInt(0), bobStake)
	assert.EqualValues(t, big.NewInt(800), bobBalance)

	//InValid
	//Alice has no stake at all
	err = processor.MoveStakeToAddress(aliceAddr, bobAddr)
	aliceStake, _ = processor.GetStake(aliceAddr)
	aliceBalance, _ = processor.GetBalance(aliceAddr)
	bobStake, _ = processor.GetStake(bobAddr)
	bobBalance, _ = processor.GetBalance(bobAddr)
	assert.EqualValues(t, g_error.ErrStakeNotEnough, err)

	//Error
	err = processor.MoveStakeToAddress(common.HexToAddress("123"), bobAddr)
	assert.Error(t, err)

	processor.AddStake(aliceAddr, big.NewInt(100))
	err = processor.MoveStakeToAddress(aliceAddr, common.HexToAddress("123"))
	assert.NoError(t, err)
}

func TestAccountStateDB_processStakeTx(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	tx := getTestCancelTransaction(0, key1)
	err = processor.processStakeTx(tx)
	assert.Equal(t, g_error.ErrTxTypeNotMatch, err)

	tx = getTestRegisterTransaction(0, key1, big.NewInt(10))
	err = processor.processStakeTx(tx)
	assert.Equal(t, g_error.ErrAccountNotExist, err)

	key1, _ = createKey()
	tx = getTestRegisterTransaction(0, key1, big.NewInt(10))
	err = processor.processStakeTx(tx)
	assert.Equal(t, g_error.ErrStakeNotEnough, err)

	tx = getTestRegisterTransaction(0, key1, big.NewInt(1e7))
	err = processor.processStakeTx(tx)
	assert.Equal(t, g_error.ErrBalanceNotEnough, err)

	tx = getTestRegisterTransaction(0, key1, big.NewInt(100))
	err = processor.processStakeTx(tx)
	assert.NoError(t, err)
}

func TestAccountStateDB_processCancelTx(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
	err = processor.processCancelTx(tx, 1)
	assert.Equal(t, g_error.ErrTxTypeNotMatch, err)

	tx = getTestCancelTransaction(0, key1)
	err = processor.processCancelTx(tx, 1)
	assert.Equal(t, g_error.ErrAccountNotExist, err)

	key1, _ = createKey()
	tx = getTestCancelTransaction(0, key1)
	err = processor.processCancelTx(tx, 1)
	assert.Equal(t, g_error.StateSendRegisterTxFirst, err)

	// add alice stake set last elect num
	err = processor.AddStake(aliceAddr, big.NewInt(1e4))
	assert.NoError(t, err)

	// delete alice last elect num from tree
	err = processor.blockStateTrie.TryDelete(GetLastElectKey(aliceAddr))
	assert.NoError(t, err)
	err = processor.processCancelTx(tx, 1)
	assert.Equal(t, "EOF", err.Error())

	// set alice last elect num
	err = processor.SetLastElect(aliceAddr, uint64(1))
	assert.NoError(t, err)

	err = processor.processCancelTx(tx, 1)
	assert.Equal(t, g_error.StateSendRegisterTxFirst, err)

	// no error
	err = processor.SetLastElect(aliceAddr, uint64(0))
	assert.NoError(t, err)
	err = processor.processCancelTx(tx, 1)
	assert.NoError(t, err)
}

func TestAccountStateDB_processUnStakeTx(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
	err = processor.processUnStakeTx(tx)
	assert.Equal(t, g_error.ErrTxTypeNotMatch, err)

	tx = getTestUnStakeTransaction(0, key1)
	err = processor.processUnStakeTx(tx)
	assert.Equal(t, g_error.ErrAccountNotExist, err)

	key1, _ = createKey()
	tx = getTestUnStakeTransaction(0, key1)
	err = processor.processUnStakeTx(tx)
	assert.Equal(t, g_error.StateSendRegisterTxFirst, err)

	// add alice stake set last elect num
	err = processor.AddStake(aliceAddr, big.NewInt(1e4))
	assert.NoError(t, err)

	// delete alice last elect num from tree
	err = processor.blockStateTrie.TryDelete(GetLastElectKey(aliceAddr))
	assert.NoError(t, err)
	err = processor.processUnStakeTx(tx)
	assert.Error(t, err)

	// set alice last elect num
	err = processor.SetLastElect(aliceAddr, uint64(0))
	assert.NoError(t, err)

	err = processor.processUnStakeTx(tx)
	assert.Equal(t, g_error.StateSendCancelTxFirst, err)

	// no error
	err = processor.SetLastElect(aliceAddr, uint64(1))
	assert.NoError(t, err)
	err = processor.processUnStakeTx(tx)
	assert.NoError(t, err)
}

func TestAccountStateDB_processEvidenceTx_Error(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
	err = processor.processEvidenceTx(tx)
	assert.Equal(t, g_error.ErrTxTypeNotMatch, err)

	tx = getTestEvidenceTransaction(0, key1, common.HexToAddress("123"), &model.VoteMsg{}, &model.VoteMsg{})
	err = processor.processEvidenceTx(tx)
	assert.Equal(t, g_error.ErrReceiverNotExist, err)

	key1, _ = createKey()
	err = processor.SetStake(bobAddr, big.NewInt(100))
	assert.NoError(t, err)
	tx = getTestEvidenceTransaction(0, key1, bobAddr, &model.VoteMsg{}, &model.VoteMsg{})
	err = processor.processEvidenceTx(tx)
	assert.NoError(t, err)
}
