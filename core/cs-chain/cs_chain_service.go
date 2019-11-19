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

package cs_chain

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/g-event"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/common/g-timer"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/hashicorp/golang-lru"
	"math/big"
	"sort"
	"sync"
	"time"
)

const (
	maxFutureBlocks        = 256
	maxFutureBlocksTestEnv = 30
)

var (
	numLowBlockToReturnErr uint64 = 10
	GenesisSetUp           bool
)

func IsSetUpGenesis() bool {
	// setup genesis if not uint test or don't want to ignore
	if !util.IsTestEnv() || GenesisSetUp {
		return true
	}
	return false
}

type CsChainServiceConfig struct {
	CacheDB CacheDB
	TxPool  TxPool
}

type CsChainService struct {
	*CsChainServiceConfig
	*CacheChainState

	saveBlockLock   sync.RWMutex
	getVerifierLock sync.RWMutex

	wg sync.WaitGroup

	FutureBlocks *lru.Cache

	initOnce sync.Once

	Quit chan struct{}
}

func NewCsChainService(config *CsChainServiceConfig, cs *chain_state.ChainState) *CsChainService {
	futureBlocks, _ := lru.New(maxFutureBlocks)

	if chain_config.GetCurBootsEnv() == "test" {
		futureBlocks, _ = lru.New(maxFutureBlocksTestEnv)
	}

	//cache chain state
	ccs, err := NewCacheChainState(cs)
	if err != nil {
		panic(err)
	}
	// this step is needed, otherwise the cache can not be used
	cs.ChainStateConfig.WriterFactory.SetChain(ccs)

	service := &CsChainService{
		CsChainServiceConfig: config,
		CacheChainState:      ccs,
		FutureBlocks:         futureBlocks,
		Quit:                 make(chan struct{}),
	}

	if err := service.initService(); err != nil {
		panic(err)
	}

	return service
}

func (cs *CsChainService) Stop() {
	close(cs.Quit)
	cs.wg.Wait()
	log.Info("Blockchain manager stopped")
}

func (cs *CsChainService) CurrentBalance(address common.Address) *big.Int {
	curState, err := cs.CurrentState()
	if err != nil {
		log.Warn("get current state failed", "err", err)
		return nil
	}
	balance, err := curState.GetBalance(address)
	if err != nil {
		log.Info("get current balance failed", "err", err)
		return nil
	}
	log.Info("call current balance", "address", address.Hex(), "balance", balance)
	return balance
}

func (cs *CsChainService) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	state, err := cs.CurrentState()
	if err != nil {
		return 0, err
	}
	nonce, err = state.GetNonce(addr)
	if err != nil {
		return 0, err
	}

	return nonce, nil
}

//func (cs *CsChainService) GetVerifiers(slot uint64) []common.Address {
//	if vs, ok := cs.cachedVerifiers.Get(slot); ok {
//		return vs.([]common.Address)
//	}
//
//	// check round
//	config := cs.ChainState.ChainConfig
//	defaultVerifiers := chain.VerifierAddress[:config.VerifierNumber]
//
//	if slot < config.SlotMargin {
//		// replace by configured verifiers
//		return defaultVerifiers
//	}
//
//	cs.getVerifierLock.Lock()
//	defer cs.getVerifierLock.Unlock()
//
//	num := cs.numBeforeLastBySlot()
//	tmpB := cs.GetBlockByNumber(num)
//	if tmpB == nil {
//		panic(fmt.Sprintf("can't get block, num: %v", num))
//	}
//	// the slot in this function is not the same as current slot, because the block passed here is 2 rounds before.
//	cs.CalVerifiers(tmpB)
//	if vs, ok := cs.cachedVerifiers.Get(slot); ok {
//		return vs.([]common.Address)
//	}
//	panic(fmt.Sprintf("calVerifiers not gen cache. block num: %v, slot: %v", num, slot))
//}

