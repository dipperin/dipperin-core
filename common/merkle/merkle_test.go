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

package merkle

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestTxNum           = 19
	TestTxIndex         = 9
	TestTxMutatedIndex1 = 7
	TestTxMutatedIndex2 = 8
)

type b struct {
	slice1 []byte
	map1   map[string]byte
	array1 [16]byte
}
type a struct {
	test1 b
	test2 []byte
}

//self test
func TestNew(t *testing.T) {
	//testdata := new(a)

	//cslog.Printf("the testdata is: %v",testdata)
	//cslog.Printf("the test1 slice1 is: %p",testdata.test1.slice1)
	//testdata.test1.slice1[:] = []byte{1,2,3}
	//testdata.test1.map1["key"] = 1
	//testda
	//cslog.Printf("the test2 address is: 0x%x",testdata.test2)
}

//self test slice and array
func TestArraySlice(t *testing.T) {
	var testArray [5]int

	testSlice := make([]int, 0)
	testSlice = []int{1, 2, 3, 4, 5, 6, 7}

	//cslog.Printf("the testSlice is:%v the test Array is:%v",testSlice,testArray)

	//a:=testArray[:]

	//cslog.Printf("the testArray is:%v the testArray type is:%v",a,reflect.TypeOf(a))

	copy(testArray[:], testSlice)

	//cslog.Printf("the testArray is:%v the testArray type is:%v",testArray,reflect.TypeOf(testArray))
}

// getTxDataHash get the transaction txdata hash value
func getTxDataHash(TestMutateFlag bool) []common.Hash {
	// construct transaction hash data
	var txNum uint = TestTxNum

	tmpByte := make([]byte, 32)
	TxData := make([]common.Hash, txNum)

	//cslog.Debug().Msg("the TxData hash is: ")
	for i := range TxData {
		// Test the same conflict situation of two transactions hash
		if TestMutateFlag == true {
			if (i == TestTxMutatedIndex1) || (i == TestTxMutatedIndex2) {
				tmpByte[0] = TestTxMutatedIndex1
			}
		} else {
			tmpByte[0] = byte(i)
		}

		TxData[i] = cs_crypto.Keccak256Hash(tmpByte)
		//cslog.Printf("%v: %v ",i,TxData[i].Hex())
	}

	return TxData
}

func testComputeBranch() (pRoot common.Hash, pBranch []common.Hash) {

	TxData := getTxDataHash(false)

	// Calculate the Merkel root value and the accompanying path of the transaction based on the given transaction data hash value, transaction index
	var index uint32 = TestTxIndex
	var mutated = false
	pRoot, pBranch = MerkleComputation(TxData, index, &mutated)
	//cslog.Printf("the pRoot is: %v",pRoot.Hex())
	//cslog.Debug().Msg("the pBranch is: ")
	//for j,tmpHash := range pBranch{
	//	cslog.Printf("%v: %v ",j,tmpHash.Hex())
	//}
	return pRoot, pBranch
}

func TestComputeMerkleBranch(t *testing.T) {
	testComputeBranch()

	TxData := getTxDataHash(false)
	// Calculate the Merkel root value and the accompanying path of the transaction based on the given transaction data hash value, transaction index
	var index uint32 = TestTxIndex
	var mutated = false
	pr, _ := MerkleComputation(TxData, index, nil)
	assert.True(t, pr.IsEmpty())
	pr, _ = MerkleComputation([]common.Hash{}, index, &mutated)
	assert.True(t, pr.IsEmpty())
	var testH []common.Hash
	for i := 0; i < MaxRouteNumber + 2; i++ {
		testH = append(testH, common.Hash{})
	}
	pr, _ = MerkleComputation(testH, index, &mutated)
	assert.True(t, pr.IsEmpty())

	pr, _ = MerkleComputation(TxData, uint32(len(TxData) + 2), &mutated)
	assert.True(t, pr.IsEmpty())
}

func TestComputeMerkleRootFromBranch(t *testing.T) {
	//var index uint32 = TestTxIndex
	//
	//TxData := getTxDataHash(false)
	//
	//pRoot,pBranch := testComputeBranch()

	//according to pbranch calculation merkelgen
	//pRoot2:=ComputeMerkleRootFromBranch(pBranch,TxData[index],index)

	//if pRoot != pRoot2{
	//	cslog.Debug().Msg("the calculated RootHash is different")
	//	cslog.Printf("the pRoot is: %v",pRoot)
	//	cslog.Printf("the pRoot2 is: %v",pRoot2)
	//}else{
	//	cslog.Debug().Msg("test ComputeMerkleRootFromBranch OK")
	//}
}

func TestBlockMerkleRoot(t *testing.T) {
	var pMutated = false
	TxData := getTxDataHash(false)
	ComputeMerkleRoot(TxData, &pMutated)

	//if pMutated != false{
	//	cslog.Debug().Msg("test normal BlockMerkleRoo fail")
	//}else{
	//	cslog.Debug().Msg("test normal BlockMerkleRoo success")
	//}
	//
	//TxData = getTxDataHash(true)
	//ComputeMerkleRoot(TxData,&pMutated)
	//
	//if pMutated != true{
	//	cslog.Debug().Msg("test mutated BlockMerkleRoo fail")
	//}else{
	//	cslog.Debug().Msg("test mutated BlockMerkleRoo success")
	//}

}
