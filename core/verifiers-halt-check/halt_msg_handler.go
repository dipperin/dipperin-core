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
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
)

func NewHaltHandler(conf ProposalGeneratorConfig) *VBHaltHandler {
	log.Halt.Info("NewHaltHandler start~~~~~~~~~~~~~~~~~~~~~~~~`")
	return &VBHaltHandler{
		pgConfig:                 conf,
		proposalMessagesByOthers: make([]ProposalMsg, 0),
		aliveVerVotes:            make(map[common.Address]model.VoteMsg),
	}
}

// verifier boot node halt handler
// A new one is generated after each timeout. This way resetting the state is not needed after finishing processing.
type VBHaltHandler struct {
	pgConfig ProposalGeneratorConfig

	// local proposal message
	proposalMsg *ProposalMsg
	// proposal message received from other peers
	proposalMessagesByOthers []ProposalMsg
	// the chosen minimum proposal
	minProposalMsg ProposalMsg

	// the voteMsg from alive verifier
	aliveVerVotes map[common.Address]model.VoteMsg
}

// generate proposalMsg
func (handler *VBHaltHandler) ProposeEmptyBlock() (pm ProposalMsg, err error) {
	log.Halt.Info("VBHaltHandler ProposeEmptyBlock start~~~~~~")
	handler.proposalMsg, err = GenProposalMsg(handler.pgConfig)
	if err != nil {
		log.Halt.Info("VBHaltHandler GenProposalMsg error", "err", err)
		return ProposalMsg{}, err
	}
	pm = *handler.proposalMsg
	handler.minProposalMsg = *handler.proposalMsg

	log.Halt.Info("VBHaltHandler ProposeEmptyBlock end~~~~~~")
	return
}

// receive msg of other VBoot
func (handler *VBHaltHandler) OnNewProposalMsg(msg ProposalMsg) error {
	// check whether this msg is already received
	for _, m := range handler.proposalMessagesByOthers {
		if m.VoteMsg.Witness.Address.IsEqual(msg.VoteMsg.Witness.Address) {
			return g_error.AlreadyHaveVoteMsgError
		}
	}

	log.Halt.Info("the received msg is:", "receivedMsg", msg)
	log.Halt.Info("the own proposal is:", "ownProposal", handler.proposalMsg)
	log.Halt.Info("the own proposal is:", "ownProposal", *handler.proposalMsg)

	// whether the height matches
	if msg.EmptyBlock.Number() != handler.proposalMsg.EmptyBlock.Number() {
		return g_error.EmptyBlockNumberNotMatchError
	}

	//check proposal valid
	if err := checkProposalValid(msg); err != nil {
		return err
	}

	// select the minimal msg
	handler.minProposalMsg = selectEmptyProposal(handler.minProposalMsg, msg)
	log.Halt.Info("the handler minProposalMsg is:", "hash", handler.minProposalMsg.EmptyBlock.Hash().Hex())

	return nil
}

func (handler *VBHaltHandler) GetProposalMsg() *ProposalMsg {
	return handler.proposalMsg
}

func (handler *VBHaltHandler) GetOtherProposalMessages() []ProposalMsg {
	return handler.proposalMessagesByOthers
}

func (handler *VBHaltHandler) GetMinProposalMsg() ProposalMsg {
	return handler.minProposalMsg
}

//handle the received proposals and change the haltHandler status
func (handler *VBHaltHandler) HandlerProposalMessages(msg ProposalMsg, selectedProposal chan ProposalMsg) error {
	chainConfig := chain_config.GetChainConfig()
	err := handler.OnNewProposalMsg(msg)
	if err != nil {
		return err
	}

	handler.proposalMessagesByOthers = append(handler.proposalMessagesByOthers, msg)
	log.Halt.Info("the handler votesLen is:", "len", handler.VotesLen(), "verBootNodeNumber", chainConfig.VerifierBootNodeNumber)
	if handler.VotesLen() == chainConfig.VerifierBootNodeNumber {
		//collect all empty block, send the block with minimal hash to verifier
		//only minimal hash node send the block
		log.Halt.Info("the own proposal block hash is:", "hash", handler.proposalMsg.EmptyBlock.Hash(), "minimalHash", handler.minProposalMsg.EmptyBlock.Hash())
		if handler.proposalMsg.EmptyBlock.Hash() == handler.minProposalMsg.EmptyBlock.Hash() {
			selectedProposal <- handler.minProposalMsg
		}
		return nil
	}
	return g_error.ProposeNotEnough
}

func (handler *VBHaltHandler) HandlerAliveVerVotes(vote model.VoteMsg, currentVerifiers []common.Address) error {
	err := vote.HaltedVoteValid(currentVerifiers)
	if err != nil {
		log.Halt.Error("the aliveVerifierVote received from alive verifier is invalid", "err", err)
		return err
	}

	log.Halt.Info("the vote is:", "vote", vote.GetAddress().Hex())
	//log.Halt.Info("the own proposal is:","ownProposal",*handler.proposalMsg)
	if vote.BlockID != handler.proposalMsg.EmptyBlock.Hash() {
		log.Halt.Error("the vote block hash error")
		return g_error.AliveVoteBlockHashError
	}

	log.Halt.Info("the vote witness address is:", "address", vote.Witness.Address)
	if _, ok := handler.aliveVerVotes[vote.Witness.Address]; !ok {
		handler.aliveVerVotes[vote.Witness.Address] = vote
	}

	return nil
}

func (handler *VBHaltHandler) VotesLen() int {
	return len(handler.proposalMessagesByOthers) + 1
}

func (handler *VBHaltHandler) MinProposalMsg() ProposalMsg {
	return handler.minProposalMsg
}

//alive verifier halt handler
type AliveVerHaltHandler struct {
	signHashFunc     SignHashFunc
	ownAddress       common.Address
	receivedProposal ProposalMsg
	ownVote          model.VoteMsg
}

func NewAliveVerHaltHandler(signFunc SignHashFunc, addr common.Address) *AliveVerHaltHandler {
	return &AliveVerHaltHandler{signHashFunc: signFunc, ownAddress: addr}
}

func (handler *AliveVerHaltHandler) OnMinimalHashBlock(selectedProposal ProposalMsg) (*model.VoteMsg, error) {
	//check minimalHash proposal valid
	if err := checkProposalValid(selectedProposal); err != nil {
		return nil, err
	}

	handler.receivedProposal = selectedProposal

	return GenVoteMsg(&selectedProposal.EmptyBlock, handler.signHashFunc, handler.ownAddress, model.AliveVerifierVoteMessage)

}
