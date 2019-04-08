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

package accounts_test

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/tests/wallet"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMakeWalletSigner(t *testing.T) {
	testWallet, testWalletManager, err := wallet.GetTestWalletManager()
	assert.NoError(t, err)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)
	testSigner := accounts.MakeWalletSigner(testAccounts[0].Address, testWalletManager)
	assert.NotEqual(t, &accounts.WalletSigner{}, testSigner)
}

func TestWalletSigner_Evaluate(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)
	_, _, err = testSigner.Evaluate(accounts.Account{Address: testSigner.GetAddress()}, wallet.TestSeed)
	assert.NoError(t, err)

	_, _, err = testSigner.Evaluate(accounts.Account{Address: wallet.TestAddress}, wallet.TestSeed)
	assert.Equal(t,accounts.ErrNotFindWallet, err)

	os.Remove(wallet.Path)
}

func TestWalletSigner_GetAddress(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)
	addr := testSigner.GetAddress()
	assert.NotEqual(t,common.Address{},addr)

	testSigner=&accounts.WalletSigner{}
	addr = testSigner.GetAddress()
	assert.Equal(t,common.Address{},addr)

	os.Remove(wallet.Path)
}

func TestWalletSigner_PublicKey(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	pk:=testSigner.PublicKey()
	addr := cs_crypto.GetNormalAddress(*pk)
	assert.Equal(t,testSigner.GetAddress(),addr)
	os.Remove(wallet.Path)
}

func TestWalletSigner_SetBaseAddress(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	testSigner.SetBaseAddress(wallet.TestAddress)
	assert.Equal(t,wallet.TestAddress,testSigner.GetAddress())
	os.Remove(wallet.Path)
}

func TestWalletSigner_SignHash(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	_,err=testSigner.SignHash(wallet.TestHashData[:])
	assert.NoError(t,err)
	os.Remove(wallet.Path)
}

func TestWalletSigner_ValidSign(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	//test travis
	//assert.Error(t,err)

	signature,err:=testSigner.SignHash(wallet.TestHashData[:])
	assert.NoError(t,err)

	pk := crypto.FromECDSAPub(testSigner.PublicKey())

	err = testSigner.ValidSign(wallet.TestHashData[:],pk,[]byte{})
	assert.Equal(t,accounts.ErrEmptySign,err)

	err = testSigner.ValidSign(wallet.TestHashData[:],pk,signature[:10])
	assert.Equal(t,accounts.ErrSignatureInvalid,err)

	err = testSigner.ValidSign(wallet.TestHashData[:],pk,signature)
	assert.NoError(t,err)
	os.Remove(wallet.Path)
}
