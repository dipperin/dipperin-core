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
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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
var password = "12345678"
var passPhrase = "12345678"
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

	err = testWallet.paddingWalletInfo(testMnemonic, password, passPhrase, testKdfPara)
	assert.NoError(t, err)

	err = testWallet.paddingWalletInfo(errTestMnemonic, password, passPhrase, testKdfPara)
	assert.Error(t, err)

	testKdfPara.KDFParams["kdfType"] = "bcrypt"
	err = testWallet.paddingWalletInfo(testMnemonic, password, passPhrase, testKdfPara)
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
	_, _, err := establishSoftWallet(TestErrWalletPath, walletName, password, passPhrase)
	assert.Error(t, err)

	err = os.RemoveAll(path)
	//assert.NoError(t,err)
	_, _, err = establishSoftWallet(path, walletName, password, passPhrase)
	assert.NoError(t, err)

	_, _, err = establishSoftWallet(path, walletName, "", passPhrase)
	assert.Equal(t, accounts.ErrPasswordOrPassPhraseIllegal, err)
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
	assert.Equal(t, accounts.ErrPasswordOrPassPhraseIllegal, err)

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
	password := "12345678"
	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)

	testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	log.DLogger.Info("the main account address is:", zap.String("address", testWallet.walletInfo.Accounts[0].Address.Hex()))
}

