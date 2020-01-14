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

package cschain

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/cschain/chainstate"
	"github.com/dipperin/dipperin-core/core/cschain/chainwriter/middleware"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hashicorp/golang-lru"
	"go.uber.org/zap"
	"sync/atomic"
)

const (
	bodyCacheLimit     = 256
	blockCacheLimit    = 256
	headerCacheLimit   = 512
	numberCacheLimit   = 2048
	verifierCacheLimit = 12
	slotCacheLimit     = 1024 * 5

	// test env have not enough mem
	bodyCacheLimitTestEnv   = 30
	blockCacheLimitTestEnv  = 30
	headerCacheLimitTestEnv = 60
	numberCacheLimitTestEnv = 512
	slotCacheLimitTestEnv   = 1024
)

func NewCacheChainState(cs *chainstate.ChainState) (*CacheChainState, error) {
	bodyCache, _ := lru.New(bodyCacheLimit)
	bodyRLPCache, _ := lru.New(bodyCacheLimit)
	blockCache, _ := lru.New(blockCacheLimit)
	headerCache, _ := lru.New(headerCacheLimit)
	numberCache, _ := lru.New(numberCacheLimit)
	cachedVerifiers, _ := lru.New(verifierCacheLimit)
	slotCache, _ := lru.New(slotCacheLimit)

	if chainconfig.GetCurBootsEnv() == "test" {
		bodyCache, _ = lru.New(bodyCacheLimitTestEnv)
		bodyRLPCache, _ = lru.New(bodyCacheLimitTestEnv)
		blockCache, _ = lru.New(blockCacheLimitTestEnv)
		headerCache, _ = lru.New(headerCacheLimitTestEnv)
		numberCache, _ = lru.New(numberCacheLimitTestEnv)
		slotCache, _ = lru.New(slotCacheLimitTestEnv)
	}

	ccs := &CacheChainState{
		ChainState:   cs,
		bodyCache:    bodyCache,
		bodyRLPCache: bodyRLPCache,
		blockCache:   blockCache,
		//FutureBlocks:    FutureBlocks,
		headerCache:     headerCache,
		numberCache:     numberCache,
		cachedVerifiers: cachedVerifiers,
		slotCache:       slotCache,
	}

	return ccs, nil
}

type CacheChainState struct {
	*chainstate.ChainState

	blockCache      *lru.Cache
	numberCache     *lru.Cache
	bodyCache       *lru.Cache
	bodyRLPCache    *lru.Cache
	headerCache     *lru.Cache
	cachedVerifiers *lru.Cache
	slotCache       *lru.Cache

	genesisBlock  model.AbstractBlock
	currentBlock  atomic.Value
	currentHeader atomic.Value
}

func (chain *CacheChainState) SaveBftBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	err := chain.WriterFactory.NewWriter(middleware.NewBftBlockContext(block, seenCommits, chain)).SaveBlock()
	if err != nil {
		log.DLogger.Error("SaveBftBlock err", zap.Error(err))
		return err
	}

	chain.currentBlock.Store(block)
	chain.currentHeader.Store(block.Header())

	return nil
}

// Deprecated: use SaveBftBlock
//func (chain *CacheChainState) SaveBlock(block model.AbstractBlock) error {
//	panic("cache chain state no allow use chain state save block")
//}

func (chain *CacheChainState) Genesis() model.AbstractBlock {
	if chain.genesisBlock != nil {
		return chain.genesisBlock
	}

	if genesisBlock := chain.ChainState.Genesis(); genesisBlock != nil {
		chain.genesisBlock = genesisBlock
		return genesisBlock
	}

	log.DLogger.Error("chain state can't get genesis block")

	return nil
}

func (chain *CacheChainState) CurrentBlock() model.AbstractBlock {
	curB := chain.currentBlock.Load()
	if curB == nil {
		curB = chain.ChainState.CurrentBlock()
		if curB == nil {
			return nil
		}
		chain.currentBlock.Store(curB)
		chain.currentHeader.Store(curB.(model.AbstractBlock).Header())
	}

	return curB.(model.AbstractBlock)
}

func (chain *CacheChainState) GetBody(hash common.Hash) model.AbstractBody {
	if cached, ok := chain.bodyCache.Get(hash); ok {
		return cached.(model.AbstractBody)
	}

	body := chain.ChainState.GetBody(hash)
	if body == nil {
		return nil
	}
	chain.bodyCache.Add(hash, body)

	return body
}

