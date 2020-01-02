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
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"time"
)

func GenVoteMsg(emptyBlock *model.Block, signFunc SignHashFunc, addr common.Address, voteType model.VoteMsgType) (*model.VoteMsg, error) {
	//generate empty block verification and send to the verification boot node
	vote := &model.VoteMsg{
		Height:    emptyBlock.Number(),
		BlockID:   emptyBlock.Hash(),
		VoteType:  voteType,
		Timestamp: time.Now(),
	}

	log.DLogger.Info("the voteMsg blockID is", zap.String("BlockID", vote.BlockID.Hex()), zap.Uint64("height", vote.Height))
	// sign msg
	log.DLogger.Info("generate empty vote", zap.Any("address", addr))
	sign, err := signFunc(vote.Hash().Bytes())
	if err != nil {
		log.DLogger.Warn("sign aliveVerifierVote msg failed", zap.Error(err))
		return nil, err
	}
	vote.Witness = &model.WitMsg{
		Address: addr,
		Sign:    sign,
	}

	return vote, nil
}

func checkProposalValid(proposal ProposalMsg) error {

	if proposal.EmptyBlock.Hash() != proposal.VoteMsg.BlockID {
		log.DLogger.Warn("the proposal empty block hash is different from VoteMsg", zap.String("blockHash", proposal.EmptyBlock.Hash().Hex()), zap.String("voteMsgBlockId", proposal.VoteMsg.BlockID.Hex()))
		return g_error.VoteMsgBlockHashNotMatchError
	}

	err := proposal.VoteMsg.HaltedVoteValid([]common.Address{})
	if err != nil {
		log.DLogger.Error("the proposal VoteMsg is invalid", zap.Error(err))
		return err
	}
	return nil
}

//select a emptyBlockProposal from emptyBlockProposals because the verifications is different in each block
func selectEmptyProposal(proposalA, proposalB ProposalMsg) ProposalMsg {
	if proposalA.EmptyBlock.Hash().Hex() < proposalB.EmptyBlock.Hash().Hex() {
		return proposalA
	} else if proposalA.EmptyBlock.Hash().Hex() > proposalB.EmptyBlock.Hash().Hex() {
		return proposalB
	} else {
		panic("the proposalA and BlockB hash isn't different")
	}
}
