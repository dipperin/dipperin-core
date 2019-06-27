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
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

func TestAccountStateDB_Commit(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	root := processor.PreStateRoot()
	assert.Equal(t, common.Hash{}, root)

	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)

	err = processor.AddBalance(aliceAddr, big.NewInt(2000))
	assert.NoError(t, err)
	err = processor.AddNonce(aliceAddr, 10)
	assert.NoError(t, err)

	fRoot, err := processor.Finalise()
	assert.NoError(t, err)
	savedRoot, err := processor.Commit()
	assert.Equal(t, fRoot, savedRoot)
}

func TestAccountStateDB_GetAccountState(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	alice, _ := processor.GetAccountState(aliceAddr)
	assert.Equal(t, alice.Nonce, uint64(0))
	assert.Equal(t, big.NewInt(9e6),alice.Balance)
}

func TestAccountStateDB_PutContract(t *testing.T) {
	cAddr := common.HexToAddress("0x3213123af")
	processor, err := NewAccountStateDB(common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)

	v, err := processor.GetContract(cAddr, reflect.TypeOf(nil))
	assert.Error(t, err)

	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)
	processor, err = NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	err = processor.PutContract(cAddr, reflect.ValueOf(nil))
	assert.Error(t, err)

	c := erc20{
		Owners: []string{"123", "234"},
		Balance: map[string]*big.Int{
			"123": big.NewInt(1e4),
			"234": big.NewInt(4e4),
		},
		Name: "jk",
		Dis:  10002,
	}

	err = processor.PutContract(cAddr, reflect.ValueOf(&c))
	assert.NoError(t, err)
	assert.False(t, processor.ContractExist(cAddr))

	_, err = processor.Commit()
	assert.NoError(t, err)
	assert.True(t, processor.ContractExist(cAddr))

	v, err = processor.GetContract(cAddr, reflect.TypeOf(c))
	assert.NoError(t, err)

	delete(processor.contractData, cAddr)
	v, err = processor.GetContract(cAddr, reflect.TypeOf(c))
	assert.NoError(t, err)

	contract := v.Interface().(*erc20)
	assert.Equal(t, c.Name, contract.Name)
	assert.Equal(t, c.Dis, contract.Dis)
}

