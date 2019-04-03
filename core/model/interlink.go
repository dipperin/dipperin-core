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
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type InterLink []common.Hash

//var maxhash = common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

func HashLevel(hash common.Hash, maxHash common.Hash) int {
	if hash.Cmp(common.Hash{}) == 0 {
		return 0
	}
	div := big.NewInt(0).Div(maxHash.Big(), hash.Big())
	count := 0
	for div.Cmp(big.NewInt(0)) != 0 {
		div.Div(div, big.NewInt(2))
		count++
	}
	return count
}

//assuming genesis blocks level is 0,other level,relative to the paper description, move one layer as a whole.
func NewInterLink(preBlockLinks InterLink, curBlock AbstractBlock) InterLink {
	if curBlock.Number() == uint64(1) {
		// assuming genesis block has 0 level
		linkList := []common.Hash{curBlock.PreHash(), curBlock.PreHash()}
		return linkList
	}

	bound := len(preBlockLinks)
	newLinks := make([]common.Hash, bound)
	copy(newLinks, preBlockLinks)

	maxLevel := 0
	if !curBlock.IsSpecial() {
		maxLevel = HashLevel(curBlock.PreHash(), curBlock.Difficulty().DiffToTarget())
	}

	//fmt.Println("maxLevel", maxLevel)
	for i := 1; i <= maxLevel; i ++ {
		if i < bound {
			newLinks[i] = curBlock.PreHash()
		} else {
			newLinks = append(newLinks, curBlock.PreHash())
		}
	}
	//fmt.Println(newLinks)
	return newLinks
}

func (l InterLink) GetKey(i int) []byte {
	keybuf := new(bytes.Buffer)
	rlp.Encode(keybuf, uint(i))
	res := keybuf.Bytes()
	return res

}

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (l InterLink) GetRlp(i int) []byte {
	//enc, _ := rlp.EncodeToBytes(l[i])
	return l[i].Bytes()
}

// Len returns the length of s.
func (l InterLink) Len() int { return len(l) }