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
	"github.com/stretchr/testify/assert"
	"testing"
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

	type result struct {
		err1 error
		err2 error
	}

	testCases := []struct {
		name   string
		given  func() (error,error)
		expect result
	}{
		{
			name:"no error",
			given: func() (error,error) {
				err1 := processor.ProcessVerification(verify, 0)
				err2 := processor.ProcessVerifierNumber(aliceAddr)
				return err1, err2
			},
			expect:result{nil,nil},
		},
		{
			name:"ErrAccountNotExist",
			given: func() (error, error) {
				processor, err = NewAccountStateDB(common.Hash{}, fakeStateStorage{getErr: TrieError})
				assert.NoError(t, err)
				err1 := processor.ProcessVerification(verify, 0)
				err2 := processor.ProcessVerifierNumber(common.HexToAddress("234"))
				return err1, err2
			},
			expect:result{gerror.ErrAccountNotExist,gerror.ErrAccountNotExist},
		},
	}

	for _,tc:=range testCases{
		err1,err2:=tc.given()
		if err1!=nil && err2!=nil{
			assert.Equal(t,tc.expect.err1,err1)
			assert.Equal(t,tc.expect.err2,err2)
		}else {
			assert.NoError(t,err1)
			assert.NoError(t,err2)
		}
	}
}

func TestAccountStateDB_ProcessPerformance(t *testing.T) {
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))
	assert.NoError(t, err)

	aliceOriginalPerformance, err := processor.GetPerformance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(30), aliceOriginalPerformance)

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"ErrAccountNotExist",
			given: func() error {
				err := processor.ProcessPerformance(common.HexToAddress("123"), reward)
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
		{
			name:"no error",
			given: func() error {
				err := processor.ProcessPerformance(aliceAddr, reward)
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

	performance, err := processor.GetPerformance(aliceAddr)
	assert.NoError(t, err)
	assert.EqualValues(t, uint64(31), performance)
}

