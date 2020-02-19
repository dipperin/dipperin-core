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

// context must have coinbaseAddress common.Address, workBuilder workBuilder, blockBuilder blockBuilder, blockSubmitter blockSubmitter) (Master, MasterServer
// make a mine master
func MakeMineMaster(config MineConfig) (Master, MasterServer) {
	// interface of minemaster
	master := newMaster(config)

	//log.Info("the coinbaseAddress is: ","coinbaseAddress", config.CoinbaseAddress.Hex())

	dispatcher := newWorkDispatcher(config, master.Workers)
	// manage worker's work and submit work to broadcast
	manager := newDefaultWorkManager(config)
	// communicator for master with workers
	server := newServer(master, manager, dispatcher.curWorkBlock)
	master.workManager = manager
	// set depends
	master.setWorkDispatcher(dispatcher)
	return master, server
}
