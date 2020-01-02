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
	"encoding/json"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/accounts"
	"go.uber.org/zap"
	"math/big"
	"sync"
)

const (
	WalletEntropyLength = 256
	WalletDefaultPath   = "/tmp/CSWallet"
	WalletDefaultName   = "CSWallet"
)

const (
	CloseWallet = iota
	RestoreWallet
	EstablishWallet
)

const (
	// m / purpose' / coin_type' / account' / change / address_index
	DefaultDerivedPath       = "m/44'/709394'/0'/0"
	DefaultDerivedPathLength = 4
	DefaultAccountValue      = 0
	AccountValueIndex        = 2
	SyncAccountNumber        = 20
	AddressIndexStartValue   = 1
)

type WalletInfo struct {
	Accounts   []accounts.Account                         // used account included in the wallet
	Paths      map[common.Address]accounts.DerivationPath // derived path of the account in the wallet
	ExtendKeys map[common.Address]ExtendedKey             //Key information corresponding to the account address in the wallet
	Balances   map[common.Address]*big.Int                //The balance corresponding to each account address in the wallet
	Nonce      map[common.Address]uint64                  //The nonce value corresponding to each account address in the wallet
	//the used largest index in wallet derivation path, the key is the changeValue to identify the derived path, and the value is the largest index used.
	DerivedPathIndex map[uint32]uint32
	Seed             []byte //Wallet seed

	//Get the balance and nonce value corresponding to the address
	lock sync.RWMutex
}

func NewHdWalletInfo() (info *WalletInfo) {

	walletInfo := WalletInfo{
		Accounts:         make([]accounts.Account, 0),
		Paths:            make(map[common.Address]accounts.DerivationPath, 0),
		ExtendKeys:       make(map[common.Address]ExtendedKey),
		Balances:         make(map[common.Address]*big.Int, 0),
		Nonce:            make(map[common.Address]uint64, 0),
		DerivedPathIndex: make(map[uint32]uint32, 0),
		Seed:             make([]byte, 0),
	}

	return &walletInfo
}

type ExtendedKeyJson struct {
	Key       []byte `json:"Key"`    // This will be the pubkey for extended pub keys
	PubKey    []byte `json:"PubKey"` // This will only be set for extended priv keys
	ChainCode []byte `json:"ChainCode"`
	Depth     uint8  `json:"Depth"`
	ParentFP  []byte `json:"ParentFP"`
	ChildNum  uint32 `json:"ChildNum"`
	Version   []byte `json:"Version"`
	IsPrivate bool   `json:"IsPrivate"`
}

type WalletInfoJson struct {
	Accounts         []accounts.Account                 `json:"accounts"`
	Paths            map[string]accounts.DerivationPath `json:"paths"`
	ExtendKeys       map[string]ExtendedKeyJson         `json:"extend_keys"`
	Balances         map[string]*big.Int                `json:"balances"`
	Nonce            map[string]uint64
	DerivedPathIndex map[uint32]uint32
	Seed             []byte `json:"seed"`
}

func NewHdWalletInfoJson() (jsonInfo *WalletInfoJson) {
	w := &WalletInfoJson{
		Accounts:         make([]accounts.Account, 0),
		Paths:            make(map[string]accounts.DerivationPath, 0),
		ExtendKeys:       make(map[string]ExtendedKeyJson, 0),
		Balances:         make(map[string]*big.Int, 0),
		Nonce:            make(map[string]uint64, 0),
		DerivedPathIndex: make(map[uint32]uint32, 0),
		Seed:             make([]byte, 0),
	}
	return w
}

//Json encode for HdWalletInfo, converts key to string type in map
func (w WalletInfo) HdWalletInfoEncodeJson() (encodeData []byte, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	tmpData := NewHdWalletInfoJson()

	tmpData.Seed = w.Seed

	for _, account := range w.Accounts {
		tmpData.Accounts = append(tmpData.Accounts, account)

		tmpData.Paths[account.Address.Hex()] = w.Paths[account.Address]
		tmpData.ExtendKeys[account.Address.Hex()] = ExtendedKeyJson{
			Key:       w.ExtendKeys[account.Address].key,
			PubKey:    w.ExtendKeys[account.Address].pubKey,
			ChainCode: w.ExtendKeys[account.Address].chainCode,
			Depth:     w.ExtendKeys[account.Address].depth,
			ParentFP:  w.ExtendKeys[account.Address].parentFP,
			ChildNum:  w.ExtendKeys[account.Address].childNum,
			Version:   w.ExtendKeys[account.Address].version,
			IsPrivate: w.ExtendKeys[account.Address].isPrivate,
		}
		tmpData.Balances[account.Address.Hex()] = w.Balances[account.Address]
		tmpData.Nonce[account.Address.Hex()] = w.Nonce[account.Address]
		tmpData.DerivedPathIndex = w.DerivedPathIndex
	}

	return json.Marshal(tmpData)

}