func TestAccountStateDB_Snapshot(t *testing.T) {
	processor := createStateProcessor(t)

	processor.NewAccountState(aliceAddr)
	processor.NewAccountState(bobAddr)
	id := processor.Snapshot()

	performance, _ := processor.GetPerformance(aliceAddr)
	assert.Equal(t, uint64(30), performance)

	processor.AddBalance(aliceAddr, big.NewInt(1000))
	processor.AddStake(aliceAddr, big.NewInt(1000))
	processor.AddNonce(aliceAddr, uint64(5))
	processor.SetLastElect(aliceAddr, uint64(5))
	processor.SetCommitNum(aliceAddr, uint64(5))
	processor.SetPerformance(aliceAddr, uint64(5))
	processor.SetVerifyNum(aliceAddr, uint64(5))
	processor.SetHashLock(aliceAddr, common.HexToHash("123"))
	processor.SetDataRoot(aliceAddr, common.HexToHash("123"))
	//processor.setContractRoot(aliceAddr, common.HexToHash("123"))
	processor.SetTimeLock(aliceAddr, big.NewInt(10))
	processor.DeleteAccountState(bobAddr)

	balance, _ := processor.GetBalance(aliceAddr)
	stake, _ := processor.GetStake(aliceAddr)
	nonce, _ := processor.GetNonce(aliceAddr)
	lastElect, _ := processor.GetLastElect(aliceAddr)
	commit, _ := processor.GetCommitNum(aliceAddr)
	performance, _ = processor.GetPerformance(aliceAddr)
	verify, _ := processor.GetVerifyNum(aliceAddr)
	hashLock, _ := processor.GetHashLock(aliceAddr)
	timeLock, _ := processor.GetTimeLock(aliceAddr)
	dataRoot, _ := processor.GetDataRoot(aliceAddr)
	//contractRoot, _ := processor.GetContractRoot(aliceAddr)

	assert.Equal(t, big.NewInt(1000), balance)
	assert.Equal(t, big.NewInt(1000), stake)
	assert.Equal(t, uint64(5), nonce)
	assert.Equal(t, uint64(5), lastElect)
	assert.Equal(t, uint64(5), commit)
	assert.Equal(t, uint64(5), performance)
	assert.Equal(t, uint64(5), verify)
	assert.Equal(t, common.HexToHash("123"), hashLock)
	assert.Equal(t, common.HexToHash("123"), dataRoot)
	//assert.Equal(t, common.HexToHash("123"), contractRoot)
	assert.Equal(t, big.NewInt(10), timeLock)
	assert.True(t, processor.IsEmptyAccount(bobAddr))
	assert.False(t, processor.IsEmptyAccount(aliceAddr))

	processor.RevertToSnapshot(id)

	balance, _ = processor.GetBalance(aliceAddr)
	stake, _ = processor.GetStake(aliceAddr)
	nonce, _ = processor.GetNonce(aliceAddr)
	lastElect, _ = processor.GetLastElect(aliceAddr)
	commit, _ = processor.GetCommitNum(aliceAddr)
	performance, _ = processor.GetPerformance(aliceAddr)
	verify, _ = processor.GetVerifyNum(aliceAddr)
	hashLock, _ = processor.GetHashLock(aliceAddr)
	timeLock, _ = processor.GetTimeLock(aliceAddr)
	dataRoot, _ = processor.GetDataRoot(aliceAddr)
	//contractRoot, _ = processor.GetContractRoot(aliceAddr)

	assert.Equal(t, big.NewInt(0), balance)
	assert.Equal(t, big.NewInt(0), stake)
	assert.Equal(t, uint64(0), nonce)
	assert.Equal(t, uint64(0), lastElect)
	assert.Equal(t, uint64(0), commit)
	assert.Equal(t, uint64(30), performance)
	assert.Equal(t, uint64(0), verify)
	assert.Equal(t, common.Hash{}, hashLock)
	assert.Equal(t, common.Hash{}, dataRoot)
	//assert.Equal(t, common.Hash{}, contractRoot)
	assert.Equal(t, big.NewInt(0), timeLock)
	assert.False(t, processor.IsEmptyAccount(bobAddr))
	assert.False(t, processor.IsEmptyAccount(aliceAddr))
}

func TestMakeGenesisAccountStateProcessor(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)
	processor, err := MakeGenesisAccountStateProcessor(storage)
	assert.NoError(t, err)
	assert.NotNil(t, processor)
}


//Test verifier transaction process
//AddPeerSet, Elect, Evidence, Cancel

