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

package ver_halt_check

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/verifiers-halt-check"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
	"time"
)

// Test special block and normal block replacement logic
// 1. If you receive a normal block but there is a special block in same height, it failed to be saved.
// 2. If you receive a special block but there is a normal block in same height, it will replace the normal block.
// 3. If you receive a special block but there is a normal block in height+1, it failed to be saved.

var testVerBootAccounts []tests.Account
var testFee = economy_model.GetMinimumTxFee(50)

func init() {
	var err error
	testVerBootAccounts, err = tests.ChangeVerBootNodeAddress()
	if err != nil {
		panic("change verifier boot node address error for test")
	}
}

type fakeCacheDB struct{
	commits map[uint64]model.Verifications
}

func (c *fakeCacheDB) GetSeenCommits(blockHeight uint64, blockHash common.Hash) (result []model.AbstractVerification, err error) {
	if _,ok := c.commits[blockHeight];ok{
		return c.commits[blockHeight],nil
	}
	return nil, errors.New("not commits")
}

func (c *fakeCacheDB) SaveSeenCommits(blockHeight uint64, blockHash common.Hash, commits []model.AbstractVerification) error {
	c.commits[blockHeight] = commits
	return nil
}

type fakeTxPool struct{}

func (t *fakeTxPool) Reset(oldHead, newHead *model.Header) {
	return
}

type fakeWalletSigner struct{}

func (fakeWalletSigner) GetAddress() common.Address {
	return testVerBootAccounts[0].Address()
}

func (fakeWalletSigner) SignHash(hash []byte) ([]byte, error) {
	return testVerBootAccounts[0].SignHash(hash)
}

func (fakeWalletSigner) PublicKey() *ecdsa.PublicKey {
	return &testVerBootAccounts[0].Pk.PublicKey
}

func (fakeWalletSigner) ValidSign(hash []byte, pubKey []byte, sign []byte) error {
	if len(sign) == 0 {
		return accounts.ErrEmptySign
	}
	if crypto.VerifySignature(pubKey,hash,sign[:len(sign)-1]) == true{
		return nil
	}else{
		return accounts.ErrSignatureInvalid
	}
}

func (fakeWalletSigner) Evaluate(account accounts.Account, seed []byte) (index [32]byte, proof []byte, err error) {
	index,proof = crypto.Evaluate(testVerBootAccounts[0].Pk, seed)
	return index,proof,nil
}


type testNeedConfig struct {
	*cs_chain.CsChainService
	*tests.GenesisEnv
	*tests.BlockBuilder
	*verifiers_halt_check.StateHandler
}

func getTestChainEnv(db cs_chain.CacheDB, pool cs_chain.TxPool) (*cs_chain.CsChainService, *tests.GenesisEnv, *tests.TxBuilder, *tests.BlockBuilder) {
	f := chain_writer.NewChainWriterFactory()
	cConf := chain_config.GetChainConfig()
	cConf.VerifierNumber = 4

	chainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir: "",
		WriterFactory: f,
		ChainConfig: cConf,
	})

	// 主要是初始化默认验证者在其中，外边测试调用它的方法对block投票
	attackEnv := tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), nil)
	txB := &tests.TxBuilder{
		Nonce:  1,
		To:     common.HexToAddress(fmt.Sprintf("0x123a%v", 1)),
		Amount: big.NewInt(1),
		Fee:    testFee,
		Pk:     attackEnv.DefaultVerifiers()[0].Pk,
	}
	bb := &tests.BlockBuilder{
		ChainState: chainState,
		PreBlock:     chainState.CurrentBlock(),
		MinerPk:      attackEnv.Miner().Pk,
	}

	ccs := cs_chain.NewCsChainService(&cs_chain.CsChainServiceConfig{
		CacheDB: db,
		TxPool: pool,
	}, chainState)
	f.SetChain(ccs.CacheChainState)
	return ccs, attackEnv, txB, bb
}