func (cs *CsChainService) GetSeenCommit(height uint64) []model.AbstractVerification {
	if height == 0 {
		return nil
	}

	log.Info("load seen commits", "height", height)

	result, err := cs.CacheDB.GetSeenCommits(height, common.Hash{})

	if err != nil {
		log.Warn("read seen commits failed", "err", err)
	}

	return result

}

func (cs *CsChainService) SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	cs.wg.Add(1)
	defer cs.wg.Done()

	timer := g_metrics.NewTimer(g_metrics.InsertOneBlockDuration)
	defer timer.ObserveDuration()

	cs.saveBlockLock.Lock()
	defer cs.saveBlockLock.Unlock()

	if err := cs.checkBftBlock(block, seenCommits); err != nil {
		return err
	}

	if err := cs.saveBftBlock(block, seenCommits); err != nil {
		g_metrics.Add(g_metrics.FailedInsertBlockCount, "", 1)
		return err
	}

	curHeight := cs.CurrentBlock().Number()
	g_metrics.Set(g_metrics.CurChainHeight, "", float64(curHeight))
	log.PBft.Debug("Save Block Success", "block height", block.Number(), "chain height", curHeight)

	//metric tps
	if curHeight != 0 {
		lastBlock := cs.GetBlockByNumber(curHeight - 1)
		nowTimestamp := block.Timestamp().Int64()
		lastTimestamp := lastBlock.Timestamp().Int64()
		totalSec := float64(nowTimestamp-lastTimestamp) / float64(1e9)
		tps := float64(block.TxCount()) / totalSec
		g_metrics.Set(g_metrics.TpsValue, "", tps)

		g_metrics.Set(g_metrics.BlockTxNumber, "", float64(block.TxCount()))
	}
	return nil
}

func (cs *CsChainService) checkBftBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	// todo this can be optimized in middleware
	if block.Number() <= cs.CurrentBlock().Number() {
		log.PBft.Debug("fullChain#SaveBlock  Save previous height", "chain height", cs.CurrentBlock().Number(), "block height", block.Number())
		if cs.CurrentBlock().Number()-block.Number() > numLowBlockToReturnErr {
			return g_error.ErrAlreadyHaveThisBlock
		}

		return nil
	}

	return nil
}

func (cs *CsChainService) saveBftBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	oldCurrentHead := cs.CurrentHeader().(*model.Header)
	tmpBlock := cs.GetBlockByNumber(block.Number())
	var txs []model.AbstractTransaction
	if tmpBlock != nil {
		transactions := tmpBlock.GetTransactions()
		txs = make([]model.AbstractTransaction, len(transactions))
		util.InterfaceSliceCopy(txs, transactions)
	}

	// save block
	err := cs.SaveBftBlock(block, seenCommits)
	switch err {
	case nil:
		if err = cs.CacheDB.SaveSeenCommits(block.Number(), common.Hash{}, seenCommits); err != nil {
			log.PBft.Error("save seenCommits failed", "err", err)
			return err
		}

		// update tx pool if insert block
		newCurrentHeader := cs.CurrentHeader().(*model.Header)
		newCurrentBlock := cs.CurrentBlock()
		cs.TxPool.Reset(oldCurrentHead, newCurrentHeader)
		if newCurrentBlock.IsSpecial() && newCurrentBlock.Hash() == block.Hash() {
			// roll back txs
			cs.TxPool.AddRemotes(txs)

			// roll back block number
			for i := newCurrentBlock.Number(); i < oldCurrentHead.Number; i++ {
				cs.ChainState.Rollback(i + 1)
			}
		}

		// check future block
		cs.FutureBlocks.Remove(block.Hash())

		// insert success then calculate verifiers
		if cs.IsChangePoint(block, false) {
			cs.CalVerifiers(block)
		}

		sendBlock := block.(*model.Block)
		g_event.Send(g_event.NewBlockInsertEvent, *sendBlock)

		return nil

	case g_error.ErrFutureBlock:
		cs.FutureBlocks.Add(block.Hash(), &futureBlock{block: block, seenCommits: seenCommits})
		return nil

	default:
		log.Error("CsChainService saveBftBlock error", "err", err)
		return err
	}
}

