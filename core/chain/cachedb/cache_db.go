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
	"encoding/binary"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

var cacheDataDecoder CacheDataDecoder = &BFTCacheDataDecoder{}

func SetCacheDataDecoder(d CacheDataDecoder) {
	cacheDataDecoder = d
}

func NewCacheDB(db ethdb.Database) *CacheDB {
	return &CacheDB{
		db: db,
	}
}

type CacheDB struct {
	db ethdb.Database
}

// hash must be empty, if only use height for tag
func (cache *CacheDB) GetSeenCommits(blockHeight uint64, blockHash common.Hash) (result []model.AbstractVerification, err error) {
	var data []byte
	data, err = cache.get(seenCommitsKey(blockHeight, blockHash))
	if err != nil {
		log.Info("get seen commits failed", "height", blockHeight)
		return
	}
	if blockHeight > 1 {
		if len(data) == 0 {
			pbft_log.Debug("Can not get seen commits", "height", blockHeight)
			return nil, errors.New("Can not get seen commits")
		}
	}
	return cacheDataDecoder.DecodeSeenCommits(data)
}

// hash must be empty, if only use height for tag
func (cache *CacheDB) SaveSeenCommits(blockHeight uint64, blockHash common.Hash, commits []model.AbstractVerification) error {
	return cache.save(seenCommitsKey(blockHeight, blockHash), commits)
}

func (cache *CacheDB) save(key []byte, data interface{}) error {
	dataB, err := rlp.EncodeToBytes(data)
	if err != nil {
		return err
	}
	return cache.db.Put(key, dataB)
}

// get bytes
func (cache *CacheDB) get(key []byte) ([]byte, error) {
	return cache.db.Get(key)
}

// decode to result
func (cache *CacheDB) load(key []byte, result interface{}) error {
	dataB, err := cache.db.Get(key)
	if err != nil {
		return err
	}
	return rlp.DecodeBytes(dataB, result)
}

func seenCommitsKey(blockHeight uint64, blockHash common.Hash) []byte {
	return append(append([]byte("seen_commits"), encodeNumber(blockHeight)...), blockHash.Bytes()..., )
}

func encodeNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}
