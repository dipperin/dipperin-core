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

package cachedb

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"errors"
)

var (
	TestError = errors.New("test error")
)

func TestCacheDB_GetSeenCommits(t *testing.T) {
	SetCacheDataDecoder(&BFTCacheDataDecoder{})

	db := ethdb.NewMemDatabase()
	cacheDB := NewCacheDB(db)

	commits := createCommits()
	err := cacheDB.SaveSeenCommits(1, common.Hash{}, commits)
	assert.NoError(t, err)

	vs, err := cacheDB.GetSeenCommits(1, common.Hash{})
	assert.NoError(t, err)
	assert.Len(t, vs, 1)

	var result []*model.VoteMsg
	err = cacheDB.load(seenCommitsKey(1, common.Hash{}), &result)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCacheDB_GetSeenCommits_Error(t *testing.T) {
	SetCacheDataDecoder(&BFTCacheDataDecoder{})
	cacheDB := NewCacheDB(fakeDataBase{err: TestError})

	vs, err := cacheDB.GetSeenCommits(1, common.Hash{})
	assert.Equal(t, TestError, err)
	assert.Nil(t, vs)

	var result []*model.VoteMsg
	err = cacheDB.load(seenCommitsKey(1, common.Hash{}), result)
	assert.Equal(t, TestError, err)

	err = cacheDB.save(seenCommitsKey(3, common.Hash{}), 3)
	assert.Error(t, err)

	cacheDB = NewCacheDB(fakeDataBase{})

	vs, err = cacheDB.GetSeenCommits(3, common.Hash{})
	assert.Error(t, err)
	assert.Nil(t, vs)
}

func TestBFTCacheDataDecoder_DecodeSeenCommits(t *testing.T) {
	decoder := BFTCacheDataDecoder{}
	result, err := decoder.DecodeSeenCommits([]byte{})
	assert.Error(t, err)
	assert.Equal(t, []model.AbstractVerification{}, result)

	result, err = decoder.DecodeSeenCommits([]byte{1, 2, 3})
	assert.Error(t, err)
	assert.Equal(t, []model.AbstractVerification{}, result)
}

func createCommits() []model.AbstractVerification {
	return []model.AbstractVerification{
		&model.VoteMsg{
			Height:    1,
			Timestamp: time.Now(),
			Witness:   &model.WitMsg{},
		},
	}
}

type fakeDataBase struct {
	err error
}

func (data fakeDataBase) Put(key []byte, value []byte) error {
	return data.err
}

func (data fakeDataBase) Delete(key []byte) error {
	panic("implement me")
}

func (data fakeDataBase) Get(key []byte) ([]byte, error) {
	return nil, data.err
}

func (data fakeDataBase) Has(key []byte) (bool, error) {
	panic("implement me")
}

func (data fakeDataBase) Close() {
	panic("implement me")
}

func (data fakeDataBase) NewBatch() ethdb.Batch {
	panic("implement me")
}
