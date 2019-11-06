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
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/go-bip39"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/tidwall/gjson"
	"golang.org/x/crypto/scrypt"
	"io"
	"regexp"
)

//default kdf parameter
const (
	KDF string = "Scrypt"
	// StandardScryptN is the N parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	WalletStandardScryptN = 1 << 18

	// StandardScryptP is the P parameter of Scrypt encryption algorithm, using 256MB
	// memory and taking approximately 1s CPU time on a modern processor.
	WalletStandardScryptP = 1

	// LightScryptN is the N parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	WalletLightScryptN = 1 << 12

	// LightScryptP is the P parameter of Scrypt encryption algorithm, using 4MB
	// memory and taking approximately 100ms CPU time on a modern processor.
	WalletLightScryptP = 6

	WalletscryptR      = 8
	WalletscryptDKLen  = 32
	WalletMacCipherLen = 32
)

//use AES-128-CBC
const (
	symmetricEncryptLen = 16
	symmetricKeyLen     = 32
	SymmetricAlgType    = "AES-256"
	SymmetricAlgMode    = "CBC"
)

//wallet plaintext
type WalletPlaintext struct {
	plainLen  uint32
	plaintext []byte
}

//wallet cipher
type WalletCipher struct {
	Cipher    string `json:"Cipher"`
	MacCipher string `json:"MacCipher"`
}

//symmetric encryption algorithm parameters
type SymmetricAlgorithm struct {
	AlgType  string                    `json:"AlgType"`
	ModeType string                    `json:"ModeType"`
	IV       [symmetricEncryptLen]byte `json:"IV"`
}

//key derivation algorithm parameters
type KDFParameter struct {
	KDF       string                 `json:"kdf"`
	KDFParams map[string]interface{} `json:"kdfparams"`
}

//encrypted keyData
type EncryptKey struct {
	encryptKey [symmetricKeyLen]byte
	macKey     [symmetricKeyLen]byte
}

//encryption parameters
type EncryptParameter struct {
	SymmetricAlgorithm
	KDFParameter
}

//wallet file content
type WalletFileContent struct {
	WalletCipher
	EncryptParameter
}

//generating mnemonic
func GenerateMnemonic(bitSize int) (mnemonic string, err error) {

	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}

	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

//Derived encrypted key and mac key based on password and KDF parameters
func GenSymKeyFromPassword(password string, kdfPara KDFParameter) (sysKey EncryptKey, err error) {

	authKey := []byte(password)

	gj := gjson.ParseBytes(util.StringifyJsonToBytes(kdfPara.KDFParams))

	//currently only supports scrypt derived keys
	if gj.Get("kdfType").String() != KDF {
		return EncryptKey{}, accounts.ErrNotSupported
	}

	keyLen := gj.Get("keyLen").Int()
	if keyLen != symmetricKeyLen {
		return EncryptKey{}, accounts.ErrInvalidKDFParameter
	}

	saltString := gj.Get("salt").String()
	salt, err := hex.DecodeString(saltString)
	if err != nil {
		return EncryptKey{}, err
	}

	//Generate the key used to encrypt the wallet data according to the password entered by the user.
	deriveKey, result := scrypt.Key(authKey, salt, int(gj.Get("n").Int()), int(gj.Get("r").Int()), int(gj.Get("p").Int()), int(keyLen))
	if (result != nil) || (len(deriveKey) != symmetricKeyLen) {
		return EncryptKey{}, accounts.ErrDeriveKey
	}

	//Mac key and encrypt key first use the same key value encryption and decryption using AES-128-CBC
	copy(sysKey.encryptKey[:], deriveKey[:])
	copy(sysKey.macKey[:], deriveKey[:])

	return sysKey, nil
}

//Obtain account information based on the extended key
func GetAccountFromExtendedKey(keyData *ExtendedKey) (account accounts.Account, err error) {

	tmpPk, err := keyData.ECPubKey()
	if err != nil {
		return accounts.Account{}, err
	}

	p := ecdsa.PublicKey{
		Curve: tmpPk.Curve,
		X:     tmpPk.X,
		Y:     tmpPk.Y,
	}

	account.Address = cs_crypto.GetNormalAddress(p)
	return account, nil
}

