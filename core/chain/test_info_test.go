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

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"reflect"
	"testing"
)

var (
	TrieError      = errors.New("trie error")
	EconomyErr     = errors.New("economy model error")
	foundationAddr common.Address

	testPriv1 = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	aliceAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	bobAddr   = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
)

func createTestStateDB(t *testing.T) (ethdb.Database, common.Hash) {
	db := ethdb.NewMemDatabase()
	tdb := state_processor.NewStateStorageWithCache(db)
	processor, err := NewBlockProcessor(nil, common.Hash{}, tdb)
	assert.NoError(t, err)

	// setup genesis block
	var earlyTokenContract contract.EarlyRewardContract
	err = util.ParseJson(contract.EarlyRewardContractStr, &earlyTokenContract)
	assert.NoError(t, err)

	// add foundationAddress account
	foundationAddr = earlyTokenContract.Owner
	err = processor.NewAccountState(foundationAddr)
	assert.NoError(t, err)
	err = processor.SetBalance(foundationAddr, big.NewInt(0).Mul(big.NewInt(999999), big.NewInt(consts.DIP)))
	assert.NoError(t, err)

	err = processor.PutContract(contract.EarlyContractAddress, reflect.ValueOf(&earlyTokenContract))
	assert.NoError(t, err)

	// add verifiers account
	for i := 0; i < len(VerifierAddress); i++ {
		err = processor.NewAccountState(VerifierAddress[i])
		assert.NoError(t, err)
	}

	// add alice account
	err = processor.NewAccountState(aliceAddr)
	assert.NoError(t, err)
	err = processor.SetBalance(aliceAddr, big.NewInt(0).Mul(big.NewInt(999999), big.NewInt(consts.DIP)))
	assert.NoError(t, err)

	root, err := processor.Commit()
	assert.NoError(t, err)

	tdb.TrieDB().Commit(root, false)
	return db, root
}

func createUnNormalTx() *model.Transaction {
	key1, _ := crypto.HexToECDSA(testPriv1)
	fs1 := model.NewMercurySigner(big.NewInt(1))
	tx := model.NewUnNormalTransaction(0, big.NewInt(1000), big.NewInt(1), model2.TxGas)
	signedTx, _ := tx.SignTx(key1, fs1)
	return signedTx
}

func createBlock(number uint64) *model.Block {
	block := model.CreateBlock(number, common.HexToHash("123456"), 10)
	votes := []model.AbstractVerification{
		createSignedVote2(number, block.Hash(), model.VoteMessage, testPriv1, aliceAddr),
	}
	block.SetVerifications(votes)
	return block
}

func createBlockWithoutCoinBase() *model.Block {
	header := model.NewHeader(1, 0, common.Hash{}, common.HexToHash("123456"), common.HexToDiff("1fffffff"), big.NewInt(0), common.Address{}, common.BlockNonce{})
	return model.NewBlock(header, nil, nil)
}

func createSignedVote2(num uint64, blockId common.Hash, voteType model.VoteMsgType, testPriv string, address common.Address) *model.VoteMsg {
	voteA := model.NewVoteMsg(num, num, blockId, voteType)
	hash := common.RlpHashKeccak256(voteA)
	key, _ := crypto.HexToECDSA(testPriv)
	sign, _ := crypto.Sign(hash.Bytes(), key)
	voteA.Witness.Address = address
	voteA.Witness.Sign = sign
	return voteA
}

func createGenesis() *Genesis {
	db := ethdb.NewMemDatabase()
	storage := state_processor.NewStateStorageWithCache(db)

	stateProcessor, _ := state_processor.MakeGenesisAccountStateProcessor(storage)
	registerProcessor, _ := registerdb.MakeGenesisRegisterProcessor(storage)

	// setup genesis block
	chainDb := chaindb.NewChainDB(db, model.MakeDefaultBlockDecoder())
	return DefaultGenesisBlock(chainDb, stateProcessor, registerProcessor,
		chain_config.GetChainConfig())
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		log.Info("stat temp dir error,maybe is not exist, maybe not")
		if os.IsNotExist(err) {
			log.Info("temp dir is not exist")
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Info(fmt.Sprintf("mkdir failed![%v]\n", err))
				return false
			}
			log.Info("make path sucessful")
			return true
		}
		log.Info("stat file error")
		return false
	}
	return true
}

