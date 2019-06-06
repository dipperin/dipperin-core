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
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
)

func ValidGasUsedAndReceipts(c *BlockContext)Middleware{
	return func() error {
		if c.Block.IsSpecial() {
			return c.Next()
		}
		curBlock := c.Chain.CurrentBlock()
		log.Info("insert block receipts", "cur number", curBlock.Number(), "new number", c.Block.Number())
		receipts := make(model2.Receipts, 0, c.Block.TxCount())
		var accumulatedGas uint64
		if err := c.Block.TxIterator(func(i int, transaction model.AbstractTransaction) error {
			receipt, err := transaction.GetReceipt()
			if err != nil {
				return err
			}
			accumulatedGas += receipt.GasUsed
			receipts = append(receipts, receipt)
			return nil
		}); err != nil {
			return err
		}

		//check receipt hash
		receiptHash := model.DeriveSha(receipts)
		if receiptHash != c.Block.GetReceiptHash() {
			log.Error("InsertReceipts#receiptHash", "receiptHash", receiptHash, "block.ReciptHash", c.Block.GetReceiptHash())
			return g_error.ReceiptHashError
		}

		if accumulatedGas != c.Block.Header().GetGasUsed(){
			log.Info("the block gas info is:","accumulatedGas",accumulatedGas,"headerGasUsed",c.Block.Header().GetGasUsed())
			return g_error.ErrGasUsedIsInvalid
		}

		//check accumulated Gas
		if accumulatedGas > c.Block.Header().GetGasLimit(){
			return g_error.ErrTxGasIsOverRanging
		}

		//padding receipts
		c.receipts =receipts
		return c.Next()
	}
}
