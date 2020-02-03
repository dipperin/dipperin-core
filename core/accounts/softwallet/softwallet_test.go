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

package softwallet

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/go-bip39"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"math/big"
	"os"
	"testing"
)

//test data
var errSeed = []byte{0x01, 0x02, 0x01}
var testSeed = []byte{0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02,
	0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02, 0x01, 0x02}
var testAddress = common.Address{0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12,
	0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12, 0x12}

//signed hash value
var testHashData = [32]byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}

var walletName = "testSoftWallet"
var path = util.HomeDir() + "/testSoftWallet"
var password = "12345678"
var passPhrase = "12345678"

var errAccount = accountsbase.Account{
	Address: testAddress,
}
var testIdentifier = accountsbase.WalletIdentifier{
	WalletType: accountsbase.SoftWallet,
	Path:       path,
	WalletName: walletName,
}


func TestSoftWallet_paddingWalletInfo(t *testing.T) {

	testMnemonic := "fragile lift oak super joy dust erode give female indoor pass throw dirt gather wedding pyramid box umbrella chimney air middle civil essence dentist"
	errTestMnemonic := "fragile lift oak super joy dust erode give female indoor pass throw dirt gather wedding pyramid box umbrella chimney air middle civil essence"

	testKdfPara := &KDFParameter{
		KDFParams: map[string]interface{}{
			"n":       WalletLightScryptN,
			"p":       WalletLightScryptP,
			"kdfType": KDF,
			"r":       WalletscryptR,
			"keyLen":  WalletscryptDKLen,
		},
	}

	testCases := []struct{
		name string
		given func() (string, string, string, *KDFParameter)
		expect error
	}{
		{
			name:"paddingWalletInfoRight",
			given:func() (string, string, string, *KDFParameter){
				return testMnemonic, password, passPhrase, testKdfPara
			},
			expect:nil,
		},
		{
			name:"errTestMnemonic",
			given:func() (string, string, string, *KDFParameter){
				return errTestMnemonic, password, passPhrase, testKdfPara
			},
			expect:bip39.ErrInvalidMnemonic,
		},
		{
			name:"errKDFParamsType",
			given:func() (string, string, string, *KDFParameter){
				testKdfPara.KDFParams["kdfType"] = "bcrypt"
				return testMnemonic, password, passPhrase, testKdfPara
			},
			expect:gerror.ErrNotSupported,
		},
	}


	for _, tc := range testCases {
		t.Log(tc.name)
		testMnemonic, password, passPhrase, testKdfPara := tc.given()
		testWallet, err := NewSoftWallet()
		assert.NoError(t, err)

		err = testWallet.paddingWalletInfo(testMnemonic, password, passPhrase, testKdfPara)
		assert.Equal(t, tc.expect, err)
	}

}

func establishSoftWallet(path, walletName, password, passPhrase string) (string, *SoftWallet, error) {
	testWallet, err := NewSoftWallet()
	if err != nil {
		return "", nil, err
	}

	testWallet.Identifier.WalletName = walletName
	testWallet.Identifier.Path = path

	os.Remove(path)

	mnemonic, err := testWallet.Establish(path, walletName, password, passPhrase)
	if err != nil {
		return "", nil, err
	}

	//mnemonic = strings.Replace(mnemonic, " ", ",", -1)
	log.DLogger.Info("EstablishWallet mnemonic is:", zap.String("mnemonic", mnemonic))
	return mnemonic, testWallet, nil
}

func GetTestWallet() (*SoftWallet, error) {
	_, softWallet, err := establishSoftWallet(path, walletName, password, passPhrase)
	if err != nil {
		return nil, err
	}

	return softWallet, nil
}

