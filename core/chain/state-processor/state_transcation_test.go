package state_processor

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestApplyMessage(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	tx := createContractTx(WASMPath, AbiPath, 0, testGasLimit)
	msg, err := tx.AsMessage(false)
	assert.NoError(t, err)

	db, root := CreateTestStateDB()
	testVm := getTestVm(db, root)
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
	msg, err = tx.AsMessage(false)
	assert.NoError(t, err)

	result, usedGas, failed, _, err = ApplyMessage(testVm, &msg, &gasPool)
	resp := utils.Align32BytesConverter(result, "string")
	assert.NoError(t, err)
	assert.Equal(t, false, failed)
	assert.NotNil(t, usedGas)
	assert.Equal(t, string(name), resp)
}

func TestApplyMessage_Error(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	gasLimit := g_testData.TestGasLimit
	tx := createContractTx(WASMPath, AbiPath, 0, gasLimit)
	msg, err := tx.AsMessage(true)
	assert.NoError(t, err)

	db := ethdb.NewMemDatabase()
	testVm := getTestVm(db, common.Hash{})
	gasPool := gasLimit / 2
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	assert.Equal(t, "account does not exist", err.Error())

	tdb := NewStateStorageWithCache(db)
	processor, _ := NewAccountStateDB(common.Hash{}, tdb)
	processor.NewAccountState(aliceAddr)
	processor.AddNonce(aliceAddr, 1)
	root, _ := processor.Commit()
	testVm = getTestVm(db, root)
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	assert.Equal(t, g_error.ErrNonceTooLow, err)

	tx = createContractTx(WASMPath, AbiPath, 2, gasLimit)
	msg, err = tx.AsMessage(true)
	assert.NoError(t, err)
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	assert.Equal(t, g_error.ErrNonceTooHigh, err)

	tx = createContractTx(WASMPath, AbiPath, 1, gasLimit)
	msg, err = tx.AsMessage(true)
	assert.NoError(t, err)
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	assert.Equal(t, g_error.ErrInsufficientBalanceForGas, err)

	processor, _ = NewAccountStateDB(root, tdb)
	processor.AddBalance(aliceAddr, big.NewInt(0).SetUint64(7000000))
	root, _ = processor.Commit()
	testVm = getTestVm(db, root)
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	assert.Equal(t, g_error.ErrGasLimitReached, err)

	gasPool = gasLimit
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	assert.Equal(t, g_error.ErrOutOfGas, err)

	gasLimit = g_testData.TestGasLimit * 50
	gasPool = gasLimit * 10
	tx = createContractTx(WASMPath, AbiPath, 1, gasLimit)
	msg, err = tx.AsMessage(true)
	assert.NoError(t, err)
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	assert.NoError(t, err)

	name := []byte("ApplyMsg")
	params := [][]byte{name}
	to := cs_crypto.CreateContractAddress(aliceAddr, 1)
	tx = callContractTx(&to, "returnString", params, 2)
	msg, err = tx.AsMessage(true)
	assert.NoError(t, err)
	_, _, _, _, err = ApplyMessage(testVm, &msg, &gasPool)
	//assert.Equal(t, g_error.ErrInsufficientBalance, err)
	assert.NoError(t, err)
}

func BenchmarkApplyMessage_Create(b *testing.B) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	tx := createContractTx(WASMPath, AbiPath, 0, testGasLimit)
	msg, err := tx.AsMessage(true)
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		db, root := CreateTestStateDB()
		testVm := getTestVm(db, root)
		gasPool := uint64(5 * testGasLimit)
		_, usedGas, failed, _, err := ApplyMessage(testVm, &msg, &gasPool)
		assert.NoError(b, err)
		assert.False(b, failed)
		assert.NotNil(b, usedGas)
	}
}

func BenchmarkApplyMessage_Call(b *testing.B) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)

	// create tx
	tx1 := createContractTx(WASMPath, AbiPath, 0, testGasLimit)
	msg1, err := tx1.AsMessage(true)
	assert.NoError(b, err)

	// call tx
	name := []byte("ApplyMsg")
	params := [][]byte{name}
	to := cs_crypto.CreateContractAddress(aliceAddr, 0)
	tx2 := callContractTx(&to, "returnString", params, 1)
	msg2, err := tx2.AsMessage(true)
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		db, root := CreateTestStateDB()
		testVm := getTestVm(db, root)
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
