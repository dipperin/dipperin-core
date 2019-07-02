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

package common

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash_Func(t *testing.T) {

	testhash1 := RlpHashKeccak256("Hello World!")
	testhash2 := RlpHashKeccak256("Hello World!")
	//testhash1 is initially equal to testhash2
	assert.Equal(t, true, testhash1.IsEqual(testhash2))
	//Testhash1 is emptied, check if the emptying is successful
	testhash1.Clear()
	if !testhash1.IsEqual(Hash{}) {
		t.Errorf("Hash clear does not work")
	}
	//Testhash1 is emptied and should not equal testhash2
	if testhash1.IsEqual(testhash2) {
		t.Errorf("Address1 should be empty and not equal to the old one")
	}
}

func TestAddress_Func(t *testing.T) {
	testaddr1 := HexToAddress("0x000F9328D55ccb3FCe531f199382339f0E576ee840B1")
	testaddr2 := HexToAddress("0x000F9328D55ccb3FCe531f199382339f0E576ee840B1")
	if !testaddr1.IsEqual(testaddr2) {
		t.Errorf("Initial address should be equal")
	}
	testaddr1.Clear()
	if !testaddr1.IsEqual(Address{}) {
		t.Errorf("Address clear does not work")
	}
	if testaddr1.IsEqual(testaddr2) {
		t.Errorf("Address1 should be empty and not equal to the old one ")
	}

}

func TestBlockNonce_FromInt(t *testing.T) {
	assert.Equal(t, "0x7274152500000000000000000000000000000000000000000000000000000000", BlockNonceFromInt(1920210213).Hex())
}

func BenchmarkBlockNonceFromInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if BlockNonceFromInt(1920210213).Hex() != "0x72741525" {
			b.Fatal("panic")
		}
	}
}

func TestDifficulty_DiffToTarget(t *testing.T) {
	testdiff := HexToDiff("0x1743eca9")
	////testdiff := HexToDiff("0x2043eca9")
	assert.Equal(t, HexToHash("0x00000000000000000043eca90000000000000000000000000000000000000000"), testdiff.DiffToTarget())
	//testdiff := HexToDiff("0x1d00ffff")
	//fmt.Println(testdiff.DiffToTarget().Hex())
	testdiff2 := HexToDiff("0x172a4e2f")
	testdiff3 := HexToDiff("0x1d00ffff")
	a := big.NewInt(0).Div(testdiff3.Big(), testdiff2.Big())
	assert.Equal(t, big.NewInt(6653303141405), a)
	b := big.NewInt(0).Mul(a, big.NewInt(1<<32))
	c := big.NewInt(0).Div(b, big.NewInt(600))
	d, _ := new(big.Int).SetString("47626199004514230818", 10)
	assert.Equal(t, d, c)
}

func BenchmarkDifficulty_DiffToTarget(b *testing.B) {
	testdiff := HexToDiff("0x1743eca9")
	testhash := HexToHash("0x00000000000000000043eca90000000000000000000000000000000000000000")
	for i := 0; i < b.N; i++ {
		if testhash.Cmp(testdiff.DiffToTarget()) != 0 {
			b.Fatal("panic")
		}
	}
}

func TestHash_Cmp(t *testing.T) {
	testhash1 := RlpHashKeccak256("Hello World!")
	testhash2 := RlpHashKeccak256("Hello World?")
	assert.Equal(t, 1, testhash1.Cmp(testhash2))
	assert.Equal(t, 0, testhash1.Cmp(testhash1))
	assert.Equal(t, -1, testhash2.Cmp(testhash1))
}

func BenchmarkHash_Cmp(b *testing.B) {
	testhash1 := RlpHashKeccak256("Hello World!")
	testhash2 := RlpHashKeccak256("Hello World?")
	for i := 0; i < b.N; i++ {
		if testhash1.Cmp(testhash2) != 1 || testhash2.Cmp(testhash1) != -1 {
			b.Fatal("panic")
		}
	}

}

func TestCopyHash(t *testing.T) {
	a := HexToHash("11")
	b := CopyHash(&a)
	fmt.Println(a)
	fmt.Println(b)
	z := HexToHash("2")
	b = &z
	fmt.Println(a)
	fmt.Println(b)
}

func TestBigToDiff(t *testing.T) {
	//diff1:=HexToDiff("0x1d001fff")
	//fmt.Println(diff1.DiffToTarget().Hex())
	//fmt.Println(diff1.Big())
	//fmt.Println(BigToDiff(diff1.Big()).Hex())
	mainPowLimit := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 253), big.NewInt(1))
	fmt.Println(mainPowLimit)
	diff := BigToDiff(mainPowLimit)
	fmt.Println(diff.Hex())
	fmt.Println(diff.DiffToTarget().Hex())
}

//func BytePlusOne(array []byte)(res []byte){
//	var index = len(array)- 1
//	for index >= 0 {
//		if array[index] < 255 {
//			array[index]++
//			break
//		} else {
//			array[index] = 0
//			index--
//		}
//	}
//	res=array
//	return
//}
//
//func TestRlpHashKeccak256(t *testing.T) {
//     testbyte:=[24]byte{}
//      a:=BytePlusOne(testbyte[:])
//      for i := 0; i < 10001; i++ {
//		  a=BytePlusOne(testbyte[:])
//	  }
//      fmt.Println(a)
//     fmt.Println(big.NewInt(1000).Bytes())
//}

