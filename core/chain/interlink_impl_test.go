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

package chain

//import (
//	"fmt"
//	"github.com/dipperin/dipperin-core/common"
//	"github.com/dipperin/dipperin-core/core/model"
//	"github.com/dipperin/dipperin-core/third-party/log"
//	"github.com/stretchr/testify/assert"
//	"math/big"
//	"testing"
//)
//
//type testSetup struct {
//	reader  *FakeFullChainInterlink
//	genesis model.AbstractBlock
//}
//
//func setUpInterlinkTest(t *testing.T) *testSetup {
//	genesis := reader.Genesis()
//	preHash := genesis.Hash()
//	log.Info("the genesis hash is:","hash",preHash.Hex())
//	preBlock := genesis
//	//tReader := NewFakeFullChain()
//
//	for i := 1; i < 1000; i++ {
//		header := model.NewHeader(1, uint64(i), preHash, common.HexToHash("1111"), minDiff, big.NewInt(324234), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
//		block := model.NewBlockWithLink(header, nil, nil, preBlock.GetInterlinks())
//
//		err := reader.SaveBlock(block, nil)
//		assert.NoError(t, err)
//
//		preBlock = block
//		preHash = block.Hash()
//	}
//
//	return &testSetup{
//		reader:  reader,
//		genesis: genesis,
//	}
//}
//
//func TestProof(t *testing.T) {
//	setUpInterlinkTest(t)
//	p, err := GetSuffixProof(reader, 800, 2, 6)
//
//	assert.NoError(t, err)
//	assert.NoError(t, p.Valid(reader))
//
//	// for inspection
//	for _, h := range p.Prefix {
//		fmt.Println(h.Hash(), h.GetNumber())
//		b := reader.GetBlockByNumber(h.GetNumber())
//		for _, inter := range b.GetInterlinks() {
//			fmt.Println(inter)
//		}
//		fmt.Println()
//	}
//}
//
//func TestVerify(t *testing.T) {
//	setUpInterlinkTest(t)
//	reader1 := NewFakeFullChain()
//
//	preHash := TestGenesis.Hash()
//	preBlock := TestGenesis
//
//	for i := uint64(0); i < reader.CurrentBlock().Number(); i++ {
//		if i > 300 {
//			header := model.NewHeader(1, uint64(i), preHash, common.HexToHash("2222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//			block := model.NewBlockWithLink(header, nil, nil, preBlock.GetInterlinks())
//
//			err := reader1.SaveBlock(block, nil)
//			assert.NoError(t, err)
//
//			preHash = block.Hash()
//			preBlock = block
//		} else {
//			b := reader.GetBlockByNumber(i)
//			err := reader1.SaveBlock(b, b.GetVerifications())
//			assert.NoError(t, err)
//
//			preHash = b.Hash()
//			preBlock = b
//		}
//	}
//
//	m := 7
//	k := uint64(2)
//
//	p, err := GetSuffixProof(reader, 500, m, k)
//	//fmt.Println("prefix", len(p.Prefix))
//	assert.NoError(t, err)
//	assert.NoError(t, p.Valid(reader))
//
//	p1, err := GetSuffixProof(reader1, 490, m, k)
//	//p1, storageErr := GetGoodSuffixProof(reader1, 4900, m, k, 0.4)
//	//fmt.Println("prefix", len(p1.Prefix))
//
//	assert.NoError(t, err)
//	assert.NoError(t, p1.Valid(reader1))
//
//	assert.True(t, uint64(len(p.Suffix)) == k)
//	assert.True(t, uint64(len(p1.Suffix)) == k)
//
//	//assert.True(t, BlockGreater(p.Prefix, LightProofs{TestGenesis.Header()}, m))
//	//fmt.Println(BlockGreater(p.Prefix, p1.Prefix, m))
//	//fmt.Println(BlockGreater(headers, headers1, m))
//	//for inspection
//	//for _, h := range p.Prefix {
//	//	fmt.Println("Number", h.GetNumber(), h.Hash())
//	//}
//	//fmt.Println("-----------------------------")
//	//for _, h := range p1.Prefix {
//	//	fmt.Println("Number", h.GetNumber(), h.Hash())
//	//}
//}
//
//func TestGood(t *testing.T) {
//	setUpInterlinkTest(t)
//
//	m := 5
//	k := uint64(2)
//
//	p, err := GetSuffixProof(reader, 500, m, k)
//	//fmt.Println("prefix", len(p.Prefix))
//	assert.NoError(t, err)
//	assert.NoError(t, p.Valid(reader))
//
//	p1, err := GetGoodSuffixProof(reader, 500, m, k, 0.7)
//	//fmt.Println("prefix", len(p.Prefix))
//	assert.NoError(t, err)
//	assert.NoError(t, p.Valid(reader))
//
//	assert.True(t, uint64(len(p.Suffix)) == k)
//	assert.True(t, uint64(len(p1.Suffix)) == k)
//
//	//assert.True(t, BlockGreater(p.Prefix, LightProofs{TestGenesis.Header()}, m))
//	fmt.Println(BlockGreater(p.Prefix, p1.Prefix, m))
//
//	//for inspection
//	//for _, h := range p.Prefix {
//	//	fmt.Println("Number", h.GetNumber(), h.Hash())
//	//}
//	//fmt.Println("-----------------------------")
//	//for _, h := range p1.Prefix {
//	//	fmt.Println("Number", h.GetNumber(), h.Hash())
//	//}
//}
//
//func TestGoodVerify(t *testing.T) {
//	setUpInterlinkTest(t)
//	reader1 := NewFakeFullChain()
//
//	preHash := TestGenesis.Hash()
//	preBlock := TestGenesis
//
//	for i := uint64(1); i < reader.CurrentBlock().Number(); i++ {
//		if i > 300 {
//			header := model.NewHeader(1, uint64(i), preHash, common.HexToHash("2222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//			block := model.NewBlockWithLink(header, nil, nil, preBlock.GetInterlinks())
//
//			err := reader1.SaveBlock(block, nil)
//			assert.NoError(t, err)
//
//			preHash = block.Hash()
//			preBlock = block
//		} else {
//			b := reader.GetBlockByNumber(i)
//			preHash = b.Hash()
//			preBlock = b
//			err := reader1.SaveBlock(b, b.GetVerifications())
//			assert.NoError(t, err)
//		}
//	}
//
//	m := 7
//	k := uint64(2)
//	delta := float64(0.7)
//	p, err := GetGoodSuffixProof(reader, 500, m, k, delta)
//	//fmt.Println("prefix", len(p.Prefix))
//	assert.NoError(t, err)
//	assert.NoError(t, p.Valid(reader))
//
//	p1, err := GetGoodSuffixProof(reader1, 400, m, k, delta)
//	//p1, storageErr := GetGoodSuffixProof(reader1, 4900, m, k, 0.4)
//	//fmt.Println("prefix", len(p1.Prefix))
//
//	assert.NoError(t, err)
//	assert.NoError(t, p1.Valid(reader1))
//
//	assert.True(t, uint64(len(p.Suffix)) == k)
//	assert.True(t, uint64(len(p1.Suffix)) == k)
//
//	//assert.True(t, BlockGreater(p.Prefix, LightProofs{TestGenesis.Header()}, m))
//	//fmt.Println(BlockGreater(p.Prefix, p1.Prefix, m))
//	//fmt.Println(BlockGreater(headers, headers1, m))
//	//for inspection
//	//for _, h := range p.Prefix {
//	//	fmt.Println("Number", h.GetNumber(), h.Hash())
//	//}
//	//fmt.Println("-----------------------------")
//	//for _, h := range p1.Prefix {
//	//	fmt.Println("Number", h.GetNumber(), h.Hash())
//	//}
//}
//
//func TestAfterTarget(t *testing.T) {
//	//header1 := model.NewHeader(1, uint64(1), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	header2 := model.NewHeader(1, uint64(2), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	header3 := model.NewHeader(1, uint64(3), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	//header4 := model.NewHeader(1, uint64(4), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	header5 := model.NewHeader(1, uint64(5), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	header6 := model.NewHeader(1, uint64(6), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	header7 := model.NewHeader(1, uint64(7), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	//header8 := model.NewHeader(1, uint64(8), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//	//header9 := model.NewHeader(1, uint64(9), TestGenesis.Hash(), common.HexToHash("22222"), minDiff, big.NewInt(1232131), common.HexToAddress("00ff1acd"), common.BlockNonceFromInt(23322))
//
//	headers := LightProofs{NewLightProof(header2, nil), NewLightProof(header3, nil), NewLightProof(header6, nil), NewLightProof(header7, nil)}
//
//	fmt.Println(AfterTarget(headers, header5))
//}
//
//func TestGetInfixProof(t *testing.T) {
//	setUpInterlinkTest(t)
//
//	m := 2
//	k := uint64(2)
//	tar := uint64(300)
//
//	p, err := GetInfixProof(reader, 800, m, k, tar)
//	assert.NoError(t, err)
//	//for _, a := range p.Prefix {
//	//	fmt.Println(a.Hash(), a.GetNumber())
//	//}
//	assert.NoError(t, p.Valid(reader))
//}
//
//func TestVerifyInfix(t *testing.T) {
//	setUpInterlinkTest(t)
//
//	m := 2
//	k := uint64(2)
//	tar := uint64(300)
//
//	p, err := GetInfixProof(reader, 800, m, k, tar)
//	assert.NoError(t, err)
//	//fmt.Println(p.Prefix)
//	assert.NoError(t, p.Valid(reader))
//	D := func(proofs LightProofs) bool {
//		//for _, a := range proofs {
//		//	fmt.Println(a.Hash(), a.GetNumber())
//		//}
//
//		return true
//	}
//
//	VerifyInfix(Proofs{*p}, k, m, D)
//}
//
//func TestVerifySuffix(t *testing.T) {
//	setUpInterlinkTest(t)
//	m := 4
//	k := uint64(6)
//	tar := uint64(300)
//
//	p, err := GetInfixProof(reader, 800, m, k, tar)
//	assert.NoError(t, err)
//
//	Q := func(proofs LightProofs) bool {
//		//for _, a := range proofs {
//		//	fmt.Println(a.Hash(), a.GetNumber())
//		//}
//
//		return true
//	}
//
//	VerifySuffix(Proofs{*p}, k, m, Q)
//}
