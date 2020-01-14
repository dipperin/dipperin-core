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

package model

import (
	"github.com/dipperin/dipperin-core/common"
	"math/big"
)

var (
	AlicePriv    = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	BobPriv      = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	AliceAddr    = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	BobAddr      = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
	CharlieAddr  = common.HexToAddress("0x00007dbbf084F4a6CcC070568f7674d4c2CE8CD2709E")
	ContractAddr = common.HexToAddress("0x0014B5Df12F50295469Fe33951403b8f4E63231Ef488")
	TestZeroValue = big.NewInt(0)
)
