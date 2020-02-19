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
	"github.com/dipperin/dipperin-core/third-party/log"
)

func InsertBlock(c *BlockContext) Middleware {
	return func() error {
		curBlock := c.Chain.CurrentBlock()
		log.Middleware.Info("InsertBlock start", "curNum", curBlock.Number(), "blockNum", c.Block.Number())

		// roll back chain if insert special block
		if c.Block.IsSpecial() {
			if err := c.Chain.Rollback(c.Block.Number()); err != nil {
				return err
			}
		}

		log.Info("insert block", "cur number", curBlock.Number(), "new number", c.Block.Number())
		// check block number
		if c.Chain.CurrentBlock().Number()+1 != c.Block.Number() {
			return g_error.ErrInvalidBlockNum
		}

		if err := c.Chain.GetChainDB().InsertBlock(c.Block); err != nil {
			return err
		}
		log.Info("insert block successful", "num", c.Block.Number())
		//currentBlock := c.Chain.CurrentBlock()
		//log.Info("the currentBlock number is~~~~~~~~~~~~~`:","number",currentBlock.Number())

		//insert receipts
		if !c.Block.IsSpecial() {
			if err := c.Chain.GetChainDB().SaveReceipts(c.Block.Hash(), c.Block.Number(), c.receipts); err != nil {
				return err
			}
			log.Info("insert receipts successful", "num", c.Block.Number())
		}

		log.Middleware.Info("InsertBlock success")
		return c.Next()
	}
}
