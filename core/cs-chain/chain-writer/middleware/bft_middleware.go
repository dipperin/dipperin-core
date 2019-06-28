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
	"github.com/dipperin/dipperin-core/core/model"
)

type BftBlockContext struct {
	BlockContext

	Votes []model.AbstractVerification
}

func NewBftBlockContext(b model.AbstractBlock, votes []model.AbstractVerification, chain ChainInterface) *BftBlockContext {
	bc := &BftBlockContext{}
	bc.index = -1

	bc.Block = b
	bc.Votes = votes
	bc.Chain = chain
	return bc
}

type BftBlockContextWithoutVotes struct {
	BlockContext
}

func NewBftBlockContextWithoutVotes(b model.AbstractBlock, chain ChainInterface) *BftBlockContextWithoutVotes {
	bc := &BftBlockContextWithoutVotes{}
	bc.index = -1

	bc.Block = b
	bc.Chain = chain
	return bc
}

func NewBftBlockValidator(chain ChainInterface) *BftBlockValidator {
	return &BftBlockValidator{ Chain: chain }
}

type BftBlockValidator struct {
	Chain ChainInterface
}

func (v *BftBlockValidator) Valid(b model.AbstractBlock) error {
	c := NewBlockContext(b, v.Chain)

	c.Use(ValidateBlockNumber(c))
	if !model.IsIgnoreDifficultyValidation() {
		c.Use(ValidateBlockDifficulty(c))
	}
	c.Use(ValidateBlockVersion(c))
	c.Use(ValidateBlockSize(c))
	c.Use(ValidateBlockHash(c))
	c.Use(ValidateBlockCoinBase(c))
	c.Use(ValidateSeed(c))
	c.Use(ValidateBlockTime(c))
	c.Use(ValidateGasLimit(c))
	c.Use(ValidateBlockTxs(c))
	c.Use(ValidateVotesForBFT(c))

	c.Use(ValidStateRoot(c))
	c.Use(ValidGasUsedAndReceipts(c))
	c.Use(ValidBlockVerifier(c))
	return c.Process()
}

func (v *BftBlockValidator) FullValid(b model.AbstractBlock) error {
	return v.Valid(b)
}