func getTestNeedConfig() testNeedConfig{
	cMock := &fakeCacheDB{
		commits:make(map[uint64]model.Verifications,0),
	}
	pMock := &fakeTxPool{}
	ccs, gEnv, _, bB := getTestChainEnv(cMock, pMock)
	testEconomy := economy_model.MakeDipperinEconomyModel(ccs,economy_model.DIPProportion)
	stateHandler :=verifiers_halt_check.MakeHaltCheckStateHandler(ccs,&fakeWalletSigner{},testEconomy)

	return testNeedConfig{
		ccs,
		gEnv,
		bB,
		stateHandler,
	}
}

func generateTestConfigAndBlocks() (config testNeedConfig,normalBlock model.AbstractBlock,proposal verifiers_halt_check.ProposalMsg, err error){
	testConf := getTestNeedConfig()

	block := testConf.Build()
	err = testConf.SaveBlock(block, testConf.VoteBlock(4, 1, block))
	if err !=nil {
		return testNeedConfig{},nil,verifiers_halt_check.ProposalMsg{},err
	}

	//generate empty block
	proposalConfig,err :=testConf.GenProposalConfig(model.VerBootNodeVoteMessage)
	if err !=nil {
		return testNeedConfig{},nil,verifiers_halt_check.ProposalMsg{},err
	}

	haltHandler := verifiers_halt_check.NewHaltHandler(proposalConfig)
	proposal,err = haltHandler.ProposeEmptyBlock()
	if err !=nil {
		return testNeedConfig{},nil,verifiers_halt_check.ProposalMsg{},err
	}

	testConf.SetPreBlock(testConf.CurrentBlock())
	testConf.SetVerifivations(testConf.GetSeenCommit(testConf.CurrentBlock().Number()))
	sameHeightBlock := testConf.Build()

	return testConf,sameHeightBlock,proposal,nil
}

func TestSaveEmptyBlock(t *testing.T) {
	testConf,_,proposal,err :=generateTestConfigAndBlocks()
	assert.NoError(t,err)
	//save empty block
	err = testConf.SaveBlock(&proposal.EmptyBlock,model.Verifications{&proposal.VoteMsg})
	assert.NoError(t,err)
}

func TestSaveSameHeightNormalBlock(t *testing.T) {
	testConf,sameHeightBlock,proposal,err :=generateTestConfigAndBlocks()
	assert.NoError(t,err)

	//save empty block
	err = testConf.SaveBlock(&proposal.EmptyBlock,model.Verifications{&proposal.VoteMsg})
	assert.NoError(t,err)

	//save normal block
	err = testConf.SaveBlock(sameHeightBlock,testConf.VoteBlock(4,1,sameHeightBlock))
	assert.Error(t,err)
}

func TestSaveSameHeightEmptyBlock(t *testing.T){
	testConf,sameHeightBlock,proposal,err :=generateTestConfigAndBlocks()
	assert.NoError(t,err)

	//save normal block
	err = testConf.SaveBlock(sameHeightBlock,testConf.VoteBlock(4,1,sameHeightBlock))
	assert.NoError(t,err)
	assert.Equal(t,sameHeightBlock,testConf.CurrentBlock())

	//save empty block
	err = testConf.SaveBlock(&proposal.EmptyBlock,model.Verifications{&proposal.VoteMsg})
	assert.NoError(t,err)
	assert.Equal(t,&proposal.EmptyBlock,testConf.CurrentBlock())
}

func TestSaveLowHeightEmptyBlock(t *testing.T){
	testConf,sameHeightBlock,proposal,err :=generateTestConfigAndBlocks()
	assert.NoError(t,err)

	//save normal block
	err = testConf.SaveBlock(sameHeightBlock,testConf.VoteBlock(4,1,sameHeightBlock))
	assert.NoError(t,err)

	//generate and save another block
	testConf.SetPreBlock(testConf.CurrentBlock())
	testConf.SetVerifivations(testConf.GetSeenCommit(testConf.CurrentBlock().Number()))
	newHeightBlock := testConf.Build()
	err = testConf.SaveBlock(newHeightBlock,testConf.VoteBlock(4,1,newHeightBlock))
	assert.NoError(t,err)
	assert.Equal(t,proposal.EmptyBlock.Number(),testConf.CurrentBlock().Number()-1)

	//save empty block with low height
	err = testConf.SaveBlock(&proposal.EmptyBlock,model.Verifications{&proposal.VoteMsg})
	assert.Error(t,err)
}

