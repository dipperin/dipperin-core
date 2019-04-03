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

package verifiers_halt_check

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var bootNodeIndex int

func init() {
	log.Info("change ver boot node address for test")
	var err error
	testVerBootAccounts, err = tests.ChangeVerBootNodeAddress()
	if err != nil {
		panic("change verifier boot node address error for test")
	}
}
func verBootSignForTest(hash []byte) ([]byte, error) {

	sk := testVerBootAccounts[bootNodeIndex].Pk
	return crypto.Sign(hash, sk)
}

func stateProcess(block model.AbstractBlock, preStateRoot, preRegisterRoot common.Hash) (stateRoot, registerRoot common.Hash, err error) {
	return common.Hash{}, common.Hash{}, nil
}

func MakeTestProposalConfig(voteType model.VoteMsgType, verBootIndex int) ProposalGeneratorConfig {
	bootNodeIndex = verBootIndex
	return ProposalGeneratorConfig{
		CurBlock:         factory.CreateBlock(2),
		SignHashFunc:     verBootSignForTest,
		ProcessStateFunc: stateProcess,
		VoteType:         voteType,
		PubKey:           crypto.FromECDSAPub(&testVerBootAccounts[bootNodeIndex].Pk.PublicKey),
	}
}

func MakeTestProposalMsg(verBootIndex int) (*ProposalMsg, error) {
	config := MakeTestProposalConfig(model.VerBootNodeVoteMessage, verBootIndex)
	return GenProposalMsg(config)
}

func TestStateHandler_SaveFinalEmptyBlock(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mockChainReader := NewMockNeedChainReaderFunction(c)
	mockWalletSigner := NewMockNeedWalletSigner(c)
	mockEconomyModel := NewMockEconomyModel(c)

	stateHandler := MakeHaltCheckStateHandler(mockChainReader, mockWalletSigner, mockEconomyModel)

	proposal ,err:= MakeTestProposalMsg(0)
	assert.NoError(t,err)

	vote1 := model.NewVoteMsg(proposal.VoteMsg.Height,0,proposal.VoteMsg.BlockID,model.AliveVerifierVoteMessage)
	vote1.Witness.Address = chain_config.MercuryVerifierAddress[0]
	vote1.Timestamp = proposal.VoteMsg.Timestamp

	//map random so just test one vote
/*	vote2 := model.NewVoteMsg(proposal.VoteMsg.Height,0,proposal.VoteMsg.BlockID,model.AliveVerifierVoteMessage)
	vote2.Witness.Address = chain_config.MercuryVerifierAddress[1]
	vote2.Timestamp = proposal.VoteMsg.Timestamp*/

	votes := make(map[common.Address]model.VoteMsg,0)
	votes[vote1.GetAddress()]=*vote1
	//votes[vote2.GetAddress()]=*vote2


	//verifications := []model.AbstractVerification{&proposal.VoteMsg,vote1,vote2}
	verifications := []model.AbstractVerification{&proposal.VoteMsg,vote1}
	mockChainReader.EXPECT().SaveBlock(&proposal.EmptyBlock,verifications).Return(nil)

	err = stateHandler.SaveFinalEmptyBlock(*proposal,votes)
	assert.NoError(t, err)
}
