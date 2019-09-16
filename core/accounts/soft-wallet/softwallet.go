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
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	crypto2 "github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/go-bip39"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

type SoftWallet struct {
	walletInfo   WalletInfo //clear wallet data
	symmetricKey EncryptKey //encrypt key and mac key

	walletFileInfo WalletFileContent //wallet file data

	status string       //wallet status　"open"　or "close"
	mu     sync.RWMutex //wallet operation lock

	Identifier accounts.WalletIdentifier //Wallet identifier
}

func NewSoftWallet() (*SoftWallet, error) {
	tmpWalletInfo := NewHdWalletInfo()
	wallet := &SoftWallet{
		walletInfo:   *tmpWalletInfo,
		symmetricKey: EncryptKey{},
		walletFileInfo: WalletFileContent{
			WalletCipher{
				Cipher: "", MacCipher: "",
			},
			EncryptParameter{
				SymmetricAlgorithm{},
				KDFParameter{KDF: "", KDFParams: make(map[string]interface{}, 0)},
			},
		},
		status:     accounts.Closed,
		mu:         sync.RWMutex{},
		Identifier: accounts.WalletIdentifier{WalletType: accounts.SoftWallet, Path: "", WalletName: ""},
	}

	return wallet, nil
}

//Generate relevant wallet information according to mnemonics, passwords, mnemonic passwords, wallet file storage paths, and KDF key-derived parameters
func (w *SoftWallet) paddingWalletInfo(mnemonic, password, passPhrase string, kdfPara *KDFParameter) (err error) {
	//log.Debug("the kdfPara is: ","kdfPara",kdfPara)
	if kdfPara == nil {
		//If no parameters are passed in when creating a new one, use the default value.
		w.walletFileInfo.KDFParams["n"] = WalletLightScryptN
		w.walletFileInfo.KDFParams["p"] = WalletLightScryptP
		w.walletFileInfo.KDFParams["kdfType"] = KDF
		w.walletFileInfo.KDFParams["r"] = WalletscryptR
		w.walletFileInfo.KDFParams["keyLen"] = WalletscryptDKLen
		//randomly generated salt value
		tmpSalt := cspRngEntropy(32)
		//log.Debug("the generate tmpSalt is: ","tmpSalt",tmpSalt)
		w.walletFileInfo.KDFParams["salt"] = hex.EncodeToString(tmpSalt)
	} else {
		w.walletFileInfo.KDFParams = kdfPara.KDFParams
	}

	//Derived key according to password and KDF parameters
	w.symmetricKey, err = GenSymKeyFromPassword(password, w.walletFileInfo.EncryptParameter.KDFParameter)
	if err != nil {
		return err
	}

	//Randomly generate the iv initial vector used for encryption
	copy(w.walletFileInfo.IV[:], cspRngEntropy(symmetricEncryptLen))
	w.walletFileInfo.AlgType = SymmetricAlgType
	w.walletFileInfo.ModeType = SymmetricAlgMode

	//generate a seed based on the mnemonic
	w.walletInfo.Seed, err = bip39.NewSeedWithErrorChecking(mnemonic, passPhrase)
	if err != nil {
		return err
	}

	//log.Debug("the seed is: ","seed",w.WalletInfo.Seed)
	//get account based on master key
	var tmpPath accounts.DerivationPath
	extKey, tmpPath, err := w.walletInfo.GenerateKeyFromSeedAndPath(DefaultDerivedPath, AddressIndexStartValue)
	if err != nil {
		ClearSensitiveData(extKey)
		return err
	}

	account, err := GetAccountFromExtendedKey(extKey)
	if err != nil {
		ClearSensitiveData(extKey)
		return err
	}

	w.walletInfo.Accounts = append(w.walletInfo.Accounts, account)
	w.walletInfo.Paths[account.Address] = tmpPath
	w.walletInfo.DerivedPathIndex[DefaultAccountValue] = AddressIndexStartValue

	//log.Debug("the extKey is: ","extKey",*extKey)

	w.walletInfo.ExtendKeys[account.Address] = *extKey

	//log.Debug("wallet info after padding","wallet",*w)

	return nil
}

