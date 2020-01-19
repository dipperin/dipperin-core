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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/economymodel"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
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

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"Alice stake 30",
			given: func() error {
				err := processor.Stake(aliceAddr, big.NewInt(30))
				aliceStake, _ := processor.GetStake(aliceAddr)
				aliceBalance, _ := processor.GetBalance(aliceAddr)
				assert.EqualValues(t, big.NewInt(30), aliceStake)
				assert.EqualValues(t, big.NewInt(8999970), aliceBalance)
				return err
			},
			expect:result{nil},
		},
		{
			name:"Not enough money",
			given: func() error {
				err := processor.Stake(aliceAddr, big.NewInt(1e7))
				return err
			},
			expect:result{gerror.ErrBalanceNotEnough},
		},
		{
			name:"Account not exit",
			given: func() error {
				err := processor.Stake(common.HexToAddress("test"), big.NewInt(20))
				return err
			},
			expect:result{gerror.ErrBalanceNotEnough},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
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

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"aliceAddr unStake",
			given: func() error {
				err := processor.UnStake(aliceAddr)
				return err
			},
			expect:result{nil},
		},
		{
			name:"unknown address",
			given: func() error {
				err := processor.UnStake(common.HexToAddress("123"))
				return err
			},
			expect:result{errors.New("account does not exist")},
		},
		{
			name:"Bob has no stake at all",
			given: func() error {
				bobStake, _ := processor.GetStake(bobAddr)
				assert.EqualValues(t, big.NewInt(0), bobStake)
				err := processor.UnStake(bobAddr)
				return err
			},
			expect:result{gerror.ErrStakeNotEnough},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
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

	type result struct {
		isNot bool
		err error
	}

	testCases := []struct {
		name   string
		given  func() (bool,error)
		expect result
	}{
		{
			name:"Move stake from alice's stake to bob's balance",
			given: func() (bool, error) {
				err := processor.MoveStakeToAddress(aliceAddr, bobAddr)
				aliceStake, _ := processor.GetStake(aliceAddr)
				aliceBalance, _ := processor.GetBalance(aliceAddr)
				bobStake, _ := processor.GetStake(bobAddr)
				bobBalance, _ := processor.GetBalance(bobAddr)
				return aliceStake.String() == big.NewInt(0).String()&& aliceBalance.String()==big.NewInt(8999200).String()&& bobStake.String()==big.NewInt(0).String() && bobBalance.String()==big.NewInt(800).String(), err
			},
			expect:result{true,nil},
		},
		{
			name:"Alice has no stake at all",
			given: func() (bool, error) {
				err := processor.MoveStakeToAddress(aliceAddr, bobAddr)
				return false, err
			},
			expect:result{false,gerror.ErrStakeNotEnough},
		},
		{
			name:"a unknown address to bob",
			given: func() (bool, error) {
				err := processor.MoveStakeToAddress(common.HexToAddress("123"), bobAddr)
				return false, err
			},
			expect:result{false,errors.New("account does not exist")},
		},
		{
			name:"add balance to alice",
			given: func() (bool, error) {
				processor.AddStake(aliceAddr, big.NewInt(100))
				err := processor.MoveStakeToAddress(aliceAddr, common.HexToAddress("123"))
				return true, err
			},
			expect:result{true,nil},
		},
	}

	for _,tc:=range testCases{
		is,err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.Equal(t,tc.expect.isNot,is)
		}
	}
}

