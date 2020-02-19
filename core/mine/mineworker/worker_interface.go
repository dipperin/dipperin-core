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
)

// external operation to the miner
type Worker interface {
	// start mining and register to master
	Start()
	// stop mining and logout from master
	Stop()
	// set coinbase
	SetCoinbaseAddress(address common.Address)
	// consult the current coinbase
	CurrentCoinbaseAddress() common.Address
}

// Work comes with the function of modifying the random number and verifying compliance
type workExecutor interface {
	// Modify the nonce and verify that it is qualified. If it passes, it returns true.
	ChangeNonce() bool
	// submit the mission
	Submit()
}

// the class to execute mining
type miner interface {
	startMine()
	stopMine()
	receiveWork(work workExecutor)
}
