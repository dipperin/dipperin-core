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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var TestSeed = []byte{0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02,
	0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02}
var TestAddress = common.Address{0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12,
	0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12}
var Path = filepath.Join(util.HomeDir(), "/tmp/testSoftWallet1")

func getWalletSigner(t *testing.T)(*gomock.Controller, *accountsbase.MockWallet, *WalletManager, *WalletSigner)  {
	ctrl, wallet, walletManager := getWalletAndWalletManager(t)

	wallet.EXPECT().Accounts().Return([]accountsbase.Account{accountsbase.Account{Address:TestAddress},},nil)
	testAccounts, err := wallet.Accounts()
	assert.NoError(t, err)
	testSigner := MakeWalletSigner(testAccounts[0].Address, walletManager)
	return ctrl,wallet,walletManager,testSigner
}

func TestMakeWalletSigner(t *testing.T) {
	ctrl, _, _ , walletSigner := getWalletSigner(t)
	ctrl.Finish()

	assert.NotEqual(t, &WalletSigner{}, walletSigner)
}


func TestWalletSigner_Evaluate(t *testing.T) {
	ctrl, wallet, _ , walletSigner := getWalletSigner(t)
	defer ctrl.Finish()
	wallet.EXPECT().Contains(accountsbase.Account{Address:TestAddress}).Return(true, nil)
	wallet.EXPECT().Contains(accountsbase.Account{Address:TestAddress}).Return(true, gerror.ErrNotFindWallet)
	wallet.EXPECT().Evaluate(accountsbase.Account{TestAddress},TestSeed).Return([32]byte{},[]byte{},nil)

	_, _, err := walletSigner.Evaluate(accountsbase.Account{Address: walletSigner.GetAddress()}, TestSeed)
	assert.NoError(t, err)

	_, _, err = walletSigner.Evaluate(accountsbase.Account{Address: TestAddress}, TestSeed)
	assert.Equal(t, gerror.ErrNotFindWallet, err)

	os.Remove(Path)
}
/*
func TestWalletSigner_GetAddress(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)
	addr := testSigner.GetAddress()
	assert.NotEqual(t, common.Address{}, addr)

	testSigner = &accounts.WalletSigner{}
	addr = testSigner.GetAddress()
	assert.Equal(t, common.Address{}, addr)

	os.Remove(wallet.Path)
}

func TestWalletSigner_PublicKey(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	pk := testSigner.PublicKey()
	addr := cs_crypto.GetNormalAddress(*pk)
	assert.Equal(t, testSigner.GetAddress(), addr)
	os.Remove(wallet.Path)
}

func TestWalletSigner_SetBaseAddress(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	testSigner.SetBaseAddress(wallet.TestAddress)
	assert.Equal(t, wallet.TestAddress, testSigner.GetAddress())
	os.Remove(wallet.Path)
}

func TestWalletSigner_SignHash(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	_, err = testSigner.SignHash(wallet.TestHashData[:])
	assert.NoError(t, err)
	os.Remove(wallet.Path)
}

func TestWalletSigner_ValidSign(t *testing.T) {
	testSigner, err := wallet.GetTestWalletSigner()
	assert.NoError(t, err)

	//test travis
	//assert.Error(t,err)

	signature, err := testSigner.SignHash(wallet.TestHashData[:])
	assert.NoError(t, err)

	pk := crypto.FromECDSAPub(testSigner.PublicKey())

	err = testSigner.ValidSign(wallet.TestHashData[:], pk, []byte{})
	assert.Equal(t, accounts.ErrEmptySign, err)

	err = testSigner.ValidSign(wallet.TestHashData[:], pk, signature[:10])
	assert.Equal(t, accounts.ErrSignatureInvalid, err)

	err = testSigner.ValidSign(wallet.TestHashData[:], pk, signature)
	assert.NoError(t, err)
	os.Remove(wallet.Path)
}*/
