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

package gerror

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common/consts"
)

var (
	ErrNotSupportUsbWallet = errors.New("not support USB wallet")

	ErrNotFindWallet = errors.New("not find the wallet")

	ErrInvalidKDFParameter = errors.New("invalid KDFParameter")

	ErrDeriveKey = errors.New("symmetric key derivation error")

	ErrAESInvalidParameter = errors.New("AES operation parameter error")

	ErrAESDecryption = errors.New("AES decryption operation error")

	ErrMacAuthentication = errors.New("MAC authentication error")

	ErrWalletNotOpen = errors.New("wallet isn't open")

	ErrWalletFileExist = errors.New("wallet file exist")

	ErrWalletFileNotExist = errors.New("wallet file doesn't exist")

	ErrWalletSendTransaction = errors.New("wallet send transaction error")

	ErrAddressBalanceNotEnough = errors.New("the address balance isn't enough when wallet send transaction")

	ErrWalletPasswordNotValid = errors.New("wallet password error")

	ErrDeleteWalletFile = errors.New("delete wallet error")

	ErrNotSupported = errors.New("not supported")

	ErrInvalidAddress = errors.New("invalid address")

	ErrAnalysisDerivedPath = errors.New("invalid Derived Path")

	ErrInvalidDerivedPath = errors.New("invalid derived path")

	ErrPasswordIsNil = errors.New("password is nil")

	ErrPasswordOrPassPhraseIllegal = errors.New("password or passPhrase illegal, must between 8 and 24, and no chinese , no spaces!!! ")

	ErrWalletPathError = errors.New("the path should be in the home path")

	ErrEmptySign = errors.New("empty sign")

	ErrSignatureInvalid = errors.New("verify signature fail")

	ErrWalletManagerIsEmpty = errors.New("there isn't a wallet in wallet manager")

	ErrWalletManagerNotRunning = errors.New("the wallet manager isn't running")
	// ErrDeriveHardFromPublic describes an error in which the caller
	// attempted to derive a hardened extended key from a public key.
	ErrDeriveHardFromPublic = errors.New("cannot derive a hardened key " +
		"from a public key")

	// ErrDeriveBeyondMaxDepth describes an error in which the caller
	// has attempted to derive more than 255 keys from a root key.
	ErrDeriveBeyondMaxDepth = errors.New("cannot derive a key with more than " +
		"255 indices in its path")

	// ErrNotPrivExtKey describes an error in which the caller attempted
	// to extract a private key from a public extended key.
	ErrNotPrivExtKey = errors.New("unable to create private keys from a " +
		"public extended key")

	// ErrInvalidChild describes an error in which the child at a specific
	// index is invalid due to the derived key falling outside of the valid
	// range for secp256k1 private keys.  This error indicates the caller
	// should simply ignore the invalid child extended key at this index and
	// increment to the next index.
	ErrInvalidChild = errors.New("the extended key at this index is invalid")

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// secp256k1 private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", consts.MinSeedBytes*8, consts.MaxSeedBytes*8)

	// ErrBadChecksum describes an error in which the checksum encoded with
	// a serialized extended key does not match the calculated value.
	ErrBadChecksum = errors.New("bad extended key checksum")

	// ErrInvalidKeyLen describes an error in which the provided serialized
	// key is not the expected length.
	ErrInvalidKeyLen = errors.New("the provided serialized extended key " +
		"length is invalid")
)

