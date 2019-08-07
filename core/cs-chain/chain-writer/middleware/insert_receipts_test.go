package middleware

import (
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidGasUsedAndReceipts(t *testing.T) {
	gasLimit := chain_config.BlockGasLimit

	receipt1 := &model2.Receipt{GasUsed: model2.TxGas, CumulativeGasUsed: model2.TxGas * 1}
	receipt2 := &model2.Receipt{GasUsed: model2.TxGas, CumulativeGasUsed: model2.TxGas * 2}
	receipt3 := &model2.Receipt{GasUsed: model2.TxGas, CumulativeGasUsed: model2.TxGas * 3}
	receipts := model2.Receipts{receipt1, receipt2, receipt3}
	txs := []model.AbstractTransaction{
		&fakeTx{GasLimit: g_testData.TestGasLimit, Receipt: receipt1},
		&fakeTx{GasLimit: g_testData.TestGasLimit, Receipt: receipt2},
		&fakeTx{GasLimit: g_testData.TestGasLimit, Receipt: receipt3},
	}

	testBlock := &fakeBlock{
		GasLimit:    chain_config.BlockGasLimit,
		ReceiptHash: model.DeriveSha(receipts),
		GasUsed:     receipt3.CumulativeGasUsed,
		txs:         txs,
	}

	assert.NoError(t, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())

}

func TestValidGasUsedAndReceipts_Error(t *testing.T) {
	// is special block
	testBlock := &fakeBlock{isSpecial: true}
	assert.NoError(t, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
	})())

	// empty receipt
	gasLimit := chain_config.BlockGasLimit
	tx := &fakeTx{}
	txs := []model.AbstractTransaction{tx}
	testBlock = &fakeBlock{txs: txs}
	assert.Equal(t, g_error.ErrEmptyReceipt, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())

	// receipt hash not match
	receipt := &model2.Receipt{GasUsed: model2.TxGas, CumulativeGasUsed: model2.TxGas * 1}
	tx.Receipt = receipt
	assert.Equal(t, g_error.ReceiptHashError, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())

	// gasUsed is invalid
	testBlock.ReceiptHash = model.DeriveSha(model2.Receipts{receipt})
	assert.Equal(t, g_error.ErrGasUsedIsInvalid, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())

	// gasUsed is over-ranging
	testBlock.GasUsed = model2.TxGas * 1
	assert.Equal(t, g_error.ErrTxGasIsOverRanging, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
}
