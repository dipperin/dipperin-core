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

package chain_state

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
)

func (cs *ChainState) BuildRegisterProcessor(preBlockRegisterRoot common.Hash) (*registerdb.RegisterDB, error) {
	return registerdb.NewRegisterDB(preBlockRegisterRoot, cs.StateStorage, cs)
}

func (cs *ChainState) CurrentSeed() (common.Hash, uint64) {

	slot := cs.GetSlot(cs.CurrentBlock())
	seedNum := cs.NumBeforeLastBySlot(*slot)
	if seedNum == nil {
		panic("can't get future seed number")
	}
	seed := cs.GetBlockByNumber(*seedNum).Seed()
	return seed, cs.CurrentBlock().Number()
}

func (cs *ChainState) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	if block.Number() < 1 {
		return false
	}

	preBlock := cs.GetBlockByNumber(block.Number() - 1)
	process, err := cs.BuildRegisterProcessor(preBlock.GetRegisterRoot())
	point, err := process.GetLastChangePoint()
	if err != nil {
		log.Error("GetLastChangePoint failed", "err", err)
		return false
	}

	// The distance from the current block to the previous turning point, 0 for the first round of turning points, and the last block in each round for the subsequent rounds. For example: SlotSize is 10, then the second round of points is 9
	diff := block.Number() - point
	if point == 0 {
		diff += 1
	}
	slotSize := cs.ChainConfig.SlotSize

	// the mineMaster packaged Block diff and nonce is 0
	if isProcessPackageBlock {
		if diff == slotSize {
			return true
		}
	} else {
		if block.IsSpecial() || diff == slotSize {
			return true
		}
	}
	return false
}

func (cs *ChainState) GetLastChangePoint(block model.AbstractBlock) *uint64 {
	if block.Number() < 1 {
		num := uint64(0)
		return &num
	}

	preBlock := cs.GetBlockByNumber(block.Number() - 1)
	if preBlock.IsSpecial() {
		num := preBlock.Number()
		return &num
	}

	process, _ := cs.BuildRegisterProcessor(preBlock.GetRegisterRoot())
	point, err := process.GetLastChangePoint()
	if err != nil {
		log.Error("get last change point failed", "err", err)
		return nil
	}

	if process.IsChangePoint(preBlock, false) {
		num := preBlock.Number()
		return &num
	}
	return &point
}

// This layer of queries does not use the cache, so try to minimize the query call
func (cs *ChainState) GetSlotByNum(num uint64) *uint64 {
	if num < 1 {
		num := uint64(0)
		return &num
	}

	preBlock := cs.GetBlockByNumber(num - 1)
	if preBlock == nil {
		return nil
	}
	process, err := cs.BuildRegisterProcessor(preBlock.GetRegisterRoot())
	if err != nil {
		log.Error("can't BuildRegisterProcessor", "r root", preBlock.GetRegisterRoot())
		return nil
	}
	slot, err := process.GetSlot()
	if err != nil {
		log.Error("get slot failed", "err", err)
		return nil
	}

	if process.IsChangePoint(preBlock, false) {
		slot += 1
		return &slot
	}

	return &slot
}

func (cs *ChainState) GetSlot(block model.AbstractBlock) *uint64 {
	return cs.GetSlotByNum(block.Number())
}

func (cs *ChainState) GetCurrVerifiers() []common.Address {
	cb := cs.CurrentBlock()
	if cb == nil {
		log.Error("can't no get current block")
		return nil
	}

	slot := cs.GetSlot(cb)
	return cs.GetVerifiers(*slot)
}

func (cs *ChainState) GetVerifiers(slot uint64) []common.Address {
	// check round
	config := cs.GetChainConfig()
	defaultVerifiers := chain.VerifierAddress[:config.VerifierNumber]

	if slot < config.SlotMargin {
		// replace by configured verifiers
		return defaultVerifiers
	}

	num := cs.NumBeforeLastBySlot(slot)
	if num == nil {
		log.Error("get verifiers error", "slot", slot, "num", num)
		panic("can't get block number before the last ")
	}
	tmpB := cs.GetBlockByNumber(*num)

	if tmpB == nil {
		panic(fmt.Sprintf("can't get block, num: %v", num))
	}
	return cs.CalVerifiers(tmpB)
}

func (cs *ChainState) GetNextVerifiers() []common.Address {
	cb := cs.CurrentBlock()
	if cb == nil {
		log.Error("can't no get current block")
		return nil
	}
	slot := cs.GetSlot(cb)
	return cs.GetVerifiers(*slot + 1)
}

// Get the last block of the last two rounds. Its seed is needed to calculate verifiers
// maybe need a cache
func (cs *ChainState) NumBeforeLastBySlot(slot uint64) *uint64 {
	margin := cs.GetChainConfig().SlotMargin
	if slot < margin {
		num := uint64(0)
		return &num
	}

	// Cannot count the slot which lasts more than two rounds
	curBlock := cs.CurrentBlock()
	curSlot := cs.GetSlot(curBlock)
	if *curSlot < (slot - margin) {
		return nil
	}
	return cs.GetNumBySlot(slot - margin)
}

