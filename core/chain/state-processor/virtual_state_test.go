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
	"testing"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
)

func TestManageState(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)
	accountDB, err := NewAccountStateDB(common.Hash{}, storage)
	assert.NoError(t, err)

	manageState := ManageState(accountDB)
	assert.NotNil(t, manageState)

	manageState.SetState(accountDB)
	nonce := manageState.NewNonce(aliceAddr)
	assert.Equal(t, uint64(0), nonce)
	assert.True(t, manageState.HasAccount(aliceAddr))

	nonce = manageState.NewNonce(aliceAddr)
	assert.Equal(t, uint64(1), nonce)

	manageState.SetNonce(aliceAddr, 5)
	nonce = manageState.GetNonce(aliceAddr)
	assert.Equal(t, uint64(5), nonce)

	manageState.RemoveNonce(aliceAddr, 5)
	nonce = manageState.GetNonce(aliceAddr)
	assert.Equal(t, uint64(5), nonce)

	nonce = manageState.GetNonce(bobAddr)
	assert.Equal(t, uint64(0), nonce)

	manageState.SetNonce(common.HexToAddress("123"), 5)
	nonce = manageState.GetNonce(common.HexToAddress("123"))
	assert.Equal(t, uint64(5), nonce)
}

func TestManagedState_NewNonce(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)
	accountDB, err := NewAccountStateDB(common.Hash{}, storage)
	assert.NoError(t, err)

	err = accountDB.NewAccountState(aliceAddr)
	assert.NoError(t, err)

	manageState := ManageState(accountDB)
	assert.NotNil(t, manageState)
	nonce := manageState.NewNonce(aliceAddr)
	assert.Equal(t, uint64(0), nonce)

	manageState.AddNonce(aliceAddr, 5)
	nonce = manageState.NewNonce(aliceAddr)
	assert.Equal(t, uint64(5), nonce)

	manageState.getAccount(aliceAddr).nonces[0] = false
	nonce = manageState.NewNonce(aliceAddr)
	assert.Equal(t, uint64(5), nonce)
}
