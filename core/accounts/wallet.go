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
	"crypto/ecdsa"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"math/big"
)

type WalletType int

const (
	SoftWallet WalletType = iota

	LedgerWallet

	TrezorWallet
)

//wallet status
const (
	Opened = "Opened"
	Closed = "Closed"
)

type WalletIdentifier struct {
	WalletType `json:"walletType"`
	//wallet file path
	Path       string `json:"path"`
	WalletName string `json:"walletName"`
}

func (w *WalletIdentifier) String() string  {
	return fmt.Sprintf(
		`
     WalletType:   %d
	 Path:         %s
     WalletName:   %s
    `, w.WalletType, w.Path, w.WalletName)
}

type Account struct {
	Address common.Address
}

type AddressInfoReader interface {
	CurrentBalance(address common.Address) *big.Int
	GetTransactionNonce(addr common.Address) (nonce uint64, err error)
}

type Wallet interface {

	//get wallet identifier include type and file path
	GetWalletIdentifier() (WalletIdentifier, error)

	//get wallet status to judge if is locked
	Status() (string, error)

	//establish wallet according to password,return mnemonic
	Establish(path, name, password, passPhrase string) (string, error)

	//restore wallet from mnemonic
	RestoreWallet(path, name, password, passPhrase, mnemonic string, GetAddressRelatedInfo AddressInfoReader) (err error)

	//open
	Open(path, name, password string) error

	//close
	Close() error

	//padding address nonce
	PaddingAddressNonce(GetAddressRelatedInfo AddressInfoReader) (err error)

	//get address nonce
	GetAddressNonce(address common.Address) (nonce uint64, err error)

	//set address nonce
	SetAddressNonce(address common.Address, nonce uint64) (err error)

	//return the accounts in the wallet
	Accounts() ([]Account, error)

	//check if the account is in the wallet
	Contains(account Account) (bool, error)

	//generate new account according to the derived path
	Derive(path DerivationPath, pin bool) (Account, error)

	//find the used account of base and add to the wallet
	SelfDerive(base DerivationPath) error

	//sign hash
	SignHash(account Account, hash []byte) ([]byte, error)

	//get pk form account
	GetPKFromAddress(account Account) (*ecdsa.PublicKey, error)

	//get sk form address
	GetSKFromAddress(address common.Address) (*ecdsa.PrivateKey, error)

	//sign transaction
	SignTx(account Account, tx *model.Transaction, chainID *big.Int) (*model.Transaction, error)

	//generate vrf proof
	Evaluate(account Account, seed []byte) (index [32]byte, proof []byte, err error)
}
