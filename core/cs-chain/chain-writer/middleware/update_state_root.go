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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/mpt_log"
)

func UpdateStateRoot(c *BlockContext) Middleware {
	return func() error {
		processor, err := validStateRoot(c)
		if err != nil {
			return err
		}

		if _, err := processor.Commit(); err != nil {
			return err
		}
		log.Info("commit state root successful")
		return c.Next()
	}
}

func ValidStateRoot(c *BlockContext) Middleware {
	return func() error {
		if _, err := validStateRoot(c); err != nil {
			return err
		}
		return c.Next()
	}
}

func validStateRoot(c *BlockContext) (*chain.BlockProcessor, error) {
	// check state
	preBlockHeight := c.Block.Number() - 1
	preBlock := c.Chain.GetBlockByNumber(preBlockHeight)

	mpt_log.Log.Debug("the preBlock stateRoot is:", "preBlockStateRoot", preBlock.StateRoot().Hex())
	processor, gErr := c.Chain.BlockProcessor(preBlock.StateRoot())
	if gErr != nil {
		return nil, gErr
	}

	mpt_log.Log.Debug("process the block", "blockId", c.Block.Hash().Hex())
	if err := processor.Process(c.Block, c.Chain.GetEconomyModel()); err != nil {
		return nil, err
	}

	roots, err := processor.Finalise()
	if err != nil {
		return nil, err
	}

	if !roots.IsEqual(c.Block.StateRoot()) {
		mpt_log.Log.Debug("state root not match", "got", roots.Hex(), "in block", c.Block.StateRoot().Hex())
		log.Error("state root not match", "got", roots.Hex(), "in block", c.Block.StateRoot().Hex())
		//fmt.Println("state root check not match", c.Block)
		return nil, errors.New("state root not match")
	}

	//check reciptHash

	return processor, nil
}

type BlockProcessor func(root common.Hash) (*chain.BlockProcessor, error)

func ValidSateRootForTest(preStateRoot common.Hash, economyModel economy_model.EconomyModel, blockProcess BlockProcessor, processBlock model.AbstractBlock) error {
	mpt_log.Log.Debug("the preBlock stateRoot is:", "preBlockStateRoot", preStateRoot.Hex())
	processor, gErr := blockProcess(preStateRoot)
	if gErr != nil {
		return gErr
	}

	mpt_log.Log.Debug("process the block", "blockId", processBlock.Hash().Hex())
	if err := processor.Process(processBlock, economyModel); err != nil {
		return err
	}

	roots, err := processor.Finalise()
	if err != nil {
		return err
	}

	if !roots.IsEqual(processBlock.StateRoot()) {
		mpt_log.Log.Debug("state root not match", "got", roots.Hex(), "in block", processBlock.StateRoot().Hex())
		log.Debug("state root not match", "got", roots.Hex(), "in block", processBlock.StateRoot().Hex())
		return errors.New("state root not match")
	}

	return nil
}
