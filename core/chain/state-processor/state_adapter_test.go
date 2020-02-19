package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
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
	//assert.Error(t, func() {state.AddBalance(aliceAddr, big.NewInt(100)) })
	//assert.Panics(t, func() { state.SubBalance(aliceAddr, big.NewInt(100)) })
	//assert.Panics(t, func() { state.GetBalance(aliceAddr) })
	//assert.Panics(t, func() { state.AddNonce(aliceAddr, uint64(1)) })
	assert.Panics(t, func() { state.SetAbi(aliceAddr, nil) })
	assert.Panics(t, func() { state.SetCode(aliceAddr, nil) })
	assert.Panics(t, func() { state.AddLog(nil) })
	assert.Equal(t, common.Hash{}, state.GetCodeHash(aliceAddr))
	assert.Equal(t, common.Hash{}, state.GetAbiHash(aliceAddr))
	assert.Nil(t, state.GetCode(aliceAddr))
	assert.Nil(t, state.GetAbi(aliceAddr))
}
