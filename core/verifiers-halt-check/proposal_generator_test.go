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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/verifiers-halt-check"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

var bootNodeIndex int
var testVerBootAccounts []tests.Account

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

func MakeTestProposalConfig(voteType model.VoteMsgType, verBootIndex int) verifiers_halt_check.ProposalGeneratorConfig {
	bootNodeIndex = verBootIndex
	return verifiers_halt_check.ProposalGeneratorConfig{
		CurBlock:         factory.CreateBlock(2),
		SignHashFunc:     verBootSignForTest,
		ProcessStateFunc: stateProcess,
		VoteType:         voteType,
		PubKey:           crypto.FromECDSAPub(&testVerBootAccounts[bootNodeIndex].Pk.PublicKey),
	}
}

func MakeTestProposalMsg(verBootIndex int) (*verifiers_halt_check.ProposalMsg, error) {
	config := MakeTestProposalConfig(model.VerBootNodeVoteMessage, verBootIndex)
	return verifiers_halt_check.GenProposalMsg(config)
}

func TestProposalGenerator_GenProposalMsg(t *testing.T) {
	_, err := MakeTestProposalMsg(0)
	assert.NoError(t, err)
}
