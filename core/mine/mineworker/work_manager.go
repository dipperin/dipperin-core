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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/third-party/log"
)

var (
	UnknownMsgCodeErr = errors.New("unknown msg code")
)

func newWorkManager(msgSender msgSender, getMinersFunc getMinersFunc, getCoinbaseAddressFunc getCoinbaseAddressFunc) *workManager {
	return &workManager{
		msgSender:              msgSender,
		getMinersFunc:          getMinersFunc,
		getCoinbaseAddressFunc: getCoinbaseAddressFunc,
		executorBuilder:        NewDefaultExecutorBuilder(),
	}
}

type getMinersFunc func() []miner

type getCoinbaseAddressFunc func() common.Address

type msgSender interface {
	SendMsg(code uint64, msg interface{}) error
}
type workMsg interface {
	MsgCode() int
	Decode(result interface{}) error
}

// build executors
type executorBuilder interface {
	CreateExecutor(msg workMsg, workCount int, submitter workSubmitter) (result []workExecutor, err error)
}

type workManager struct {
	getMinersFunc          getMinersFunc
	getCoinbaseAddressFunc getCoinbaseAddressFunc
	msgSender              msgSender
	executorBuilder        executorBuilder
}

func (workManager *workManager) OnNewWork(msg workMsg) {
	log.Info("work manager receive new work", "msg code", msg.MsgCode())
	miners := workManager.getMinersFunc()
	executors, err := workManager.executorBuilder.CreateExecutor(msg, len(miners), workManager)
	if err != nil {
		log.Warn("build executor for msg failed", "err", err)
		return
	}
	log.Info("dispatch work to miners", "executors len", len(executors))
	for i, executor := range executors {
		miners[i].receiveWork(executor)
	}
}

func (workManager *workManager) SubmitWork(work minemsg.Work) {
	work.SetWorkerCoinbaseAddress(workManager.getCoinbaseAddressFunc())
	if err := workManager.msgSender.SendMsg(minemsg.SubmitDefaultWorkMsg, work); err != nil {
		log.Warn("submit work failed", "err", err)
	}
}
