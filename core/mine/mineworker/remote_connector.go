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
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chaincommunication"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/third_party/p2p"
)

func newRemoteConnector() *RemoteConnector {
	return &RemoteConnector{}
}

type RemoteConnector struct {
	peer chaincommunication.PmAbstractPeer

	// set when build
	worker   Worker
	receiver msgReceiver
}

func (conn *RemoteConnector) SetMineMasterPeer(peer chaincommunication.PmAbstractPeer) {
	conn.peer = peer
	log.DLogger.Info("connect master, do register")
	// peer connected do register
	go conn.Register()
	conn.worker.Start()
}

func (conn *RemoteConnector) Register() error {
	return conn.peer.SendMsg(minemsg.RegisterMsg, minemsg.Register{Coinbase: conn.worker.CurrentCoinbaseAddress()})
}

func (conn *RemoteConnector) UnRegister() {
	conn.peer.SendMsg(minemsg.UnRegisterMsg, "")
}

func (conn *RemoteConnector) SendMsg(code uint64, msg interface{}) error {
	return conn.peer.SendMsg(code, msg)
}

func (conn *RemoteConnector) OnNewMsg(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {
	switch msg.Code {
	case minemsg.StartMineMsg:
		conn.worker.Start()
	case minemsg.StopMineMsg:
		conn.worker.Stop()
	case minemsg.WaitForCommitMsg:
		// todo deal this msg

	case minemsg.NewDefaultWorkMsg:
		var work minemsg.DefaultWork
		if err := msg.Decode(&work); err != nil {
			return err
		}
		msg := &localWorkMsg{code: int(msg.Code), work: &work}
		conn.receiver.OnNewWork(msg)
	}
	return nil
}
