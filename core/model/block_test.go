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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestBlock_GetTransactionFees(t *testing.T) {
	block := CreateBlock(1, common.Hash{}, 10)
	fee := block.GetTransactionFees()
	assert.Equal(t, big.NewInt(210000), fee)
}

func TestBlock_EncodeToIBLT(t *testing.T) {
	block := CreateBlock(1, common.Hash{}, 150)
	txsMap := make(map[common.Hash]struct{})
	for _, tx := range block.GetTransactions() {
		txsMap[tx.CalTxId()] = struct{}{}
	}

	bloom := block.EncodeToIBLT()
	recovered, _, err := bloom.ListRLP()
	assert.NoError(t, err)

	var tx Transaction
	assert.Equal(t, len(txsMap), len(recovered))

	for _, re := range recovered {
		err = rlp.DecodeBytes(re, &tx)
		assert.NoError(t, err)
		assert.NotNil(t, txsMap[tx.CalTxId()])
	}
}

func TestBlock_BloomFilter(t *testing.T) {
	header := NewHeader(1, 100, common.HexToHash("1111"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(324234), aliceAddr, common.BlockNonceFromInt(432423))

	var (
		aliceTxs  = 50
		bobTxs    = 500
		commonTxs = 2000
	)
	txs := CreateSignedTxList(aliceTxs + bobTxs + commonTxs)
	aTxs := txs[:aliceTxs]
	bTxs := txs[aliceTxs : bobTxs+aliceTxs]
	cTxs := txs[bobTxs+aliceTxs:]

	// alice pool
	blockTxs := append(aTxs, cTxs...)
	block := NewBlock(header, blockTxs, nil)
	bloom := block.EncodeToIBLT()
	rootSent := DeriveSha(Transactions(blockTxs))

	// bob pool
	bPool := append(bTxs, cTxs...)
	bPoolMap := make(map[common.Hash]*Transaction)
	for _, tx := range bPool {
		bPoolMap[tx.CalTxId()] = tx
	}

	blockTxMap := make(map[common.Hash]*Transaction)
	for _, tx := range blockTxs {
		blockTxMap[tx.CalTxId()] = tx
	}

	recovered, err := bloom.Recover(bPoolMap)
	assert.NoError(t, err)
	assert.Equal(t, len(blockTxs), len(recovered))
	assert.Equal(t, aliceTxs+commonTxs, len(blockTxs))

	var txsRecovered []*Transaction
	for _, re := range recovered {
		var trans Transaction

		err = rlp.DecodeBytes(re, &trans)
		assert.NoError(t, err)

		assert.NotNil(t, blockTxMap[trans.CalTxId()])

		txsRecovered = append(txsRecovered, &trans)
	}

	rootGet := DeriveSha(Transactions(txsRecovered))
	assert.Equal(t, rootSent, rootGet)
}

// test rlpHash
func Test_RlpHash(t *testing.T) {
	// calculate TxId = hash(tx.data+address)
	testNonce := 4
	from := common.HexToAddress("0x0000d28Eb0154A96F4af6E631766939593554c7E5577")
	to := common.HexToAddress("0x0000E21391AA1ccAcb7c8E7E2E645Bb2cF811fe1E30D")
	value := big.NewInt(100000000000000)
	log.Info("the value is:", "value", hexutil.Encode(value.Bytes()))

	transactionFee := big.NewInt(100000000)
	log.Info("the transactionFee is:", "transactionFee", hexutil.Encode(transactionFee.Bytes()))

	data := make([]byte, 0)
	tx := NewTransaction(uint64(testNonce), to, value, g_testData.TestGasPrice, g_testData.TestGasLimit, data)
	txId, err := rlpHash([]interface{}{tx.data, from})
	assert.NoError(t, err)
	log.Debug("the txId is :", "txId", txId.Hex())

	/*encodeBytes ,err:= rlp.EncodeToBytes(tx.data)
	assert.NoError(t,err)
	log.Debug("the encodeBytes is :","encodeBytes",hexutil.Encode(encodeBytes))*/

	hashSrcData, err := hexutil.Decode("0xf83fe704960000e21391aa1ccacb7c8e7e2e645bb2cf811fe1e30d8080865af3107a40008405f5e10080960000d28eb0154a96f4af6e631766939593554c7e5577")
	assert.NoError(t, err)

	txId2 := cs_crypto.Keccak256Hash(hashSrcData)
	log.Debug("the txId2 is :", "txId2", txId2.Hex())
}

func TestRunWorkMap(t *testing.T) {
	txs1 := CreateSignedTxList(4000)
	reqEstimator := iblt.NewHybridEstimator(iblt.NewHybridEstimatorConfig())
	estimator := iblt.NewHybridEstimator(reqEstimator.Config())
	mapEstimator := newMapWorkHybridEstimator(estimator)

	start := time.Now()
	for _, tx := range txs1 {
		estimator.EncodeByte(tx.CalTxId().Bytes())
	}
	fmt.Println(time.Now().Sub(start))

	txs2 := CreateSignedTxList(4000)
	start = time.Now()
	err := RunWorkMap(mapEstimator, txs2)
	assert.NoError(t, err)
	fmt.Println(time.Now().Sub(start))

	estimatedConfig := estimator.DeriveConfig(reqEstimator)
	bloomConfig := iblt.DeriveBloomConfig(len(txs1))
	invBloom := iblt.NewGraphene(estimatedConfig, bloomConfig)
	mapBloom := newMapWorkInvBloom(invBloom)

	start = time.Now()
	for _, tx := range txs1 {
		invBloom.InsertRLP(tx.CalTxId(), tx)
		invBloom.Bloom().Digest(tx.CalTxId().Bytes())
	}
	fmt.Println(time.Now().Sub(start))

	start = time.Now()
	err = RunWorkMap(mapBloom, txs2)
	assert.NoError(t, err)
	fmt.Println(time.Now().Sub(start))
}

func TestHeader_IsEqual(t *testing.T) {
	header := newTestHeader()
	assert.Panics(t, func() {
		header.IsEqual(header)
	})
}

func TestHeader_CoinBaseAddress(t *testing.T) {
	header := newTestHeader()
	assert.Equal(t, header.CoinBase, header.CoinBaseAddress())
}

func TestHeader_GetStateRoot(t *testing.T) {
	block := CreateBlock(1, common.Hash{}, 10000)
	block.SetStateRoot(common.HexToHash("123"))

	stateRoot := block.header.GetStateRoot()
	assert.Equal(t, common.HexToHash("123"), stateRoot)
}

func newTestHeader() (h *Header) {
	h = NewHeader(1, 100, common.HexToHash("001010001010"), common.HexToHash("10111101011"), common.HexToDiff("1111ffff"), big.NewInt(10100), aliceAddr, common.BlockNonceFromInt(100))
	return
}

func TestNewHeader(t *testing.T) {
	h := newTestHeader()
	assert.NotNil(t, h)
	assert.Equal(t, h.Version, uint64(1), h.Number, uint64(100))
}

func TestCopyHeader(t *testing.T) {
	h := newTestHeader()
	ch := CopyHeader(h)
	assert.NotNil(t, ch)
}

func Test_writeCounter_Write(t *testing.T) {
	counter := writeCounter(100)
	result, err := counter.Write([]byte{123})
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestHeader_RlpBlockWithoutNonce(t *testing.T) {
	h := newTestHeader()
	rb := h.RlpBlockWithoutNonce()
	assert.NotNil(t, rb)
}

func TestHeader_Hash(t *testing.T) {
	h := newTestHeader()
	hash := h.Hash()
	assert.NotNil(t, hash)
}

func TestHeader_GetNumber(t *testing.T) {
	h := newTestHeader()
	num := h.GetNumber()
	assert.Equal(t, h.Number, num)
}

func TestHeader_GetSeed(t *testing.T) {
	h := newTestHeader()
	seed := h.GetSeed()
	assert.Equal(t, h.Seed, seed)
}

func TestHeader_GetProof(t *testing.T) {
	h := newTestHeader()
	result := h.GetProof()
	assert.Equal(t, h.Proof, result)
}

func TestHeader_GetMinerPubKey(t *testing.T) {
	h := newTestHeader()
	result := h.GetMinerPubKey()

	var key *ecdsa.PublicKey
	assert.Equal(t, key, result)
}

func TestHeader_GetPreHash(t *testing.T) {
	h := newTestHeader()
	result := h.GetPreHash()
	assert.Equal(t, h.PreHash, result)
}

func TestHeader_GetInterLinkRoot(t *testing.T) {
	h := newTestHeader()
	result := h.GetInterLinkRoot()
	assert.Equal(t, h.InterlinkRoot, result)
}

func TestHeader_GetDifficulty(t *testing.T) {
	h := newTestHeader()
	result := h.GetDifficulty()
	assert.Equal(t, h.Diff, result)
}

func TestHeader_GetRegisterRoot(t *testing.T) {
	h := newTestHeader()
	result := h.GetRegisterRoot()
	assert.Equal(t, h.RegisterRoot, result)
}

func TestHeader_SetRegisterRoot(t *testing.T) {
	h := newTestHeader()
	h.SetRegisterRoot(common.HexToHash("123"))
	assert.Equal(t, common.HexToHash("123"), h.RegisterRoot)
	result := h.GetRegisterRoot()
	assert.Equal(t, h.RegisterRoot, result)
}

func TestHeader_DuplicateHeader(t *testing.T) {
	h := newTestHeader()
	c := h.DuplicateHeader()
	assert.Equal(t, h, c)
}

func TestHeader_SetVerificationRoot(t *testing.T) {
	h := newTestHeader()
	h.SetVerificationRoot(common.HexToHash("123"))
	assert.Equal(t, common.HexToHash("123"), h.VerificationRoot)
}

func TestHeader_Size(t *testing.T) {
	h := newTestHeader()
	result := h.Size()
	assert.NotEqual(t, 0, result)
}

func TestHeader_HashWithoutNonce(t *testing.T) {
	h := newTestHeader()
	result := h.HashWithoutNonce()
	assert.NotNil(t, result)
}

func Test_rlpHash(t *testing.T) {
	h := newTestHeader()
	rlp, err := rlpHash(h)
	assert.NoError(t, err)
	assert.NotNil(t, rlp)
}

func TestHeader_String(t *testing.T) {
	h := newTestHeader()
	str := h.String()
	fmt.Println(str)
	assert.NotEqual(t, 0, len(str))
}

func newTestBody() (b *Body) {
	b = &Body{
		[]*Transaction{CreateSignedTx(0, big.NewInt(10000))},
		nil,
		[]common.Hash{common.HexToHash("123")},
	}
	return
}

func TestBody_GetTxsSize(t *testing.T) {
	b := newTestBody()
	result := b.GetTxsSize()
	assert.Equal(t, len(b.Txs), result)
}

func TestBody_GetTxByIndex(t *testing.T) {
	b := newTestBody()
	result := b.GetTxByIndex(0)
	assert.Equal(t, b.Txs[0], result)
}

func TestBody_GetInterLinks(t *testing.T) {
	b := newTestBody()
	result := b.GetInterLinks()
	assert.Equal(t, b.Inters, result)
}

func TestBody_EncodeRlpToBytes(t *testing.T) {
	b := newTestBody()
	_, err := b.EncodeRlpToBytes()
	assert.NoError(t, err)
}

func TestBlock_IsSpecial(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	assert.Equal(t, false, block.IsSpecial())

	block.SetDifficulty(common.Difficulty{})
	block.SetNonce(common.BlockNonce{})
	block.SetTimeStamp(big.NewInt(100))
	assert.Equal(t, true, block.IsSpecial())
}

func TestBlock_GetRegisterRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.GetRegisterRoot()
	assert.Equal(t, block.header.RegisterRoot, result)
}

func TestBlock_SetRegisterRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetRegisterRoot(common.HexToHash("000000000001"))
	result := block.GetRegisterRoot()
	assert.Equal(t, block.header.RegisterRoot, result)
}