type fakeAccountDBChain struct {
	state *state_processor.AccountStateDB
}

func (dbChain fakeAccountDBChain) CurrentBlock() model.AbstractBlock {
	panic("implement me")
}

func (dbChain fakeAccountDBChain) GetBlockByNumber(number uint64) model.AbstractBlock {
	if number == 0 || number == 1 {
		return nil
	}

	if number == 9 {
		block := factory.CreateSpecialBlock(9)
		//votes := []model.AbstractVerification{
		//	createSignedVote2(number, block.Hash(), model.VoteMessage, testPriv1, aliceAddr),
		//	createSignedVote2(number, block.Hash(), model.VoteMessage, testPriv1, aliceAddr),
		//}
		//block.SetVerifications(votes)
		return block
	}

	return createBlock(number)
}

func (dbChain fakeAccountDBChain) GetVerifiers(round uint64) []common.Address {
	return VerifierAddress
}

func (dbChain fakeAccountDBChain) StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error) {
	if num == 10 {
		db := ethdb.NewMemDatabase()
		storage := state_processor.NewStateStorageWithCache(db)
		return state_processor.NewAccountStateDB(common.Hash{}, storage)
	}
	return dbChain.state, nil
}

func (dbChain fakeAccountDBChain) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	if block.Number() == 9 {
		return false
	}
	return true
}

func (dbChain fakeAccountDBChain) GetLastChangePoint(block model.AbstractBlock) *uint64 {
	num := block.Number() - 1
	return &num
}

func (dbChain fakeAccountDBChain) GetSlot(block model.AbstractBlock) *uint64 {
	slot := uint64(5)
	return &slot
}

type fakeStateStorage struct {
	storageErr      error
	getErr          error
	setErr          error
	errKey          string
	contractBalance *big.Int
}

func (storage fakeStateStorage) OpenTrie(root common.Hash) (state_processor.StateTrie, error) {
	return fakeTrie{
			getErr:          storage.getErr,
			setErr:          storage.setErr,
			errKey:          storage.errKey,
			contractBalance: storage.contractBalance,
		},
		storage.storageErr
}

func (storage fakeStateStorage) OpenStorageTrie(addrHash, root common.Hash) (state_processor.StateTrie, error) {
	panic("implement me")
}

func (storage fakeStateStorage) CopyTrie(state_processor.StateTrie) state_processor.StateTrie {
	panic("implement me")
}

func (storage fakeStateStorage) TrieDB() *trie.Database {
	panic("implement me")
}

func (storage fakeStateStorage) DiskDB() ethdb.Database {
	return ethdb.NewMemDatabase()
}

type fakeTrie struct {
	getErr          error
	setErr          error
	errKey          string
	contractBalance *big.Int
}

func (t fakeTrie) TryGet(key []byte) ([]byte, error) {
	if t.errKey == string(key[22:]) {
		return nil, TrieError
	}
	if t.contractBalance != nil {
		result, _ := rlp.EncodeToBytes(t.contractBalance)
		return result, t.getErr
	}
	return []byte{128}, t.getErr
}

func (t fakeTrie) TryUpdate(key, value []byte) error {
	return t.setErr
}

func (t fakeTrie) TryDelete(key []byte) error {
	panic("implement me")
}

func (t fakeTrie) Commit(onleaf trie.LeafCallback) (common.Hash, error) {
	return common.Hash{}, TrieError
}

func (t fakeTrie) Hash() common.Hash {
	return common.Hash{}
}

