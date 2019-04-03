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
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/third-party/log"
	"sync/atomic"
)

func newRemoteWorker(peer chain_communication.PmAbstractPeer, curCoinbaseAddr common.Address, workerId WorkerId) *remoteWorker {
	if peer == nil || curCoinbaseAddr.IsEqual(common.Address{}) || workerId == "" {
		log.Warn("invalid remote worker", "curCoinbaseAddr", curCoinbaseAddr, "workerId", workerId)
		return nil
	}
	worker := &remoteWorker{
		peer:     peer,
		workerId: workerId,
	}
	worker.curCoinbaseAddr.Store(curCoinbaseAddr)
	return worker
}

type remoteWorker struct {
	peer chain_communication.PmAbstractPeer

	curCoinbaseAddr atomic.Value
	workerId        WorkerId
}

func (worker *remoteWorker) SetCoinbase(coinbase common.Address) {
	worker.curCoinbaseAddr.Store(coinbase)
}

func (worker *remoteWorker) CurrentCoinbaseAddress() common.Address {
	if addr := worker.curCoinbaseAddr.Load(); addr != nil {
		return addr.(common.Address)
	}
	return common.Address{}
}

func (worker *remoteWorker) Start() {
	worker.peer.SendMsg(minemsg.StartMineMsg, "")
}

func (worker *remoteWorker) Stop() {
	worker.peer.SendMsg(minemsg.StopMineMsg, "")
}

func (worker *remoteWorker) GetId() WorkerId {
	return worker.workerId
}

func (worker *remoteWorker) SendNewWork(msgCode int, work minemsg.Work) {
	worker.peer.SendMsg(uint64(msgCode), work)
}

func (worker *remoteWorker) WaitForCommit() {
	worker.peer.SendMsg(minemsg.WaitForCommitMsg, "")
}