func (chain *CacheChainState) GetBodyRLP(hash common.Hash) rlp.RawValue {
	if cached, ok := chain.bodyRLPCache.Get(hash); ok {
		return cached.(rlp.RawValue)
	}

	bodyRLP := chain.ChainState.GetBodyRLP(hash)
	if bodyRLP == nil || len(bodyRLP) == 0 {
		return nil
	}
	chain.bodyRLPCache.Add(hash, bodyRLP)

	return bodyRLP
}

func (chain *CacheChainState) GetBlock(hash common.Hash, number uint64) model.AbstractBlock {
	if block, ok := chain.blockCache.Get(hash); ok {
		return block.(model.AbstractBlock)
	}

	block := chain.ChainState.GetBlock(hash, number)
	if block == nil {
		return nil
	}
	chain.blockCache.Add(block.Hash(), block)

	return block
}

func (chain *CacheChainState) HasBlock(hash common.Hash, number uint64) bool {
	if chain.blockCache.Contains(hash) {
		return true
	}

	return chain.ChainState.HasBlock(hash, number)
}

func (chain *CacheChainState) GetHeader(hash common.Hash, number uint64) model.AbstractHeader {
	if header, ok := chain.headerCache.Get(hash); ok {
		return header.(model.AbstractHeader)
	}

	header := chain.ChainState.GetHeader(hash, number)
	if header == nil {
		return nil
	}

	chain.headerCache.Add(hash, header)
	return header
}

func (chain *CacheChainState) HasHeader(hash common.Hash, number uint64) bool {
	if chain.numberCache.Contains(hash) || chain.headerCache.Contains(hash) {
		return true
	}

	return chain.ChainState.HasHeader(hash, number)
}

func (chain *CacheChainState) GetBlockNumber(hash common.Hash) *uint64 {
	if cached, ok := chain.numberCache.Get(hash); ok {
		number := cached.(uint64)
		return &number
	}

	number := chain.ChainState.GetBlockNumber(hash)
	if number != nil {
		chain.numberCache.Add(hash, *number)
	}

	return number
}

func (chain *CacheChainState) GetVerifiers(slot uint64) []common.Address {
	if vs, ok := chain.cachedVerifiers.Get(slot); ok {
		return vs.([]common.Address)
	}

	vs := chain.ChainState.GetVerifiers(slot)
	if len(vs) > 0 {
		chain.cachedVerifiers.Add(slot, vs)
	}

	return vs
}

func (chain *CacheChainState) CalVerifiers(block model.AbstractBlock) {
	vs := chain.ChainState.CalVerifiers(block)
	if len(vs) > 0 {
		slot := chain.GetSlotByNum(block.Number())
		chain.cachedVerifiers.Add(*slot+chain.ChainConfig.SlotMargin, vs)
	}
}

func (chain *CacheChainState) GetSlotByNum(num uint64) *uint64 {
	if s, ok := chain.slotCache.Get(num); ok {
		return s.(*uint64)
	}
	if x := chain.ChainState.GetSlotByNum(num); x != nil {
		chain.slotCache.Add(num, x)
		return x
	}
	return nil
}

func (chain *CacheChainState) GetSlot(block model.AbstractBlock) *uint64 {
	num := block.Number()
	if s, ok := chain.slotCache.Get(num); ok {
		return s.(*uint64)
	}
	if x := chain.ChainState.GetSlot(block); x != nil {
		chain.slotCache.Add(num, x)
		return x
	}
	return nil
}

func (chain *CacheChainState) GetBlockByHash(hash common.Hash) model.AbstractBlock {
	number := chain.GetBlockNumber(hash)
	if number == nil {
		return nil
	}

	return chain.GetBlock(hash, *number)
}

func (chain *CacheChainState) Rollback(target uint64) error {
	curBlock := chain.CurrentBlock()
	if target == curBlock.Number()+1 {
		return nil
	}

	rollBackNum := chain.GetChainConfig().RollBackNum
	if target > curBlock.Number()+1 || target+rollBackNum <= curBlock.Number() {
		return gerror.ErrTargetOutOfRange
	}

	// roll back current block
	tarBlock := chain.GetBlockByNumber(target - 1)
	if tarBlock == nil {
		return gerror.ErrPreTargetBlockIsNil
	} else {
		chain.currentBlock.Store(tarBlock)
		chain.currentHeader.Store(tarBlock.Header())
		return nil
	}
}
