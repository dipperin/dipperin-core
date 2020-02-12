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

package components

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"fmt"
)

type poolEventNotifier interface {
	BlockPoolNotEmpty()
}

//type NodeContext interface {
//	ChainReader() state_processor.ChainReader
//}

//go:generate mockgen -destination=./pool_mock_test.go -package=components github.com/dipperin/dipperin-core/core/csbft/components Blockpoolconfig
type Blockpoolconfig interface {
	CurrentBlock() model.AbstractBlock
}

type BlockPool struct {
	// the height of the block to be saved
	height uint64
	blocks []model.AbstractBlock

	poolEventNotifier poolEventNotifier
	//context NodeContext
	Blockpoolconfig

	newHeightChan chan uint64
	newBlockChan  chan newBlockWithResultErr
	getterChan    chan *blockPoolGetter
	stopChan      chan struct{}
	rmBlockChan   chan common.Hash
}

type newBlockWithResultErr struct {
	block      model.AbstractBlock
	resultChan chan error
}

func NewBlockPool(height uint64, eventNotifier poolEventNotifier) *BlockPool {
	return &BlockPool{
		height:            height,
		blocks:            []model.AbstractBlock{},
		poolEventNotifier: eventNotifier,
		Blockpoolconfig:   nil,

		newHeightChan: make(chan uint64, 5),
		newBlockChan:  make(chan newBlockWithResultErr, 5),
		getterChan:    make(chan *blockPoolGetter, 5),
		rmBlockChan:   make(chan common.Hash),
		stopChan:      nil,
	}
}

//func (p *BlockPool) SetNodeContext(context NodeContext) {
//	p.context = context
//}
func (p *BlockPool) SetNodeConfig(config Blockpoolconfig) {
	p.Blockpoolconfig = config
}

// only init will call here, do not need lock
func (p *BlockPool) SetPoolEventNotifier(eventNotifier poolEventNotifier) {
	p.poolEventNotifier = eventNotifier
}

func (p *BlockPool) Start() error {
	if p.stopChan != nil {
		return errors.New("block pool already started")
	}

	if p.Blockpoolconfig != nil {
		//p.height = p.Blockpoolconfig.ChainReader.CurrentBlock().Number() + 1
		p.height = p.Blockpoolconfig.CurrentBlock().Number() + 1
	}
	p.stopChan = make(chan struct{})
	go p.loop()
	return nil
}

func (p *BlockPool) Stop() {
	if p.stopChan == nil {
		return
	}
	close(p.stopChan)
	p.stopChan = nil
}

func (p *BlockPool) IsEmpty() bool {
	return len(p.blocks) == 0
}

func (p *BlockPool) loop() {
	for {
		select {
		case h := <-p.newHeightChan:
			p.doNewHeight(h)
		case b := <-p.newBlockChan:
			p.doAddBlock(b)
		case c := <-p.getterChan:
			fmt.Println("gggggg")
			p.doGetBlock(c)
		case h := <-p.rmBlockChan:
			p.doRemoveBlock(h)
		case <-p.stopChan:
			fmt.Println("wwwww")
			return
		}
	}
}

func (p *BlockPool) doRemoveBlock(h common.Hash) {
	for i, b := range p.blocks {
		if b.Hash().IsEqual(h) {
			p.blocks = append(p.blocks[:i], p.blocks[i+1:]...)
			return
		}
	}
}

func (p *BlockPool) IsRunning() bool {
	return p.stopChan != nil
}

func (p *BlockPool) RemoveBlock(h common.Hash) {
	if p.IsRunning() {
		p.rmBlockChan <- h
	}
}

func (p *BlockPool) NewHeight(h uint64) {
	if p.IsRunning() {
		p.newHeightChan <- h
	}
}

// modify the height and empty blocks
func (p *BlockPool) doNewHeight(h uint64) {
	log.DLogger.Info("Update pool height", zap.Uint64("original height", p.height), zap.Uint64("new height", h))
	if h < p.height {
		log.DLogger.Warn("call block pool change to new height, but new height is lower than cur block pool height", zap.Uint64("pool height", p.height), zap.Uint64("new h", h))
		return
	}

	p.height = h
	p.blocks = []model.AbstractBlock{}
	log.DLogger.Debug("block pool", zap.Int("len", len(p.blocks)))
}

