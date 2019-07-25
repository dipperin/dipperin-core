package state_processor

import (
	"fmt"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyMessage(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	tx := createContractTx(WASMPath, AbiPath, 0)
	msg, err := tx.AsMessage()
	assert.NoError(t, err)

	testVm := getTestVm()
	gasPool := uint64(5 * testGasLimit)
	result, usedGas, failed, _, err := ApplyMessage(testVm, &msg, &gasPool)
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)

	fmt.Println("----------------------------------")

	name := []byte("ApplyMsg")
	params := [][]byte{name}
	to := cs_crypto.CreateContractAddress(aliceAddr, 0)
	tx = callContractTx(&to, "returnString", params, 1)
	msg, err = tx.AsMessage()
	assert.NoError(t, err)

	result, usedGas, failed, _, err = ApplyMessage(testVm, &msg, &gasPool)
	resp := utils.Align32BytesConverter(result, "string")
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)
	assert.Equal(t, string(name), resp)
}

func BenchmarkApplyMessage_Create(b *testing.B) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	tx := createContractTx(WASMPath, AbiPath, 0)
	msg, err := tx.AsMessage()
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		testVm := getTestVm()
		gasPool := uint64(5 * testGasLimit)
		_, usedGas, failed, _, err := ApplyMessage(testVm, &msg, &gasPool)
		assert.NoError(b, err)
		assert.False(b, failed)
		assert.NotNil(b, usedGas)
	}

	fmt.Println("----------------------------------")
}

func BenchmarkApplyMessage_Call(b *testing.B) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)

	// create tx
	tx1 := createContractTx(WASMPath, AbiPath, 0)
	msg1, err := tx1.AsMessage()
	assert.NoError(b, err)

	// call tx
	name := []byte("ApplyMsg")
	num := utils.Int64ToBytes(234)
	params := [][]byte{name, num}
	to := cs_crypto.CreateContractAddress(aliceAddr, 0)
	tx2 := callContractTx(&to, "hello", params, 1)
	msg2, err := tx2.AsMessage()
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		testVm := getTestVm()
		gasPool := uint64(5 * testGasLimit)

		_, usedGas, failed, _, innerErr := ApplyMessage(testVm, &msg1, &gasPool)
		assert.NoError(b, innerErr)
		assert.False(b, failed)
		assert.NotNil(b, usedGas)

		fmt.Println("----------------------------------")

		_, usedGas, failed, _, innerErr = ApplyMessage(testVm, &msg2, &gasPool)
		assert.NoError(b, innerErr)
		assert.False(b, failed)
		assert.NotNil(b, usedGas)
	}
}
