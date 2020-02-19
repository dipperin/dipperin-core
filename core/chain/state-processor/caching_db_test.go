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
	trie2 "github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStateStorageWithCache(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)

	trie, err := storage.OpenTrie(common.Hash{})
	assert.NoError(t, err)
	assert.NotNil(t, trie)

	root, err := trie.Commit(nil)
	assert.NoError(t, err)

	trie, err = storage.OpenTrie(root)
	assert.NoError(t, err)
	assert.NotNil(t, trie)

	err = trie.Prove([]byte{}, 0, db)
	assert.NoError(t, err)

	trie, err = storage.OpenTrie(common.HexToHash("123"))
	assert.Error(t, err)
	assert.Nil(t, trie)

	storageTrie, err := storage.OpenStorageTrie(common.Hash{}, common.Hash{})
	assert.NoError(t, err)
	assert.NotNil(t, storageTrie)

}

func TestCachingDB_CopyTrie(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)

	trie, err := storage.OpenTrie(common.Hash{})
	assert.NoError(t, err)
	assert.NotNil(t, trie)
	assert.NotNil(t, storage.CopyTrie(trie))

	secureTrie, err := trie2.NewSecure(common.Hash{}, storage.TrieDB(), 0)
	assert.NoError(t, err)
	assert.NotNil(t, secureTrie)
	assert.NotNil(t, storage.CopyTrie(secureTrie))
}

func TestAccount_Contract(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)
	cache := storage.(*cachingDB)

	hash := common.HexToHash("123")
	blob := []byte{1}

	code, err := cache.ContractCode(common.Hash{}, hash)
	assert.Error(t, err)
	assert.Nil(t, code)

	num, err := cache.ContractCodeSize(common.Hash{}, hash)
	assert.Error(t, err)
	assert.Equal(t, 0, num)

	cache.db.InsertBlob(hash, blob)
	code, err = cache.ContractCode(common.Hash{}, hash)
	assert.NoError(t, err)
	assert.NotNil(t, code)

	num, err = cache.ContractCodeSize(common.Hash{}, hash)
	assert.NoError(t, err)
	assert.Equal(t, 1, num)
}

func TestCachedTrie_Commit(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)
	cache := storage.(*cachingDB)

	trie, err := storage.OpenTrie(common.Hash{})
	assert.NoError(t, err)
	assert.NotNil(t, trie)

	for i := 0; i < maxPastTries+1; i++ {
		_, err = trie.Commit(nil)
		assert.NoError(t, err)
	}
	assert.Len(t, cache.pastTries, maxPastTries)
}
