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
package registerdb

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third_party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"math/big"
	"time"
)

var(
	testPriv1 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	testPriv2 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	AlicePrivV = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	BobPrivV = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	AliceAddrV = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	BobAddrV = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
	alicePriv  = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	SlotError  = errors.New("slot test error")
	PointError = errors.New("point test error")
	TestError  = errors.New("test error")
)

func CreateBlock(number uint64) *model.Block {

	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))

	tx1, tx2 := CreateTestTx()
	txs := []*model.Transaction{tx1, tx2}
	vfy := []model.AbstractVerification{}

	b := model.NewBlock(header, txs, vfy)

	return b
}

func CreateKey() (*ecdsa.PrivateKey, *ecdsa.PrivateKey) {
	key1, err1 := crypto.HexToECDSA(testPriv1)
	key2, err2 := crypto.HexToECDSA(testPriv2)
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	return key1, key2
}

func CreateTestTx() (*model.Transaction, *model.Transaction) {
	key1, key2 := CreateKey()
	fs1 := model.NewSigner(big.NewInt(1))
	fs2 := model.NewSigner(big.NewInt(3))
	testTx1 := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000), model.TestGasPrice, model.TestGasLimit, []byte{})
	gasUsed, _ := model.IntrinsicGas(testTx1.ExtraData(), false, false)
	testTx1.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), testTx1.GetGasPrice()))
	testTx1.SignTx(key1, fs1)

	testTx2 := model.NewTransaction(10, common.HexToAddress("0121321432423534534535"), big.NewInt(20000), model.TestGasPrice, model.TestGasLimit, []byte{})
	gasUsed, _ = model.IntrinsicGas(testTx2.ExtraData(), false, false)
	testTx2.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), testTx2.GetGasPrice()))
	testTx2.SignTx(key2, fs2)
	return testTx1, testTx2
}

func CreateSpecialBlock(number uint64) *model.Block {

	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.Difficulty{0}, big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonce{0})
	b := model.NewBlock(header, []*model.Transaction{}, []model.AbstractVerification{})

	return b
}

func CreatKeyV() (*ecdsa.PrivateKey, *ecdsa.PrivateKey) {
	key1, err1 := crypto.HexToECDSA(AlicePrivV)
	key2, err2 := crypto.HexToECDSA(BobPrivV)
	if err1 != nil || err2 != nil {
		return nil, nil
	}
	return key1, key2
}

func CreatTestTxV() (*model.Transaction, *model.Transaction) {
	key1, key2 := CreatKeyV()
	fs1 := model.NewSigner(big.NewInt(1))
	fs2 := model.NewSigner(big.NewInt(3))
	alice := cs_crypto.GetNormalAddress(key1.PublicKey)
	bob := cs_crypto.GetNormalAddress(key2.PublicKey)
	hashkey := []byte("123")
	hashlock := cs_crypto.Keccak256Hash(hashkey)
	testtx1 := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000), model.TestGasPrice, model.TestGasLimit, []byte{})
	testtx1.SignTx(key1, fs1)
	testtx2 := model.CreateRawLockTx(1, hashlock, big.NewInt(34564), big.NewInt(10000), model.TestGasPrice, model.TestGasLimit, alice, bob)
	testtx2.SignTx(key2, fs2)
	return testtx1, testtx2
}

