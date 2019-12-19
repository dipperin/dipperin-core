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
	"fmt"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"reflect"
	"time"
)

/*
	validate_block typical usage:
	bc := InitBlockContext(block)
	m := &BlockProcessor{bc, chain}
	m.Use(ValidateBlockNumber(bc), ValidateBlockDifficulty(bc), ValidateBlockTxs(bc), )
	m.Use(VrfCheckCommits(bc), )
	m.Use(UpdateState(bc))
	m.Use(InsertBlock(bc))

*/
func ValidateBlockNumber(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("the save block info is:", "block", c.Block)
		if c.Chain == nil || c.Block == nil {
			return g_error.ErrChainOrBlockIsNil
		}

		rollBackNum := c.Chain.GetChainConfig().RollBackNum
		curBlock := c.Chain.CurrentBlock()
		if curBlock.Number() >= c.Block.Number()+rollBackNum {
			return g_error.ErrBlockHeightTooLow
		}

		if curBlock.Number() >= c.Block.Number() && !c.Block.IsSpecial() {
			return g_error.ErrNormalBlockHeightTooLow
		}

		if curBlock.Number()+1 < c.Block.Number() {
			max := big.NewInt(time.Now().Add(time.Second * maxTimeFutureBlocks).UnixNano())
			cmpResult := c.Block.Timestamp().Cmp(max)
			//log.Info("check future block valid", "cmp result", cmpResult, "block timestamp", block.Timestamp().String(), "max to", max.String())
			if cmpResult > 0 {
				return g_error.ErrFutureBlockTooFarAway
			}
			return g_error.ErrFutureBlock
		}
		log.Middleware.Info("ValidateBlockNumber success")
		return c.Next()
	}
}

func ValidateBlockHash(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("ValidateBlockHash start")
		preBlock := c.Chain.GetBlockByNumber(c.Block.Number() - 1)
		preRv := reflect.ValueOf(preBlock)
		if !preRv.IsValid() || preRv.IsNil() {
			return g_error.ErrPreBlockIsNil
		}

		if !c.Block.PreHash().IsEqual(preBlock.Hash()) {
			log.Error("pre block hash not match", "block pre hash", c.Block.PreHash().Hex(),
				"pre block hash", preBlock.Hash().Hex())
			return g_error.ErrPreBlockHashNotMatch
		}
		log.Middleware.Info("ValidateBlockHash end")
		return c.Next()
	}
}

func ValidateBlockDifficulty(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("ValidateBlockDifficulty start")
		if c.Block.IsSpecial() {
			log.Middleware.Info("ValidateBlockDifficulty the block is special")
			return c.Next()

		}

		preBlockHeight := c.Block.Number() - 1
		preSpanH := model.LastPeriodBlockNum(preBlockHeight)
		if preSpanH == 0 {
			preSpanH = 1
		}

		//get the first block in the preBlock period
		preSpanBlock := c.Chain.GetBlockByNumber(preSpanH)
		//lastBlock := c.Chain.GetBlockByNumber(preBlockHeight)

		//find the neighbor normal block
		findBlock := c.Chain.GetLatestNormalBlock()

		targetDiff := model.NewCalNewWorkDiff(preSpanBlock, findBlock, c.Block.Number()-1)
		//targetDiff := model.CalNewWorkDiff(preSpanBlock, lastBlock)

		if !targetDiff.Equal(c.Block.Difficulty()) {
			log.Error("the c.Block number is:", "number", c.Block.Number())
			log.Error("valid difficulty error", "targetDiff", targetDiff.Hex(), "blockDifficulty", c.Block.Difficulty().Hex())
			return g_error.ErrInvalidDiff
		}

		// valid block hash for difficulty
		log.Info("ValidateBlockDifficulty", "calculate difficulty", c.Block.RefreshHashCache().Hex(), "block difficulty", c.Block.Difficulty().DiffToTarget().Hex())
		if !c.Block.RefreshHashCache().ValidHashForDifficulty(c.Block.Difficulty()) {
			log.Error("ValidateBlockDifficulty failed")
			fmt.Println(c.Block.Header().(*model.Header).String())
			return g_error.ErrInvalidHashDiff
		}
		log.Middleware.Info("ValidateBlockDifficulty success")
		return c.Next()
	}
}

