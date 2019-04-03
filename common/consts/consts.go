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


package consts

import (
	"math/big"
)

const (
    MinAmount = 0
)
var (
	MaxAmount = big.NewInt(0).Mul(big.NewInt(1000000000000), big.NewInt(DIP))
)
// currency unit
const (
	// gray but other places to use
	DIP_Units = 1
	DIP_Thousands = 1000
	DIP_Millions = 1000 * 1000
	DIP = 1000 * 1000 * 1000
)

// coin digits
const(
	DIPDecimalBits = 9
)

// ninimum currency unit name
const (
	CoinDIPName 		 = " DIP"
)

// contract name configuration
const (
	ERC20TypeName = "ERC20"
	EarlyTokenTypeName = "EarlyReward"
)