func (p *BlockPool) AddBlock(b model.AbstractBlock) error {
	if p.IsRunning() {
		resultChan := make(chan error)
		p.newBlockChan <- newBlockWithResultErr{block: b, resultChan: resultChan}
		return <-resultChan
	}
	return errors.New("block pool not running")
}

func (p *BlockPool) GetBlockByHash(h common.Hash) model.AbstractBlock {
	resultC := make(chan model.AbstractBlock)
	p.getBlock(&blockPoolGetter{
		blockHash:  h,
		resultChan: resultC,
	})
	return <-resultC
}

/*

The block here is better without checking the legitimacy. If it is illegally excluded, then when the proposal receives the block, it will consume the bandwidth and pull the block once.
Check if the height matches, if it doesn't match, it won't accept.
Exclude duplicate blocks based on hashExclude duplicate blocks based on hash.

*/
func (p *BlockPool) doAddBlock(nb newBlockWithResultErr) {
	b := nb.block

	log.DLogger.Debug("Pool received a block of height", zap.Uint64("height", b.Number()), zap.Uint64("pool height", p.height))
	if b.Number() != p.height {
		log.DLogger.Debug("receive invalid height block", zap.Uint64("b", b.Number()), zap.Uint64("p", p.height))

		nb.resultChan <- errors.New("invalid height block")
		return
	}
	for _, oldB := range p.blocks {
		// delete repeated block
		if oldB.Hash().IsEqual(b.Hash()) {
			//log.DLogger.Info("receive dul block")

			nb.resultChan <- errors.New("dul block")
			return
		}
		log.DLogger.Info("the oldB in block pool", zap.String("blockHash", oldB.Hash().Hex()))
	}
	log.DLogger.Info("the add block in block pool", zap.String("blockHash", b.Hash().Hex()))
	p.blocks = append(p.blocks, b)
	log.DLogger.Debug("pool length", zap.Uint64("height", p.height), zap.Int("len", len(p.blocks)))

	// send result
	nb.resultChan <- nil

	if len(p.blocks) == 1 {
		// notify have new block in pool
		p.poolEventNotifier.BlockPoolNotEmpty()
	}
}

func (p *BlockPool) GetProposalBlock() model.AbstractBlock {
	log.DLogger.Info("[GetProposalBlock] start~~~~~~~~~~~~~~~")
	defer log.DLogger.Info("[GetProposalBlock] end~~~~~~~~~~~~~~~")
	resultC := make(chan model.AbstractBlock)
	p.getBlock(&blockPoolGetter{
		resultChan: resultC,
	})
	select {
	case block := <-resultC:
		if block != nil {
			p.RemoveBlock(block.Hash())
		}
		return block
	}
}
func (p *BlockPool) getBlock(getter *blockPoolGetter) {
	if p.IsRunning() {
		p.getterChan <- getter
		return
	}
	log.DLogger.Warn("call get block from pool, but pool not started")
}

/*

Get the block, if the hash is not passed, it means that it's the proposal who takes the proposal block.
If the hash is passed, it means that the block matching the master is obtained.

*/
func (p *BlockPool) doGetBlock(getter *blockPoolGetter) {
	var result model.AbstractBlock = nil
	log.DLogger.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~get Block From blockPool~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	log.DLogger.Info("doGetBlock the blockPool len is:", zap.Int("len", len(p.blocks)))
	log.DLogger.Info("want get block hash is:", zap.String("id", getter.blockHash.Hex()))
	// proposer get first block
	if getter.blockHash.IsEqual(common.Hash{}) {
		if len(p.blocks) == 0 {
			result = nil
			fmt.Println("yyyyyy11111111")
		} else {
			result = p.blocks[0]
			fmt.Println("yyyyyy22222222")
		}
		// get match hash block
	} else {
		for _, b := range p.blocks {
			log.DLogger.Info("block in block Pool", zap.String("id", b.Hash().Hex()))
			if b.Hash().IsEqual(getter.blockHash) {
				result = b
				break
			}
		}
	}
	getter.resultChan <- result
}

type blockPoolGetter struct {
	blockHash  common.Hash
	resultChan chan model.AbstractBlock
}
