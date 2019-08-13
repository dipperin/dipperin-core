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
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"math/big"
)

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
	slotSize := chain_config.GetChainConfig().SlotSize
	slot := block.Number() / slotSize
	return &slot
}

func (chain *fakeChainReader) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	return true
}

func (chain *fakeChainReader) GetSlotByNum(num uint64) *uint64 {
	slotSize := chain_config.GetChainConfig().SlotSize
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
	slotSize := chain_config.GetChainConfig().SlotSize
	if number == uint64(4*slotSize+1) {
		return CreateSpecialBlock(number)
	}
	return CreateBlock(number)
}

func (chain *fakeChainReader) CurrentState() (*state_processor.AccountStateDB, error) {
	_, state, _ := CreateTestStateDBV()
	return state, nil
}

func (chain *fakeChainReader) StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error) {
	_, state, _ := CreateTestStateDBV()
	return state, nil
}

func newFakeEconomyModel() *fakeEconomyModel {
	return &fakeEconomyModel{}
}

type fakeEconomyModel struct{}

func (em fakeEconomyModel) GetBlockYear(blockNumber uint64) (uint64, error) {
	panic("implement me")
}

func (em fakeEconomyModel) GetOneBlockTotalDIPReward(blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (em fakeEconomyModel) GetMineMasterDIPReward(block model.AbstractBlock) (*big.Int, error) {
	panic("implement me")
}

func (em fakeEconomyModel) GetVerifierDIPReward(PreBlock model.AbstractBlock) (map[economy_model.VerifierType]*big.Int, error) {
	panic("implement me")
}

func (em fakeEconomyModel) GetInvestorInitBalance() map[common.Address]*big.Int {
	panic("implement me")
}

func (em fakeEconomyModel) GetDeveloperInitBalance() map[common.Address]*big.Int {
	panic("implement me")
}

func (em fakeEconomyModel) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (em fakeEconomyModel) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (em fakeEconomyModel) GetFoundation() economy_model.Foundation {
	panic("implement me")
}

func (em fakeEconomyModel) CheckAddressType(address common.Address) economy_model.EconomyModelAddress {
	panic("implement me")
}

func (em fakeEconomyModel) GetDiffVerifierAddress(preBlock, Block model.AbstractBlock) (map[economy_model.VerifierType][]common.Address, error) {
	panic("implement me")
}

func (em fakeEconomyModel) GetAddressLockMoney(address common.Address, blockNumber uint64) (*big.Int, error) {
	return big.NewInt(100), nil
}

var AlicePrivV = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
var BobPrivV = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
var AliceAddrV = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
var BobAddrV = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")

func InitTestV() (alice, bob *ecdsa.PrivateKey, fr *fakeChainReader) {
	alice, bob = CreatKeyV()
	fr = &fakeChainReader{currentBlock: nil}
	return

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
	testtx1 := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"), big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	testtx1.SignTx(key1, fs1)
	testtx2 := model.CreateRawLockTx(1, hashlock, big.NewInt(34564), big.NewInt(10000), g_testData.TestGasPrice, g_testData.TestGasLimit, alice, bob)
	testtx2.SignTx(key2, fs2)
	return testtx1, testtx2
}

func CreateBlockV() (*model.Block, *model.Block) {
	header1 := model.NewHeader(1, 1, common.HexToHash("1111"), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(324234), common.HexToAddress("032f14"), common.BlockNonceFromInt(432423))
	tx1, tx2 := CreatTestTxV()
	txs1 := []*model.Transaction{tx1}
	var msg1 []model.AbstractVerification
	block1 := model.NewBlock(header1, txs1, msg1)
	slotSize := chain_config.GetChainConfig().SlotSize
	header2 := model.NewHeader(1, 10*slotSize, block1.Hash(), common.HexToHash("1111"), common.HexToDiff("1fffffff"), big.NewInt(324254), common.HexToAddress("032f14fd2"), common.BlockNonceFromInt(432424))
	txs2 := []*model.Transaction{tx2}
	var msg2 []model.AbstractVerification
	block2 := model.NewBlock(header2, txs2, msg2)
	return block1, block2
}

func CreateTestStateDBV() (state_processor.StateStorage, *state_processor.AccountStateDB, common.Hash) {
	db := ethdb.NewMemDatabase()
	//todo New method does not take MPT tree from the underlying database
	tdb := state_processor.NewStateStorageWithCache(db)
	processor, err := state_processor.NewAccountStateDB(common.Hash{}, tdb)
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

func CreateEmptyBlockByPH(number uint64, preHash common.Hash) *model.Block {
	coinbase := common.HexToAddress("0x223432423234")
	header := &model.Header{Number: number, PreHash: preHash, CoinBase: coinbase}
	txs := []*model.Transaction{}
	vfy := []model.AbstractVerification{}
	b := model.NewBlock(header, txs, vfy)
	return b
}

func CreateOneTxBlockByPH(number uint64, preHash common.Hash, transaction *model.Transaction) *model.Block {
	coinbase := common.HexToAddress("0x223432423234")
	header := &model.Header{Number: number, PreHash: preHash, CoinBase: coinbase, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}
	txs := []*model.Transaction{transaction}
	vfy := []model.AbstractVerification{}
	b := model.NewBlock(header, txs, vfy)
	return b
}
func CreateTxsBlockByPH(number uint64, preHash common.Hash, transactions []*model.Transaction) *model.Block {

	header := &model.Header{Number: number, PreHash: preHash}
	txs := transactions
	vfy := []model.AbstractVerification{}
	b := model.NewBlock(header, txs, vfy)
	return b
}
