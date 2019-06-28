package state_processor

import (
	"fmt"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyMessage(t *testing.T) {
	var testPath = "/home/qydev/go/src/github.com/dipperin/dipperin-core/core/vm/event"
	tx := createContractTx(t, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	msg, err := tx.AsMessage()
	assert.NoError(t, err)

	testVm := getTestVm()

	gasPool := uint64(5*gasLimit)
	result, usedGas, failed,_, err := ApplyMessage(testVm, msg, &gasPool)
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)

	fmt.Println("----------------------------------")

	name := []byte("ApplyMsg")
	num := utils.Int64ToBytes(234)
	params := [][]byte{name, num}
	to := cs_crypto.CreateContractAddress(aliceAddr, 0)
	tx = callContractTx(t, &to, "hello", params, 0)
	msg, err = tx.AsMessage()
	assert.NoError(t, err)

	result, usedGas, failed,_, err = ApplyMessage(testVm, msg, &gasPool)
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)
	assert.Equal(t, make([]byte, 32), result)
}
