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

package chain_writer

import (
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/core/model"
)

type PowChainWriter struct {
	context *middleware.BlockContext
	chain   middleware.ChainInterface
}

func NewPowChainWriter(context *middleware.BlockContext, chain middleware.ChainInterface) *PowChainWriter {
	return &PowChainWriter{context: context, chain: chain}
}

func (cw *PowChainWriter) SaveBlock() error {
	c := cw.context

	c.Use(middleware.ValidateBlockNumber(c))
	if !model.IsIgnoreDifficultyValidation() {
		c.Use(middleware.ValidateBlockDifficulty(c))
	}
	c.Use(middleware.ValidateBlockHash(c))
	c.Use(middleware.ValidateBlockTxs(c))
	c.Use(middleware.UpdateStateRoot(c))
	c.Use(middleware.InsertBlock(c))

	//Call BlockProcessor.Process
	return c.Process()
}
