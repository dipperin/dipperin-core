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

package accounts

import (
	"errors"
)

var ErrNotSupportUsbWallet = errors.New("not support USB wallet")

var ErrNotFindWallet = errors.New("not find the wallet")

var ErrInvalidKDFParameter = errors.New("invalid KDFParameter")

var ErrDeriveKey = errors.New("symmetric key derivation error")

var ErrAESInvalidParameter = errors.New("AES operation parameter error")

var ErrAESDecryption = errors.New("AES decryption operation error")

var ErrMacAuthentication = errors.New("MAC authentication error")

var ErrWalletNotOpen = errors.New("wallet isn't open")

var ErrWalletFileExist = errors.New("wallet file exist")

var ErrWalletFileNotExist = errors.New("wallet file doesn't exist")

var ErrWalletSendTransaction = errors.New("wallet send transaction error")

var ErrAddressBalanceNotEnough = errors.New("the address balance isn't enough when wallet send transaction")

var ErrWalletPasswordNotValid = errors.New("wallet password error")

var ErrDeleteWalletFile = errors.New("delete wallet error")

var ErrNotSupported = errors.New("not supported")

var ErrInvalidAddress = errors.New("invalid address")

var ErrAnalysisDerivedPath = errors.New("invalid Derived Path")

var ErrInvalidDerivedPath = errors.New("invalid derived path")

var ErrPasswordIsNil = errors.New("password is nil")

var ErrWalletPathError = errors.New("the path should be in the home path")

var ErrEmptySign = errors.New("empty sign")

var ErrSignatureInvalid = errors.New("verify signature fail")
