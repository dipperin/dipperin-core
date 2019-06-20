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
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
)

// context must have workBuilder workBuilder, blockBuilder blockBuilder, curCoinbaseAddressFunc curCoinbaseAddressFunc
func newWorkDispatcher(config MineConfig, getWorkersFunc getWorkersFunc) *workDispatcher {
	return &workDispatcher{
		MineConfig:     config,
		getWorkersFunc: getWorkersFunc,
	}
}

type getWorkersFunc func() map[WorkerId]WorkerForMaster

type workDispatcher struct {
	MineConfig

	curBlock       model.AbstractBlock
	getWorkersFunc getWorkersFunc
	//workBuilder    workBuilder
	//blockBuilder   blockBuilder
	//curCoinbaseAddressFunc curCoinbaseAddressFunc
}

func (dispatcher *workDispatcher) onNewBlock(block model.AbstractBlock) error {
	// new block num equal or bigger than cur work block num, reset work and dispatch a new
	if dispatcher.curWorkBlock() != nil {
		if block.Number() < dispatcher.curWorkBlock().Number() {
			//log.Warn("new block is smaller than cur work, nothing to do", "block num", block.Number())
			return fmt.Errorf("new block is smaller than cur work, nothing to do block num: %v", block.Number())
		}
	}

	//log.Info("mine dispatcher receive new block", "b num", block.Number(), "cur b num", dispatcher.curWorkBlock().Number())

	if err := dispatcher.dispatchNewWork(); err != nil {
		//log.Warn("new block come in, but dispatch work failed", "err", err, "worker len", len(dispatcher.getWorkersFunc()))
		return fmt.Errorf("new block come in, but dispatch work failed err: %v, worker len: %v", err, len(dispatcher.getWorkersFunc()))
	}
	return nil
}

func (dispatcher *workDispatcher) dispatchNewWork() error {
	pbft_log.Debug("dispatch mine work")
	workers := dispatcher.getWorkersFunc()
	workersLen := len(workers)
	if workersLen == 0 {
		log.Error("no worker to dispatch work", "worker len", workersLen)
		return errors.New("no worker to dispatch work")
	}
	workMsgCode, works := dispatcher.makeNewWorks(workersLen)

	if len(works) < workersLen {
		log.Error("can't dispatch work, gen works failed", "works len", len(works))
		return errors.New("can't dispatch work, gen works failed")
	}

	//log.Debug("dispatch new work", "workers len", workersLen)
	i := 0
	for _, w := range workers {
		w.SendNewWork(workMsgCode, works[i])
		i++
	}
	pbft_log.Debug("finish dispatch mine work")
	log.Info("finish dispatch work")
	return nil
}

func (dispatcher *workDispatcher) makeNewWorks(workerLen int) (workMsgCode int, works []minemsg.Work) {
	pbft_log.Debug("make new works")
	coinBaseAddr := dispatcher.GetCoinbaseAddr()
	gasFloor := dispatcher.GetGasFloor()
	gasCeil := dispatcher.GetGasCeil()
	dispatcher.curBlock = dispatcher.BlockBuilder.BuildWaitPackBlock(coinBaseAddr, gasFloor, gasCeil)
	mineWorkBuilder := minemsg.MakeDefaultWorkBuilder()
	return mineWorkBuilder.BuildWorks(dispatcher.curBlock, workerLen)
}

func (dispatcher *workDispatcher) curWorkBlock() model.AbstractBlock {
	return dispatcher.curBlock
}
