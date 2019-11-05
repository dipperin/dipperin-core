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

package soft_wallet

import (
	"encoding/json"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	user2 "os/user"
	"strconv"
	"sync/atomic"
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
var password = "123456"
var passPhrase = ""
var TestErrWalletPath = "/tmp/testSoftWallet"
var errAccount = accounts.Account{
	Address: testAddress,
}
var testIdentifier = accounts.WalletIdentifier{
	WalletType: accounts.SoftWallet,
	Path:       path,
	WalletName: walletName,
}

type testAccountStatus struct {
}

func (*testAccountStatus) CurrentBalance(address common.Address) *big.Int {
	return big.NewInt(0)
}

func (*testAccountStatus) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	return 0, nil
}

var tmpAccountStatus testAccountStatus

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

	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)

	err = testWallet.paddingWalletInfo(testMnemonic, "123", "", testKdfPara)
	assert.NoError(t, err)

	err = testWallet.paddingWalletInfo(errTestMnemonic, "123", "", testKdfPara)
	assert.Error(t, err)

	testKdfPara.KDFParams["kdfType"] = "bcrypt"
	err = testWallet.paddingWalletInfo(testMnemonic, "123", "", testKdfPara)
	assert.Error(t, err)

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
	log.Info("EstablishWallet mnemonic is:", "mnemonic", mnemonic)
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
	_, _, err := establishSoftWallet(TestErrWalletPath, walletName, password, passPhrase)
	assert.Error(t, err)

	err = os.RemoveAll(path)
	//assert.NoError(t,err)
	_, _, err = establishSoftWallet(path, walletName, password, passPhrase)
	assert.NoError(t, err)

	_, _, err = establishSoftWallet(path, walletName, "", passPhrase)
	assert.Equal(t, accounts.ErrPasswordIsNil, err)
	os.RemoveAll(path)

	_, _, err = establishSoftWallet(path, walletName, password, passPhrase)
	assert.NoError(t, err)
}

func TestSoftWallet_Open(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	err = testWallet.Open("/tmp", walletName, password)
	assert.Equal(t, accounts.ErrWalletPathError, err)

	err = testWallet.Open(path, walletName, "")
	assert.Equal(t, accounts.ErrPasswordIsNil, err)

	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	err = testWallet.Close()
	assert.NoError(t, err)

	errPassword := "343543564"
	err = testWallet.Open(path, walletName, errPassword)

	assert.Error(t, err, accounts.ErrWalletPasswordNotValid)

	os.Remove(path)
}

func TestEconomyOpenWallet(t *testing.T) {
	t.Skip()
	walletName := "InvestorWallet0"
	path := "/home/qydev/economyWallet/InvestorWallet0"
	password := "123"
	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)

	testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	log.Info("the main account address is:", "address", testWallet.walletInfo.Accounts[0].Address.Hex())
}

func TestGetWalletPrivateKey(t *testing.T) {
	t.Skip()
	path := "/home/qydev/tmp/dipperin_apps/default_v0/CSWallet"
	password := "123"
	walletName := "CSWallet"

	log.Info("the path is:", "path", path)

	cmd2 := exec.Command("pwd")
	output, err := cmd2.Output()
	assert.NoError(t, err)
	log.Info("the output is1:", "output", string(output))

	cmd := exec.Command("ls")
	output, err = cmd.Output()
	assert.NoError(t, err)
	log.Info("the output is:", "output", string(output))

	assert.NoError(t, err)

	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)
	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	log.Info("the account info is: ", "sk", testWallet.walletInfo.ExtendKeys)
}