func TestBlock_GetBlockTxsBloom(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.GetBlockTxsBloom()
	assert.NotNil(t, result)
}

func TestBlock_getBlockInvBloom(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	c := iblt.NewHybridEstimatorConfig()
	estimator := iblt.NewHybridEstimator(c)

	result := block.getBlockInvBloom(estimator)
	assert.NotNil(t, result)
}

func TestBlock_GetEiBloomBlockData(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 2)
	c := iblt.NewHybridEstimatorConfig()
	estimator := iblt.NewHybridEstimator(c)

	result := block.GetEiBloomBlockData(estimator)
	assert.NotNil(t, result)

	DefaultTxs = 0
	result = block.GetEiBloomBlockData(estimator)
	assert.NotNil(t, result)
}

func TestBlock_SetVerifications(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetVerifications([]AbstractVerification{})
	assert.NotNil(t, block.body.Vers)
}

func TestBlock_GetVerifications(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetVerifications([]AbstractVerification{})
	result := block.GetVerifications()
	assert.NotNil(t, result)
}

func TestBlock_GetTransactions(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)

	result := block.GetTransactions()
	assert.Len(t, result, 1)
}

func TestBlock_GetInterlinks(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)

	result := block.GetInterlinks()
	assert.NotNil(t, result)
}