func TestAccountStateDB_processStakeTx(t *testing.T) {
	minRegisterValue := economymodel.MiniPledgeValue
	db, root := CreateTestStateDBWithMutableBalance(minRegisterValue)
	twiceMiniPledgeValue := big.NewInt(0).Mul(economymodel.MiniPledgeValue, big.NewInt(2))
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"process Stake CancelTransaction",
			given: func() error {
				tx := getTestCancelTransaction(0, key1)
				err := processor.processStakeTx(tx)
				return err
			},
			expect:result{gerror.ErrTxTypeNotMatch},
		},
		{
			name:"process Stake unknown account RegisterTransaction",
			given: func() error {
				tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
				err := processor.processStakeTx(tx)
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
		{
			name:"process Stake 10 amount RegisterTransaction",
			given: func() error {
				key1, _ = createKey()
				tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
				err := processor.processStakeTx(tx)
				return err
			},
			expect:result{gerror.ErrStakeNotEnough},
		},
		{
			name:"process Stake twiceMiniPledgeValue amount RegisterTransaction",
			given: func() error {
				tx := getTestRegisterTransaction(0, key1, twiceMiniPledgeValue)
				err := processor.processStakeTx(tx)
				return err
			},
			expect:result{gerror.ErrBalanceNotEnough},
		},
		{
			name:"process Stake minRegisterValue amount RegisterTransaction",
			given: func() error {
				tx := getTestRegisterTransaction(0, key1, minRegisterValue)
				err := processor.processStakeTx(tx)
				return err
			},
			expect:result{nil},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}

func TestAccountStateDB_processCancelTx(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"ErrTxTypeNotMatch",
			given: func() error {
				tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
				err := processor.processCancelTx(tx, 1)
				return err
			},
			expect:result{gerror.ErrTxTypeNotMatch},
		},
		{
			name:"ErrAccountNotExist",
			given: func() error {
				tx := getTestCancelTransaction(0, key1)
				err := processor.processCancelTx(tx, 1)
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
		{
			name:"StateSendRegisterTxFirst",
			given: func() error {
				key1, _ = createKey()
				tx := getTestCancelTransaction(0, key1)
				err := processor.processCancelTx(tx, 1)
				return err
			},
			expect:result{gerror.StateSendRegisterTxFirst},
		},
		{
			name:"delete alice last elect num from tree",
			given: func() error {
				tx := getTestCancelTransaction(0, key1)
				err := processor.AddStake(aliceAddr, big.NewInt(1e4))
				assert.NoError(t, err)
				err = processor.blockStateTrie.TryDelete(GetLastElectKey(aliceAddr))
				assert.NoError(t, err)
				err = processor.processCancelTx(tx, 1)
				return err
			},
			expect:result{errors.New("EOF")},
		},
		{
			name:"set alice last elect num",
			given: func() error {
				tx := getTestCancelTransaction(0, key1)
				err := processor.SetLastElect(aliceAddr, uint64(1))
				assert.NoError(t, err)
				err = processor.processCancelTx(tx, 1)
				return err
			},
			expect:result{gerror.StateSendRegisterTxFirst},
		},
		{
			name:"no error",
			given: func() error {
				tx := getTestCancelTransaction(0, key1)
				err := processor.SetLastElect(aliceAddr, uint64(0))
				assert.NoError(t, err)
				err = processor.processCancelTx(tx, 1)
				return err
			},
			expect:result{nil},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}

func TestAccountStateDB_processUnStakeTx(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"ErrTxTypeNotMatch",
			given: func() error {
				tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
				err := processor.processUnStakeTx(tx)
				return err
			},
			expect:result{gerror.ErrTxTypeNotMatch},
		},
		{
			name:"ErrAccountNotExist",
			given: func() error {
				tx := getTestUnStakeTransaction(0, key1)
				err := processor.processUnStakeTx(tx)
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
		{
			name:"StateSendRegisterTxFirst",
			given: func() error {
				key1, _ = createKey()
				tx := getTestUnStakeTransaction(0, key1)
				err = processor.processUnStakeTx(tx)
				return err
			},
			expect:result{gerror.StateSendRegisterTxFirst},
		},
		{
			name:"delete alice last elect num from tree",
			given: func() error {
				tx := getTestUnStakeTransaction(0, key1)
				err := processor.AddStake(aliceAddr, big.NewInt(1e4))
				assert.NoError(t, err)
				err = processor.blockStateTrie.TryDelete(GetLastElectKey(aliceAddr))
				err = processor.processUnStakeTx(tx)
				return err
			},
			expect:result{errors.New("EOF")},
		},
		{
			name:"set alice last elect num",
			given: func() error {
				tx := getTestUnStakeTransaction(0, key1)
				err := processor.SetLastElect(aliceAddr, uint64(0))
				assert.NoError(t, err)
				err = processor.processUnStakeTx(tx)
				return err
			},
			expect:result{gerror.StateSendCancelTxFirst},
		},
		{
			name:"no error",
			given: func() error {
				tx := getTestUnStakeTransaction(0, key1)
				err := processor.SetLastElect(aliceAddr, uint64(1))
				assert.NoError(t, err)
				err = processor.processUnStakeTx(tx)
				return err
			},
			expect:result{nil},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}

func TestAccountStateDB_processEvidenceTx_Error(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, _ := NewAccountStateDB(root, NewStateStorageWithCache(db))

	key1, err := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232033")
	assert.NoError(t, err)

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"ErrTxTypeNotMatch",
			given: func() error {
				tx := getTestRegisterTransaction(0, key1, big.NewInt(10))
				err = processor.processEvidenceTx(tx)
				return err
			},
			expect:result{gerror.ErrTxTypeNotMatch},
		},
		{
			name:"ErrReceiverNotExist",
			given: func() error {
				tx := getTestEvidenceTransaction(0, key1, common.HexToAddress("123"), &model.VoteMsg{}, &model.VoteMsg{})
				err = processor.processEvidenceTx(tx)
				return err
			},
			expect:result{gerror.ErrReceiverNotExist},
		},
		{
			name:"no error",
			given: func() error {
				key1, _ = createKey()
				err := processor.SetStake(bobAddr, big.NewInt(100))
				assert.NoError(t, err)
				tx := getTestEvidenceTransaction(0, key1, bobAddr, &model.VoteMsg{}, &model.VoteMsg{})
				err = processor.processEvidenceTx(tx)
				return err
			},
			expect:result{nil},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}

