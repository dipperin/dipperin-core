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
	"crypto/aes"
	"crypto/cipher"
	"github.com/dipperin/dipperin-core/core/accounts"
)

//AES CBC encryption
func AesEncryptCBC(iv []byte, key []byte,plaintext []byte) (cipherText []byte,err error){

	tmpCipher := make([]byte,len(plaintext))

	c, err := aes.NewCipher(key)
	if err != nil  {
		return nil,err
	}

	encrypt := cipher.NewCBCEncrypter(c, iv)

	encrypt.CryptBlocks(tmpCipher,plaintext)

	cipherText = tmpCipher

	return cipherText,nil
}

//AES CBC decryption
func AesDecryptCBC(iv []byte, key []byte,cipherText []byte) (plaintext []byte,err error){

	if len(cipherText)%16 !=0{
		return nil,accounts.ErrAESInvalidParameter
	}

	tmpPlaintext := make([]byte,len(cipherText))
	c, err := aes.NewCipher(key)
	if err != nil  {
		return nil,err
	}

	decrypt := cipher.NewCBCDecrypter(c, iv)

	decrypt.CryptBlocks(tmpPlaintext,cipherText)

	plaintext = tmpPlaintext

	return plaintext,nil
}