func TestBlock_GetAbsTransactions(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 2)

	result := block.GetAbsTransactions()
	assert.Len(t, result, 2)
}

func TestBlock_SetNonce(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetNonce(common.BlockNonceFromHex("85"))
	assert.Equal(t, common.BlockNonceFromHex("85"), block.header.Nonce)
}

func TestBlock_FormatForRpc(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.FormatForRpc()
	assert.Equal(t, nil, result)
}

func TestBlock_SetStateRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetStateRoot(common.HexToHash("01010101010101010101"))
	assert.Equal(t, common.HexToHash("01010101010101010101"), block.header.StateRoot)
}

func TestBlock_SetInterLinkRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetInterLinkRoot(common.HexToHash("01010101010101010101"))
	assert.Equal(t, common.HexToHash("01010101010101010101"), block.header.InterlinkRoot)
}

func TestBlock_GetInterLinkRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetInterLinkRoot(common.HexToHash("01010101010101010101"))
	result := block.GetInterLinkRoot()
	assert.Equal(t, block.header.InterlinkRoot, result)
}

func TestBlock_SetInterLinks(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetInterLinks(InterLink{})
	assert.Equal(t, InterLink{}, block.body.Inters)
}

func TestBlock_TxIterator(t *testing.T) {
	block := CreateBlock(0, common.HexToHash("123"), 1)
	err := block.TxIterator(func(i int, transaction AbstractTransaction) error {
		assert.Equal(t, uint64(0), transaction.Nonce())
		return nil
	})
	assert.NoError(t, err)

	err = block.TxIterator(func(i int, transaction AbstractTransaction) error {
		assert.Equal(t, uint64(0), transaction.Nonce())
		return errors.New("iterator error")
	})
	assert.Equal(t, "iterator error", err.Error())
}