func TestSoftWallet_Establish(t *testing.T) {

	testCases := []struct{
		name string
		given func() (string,string,string,string)
		expect error
	}{
		{
			name:"ErrWalletPath",
			given: func() (string,string,string,string) {
				TestErrWalletPath := "/tmp/testSoftWallet"
				return TestErrWalletPath, walletName, password, passPhrase
			},
			expect:gerror.ErrWalletPathError,
		},
		{
			name:"establishSoftWalletRight",
			given: func() (string,string,string,string) {
				return path, walletName, password, passPhrase
			},
			expect:nil,
		},
		{
			name:"ErrPasswordOrPassPhraseIllegal",
			given: func() (string,string,string,string) {
				return path, walletName, "", passPhrase
			},
			expect:gerror.ErrPasswordOrPassPhraseIllegal,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		walletPath, walletName, password, passPhrase := tc.given()
		_, _, err := establishSoftWallet(walletPath, walletName, password, passPhrase)
		assert.Equal(t, tc.expect, err)
	}

}

func TestSoftWallet_Open(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	testCases := []struct{
		name string
		given func() (string,string,string)
		expect error
	}{
		{
			name:"ErrWalletPathError",
			given: func() (string,string,string) {
				return "/tmp", walletName, password
			},
			expect:gerror.ErrWalletPathError,
		},
		{
			name:"ErrPasswordOrPassPhraseIllegal",
			given: func() (string,string,string) {
				return path, walletName, ""
			},
			expect:gerror.ErrPasswordOrPassPhraseIllegal,
		},
		{
			name:"OpenWalletRight",
			given: func() (string,string,string) {
				return path, walletName, password
			},
			expect:nil,
		},
		{
			name:"ErrWalletPasswordNotValid",
			given: func() (string,string,string) {
				errPassword := "343543564"
				//return  testWallet.Open(path, walletName, errPassword)
				return  path, walletName, errPassword
			},
			expect:gerror.ErrWalletPasswordNotValid,
		},

	}

	for _, tc := range testCases {
		t.Log(tc.name)
		path, walletName, password := tc.given()
		err := testWallet.Open(path, walletName, password)
		assert.Equal(t, tc.expect, err)
	}
	os.Remove(path)

}


func TestSoftWallet_Accounts(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	accounts, err := testWallet.Accounts()
	assert.NoError(t, err)

	log.DLogger.Debug("the wallet accounts is: ", zap.Any("accounts", accounts))

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_SignHash(t *testing.T) {

	testCases := []struct{
		name string
		given func() (*SoftWallet, accountsbase.Account, []byte)
		expect error
	}{
		{
			name:"ErrWalletNotOpen",
			given: func() (*SoftWallet, accountsbase.Account, []byte) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				testWallet.Close()
				return testWallet,testWallet.walletInfo.Accounts[0], testHashData[:]
			},
			expect:gerror.ErrWalletNotOpen,
		},
		{
			name:"ErrInvalidAddress",
			given: func() (*SoftWallet, accountsbase.Account, []byte) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				err = testWallet.Open(path, walletName, password)
				assert.NoError(t, err)
				return testWallet,errAccount, testHashData[:]
			},
			expect:gerror.ErrInvalidAddress,
		},
		{
			name:"SignHashRight",
			given: func() (*SoftWallet, accountsbase.Account, []byte) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				err = testWallet.Open(path, walletName, password)
				assert.NoError(t, err)
				return testWallet,testWallet.walletInfo.Accounts[0], testHashData[:]
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet,account, hash := tc.given()
		_, err := wallet.SignHash(account, hash)
		assert.Equal(t, tc.expect, err)
	}

	os.Remove(path)
}

func TestSoftWallet_SignTx(t *testing.T) {
	//signed transaction information
	testTx := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000), model.TestGasPrice, model.TestGasLimit, []byte{})

	testCases := []struct{
		name string
		given func() (*SoftWallet, accountsbase.Account, *model.Transaction, *big.Int)
		expect error
	}{
		{
			name:"ErrWalletNotOpen",
			given: func() (*SoftWallet, accountsbase.Account, *model.Transaction, *big.Int) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				testWallet.Close()
				return testWallet, testWallet.walletInfo.Accounts[0], testTx, nil
			},
			expect:gerror.ErrWalletNotOpen,
		},
		{
			name:"ErrInvalidAddress",
			given: func() (*SoftWallet, accountsbase.Account, *model.Transaction, *big.Int) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				err = testWallet.Open(path, walletName, password)
				assert.NoError(t, err)
				return testWallet,errAccount, testTx, nil
			},
			expect:gerror.ErrInvalidAddress,
		},
		{
			name:"SignHashRight",
			given: func() (*SoftWallet, accountsbase.Account, *model.Transaction, *big.Int) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				err = testWallet.Open(path, walletName, password)
				assert.NoError(t, err)
				return testWallet,testWallet.walletInfo.Accounts[0], testTx, nil
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet, account, tx, chainID := tc.given()
		_, err := wallet.SignTx(account, tx, chainID)
		assert.Equal(t, tc.expect, err)
	}

	os.Remove(path)
}

