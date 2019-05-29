package state_processor

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
	"math/big"
)

func TestApplyMessage(t *testing.T) {
	var testPath = "/home/qydev/go/src/github.com/dipperin/dipperin-core/core/vm/event"
	tx := createContractTx(t, nil, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	msg, err := tx.AsMessage()
	assert.NoError(t, err)

	a := make(map[common.Address]*big.Int)
	testVm := getTestVm(a)

	result, usedGas, failed, err := ApplyMessage(testVm, msg, 5*gasLimit, fakeAccountStateTx{account:a})
	assert.NoError(t, err)
	assert.False(t, failed)
	assert.NotNil(t, usedGas)
	data := getTxData(t, testPath+"/event.wasm", testPath+"/event.cpp.abi.json")
	assert.Equal(t, data, result)
}