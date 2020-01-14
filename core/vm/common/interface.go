// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package common

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"math/big"
)

//go:generate mockgen -destination=./../../tests/mock/vm/state_db_mock.go -package=vm_mock github.com/dipperin/dipperin-core/core/vm StateDB
type StateDB interface {
	GetBalance(common.Address) *big.Int
	CreateAccount(common.Address) error
	GetNonce(common.Address) (uint64, error)
	AddNonce(common.Address, uint64)
	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetAbiHash(common.Address) common.Hash
	GetAbi(common.Address) []byte
	SetAbi(common.Address, []byte)
	AddBalance(addr common.Address, amount *big.Int) error
	SubBalance(addr common.Address, amount *big.Int) error
	// todo: hash -> bytes
	GetState(common.Address, []byte) []byte
	SetState(common.Address, []byte, []byte)
	AddLog(addedLog *model.Log)
	GetLogs(txHash common.Hash) []*model.Log

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	RevertToSnapshot(int)
	Snapshot() int
}