func TestTypeXX(t *testing.T) {
	x := AddressTypeNormal
	assert.Equal(t, "normal transaction", (TxType)(x).String())
	x = AddressTypeCross
	assert.Equal(t, "cross chain transaction", (TxType)(x).String())
	x = AddressTypeStake
	assert.Equal(t, "stake transaction", (TxType)(x).String())
	x = AddressTypeCancel
	assert.Equal(t, "cancel transaction", (TxType)(x).String())
	x = AddressTypeUnStake
	assert.Equal(t, "unstake transaction", (TxType)(x).String())
	x = AddressTypeEvidence
	assert.Equal(t, "evidence transaction", (TxType)(x).String())
	x = AddressTypeERC20
	assert.Equal(t, "erc20 transaction", (TxType)(x).String())
	x = 0x999
	assert.Contains(t, (TxType)(x).String(), "unkonw")

	h := &Hash{}
	assert.True(t, CopyHash(h).IsEmpty())

	h1 := Hash{0x12}
	assert.Equal(t, `"0x1200000000000000000000000000000000000000000000000000000000000000"`, util.StringifyJson(h1))
	assert.NoError(t, util.ParseJson(`"0x1200000000000000000000000000000000000000000000000000000000000000"`, &h1))

	xq, err := hexutil.Decode("0x1200000000000000000000000000000000000000000000000000000000000000123f")
	assert.NoError(t, err)
	h.SetBytes(xq)

	assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000123f", h.HexWithout0x())
	h.Str()
	h.Bytes()
	assert.Equal(t, Hash{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30}, StringToHash("0x1200000000000000000000000000000000000000000000000000000000000000"))
	assert.Equal(t, Hash{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc}, BigToHash(big.NewInt(12)))

	assert.True(t, h.ValidHashForDifficulty(HexToDiff("0x1e566611")))
	assert.False(t, Hash{0xff}.ValidHashForDifficulty(HexToDiff("0x1e566611")))
}

func TestAddressFuncs(t *testing.T) {
	addr := Address{0x12}
	assert.Equal(t, "0x12000000000000000000000000000000000000000000", addr.String())

	util.StringifyJson(addr)
	assert.NoError(t, util.ParseJson(`"0x123f"`, &addr))
	addr.GetAddressType()
	addr.GetAddressTypeStr()
	assert.False(t, addr.IsEmpty())
	assert.True(t, Address{}.IsEmpty())
	addr.IsEqualWithoutType(HexToAddress("0x123"))
	addr.Str()
	addr.Big()
	addr.Bytes()
	addr.Hash()
	Address{0x12}.InSlice([]Address{{0x12}})
	Address{0x12}.InSlice([]Address{{}})
	addr.SetBytes([]byte{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x11, 0x12})

	assert.Equal(t, "Normal", HexToAddress("0x00005033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())
	assert.Equal(t, consts.ERC20TypeName, HexToAddress("0x00105033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())
	assert.Equal(t, "CrossChain", HexToAddress("0x00015033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())
	assert.Equal(t, "Stake", HexToAddress("0x00025033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())
	assert.Equal(t, "Cancel", HexToAddress("0x00035033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())
	assert.Equal(t, "UnStake", HexToAddress("0x00045033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())
	assert.Equal(t, "Evidence", HexToAddress("0x00055033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())
	assert.Equal(t, consts.EarlyTokenTypeName, HexToAddress("0x00115033874289F4F823A896700D94274683535cF0E1").GetAddressTypeStr())

	assert.False(t, StringToAddress("0x00012").IsEmpty())
	assert.False(t, BigToAddress(big.NewInt(11)).IsEmpty())

}

func TestDiffFuncs(t *testing.T) {
	diff := Difficulty{0x12}
	var diff2 Difficulty
	assert.NoError(t, util.ParseJson(util.StringifyJson(diff), &diff2))
	assert.True(t, diff.Equal(diff2))

	diff.SetBytes([]byte{0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x11, 0x12})
	assert.Equal(t, Difficulty{0x30, 0x30, 0x11, 0x12}, diff)

	assert.Equal(t, "12000000", diff2.HexWithout0x())
	assert.Equal(t, "\x12\x00\x00\x00", diff2.Str())
	assert.Equal(t, []byte{0x12, 0x0, 0x0, 0x0}, diff2.Bytes())

	diff3 := Difficulty{}
	assert.Panics(t, func() {
		diff3.DiffToTarget()
	})
}

func TestBlockNonceFuncs(t *testing.T) {
	n1 := BlockNonce{0x1}
	var n2 BlockNonce
	assert.NoError(t, util.ParseJson(util.StringifyJson(n1), &n2))
	assert.True(t, n1.IsEqual(n2))

	assert.Equal(t, BlockNonce{0x12, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, BlockNonceFromHex("0x12"))
	assert.Equal(t, BlockNonce{0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, EncodeNonce(1))
}