//wallet data json decoding interface
func (w *WalletInfo) HdWalletInfoDecodeJson(decodeData []byte) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	tmpData := NewHdWalletInfoJson()
	//log.DLogger.Debug(string(decodeData))
	err = json.Unmarshal(decodeData, tmpData)
	if err != nil {
		return err
	}

	w.Seed = tmpData.Seed

	for _, account := range tmpData.Accounts {
		w.Paths[account.Address] = tmpData.Paths[string(account.Address[:])]

		w.Accounts = append(w.Accounts, account)

		w.ExtendKeys[account.Address] = ExtendedKey{
			key:       tmpData.ExtendKeys[account.Address.Hex()].Key,
			pubKey:    tmpData.ExtendKeys[account.Address.Hex()].PubKey,
			chainCode: tmpData.ExtendKeys[account.Address.Hex()].ChainCode,
			depth:     tmpData.ExtendKeys[account.Address.Hex()].Depth,
			parentFP:  tmpData.ExtendKeys[account.Address.Hex()].ParentFP,
			childNum:  tmpData.ExtendKeys[account.Address.Hex()].ChildNum,
			version:   tmpData.ExtendKeys[account.Address.Hex()].Version,
			isPrivate: tmpData.ExtendKeys[account.Address.Hex()].IsPrivate,
		}
		w.Balances[account.Address] = tmpData.Balances[string(account.Address[:])]
		w.Nonce[account.Address] = tmpData.Nonce[string(account.Address[:])]
		w.DerivedPathIndex = tmpData.DerivedPathIndex
	}

	return nil
}

//Generate a master key according to the seed. Then generate a key corresponding to the index on the specified derived path according to the master key
func (w *WalletInfo) GenerateKeyFromSeedAndPath(derivedPath string, index uint32) (derivedKey *ExtendedKey, Path accounts.DerivationPath, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	log.DLogger.Debug("the wallet seed is：", zap.String("seed", hexutil.Encode(w.Seed)))
	//generate master key according to seed
	extKey, err := NewMaster(w.Seed, &DipperinChainCfg)
	if err != nil {
		return nil, accounts.DerivationPath{}, err
	}

	//Get the first account information on the default derived path based on the master key
	Path, err = accounts.ParseDerivationPath(derivedPath)
	if err != nil {
		ClearSensitiveData(extKey)
		return nil, accounts.DerivationPath{}, err
	}
	Path = append(Path, index)

	log.DLogger.Debug("the Path is: ", zap.Any("path", Path))
	for _, value := range Path {
		extKey, err = extKey.Child(value)
		if err != nil {
			return nil, accounts.DerivationPath{}, err
		}
	}
	return extKey, Path, nil
}

//get the  first 20 account according to the derived path, query whether it is used, and if it is used, join the recovered wallet.
func (w *WalletInfo) paddingUsedAccount(GetAddressRelatedInfo accounts.AddressInfoReader) (err error) {
	tmpPath, err := accounts.ParseDerivationPath(DefaultDerivedPath)
	if err != nil {
		return err
	}

	tmpPath = append(tmpPath, w.DerivedPathIndex[DefaultAccountValue])
	for i := 0; i < SyncAccountNumber; i++ {
		//derived path to use
		tmpPath[len(tmpPath)-1] = w.DerivedPathIndex[DefaultAccountValue] + 1

		log.DLogger.Info("the tmpPath is:", zap.String("tmpPath", tmpPath.String()))
		//determine if the derived path is legal
		isValid, err := CheckDerivedPathValid(tmpPath)
		if err != nil || !isValid {
			return accounts.ErrInvalidDerivedPath
		}

		//Generate derived keys based on incoming derived paths and wallet seeds
		extKey, err := NewMaster(w.Seed, &DipperinChainCfg)
		if err != nil {
			return err
		}

		log.DLogger.Info("Derive tmpPath is:", zap.Any("tmpPath", tmpPath))
		//Generate derived keys based on path parameters and master key
		for _, value := range tmpPath {
			var err error
			extKey, err = extKey.Child(value)
			if err != nil {
				return err
			}
		}

		account, err := GetAccountFromExtendedKey(extKey)
		if err != nil {
			return err
		}

		//check if the derived account is used
		nonce, err := GetAddressRelatedInfo.GetTransactionNonce(account.Address)
		if err != nil {
			if err == g_error.ErrAccountNotExist {
				break
			} else {
				return err
			}
		} else {
			log.DLogger.Info("the used address is:", zap.String("address", account.Address.Hex()))
			balance := GetAddressRelatedInfo.CurrentBalance(account.Address)
			w.Accounts = append(w.Accounts, account)
			w.Paths[account.Address] = tmpPath
			w.Nonce[account.Address] = nonce
			w.Balances[account.Address] = balance
			w.ExtendKeys[account.Address] = *extKey
			w.DerivedPathIndex[DefaultAccountValue] += 1
		}
	}
	return nil
}

//Obtain the private key data corresponding to the address
func (w *WalletInfo) getSkFromAddress(address common.Address) (sk *ecdsa.PrivateKey, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	tmpKey, ok := w.ExtendKeys[address]
	if !ok {
		return nil, accounts.ErrInvalidAddress
	}

	//take out the private key in the extension key
	privateKey, err := tmpKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	tmpSk := &ecdsa.PrivateKey{
		PublicKey: privateKey.PublicKey,
		D:         privateKey.D,
	}

	return tmpSk, nil
}

func (w *WalletInfo) PaddingAddressNonce(GetAddressRelatedInfo accounts.AddressInfoReader) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	for _, account := range w.Accounts {
		currentNonce, err := GetAddressRelatedInfo.GetTransactionNonce(account.Address)
		log.DLogger.Error("PaddingAddressNonce GetTransactionNonce error", zap.Error(err))
		log.DLogger.Info("the padding address is:", zap.String("address", account.Address.Hex()))
		log.DLogger.Info("PaddingAddressNonce is: ", zap.Uint64("currentNonce", currentNonce))
		w.Nonce[account.Address] = currentNonce
	}
	return nil
}

func (w *WalletInfo) GetAddressNonce(address common.Address) (nonce uint64, err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.Nonce[address], nil
}

//add nonce when send transaction
func (w *WalletInfo) SetAddressNonce(address common.Address, nonce uint64) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.Nonce[address] = nonce
	return nil
}
