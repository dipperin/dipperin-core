package gasprice

import (
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewOracle(t *testing.T) {
	config := GasPriceConfig{
		Blocks:     -20,
		Percentile: -60,
		Default:    big.NewInt(chainconfig.DefaultGasPrice),
	}
	oracle := NewOracle(nil, config)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 0, oracle.percentile)
	
	config.Percentile = 200
	oracle = NewOracle(nil, config)
	assert.Equal(t, 100, oracle.percentile)
}

func TestOracle_SuggestPrice(t *testing.T) {
	// TODO: need env of tests
	// csChain := createCsChain(nil)
	// config := GasPriceConfig{
	// 	Blocks:     20,
	// 	Percentile: 60,
	// 	Default:    big.NewInt(chainconfig.DefaultGasPrice),
	// }
	//
	// insertBlockToChain(t, csChain, 5, nil)
	// oracle := NewOracle(csChain, config)
	// gasPrice, err := oracle.SuggestPrice()
	// assert.NoError(t, err)
	// assert.Equal(t, big.NewInt(1), gasPrice)
	//
	// // insert 10 blocks with txs
	// var txs1 []*model.Transaction
	// for i := 0; i < 50; i++ {
	// 	tx := createSignedTx3(uint64(i), big.NewInt(0), big.NewInt(int64(i)))
	// 	txs1 = append(txs1, tx)
	// }
	// for i := 0; i < 10; i++ {
	// 	insertBlockToChain(t, csChain, 1, txs1[i*5:(i+1)*5])
	// }
	// gasPrice, err = oracle.SuggestPrice()
	// assert.NoError(t, err)
	// assert.Equal(t, big.NewInt(25), gasPrice)
	//
	// // insert 15 blocks without txs
	// insertBlockToChain(t, csChain, 15, nil)
	// gasPrice, err = oracle.SuggestPrice()
	// assert.NoError(t, err)
	// assert.Equal(t, big.NewInt(25), gasPrice)
	//
	// // insert 10 blocks with txs
	// var txs2 []*model.Transaction
	// for i := 0; i < 50; i++ {
	// 	tx := createSignedTx3(uint64(i+50), big.NewInt(0), big.NewInt(int64(100-i)))
	// 	txs2 = append(txs2, tx)
	// }
	// for i := 0; i < 10; i++ {
	// 	insertBlockToChain(t, csChain, 1, txs2[i*5:(i+1)*5])
	// }
	// gasPrice, err = oracle.SuggestPrice()
	// assert.NoError(t, err)
	// assert.Equal(t, big.NewInt(76), gasPrice)
	//
	// // get price again
	// gasPrice, err = oracle.SuggestPrice()
	// assert.NoError(t, err)
	// assert.Equal(t, big.NewInt(76), gasPrice)
}

// TODO: need env of tests
// MARK: info for test
// var (
// 	alicePriv = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
// 	aliceAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
// )
//
// func createCsChain(accounts []tests.Account) *chainstate.ChainState {
// 	f := chainwriter.NewChainWriterFactory()
// 	chainState := chainstate.NewChainState(&chainstate.ChainStateConfig{
// 		DataDir:       "",
// 		WriterFactory: f,
// 		ChainConfig:   chainconfig.GetChainConfig(),
// 	})
// 	f.SetChain(chainState)
//
// 	// Mainly the default initial verifier is in it, the outer test need to call it to vote on the block
// 	tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), accounts)
// 	return chainState
// }
//
// func createSignedTx3(nonce uint64, amount, gasPrice *big.Int) *model.Transaction {
// 	verifiers, _ := tests.ChangeVerifierAddress(nil)
// 	fs1 := model.NewSigner(big.NewInt(1))
// 	tx := model.NewTransaction(nonce, aliceAddr, amount, gasPrice, g_testData.TestGasLimit, []byte{})
// 	signedTx, _ := tx.SignTx(verifiers[0].Pk, fs1)
// 	return signedTx
// }
//
// func createBlock(chain *chainstate.ChainState, txs []*model.Transaction, votes []model.AbstractVerification) model.AbstractBlock {
// 	key1, _ := crypto.HexToECDSA(alicePriv)
// 	bb := &tests.BlockBuilder{
// 		ChainState: chain,
// 		PreBlock:   chain.CurrentBlock(),
// 		Txs:        txs,
// 		Vers:       votes,
// 		MinerPk:    key1,
// 	}
// 	return bb.Build()
// }
//
// func createVerifiersVotes(block model.AbstractBlock, votesNum int, testAccounts []tests.Account) (votes []model.AbstractVerification) {
// 	testVerifierAccounts, _ := tests.ChangeVerifierAddress(testAccounts)
// 	for i := 0; i < votesNum; i++ {
// 		voteA := model.NewVoteMsg(block.Number(), uint64(0), block.Hash(), model.VoteMessage)
// 		sign, _ := crypto.Sign(voteA.Hash().Bytes(), testVerifierAccounts[i].Pk)
// 		voteA.Witness.Address = testVerifierAccounts[i].Address()
// 		voteA.Witness.Sign = sign
// 		votes = append(votes, voteA)
// 	}
// 	return
// }
//
// func insertBlockToChain(t *testing.T, chain *chainstate.ChainState, num int, txs []*model.Transaction) {
// 	curNum := int(chain.CurrentBlock().Number())
// 	config := chainconfig.GetChainConfig()
// 	for i := curNum; i < curNum+num; i++ {
// 		curBlock := chain.CurrentBlock()
// 		var block model.AbstractBlock
// 		if curBlock.Number() == 0 {
// 			block = createBlock(chain, txs, nil)
// 		} else {
//
// 			// votes for curBlock on chain
// 			curBlockVotes := createVerifiersVotes(curBlock, config.VerifierNumber*2/3+1, nil)
// 			block = createBlock(chain, txs, curBlockVotes)
// 		}
//
// 		// votes for build block
// 		votes := createVerifiersVotes(block, config.VerifierNumber*2/3+1, nil)
// 		err := chain.SaveBftBlock(block, votes)
// 		assert.NoError(t, err)
// 		assert.Equal(t, uint64(i+1), chain.CurrentBlock().Number())
// 		assert.Equal(t, false, chain.CurrentBlock().IsSpecial())
// 	}
// }
