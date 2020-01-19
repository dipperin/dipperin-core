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
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManageState(t *testing.T) {
	db := ethdb.NewMemDatabase()
	storage := NewStateStorageWithCache(db)
	accountDB, err := NewAccountStateDB(common.Hash{}, storage)
	assert.NoError(t, err)

	manageState := ManageState(accountDB)
	assert.NotNil(t, manageState)
	manageState.SetState(accountDB)

	type result struct {
		nonce uint64
	}

	testCases := []struct {
		name   string
		given  func() uint64
		expect result
	}{
		{
			name:"new nonce",
			given: func() uint64 {
				nonce := manageState.NewNonce(aliceAddr)
				return nonce
			},
			expect:result{uint64(0)},
		},
		{
			name:"nonce grow up",
			given: func() uint64 {
				nonce := manageState.NewNonce(aliceAddr)
				return nonce
			},
			expect:result{uint64(1)},
		},
		{
			name:"set nonce 5",
			given: func() uint64 {
				manageState.SetNonce(aliceAddr, 5)
				nonce := manageState.GetNonce(aliceAddr)
				return nonce
			},
			expect:result{uint64(5)},
		},
		{
			name:"remove nonce 5",
			given: func() uint64 {
				manageState.RemoveNonce(aliceAddr, 5)
				nonce := manageState.GetNonce(aliceAddr)
				return nonce
			},
			expect:result{uint64(5)},
		},
		{
			name:"change another address's nonce",
			given: func() uint64 {
				nonce := manageState.GetNonce(bobAddr)
				return nonce
			},
			expect:result{uint64(0)},
		},
		{
			name:"set a new address's nonce",
			given: func() uint64 {
				manageState.SetNonce(common.HexToAddress("123"), 5)
				nonce := manageState.GetNonce(common.HexToAddress("123"))
				return nonce
			},
			expect:result{uint64(5)},
		},
	}

	for _,tc:=range testCases{
		nonce:=tc.given()
		assert.Equal(t,tc.expect.nonce,nonce)
	}
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

	type result struct {
		nonce uint64
	}

	testCases := []struct {
		name   string
		given  func() uint64
		expect result
	}{
		{
			name:"nonce 0",
			given: func() uint64 {
				nonce := manageState.NewNonce(aliceAddr)
				return nonce
			},
			expect:result{uint64(0)},
		},
		{
			name:"set nonce 5",
			given: func() uint64 {
				manageState.AddNonce(aliceAddr, 5)
				nonce := manageState.NewNonce(aliceAddr)
				return nonce
			},
			expect:result{uint64(5)},
		},
		{
			name:"set nonce getAccount false",
			given: func() uint64 {
				manageState.getAccount(aliceAddr).nonces[0] = false
				nonce := manageState.NewNonce(aliceAddr)
				return nonce
			},
			expect:result{uint64(5)},
		},
	}

	for _,tc:=range testCases{
		nonce:=tc.given()
		assert.Equal(t,tc.expect.nonce,nonce)
	}
}

