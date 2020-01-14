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
	"sync"
)

type virtualAccount struct {
	account *account
	nstart  uint64
	nonces  []bool
}

type ManagedState struct {
	*AccountStateDB

	mu sync.RWMutex

	accounts map[common.Address]*virtualAccount
}

// ManagedState returns a new managed state with the statedb as it's backing layer
func ManageState(statedb *AccountStateDB) *ManagedState {
	return &ManagedState{
		AccountStateDB: statedb.Copy(),
		accounts:       make(map[common.Address]*virtualAccount),
	}
}

// SetState sets the backing layer of the managed state
func (ms *ManagedState) SetState(statedb *AccountStateDB) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.AccountStateDB = statedb
}

// RemoveNonce removed the nonce from the managed state and all future pending nonces
func (ms *ManagedState) RemoveNonce(addr common.Address, n uint64) {
	if !ms.IsEmptyAccount(addr) {
		ms.mu.Lock()
		defer ms.mu.Unlock()

		account := ms.getAccount(addr)
		if n-account.nstart <= uint64(len(account.nonces)) {
			reslice := make([]bool, n-account.nstart)
			copy(reslice, account.nonces[:n-account.nstart])
			account.nonces = reslice
		}
	}
}

// NewNonce returns the new canonical nonce for the managed account
func (ms *ManagedState) NewNonce(addr common.Address) uint64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	account := ms.getAccount(addr)
	for i, nonce := range account.nonces {
		if !nonce {
			return account.nstart + uint64(i)
		}
	}
	account.nonces = append(account.nonces, true)

	return uint64(len(account.nonces)-1) + account.nstart
}

// GetNonce returns the canonical nonce for the managed or unmanaged account.
//
// Because GetNonce mutates the DB, we must take a write lock.
func (ms *ManagedState) GetNonce(addr common.Address) uint64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.hasAccount(addr) {
		account := ms.getAccount(addr)
		return uint64(len(account.nonces)) + account.nstart
	} else {
		nonce, _ := ms.AccountStateDB.GetNonce(addr)
		return nonce
	}
}

// SetNonce sets the new canonical nonce for the managed state
func (ms *ManagedState) SetNonce(addr common.Address, nonce uint64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.IsEmptyAccount(addr) {
		//todo may need modify the state change list
		so, _ := ms.newAccountState(addr)
		so.setNonce(nonce)
		ms.accounts[addr] = newAccount(so)
	} else {
		so, _ := ms.GetAccountState(addr)
		so.setNonce(nonce)
		ms.accounts[addr] = newAccount(so)
	}

}

// HasAccount returns whether the given address is managed or not
func (ms *ManagedState) HasAccount(addr common.Address) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.hasAccount(addr)
}

func (ms *ManagedState) hasAccount(addr common.Address) bool {
	_, ok := ms.accounts[addr]
	return ok
}

// populate the managed state
func (ms *ManagedState) getAccount(addr common.Address) *virtualAccount {
	if account, ok := ms.accounts[addr]; !ok {
		if ms.IsEmptyAccount(addr) {
			//todo may need modify the state change list
			so, _ := ms.newAccountState(addr)
			ms.accounts[addr] = newAccount(so)
		} else {
			so, _ := ms.GetAccountState(addr)
			ms.accounts[addr] = newAccount(so)
		}
	} else {
		// Always make sure the state account nonce isn't actually higher
		// than the tracked one.
		so, _ := ms.GetAccountState(addr)
		if so != nil && uint64(len(account.nonces))+account.nstart < so.getNonce() {
			ms.accounts[addr] = newAccount(so)
		}

	}
	return ms.accounts[addr]
}

func newAccount(account *account) *virtualAccount {
	return &virtualAccount{account, account.getNonce(), nil}
}
