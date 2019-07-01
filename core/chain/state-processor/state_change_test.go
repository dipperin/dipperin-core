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
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"testing"
	"math/big"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common/g-error"
)

func TestStateChange(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor.NewAccountState(aliceAddr)
	processor.NewAccountState(bobAddr)
	for i := 0; i < 100; i++ {
		processor.AddBalance(aliceAddr, big.NewInt(int64(i)))
	}
	for j := 0; j < 100; j++ {
		processor.AddNonce(bobAddr, uint64(j))
	}
	enc, _ := rlp.EncodeToBytes(big.NewInt(1000000))
	processor.blockStateTrie.TryUpdate(GetBalanceKey(aliceAddr), enc)
	//processor.stateChangeList.recover(processor)
	newscl := processor.stateChangeList.digest()
	processor.stateChangeList = newscl
	processor.stateChangeList.recover(processor)
	alicebalance, err := processor.GetBalance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, alicebalance, big.NewInt(4950))
	bobnonce, err := processor.GetNonce(bobAddr)
	assert.NoError(t, err)
	assert.Equal(t, bobnonce, uint64(4950))
}

func TestStateChangeList_DecodeRLP(t *testing.T) {
	db := ethdb.NewMemDatabase()
	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor.NewAccountState(aliceAddr)
	processor.NewAccountState(bobAddr)
	processor.NewAccountState(charlieAddr)
	processor.AddBalance(charlieAddr, big.NewInt(500))
	processor.DeleteAccountState(charlieAddr)
	for i := 0; i < 100; i++ {
		processor.AddBalance(aliceAddr, big.NewInt(int64(i)))
	}
	for j := 0; j < 100; j++ {
		processor.AddNonce(bobAddr, uint64(j))
	}
	sclsent := processor.stateChangeList.digest()

	enc, err := rlp.EncodeToBytes(sclsent)
	assert.NoError(t, err)

	var sclget StateChangeList
	err2 := rlp.DecodeBytes(enc, &sclget)
	assert.NoError(t, err2)

	processor2, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor2.stateChangeList = &sclget
	processor2.stateChangeList.recover(processor2)
	_, charlieerr := processor2.GetAccountState(charlieAddr)
	assert.Equal(t, charlieerr, g_error.AccountNotExist)
	alicebalance, err := processor2.GetBalance(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, alicebalance, big.NewInt(4950))
	bobnonce, err := processor2.GetNonce(bobAddr)
	assert.NoError(t, err)
	assert.Equal(t, bobnonce, uint64(4950))
}
