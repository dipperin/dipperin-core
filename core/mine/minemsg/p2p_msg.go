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


package minemsg

import "github.com/dipperin/dipperin-core/common"

const (
	// msg sent to worker by master
	StartMineMsg = 0x1
	StopMineMsg = 0x2
	WaitForCommitMsg = 0x3

	// work msg(one dispatch, one submit)
	NewDefaultWorkMsg = 0x10
	SubmitDefaultWorkMsg = 0x11

	// msg sent to master by worker
	RegisterMsg = 0x50
	UnRegisterMsg = 0x51
	SetCurrentCoinbaseMsg = 0x52
)

type Register struct {
	Coinbase common.Address
}

type SetCurrentCoinbase struct {
	Coinbase common.Address
}