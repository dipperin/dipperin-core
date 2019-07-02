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

package factory

import (
	"encoding/hex"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"testing"
)

func TestNewFBlock(t *testing.T) {
	//// new factory block
	//mockBlock := NewFBlock(100, common.HexToHash("perBlockHash"), common.HexToDiff("1fffffff"),
	//	common.StringToAddress("coinbase")).CreateBlock()
	//
	//assert.NotNil(t, mockBlock)
	//assert.Equal(t, uint64(100), mockBlock.Number())
}

func TestNewFTx(t *testing.T) {
	//fTx := NewFTx("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031", 0)
	//
	//txs := fTx.CreateTx(common.StringToAddress("testAddress"), big.NewInt(10), big.NewInt(3), nil).
	//	CreateTx(common.StringToAddress("testAddress"), big.NewInt(10), big.NewInt(3), nil).
	//	GetTxs()
	//
	//assert.Equal(t, 2, len(txs))
}

func TestNewBlock(t *testing.T) {
	//fTx := NewFTx("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031", 0)
	//txs := fTx.CreateTx(common.StringToAddress("testAddress"), big.NewInt(10), big.NewInt(3), nil).
	//	CreateTx(common.StringToAddress("testAddress"), big.NewInt(10), big.NewInt(3), nil).
	//	GetTxs()
	//
	//mockBlock := NewFBlock(100, common.HexToHash("perBlockHash"), common.HexToDiff("1fffffff"),
	//	common.StringToAddress("coinbase")).
	//	AddTxs(txs).
	//	CreateBlock()
	//
	//assert.Equal(t, 2, len(mockBlock.GetTransactions()))
}

func TestFactoryGenPrk(t *testing.T) {
	prk := factoryGenPrk()

	println(hex.EncodeToString(crypto.FromECDSA(prk)))
}

func TestNewFAddress(t *testing.T) {
	fa := NewFAddress(50)

	println(fa.GetAddress(49).Hex())
}