func TestBlock_VersIterator(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	vote := CreateSignedVote(1, 0, block.Hash(), VoteMessage)
	block.SetVerifications([]AbstractVerification{vote})

	err := block.VersIterator(func(i int, verification AbstractVerification, block AbstractBlock) error {
		return nil
	})
	assert.NoError(t, err)

	err = block.VersIterator(func(i int, verification AbstractVerification, block AbstractBlock) error {
		return errors.New("iterator error")
	})
	assert.Equal(t, "iterator error", err.Error())
}

func TestBlock_GetCoinbase(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.GetCoinbase()
	assert.NotEqual(t, nil, result)
}

func TestNewBlock(t *testing.T) {
	vote := CreateSignedVote(1, 0, common.HexToHash("123"), VoteMessage)
	block := NewBlock(&Header{}, nil, []AbstractVerification{vote})
	assert.NotNil(t, block)
}

func TestNewBlockWithLink(t *testing.T) {
	h := newTestHeader()
	txs := []*Transaction{CreateSignedTx(0, big.NewInt(10000))}
	block := NewBlockWithLink(h, txs, nil, nil)
	assert.NotNil(t, block)
}

func TestBlock_String(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	str := block.String()
	assert.NotEqual(t, uint64(0), len(str))
}

func TestBlock_RefreshHashCache(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.RefreshHashCache()
	assert.NotNil(t, result)
}

func TestBlock_Hash(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Hash()
	assert.NotNil(t, result)
}

func TestHeader_EncodeRlpToBytes(t *testing.T) {
	h := newTestHeader()
	_, err := h.EncodeRlpToBytes()
	assert.NoError(t, err)
}

func TestBlock_EncodeRlpToBytes(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	_, err := block.EncodeRlpToBytes()
	assert.NoError(t, err)
}

func TestBlock_Version(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Version()
	assert.Equal(t, block.header.Version, result)
}

func TestBlock_Number(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Number()
	assert.Equal(t, block.header.Number, result)
}

func TestBlock_PreHash(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.PreHash()
	assert.Equal(t, block.header.PreHash, result)
}

func TestBlock_Seed(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Seed()
	assert.Equal(t, block.header.Seed, result)
}

func TestBlock_Timestamp(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Timestamp()
	assert.NotNil(t, result)
}

func TestBlock_CoinBaseAddress(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.CoinBaseAddress()
	assert.Equal(t, block.header.CoinBase, result)
}

func TestBlock_Nonce(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Nonce()
	assert.Equal(t, block.header.Nonce, result)
}

func TestBlock_Difficulty(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Difficulty()
	assert.Equal(t, block.header.Diff, result)
}

func TestBlock_TxRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.TxRoot()
	assert.Equal(t, block.header.TransactionRoot, result)
}

func TestBlock_StateRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.StateRoot()
	assert.Equal(t, block.header.StateRoot, result)
}

func TestBlock_VerificationRoot(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.VerificationRoot()
	assert.Equal(t, block.header.VerificationRoot, result)
}

