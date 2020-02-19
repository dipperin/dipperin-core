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
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

var minDiff = common.HexToDiff("0x20ffffff")

func TestNewBlockProcessor(t *testing.T) {
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)
	assert.NotNil(t, processor)

	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{storageErr: TrieError})
	assert.Equal(t, TrieError, err)
	assert.Nil(t, processor)
}

func TestBlockProcessor_Process(t *testing.T) {
	db, root := createTestStateDB(t)
	tdb := state_processor.NewStateStorageWithCache(db)
	state, err := state_processor.NewAccountStateDB(root, tdb)
	assert.NoError(t, err)
	processor, err := NewBlockProcessor(fakeAccountDBChain{state: state}, root, tdb)
	assert.NoError(t, err)

	// block.Number() == 0 return
	block := model.CreateBlock(0, common.Hash{}, 0)
	err = processor.Process(block, fakeEconomyModel{})
	assert.NoError(t, err)

	// blockNum = 9 isn't change point
	block = createBlock(10)
	err = processor.Process(block, fakeEconomyModel{})
	assert.Equal(t, g_error.ErrAccountNotExist, err)

	// blockNum = 19 is change point
	block = model.CreateBlock(20, common.Hash{}, 0)
	err = processor.Process(block, fakeEconomyModel{})
	assert.NoError(t, err)
}

func TestBlockProcessor_Process_Error(t *testing.T) {
	header := model.NewHeader(1, 10, common.Hash{}, common.HexToHash("1111"), minDiff, big.NewInt(324234), common.Address{}, common.BlockNonceFromInt(432423))

	tx := createUnNormalTx()
	block := model.NewBlock(header, []*model.Transaction{tx}, nil)

	db, root := createTestStateDB(t)
	tdb := state_processor.NewStateStorageWithCache(db)
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, root, tdb)
	assert.NoError(t, err)

	// TxIterator inner error
	err = processor.Process(block, fakeEconomyModel{})
	assert.Equal(t, g_error.ErrUnknownTxType, err)

	// doReward- RewardCoinBase no coin base error
	block = model.NewBlock(header, nil, nil)
	err = processor.Process(block, fakeEconomyModel{})
	assert.Equal(t, g_error.InvalidCoinBaseAddressErr, err)
}

func TestBlockProcessor_processorCommitList_Error(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := state_processor.NewStateStorageWithCache(db)
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, tdb)
	assert.NoError(t, err)

	block := createBlock(1)
	err = processor.processCommitList(block, false)
	assert.Equal(t, g_error.NotHavePreBlockErr, err)

	block = createBlock(20)
	err = processor.processCommitList(block, false)
	assert.Equal(t, g_error.ErrAccountNotExist, err)

	for i := 0; i < len(VerifierAddress); i++ {
		err = processor.NewAccountState(VerifierAddress[i])
		assert.NoError(t, err)
	}

	err = processor.processCommitList(block, false)
	assert.Equal(t, g_error.ErrAccountNotExist, err)
}

func TestBlockProcessor_doReward_Error(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := state_processor.NewStateStorageWithCache(db)
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
	assert.Equal(t, g_error.NotHavePreBlockErr, err)
}