//Encrypt wallet plaintext data based on wallet plaintext and derived encrypted key and mac key
func EncryptWalletContent(walletPlain []byte, iv []byte, sysKey EncryptKey) (walletCipher WalletCipher, err error) {

	//padding random number when the plaintext data length is not 16 integer multiples
	encryptData := make([]byte, 4)
	//log.Debug("EncryptWalletContent 1", "encryptData len", len(encryptData))

	binary.BigEndian.PutUint32(encryptData, uint32(len(walletPlain)))

	encryptData = append(encryptData, walletPlain...)
	//log.Debug("EncryptWalletContent 2 ", "encryptData len", len(encryptData))

	calcHashSrcDataLen := len(encryptData)
	//log.Debug("EncryptWalletContent 3 ", "encryptData len", len(encryptData))

	if len(encryptData)%16 != 0 {
		padding := cspRngEntropy(16 - len(encryptData)%16)
		encryptData = append(encryptData, padding...)
	}

	//Calculate the plain text mac value
	hashValue := crypto.Keccak256(encryptData[:calcHashSrcDataLen])

	//calculate cipher value
	mac, err := AesEncryptCBC(iv, sysKey.macKey[:], hashValue)
	if err != nil {
		return walletCipher, err
	}

	walletCipher.MacCipher = hex.EncodeToString(mac)

	cipher, err := AesEncryptCBC(iv, sysKey.encryptKey[:], encryptData)
	if err != nil {
		return walletCipher, err
	}

	walletCipher.Cipher = hex.EncodeToString(cipher)

	return walletCipher, nil
}

//Decrypt and verify wallet ciphertext data based on wallet ciphertext and derived encrypted key and mac key
func DecryptWalletContent(walletCipher WalletCipher, iv []byte, sysKey EncryptKey) (walletPlain []byte, err error) {

	cipher, err := hex.DecodeString(walletCipher.Cipher)
	if err != nil {
		return nil, err
	}
	//Decrypt the wallet data according to the generated derivative key, verify the legality of the wallet file according to the mac value, and open the wallet if legal
	decryptData, err := AesDecryptCBC(iv, sysKey.encryptKey[:], cipher)
	if err != nil {
		return nil, accounts.ErrAESDecryption
	}

	truePlainLen := binary.BigEndian.Uint32(decryptData[:4])
	if int(truePlainLen) > (len(decryptData) - 4) {
		return []byte{}, accounts.ErrMacAuthentication
	}

	plainData := WalletPlaintext{
		plainLen:  binary.BigEndian.Uint32(decryptData[:4]),
		plaintext: decryptData[4 : truePlainLen+4],
	}

	macData, err := hex.DecodeString(walletCipher.MacCipher)
	if err != nil {
		return nil, err
	}

	//decrypt the mac value
	decryptMac, err := AesDecryptCBC(iv, sysKey.macKey[:], macData)
	if err != nil {
		return nil, accounts.ErrAESDecryption
	}

	//check if the mac value is legal
	hashValue := crypto.Keccak256(decryptData[:(4 + plainData.plainLen)])

	for index, value := range hashValue {
		if value != decryptMac[index] {
			return nil, accounts.ErrMacAuthentication
		}
	}

	return plainData.plaintext, nil

}

//calculate wallet cipher
func CalWalletCipher(walletInfo WalletInfo, iv []byte, sysKey EncryptKey) (walletCipher WalletCipher, err error) {

	//Encode plaintext wallet data into []byte form by json
	walletPlain, err := walletInfo.HdWalletInfoEncodeJson()

	if err != nil {
		return WalletCipher{}, err
	}

	//Calculate wallet data cipher and mac value using wallet internal derivative key
	walletCipher, err = EncryptWalletContent(walletPlain, iv, sysKey)
	if err != nil {
		return WalletCipher{}, err
	}

	return walletCipher, nil
}

//judging the legitimacy of derived paths
func CheckDerivedPathValid(path accounts.DerivationPath) (bool, error) {
	defaultPath, err := accounts.ParseDerivationPath(DefaultDerivedPath)
	if err != nil {
		return false, err
	}

	if len(path) != DefaultDerivedPathLength+1 {
		return false, nil
	}

	for i := 0; i < DefaultDerivedPathLength-1; i++ {
		if defaultPath[i] != path[i] {
			return false, nil
		}
	}

	return true, nil
}

//judge the wallet password
func CheckPassword(password string) (err error) {
	reg := regexp.MustCompile(`[0-9]|[a-z]|[A-Z]|[~!@#$%^&*()_+<>?:"{},.\\/;'[\]` + "`]")
	regNoCh := regexp.MustCompile("[\u4e00-\u9fa5]")
    strs := regNoCh.FindAllString(password, -1)

    if reg.MatchString(password) && len(strs) <= 0 {
    	if len(password) >= accounts.PasswordMin && len(password) <= accounts.PassWordMax {
			return  nil
		}
	}
    return accounts.ErrPasswordOrPassPhraseIllegal
}

//judge the incoming wallet path
func CheckWalletPath(path string) (err error) {
	homeDir := util.HomeDir()
	log.Info("the path is:", "path", path)
	log.Info("the home dir is", "homeDir", homeDir)
	if len(path) < len(homeDir) {
		return accounts.ErrWalletPathError
	}

	if path[:len(homeDir)] != homeDir {
		return accounts.ErrWalletPathError
	}

	return nil
}

func cspRngEntropy(n int) []byte {
	buf := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	return buf
}