func TestBlock_GetBloom(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetBloomLog(model.BytesToBloom(common.Hex2Bytes("0x00000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000004000000000000000000000000000000000000000000000000")))
	result := block.GetBloom()

	assert.Equal(t, *block.header.Bloom, result)
	//blockByte, err := rlp.EncodeToBytes(block)
	//assert.NoError(t, err)
	//block.DecodeRLP(blockByte)
}

func TestBlock_GetBloomLog(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.GetBloomLog()
	assert.Equal(t, block.header.BloomLogs, result)
}

func TestBlock_CoinBase(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.CoinBase()
	assert.NotNil(t, result)
}

func TestBlock_TxCount(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.TxCount()
	assert.Equal(t, len(block.body.Txs), result)
}

func TestBlock_Header(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Header()
	assert.NotNil(t, result)

	block.header = nil
	result = block.Header()
	assert.Nil(t, result)
}

func TestBlock_Body(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	result := block.Body()
	assert.NotNil(t, result)
}

func TestBlock_Verifications(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	block.SetVerifications([]AbstractVerification{})
	result := block.Verifications()
	assert.NotNil(t, result)
}

func TestBlock_Transaction(t *testing.T) {
	block := CreateBlock(0, common.Hash{}, 1)
	txId := block.GetTransactions()[0].CalTxId()
	result := block.Transaction(common.HexToHash("123"))
	assert.Nil(t, result)

	result = block.Transaction(txId)
	assert.Equal(t, txId, result.CalTxId())
}

func TestBlockBy_Sort(t *testing.T) {
	block1 := CreateBlock(0, common.Hash{}, 2)
	block2 := CreateBlock(1, block1.Hash(), 2)
	bs := Blocks{block1, block2}

	BlockBy(func(b1, b2 *Block) bool { return true }).Sort(bs)
	assert.Equal(t, block1.Hash(), bs[1].Hash())
	assert.Equal(t, block2.Hash(), bs[0].Hash())
}

func Test_blockSorter_Len(t *testing.T) {
	block1 := CreateBlock(0, common.Hash{}, 2)
	block2 := CreateBlock(1, block1.Hash(), 2)
	bs := Blocks{block1, block2}

	sorter := blockSorter{bs, func(b1, b2 *Block) bool { return true }}
	length := sorter.Len()
	assert.Equal(t, 2, length)
}

func Test_blockSorter_Swap(t *testing.T) {
	block1 := CreateBlock(0, common.Hash{}, 2)
	block2 := CreateBlock(1, block1.Hash(), 2)
	bs := Blocks{block1, block2}

	sorter := blockSorter{bs, func(b1, b2 *Block) bool { return true }}
	sorter.Swap(0, 1)
	assert.Equal(t, block1.Hash(), sorter.blocks[1].Hash())
	assert.Equal(t, block2.Hash(), sorter.blocks[0].Hash())
}

func Test_blockSorter_Less(t *testing.T) {
	block1 := CreateBlock(0, common.Hash{}, 2)
	block2 := CreateBlock(1, block1.Hash(), 2)
	bs := Blocks{block1, block2}

	sorter := blockSorter{bs, func(b1, b2 *Block) bool { return true }}
	result := sorter.Less(0, 1)
	assert.Equal(t, true, result)
}

func creatBlockWithAllTx(n int, t *testing.T) *Block {
	header := NewHeader(1, 0, common.Hash{}, common.HexToHash("123456"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), aliceAddr, common.BlockNonce{})

	keyAlice, _ := CreateKey()
	ms := NewSigner(big.NewInt(1))
	tempTx := NewTransaction(uint64(0), bobAddr, big.NewInt(1000), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	tempTx.SignTx(keyAlice, ms)
	var res []*Transaction
	for i := 0; i < n; i++ {
		res = append(res, tempTx)
	}

	var voteList []AbstractVerification
	voteMsg := CreateSignedVote(0, 0, common.Hash{}, VoteMessage)
	for i := 0; i < (chain_config.GetChainConfig().VerifierNumber*2/3 + 1); i++ {
		voteList = append(voteList, voteMsg)
	}

	return NewBlock(header, res, voteList)
}

func Test_BlockTxNumber(t *testing.T) {
	//t.Skip()
	maxNormalTxNumber := chain_config.BlockGasLimit / model.TxGas
	assert.Equal(t, 160000, int(maxNormalTxNumber))

	tmpBlock := creatBlockWithAllTx(int(maxNormalTxNumber), t)
	blockByte, err := tmpBlock.EncodeRlpToBytes()
	log.Info("the block size is:", "size", len(blockByte))
	log.Info("the tx number is:", "txNumber", tmpBlock.Body().GetTxsSize())
	assert.NoError(t, err)
}
