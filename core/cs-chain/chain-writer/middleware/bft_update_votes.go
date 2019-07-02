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
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/third-party/log"
)

func UpdateBlockVerifier(c *BlockContext) Middleware {
	return func() error {
		registerPro, err := validBlockVerifier(c)
		if err != nil {
			return err
		}

		if _, err := registerPro.Commit(); err != nil {
			return err
		}
		log.Info("commit register root successful")
		return c.Next()
	}
}

func ValidBlockVerifier(c *BlockContext) Middleware {
	return func() error {
		if _, err := validBlockVerifier(c); err != nil {
			return err
		}
		return c.Next()
	}
}

// Two methods will be called concurrently, one is SaveBlock and the other is FullValid in PBFT
func validBlockVerifier(c *BlockContext) (*registerdb.RegisterDB, error) {
	// check register
	preBlockHeight := c.Block.Number() - 1
	preBlock := c.Chain.GetBlockByNumber(preBlockHeight)
	registerPro, gErr := c.Chain.BuildRegisterProcessor(preBlock.Header().GetRegisterRoot())
	if gErr != nil {
		return nil, gErr
	}

	if err := registerPro.Process(c.Block); err != nil {
		return nil, err
	}

	registerRoot := registerPro.Finalise()
	if !registerRoot.IsEqual(c.Block.GetRegisterRoot()) {
		return nil, g_error.ErrRegisterRootNotMatch
	}

	return registerPro, nil
}

func NextRoundVerifier(c *BlockContext) Middleware {
	return func() error {
		chain := c.Chain
		// check register
		// insert success then calculate verifiers
		slot := chain.GetSlot(c.Block)
		if chain.IsChangePoint(c.Block, false) {
			chain.GetVerifiers(*slot + chain.GetChainConfig().SlotMargin)
		}

		return c.Next()
	}
}
