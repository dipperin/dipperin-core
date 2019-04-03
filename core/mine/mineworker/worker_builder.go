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
)

// make a local worker
func MakeLocalWorker(coinbaseAddr common.Address, workerCount int, master minemaster.MasterServer) Worker {
	// init conn
	conn := newLocalConnector("local", master)

	// init worker
	w := newWorker(coinbaseAddr, workerCount, conn)

	// init work manager
	manager := newWorkManager(conn, w.Miners, w.CurrentCoinbaseAddress)

	// set conn msg send to work manager
	conn.receiver = manager
	// set worker for receive stop msg
	conn.worker = w

	return w
}

// make a remote worker
func MakeRemoteWorker(coinbaseAddr common.Address, workerCount int) (Worker, *RemoteConnector) {
	// init conn
	conn := newRemoteConnector()

	// init worker
	w := newWorker(coinbaseAddr, workerCount, conn)

	// init work manager
	manager := newWorkManager(conn, w.Miners, w.CurrentCoinbaseAddress)

	// set conn msg send to work manager
	conn.receiver = manager
	// set worker for receive stop msg
	conn.worker = w

	return w, conn
}