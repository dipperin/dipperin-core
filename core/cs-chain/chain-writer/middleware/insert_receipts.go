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
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"go.uber.org/zap"
)

func ValidGasUsedAndReceipts(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidGasUsedAndReceipts start")
		if c.Block.IsSpecial() {
			log.DLogger.Info("ValidGasUsedAndReceipts the block is empty block")
			return c.Next()
		}
		curBlock := c.Chain.CurrentBlock()
		log.DLogger.Info("Insert block receipts", zap.Uint64("cur number", curBlock.Number()), zap.Uint64("new number", c.Block.Number()))
		receipts := make(model2.Receipts, 0, c.Block.TxCount())
		var accumulatedGas uint64
		if err := c.Block.TxIterator(func(i int, transaction model.AbstractTransaction) error {
			receipt := transaction.GetReceipt()
			if receipt == nil {
				return g_error.ErrTxReceiptIsNil
			}
			accumulatedGas = receipt.CumulativeGasUsed
			receipts = append(receipts, receipt)
			return nil
		}); err != nil {
			return err
		}

		//check receipt hash
		receiptHash := model.DeriveSha(receipts)
		if receiptHash != c.Block.GetReceiptHash() {
			log.DLogger.Error("InsertReceipts receiptHash not match", zap.Any("receiptHash", receiptHash), zap.Any("block.ReciptHash", c.Block.GetReceiptHash()))
			return g_error.ErrReceiptHashNotMatch
		}

		if accumulatedGas != c.Block.Header().GetGasUsed() {
			log.DLogger.Error("InsertReceipts accumulatedGas not match", zap.Uint64("accumulatedGas", accumulatedGas), zap.Uint64("headerGasUsed", c.Block.Header().GetGasUsed()))
			return g_error.ErrInvalidHeaderGasUsed
		}

		//check accumulated Gas
		if accumulatedGas > c.Block.Header().GetGasLimit() {
			return g_error.ErrHeaderGasUsedOverRanging
		}

		//padding receipts
		c.receipts = receipts
		log.DLogger.Info("ValidGasUsedAndReceipts success")
		return c.Next()
	}
}