//Encrypt the wallet plaintext data and write it to the file when creating, closing, or restoring the wallet
func (w *SoftWallet) encryptWalletAndWriteFile(operation int) (err error) {
	//Generate wallet cipher and write it to the wallet storage path
	w.walletFileInfo.WalletCipher, err = CalWalletCipher(w.walletInfo, w.walletFileInfo.IV[:], w.symmetricKey)
	if err != nil {
		return err
	}

	log.Debug("encryptWalletAndWriteFile 1")

	//log.Debug("===establish: the walletFileInfo salt is: ","salt",w.walletFileInfo.KDFParams["salt"])
	//write wallet cipher to local file
	writeData, err := json.Marshal(w.walletFileInfo)
	if err != nil {
		return err
	}

	var walletPath string
	if w.Identifier.Path == "" {
		walletPath = WalletDefaultPath
	} else {
		walletPath = w.Identifier.Path
	}

	//When closing the wallet operation, it is judged whether the wallet file exists, and the error is returned if it no longer exists.
	exist, _ := PathExists(walletPath)
	if operation == CloseWallet {
		if exist == false {
			//wallet file does not exist
			return accounts.ErrWalletFileNotExist
		}
	} else {
		if exist == true {
			//file already exists when creating a new wallet
			return accounts.ErrWalletFileExist
		} else {
			path := filepath.Dir(walletPath)
			os.MkdirAll(path, 0766)
		}
	}
	log.Debug("write walletPath", "walletPath", walletPath)
	ioutil.WriteFile(walletPath, writeData, 0666)

	return nil
}

func (w *SoftWallet) decryptWallet(password string) (passwordValid bool, walletPlain []byte, keyData EncryptKey, err error) {
	var walletPath string
	if w.Identifier.Path == "" {
		walletPath = WalletDefaultPath
	} else {
		walletPath = w.Identifier.Path
	}

	//Read wallet cipher and encryption parameter data according to wallet path
	walletJsonData, err := ioutil.ReadFile(walletPath)
	if err != nil {
		log.Info("the err is:", "err", err)
		return
	}

	err = json.Unmarshal(walletJsonData, &w.walletFileInfo)
	if err != nil {
		return
	}

	gj := gjson.ParseBytes(util.StringifyJsonToBytes(w.walletFileInfo.KDFParams))
	w.walletFileInfo.KDF = gj.Get("kdf").String()
	w.walletFileInfo.KDFParams["kdfType"] = gj.Get("kdfType").String()
	w.walletFileInfo.KDFParams["keyLen"] = gj.Get("keyLen").Int()
	w.walletFileInfo.KDFParams["n"] = gj.Get("n").String()
	w.walletFileInfo.KDFParams["r"] = gj.Get("r").String()
	w.walletFileInfo.KDFParams["p"] = gj.Get("p").String()

	//Derive encrypt key and mac key according to password
	keyData, err = GenSymKeyFromPassword(password, w.walletFileInfo.KDFParameter)
	if err != nil {
		return
	}

	//decrypt wallet plaintext
	WalletPlain, err1 := DecryptWalletContent(w.walletFileInfo.WalletCipher, w.walletFileInfo.IV[:], keyData)
	if err1 != nil {
		log.Warn("decrypt wallet failed", "err", err1)
		err = accounts.ErrWalletPasswordNotValid
		return
	}

	return true, WalletPlain, keyData, nil
}

//return the soft wallet identifier
func (w *SoftWallet) GetWalletIdentifier() (accounts.WalletIdentifier, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != accounts.Opened {
		return accounts.WalletIdentifier{}, accounts.ErrWalletNotOpen
	}

	return w.Identifier, nil
}