func TestSoftWallet_Evaluate(t *testing.T) {
	testCases := []struct{
		name string
		given func() (*SoftWallet, accountsbase.Account, []byte)
		expect error
	}{
		{
			name:"ErrWalletNotOpen",
			given: func() (*SoftWallet, accountsbase.Account, []byte) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				testWallet.Close()
				return testWallet,testWallet.walletInfo.Accounts[0], testSeed
			},
			expect:gerror.ErrWalletNotOpen,
		},
		{
			name:"ErrInvalidAddress",
			given: func() (*SoftWallet, accountsbase.Account, []byte) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				err = testWallet.Open(path, walletName, password)
				assert.NoError(t, err)
				return testWallet,errAccount, testSeed
			},
			expect:gerror.ErrInvalidAddress,
		},
		{
			name:"SignHashRight",
			given: func() (*SoftWallet, accountsbase.Account, []byte) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				err = testWallet.Open(path, walletName, password)
				assert.NoError(t, err)
				return testWallet,testWallet.walletInfo.Accounts[0], testSeed
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet, account, seed := tc.given()
		_,_, err := wallet.Evaluate(account, seed)
		assert.Equal(t, tc.expect, err)
	}

	os.Remove(path)
}


func TestSoftWallet_Derive(t *testing.T) {
	testCases := []struct{
		name string
		given func() (*SoftWallet, accountsbase.DerivationPath, bool)
		expect error
	}{
		{
			name:"ErrInvalidDerivedPath",
			given: func() (*SoftWallet, accountsbase.DerivationPath, bool) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)

				testWallet.Open(path, walletName, password)
				assert.NoError(t, err)

				var testDrivePath = accountsbase.DerivationPath{0, 1, 1, 0}
				return testWallet,testDrivePath, true
			},
			expect:gerror.ErrInvalidDerivedPath,
		},
		{
			name:"DeriveRight",
			given: func() (*SoftWallet, accountsbase.DerivationPath, bool) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)

				testWallet.Open(path, walletName, password)
				assert.NoError(t, err)

				testDrivePath := accountsbase.DerivationPath{0x80000000 + 44, 0x80000000 + 709394, 0x80000000 + 0, 0}
				return testWallet,testDrivePath,true
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet, drivePath, save := tc.given()
		derivedAccount, err := wallet.Derive(drivePath, save)
		if err != nil {
			assert.Equal(t, tc.expect, err)
		}else {
			contain, err := wallet.Contains(derivedAccount)
			assert.NoError(t, err)
			assert.Equal(t, true, contain)
		}
	}

	os.Remove(path)
}

