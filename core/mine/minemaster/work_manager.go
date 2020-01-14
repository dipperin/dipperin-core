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

package minemaster

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"math/big"
	"sync"
)

func newDefaultWorkManager(config MineConfig) *defaultWorkManager {
	return &defaultWorkManager{
		MineConfig: config,

		performance: make(map[common.Address]workerPerformance),
		reward:      make(map[common.Address]*big.Int),
		totalReward: new(big.Int),
	}
}

type defaultWorkManager struct {
	MineConfig

	submitBlockLock sync.Mutex

	performance map[common.Address]workerPerformance
	reward      map[common.Address]*big.Int

	// wallet sums up all the rewards that this minemaster had received
	totalReward *big.Int
}

func (manager *defaultWorkManager) subtractPerformance(address common.Address, performance uint64) {
	if p := manager.performance[address]; p != nil {
		op := p.getPerformance()
		if op > performance {
			manager.performance[address].setPerformance(op - performance)
		} else {
			log.DLogger.Debug("reward is less than current performance", zap.Uint64("performance", performance), zap.Uint64("current performance", op))
		}
	} else {
		log.DLogger.Debug("address is invalid", zap.Any("address", address))
	}

}

func (manager *defaultWorkManager) subtractReward(address common.Address, reward *big.Int) {
	if r := manager.reward[address]; r != nil {
		if r.Cmp(reward) > 0 {
			newPer := r.Sub(r, reward)
			manager.clearReward(address)
			manager.reward[address] = newPer
		} else {
			log.DLogger.Debug("reward is less than current reward", zap.Any("reward", reward), zap.Any("current reward", r))
		}
	} else {
		log.DLogger.Debug("address is invalid", zap.Any("address", address))
	}
}

func (manager *defaultWorkManager) clearPerformance(address common.Address) {
	delete(manager.performance, address)
}

func (manager *defaultWorkManager) clearReward(address common.Address) {
	delete(manager.reward, address)
}

func (manager *defaultWorkManager) getReward(address common.Address) *big.Int {
	if manager.reward[address] == nil {
		return big.NewInt(0)
	}
	return manager.reward[address]
}

// divideReward updates the reward distribution every time it is called.
// it refreshes the reward map according to the performance and totalReward
func (manager *defaultWorkManager) divideReward(coinbase *big.Int) map[common.Address]*big.Int {
	res := make(map[common.Address]*big.Int)

	manager.totalReward.Add(manager.totalReward, coinbase)
	totalPerformance := uint64(0)

	// divide equally
	for _, v := range manager.performance {
		totalPerformance += v.getPerformance()
	}

	performance := big.NewInt(int64(totalPerformance))
	// could not avoid two loops
	for k, v := range manager.performance {
		num := big.NewInt(int64(v.getPerformance()))
		num.Mul(num, manager.totalReward)
		res[k] = num.Quo(num, performance)
	}

	// refresh the reward map
	manager.reward = res

	return manager.reward
}

func (manager *defaultWorkManager) onNewBlock(block model.AbstractBlock) {
	coinbase := block.CoinBase()
	txFees := block.GetTransactionFees()

	manager.divideReward(coinbase.Add(coinbase, txFees))
}

func (manager *defaultWorkManager) getPerformance(address common.Address) uint64 {
	// now performance equals reward
	if manager.performance[address] == nil {
		return 0
	}
	return manager.performance[address].getPerformance()
}

// defaultPerformance records the performance of a worker based on how many
// blocks a specific miner has mined.
type defaultPerformance struct {
	blocksMined uint64
}

func (p defaultPerformance) getPerformance() uint64 {
	return p.blocksMined
}

// updatePerformance increments the count for how much block a miner has mined
func (p *defaultPerformance) updatePerformance() {
	p.blocksMined++
}

func (p *defaultPerformance) setPerformance(b uint64) {
	p.blocksMined = b
}

func newDefaultPerformance() *defaultPerformance {
	return &defaultPerformance{}
}

func (manager *defaultWorkManager) submitBlock(workerAddress common.Address, block model.AbstractBlock) {
	// TODO: do something to the block, otherwise change the function signature
	manager.submitBlockLock.Lock()
	defer manager.submitBlockLock.Unlock()
	//pbft_log.DLogger.Debug("submitBlock","block id",block.Number(),"block txs",block.TxCount())

	// todo: here should be change a way to judge
	if manager.performance[workerAddress] == nil {
		manager.performance[workerAddress] = newDefaultPerformance()
	}
	manager.performance[workerAddress].updatePerformance()

	// broadcast block
	//pbft_log.DLogger.Debug("submitBlock broad cast block","block id",block.Number(),"block txs",block.TxCount())
	manager.BlockBroadcaster.BroadcastMinedBlock(block)
}
