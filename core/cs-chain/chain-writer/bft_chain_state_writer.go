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
	"github.com/dipperin/dipperin-core/third-party/log"
)

type BftChainWriter struct {
	context *middleware.BftBlockContext
	chain   middleware.ChainInterface
}

func NewBftChainWriter(context *middleware.BftBlockContext, chain middleware.ChainInterface) *BftChainWriter {
	return &BftChainWriter{context: context, chain: chain}
}

func (cw *BftChainWriter) SaveBlock() error {
	log.Info("bftChainWriter save block")

	//Create a BlockProcessor
	c := cw.context

	c.Use(middleware.ValidateBlockNumber(&c.BlockContext))
	if !model.IsIgnoreDifficultyValidation() {
		c.Use(middleware.ValidateBlockDifficulty(&c.BlockContext))
	}
	c.Use(middleware.ValidateBlockVersion(&c.BlockContext))
	//c.Use(middleware.ValidateBlockSize(&c.BlockContext))
	c.Use(middleware.ValidateBlockHash(&c.BlockContext))
	c.Use(middleware.ValidateBlockCoinBase(&c.BlockContext))
	c.Use(middleware.ValidateSeed(&c.BlockContext))
	c.Use(middleware.ValidateBlockTime(&c.BlockContext))
	c.Use(middleware.ValidateGasLimit(&c.BlockContext))
	c.Use(middleware.ValidateBlockTxs(&c.BlockContext))
	c.Use(middleware.ValidateVotes(c))
	c.Use(middleware.UpdateStateRoot(&c.BlockContext))
	c.Use(middleware.UpdateBlockVerifier(&c.BlockContext))
	c.Use(middleware.ValidGasUsedAndReceipts(&c.BlockContext))
	c.Use(middleware.InsertBlock(&c.BlockContext))
	c.Use(middleware.NextRoundVerifier(&c.BlockContext))
	//Call BlockProcessor.Process

	err := c.Process()
	if err != nil {
		log.Error("bft save block failed!", "err", err)
	}

	return err
}

type BftChainWriterWithoutVotes struct {
	context *middleware.BftBlockContextWithoutVotes
	chain   middleware.ChainInterface
}

func NewBftChainWriterWithoutVotes(context *middleware.BftBlockContextWithoutVotes, chain middleware.ChainInterface) *BftChainWriterWithoutVotes {
	return &BftChainWriterWithoutVotes{context: context, chain: chain}
}

func (cw *BftChainWriterWithoutVotes) SaveBlock() error {
	log.Info("bftChainWriter save block")
	//Create a BlockProcessor
	c := cw.context

	c.Use(middleware.ValidateBlockNumber(&c.BlockContext))
	if !model.IsIgnoreDifficultyValidation() {
		c.Use(middleware.ValidateBlockDifficulty(&c.BlockContext))
	}
	c.Use(middleware.ValidateBlockVersion(&c.BlockContext))
	//c.Use(middleware.ValidateBlockSize(&c.BlockContext))
	c.Use(middleware.ValidateBlockHash(&c.BlockContext))
	c.Use(middleware.ValidateBlockCoinBase(&c.BlockContext))
	c.Use(middleware.ValidateSeed(&c.BlockContext))

	c.Use(middleware.ValidateBlockTxs(&c.BlockContext))

	c.Use(middleware.UpdateStateRoot(&c.BlockContext))

	c.Use(middleware.UpdateBlockVerifier(&c.BlockContext))
	c.Use(middleware.InsertBlock(&c.BlockContext))

	// after insert block, update verifier
	c.Use(middleware.NextRoundVerifier(&c.BlockContext))
	//Call BlockProcessor.Process

	err := c.Process()
	if err != nil {
		log.Error("bft save block failed", "err", err)
	}

	return err
}
