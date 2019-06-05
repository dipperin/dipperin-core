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
	"fmt"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
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
		if c.Chain == nil || c.Block == nil {
			fmt.Println(c.Chain == nil, c.Block == nil)
			return errors.New("chain or block cannot be null")
		}
		cur := c.Chain.CurrentBlock()
		if cur.Number() > c.Block.Number() {
			return g_error.ErrBlockHeightTooLow
		}
		if cur.Number() == c.Block.Number() && !c.Block.IsSpecial() {
			return g_error.ErrBlockHeightIsCurrentAndIsNotSpecial
		}

		// if insert same height special block, continue
		if c.Block.Number() == cur.Number() && c.Block.IsSpecial() {
			log.Info("insert same height special block", "num", c.Block.Number())
			return c.Next()
		}

		if cur.Number()+1 != c.Block.Number() {
			max := big.NewInt(time.Now().Add(time.Second * maxTimeFutureBlocks).UnixNano())
			cmpResult := c.Block.Timestamp().Cmp(max)
			//log.Info("check future block valid", "cmp result", cmpResult, "block timestamp", block.Timestamp().String(), "max to", max.String())
			if cmpResult > 0 {
				return g_error.ErrFutureBlockTooFarAway
			}
			return g_error.ErrFutureBlock
		}
		return c.Next()
	}
}

func ValidateBlockHash(c *BlockContext) Middleware {
	return func() error {
		preBlock := c.Chain.GetBlockByNumber(c.Block.Number() - 1)
		preRv := reflect.ValueOf(preBlock)
		if  !preRv.IsValid() || preRv.IsNil() {
			return g_error.ErrPreBlockIsNil
		}

		if !c.Block.PreHash().IsEqual(preBlock.Hash()) {
			//fmt.Println("pre block", preBlock, preBlock.Hash())
			//fmt.Println("c.Block", c.Block, c.Block.Hash())
			log.Error("pre block hash not match", "block pre hash", c.Block.PreHash().Hex(),
				"pre block hash", preBlock.Hash().Hex())
			return g_error.ErrPreBlockHashNotMatch
		}

		return c.Next()
	}
}

func ValidateBlockSize(c *BlockContext) Middleware {
	return func() error {
		bb, err := rlp.EncodeToBytes(c.Block)
		if err != nil {
			return err
		}
		if len(bb) > chain_config.MaxBlockSize {
			return g_error.ErrBlockSizeTooLarge
		}
		return c.Next()
	}
}

func ValidateBlockDifficulty(c *BlockContext) Middleware {
	return func() error {
		if c.Block.IsSpecial() {
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
			return g_error.ErrInvalidDiff
		}

		// valid block hash for difficulty
		if !c.Block.RefreshHashCache().ValidHashForDifficulty(c.Block.Difficulty()) {
			log.Info("ValidateBlockDifficulty failed", "block", c.Block.Header().(*model.Header).String())
			return g_error. ErrWrongHashDiff
		}
		return c.Next()
	}
}

func ValidateBlockCoinBase(c *BlockContext) Middleware {
	return func() error {
		if c.Block.IsSpecial() {
			if !model.CheckAddressIsVerifierBootNode(c.Block.CoinBaseAddress()) {
				return g_error.ErrSpecialInvalidCoinBase
			}
		}

		return c.Next()
	}
}

func ValidateSeed(c *BlockContext) Middleware {
	return func() error {
		block := c.Block
		preBlockHeight := block.Number() - 1
		preBlock := c.Chain.GetBlockByNumber(preBlockHeight)

		seed := preBlock.Header().GetSeed().Bytes()
		proof := block.Header().GetProof()
		pk := block.Header().GetMinerPubKey()

		if pk == nil {
			return g_error.ErrNotGetPk
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
			return g_error.ErrPkNotIsCoinBase
		}
		return c.Next()
	}
}

func ValidateBlockVersion(c *BlockContext) Middleware {
	return func() error {
		if c.Block.Version() != c.Chain.GetChainConfig().Version {
			return g_error.ErrBlockVer
		}
		return c.Next()
	}
}


func ValidateBlockTime(c *BlockContext) Middleware {
	return func() error {
		blockTime := c.Block.Timestamp().Int64()
		if time.Now().Add(c.Chain.GetChainConfig().BlockTimeRestriction).UnixNano() < blockTime {
			return g_error.ErrBlockTimeStamp
		}
		return c.Next()
	}
}
