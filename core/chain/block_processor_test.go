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

package chain

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

func TestNewBlockProcessor(t *testing.T) {
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)
	assert.NotNil(t, processor)

	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{storageErr: TrieError})
	assert.Equal(t, TrieError, err)
	assert.Nil(t, processor)
}

func TestBlockProcessor_Process(t *testing.T) {
	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"NewAccountStateDB and NewBlockProcessor",
			given: func() error {
				db, root := createTestStateDB(t)
				tdb := stateprocessor.NewStateStorageWithCache(db)
				state, err := stateprocessor.NewAccountStateDB(root, tdb)
				_, err = NewBlockProcessor(fakeAccountDBChain{state: state}, root, tdb)
				return err
			},
			expect:result{nil},
		},
		{
			name:"return block number is zero",
			given: func() error {
				db, root := createTestStateDB(t)
				tdb := stateprocessor.NewStateStorageWithCache(db)
				state, _ := stateprocessor.NewAccountStateDB(root, tdb)
				processor, err := NewBlockProcessor(fakeAccountDBChain{state: state}, root, tdb)
				block := model.CreateBlock(0, common.Hash{}, 0)
				err = processor.Process(block, fakeEconomyModel{})
				return err
			},
			expect:result{nil},
		},
		{
			name:"block number 9 isn't change point",
			given: func() error {
				db, root := createTestStateDB(t)
				tdb := stateprocessor.NewStateStorageWithCache(db)
				state, _ := stateprocessor.NewAccountStateDB(root, tdb)
				processor, err := NewBlockProcessor(fakeAccountDBChain{state: state}, root, tdb)
				block := createBlock(10)
				err = processor.Process(block, fakeEconomyModel{})
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
		{
			name:"block number 19 is change point",
			given: func() error {
				db, root := createTestStateDB(t)
				tdb := stateprocessor.NewStateStorageWithCache(db)
				state, _ := stateprocessor.NewAccountStateDB(root, tdb)
				processor, err := NewBlockProcessor(fakeAccountDBChain{state: state}, root, tdb)
				block := model.CreateBlock(20, common.Hash{}, 0)
				err = processor.Process(block, fakeEconomyModel{})
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

func TestBlockProcessor_Process_Error(t *testing.T) {
	header := model.NewHeader(1, 10, common.Hash{}, common.HexToHash("1111"), minDiff, big.NewInt(324234), common.Address{}, common.BlockNonceFromInt(432423))
	tx := createUnNormalTx()
	block := model.NewBlock(header, []*model.Transaction{tx}, nil)

	db, root := createTestStateDB(t)
	tdb := stateprocessor.NewStateStorageWithCache(db)
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, root, tdb)
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
			name:"TxIterator inner error",
			given: func() error {
				err = processor.Process(block, fakeEconomyModel{})
				return err
			},
			expect:result{gerror.ErrUnknownTxType},
		},
		{
			name:"doReward- RewardCoinBase no coin base error",
			given: func() error {
				block = model.NewBlock(header, nil, nil)
				err = processor.Process(block, fakeEconomyModel{})
				return err
			},
			expect:result{gerror.InvalidCoinBaseAddressErr},
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

func TestBlockProcessor_processorCommitList_Error(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := stateprocessor.NewStateStorageWithCache(db)
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, tdb)
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
			name:"NotHavePreBlock Error",
			given: func() error {
				block := createBlock(1)
				err = processor.processCommitList(block, false)
				return err
			},
			expect:result{gerror.NotHavePreBlockErr},
		},
		{
			name:"Account Not Exist",
			given: func() error {
				block := createBlock(20)
				err = processor.processCommitList(block, false)
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
		{
			name:"NewAccountState AccountNotExist",
			given: func() error {
				block := createBlock(20)
				for i := 0; i < len(VerifierAddress); i++ {
					err = processor.NewAccountState(VerifierAddress[i])
				}
				err = processor.processCommitList(block, false)
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
	}

	for _,tc:=range testCases {
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}

func TestBlockProcessor_doReward_Error(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := stateprocessor.NewStateStorageWithCache(db)
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, tdb)
	processor.economyModel = fakeEconomyModel{}
	assert.NoError(t, err)

	block := createBlock(10)
	err = processor.doRewards(block)
	assert.Error(t, err)

	var earlyTokenContract contract.EarlyRewardContract
	err = processor.PutContract(contract.EarlyContractAddress, reflect.ValueOf(&earlyTokenContract))
	assert.NoError(t, err)
	block = createBlock(2)
	err = processor.doRewards(block)
	assert.Equal(t, gerror.NotHavePreBlockErr, err)
}