func TestSoftWallet_RestoreWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	accountStatus := accountsbase.NewMockAddressInfoReader(ctrl)
	mnemonic := "chicken coconut winner february brown topple pond bird endless salt filter journey mass ramp milk tuition card seat worth school length rain slice ozone"
	walletName := "RemainRewardWallet4"


	testCases := []struct{
		name string
		given func() (*accountsbase.MockAddressInfoReader, string, string)
		expect error
	}{
		{
			name:"ErrWalletPathError",
			given: func() (*accountsbase.MockAddressInfoReader, string,string){
				return accountStatus,"/tmp", password
			},
			expect:gerror.ErrWalletPathError,
		},
		{
			name:"ErrPasswordOrPassPhraseIllegal",
			given: func() (*accountsbase.MockAddressInfoReader, string,string) {
				return accountStatus, path, ""
			},
			expect:gerror.ErrPasswordOrPassPhraseIllegal,
		},
		{
			name:"RestoreWalletRight",
			given: func() (*accountsbase.MockAddressInfoReader, string,string) {
				path := util.HomeDir() + "/testSoftWallet/RemainRewardWallet4"
				password = "12345678"
				accountStatus.EXPECT().GetTransactionNonce(common.HexToAddress("0x00001ac2a396f7100C4B2838A171B68d654B9B56B0c1")).Return(uint64(0), nil).AnyTimes()
				accountStatus.EXPECT().GetTransactionNonce(common.HexToAddress("0x0000F2415f9280d7Cb6eF45046a189dA46FBacb7dF5D")).Return(uint64(0), gerror.ErrAccountNotExist).AnyTimes()
				accountStatus.EXPECT().CurrentBalance(common.HexToAddress("0x00001ac2a396f7100C4B2838A171B68d654B9B56B0c1")).Return(big.NewInt(int64(100))).AnyTimes()
				accountStatus.EXPECT().CurrentBalance(common.HexToAddress("0x0000F2415f9280d7Cb6eF45046a189dA46FBacb7dF5D")).Return(big.NewInt(int64(100))).AnyTimes()
				return accountStatus, path, password
			},
			expect:nil,
		},
	}
	for _, tc := range testCases {
		testWallet, err := NewSoftWallet()
		assert.NoError(t, err)
		t.Log(tc.name)
		accountStatusTmp, tempPath, passwd := tc.given()
		err = testWallet.RestoreWallet(tempPath, walletName,  passwd, passPhrase, mnemonic, accountStatusTmp)
		assert.Equal(t, tc.expect, err)
		os.Remove(tempPath)
	}
}

func TestSoftWallet_Status(t *testing.T) {
	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)

	status, err := testWallet.Status()
	assert.NoError(t, err)
	assert.Equal(t, accountsbase.Closed, status)

}

