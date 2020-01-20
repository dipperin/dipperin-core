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
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestFullstate(t *testing.T) {
	db := ethdb.NewMemDatabase()
	processor, _ := NewAccountStateDB(common.Hash{}, NewStateStorageWithCache(db))
	state := NewFullState(processor)

	// set alice info
	state.CreateAccount(aliceAddr)
	state.AddLog(&model.Log{TxHash: common.HexToHash("txHash")})

	// snap shot
	snapShot := state.Snapshot()

	state.AddBalance(aliceAddr, big.NewInt(500))
	state.SubBalance(aliceAddr, big.NewInt(200))
	state.AddNonce(aliceAddr, 1)
	state.AddLog(&model.Log{TxHash: common.HexToHash("txHash")})
	state.SetAbi(aliceAddr, []byte{123})
	state.SetCode(aliceAddr, []byte{234})
	state.SetState(aliceAddr, []byte("key"), []byte("value"))

	// assert alice info
	assert.Equal(t, true, state.Exist(aliceAddr))
	assert.Equal(t, big.NewInt(300), state.GetBalance(aliceAddr))
	nonce, err := state.GetNonce(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), nonce)

	logs := state.GetLogs(common.HexToHash("txHash"))
	assert.Equal(t, 2, len(logs))
	assert.Equal(t, []byte{123}, state.GetAbi(aliceAddr))
	assert.Equal(t, []byte{234}, state.GetCode(aliceAddr))
	assert.Equal(t, cs_crypto.Keccak256Hash([]byte{123}), state.GetAbiHash(aliceAddr))
	assert.Equal(t, cs_crypto.Keccak256Hash([]byte{234}), state.GetCodeHash(aliceAddr))
	assert.Equal(t, []byte("value"), state.GetState(aliceAddr, []byte("key")))

	// revert to snap shot
	state.RevertToSnapshot(snapShot)
	assert.Equal(t, true, state.Exist(aliceAddr))
	assert.Equal(t, big.NewInt(0), state.GetBalance(aliceAddr))
	nonce, err = state.GetNonce(aliceAddr)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), nonce)

	logs = state.GetLogs(common.HexToHash("txHash"))
	assert.Equal(t, 1, len(logs))
	assert.Equal(t, []byte{}, state.GetAbi(aliceAddr))
	assert.Equal(t, []byte{}, state.GetCode(aliceAddr))
	assert.Equal(t, vm.EmptyCodeHash, state.GetAbiHash(aliceAddr))
	assert.Equal(t, vm.EmptyCodeHash, state.GetCodeHash(aliceAddr))
	assert.Equal(t, []byte(nil), state.GetState(aliceAddr, []byte("key")))
}

func TestNewFullState_Error(t *testing.T) {
	db := ethdb.NewMemDatabase()
	processor, _ := NewAccountStateDB(common.Hash{}, NewStateStorageWithCache(db))
	state := NewFullState(processor)
	assert.Panics(t, func() { state.SetAbi(aliceAddr, nil) })
	assert.Panics(t, func() { state.SetCode(aliceAddr, nil) })
	assert.Panics(t, func() { state.AddLog(nil) })
	assert.Equal(t, common.Hash{}, state.GetCodeHash(aliceAddr))
	assert.Equal(t, common.Hash{}, state.GetAbiHash(aliceAddr))
	assert.Nil(t, state.GetCode(aliceAddr))
	assert.Nil(t, state.GetAbi(aliceAddr))
}
