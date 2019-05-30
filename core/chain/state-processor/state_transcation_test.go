package state_processor

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
	"math/big"
	"github.com/dipperin/dipperin-core/common/vmcommon"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
)

func TestApplyMessage(t *testing.T) {
	var testPath = "/home/qydev/go/src/github.com/dipperin/dipperin-core/core/vm/event"
	tx := createContractTx(t, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	msg, err := tx.AsMessage()
	assert.NoError(t, err)

	a := make(map[common.Address]*big.Int)
	c := make(map[common.Address][]byte)
	testVm := getTestVm(a, c)

	gasPool := uint64(5*gasLimit)
	result, usedGas, failed, err := ApplyMessage(testVm, msg, &gasPool)
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)
	data := getContractCode(t, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	assert.Equal(t, data, result)

	fmt.Println("----------------------------------")

	name := []byte("ApplyMsg")
	num := vmcommon.Int64ToBytes(234)
	params := [][]byte{name, num}
	to := cs_crypto.CreateContractAddress(aliceAddr, 0)
	tx = callContractTx(t, &to, "hello", params)
	msg, err = tx.AsMessage()
	assert.NoError(t, err)

	result, usedGas, failed, err = ApplyMessage(testVm, msg, &gasPool)
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)
	assert.Equal(t, make([]byte, 32), result)
}
