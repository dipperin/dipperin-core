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


package model

import (
	"bytes"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/rlp"
)

//var EmptyRoot = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

type DerivableList interface {
	Len() int
	GetKey(i int) []byte
	GetRlp(i int) []byte
}

type DeriveShaInterface interface {
	DeriveSha(list DerivableList) common.Hash
}

func DeriveSha(list DerivableList) common.Hash {
	tree := new(trie.Trie)
	for i := 0; i < list.Len(); i++ {
		tree.Update(list.GetKey(i), list.GetRlp(i))
	}
	return tree.Hash()
}

// merkle root for AbstractVerification
type Verifications []AbstractVerification

func (verf Verifications) GetKey(i int) []byte {
	keybuf := new(bytes.Buffer)
	rlp.Encode(keybuf, uint(i))
	res := keybuf.Bytes()
	return res

}

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (verf Verifications) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(verf[i])
	return enc
}

// Len returns the length of s.
func (verf Verifications) Len() int { return len(verf) }


// merkle root for AbstractTransaction
type AbsTransactions []AbstractTransaction

func (verf AbsTransactions) GetKey(i int) []byte {
	res := verf[i].CalTxId().Bytes()
	return res

}

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (verf AbsTransactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(verf[i])
	return enc
}

// Len returns the length of s.
func (verf AbsTransactions) Len() int { return len(verf) }