func TestSoftWallet_Accounts(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	accounts, err := testWallet.Accounts()
	assert.NoError(t, err)

	log.Debug("the wallet accounts is: ", "accounts", accounts)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_SignHash(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	testWallet.Close()
	_, err = testWallet.SignHash(testWallet.walletInfo.Accounts[0], testHashData[:])
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	_, err = testWallet.SignHash(errAccount, testHashData[:])
	assert.Equal(t, accounts.ErrInvalidAddress, err)

	testSignData, err := testWallet.SignHash(testWallet.walletInfo.Accounts[0], testHashData[:])
	assert.NoError(t, err)

	log.Debug("the testSignData is: ", "testSignData", testSignData)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_SignTx(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	//signed transaction information
	testTx := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})

	testWallet.Close()
	_, err = testWallet.SignTx(testWallet.walletInfo.Accounts[0], testTx, nil)
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	_, err = testWallet.SignTx(errAccount, testTx, nil)
	assert.Equal(t, accounts.ErrInvalidAddress, err)

	signTxResult, err := testWallet.SignTx(testWallet.walletInfo.Accounts[0], testTx, nil)
	assert.NoError(t, err)

	log.Debug("TestSoftWallet_SignTx end", "err", err)
	log.Debug("the signTxResult is", "signTxResult", signTxResult.CalTxId().Hex())

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_Evaluate(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	testWallet.Close()
	_, _, err = testWallet.Evaluate(testWallet.walletInfo.Accounts[0], testSeed)
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	_, _, err = testWallet.Evaluate(errAccount, testSeed)
	assert.Equal(t, accounts.ErrInvalidAddress, err)

	_, _, err = testWallet.Evaluate(testWallet.walletInfo.Accounts[0], testSeed)
	assert.NoError(t, err)

	testWallet.Close()
	os.Remove(path)
}

type TestTx struct {
	test1 int
	test2 atomic.Value
}

func Test_TX(t *testing.T) {
	//testTx := TestTx{}

	testTx := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})

	log.Debug("Test_tx is:", "testTx", testTx)
}

func TestSoftWallet_Derive(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	var testDrivePath = accounts.DerivationPath{0, 1, 1, 0}
	derivedAccount, err := testWallet.Derive(testDrivePath, true)
	assert.Error(t, err, accounts.ErrInvalidDerivedPath)

	testDrivePath = accounts.DerivationPath{0x80000000 + 44, 0x80000000 + 709394, 0x80000000 + 0, 0}
	derivedAccount, err = testWallet.Derive(testDrivePath, true)
	assert.NoError(t, err)

	contain, err := testWallet.Contains(derivedAccount)
	assert.NoError(t, err)
	assert.Equal(t, true, contain)

	dpath, _ := accounts.ParseDerivationPath(DefaultDerivedPath)
	//testWallet.Derive(DefaultDerivedPath, false)
	log.Info("TestSoftWallet_Derive", "dpath", dpath)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_RestoreWallet(t *testing.T) {
	mnemonic := "chicken coconut winner february brown topple pond bird endless salt filter journey mass ramp milk tuition card seat worth school length rain slice ozone"

	//create a new local wallet
	testWallet, err := NewSoftWallet()

	assert.NoError(t, err)

	walletName := "RemainRewardWallet4"
	path := util.HomeDir() + "/testSoftWallet/RemainRewardWallet4"
	password := "123"
	passPhrase := ""

	err = testWallet.RestoreWallet("/tmp", walletName, password, passPhrase, mnemonic, &tmpAccountStatus)
	assert.Equal(t, accounts.ErrWalletPathError, err)

	err = testWallet.RestoreWallet(path, walletName, "", passPhrase, mnemonic, &tmpAccountStatus)
	assert.Equal(t, accounts.ErrPasswordIsNil, err)

	err = testWallet.RestoreWallet(path, walletName, password, passPhrase, mnemonic, &tmpAccountStatus)
	assert.NoError(t, err)

	accounts, err := testWallet.Accounts()
	assert.NoError(t, err)
	log.Debug("the wallet accounts is: ", "accounts", accounts)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_Status(t *testing.T) {
	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)

	status, err := testWallet.Status()
	assert.NoError(t, err)
	assert.Equal(t, accounts.Closed, status)

}

func TestSoftWallet_decryptWallet(t *testing.T) {
	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)
	testWallet.decryptWallet("")
}

func TestSoftWallet_GetWalletIdentifier(t *testing.T) {
	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)
	_, err = testWallet.GetWalletIdentifier()
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	testWallet.Identifier = testIdentifier
	testWallet.status = accounts.Opened
	id, err := testWallet.GetWalletIdentifier()
	assert.NoError(t, err)
	assert.Equal(t, testIdentifier, id)
}

