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
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"reflect"
)

func newServer(mineMaster mineMaster, wManager workManager, getCurWorkBlockFunc getCurWorkBlockFunc) *server {
	return &server{
		master:              mineMaster,
		getCurWorkBlockFunc: getCurWorkBlockFunc,
		workManager:         wManager,
	}
}

type getCurWorkBlockFunc func() model.AbstractBlock

type server struct {
	master              mineMaster
	getCurWorkBlockFunc getCurWorkBlockFunc
	workManager         workManager
}

func (s *server) RegisterWorker(worker WorkerForMaster) {
	s.master.registerWorker(worker)
}

func (s *server) UnRegisterWorker(workerId WorkerId) {
	s.master.unRegisterWorker(workerId)
}

func (s *server) ReceiveMsg(workerID WorkerId, code uint64, msg interface{}) {
	switch code {
	case minemsg.SubmitDefaultWorkMsg:
		w, ok := msg.(minemsg.Work)
		if !ok {
			log.Warn("receive wrong work submit msg", "work", reflect.TypeOf(msg))
			return
		}

		// add timeout for wait new block event
		s.master.startWaitTimer()
		// TODO: verify different block difficulty
		s.onSubmitBlock(workerID, w)
	default:
		log.Debug("receive wrong msg", "code", code)
	}
}

func (s *server) onSubmitBlock(workerID WorkerId, work minemsg.Work) {

	block := s.getCurWorkBlockFunc()
	log.PBft.Debug("onSubmitBlock", "block id", block.Number(), "block txs", block.TxCount())
	if err := work.FillSealResult(block); err != nil {
		log.Warn("fill seal result failed", "err", err)
		return
	}

	log.Info("mine master before submit block", "hash", block.RefreshHashCache())
	// check block valid
	if !block.RefreshHashCache().ValidHashForDifficulty(block.Difficulty()) {
		log.Warn("master receive invalid mined block", "do unregister worker", workerID)
		//s.UnRegisterWorker(workerID)
		return
	}

	receiptHash := block.GetReceiptHash()
	//bloomLog := block.GetBloomLog()
	log.Info("server#onSubmitBlock", "receipts", receiptHash)

	//fmt.Println("mine master prepare broadcast block", util.StringifyJson(block), block.Hash())
	//log.Info("mine master receive new work", "block hash", block.Hash().Hex(), "block number", block.Number())
	s.workManager.submitBlock(work.GetWorkerCoinbaseAddress(), block)
}

// only for worker, do nothing
func (s *server) SetMineMasterPeer(peer chain_communication.PmAbstractPeer) {}

// receive worker msg
func (s *server) OnNewMsg(msg p2p.Msg, p chain_communication.PmAbstractPeer) error {
	workerId := WorkerId(p.ID())
	switch msg.Code {
	case minemsg.RegisterMsg:
		var register minemsg.Register
		if err := msg.Decode(&register); err != nil {
			return err
		}

		if remoteW := newRemoteWorker(p, register.Coinbase, workerId); remoteW == nil {
			return errors.New("invalid worker")
		} else {
			s.RegisterWorker(remoteW)
		}

		return nil
	case minemsg.UnRegisterMsg:
		log.Info("receive un register msg", "worker id", workerId)
		s.UnRegisterWorker(workerId)

	case minemsg.SetCurrentCoinbaseMsg:
		var setCoinbaseReq minemsg.SetCurrentCoinbase
		if err := msg.Decode(&setCoinbaseReq); err != nil {
			return err
		}

		if w := s.master.getWorker(workerId); w == nil {
			return nil
		} else {
			w.SetCoinbase(setCoinbaseReq.Coinbase)
		}
		return nil

	case minemsg.SubmitDefaultWorkMsg:
		var defaultWork minemsg.DefaultWork
		if err := msg.Decode(&defaultWork); err != nil {
			return err
		}
		s.ReceiveMsg(workerId, msg.Code, &defaultWork)

	default:
		log.Warn("receive unknown msg", "code", msg.Code)
		return errors.New("unknown msg")
	}
	return nil
}
