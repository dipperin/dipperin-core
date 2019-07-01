package state_processor

import (
	"fmt"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/dipperin/dipperin-core/tests/g-testData"
)

func TestApplyMessage(t *testing.T) {
	WASMPath := g_testData.GetWasmPath("event")
	AbiPath := g_testData.GetAbiPath("event")
	tx := createContractTx(t, WASMPath, AbiPath)
	msg, err := tx.AsMessage()
	assert.NoError(t, err)

	testVm := getTestVm()
	gasPool := uint64(5 * testGasLimit)
	result, usedGas, failed, _, err := ApplyMessage(testVm, msg, &gasPool)
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

	result, usedGas, failed, _, err = ApplyMessage(testVm, msg, &gasPool)
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)
	assert.Equal(t, make([]byte, 32), result)
}
