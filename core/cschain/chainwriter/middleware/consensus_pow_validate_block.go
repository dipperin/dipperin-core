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
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"go.uber.org/zap"
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
		log.DLogger.Info("the save block info is:", zap.Any("block", c.Block))
		if c.Chain == nil || c.Block == nil {
			return gerror.ErrChainOrBlockIsNil
		}

		rollBackNum := c.Chain.GetChainConfig().RollBackNum
		curBlock := c.Chain.CurrentBlock()
		if curBlock.Number() >= c.Block.Number()+rollBackNum {
			return gerror.ErrBlockHeightTooLow
		}

		if curBlock.Number() >= c.Block.Number() && !c.Block.IsSpecial() {
			return gerror.ErrNormalBlockHeightTooLow
		}

		if curBlock.Number()+1 < c.Block.Number() {
			max := big.NewInt(time.Now().Add(time.Second * maxTimeFutureBlocks).UnixNano())
			cmpResult := c.Block.Timestamp().Cmp(max)
			//log.DLogger.Info("check future block valid", "cmp result", cmpResult, "block timestamp", block.Timestamp().String(), "max to", max.String())
			if cmpResult > 0 {
				return gerror.ErrFutureBlockTooFarAway
			}
			return gerror.ErrFutureBlock
		}
		log.DLogger.Info("ValidateBlockNumber success")
		return c.Next()
	}
}

func ValidateBlockHash(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateBlockHash start")
		preBlock := c.Chain.GetBlockByNumber(c.Block.Number() - 1)
		preRv := reflect.ValueOf(preBlock)
		if !preRv.IsValid() || preRv.IsNil() {
			return gerror.ErrPreBlockIsNil
		}

		if !c.Block.PreHash().IsEqual(preBlock.Hash()) {
			log.DLogger.Error("pre block hash not match", zap.String("block pre hash", c.Block.PreHash().Hex()),
				zap.Any("pre block hash", preBlock.Hash().Hex()))
			return gerror.ErrPreBlockHashNotMatch
		}
		log.DLogger.Info("ValidateBlockHash end")
		return c.Next()
	}
}

func ValidateBlockDifficulty(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateBlockDifficulty start")
		if c.Block.IsSpecial() {
			log.DLogger.Info("ValidateBlockDifficulty the block is special")
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
			log.DLogger.Error("the c.Block number is:", zap.Uint64("number", c.Block.Number()))
			log.DLogger.Error("valid difficulty error", zap.String("targetDiff", targetDiff.Hex()), zap.String("blockDifficulty", c.Block.Difficulty().Hex()))
			return gerror.ErrInvalidDiff
		}

		// valid block hash for difficulty
		log.DLogger.Info("ValidateBlockDifficulty", zap.String("calculate difficulty", c.Block.RefreshHashCache().Hex()), zap.String("block difficulty", c.Block.Difficulty().DiffToTarget().Hex()))
		if !c.Block.RefreshHashCache().ValidHashForDifficulty(c.Block.Difficulty()) {
			log.DLogger.Error("ValidateBlockDifficulty failed")
			fmt.Println(c.Block.Header().(*model.Header).String())
			return gerror.ErrInvalidHashDiff
		}
		log.DLogger.Info("ValidateBlockDifficulty success")
		return c.Next()
	}
}

func ValidateBlockCoinBase(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateBlockCoinBase start")
		if c.Block.IsSpecial() {
			if !model.CheckAddressIsVerifierBootNode(c.Block.CoinBaseAddress()) {
				return gerror.ErrInvalidCoinBase
			}
		}
		log.DLogger.Info("ValidateBlockCoinBase success")
		return c.Next()
	}
}

func ValidateSeed(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateSeed start")
		block := c.Block
		preBlockHeight := block.Number() - 1
		log.DLogger.Info("ValidateSeed the preBlockHeight is:", zap.Uint64("preBlockHeight", preBlockHeight))
		preBlock := c.Chain.GetBlockByNumber(preBlockHeight)
		log.DLogger.Info("ValidateSeed the preBlock is:", zap.Any("preBlock", preBlock))

		seed := preBlock.Header().GetSeed().Bytes()
		proof := block.Header().GetProof()
		pk := block.Header().GetMinerPubKey()

		if pk == nil {
			return gerror.ErrPkIsNil
		}

		result, err := crypto.VRFVerify(pk, seed, proof)
		if err != nil {
			return err
		}
		if !result {
			return gerror.ErrSeedNotMatch
		}
		address := cs_crypto.GetNormalAddress(*pk)
		if !address.IsEqual(block.CoinBaseAddress()) {
			return gerror.ErrCoinBaseNotMatch
		}
		log.DLogger.Info("ValidateSeed success")
		return c.Next()
	}
}

func ValidateBlockVersion(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateBlockVersion start")
		if c.Block.Version() != c.Chain.GetChainConfig().Version {
			return gerror.ErrInvalidBlockVersion
		}
		log.DLogger.Info("ValidateBlockVersion end")
		return c.Next()
	}
}

func ValidateBlockTime(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateBlockTime start")
		blockTime := c.Block.Timestamp().Int64()
		if time.Now().Add(c.Chain.GetChainConfig().BlockTimeRestriction).UnixNano() < blockTime {
			return gerror.ErrInvalidBlockTimeStamp
		}
		log.DLogger.Info("ValidateBlockTime success")
		return c.Next()
	}
}

// valid gas limit
func ValidateGasLimit(c *BlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateGasLimit start")
		if c.Block.IsSpecial() {
			return c.Next()
		}
		currentGasLimit := c.Block.Header().GetGasLimit()
		// Verify that the gas limit is <= 2^63-1
		if currentGasLimit > chainconfig.MaxGasLimit || currentGasLimit < model.MinGasLimit {
			log.DLogger.Error("Invalid GasLimit", zap.Uint64("curGasLimit", currentGasLimit), zap.Uint64("maxGasLimit", chainconfig.MaxGasLimit), zap.Uint64("minGasLimit", model.MinGasLimit))
			return gerror.ErrInvliadHeaderGasLimit
		}
		parentGasLimit := c.Chain.GetLatestNormalBlock().Header().GetGasLimit()
		diff := int64(currentGasLimit) - int64(parentGasLimit)
		if diff < 0 {
			diff *= -1
		}
		limit := parentGasLimit / model.GasLimitBoundDivisor

		if uint64(diff) >= limit {
			log.DLogger.Error("Invalid GasLimit with parent block", zap.Uint64("curGasLimit", currentGasLimit), zap.Uint64("parentGasLimit", parentGasLimit), zap.Uint64("limitDiff", limit))
			return gerror.ErrHeaderGasLimitNotEnough
		}
		log.DLogger.Info("ValidateGasLimit success")
		return c.Next()
	}
}