type futureBlock struct {
	block       model.AbstractBlock
	seenCommits []model.AbstractVerification
}

// todo genesis check
func (cs *CsChainService) checkGenesis() {
	if cs.ChainState.StateStorage == nil || cs.ChainState.ChainDB == nil {
		panic("you need new chain state first")
	}

	genesisAccountStateProcessor, err := state_processor.MakeGenesisAccountStateProcessor(cs.ChainState.StateStorage)
	if err != nil {
		panic("open account state processor for genesis failed: " + err.Error())
	}

	genesisRegisterProcessor, err := registerdb.MakeGenesisRegisterProcessor(cs.ChainState.StateStorage)
	if err != nil {
		panic("make registerDB processor for genesis failed: " + err.Error())
	}

	// setup genesis block
	defaultGenesis := chain.DefaultGenesisBlock(cs.ChainDB, genesisAccountStateProcessor, genesisRegisterProcessor,
		cs.ChainState.ChainConfig)

	// todo no need dataDir
	if _, _, err = chain.SetupGenesisBlock(defaultGenesis); err != nil {
		panic("setup genesis block failed: " + err.Error())
	}
}

func (cs *CsChainService) InitService() error {
	var err error

	cs.initOnce.Do(func() {
		err = cs.initService()
	})

	return err
}

// mainly check the chain data here
func (cs *CsChainService) initService() error {
	if IsSetUpGenesis() {
		cs.checkGenesis()
	}

	// check genesis block
	if genesisBlock := cs.Genesis(); genesisBlock == nil {
		return g_error.ErrNoGenesis
	}

	headBlockHash := cs.ChainDB.GetHeadBlockHash()
	currentBlock := cs.GetBlockByHash(headBlockHash)

	cs.currentBlock.Store(currentBlock)
	currentHeader := currentBlock.Header()
	cs.currentHeader.Store(currentHeader)

	// Update cached verifier
	currentSlot := cs.GetSlot(currentBlock)
	if *currentSlot >= cs.GetChainConfig().SlotMargin {
		lastNum := cs.ChainState.NumBeforeLastBySlot(*currentSlot)
		if lastNum == nil {
			return g_error.ErrLastNumIsNil
		}
		cs.CalVerifiers(cs.GetBlockByNumber(*lastNum))
	}
	if *currentSlot >= 1 {
		lastPoint := cs.GetLastChangePoint(currentBlock)
		cs.CalVerifiers(cs.GetBlockByNumber(*lastPoint))
	}

	log.Info("Loaded most recent local header", "number", currentHeader.GetNumber(), "hash", currentHeader.Hash().Hex())
	log.Info("Loaded most recent local full block", "number", currentBlock.Number(), "hash", currentBlock.Hash().Hex())
	log.Info("initChain", "Chain version", cs.genesisBlock.Version())

	// handle future block
	go cs.handleFutureBlockTask()

	return nil
}

func (cs *CsChainService) handleFutureBlockTask() {
	// 5s update chain future block
	tickHandler := func() { cs.handleFutureBlock() }
	futureTimer := g_timer.SetPeriodAndRun(tickHandler, 5*time.Second)
	defer g_timer.StopWork(futureTimer)
	for {
		select {
		case <-cs.Quit:
			return
		}
	}
}

// handle chain future block set
func (cs *CsChainService) handleFutureBlock() {
	fb := make([]*futureBlock, 0, cs.FutureBlocks.Len())
	for _, hash := range cs.FutureBlocks.Keys() {
		if tmpFb, exist := cs.FutureBlocks.Peek(hash); exist {
			fb = append(fb, tmpFb.(*futureBlock))
		}
	}

	if len(fb) > 0 {
		// sort blocks
		sort.Slice(fb, func(i, j int) bool {
			return fb[i].block.Number() < fb[j].block.Number()
		})

		for i := range fb {
			tmpFb := fb[i]
			if err := cs.SaveBlock(tmpFb.block, tmpFb.seenCommits); err != nil {
				cs.FutureBlocks.Remove(tmpFb.block.Hash())
			}
		}
	}
}