func (t fakeTrie) NodeIterator(startKey []byte) trie.NodeIterator {
	panic("implement me")
}

func (t fakeTrie) GetKey([]byte) []byte {
	panic("implement me")
}

func (t fakeTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error {
	panic("implement me")
}

type fakeEconomyModel struct {
	DIPErr  error
	addrErr error
}

func (model fakeEconomyModel) GetMineMasterDIPReward(block model.AbstractBlock) (*big.Int, error) {
	if model.DIPErr != nil {
		return nil, model.DIPErr
	}
	return big.NewInt(10000), nil
}

func (model fakeEconomyModel) GetVerifierDIPReward(block model.AbstractBlock) (map[economy_model.VerifierType]*big.Int, error) {
	if model.DIPErr != nil {
		return nil, EconomyErr
	}
	rewardMap := make(map[economy_model.VerifierType]*big.Int, 3)
	rewardMap[economy_model.MasterVerifier] = big.NewInt(300)
	rewardMap[economy_model.CommitVerifier] = big.NewInt(100)
	rewardMap[economy_model.NotCommitVerifier] = big.NewInt(0)
	return rewardMap, nil
}

func (model fakeEconomyModel) GetMinimumTxFee() *big.Int {
	panic("implement me")
}

func (model fakeEconomyModel) GetInvestorInitBalance() map[common.Address]*big.Int {
	panic("implement me")
}

func (model fakeEconomyModel) GetDeveloperInitBalance() map[common.Address]*big.Int {
	panic("implement me")
}

func (model fakeEconomyModel) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (model fakeEconomyModel) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (model fakeEconomyModel) GetFoundation() economy_model.Foundation {
	return nil
}

func (model fakeEconomyModel) CheckAddressType(address common.Address) economy_model.EconomyModelAddress {
	panic("implement me")
}

func (model fakeEconomyModel) GetDiffVerifierAddress(preBlock, block model.AbstractBlock) (map[economy_model.VerifierType][]common.Address, error) {
	if model.addrErr != nil {
		return nil, model.addrErr
	}
	addressMap := make(map[economy_model.VerifierType][]common.Address, 1)
	addressMap[economy_model.MasterVerifier] = VerifierAddress[:2]
	return addressMap, nil
}

func (model fakeEconomyModel) GetAddressLockMoney(address common.Address, blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

func (model fakeEconomyModel) GetBlockYear(blockNumber uint64) (uint64, error) {
	panic("implement me")
}

func (model fakeEconomyModel) GetOneBlockTotalDIPReward(blockNumber uint64) (*big.Int, error) {
	panic("implement me")
}

type earlyContractFakeChainService struct{}

func (s *earlyContractFakeChainService) GetVerifiers(slotNum uint64) (addresses []common.Address) {
	return VerifierAddress
}

func (s *earlyContractFakeChainService) GetSlot(block model.AbstractBlock) *uint64 {
	slotSize := chain_config.GetChainConfig().SlotSize
	slot := block.Number() / slotSize
	return &slot
}

type fakeChain struct {
	block model.AbstractBlock
}

func (c fakeChain) GetBlockByNumber(number uint64) model.AbstractBlock {
	panic("implement me")
}

func (c fakeChain) CurrentBlock() model.AbstractBlock {
	return c.block
}

func (c fakeChain) Genesis() model.AbstractBlock {
	panic("implement me")
}

func (c fakeChain) GetCurrVerifiers() []common.Address {
	return []common.Address{aliceAddr}
}

func (c fakeChain) GetNextVerifiers() []common.Address {
	return []common.Address{bobAddr}
}

func (c fakeChain) GetHeaderByHash(hash common.Hash) model.AbstractHeader {
	panic("implement me")
}

func (c fakeChain) GetHeaderByNumber(number uint64) model.AbstractHeader {
	panic("implement me")
}

func (c fakeChain) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	if block == nil {
		return false
	}
	return true
}
