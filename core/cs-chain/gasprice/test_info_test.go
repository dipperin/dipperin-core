package gasprice

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	alicePriv = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	aliceAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
)

func createCsChain(accounts []tests.Account) *chain_state.ChainState {
	f := chain_writer.NewChainWriterFactory()
	chainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   chain_config.GetChainConfig(),
	})
	f.SetChain(chainState)

	// Mainly the default initial verifier is in it, the outer test need to call it to vote on the block
	tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), accounts)
	return chainState
}

func createSignedTx3(nonce uint64, amount, gasPrice *big.Int) *model.Transaction {
	verifiers, _ := tests.ChangeVerifierAddress(nil)
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewTransaction(nonce, aliceAddr, amount, gasPrice, g_testData.TestGasLimit, []byte{})
	signedTx, _ := tx.SignTx(verifiers[0].Pk, fs1)
	return signedTx
}

func createBlock(chain *chain_state.ChainState, txs []*model.Transaction, votes []model.AbstractVerification) model.AbstractBlock {
	key1, _ := crypto.HexToECDSA(alicePriv)
	bb := &tests.BlockBuilder{
		ChainState: chain,
		PreBlock:   chain.CurrentBlock(),
		Txs:        txs,
		Vers:       votes,
		MinerPk:    key1,
	}
	return bb.Build()
}

func createVerifiersVotes(block model.AbstractBlock, votesNum int, testAccounts []tests.Account) (votes []model.AbstractVerification) {
	testVerifierAccounts, _ := tests.ChangeVerifierAddress(testAccounts)
	for i := 0; i < votesNum; i++ {
		voteA := model.NewVoteMsg(block.Number(), uint64(0), block.Hash(), model.VoteMessage)
		sign, _ := crypto.Sign(voteA.Hash().Bytes(), testVerifierAccounts[i].Pk)
		voteA.Witness.Address = testVerifierAccounts[i].Address()
		voteA.Witness.Sign = sign
		votes = append(votes, voteA)
	}
	return
}

func insertBlockToChain(t *testing.T, chain *chain_state.ChainState, num int, txs []*model.Transaction) {
	curNum := int(chain.CurrentBlock().Number())
	config := chain_config.GetChainConfig()
	for i := curNum; i < curNum+num; i++ {
		curBlock := chain.CurrentBlock()
		var block model.AbstractBlock
		if curBlock.Number() == 0 {
			block = createBlock(chain, txs, nil)
		} else {

			// votes for curBlock on chain
			curBlockVotes := createVerifiersVotes(curBlock, config.VerifierNumber*2/3+1, nil)
			block = createBlock(chain, txs, curBlockVotes)
		}

		// votes for build block
		votes := createVerifiersVotes(block, config.VerifierNumber*2/3+1, nil)
		err := chain.SaveBftBlock(block, votes)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+1), chain.CurrentBlock().Number())
		assert.Equal(t, false, chain.CurrentBlock().IsSpecial())
	}
}