//return the soft wallet status
func (w *SoftWallet) Status() (string, error) {

	w.mu.Lock()
	defer w.mu.Unlock()

	return w.status, nil
}

//create a new soft wallet and return the mnemonic
func (w *SoftWallet) Establish(path, name, password, passPhrase string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	err := CheckWalletPath(path)
	if err != nil {
		return "", err
	}

	err = CheckPassword(password)
	if err != nil {
		return "", err
	}
	w.Identifier.WalletName = name
	w.Identifier.Path = path
	w.Identifier.WalletType = accounts.SoftWallet

	mnemonic, err := GenerateMnemonic(WalletEntropyLength)
	if err != nil {
		return "", err
	}

	//fill wallet related information
	//kdfPara default program built in
	err = w.paddingWalletInfo(mnemonic, password, passPhrase, nil)
	if err != nil {
		log.Info("paddingWalletInfo error", "err", err)
		return "", err
	}
	//The main account balance is 0 when creating a new wallet
	w.walletInfo.Balances[w.walletInfo.Accounts[0].Address] = big.NewInt(0)

	log.Debug("set wallet status is open")
	w.status = accounts.Opened

	err = w.encryptWalletAndWriteFile(EstablishWallet)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

//recover wallet based on mnemonic
func (w *SoftWallet) RestoreWallet(path, name, password, passPhrase, mnemonic string, GetAddressRelatedInfo accounts.AddressInfoReader) (err error) {

	w.mu.Lock()
	defer w.mu.Unlock()

	err = CheckWalletPath(path)
	if err != nil {
		return err
	}

	err = CheckPassword(password)
	if err != nil {
		return err
	}

	w.Identifier.WalletName = name
	w.Identifier.Path = path
	w.Identifier.WalletType = accounts.SoftWallet

	//fill wallet related information
	//kdfPara default program built in
	err = w.paddingWalletInfo(mnemonic, password, passPhrase, nil)
	if err != nil {
		return err
	}

	err = w.walletInfo.paddingUsedAccount(GetAddressRelatedInfo)
	if err != nil {
		return err
	}

	w.status = accounts.Opened
	// write recovered wallet data to a local file
	err = w.encryptWalletAndWriteFile(RestoreWallet)
	if err != nil {
		return err
	}

	return nil
}

//open the soft wallet according to the password
func (w *SoftWallet) Open(path, name, password string) error {
	//var keyData EncryptKey

	w.mu.Lock()
	defer w.mu.Unlock()

	err := CheckWalletPath(path)
	if err != nil {
		return err
	}

	err = CheckPassword(password)
	if err != nil {
		return err
	}

	w.Identifier.Path = path
	w.Identifier.WalletName = name

	_, WalletPlain, keyData, err := w.decryptWallet(password)
	if err != nil {
		return err
	}

	//Convert the wallet plaintext byte to struct by json decoding
	tempWalletInfo := NewHdWalletInfo()
	err = tempWalletInfo.HdWalletInfoDecodeJson(WalletPlain)
	if err != nil {
		return err
	}

	w.symmetricKey = keyData

	w.walletInfo = *tempWalletInfo
	w.status = accounts.Opened

	return nil
}

//close the soft wallet
func (w *SoftWallet) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != accounts.Opened {
		return accounts.ErrWalletNotOpen
	}

	//Calculate wallet data cipher and mac value using wallet internal derivative key
	var err error
	w.walletFileInfo.WalletCipher, err = CalWalletCipher(w.walletInfo, w.walletFileInfo.SymmetricAlgorithm.IV[:], w.symmetricKey)
	if err != nil {
		return err
	}

	//Set the wallet status to "close". Encrypt the wallet data and write to a local file.
	w.status = accounts.Closed

	err = w.encryptWalletAndWriteFile(CloseWallet)
	if err != nil {
		return err
	}

	return nil
}