func TestSaveInvalidDifficultyBlock(t *testing.T){
	testConf,sameHeightBlock,_,err :=generateTestConfigAndBlocks()
	assert.NoError(t,err)

	block := reflect.ValueOf(sameHeightBlock)
	param := []reflect.Value{
		reflect.ValueOf(common.Difficulty{}),
	}

	log.Info("the block method number is:","number",block.NumMethod())
	for i:=0;i<block.NumMethod();i++{
		log.Info("the method is:","method",block.Method(i))
	}
	block.MethodByName("SetDifficulty").Call(param)
	sameHeightBlock.RefreshHashCache()

	log.Info("the mockBlock diff is:","diff",sameHeightBlock.Difficulty().Hex())
	err = testConf.SaveBlock(sameHeightBlock,sameHeightBlock.GetVerifications())
	log.Info("the err is:","err",err)
	assert.Equal(t,g_error.ErrInvalidDiff,err)
}

func TestSaveInvalidTimeStampBlock(t *testing.T) {
	testConf, sameHeightBlock, _, err := generateTestConfigAndBlocks()
	assert.NoError(t, err)

	//set error timeStamp
	errTimeStamp := time.Now().Add(testConf.GetChainConfig().BlockTimeRestriction + time.Second).UnixNano()
	block := reflect.ValueOf(sameHeightBlock)
	param := []reflect.Value{
		reflect.ValueOf(big.NewInt(errTimeStamp)),
	}

	block.MethodByName("SetTimeStamp").Call(param)
	model.CalNonce(sameHeightBlock.(*model.Block))
	sameHeightBlock.RefreshHashCache()

	err = testConf.SaveBlock(sameHeightBlock, sameHeightBlock.GetVerifications())
	assert.Equal(t,g_error.ErrBlockTimeStamp,err)
}


//debug dipperIn-core
func TestSystemBug(t *testing.T){
	t.Skip()
	vb2 := "/home/qydev/yc/debug/err-log-20190329/b2Data"
	vb2ChainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       vb2,
		WriterFactory: chain_writer.NewChainWriterFactory(),
		ChainConfig:   chain_config.GetChainConfig(),
	})

	currentBlock := vb2ChainState.CurrentBlock()
	log.Info("the vb2 currentBlock is:","number",currentBlock.Number())
	testBlock:=vb2ChainState.GetBlockByNumber(13867)
	log.Info("the vb2 13867 block info is:","id",testBlock.Hash().Hex(),"preHash",testBlock.PreHash().Hex())
	log.Info("the vb2 13866 block info is:","id",vb2ChainState.GetBlockByNumber(13866).Hash().Hex())

	blockHash ,err:= hexutil.Decode("0x000000e9cd7031841cdbef5465b80a04b18f00e49278e56588a4a541092046d8")
	assert.NoError(t,err)
	errBlock := vb2ChainState.GetBlockByHash(common.BytesToHash(blockHash))
	assert.NotEqual(t,nil,errBlock)
	log.Info("the errBlock number is:","errBlockNumber",errBlock.Number())

	log.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	v1 := "/home/qydev/yc/debug/err-log-20190329/v1Data"
	v1ChainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       v1,
		WriterFactory: chain_writer.NewChainWriterFactory(),
		ChainConfig:   chain_config.GetChainConfig(),
	})

	currentBlock2 := v1ChainState.CurrentBlock()
	log.Info("the v1 currentBlock is:","number",currentBlock2.Number())
	testBlock2:=vb2ChainState.GetBlockByNumber(13866)
	log.Info("the v1 block 13866 blockHash is:","blockHash",testBlock2.Hash().Hex())
	log.Info("the correct block is:","correctBlock",testBlock2)

	assert.NoError(t,err)
	errBlock = v1ChainState.GetBlockByHash(common.BytesToHash(blockHash))
	assert.NotEqual(t,nil,errBlock)
	log.Info("the errBlock is:","errBlock",errBlock)
	log.Info("the errBlock number is:","errBlockNumber",errBlock.Number())


}