// The returned value is the last block of the slot. If the slot is incomplete, it returns null.
func (cs *ChainState) GetNumBySlot(slot uint64) *uint64 {
	block := cs.CurrentBlock()
	for {
		blockSlot := cs.GetSlot(block)
		if slot > *blockSlot {
			log.Error("input slot can't larger than current slot")
			return nil
		}
		if *blockSlot == slot {
			log.Info("GetNumBySlot return ", "input slot", slot, "return num", block.Number())
			if cs.IsChangePoint(block, false) {
				num := block.Number()
				return &num
			} else {
				log.Error("the last block in input slot doesn't exist")
				return nil
			}
		} else {
			num := cs.GetLastChangePoint(block)
			block = cs.GetBlockByNumber(*num)
		}
	}
}

// Calculate the verifiers after two rounds based on the last block of each round
func (cs *ChainState) CalVerifiers(block model.AbstractBlock) []common.Address {

	// get all registration data
	//pbft_log.Log.Debug("CalVerifiers", "num", block.Number())

	root := block.GetRegisterRoot()
	log.Info("the register root is:", "root", root.Hex())
	register, err := cs.BuildRegisterProcessor(root)
	if err != nil {
		pbft_log.Log.Debug("BuildRegisterProcessor failed", "err", err)
	}
	list := register.GetRegisterData()
	log.Info("the register list len is:", "len", len(list))
	//pbft_log.Log.Debug("GetRegisterData", "register data", list, "root", root)
	//log.Info("GetRegisterData", "register data", list, "root", root)

	// get top verifiers
	var topAddress []common.Address
	var topPriority []uint64
	for i := 0; i < len(list); i++ {
		priority, err := cs.calPriority(list[i], block.Number())
		if err != nil {
			pbft_log.Log.Info("calPriority", "err", err)
		}
		topAddress, topPriority = cs.getTopVerifiers(list[i], priority, topAddress, topPriority)
	}
	//pbft_log.Log.Info("getTopVerifiers", "topAddress", len(topAddress), "topPriority", topPriority)

	// angel nodes take the place
	config := cs.GetChainConfig()
	defaultVerifiers := chain.VerifierAddress[:config.VerifierNumber]
	for add := range defaultVerifiers {
		if len(topAddress) < config.VerifierNumber {
			topAddress, topPriority = cs.getTopVerifiers(defaultVerifiers[add], config.SystemVerifierPriority, topAddress, topPriority)
		}
	}

	//pbft_log.Log.Debug("Add cachedVerifiers success", "topAddress", topAddress, "slot", slot+config.SlotMargin)
	return topAddress
}

func (cs *ChainState) getTopVerifiers(address common.Address, priority uint64, topAddress []common.Address, topPriority []uint64) ([]common.Address, []uint64) {
	config := cs.GetChainConfig()
	if len(topAddress) < config.VerifierNumber {
		tmpIndex := 0
		for i := 0; i < len(topAddress); i++ {
			if topAddress[i].IsEqual(address) {
				return topAddress, topPriority
			}

			if priority < topPriority[i] {
				tmpIndex++
			}
		}

		recordAddress := append([]common.Address{}, topAddress[tmpIndex:]...)
		topAddress = append(append(topAddress[:tmpIndex], address), recordAddress...)

		recordPriority := append([]uint64{}, topPriority[tmpIndex:]...)
		topPriority = append(append(topPriority[:tmpIndex], priority), recordPriority...)

	} else {
		insertPosition := 0
		config := cs.GetChainConfig()
		if priority <= topPriority[config.VerifierNumber-1] {
			return topAddress, topPriority
		}

		//find the insert position
		for i := 0; i < len(topAddress); i++ {
			// check address
			if topAddress[i].IsEqual(address) {
				return topAddress, topPriority
			}

			if priority < topPriority[i] {
				insertPosition++
			}
		}

		//insert the priority and delete the smallest priority
		recordAddress := append([]common.Address{}, topAddress[insertPosition:config.VerifierNumber-1]...)
		topAddress = append(append(topAddress[:insertPosition], address), recordAddress...)

		recordPriority := append([]uint64{}, topPriority[insertPosition:config.VerifierNumber-1]...)
		topPriority = append(append(topPriority[:insertPosition], priority), recordPriority...)
	}
	return topAddress, topPriority
}

func (cs *ChainState) calPriority(addr common.Address, blockNum uint64) (uint64, error) {
	luck := cs.getLuck(addr, blockNum)
	state, err := cs.StateAtByBlockNumber(blockNum)
	if err != nil {
		return 0, err
	}

	accountNonce, err := state.GetNonce(addr)
	stake, err := state.GetStake(addr)
	performance, err := state.GetPerformance(addr)

	// todo take this shit to Ox Star Star
	priority, err := model.DefaultPriorityCalculator.GetElectPriority(luck, accountNonce, stake, performance)
	if err != nil {
		return 0, err
	}
	return priority, nil
}

func (cs *ChainState) getLuck(addr common.Address, blockNum uint64) common.Hash {
	seed := cs.GetBlockByNumber(blockNum).Seed()
	list := append(seed.Bytes(), addr.Bytes()...)
	return common.RlpHashKeccak256(list)
}
