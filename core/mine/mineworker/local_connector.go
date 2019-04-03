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
	"github.com/dipperin/dipperin-core/core/mine/minemaster"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/third-party/log"
	"reflect"
	"errors"
)

func newLocalConnector(wId minemaster.WorkerId, master minemaster.MasterServer) *localConnector {
	return &localConnector{
		id:          wId,
		localServer: master,
	}
}

type localConnector struct {
	// set by [func (conn *localConnector) SetWorker]
	worker Worker
	// set by [func (conn *localConnector) SetMsgReceiver]
	receiver msgReceiver

	id          minemaster.WorkerId
	localServer minemaster.MasterServer
}

func (conn *localConnector) Register() error {
	conn.localServer.RegisterWorker(conn)
	return nil
}

func (conn *localConnector) UnRegister() {
	conn.localServer.UnRegisterWorker(conn.id)
}

func (conn *localConnector) SendMsg(code uint64, msg interface{}) error {
	// interesting
	go conn.localServer.ReceiveMsg(conn.id, code, msg)
	return nil
}

func (conn *localConnector) WaitForCommit() {

}

// implement for minemaster worker

func (conn *localConnector) SetCoinbase(coinbase common.Address) {}

func (conn *localConnector) Start() {
	conn.worker.Start()
}
func (conn *localConnector) Stop() {
	conn.worker.Stop()
}
func (conn *localConnector) GetId() minemaster.WorkerId {
	return conn.id
}
func (conn *localConnector) SendNewWork(msgCode int, work minemsg.Work) {
	log.Debug("localConnector SendNewWork", "msg code", msgCode)
	msg := &localWorkMsg{code: msgCode, work: work}
	go conn.receiver.OnNewWork(msg)
}

func (conn *localConnector) CurrentCoinbaseAddress() common.Address {
	return conn.worker.CurrentCoinbaseAddress()
}

type localWorkMsg struct {
	code int
	work minemsg.Work
}

func (msg *localWorkMsg) MsgCode() int {
	return msg.code
}

var (
	WorkMsgDecodeShouldBeSameTypeErr  = errors.New("decode result arg type not match the msg")
	WorkMsgDecodeResultShouldBePtrErr = errors.New("decode result should be ptr")
)

func (msg *localWorkMsg) Decode(result interface{}) error {
	rv := reflect.ValueOf(result)
	// result should be ptr
	if rv.Kind() != reflect.Ptr {
		return WorkMsgDecodeResultShouldBePtrErr
	}
	rv = rv.Elem()
	wRv := reflect.ValueOf(msg.work)
	if wRv.Kind() == reflect.Ptr {
		wRv = wRv.Elem()
	}
	// result's type should equal msg.work's type
	if rv.Type() != wRv.Type() {
		return WorkMsgDecodeShouldBeSameTypeErr
	}
	rv.Set(wRv)
	return nil
}