/*func TestAccountStateProcessor_Process_Register(t *testing.T) {
	processor := createStateProcessor(t)
	k1, _ := createKey()
	tx := getTestRegisterTransaction(0, k1, big.NewInt(1000))

	aliceStake, _ := processor.GetStake(aliceAddr)
	aliceBalance, _ := processor.GetBalance(aliceAddr)
	assert.EqualValues(t, big.NewInt(0), aliceStake)
	assert.EqualValues(t, big.NewInt(9000000), aliceBalance)

	err := processor.ProcessTx(tx, 1)
	assert.NoError(t, err)

	aliceStake, _ = processor.GetStake(aliceAddr)
	aliceBalance, _ = processor.GetBalance(aliceAddr)
	assert.EqualValues(t, big.NewInt(1000), aliceStake)
	assert.EqualValues(t, big.NewInt(8998960), aliceBalance)
}

func TestAccountStateProcessor_Process_Evidence(t *testing.T) {
	processor := createStateProcessor(t)
	key1, _ := createKey()

	// Valid evidence lead to stake move
	voteA := model.CreateSignedVote(1, 2, common.HexToHash("0x123456"), model.VoteMessage)
	voteB := model.CreateSignedVote(1, 2, common.HexToHash("0x654321"), model.VoteMessage)

	err := processor.AddStake(bobAddr, big.NewInt(1000))
	assert.NoError(t, err)

	bobStake, _ := processor.GetStake(bobAddr)
	assert.EqualValues(t, big.NewInt(1000), bobStake)

	tx := getTestEvidenceTransaction(0, key1, bobAddr, voteA, voteB)
	err = processor.ProcessTx(tx, 1)
	assert.NoError(t, err)
	aliceStake, _ := processor.GetStake(aliceAddr)
	aliceBalance, _ := processor.GetBalance(aliceAddr)
	bobStake, _ = processor.GetStake(bobAddr)

	assert.EqualValues(t, big.NewInt(0), aliceStake)
	assert.EqualValues(t, big.NewInt(9000960), aliceBalance)
	assert.EqualValues(t, big.NewInt(0), bobStake)
}

func TestAccountStateProcessor_Process_Cancel(t *testing.T) {
	processor := createStateProcessor(t)
	k1, _ := createKey()
	tx := getTestCancelTransaction(0, k1)

	err := processor.AddStake(aliceAddr, big.NewInt(1000))
	assert.NoError(t, err)

	aliceStake, _ := processor.GetStake(aliceAddr)
	aliceLastElect, _ := processor.GetLastElect(aliceAddr)
	assert.EqualValues(t, big.NewInt(1000), aliceStake)
	assert.EqualValues(t, uint64(0), aliceLastElect)

	err = processor.ProcessTx(tx, 1)
	assert.NoError(t, err)

	aliceStake, _ = processor.GetStake(aliceAddr)
	aliceLastElect, _ = processor.GetLastElect(aliceAddr)
	assert.EqualValues(t, big.NewInt(1000), aliceStake)
	assert.EqualValues(t, uint64(1), aliceLastElect)
}

func TestAccountStateProcessor_Process_UnStake(t *testing.T) {
	processor := createStateProcessor(t)
	k1, _ := createKey()
	tx := getTestUnStakeTransaction(0, k1)

	err := processor.AddStake(aliceAddr, big.NewInt(1000))
	assert.NoError(t, err)
	err = processor.SetLastElect(aliceAddr, uint64(5))
	assert.NoError(t, err)

	aliceStake, _ := processor.GetStake(aliceAddr)
	aliceLastElect, _ := processor.GetLastElect(aliceAddr)
	assert.EqualValues(t, big.NewInt(1000), aliceStake)
	assert.EqualValues(t, uint64(5), aliceLastElect)

	err = processor.ProcessTx(tx, 1)
	assert.NoError(t, err)

	aliceStake, _ = processor.GetStake(aliceAddr)
	aliceLastElect, _ = processor.GetLastElect(aliceAddr)
	assert.EqualValues(t, big.NewInt(0), aliceStake)
	assert.EqualValues(t, uint64(5), aliceLastElect)
}
*/
func createStateProcessor(t *testing.T) *AccountStateDB {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))
	aliceOriginalStake, _ := processor.GetStake(aliceAddr)
	aliceOriginalBalance, _ := processor.GetBalance(aliceAddr)
	aliceOriginalNonce, _ := processor.GetNonce(aliceAddr)
	aliceOriginalLastElect, _ := processor.GetLastElect(aliceAddr)
	processor.newAccountState(bobAddr)

	assert.EqualValues(t, big.NewInt(0), aliceOriginalStake)
	assert.EqualValues(t, big.NewInt(9000000), aliceOriginalBalance)
	assert.EqualValues(t, uint64(0), aliceOriginalNonce)
	assert.EqualValues(t, uint64(0), aliceOriginalLastElect)
	return processor
}

/*func TestAccountStateDB_ProcessTx(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)
	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)
	err = processor.AddBalance(aliceAddr, big.NewInt(1e4))
	assert.NoError(t, err)

	tx := fakeTransaction{
		txType: common.AddressTypeCross,
		nonce:  0,
		sender: aliceAddr,
	}
	err = processor.ProcessTx(tx, 1)
	assert.Error(t, err)

	tx = fakeTransaction{
		txType: common.AddressTypeERC20,
		nonce:  1,
		sender: aliceAddr,
	}
	err = processor.ProcessTx(tx, 1)
	assert.Error(t, err)

	tx = fakeTransaction{
		txType: common.AddressTypeEarlyReward,
		nonce:  2,
		sender: aliceAddr,
	}
	err = processor.ProcessTx(tx, 1)
	assert.Error(t, err)

	tx = fakeTransaction{
		txType: 0x0099,
		nonce:  3,
		sender: aliceAddr,
	}
	err = processor.ProcessTx(tx, 1)
	assert.Equal(t, g_error.UnknownTxTypeErr, err)

	tx = fakeTransaction{err: TxError}
	err = processor.ProcessTx(tx, 1)
	assert.Equal(t, TxError, err)
}*/

