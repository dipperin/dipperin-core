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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"time"
	"fmt"
	"github.com/dipperin/dipperin-core/common/g-event"
)

var waitTimeout = 20 * time.Second

func newMaster(config MineConfig) *master {
	m := &master{
		MineConfig: config,

		workers:              map[WorkerId]WorkerForMaster{},
		registerWorkerChan:   make(chan WorkerForMaster),
		unRegisterWorkerChan: make(chan WorkerId),
		onNewBlockChan:       make(chan model.AbstractBlock),
		startTimerChan:       make(chan struct{}),
		timeoutChan:          make(chan struct{}),
	}
	return m
}

type master struct {
	MineConfig

	workers map[WorkerId]WorkerForMaster

	workDispatcher dispatcher
	workManager    workManager

	// control the reception of OnNewBlock to prevent the repeated launch of reset task
	curNewBlockHeight uint64

	registerWorkerChan   chan WorkerForMaster
	unRegisterWorkerChan chan WorkerId
	onNewBlockChan       chan model.AbstractBlock
	startTimerChan       chan struct{}
	timeoutChan          chan struct{}
	stopChan             chan struct{}

	stopTimerFunc func()
}

func (ms *master) RetrieveReward(address common.Address) {
	ms.sendReward(address)

	ms.workManager.clearReward(address)
	ms.workManager.clearPerformance(address)

}

func (ms *master) sendReward(address common.Address) {
	//reward := ms.workManager.getReward(address)
	// ms.send(reward)
	// TODO, initiate transaction
}

func (ms *master) GetReward(address common.Address) *big.Int {
	return ms.workManager.getReward(address)
}

func (ms *master) GetPerformance(address common.Address) uint64 {
	return ms.workManager.getPerformance(address)
}

/*

to handle:
1. the arrival of new block
2. waiting reset task expires

*/
func (ms *master) loop() {
	newBlockInsertChan := make(chan model.Block)
	//sb := ms.nodeContext.ChainReader().SubscribeBlockEvent(newBlockInsertChan)
	sb := g_event.Subscribe(g_event.NewBlockInsertEvent, newBlockInsertChan)
	defer sb.Unsubscribe()

	for {
		select {
		case b := <-newBlockInsertChan:
			ms.doOnNewBlock(&b)

		case worker := <-ms.registerWorkerChan:
			//if ms.workers[worker.GetId()] != nil {
			//	log.Warn("register WorkerForMaster, but WorkerForMaster already exist in mine master, replace old", "worker id", worker.GetId(), "workers len", len(ms.workers))
			//	//continue
			//}
			ms.workers[worker.GetId()] = worker

		case wId := <-ms.unRegisterWorkerChan:
			log.Info("un register worker", "w id", wId)
			if ms.workers[wId] != nil {
				ms.workers[wId].Stop()
			}
			delete(ms.workers, wId)

		case block := <-ms.onNewBlockChan:
			ms.doOnNewBlock(block)

		case <-ms.startTimerChan:
			log.Info("start mine wait timer")
			ms.stopWait()

			// set a time of 3 seconds
			ms.stopTimerFunc = util.SetTimeout(func() {
				if ms.stopped() {
					return
				}
				ms.timeoutChan <- struct{}{}
			}, waitTimeout)

		case <-ms.timeoutChan:
			//log.Info("on timeout chan")
			ms.stopWait()
			log.Warn("wait block verified timeout, dispatch a new work")
			if err := ms.workDispatcher.dispatchNewWork(); err != nil {
				log.Error(fmt.Sprintf("master dispatch work failed, err: %v, w len: %v", err, len(ms.workers)))
			}

		case <-ms.stopChan:
			//log.Info("on stop chan")

			ms.stopWait()
			return
		}
	}
	//log.Info("finish mine master loop")
}

func (ms *master) getWorker(id WorkerId) WorkerForMaster {
	return ms.workers[id]
}

// return true if the mining is ongoing
func (ms *master) Mining() bool {
	return !ms.stopped()
}

func (ms *master) MineTxCount() int {
	if ms.workDispatcher.curWorkBlock() == nil {
		return 0
	}
	return ms.workDispatcher.curWorkBlock().TxCount()
}

func (ms *master) Start() {
	log.Info("start mine master", "worker len", len(ms.workers))
	if !ms.stopped() {
		log.Info("call miner master start, bug it already started")
		return
	}
	ms.stopChan = make(chan struct{})

	log.Info("run mine master loop")
	go ms.loop()
	// wait local worker connect

	log.Info("start mine workers", "worker len", len(ms.workers))
	for _, w := range ms.workers {
		w.Start()
	}

	time.Sleep(time.Second)

	if err := ms.workDispatcher.dispatchNewWork(); err != nil {
		panic(fmt.Sprintf("master dispatch work failed, err: %v, w len: %v", err, len(ms.workers)))
	}
}

func (ms *master) Stop() {
	if ms.stopped() {
		log.Info("call miner master stop, but it already stopped")
		return
	}

	for _, w := range ms.workers {
		w.Stop()
	}
	close(ms.stopChan)
}

func (ms *master) CurrentCoinbaseAddress() (result common.Address) {
	return ms.GetCoinbaseAddr()
}

func (ms *master) SetCoinbaseAddress(addr common.Address) {
	ms.CoinbaseAddress.Store(addr)
}

func (ms *master) SetMineGasConfig(gasFloor, gasCeil uint64) {
	ms.GasFloor.Store(gasFloor)
	ms.GasCeil.Store(gasCeil)
}

func (ms *master) stopped() bool {
	return util.StopChanClosed(ms.stopChan)
}

func (ms *master) OnNewBlock(block model.AbstractBlock) {
	if ms.stopped() {
		return
	}

	ms.onNewBlockChan <- block
}

//
func (ms *master) doOnNewBlock(block model.AbstractBlock) {
	log.Info("on new block chan", "new block", block.Number(), "cur block height", ms.curNewBlockHeight)

	if block.Number() <= ms.curNewBlockHeight {
		return
	}

	if block.CoinBaseAddress().IsEqual(ms.CurrentCoinbaseAddress()) {
		// if the received block is mined by ourselves
		ms.workManager.onNewBlock(block)
	}
	if err := ms.workDispatcher.onNewBlock(block); err == nil {
		ms.curNewBlockHeight = block.Number()
		ms.stopWait()
	} else {
		log.Warn("dispatcher process new block failed", "err", err)
	}
}

func (ms *master) Workers() map[WorkerId]WorkerForMaster {
	//log.Info("get mine master workers", "workers len", len(ms.workers))
	return ms.workers
}

func (ms *master) registerWorker(worker WorkerForMaster) {
	ms.registerWorkerChan <- worker
}

func (ms *master) unRegisterWorker(workerId WorkerId) {
	ms.unRegisterWorkerChan <- workerId
}

func (ms *master) setWorkDispatcher(dispatcher dispatcher) {
	ms.workDispatcher = dispatcher
}

func (ms *master) stopWait() {
	if ms.stopTimerFunc == nil {
		log.Info("no timer to stop")
		return
	}
	log.Info("stop timer")
	ms.stopTimerFunc()
	ms.stopTimerFunc = nil
}

// if it remains waiting status 5 seconds after a block
// is submitted by miners, a new block generation is automated
func (ms *master) startWaitTimer() {
	//log.Info("start wait timer")
	if ms.stopped() {
		//log.Info("miner stopped not start timer")
		return
	}
	ms.startTimerChan <- struct{}{}
}