func TestGetWalletPrivateKey(t *testing.T) {
	t.Skip()
	path := "/home/qydev/tmp/dipperin_apps/default_v0/CSWallet"
	password := "12345678"
	walletName := "CSWallet"

	log.DLogger.Info("the path is:", zap.String("path", path))

	cmd2 := exec.Command("pwd")
	output, err := cmd2.Output()
	assert.NoError(t, err)
	log.DLogger.Info("the output is1:", zap.String("output", string(output)))

	cmd := exec.Command("ls")
	output, err = cmd.Output()
	assert.NoError(t, err)
	log.DLogger.Info("the output is:", zap.String("output", string(output)))

	assert.NoError(t, err)

	testWallet, err := NewSoftWallet()
	assert.NoError(t, err)
	err = testWallet.Open(path, walletName, password)
	assert.NoError(t, err)

	log.DLogger.Info("the account info is: ", zap.Any("sk", testWallet.walletInfo.ExtendKeys))
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

	log.DLogger.Debug("the testSignData is: ", zap.Uint8s("testSignData", testSignData))

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

	log.DLogger.Debug("TestSoftWallet_SignTx end", zap.Error(err))
	log.DLogger.Debug("the signTxResult is", zap.String("signTxResult", signTxResult.CalTxId().Hex()))

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

	log.DLogger.Debug("Test_tx is:", zap.Any("testTx", testTx))
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
	log.DLogger.Info("TestSoftWallet_Derive", zap.Any("dpath", dpath))

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
	password := "12345678"
	passPhrase := "12345678"

	err = testWallet.RestoreWallet("/tmp", walletName, password, passPhrase, mnemonic, &tmpAccountStatus)
	assert.Equal(t, accounts.ErrWalletPathError, err)

	err = testWallet.RestoreWallet(path, walletName, "", passPhrase, mnemonic, &tmpAccountStatus)
	assert.Equal(t, accounts.ErrPasswordOrPassPhraseIllegal, err)

	err = testWallet.RestoreWallet(path, walletName, password, passPhrase, mnemonic, &tmpAccountStatus)
	assert.NoError(t, err)

	accounts, err := testWallet.Accounts()
	assert.NoError(t, err)
	log.DLogger.Debug("the wallet accounts is: ", zap.Any("accounts", accounts))

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
const walletPassword = "12345678"

type walletConf struct {
	WalletCipher string
	MainAddress  string
	//PK           string
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

	log.DLogger.Info("the walletPath is:", zap.Any("walletPath", walletPath))
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
		log.DLogger.Info("the mine address is:", zap.String("address", accounts[0].Address.Hex()))

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

//generate test wallet
func TestGenerateWalletForMonitor(t *testing.T) {
	t.Skip()
	user, err := user2.Current()
	assert.NoError(t, err)

	walletPath := user.HomeDir + "/test/" + generateWalletPath

	log.DLogger.Info("the walletPath is:", zap.String("walletPath", walletPath))
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
	generateWalletNumberN := 22

	for i := 0; i < generateWalletNumberN; i++ {
		walletName := "testSoftWallet" + strconv.Itoa(i)
		path := walletPath + "/" + walletName
		passPhrase := "12345678"

		mnemonic, wallet, err := establishSoftWallet(path, walletName, walletPassword, passPhrase)
		assert.NoError(t, err)

		accounts, err := wallet.Accounts()
		//pk,_ := wallet.GetSKFromAddress(accounts[0].Address)
		assert.NoError(t, err)
		log.DLogger.Info("the mine address is:", zap.String("address", accounts[0].Address.Hex()))

		// Must be at the front, he will post the address information in this file, causing the wallet to be incorrect
		wb, err := ioutil.ReadFile(path)
		assert.NoError(t, err)
		conf = append(conf, &walletConf{
			WalletCipher: string(wb),
			MainAddress:  accounts[0].Address.Hex(),
			//PK: pk.
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
		//fmt.Println(info)
		assert.NoError(t, err)

	}

	fmt.Println("===================")
	result := util.StringifyJson(conf)
	//strings.Replace(result, `\"`, `"`, -1)
	fmt.Println(result)
}

func TestRecoverWalletInfo(t *testing.T) {
	//RecoverWalletInfoAndPrint(t, verifierBootNodeDefaultAccounts)
	//RecoverWalletInfoAndPrint(t, verifierDefaultAccounts)
}

func RecoverWalletInfoAndPrint(t *testing.T, walletInfo []walletConf) {
	var conf []*walletConf
	sw, err := NewSoftWallet()
	assert.NoError(t, err)
	for _, wi := range walletInfo {
		err = json.Unmarshal([]byte(wi.WalletCipher), &sw.walletFileInfo)
		assert.NoError(t, err)
		_, WalletPlain, keyData, err := sw.decryptWalletFromJsonData([]byte(wi.WalletCipher), "123")
		assert.NoError(t, err)
		keyData, err = GenSymKeyFromPassword(password, sw.walletFileInfo.KDFParameter)
		assert.NoError(t, err)
		hd := NewHdWalletInfo()
		err = hd.HdWalletInfoDecodeJson(WalletPlain)
		assert.NoError(t, err)
		sw.walletInfo = *hd
		sw.symmetricKey = keyData
		//sw.Identifier.Path =  "/Users/konggan/test/CSWallet" + strconv.Itoa(i)

		sw.walletFileInfo.WalletCipher, err = CalWalletCipher(sw.walletInfo, sw.walletFileInfo.IV[:], sw.symmetricKey)
		assert.NoError(t, err)
		sw.status = accounts.Opened

		writeData, err := json.Marshal(sw.walletFileInfo)
		assert.NoError(t, err)
		accounts, err := sw.Accounts()
		assert.NoError(t, err)
		conf = append(conf, &walletConf{
			WalletCipher: string(writeData),
			MainAddress:  accounts[0].Address.Hex(),
		})
	}

	fmt.Println("===================================================")
	fmt.Println(util.StringifyJson(conf))
}

func TestCryptoForPaymentChannel(t *testing.T) {
	wc := walletConf{
		WalletCipher: `{"Cipher":"990bf546229f1e4e92fab244e63070975e5ec5c28e5b9004247f1dcd7c3cebbb94d1d8f93144a2673aaae59d4b7aac97b7b6a2818d5466803cd05b3ee74f1d62947e21cd08f3b1492e20e1fb2f56c587a973664c593440122d74de0f66d0f48aa9eae9fa0dc8a1ceaf439b909fd899c9222d728ef8c1fa32e086d77c8cf5de7e0bdddf1abf159d57d14d631afd8af44a98cf66aaea2d5230decffeb77a87b2c8e0334aba7fb753b420c408f6f7d2db30fb9dbc1ef64b5c86816d03004bef84319ac8cc982c9de535581758f9ddcbe6a192748ebf4972760cf8b6c5be6ff18a7628c0816ef5b6e32c6e0f28ab56e3a00076aea76900022be16dbdf560a292d097c690f3e67910a7e80156a18eb75705914665dd82544ec9fde231fe929fabf8819b70969b2b5b83ae471f2219fb4c0487cfaeb21943235130bc73c630160720324196711515503d53a8ae93990d3b2538e49b8c1f6b12c6e9422ebe2cceced06d361ec8ceb305ea7663e0de039ffd9260eb8f9e9894106332ec8b49acaf625c738122fd76da0c5a987f7909371316ef0431530c2454d8db5719429f25d12967bf943ced28de279a1e0c5432f3dfb8e0175b854ca6eb63eb28c1d633a0ea6c6435e8508a195e90aaceee58a553a71f4f10e11a5c4242706d637c6108f8dc14e6d1e07b26ca6961ac1f42da636e124c761e2233092610646889ba1c17b8b7023b6a3e0dd637e7a151eff785ddff5fa6312056eefa1d63f10707363bc3f291ff93183d03fafb0816e0448dc3a7dd162d5ed06a994b90679c57f460a6e1272841d9d31eba48010d3608878db31699b58b60040e7dfae2e66b1675d6f2cd5b816263191c8853663aac82292abb69e12f08cf1aa145901e2b8e26acc06f8739b8a65177ad4cede772bf37eb42d51e1e7fed1512ab584e85ffc88af6cfe2c873da95b53aafa56f64d47c951b92effaef407549736036f8799ae14d6914c4cd60e2b9a1a7d54084bb8054150ea1bd7f44e08f1eeb86f870fe2bf034e024e84fede57f56804b5287c25d67a0e6f9d69b091e5761737680ac6a1009aea8c6c4d97dea3c8738d2f23e46c298d637a5730e1832cf6942ac80f49f51dcc3fd9f7ff7972966c874c74241d960d800c03ab2794d89a7544e2322fe15ecbaf318c8f14f0bcb48fe57779b78e088ff5768951ce36d9f8d7139483021f78cf0d05b3bb3a241c5627feb83f59d1620d439e3d246b4e58fed21a9bad960269ea389c38985d5b748c40839f84cf4914287b2756c880040a7b87481d2989bc371a526fd83d5df23cea4135ebc7f5d76acfbf0982e1e3d1dcbcb119fffa76a926fe3313d889e84aa019716d85e1f136687585e0868646802f5c8784b5d7f05d866ec79d3704a6e44e17723d2f63bfa004f837e97deb6d086c0da02451e55786eaeefe5e7162eccd09da0ccc63f3e5df1faf4cb0dc1c8efb988e04587fbdfa5be18b37d9685b281a14e57a8ad88085e23825c7549cb1507d879ae3dd37ec7820d4082b77ae38e9e10a4442a97cce8e311f715d724340807d37bdc526afb5d71788ccce0a601bb34d72cd7c44b4a3a2ab954f7c4d9831833e8dbf4a7194682d37bed5a76e02e8596ba2b5c6d4c093970bba6f38fd79fc03caff3e0c1fda414602d7ce6d599ef6ac4a83fec9c66c4415cc3a6ffb390b3ca91e8ee196b9185100f5fb9d4a24da2dd9a7e8fb58d0be5133816369a8057094a863444bd592d3fa50396e9e2e2b103709d7c0a30cbf2def75cae5dd69d5bc5d31d967b9dd9e3","MacCipher":"82fa9c8097b85dd8cdb6a7eabf343f66c1659851dfa8dff221a3d4cbefa61b4e","AlgType":"AES-256","ModeType":"CBC","IV":[220,80,210,231,228,225,190,88,19,155,160,115,209,141,116,75],"kdf":"","kdfparams":{"kdfType":"Scrypt","keyLen":32,"n":4096,"p":6,"r":8,"salt":"017302eb3a7ee61665b2acca4a1ba281de199ddcb7900fe7d7ac2127ca3282b6"}}`,
		MainAddress:  "0x0000F075aA5acAE20D5Bad5f6215451dBdA09c00A523",
	}

	sw, err := NewSoftWallet()
	err = json.Unmarshal([]byte(wc.WalletCipher), &sw.walletFileInfo)
	assert.NoError(t, err)

	_, WalletPlain, keyData, err := sw.decryptWalletFromJsonData([]byte(wc.WalletCipher), password)
	assert.NoError(t, err)

	//Convert the wallet plaintext byte to struct by json decoding
	tempWalletInfo := NewHdWalletInfo()
	err = tempWalletInfo.HdWalletInfoDecodeJson(WalletPlain)
	assert.NoError(t, err)
	sw.symmetricKey = keyData

	sw.walletInfo = *tempWalletInfo
	sw.status = accounts.Opened

	message := "0x000087Be6c42Ca4F12D3203E3452D198F2581Fc2D010" + "1" + "0x000087Be6c42Ca4F12D3203E3452D198F2581Fc2D010"
	mByte := crypto.Keccak256([]byte(message))
	accounts, err := sw.Accounts()
	assert.NoError(t, err)
	signByte, err := sw.SignHash(accounts[0], mByte)
	assert.NoError(t, err)

	signHex := common.Bytes2Hex(signByte)
	fmt.Println(signHex)

}

//func TestRecoverAndEncryptAgagin(t *testing.T){
//	//wc := walletConf{
//	//	WalletCipher:`{"Cipher":"24eb4f6fec03f6629663e513ebb14bf52a911dc7e312f2f3da3e4dd3b0a83c11076e5038d22bb53a1c4acb64498f6fdf8adea9c7dbd7d60239ddb46a72e6a6e0ba3fdb29144ce3827fb67e2597e30cc1456ac4427124b0e0154e6034f1f0a88b28e7f7ca78299c26e74236f65ebcb6062d10778ede3215fdb41f2bcd95efe0a4633f4b71b96ffee4d0bf6a3e9654140b7abc08e72f210e321c1c2956eb8bc8278a91a7a5da51eafec41abd3ac16182b7a6591659332c92586417da702ac2e86542146986f080947b900b751840bb11d072607f45a70d4c5369ed794fd974fa32be7057451d0c71c7d738eb498b60fc4deb9fef5a0ac54fc39ee74470c1455e9fc8abe6d1396fdeedd9d7ea2cbef48a0c7034584fa9f26f5a7785ce842e2094d9bfb5904dfcc60293bbdad74a65cda404388704fcf2515e4928fbaf3e3a77b3d4db4812b713ab2be96c1ed97187595f1102f2e4ae351fb7ef59e47a9776f007fdffab28e0607107cc225826e91ac42ce13ca1673b867a864f4eb4abcddfc516c719909060b249be7b9ea01a605e31079320f9a5c0afd6258b5388962987b844366d9d99bd04e9858a51271f06e76a1d0ac4345efce5a7fceb43ca9a631d7cf51c3aed9e41329f9e036ec2889173183fe2f80c025aed07dc462d951fa5dfe6fa4e6c08a48c0364bceb3c0e9a05024e536228b2a2105fad3d1c8e89984caa92ad01d10ccf0479a7d4ee288131d27ce1d192121ddf8d513eb01e8ba5b952b7bc76dbeeb5f27125a9d07404af5df5ba8e20240b2dd5f133ea6de87d8de39b614eb18ee00747c71a7108a3dbd9e24f9674599cff9fc95ab5fca413f488e5e047eee8c0f98548bc5686ff1561060a2ab9dd99b6399ea2043c8fe8452250b6992267c5568a9afbe2e768c9713dc49f459e2e309feda0177ae7325976974c177db12325ec03aab475cd0338904f85e0584cd4bd27bbf8d5e9f200c6099b04307049eb35d9c77c11c68a5104a42d991f97340a65d3cc8d5525ddcfbd7639b063aac29fc20edfc037c484806e614692f140835da9ff5b5df8cf134bdff11f049acb2c51647e316a82f4454250802ca9a1873654e25f5c8c87ae64496262c04255d5daee3e82d642206f1bcd4428af62af01c80788448a8b37b16329cdab9a7fb3a455675b0a2a2a7b2b4b23588e7d626900230668cef992c27b69af83e8b42cbf6d808a84bf93fe6a4ef93a593171f743842e4c30c89541fb4a603d72a3aa67d817aac45b8f653ef7dffc7bd50c20482bbc3fdd2e2a496b79daf6b5e6c497d0242807fdf8e7e199a375b770a5278d731632d6cf54ecafc05fcf55bfc56741ae14ee1cd47169ee03f4b2f0a0b8fdf529f1b5d16b3480cd7068942c2e3b84f125c08c17dab30b7a28b0332f7bd7418d6327d2383911dde6cfb1360cb5a1b0d9801f1dcb63197a943ef945efa7cb9ef64754ae662cca05517545507c3ece7adbef43c91b6e83404e9dacd05fbbe65d1a8f30253a266d6fffc68a8026497b645c5d17ff23f029d4ab563db1d77728da1dfa16204b9deb604d9e86f53097a8c9d85285349c8fa54d5cb35407760930aa46c89b013266ee337d0a373c13b773d28c7cc0c3e99281d92fb3b3cd93c77dc36409894aac3ccbee6d72e506f6eed6b69245d13f56d64e231045500cf3ace8254344d572df8d4cd89324f04849476447b4594067b7c6166f3ff5cfcead4e8dde82f85ef0742b5d08a16792202c1d7c591eea68b01229a150712b332d44c0fc6eb7455a805f846334","MacCipher":"75b5141d8a8b0da00878816f5e33ebd8119d921cb3e1e6ba6a352633b3ce4221","AlgType":"AES-256","ModeType":"CBC","IV":[182,15,4,128,3,83,190,237,105,97,118,41,104,87,220,59],"kdf":"","kdfparams":{"kdfType":"Scrypt","keyLen":32,"n":4096,"p":6,"r":8,"salt":"6b46ab38c2ca8e53b47109b062ce976f281cbb305bc52d5b4aef868b88a85b5f"}}`,
//	//	MainAddress:"0x0000cc8bD424f554c539746E849f95BF19CCa3292fd1",
//	//}
//
//	wfi := `{"Cipher":"b39010a4328d265640ab116c64f8ad8bc40caa659d302509c7d27300be990138eeedef3d1bb56cfa50ba28a77f3b521d008b181a3535b619cadb220be786ea62c8341e8dd29554ddb7ba023cb9f35f566df2f271ed890dc140005144928901a31df718c32920e232e23d5a68cb967363a511295a3e5aca065afb578d0604fcaaa9b477fb37cb955d18aba89802e7d52030ea3249f27bccc742240ab6da7e26d332c25c26f7fa688c33df77e85f6fce8f3864cb17557ea90bbe034f68e21ba556a42a9b560464aa61dfcac0ae451fa854bc04d142dbe9fd791dc18674e6abfcbbd478ec0a0151f582c37983902ea15433244169f5260ab227ed9ec27efd68dca4857b7854451a22f4db7a930515b25697b4aec455144565c5ad2e930a06425dcb06872bad1da776def89946a03ddc923e085cbe0c716295f28fe222c15bed2a1ac6258295d780e2b5ec4b6864a7cbb5f9aa927fc7519b1f65ecc37ad3b3f3c797f76f2c093ef36c39cd04fefac8b67b23d62ef5d90f9619e19afed86323dd4111ab9a9b269d6c195f94f0e3fedbacbf2581c91183fea11376f0415007144cafe5dd6f2fa6e849861fa860db5ecd001e438989b85d7657dac6af3aeec8a00bb4497fd1bc18debc61d0cca0fbd9c336bb3428fc7fa2d5ed5e4a46f4b5f79378243c1414d0c77115802393eaed78eb466819c7bf080a62e5cb00342828299f4e88145af0a0cbc47c89871d1d08e8dffb313a1160e3a99d36a4eb944517304e83d2d73c86c4de14241fa6c893f938abd59e7df518b67ee326b2835c48cf30d45ccfefb0f5dd18907072c4e8c26a0b9e9021fed7fde0efe14f50880513284199a2a62235ef5af384e4a4cfb7725eeb85309c0eed6e99e23cc2a71035f29b37c08a530f379b4592302107302e796ad8157a2592e56b3297e30ed16424b4132fa84004211857a1f63a8cd13e81ee9a2fc06e2b1af2203fe0d805db2b4283faeba4a8eeff203c5085e84429c05ef6439113a657bc","MacCipher":"7cff4b4ceb5607dcdcdda5f798912ff5d467121da2d62bbb7db0f8bf0f54c50d","AlgType":"AES-256","ModeType":"CBC","IV":[141,231,159,105,106,136,70,74,134,70,75,236,186,156,231,94],"kdf":"","kdfparams":{"kdfType":"Scrypt","keyLen":32,"n":"4096","p":"6","r":"8","salt":"310e1c65a46f8f5a2b48a53e3149555d7a21ceb16a380aa8b5e046feb2c7eda1"}}`
//	sw, err := NewSoftWallet()
//	err = json.Unmarshal([]byte(wfi), &sw.walletFileInfo)
//	assert.NoError(t, err)
//
//	err = json.Unmarshal([]byte(wc.WalletCipher), &sw.walletFileInfo)
//	assert.NoError(t, err)
//
//	_, WalletPlain, keyData, err := sw.decryptWalletFromJsonData([]byte(wc.WalletCipher), password)
//	assert.NoError(t, err)
//
//	//Convert the wallet plaintext byte to struct by json decoding
//	tempWalletInfo := NewHdWalletInfo()
//	err = tempWalletInfo.HdWalletInfoDecodeJson(WalletPlain)
//	assert.NoError(t, err)
//	sw.symmetricKey = keyData
//
//	sw.walletInfo = *tempWalletInfo
//	sw.status = accounts.Opened
//
//	message := "0x0014659c3Bd6983c5CEd1De23976d2f3907504e80dD1" + "1" + "0x00002fDc5F7489DA4877561CDe24E337136aF28800FF";
//	mByte := crypto.Keccak256([]byte(message))
//	accounts,err := sw.Accounts()
//	assert.NoError(t, err)
//	signByte , err := sw.SignHash(accounts[0],mByte )
//	assert.NoError(t, err)
//
//	signHex := common.Bytes2Hex(signByte)
//	fmt.Println(signHex)
//
//}