func TestSoftWallet_PaddingAddressNonce(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	err = testWallet.PaddingAddressNonce(&testAccountStatus{})
	assert.NoError(t, err)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_SetAddressNonce(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)

	err = testWallet.SetAddressNonce(testAccounts[0].Address, 1)
	assert.NoError(t, err)

	nonce, err := testWallet.GetAddressNonce(testAccounts[0].Address)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), nonce)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_GetPKFromAddress(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	_, err = testWallet.GetPKFromAddress(errAccount)
	assert.Equal(t, accounts.ErrInvalidAddress, err)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)

	_, err = testWallet.GetPKFromAddress(testAccounts[0])
	assert.NoError(t, err)

	testWallet.Close()
	_, err = testWallet.GetPKFromAddress(testAccounts[0])
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	os.Remove(path)
}

func TestSoftWallet_GetSKFromAddress(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	_, err = testWallet.GetSKFromAddress(errAccount.Address)
	assert.Equal(t, accounts.ErrInvalidAddress, err)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)

	_, err = testWallet.GetSKFromAddress(testAccounts[0].Address)
	assert.NoError(t, err)

	testWallet.Close()
	_, err = testWallet.GetSKFromAddress(testAccounts[0].Address)
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	os.Remove(path)
}

func TestSoftWallet_Close(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	err = testWallet.Close()
	assert.NoError(t, err)

	status, err := testWallet.Status()
	assert.NoError(t, err)
	assert.Equal(t, accounts.Closed, status)

	testWallet.Close()
	os.Remove(path)
}

func TestSoftWallet_Contains(t *testing.T) {
	testWallet, err := GetTestWallet()
	assert.NoError(t, err)

	result, err := testWallet.Contains(errAccount)
	assert.NoError(t, err)
	assert.Equal(t, false, result)

	testAccounts, err := testWallet.Accounts()
	assert.NoError(t, err)
	result, err = testWallet.Contains(testAccounts[0])
	assert.NoError(t, err)
	assert.Equal(t, true, result)

	testWallet.Close()
	result, err = testWallet.Contains(testAccounts[0])
	assert.Equal(t, accounts.ErrWalletNotOpen, err)

	os.Remove(path)
}

const generateWallet = false
const generateWalletNumber = 5
const generateWalletPath = "testWallet"
const walletPassword = "123"

type walletConf struct {
	WalletCipher string
	MainAddress  string
}

type walletInfo struct {
	Mnemonic       string `json:"mnemonic"`
	Address        string `json:"mainAddress"`
	WalletPassword string `json:"wallet_password"`
}

//generate test wallet
func TestGenerateWallet(t *testing.T) {
	if !generateWallet {
		return
	}

	user, err := user2.Current()
	assert.NoError(t, err)

	walletPath := user.HomeDir + "/" + generateWalletPath

	log.Info("the walletPath is:", "walletPath", walletPath)
	_, err = os.Stat(walletPath)
	if err == nil {
		err = os.RemoveAll(walletPath)
		assert.NoError(t, err)
	}

	err = os.Mkdir(walletPath, 0777)
	assert.NoError(t, err)
	err = os.Chmod(walletPath, 0777)
	assert.NoError(t, err)
	var conf []*walletConf

	for i := 0; i < generateWalletNumber; i++ {
		walletName := "testSoftWallet" + strconv.Itoa(i)
		path := walletPath + "/" + walletName
		passPhrase := ""

		mnemonic, wallet, err := establishSoftWallet(path, walletName, walletPassword, passPhrase)
		assert.NoError(t, err)

		accounts, err := wallet.Accounts()
		assert.NoError(t, err)
		log.Info("the mine address is:", "address", accounts[0].Address.Hex())

		// Must be at the front, he will post the address information in this file, causing the wallet to be incorrect
		wb, err := ioutil.ReadFile(path)
		assert.NoError(t, err)
		conf = append(conf, &walletConf{
			WalletCipher: string(wb),
			MainAddress:  accounts[0].Address.Hex(),
		})

		walletInfoFile := "walletInfo" + strconv.Itoa(i)
		infoPath := walletPath + "/" + walletInfoFile
		info := walletInfo{
			Mnemonic:       mnemonic,
			Address:        accounts[0].Address.Hex(),
			WalletPassword: walletPassword,
		}
		writeData, err := json.Marshal(&info)
		assert.NoError(t, err)
		err = ioutil.WriteFile(infoPath, writeData, 0666)
		assert.NoError(t, err)

	}

	fmt.Println("===================")
	fmt.Println(util.StringifyJson(conf))
}
