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


package g_error

import "errors"

var (
	AccountNotExist    = errors.New("account does not exist")
	BalanceNegErr      = errors.New("balance can not be negtive")
	NotHavePreBlockErr = errors.New("not have pre block")
	UnknownTxTypeErr   = errors.New("unknown tx type")
	InvalidVerifierAddressErr = errors.New("invalid verifier address")
	InvalidCoinBaseAddressErr = errors.New("invalid coinBase address")
	ErrNonceTooHigh = errors.New("nonce too high")
	ErrNonceTooLow = errors.New("nonce too low")
	ErrGasLimitReached = errors.New("gas limit reached")
)