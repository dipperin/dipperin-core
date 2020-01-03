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

package mineworker

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"go.uber.org/zap"
	"sync/atomic"
)

func newWorker(coinbaseAddr common.Address, workerCount int, connector connector) *worker {
	worker := &worker{connector: connector}
	worker.SetCoinbaseAddress(coinbaseAddr)
	for i := 0; i < workerCount; i++ {
		worker.miners = append(worker.miners, NewMiner())
	}
	return worker
}

type msgReceiver interface {
	OnNewWork(msg workMsg)
}

type connector interface {
	Register() error
	UnRegister()
	SendMsg(code uint64, msg interface{}) error
	//SetMsgReceiver(receiver msgReceiver)
	//SetWorker(w Worker)
}

// manager miners and coinbase
type worker struct {
	miners          []miner
	coinbaseAddress atomic.Value
	connector       connector
}

func (worker *worker) Miners() []miner {
	return worker.miners
}

func (worker *worker) Start() {
	log.DLogger.Info("call worker start mine")
	if err := worker.register(); err != nil {
		log.DLogger.Error("register to master failed, can't start worker", zap.Error(err))
		return
	}
	// start miners
	for _, m := range worker.miners {
		m.startMine()
	}
}

func (worker *worker) register() error {
	go worker.connector.Register()
	return nil
}

func (worker *worker) Stop() {
	log.DLogger.Info("stop worker's miner")
	// stop miners
	for _, m := range worker.miners {
		m.stopMine()
	}
}

func (worker *worker) unRegister() {
	worker.connector.UnRegister()
}

// use pointer to change value
func (worker *worker) SetCoinbaseAddress(address common.Address) {
	worker.coinbaseAddress.Store(address)
}

func (worker *worker) CurrentCoinbaseAddress() common.Address {
	if addr := worker.coinbaseAddress.Load(); addr != nil {
		return addr.(common.Address)
	}
	return common.Address{}
}
