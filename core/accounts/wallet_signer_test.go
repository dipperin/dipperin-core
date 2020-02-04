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
	"path/filepath"
	"testing"
)

var TestSeed = []byte{0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02,
	0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02}
var TestAddress = common.Address{0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12,
	0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12}
var Path = filepath.Join(util.HomeDir(), "/tmp/testSoftWallet1")

//签名hash值
var TestHashData = [32]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}

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

	testCases := []struct{
		name string
		given func() (*WalletSigner, common.Address)
		expect error
	} {
		{
			name:"ErrNotFindWallet",
			given: func() (*WalletSigner, common.Address) {
				ctrl, wallet, _ , walletSigner := getWalletSigner(t)
				defer ctrl.Finish()
				wallet.EXPECT().Evaluate(accountsbase.Account{TestAddress},TestSeed).Return([32]byte{},[]byte{},nil).AnyTimes()
				wallet.EXPECT().Contains(accountsbase.Account{Address:TestAddress}).Return(true, gerror.ErrNotFindWallet).AnyTimes()
				return walletSigner,TestAddress
			},
			expect:gerror.ErrNotFindWallet,
		},
		{
			name:"EvaluateRight",
			given: func() (*WalletSigner, common.Address){
				ctrl, wallet, _ , walletSigner := getWalletSigner(t)
				defer ctrl.Finish()
				wallet.EXPECT().Evaluate(accountsbase.Account{TestAddress},TestSeed).Return([32]byte{},[]byte{},nil).AnyTimes()
				wallet.EXPECT().Contains(accountsbase.Account{Address:TestAddress}).Return(true, nil).AnyTimes()
				return walletSigner,walletSigner.GetAddress()
			},
			expect:nil,
		},
	}


	for _,tc := range testCases{
		walletSigner, addr := tc.given()
		_,_, err := walletSigner.Evaluate(accountsbase.Account{Address:addr},TestSeed)
		assert.Equal(t, tc.expect, err)
	}
}

func TestWalletSigner_GetAddress(t *testing.T) {
	ctrl, _, _ , walletSigner := getWalletSigner(t)
	ctrl.Finish()

	testCases := []struct{
		name string
		given func() *WalletSigner
		expect common.Address
	}{
		{
			name:"GetAddressRight",
			given: func() *WalletSigner {
				return walletSigner
			},
			expect:TestAddress,
		},
		{
			name:"WalletSignerIsNil",
			given: func() *WalletSigner {
				return &WalletSigner{}
			},
			expect:common.Address{},
		},
	}

	for _, tc := range testCases{
		ws := tc.given()
		addr := ws.GetAddress()
		assert.Equal(t, tc.expect, addr)
	}
}


// todo
func TestWalletSigner_PublicKey(t *testing.T) {
	ctrl, wallet, _ , walletSigner := getWalletSigner(t)
	ctrl.Finish()

	wallet.EXPECT().GetPKFromAddress(accountsbase.Account{TestAddress}).Return(nil, nil)
	wallet.EXPECT().Contains(accountsbase.Account{TestAddress}).Return(true, nil).AnyTimes()

	_ = walletSigner.PublicKey()
	//addr := cs_crypto.GetNormalAddress(*pk)
	//assert.Equal(t, walletSigner.GetAddress(), addr)
}

func TestWalletSigner_SetBaseAddress(t *testing.T) {
	ctrl, _, _ , walletSigner := getWalletSigner(t)
	ctrl.Finish()

	walletSigner.SetBaseAddress(TestAddress)
	assert.Equal(t, TestAddress, walletSigner.GetAddress())
}

func TestWalletSigner_SignHash(t *testing.T) {
	ctrl, wallet, _ , walletSigner := getWalletSigner(t)
	ctrl.Finish()
	wallet.EXPECT().Contains(accountsbase.Account{TestAddress}).Return(true, nil)
	wallet.EXPECT().SignHash(accountsbase.Account{TestAddress}, TestHashData[:]).Return([]byte{}, nil)

	_, err := walletSigner.SignHash(TestHashData[:])
	assert.NoError(t, err)
}

/*func TestWalletSigner_ValidSign(t *testing.T) {
	ctrl, _, _ , walletSigner := getWalletSigner(t)
	ctrl.Finish()

	//test travis
	//assert.Error(t,err)

	signature, err := walletSigner.SignHash(wallet.TestHashData[:])
	assert.NoError(t, err)

	pk := crypto.FromECDSAPub(walletSigner.PublicKey())

	err = walletSigner.ValidSign(wallet.TestHashData[:], pk, []byte{})
	assert.Equal(t, gerror.ErrEmptySign, err)

	err = walletSigner.ValidSign(wallet.TestHashData[:], pk, signature[:10])
	assert.Equal(t, gerror.ErrSignatureInvalid, err)

	err = walletSigner.ValidSign(wallet.TestHashData[:], pk, signature)
	assert.NoError(t, err)
	os.Remove(wallet.Path)
}*/
