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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	crypto2 "github.com/dipperin/dipperin-core/third_party/crypto"
	"sync"
)

func MakeWalletSigner(addr common.Address, wm *WalletManager) *WalletSigner {
	return &WalletSigner{
		account:       accountsbase.Account{Address: addr},
		walletManager: wm,
	}
}

type WalletSigner struct {
	account       accountsbase.Account
	walletManager *WalletManager
	lock          sync.Mutex
}

func (signer *WalletSigner) GetAddress() common.Address {
	if signer == nil {
		return common.Address{}
	}
	return signer.account.Address
}

func (signer *WalletSigner) SetBaseAddress(address common.Address) {
	signer.lock.Lock()
	defer signer.lock.Unlock()
	signer.account.Address = address
}

func (signer *WalletSigner) SignHash(hash []byte) ([]byte, error) {
	//log.DLogger.Info("the signer is:","signer",signer)
	wallet, err := signer.walletManager.FindWalletFromAddress(signer.account.Address)
	if err != nil {
		return nil, err
	}
	return wallet.SignHash(signer.account, hash)
}

func (signer *WalletSigner) PublicKey() *ecdsa.PublicKey {
	wallet, err := signer.walletManager.FindWalletFromAddress(signer.account.Address)
	if err != nil {
		return nil
	}
	pk, err := wallet.GetPKFromAddress(signer.account)
	if err != nil {
		return nil
	}
	return pk
}

func (signer *WalletSigner) ValidSign(hash []byte, pubKey []byte, sign []byte) error {
	if len(sign) == 0 {
		return accountsbase.ErrEmptySign
	}
	if crypto2.VerifySignature(pubKey, hash, sign[:len(sign)-1]) == true {
		return nil
	} else {
		return accountsbase.ErrSignatureInvalid
	}
}

func (signer *WalletSigner) Evaluate(account accountsbase.Account, seed []byte) (index [32]byte, proof []byte, err error) {
	//find wallet from address
	wallet, err := signer.walletManager.FindWalletFromAddress(account.Address)
	if err != nil {
		return [32]byte{}, []byte{}, err
	}
	index, proof, err = wallet.Evaluate(account, seed)
	return index, proof, nil
}