//get a list of accounts in your soft wallet
func (w *SoftWallet) Accounts() ([]accounts.Account, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.status != accounts.Opened {
		return []accounts.Account{}, accounts.ErrWalletNotOpen
	}
	return w.walletInfo.Accounts, nil
}

//determine if the soft wallet contains an account
func (w *SoftWallet) Contains(account accounts.Account) (bool, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.status != accounts.Opened {
		return false, accounts.ErrWalletNotOpen
	}
	for _, tmpAccount := range w.walletInfo.Accounts {
		if tmpAccount == account {
			return true, nil
		}
	}
	return false, nil
}

//Generate an account based on the input derived path and add it to SoftWallet
func (w *SoftWallet) Derive(path accounts.DerivationPath, save bool) (accounts.Account, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != accounts.Opened {
		return accounts.Account{}, accounts.ErrWalletNotOpen
	}

	//use default path when no derived path is provided
	var tmpPath accounts.DerivationPath
	var err error

	log.Info("the derive path is: ", "path", path)

	if path.String() == "m" {
		tmpPath, err = accounts.ParseDerivationPath(DefaultDerivedPath)
		if err != nil {
			return accounts.Account{}, err
		}
		//used derivation path
		tmpPath = append(tmpPath, w.walletInfo.DerivedPathIndex[DefaultAccountValue]+1)
		w.walletInfo.DerivedPathIndex[DefaultAccountValue] += 1
	} else {
		tmpPath = path
		changeValue := tmpPath[AccountValueIndex]
		if value, ok := w.walletInfo.DerivedPathIndex[changeValue]; ok {
			tmpPath = append(tmpPath, value+1)
			w.walletInfo.DerivedPathIndex[DefaultAccountValue] += 1
		} else {
			tmpPath = append(tmpPath, 0)
			w.walletInfo.DerivedPathIndex[DefaultAccountValue] = 0
		}
	}

	//determine if the derived path is legal
	isValid, err := CheckDerivedPathValid(tmpPath)
	if err != nil || !isValid {
		return accounts.Account{}, accounts.ErrInvalidDerivedPath
	}

	//Generate derived keys based on incoming derived paths and wallet seeds
	extKey, err := NewMaster(w.walletInfo.Seed, &DipperinChainCfg)
	if err != nil {
		return accounts.Account{}, err
	}

	log.Info("Derive tmpPath is:", "tmpPath", tmpPath)
	//Generate derived keys based on path parameters and master key
	for _, value := range tmpPath {
		var err error
		extKey, err = extKey.Child(value)
		if err != nil {
			return accounts.Account{}, err
		}
	}

	account, err := GetAccountFromExtendedKey(extKey)
	if err != nil {
		return accounts.Account{}, err
	}

	//determine if the account already exists
	for _, tmpAccount := range w.walletInfo.Accounts {
		if tmpAccount == account {
			return account, nil
		}
	}

	//According to the incoming pin, judge whether to save it in the wallet file.
	if save == true {
		//Store the generated account public and private key in the wallet
		w.walletInfo.Accounts = append(w.walletInfo.Accounts, account)
		w.walletInfo.ExtendKeys[account.Address] = *extKey
		w.walletInfo.Paths[account.Address] = path
		w.walletInfo.Balances[account.Address] = big.NewInt(0)
	}

	//update wallet file
	w.walletFileInfo.WalletCipher, err = CalWalletCipher(w.walletInfo, w.walletFileInfo.SymmetricAlgorithm.IV[:], w.symmetricKey)
	if err != nil {
		return accounts.Account{}, err
	}
	err = w.encryptWalletAndWriteFile(CloseWallet)
	if err != nil {
		return accounts.Account{}, err
	}

	ClearSensitiveData(extKey)

	return account, nil
}

//According to the base path, query the used account from the chain and add it to the wallet.
func (w *SoftWallet) SelfDerive(base accounts.DerivationPath) error {
	return nil
}