func CreateBlockV() (*model.Block, *model.Block) {
	header1 := model.NewHeader(1, 1, common.HexToHash("1111"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(324234), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
	tx1, tx2 := CreatTestTxV()
	txs1 := []*model.Transaction{tx1}
	var msg1 []model.AbstractVerification
	block1 := model.NewBlock(header1, txs1, msg1)
	slotSize := chainconfig.GetChainConfig().SlotSize
	header2 := model.NewHeader(1, 10*slotSize, block1.Hash(), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(324254), common.HexToAddress("032f14fd2"), common.BlockNonceFromInt(432424))
	txs2 := []*model.Transaction{tx2}
	var msg2 []model.AbstractVerification
	block2 := model.NewBlock(header2, txs2, msg2)
	return block1, block2
}

func CreateTestStateDBV() (stateprocessor.StateStorage, *stateprocessor.AccountStateDB, common.Hash) {
	db := ethdb.NewMemDatabase()
	//todo New method does not take MPT tree from the underlying database
	tdb := stateprocessor.NewStateStorageWithCache(db)
	processor, err := stateprocessor.NewAccountStateDB(common.Hash{}, tdb)
	if err != nil {
		panic(err)
	}

	processor.NewAccountState(AliceAddrV)
	processor.NewAccountState(BobAddrV)
	processor.AddBalance(AliceAddrV, big.NewInt(5000))
	processor.AddBalance(BobAddrV, big.NewInt(5000))
	processor.Stake(AliceAddrV, big.NewInt(200))
	processor.Stake(BobAddrV, big.NewInt(200))
	processor.SetNonce(AliceAddrV, uint64(2))
	processor.SetLastElect(AliceAddrV, uint64(0))

	processor.SetLastElect(BobAddrV, uint64(9))

	root, err := processor.Commit()
	if err != nil {
		panic(err)
	}
	return tdb, processor, root
}
/*
ChainReader test set
*/
func NewFakeReader(block *model.Block) *fakeChainReader {
	return &fakeChainReader{
		currentBlock: block,
	}
}

type fakeChainReader struct {
	currentBlock model.AbstractBlock
}

func (chain *fakeChainReader) GetVerifiers(round uint64) []common.Address {
	panic("implement me")
}

func (chain *fakeChainReader) GetLastChangePoint(block model.AbstractBlock) *uint64 {
	panic("implement me")
}

func (chain *fakeChainReader) GetSlot(block model.AbstractBlock) *uint64 {
	slotSize := chainconfig.GetChainConfig().SlotSize
	slot := block.Number() / slotSize
	return &slot
}

func (chain *fakeChainReader) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	return true
}

func (chain *fakeChainReader) GetSlotByNum(num uint64) *uint64 {
	slotSize := chainconfig.GetChainConfig().SlotSize
	slot := num / slotSize
	return &slot
}

func (chain *fakeChainReader) CurrentBlock() model.AbstractBlock {
	if chain.currentBlock.(*model.Block) != nil {
		return chain.currentBlock
	}
	_, block2 := CreateBlockV()
	return block2
}

func (chain *fakeChainReader) GetLatestNormalBlock() model.AbstractBlock {
	return nil
}

func (chain *fakeChainReader) GetBlockByNumber(number uint64) model.AbstractBlock {
	slotSize := chainconfig.GetChainConfig().SlotSize
	if number == uint64(4*slotSize+1) {
		return CreateSpecialBlock(number)
	}
	return CreateBlock(number)
}

func (chain *fakeChainReader) CurrentState() (*stateprocessor.AccountStateDB, error) {
	_, state, _ := CreateTestStateDBV()
	return state, nil
}

func (chain *fakeChainReader) StateAtByBlockNumber(num uint64) (*stateprocessor.AccountStateDB, error) {
	_, state, _ := CreateTestStateDBV()
	return state, nil
}


func createRegisterDb(blockNum uint64) *RegisterDB {
	db := ethdb.NewMemDatabase()
	storage := stateprocessor.NewStateStorageWithCache(db)
	block := CreateBlock(blockNum)
	reader := NewFakeReader(block)
	register, _ := NewRegisterDB(common.Hash{}, storage, reader)
	return register
}

func createBlock(number uint64, txs []*model.Transaction) *model.Block {
	header := model.NewHeader(1, number, common.HexToHash("0x12312fa0929348"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
	return model.NewBlock(header, txs, []model.AbstractVerification{})
}

func createRegisterTX(nonce uint64, amount *big.Int) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewRegisterTransaction(nonce, amount, model.TestGasPrice, model.TestGasLimit)
	key, _ := crypto.HexToECDSA(alicePriv)
	signedTx, _ := tx.SignTx(key, fs1)
	return signedTx
}

func createCannelTX(nonce uint64) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewCancelTransaction(nonce, model.TestGasPrice, model.TestGasLimit)
	key, _ := crypto.HexToECDSA(alicePriv)
	signedTx, _ := tx.SignTx(key, fs1)
	return signedTx
}

type fakeTrie struct {
	slotErr  error
	pointErr error
}

func (f fakeTrie) TryGet(key []byte) ([]byte, error) {
	if bytes.Equal(key, []byte(slotKey)) {
		return nil, f.slotErr
	}

	if bytes.Equal(key, []byte(lastChangePointKey)) {
		return nil, f.pointErr
	}
	return nil, nil
}

func (f fakeTrie) TryUpdate(key, value []byte) error {
	if bytes.Equal(key, []byte(slotKey)) {
		return f.slotErr
	}

	if bytes.Equal(key, []byte(lastChangePointKey)) {
		return f.pointErr
	}
	return TestError
}

func (f fakeTrie) TryDelete(key []byte) error {
	return TestError
}

func (f fakeTrie) Commit(onleaf trie.LeafCallback) (common.Hash, error) {
	return common.Hash{}, SlotError
}

func (f fakeTrie) Hash() common.Hash {
	panic("implement me")
}

func (f fakeTrie) NodeIterator(startKey []byte) trie.NodeIterator {
	panic("implement me")
}

func (f fakeTrie) GetKey([]byte) []byte {
	panic("implement me")
}

func (f fakeTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error {
	panic("implement me")
}