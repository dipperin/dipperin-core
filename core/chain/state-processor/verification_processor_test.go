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
	"github.com/dipperin/dipperin-core/core/model"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common/g-error"
)

func TestAccountStateDB_ProcessVerification(t *testing.T) {
	// Create original data
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	// createVote(height, viewId, sequenceId uint64, blockHash common.Hash, address common.Address, privKey string)
	verify := createSignedVote(1, common.HexToHash("123456"), model.VoteMessage, testPriv1, aliceAddr)
	aliceOriginalVerifyNum, _ := processor.GetVerifyNum(aliceAddr)
	assert.EqualValues(t, uint64(0), aliceOriginalVerifyNum)
	aliceOriginalCommitNum, _ := processor.GetCommitNum(aliceAddr)
	assert.EqualValues(t, uint64(0), aliceOriginalCommitNum)
	aliceOriginalPerformance, _ := processor.GetPerformance(aliceAddr)
	assert.EqualValues(t, uint64(30), aliceOriginalPerformance)

	// Valid
	err = processor.ProcessVerification(verify,0)
	assert.NoError(t, err)
	err = processor.ProcessVerifierNumber(aliceAddr)
	assert.NoError(t, err)
	aliceNewCommitNum, err := processor.GetCommitNum(aliceAddr)
	assert.NoError(t, err)
	aliceNewVerifyNum, err := processor.GetVerifyNum(aliceAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, aliceOriginalVerifyNum + 1, aliceNewVerifyNum)
	assert.EqualValues(t, aliceOriginalCommitNum + 1, aliceNewCommitNum)

	// Invalid
	processor, err = NewAccountStateDB(common.Hash{}, fakeStateStorage{getErr:TrieError})
	assert.NoError(t, err)

	err = processor.ProcessVerification(verify,0)
	assert.Equal(t, g_error.AccountNotExist, err)
	err = processor.ProcessVerifierNumber(common.HexToAddress("234"))
	assert.Equal(t, g_error.AccountNotExist, err)
}

func TestAccountStateDB_ProcessPerformance(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	aliceOriginalPerformance, err := processor.GetPerformance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(30), aliceOriginalPerformance)

	err = processor.ProcessPerformance(common.HexToAddress("123"), reward)
	assert.Equal(t, g_error.AccountNotExist, err)

	err = processor.ProcessPerformance(aliceAddr, reward)
	assert.NoError(t, err)

	performance, err := processor.GetPerformance(aliceAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, uint64(31), performance)
}