func ValidateBlockCoinBase(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("ValidateBlockCoinBase start")
		if c.Block.IsSpecial() {
			if !model.CheckAddressIsVerifierBootNode(c.Block.CoinBaseAddress()) {
				return g_error.ErrInvalidCoinBase
			}
		}
		log.Middleware.Info("ValidateBlockCoinBase success")
		return c.Next()
	}
}

func ValidateSeed(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("ValidateSeed start")
		block := c.Block
		preBlockHeight := block.Number() - 1
		log.Middleware.Info("ValidateSeed the preBlockHeight is:", "preBlockHeight", preBlockHeight)
		preBlock := c.Chain.GetBlockByNumber(preBlockHeight)
		log.Middleware.Info("ValidateSeed the preBlock is:", "preBlock", preBlock)

		seed := preBlock.Header().GetSeed().Bytes()
		proof := block.Header().GetProof()
		pk := block.Header().GetMinerPubKey()

		if pk == nil {
			return g_error.ErrPkIsNil
		}

		result, err := crypto.VRFVerify(pk, seed, proof)
		if err != nil {
			return err
		}
		if !result {
			return g_error.ErrSeedNotMatch
		}
		address := cs_crypto.GetNormalAddress(*pk)
		if !address.IsEqual(block.CoinBaseAddress()) {
			return g_error.ErrCoinBaseNotMatch
		}
		log.Middleware.Info("ValidateSeed success")
		return c.Next()
	}
}

func ValidateBlockVersion(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("ValidateBlockVersion start")
		if c.Block.Version() != c.Chain.GetChainConfig().Version {
			return g_error.ErrInvalidBlockVersion
		}
		log.Middleware.Info("ValidateBlockVersion end")
		return c.Next()
	}
}

func ValidateBlockTime(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("ValidateBlockTime start")
		blockTime := c.Block.Timestamp().Int64()
		if time.Now().Add(c.Chain.GetChainConfig().BlockTimeRestriction).UnixNano() < blockTime {
			return g_error.ErrInvalidBlockTimeStamp
		}
		log.Middleware.Info("ValidateBlockTime success")
		return c.Next()
	}
}

// valid gas limit
func ValidateGasLimit(c *BlockContext) Middleware {
	return func() error {
		log.Middleware.Info("ValidateGasLimit start")
		if c.Block.IsSpecial() {
			return c.Next()
		}
		currentGasLimit := c.Block.Header().GetGasLimit()
		// Verify that the gas limit is <= 2^63-1
		if currentGasLimit > chain_config.MaxGasLimit || currentGasLimit < model2.MinGasLimit {
			log.Error("Invalid GasLimit", "curGasLimit", currentGasLimit, "maxGasLimit", chain_config.MaxGasLimit, "minGasLimit", model2.MinGasLimit)
			return g_error.ErrInvliadHeaderGasLimit
		}
		parentGasLimit := c.Chain.GetLatestNormalBlock().Header().GetGasLimit()
		diff := int64(currentGasLimit) - int64(parentGasLimit)
		if diff < 0 {
			diff *= -1
		}
		limit := parentGasLimit / model2.GasLimitBoundDivisor

		if uint64(diff) >= limit {
			log.Error("Invalid GasLimit with parent block", "curGasLimit", currentGasLimit, "parentGasLimit", parentGasLimit, "limitDiff", limit)
			return g_error.ErrHeaderGasLimitNotEnough
		}
		log.Middleware.Info("ValidateGasLimit success")
		return c.Next()
	}
}
