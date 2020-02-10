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

package middleware

import (
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidGasUsedAndReceipts(t *testing.T) {
	gasLimit := chainconfig.BlockGasLimit
	
	receipt1 := &model.Receipt{GasUsed: model.TxGas, CumulativeGasUsed: model.TxGas * 1}
	receipt2 := &model.Receipt{GasUsed: model.TxGas, CumulativeGasUsed: model.TxGas * 2}
	receipt3 := &model.Receipt{GasUsed: model.TxGas, CumulativeGasUsed: model.TxGas * 3}
	receipts := model.Receipts{receipt1, receipt2, receipt3}
	txs := []model.AbstractTransaction{
		&fakeTx{GasLimit: 2 * model.TxGas, Receipt: receipt1},
		&fakeTx{GasLimit: 2 * model.TxGas, Receipt: receipt2},
		&fakeTx{GasLimit: 2 * model.TxGas, Receipt: receipt3},
	}
	
	testBlock := &fakeBlock{
		GasLimit:    chainconfig.BlockGasLimit,
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
	gasLimit := chainconfig.BlockGasLimit
	tx := &fakeTx{}
	txs := []model.AbstractTransaction{tx}
	testBlock = &fakeBlock{txs: txs}
	assert.Equal(t, gerror.ErrTxReceiptIsNil, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
	
	// receipt hash not match
	receipt := &model.Receipt{GasUsed: model.TxGas, CumulativeGasUsed: model.TxGas * 1}
	tx.Receipt = receipt
	assert.Equal(t, gerror.ErrReceiptHashNotMatch, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
	
	// gasUsed is invalid
	testBlock.ReceiptHash = model.DeriveSha(model.Receipts{receipt})
	assert.Equal(t, gerror.ErrInvalidHeaderGasUsed, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
	
	// gasUsed is over-ranging
	testBlock.GasUsed = model.TxGas * 1
	assert.Equal(t, gerror.ErrHeaderGasUsedOverRanging, ValidGasUsedAndReceipts(&BlockContext{
		Block: testBlock,
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
}