func TestAccountStateDB_processBasicTx_Error(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	tx := fakeTransaction{}
	conf := &TxProcessConfig{
		Tx: &tx,
		TxFee:big.NewInt(0),
	}
	err = processor.processBasicTx(conf)
	assert.Equal(t, SenderOrReceiverIsEmptyErr, err)

	tx = fakeTransaction{sender:aliceAddr}
	conf.Tx = &tx
	err = processor.processBasicTx(conf)
	assert.Equal(t, SenderNotExistErr, err)

	err = processor.blockStateTrie.TryUpdate(GetNonceKey(aliceAddr), []byte{})
	assert.NoError(t, err)
	err = processor.processBasicTx(conf)
	assert.Equal(t, SenderNotExistErr, err)

	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)
	tx = fakeTransaction{
		sender:aliceAddr,
		nonce:1,
	}
	conf.Tx = &tx
	err = processor.processBasicTx(conf)
	assert.Equal(t, g_error.ErrTxNonceNotMatch, err)

	tx = fakeTransaction{sender:aliceAddr}
	conf.Tx = &tx
	err = processor.processBasicTx(conf)
	assert.Equal(t, g_error.BalanceNegErr, err)
}

func TestAccountStateDB_processNormalTx_Error(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	err = processor.blockStateTrie.TryUpdate(GetNonceKey(aliceAddr), []byte{})
	assert.NoError(t, err)

	tx := fakeTransaction{sender:aliceAddr}
	err = processor.processNormalTx(tx)
	assert.Equal(t, g_error.AccountNotExist, err)
}

func TestAccountStateDB_Commit_Error(t *testing.T) {
	processor, err := NewAccountStateDB(common.Hash{}, fakeStateStorage{getErr:TrieError})
	assert.NoError(t, err)

	value := 0
	processor.contractData[aliceAddr] = reflect.ValueOf(&value)
	_, err = processor.Commit()
	assert.Error(t, err)
}