//Sign the hash value with its corresponding private key based on the incoming account
func (w *SoftWallet) SignHash(account accounts.Account, hash []byte) ([]byte, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.status != accounts.Opened {
		return []byte{}, accounts.ErrWalletNotOpen
	}

	//Obtain the corresponding private key data according to the address
	tmpSk, err := w.walletInfo.getSkFromAddress(account.Address)
	if err != nil {
		return nil, err
	}

	signData, err := crypto.Sign(hash, tmpSk)
	if err != nil {
		ClearSensitiveData(tmpSk)
		return []byte{}, err
	}

	ClearSensitiveData(tmpSk)
	//sign the hash data with the private key
	return signData, nil
}

//Sign the transaction with its corresponding private key based on the incoming account
func (w *SoftWallet) SignTx(account accounts.Account, tx *model.Transaction, chainID *big.Int) (*model.Transaction, error) {
	// Transaction signature operation with the private key based on the account
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.status != accounts.Opened {
		return nil, accounts.ErrWalletNotOpen
	}

	//Obtain the corresponding private key data according to the address
	priKeys, err := w.walletInfo.getSkFromAddress(account.Address)
	if err != nil {
		return nil, err
	}

	var s model.Signer = model.NewSigner(chainID)

	signedTx, err := tx.SignTx(priKeys, s)
	if err != nil {
		return nil, err
	}

	//log.Debug("softWallet SignTx end")
	return signedTx, nil
}

//generate vrf proof using private key and seed
func (w *SoftWallet) Evaluate(account accounts.Account, seed []byte) (index [32]byte, proof []byte, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.status != accounts.Opened {
		return [32]byte{}, []byte{}, accounts.ErrWalletNotOpen
	}

	//Obtain the corresponding private key data according to the address
	priKeys, err := w.walletInfo.getSkFromAddress(account.Address)
	if err != nil {
		return [32]byte{}, []byte{}, err
	}

	index, proof = crypto2.Evaluate(priKeys, seed)

	return index, proof, nil
}

func (w *SoftWallet) GetPKFromAddress(account accounts.Account) (*ecdsa.PublicKey, error) {
	//get sk according to the address
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.status != accounts.Opened {
		return nil, accounts.ErrWalletNotOpen
	}

	if sk, ok := w.walletInfo.ExtendKeys[account.Address]; ok {
		//generate pk according to the sk
		skByte := sk.pubKeyBytes()
		return crypto.DecompressPubkey(skByte)
	} else {
		return nil, accounts.ErrInvalidAddress
	}
}

func (w *SoftWallet) GetSKFromAddress(address common.Address) (*ecdsa.PrivateKey, error) {
	//get sk according to the address
	w.mu.RLock()
	defer w.mu.RUnlock()
	if w.status != accounts.Opened {
		return nil, accounts.ErrWalletNotOpen
	}

	if sk, ok := w.walletInfo.ExtendKeys[address]; ok {
		//generate pk according to the sk
		privateKey, err := sk.ECPrivKey()
		if err != nil {
			return nil, err
		}

		result := ecdsa.PrivateKey{
			PublicKey: privateKey.PublicKey,
			D:         privateKey.D,
		}
		return &result, nil
	} else {
		return nil, accounts.ErrInvalidAddress
	}
}

func (w *SoftWallet) PaddingAddressNonce(GetAddressRelatedInfo accounts.AddressInfoReader) (err error) {
	return w.walletInfo.PaddingAddressNonce(GetAddressRelatedInfo)
}

func (w *SoftWallet) GetAddressNonce(address common.Address) (nonce uint64, err error) {
	return w.walletInfo.GetAddressNonce(address)
}

//add nonce when send transaction
func (w *SoftWallet) SetAddressNonce(address common.Address, nonce uint64) (err error) {
	return w.walletInfo.SetAddressNonce(address, nonce)
}