func TestSoftWallet_GetWalletIdentifier(t *testing.T) {
	testCases := []struct{
		name string
		given func() *SoftWallet
		expect error
	}{
		{
			name:"ErrWalletNotOpen",
			given: func() *SoftWallet {
				testWallet, err := NewSoftWallet()
				assert.NoError(t, err)
				return testWallet
			},
			expect:gerror.ErrWalletNotOpen,
		},
		{
			name:"GetWalletIdentifierRight",
			given: func() *SoftWallet {
				testWallet, err := NewSoftWallet()
				assert.NoError(t, err)
				testWallet.Identifier = testIdentifier
				testWallet.status = accountsbase.Opened
				return testWallet
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet := tc.given()
		id ,err := wallet.GetWalletIdentifier()
		if err != nil {
			assert.Equal(t, tc.expect, err)
		}else {
			assert.Equal(t, wallet.Identifier, id)
		}
		os.Remove(path)
	}
}

func TestSoftWallet_PaddingAddressNonce(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountStatus := accountsbase.NewMockAddressInfoReader(ctrl)
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	accounts, err := testWallet.Accounts()
	assert.NoError(t, err)
	for _, ac := range accounts{
		accountStatus.EXPECT().GetTransactionNonce(ac.Address).Return(uint64(0),nil)
	}

	err = testWallet.PaddingAddressNonce(accountStatus)
	assert.NoError(t, err)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_SetAddressNonce(t *testing.T) {
	testCases := []struct{
		name string
		given func() (*SoftWallet, common.Address, uint64)
		expect error
	}{
		{
			name:"SetAddressNonceRight",
			given: func() (*SoftWallet, common.Address, uint64) {
				testWallet, err := GetTestWallet()
				assert.NoError(t, err)

				testAccounts, err := testWallet.Accounts()
				assert.NoError(t, err)
				return testWallet, testAccounts[0].Address, 1
			},
			expect:nil,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet, address, nonce := tc.given()
		err := wallet.SetAddressNonce(address, nonce)
		assert.Equal(t, tc.expect, err)
		nonce, err = wallet.GetAddressNonce(address)
		assert.Equal(t, uint64(1), nonce)
		assert.NoError(t, err)
	}

	os.Remove(path)
}

func TestSoftWallet_GetPKFromAddress(t *testing.T) {
	testCases := []struct{
		name string
		given func() (*SoftWallet, accountsbase.Account)
		expect error
	}{
		{
			name:"ErrInvalidAddress",
			given: func() (*SoftWallet, accountsbase.Account) {

				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				return testWallet, errAccount
			},
			expect:gerror.ErrInvalidAddress,
		},
		{
			name:"GetPKFromAddressRight",
			given: func() (*SoftWallet, accountsbase.Account) {

				testWallet, err := GetTestWallet()
				assert.NoError(t, err)

				testAccounts, err := testWallet.Accounts()
				assert.NoError(t, err)
				return testWallet,testAccounts[0]
			},
			expect:nil,
		},
		{
			name:"ErrWalletNotOpen",
			given: func() (*SoftWallet, accountsbase.Account) {

				testWallet, err := GetTestWallet()
				assert.NoError(t, err)
				testAccounts, err := testWallet.Accounts()
				assert.NoError(t, err)

				testWallet.Close()
				return testWallet,testAccounts[0]
			},
			expect:gerror.ErrWalletNotOpen,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		wallet, account := tc.given()
		_, err := wallet.GetPKFromAddress(account)
		assert.Equal(t, tc.expect, err)
		os.Remove(path)
	}
}

func TestSoftWallet_GetSKFromAddress(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)

	testCases := []struct{
		name string
		given func() common.Address
		expect error
	}{
		{
			name:"ErrInvalidAddress",
			given: func() common.Address {

				return errAccount.Address
			},
			expect:gerror.ErrInvalidAddress,
		},
		{
			name:"GetSKFromAddressRight",
			given: func() common.Address {
				return testAccounts[0].Address
			},
			expect:nil,
		},
		{
			name:"ErrWalletNotOpen",
			given: func() common.Address {
				testWallet.Close()
				return testAccounts[0].Address
			},
			expect:gerror.ErrWalletNotOpen,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		_, err := testWallet.GetSKFromAddress(tc.given())
		assert.Equal(t, tc.expect, err)
		os.Remove(path)
	}
}

func TestSoftWallet_Close(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	err = testWallet.Close()
	assert.NoError(t, err)

	status, err := testWallet.Status()
	assert.NoError(t, err)
	assert.Equal(t, accountsbase.Closed, status)

	os.Remove(path)
}

func TestSoftWallet_Contains(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)


	testCases := []struct{
		name string
		given func() accountsbase.Account
		expect error
		contains bool
	}{
		{
			name:"ErrInvalidAddress",
			given: func() accountsbase.Account {
				return errAccount
			},
			expect:nil,
			contains:false,
		},
		{
			name:"ContainsRight",
			given: func() accountsbase.Account {
				return testAccounts[0]
			},
			expect:nil,
			contains:true,
		},
		{
			name:"ErrWalletNotOpen",
			given: func() accountsbase.Account {
				testWallet.Close()
				return  testAccounts[0]
			},
			expect:gerror.ErrWalletNotOpen,
			contains:false,
		},
	}

	for _, tc := range testCases {
		t.Log(tc.name)
		contains, err := testWallet.Contains(tc.given())
		assert.Equal(t, tc.expect,err )
		assert.Equal(t, tc.contains, contains)
	}
}