func TestAccountStateDB_SetError(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)

	processor, err := NewAccountStateDB(common.HexToHash("123"), tdb)
	assert.Error(t, err)
	processor, err = NewAccountStateDB(common.Hash{}, tdb)
	assert.NoError(t, err)

	_, err = processor.GetAccountState(aliceAddr)
	assert.Equal(t, g_error.AccountNotExist, err)
	assert.Equal(t, g_error.AccountNotExist, processor.AddBalance(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, g_error.AccountNotExist, processor.AddStake(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, g_error.AccountNotExist, processor.SubStake(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, g_error.AccountNotExist, processor.AddNonce(aliceAddr, uint64(5)))

	assert.Equal(t, g_error.AccountNotExist, processor.SetStake(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, g_error.AccountNotExist, processor.SetBalance(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, g_error.AccountNotExist, processor.SetNonce(aliceAddr, uint64(5)))
	assert.Equal(t, g_error.AccountNotExist, processor.SetLastElect(aliceAddr, uint64(5)))
	assert.Equal(t, g_error.AccountNotExist, processor.SetCommitNum(aliceAddr, uint64(5)))
	assert.Equal(t, g_error.AccountNotExist, processor.SetPerformance(aliceAddr, uint64(5)))
	assert.Equal(t, g_error.AccountNotExist, processor.SetVerifyNum(bobAddr, uint64(5)))
	assert.Equal(t, g_error.AccountNotExist, processor.SetHashLock(aliceAddr, common.HexToHash("123")))
	assert.Equal(t, g_error.AccountNotExist, processor.SetDataRoot(aliceAddr, common.HexToHash("123")))
	//processor.setContractRoot(aliceAddr, common.HexToHash("123"))
	assert.Equal(t, g_error.AccountNotExist, processor.SetTimeLock(aliceAddr, big.NewInt(10)))

	processor, err = NewAccountStateDB(common.Hash{}, fakeStateStorage{setErr:TrieError})
	assert.Equal(t, TrieError, processor.NewAccountState(aliceAddr))
	assert.Equal(t, TrieError, processor.DeleteAccountState(aliceAddr))
}

func TestAccountStateDB_SetError2(t *testing.T) {
	processor, err := NewAccountStateDB(common.Hash{}, fakeStateStorage{setErr:TrieError})
	assert.NoError(t, err)

	assert.Equal(t, TrieError, processor.AddBalance(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, TrieError, processor.SetStake(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, TrieError, processor.SetBalance(aliceAddr, big.NewInt(1000)))
	assert.Equal(t, TrieError, processor.SetNonce(aliceAddr, uint64(5)))
	assert.Equal(t, TrieError, processor.SetLastElect(aliceAddr, uint64(5)))
	assert.Equal(t, TrieError, processor.SetCommitNum(aliceAddr, uint64(5)))
	assert.Equal(t, TrieError, processor.SetPerformance(aliceAddr, uint64(5)))
	assert.Equal(t, TrieError, processor.SetVerifyNum(bobAddr, uint64(5)))
	assert.Equal(t, TrieError, processor.SetHashLock(aliceAddr, common.HexToHash("123")))
	assert.Equal(t, TrieError, processor.SetDataRoot(aliceAddr, common.HexToHash("123")))
	//processor.setContractRoot(aliceAddr, common.HexToHash("123"))
	assert.Equal(t, TrieError, processor.SetTimeLock(aliceAddr, big.NewInt(10)))
}

func TestAccountStateDB_GetError(t *testing.T) {
	processor, err := NewAccountStateDB(common.Hash{}, fakeStateStorage{getErr:TrieError})
	assert.NoError(t, err)

	_, err = processor.GetNonce(aliceAddr)
	assert.Equal(t, TrieError, err)

	processor, err = NewAccountStateDB(common.Hash{}, fakeStateStorage{getErr:TrieError, passKey:nonceKeySuffix})
	assert.NoError(t, err)

	_, err = processor.GetBalance(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetStake(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetLastElect(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetCommitNum(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetPerformance(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetVerifyNum(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetHashLock(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetTimeLock(aliceAddr)
	assert.Equal(t, TrieError, err)
	_, err = processor.GetDataRoot(aliceAddr)
	assert.Equal(t, TrieError, err)
	//contractRoot, _ := processor.GetContractRoot(aliceAddr)

	processor, err = NewAccountStateDB(common.Hash{}, fakeStateStorage{decodeErr:true})
	assert.NoError(t, err)

	_, err = processor.GetNonce(aliceAddr)
	assert.Error(t, err)

	processor, err = NewAccountStateDB(common.Hash{}, fakeStateStorage{passKey: nonceKeySuffix, decodeErr:true})
	assert.NoError(t, err)

	_, err = processor.GetBalance(aliceAddr)
	assert.Error(t, err)
	_, err = processor.GetStake(aliceAddr)
	assert.Error(t, err)
	_, err = processor.GetLastElect(aliceAddr)
	assert.Error(t, err)
	_, err = processor.GetCommitNum(aliceAddr)
	assert.Error(t, err)
	_, err = processor.GetPerformance(aliceAddr)
	assert.Error(t, err)
	_, err = processor.GetVerifyNum(aliceAddr)
	assert.Error(t, err)
	_, err = processor.GetHashLock(aliceAddr)
	assert.Error(t, err)
	_, err = processor.GetTimeLock(aliceAddr)
	assert.Error(t, err)
}

func TestGetContractAddrAndKey(t *testing.T) {
	address, key := GetContractAddrAndKey(GetNonceKey(aliceAddr))
	assert.Equal(t, aliceAddr, address)
	assert.Equal(t, nonceKeySuffix, string(key))

	address, key = GetContractAddrAndKey([]byte{})
	assert.Equal(t, common.Address{}, address)
	assert.Nil(t, key)
}