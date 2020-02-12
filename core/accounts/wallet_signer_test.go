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
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

var TestSeed = []byte{0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02,
	0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02}
var TestAddress = common.Address{0,0,242,20,126,130,14,11,206,229,189,115,83,59,176,241,164,201,156,211,232,164,}
var Path = filepath.Join(util.HomeDir(), "/tmp/testSoftWallet1")
var sign = []byte{
	17,21,169,122,83,118,144,135,206,154,17,118,109,93,227,215,151,170,7,116,68,134,239,239,102,209,137,157,96,211,212,78,17,174,183,108,154,76,150,70,134,88,31,217,118,154,107,145,22,210,206,193,112,115,8,133,153,227,164,249,63,27,44,40,1,
}





//签名hash值
var TestHashData = [32]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}

func getWalletSigner(t *testing.T)(*gomock.Controller, *MockWallet, *WalletManager, *WalletSigner)  {
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

func TestWalletSigner_ValidSign(t *testing.T) {
	ctrl, wallet, _ , walletSigner := getWalletSigner(t)
	ctrl.Finish()

	var pk = []byte{
		4,247,158,61,109,132,5,231,184,173,19,231,66,132,176,212,161,94,87,186,194,111,134,150,214,35,173,96,61,192,230,207,37,184,111,138,162,5,161,129,52,52,154,135,74,134,72,128,249,43,84,99,97,131,216,232,247,212,17,125,167,17,112,187,50,
	}

	wallet.EXPECT().Contains(accountsbase.Account{TestAddress}).Return(true, nil).AnyTimes()
	wallet.EXPECT().SignHash(accountsbase.Account{TestAddress}, TestHashData[:]).Return(sign, nil)
	pkTemp, err  := crypto.UnmarshalPubkey(pk)
	assert.NoError(t, err)
	wallet.EXPECT().GetPKFromAddress(accountsbase.Account{TestAddress}).Return(pkTemp,nil)
	signature, err := walletSigner.SignHash(TestHashData[:])
	assert.NoError(t, err)


	testCases := []struct{
		name string
		given []byte
		expect error
	}{
		{
			name:"ErrEmptySign",
			given:[]byte{},
			expect:gerror.ErrEmptySign,
		},
		{
			name:"ErrSignatureInvalid",
			given:signature[:10],
			expect:gerror.ErrSignatureInvalid,
		},
		{
			name:"ErrEmptySign",
			given:signature,
			expect:nil,
		},
	}

	for _,tc := range testCases{
		input := tc.given
		err := walletSigner.ValidSign(TestHashData[:],pk, input)
		assert.Equal(t, tc.expect, err)
	}
}
