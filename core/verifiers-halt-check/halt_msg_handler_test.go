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

package verifiers_halt_check_test

import (
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/verifiers-halt-check"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testVerAccounts []tests.Account
var verNodeIndex int

func init() {
	log.Info("change ver boot node address for test")
	var err error
	testVerAccounts, err = tests.ChangeVerifierAddress(nil)
	if err != nil {
		panic("change verifier boot node address error for test")
	}
}

func verSignForTest(hash []byte) ([]byte, error) {
	sk := testVerAccounts[verNodeIndex].Pk
	return crypto.Sign(hash, sk)
}

func MakeTestHaltHandler(verBootIndex int) *verifiers_halt_check.VBHaltHandler {
	config := MakeTestProposalConfig(model.VerBootNodeVoteMessage, verBootIndex)
	return verifiers_halt_check.NewHaltHandler(config)
}

func ProposeEmptyBlockForTest(verBootIndex int) (handler *verifiers_halt_check.VBHaltHandler, err error) {
	testHaltHandler := MakeTestHaltHandler(verBootIndex)
	_, err = testHaltHandler.ProposeEmptyBlock()
	if err != nil {
		return nil, err
	}

	return testHaltHandler, nil
}

func TestVBHaltHandler_ProposeEmptyBlock(t *testing.T) {
	_, err := ProposeEmptyBlockForTest(0)
	assert.NoError(t, err)
}

func TestVBHaltHandler_OnNewProposalMsg(t *testing.T) {
	testMsg, err := MakeTestProposalMsg(1)
	assert.NoError(t, err)

	testHandler, err := ProposeEmptyBlockForTest(0)
	assert.NoError(t, err)

	err = testHandler.OnNewProposalMsg(*testMsg)
	assert.NoError(t, err)
}

func TestVBHaltHandler_HandlerProposalMessages(t *testing.T) {
	chainConfig := chain_config.GetChainConfig()
	testMessages := make([]verifiers_halt_check.ProposalMsg, chainConfig.VerifierBootNodeNumber-1)

	//generate other proposal msg
	var err error
	var testMinProposal verifiers_halt_check.ProposalMsg
	for i := 0; i < (chainConfig.VerifierBootNodeNumber - 1); i++ {
		proposal, err := MakeTestProposalMsg(i + 1)
		if i == 0 {
			testMinProposal = *proposal
		} else {
			if testMinProposal.EmptyBlock.Hash().Hex() > proposal.EmptyBlock.Hash().Hex() {
				testMinProposal = *proposal
			}
		}

		assert.NoError(t, err)
		testMessages[i] = *proposal
		//log.Info("the other proposal hash is:","hash",proposal.EmptyBlock.Hash().Hex())
	}

	//propose empty block
	testHandler, err := ProposeEmptyBlockForTest(0)
	assert.NoError(t, err)
	if testMinProposal.EmptyBlock.Hash().Hex() > testHandler.GetProposalMsg().EmptyBlock.Hash().Hex() {
		testMinProposal = *testHandler.GetProposalMsg()
	}
	//log.Info("the own proposal hash is:","hash",testHandler.GetProposalMsg().EmptyBlock.Hash().Hex())

	//test handler proposal messages
	testChan := make(chan verifiers_halt_check.ProposalMsg, 0)
	go func() {
		for _, msg := range testMessages {
			testHandler.HandlerProposalMessages(msg, testChan)
		}
	}()

	readProposal := verifiers_halt_check.ProposalMsg{}
Loop:
	for {
		select {
		case readProposal = <-testChan:
			break Loop
		case <-time.After(100 * time.Millisecond):
			break Loop
		}
	}

	//assert the proposal message handler result
	minimalProposal := testHandler.GetMinProposalMsg()
	ownProposal := *testHandler.GetProposalMsg()
	otherProposals := testHandler.GetOtherProposalMessages()

	assert.EqualValues(t, chainConfig.VerifierBootNodeNumber-1, len(otherProposals))
	assert.EqualValues(t, testMinProposal.EmptyBlock.Hash().Hex(), minimalProposal.EmptyBlock.Hash().Hex())
	nilProposal := verifiers_halt_check.ProposalMsg{}
	if readProposal != nilProposal {
		assert.EqualValues(t, ownProposal, readProposal)
		assert.EqualValues(t, ownProposal, minimalProposal)
	}
}

func TestVBHaltHandler_HandlerAliveVerVotes(t *testing.T) {
	testHandler, err := ProposeEmptyBlockForTest(0)
	assert.NoError(t, err)

	testVerVote, err := verifiers_halt_check.GenVoteMsg(&testHandler.GetProposalMsg().EmptyBlock, verSignForTest, testVerAccounts[verNodeIndex].Address(), model.AliveVerifierVoteMessage)
	assert.NoError(t, err)

	err = testHandler.HandlerAliveVerVotes(*testVerVote, chain.VerifierAddress)
	assert.NoError(t, err)
}

func TestAliveVerHaltHandler_OnMinimalHashBlock(t *testing.T) {
	testHandler, err := ProposeEmptyBlockForTest(0)
	assert.NoError(t, err)

	testAliveVerHandler := verifiers_halt_check.NewAliveVerHaltHandler(verSignForTest, testVerAccounts[verNodeIndex].Address())

	_, err = testAliveVerHandler.OnMinimalHashBlock(*testHandler.GetProposalMsg())
	assert.NoError(t, err)
}

type A struct {
	B *int
	C int
}

/*func TestGoPriority(t *testing.T) {
	testB := 0
	testA := &A{
		B: &testB,
		C: 1,
	}

	log.Info("the *testA.B is:","value",*testA.B)
	log.Info("the *testA.B is:","value",(*testA).B)
}
*/
