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

package state_machine

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"time"
)

var VerifierNumber uint64 = 22

type BftState struct {
	Height      uint64
	Round       uint64
	Step        model2.RoundStepType
	CommitRound uint64

	BlockPoolNotEmpty bool
	CurVerifiers      []common.Address
	LockedBlock       model.AbstractBlock
	LockedRound       uint64

	NewRound      *NewRoundSet
	Proposal      *ProposalSet
	ProposalBlock *BlockSet
	PreVotes      *VoteSet
	Votes         *VoteSet
}

// Events Process
// when a new block is inserted, the height should be changed.
func (bs *BftState) OnNewHeight(height uint64, round uint64, verifiers []common.Address) {
	log.PBft.Debug("New Height Called", "enter height", height)
	bs.enterNewHeight(height, round, verifiers)
}

//When receive a block on this height
func (bs *BftState) OnBlockPoolNotEmpty() {
	bs.BlockPoolNotEmpty = true
	if bs.Step == model2.RoundStepNewHeight {
		bs.enterNewRound(bs.Height, bs.Round)
	}
}

//When receive a NewRound message
func (bs *BftState) OnNewRound(r *model2.NewRoundMsg) {
	log.PBft.Info("[BftState-OnNewRound]", "ownHeight", bs.Height, "ownBlockPoolNotEmpty", bs.BlockPoolNotEmpty)
	if bs.Height != r.Height {
		log.PBft.Warn("height is different", "rHeight", r.Height, "ownHeight", bs.Height)
		return
	}

	if err := bs.NewRound.Add(r); err != nil {
		log.PBft.Info("[BftState-OnNewRound] err", "err", err)
		return
	}

	// check should catch up?
	if bs.Round < bs.NewRound.shouldCatchUpTo() {
		bs.enterNewRound(bs.Height, bs.NewRound.shouldCatchUpTo())
	}

	bs.tryEnterPropose()
}

//When receive a proposal
func (bs *BftState) OnNewProposal(p *model2.Proposal, block model.AbstractBlock) {
	log.PBft.Info("[BftState-OnNewProposal]", "pRound", p.Round, "pHeight", p.Height, "pBlockId", p.BlockID.Hex(), "ownRound", bs.Round, "ownHeight", bs.Height)
	if !bs.validProposal(p) {
		log.PBft.Error("validProposal err")
		return
	}

	if !p.BlockID.IsEqual(block.Hash()) {
		log.PBft.Error("the proposal block error", "p.BlockID", p.BlockID.Hex(), "blockId", block.Hash().Hex())
		return
	}

	// now we have a valid proposal message
	bs.Proposal.Add(p)

	log.PBft.Info("get valid proposal", "height", p.Height, "block", p.BlockID.Hex(), "round", p.Round)
	bs.ProposalBlock.AddBlock(block, p.Round)
	if bs.Step == model2.RoundStepPropose && bs.Round == p.Round {
		bs.enterPreVote(p, block)
	}
}

//When receive a prevote
func (bs *BftState) OnPreVote(pv *model.VoteMsg) {
	log.PBft.Info("[BftState-OnPreVote]")
	if pv.Height != bs.Height {
		log.PBft.Warn("height is different", "pvHeight", pv.Height, "ownHeight", bs.Height)
		return
	}
	if err := bs.PreVotes.AddVote(pv); err != nil {
		log.PBft.Warn("receive invalid pre vote", "err", err)
		return
	}

	roundBlock := bs.PreVotes.VotesEnough(pv.Round)
	//fmt.Println("onprevote","who",reflect.ValueOf(bs).Pointer(),"pv",roundBlock)
	log.PBft.Info("the round info is:", "pv.Round", pv.Round, "bs.Round", bs.Round, "bs.LockedRound", bs.LockedRound)
	// Release block lock
	if bs.LockedBlock != nil && !roundBlock.IsEqual(common.Hash{}) && pv.Round >= bs.Round && bs.LockedRound < pv.Round {
		log.PBft.Info("the locked and round Block is:", "lockedBlock", bs.LockedBlock.Hash().Hex(), "roundBlock", roundBlock.Hex())
		//Fixme Deleted code: if !bs.LockedBlock.Hash().IsEqual(roundBlock){bs.LockedBlock = nil}
		//When bs.Round = 13, LockedRound=13, lockedBlock=X. Receive 2/3 vote on round 14, block X. Should unlock X, and lock X after.
		bs.LockedBlock = nil
	}

	// Add lock
	// Fixme if pv.Round > bs.Round should I do this?
	if !roundBlock.IsEqual(common.Hash{}) && bs.LockedBlock == nil {
		block := bs.ProposalBlock.GetBlock(pv.Round)
		if block != nil {
			if block.Hash().IsEqual(roundBlock) {
				bs.LockedBlock = block
				bs.LockedRound = pv.Round
				log.PBft.Debug("[BftState-LockBlock]", "LockedRound", bs.LockedRound, "block", block.Hash().Hex())
				bs.enterPreCommit(pv.Round)
			}
		} else {
			log.PBft.Error("the proposal Block is nil")
		}
	}
}

//receive vote msg
func (bs *BftState) OnVote(v *model.VoteMsg) (common.Hash, []model.AbstractVerification) {
	log.PBft.Info("[BftState-OnVote]")
	if v.Height != bs.Height {
		log.PBft.Error("[BftState-OnVote] ignore wrong height vote", "voteHeight", v.Height, "ownHeight", bs.Height)
		return common.Hash{}, nil
	}

	// Add a valid vote
	if err := bs.Votes.AddVote(v); err != nil {
		log.PBft.Error("add vote err", "err", err)
		return common.Hash{}, nil
	}

	// check have enough votes, reject very past commits
	maj32Block := bs.Votes.VotesEnough(v.Round)
	if maj32Block.IsEqual(common.Hash{}) || v.Round < bs.Round {
		log.PBft.Error("[BftState-OnVote] cannot final block", "maj32Block == nil", maj32Block.IsEqual(common.Hash{}), "v.Round < bs.Round", v.Round < bs.Round)
		return common.Hash{}, nil
	}

	// select correct voteMsgs
	resultVotes := bs.Votes.FinalVerifications(v.Round)
	log.PBft.Info("[BftState-OnVote] can final block", "voteRound", v.Round, "ownRound", bs.Round)
	return maj32Block, resultVotes
}

// State Change

/*
Enter New Height
Called by: OnNewHeight, OnVote
Enter:
	- When download block from others, enter a new height.
	- When commit a block on last height, enter a new height
*/
func (bs *BftState) enterNewHeight(newHeight uint64, round uint64, verifiers []common.Address) {
	log.PBft.Info(fmt.Sprintf("[EnterNewHeight], last height info (H: %v, R: %v), new height (H: %v, R: %v)", bs.Height, bs.Round, newHeight, round))
	log.Health.Info("bft state enter new height", "h", bs.Height, "r", bs.Round)

	bs.Height = newHeight
	bs.Round = round
	bs.Step = model2.RoundStepNewHeight
	bs.BlockPoolNotEmpty = false
	bs.CurVerifiers = verifiers
	bs.NewRound = NewNRoundSet(newHeight, verifiers)
	bs.Proposal = NewProposalSet()
	bs.ProposalBlock = NewBlockSet()
	bs.PreVotes = NewVoteSet(newHeight, verifiers)
	bs.Votes = NewVoteSet(newHeight, verifiers)
	bs.LockedBlock = nil
	bs.LockedRound = uint64(0)
}

/*
Enter New Round
*/
func (bs *BftState) enterNewRound(height, round uint64) {
	bs.Step = model2.RoundStepNewRound
	// if the status is moved from new height, then the state should not be added by 1, otherwise it would not be the 0th verifier who does the proposition of block
	if height != bs.Height {
		//panic("Shouldn't call new round before new height")
		log.PBft.Info("[BftState-enterNewRound]:error", "height", height, "ownHeight", bs.Height)
		return
	}
	log.PBft.Info(fmt.Sprintf("[EnterNewRound], (H: %v, R: %v, S: %v)", bs.Height, bs.Round, bs.Step))
	bs.Round = round
}

/*
Enter propose
Called by: tryStartNewRound, OnNewProposal
Enter:
- have 32 majority agree the round
- received a right proposal
*/
func (bs *BftState) enterPropose(round uint64) {
	bs.Step = model2.RoundStepPropose
	log.PBft.Debug(fmt.Sprintf("[EnterPropose] (H: %v, R: %v, S:%v)", bs.Height, bs.Round, bs.Step))
	log.Health.Info("bft state enter new propose", "h", bs.Height, "r", bs.Round)
	//Do propose outside
	//Already get the proposal?
	proposal := bs.Proposal.GetProposal(round)
	proposalBlock := bs.ProposalBlock.Blocks[round]
	if proposal != nil && proposalBlock != nil {
		if proposal.BlockID.IsEqual(proposalBlock.Hash()) {
			bs.enterPreVote(proposal, proposalBlock)
			return
		}
	}
}

/*
Enter prevote
Enter:
- when received proposal and block
?? 3 prevotes received before reception of propose
*/
func (bs *BftState) enterPreVote(p *model2.Proposal, b model.AbstractBlock) {
	if p == nil || bs.Round != p.Round || bs.Height != p.Height || bs.Step != model2.RoundStepPropose {
		return
	}
	log.Health.Info("bft state enter new pre vote", "h", bs.Height, "r", bs.Round)
	// Update state
	bs.Step = model2.RoundStepPreVote
	log.PBft.Info(fmt.Sprintf("[EnterPreVote], (H: %v, R: %v, S: %v)", bs.Height, bs.Round, bs.Step))
}

/*
Enter Precommit
Called by: onPrevote
Enter:
- Get 2/3 majority vote.
*/
func (bs *BftState) enterPreCommit(round uint64) {
	if bs.Round != round || bs.Step != model2.RoundStepPreVote {
		return
	}
	log.Health.Info("bft state enter new pre commit", "h", bs.Height, "r", bs.Round)
	bs.Step = model2.RoundStepPreCommit
	log.PBft.Info(fmt.Sprintf("[EnterPrecommit], (H: %v, R: %v, S: %v)", bs.Height, bs.Round, bs.Step))
}

// Timeout actions

/*
1. majority not agree with the box -> next round directly
2. majority offline -> next round directly
3. when I am offline but others reach consensus -> enter new round, but cannot receive the corresponding proposal. wait for downloader for a sync
4. when I am offline and others cannot reach consensus -> enter new round and continue the current mechanism
*/
func (bs *BftState) OnProposeTimeout() {
	log.PBft.Debug("StateHandler#OnProposeTimeout Propose get time out", "h.roundState.Height", bs.Height, "h.roundState.Round", bs.Round, "h.roundState.Step", bs.Step.String())
	if bs.Step == model2.RoundStepPropose {
		bs.enterNewRound(bs.Height, bs.Round+1)
	}
}

func (bs *BftState) OnPreVoteTimeout() {
	log.PBft.Info("StateHandler#OnPreVoteTimeout enterNewRound", "h.roundState.Height", bs.Height, "h.roundState.Round", bs.Round)
	bs.enterNewRound(bs.Height, bs.Round+1)
}

func (bs *BftState) OnPreCommitTimeout() {
	log.PBft.Info("StateHandler#OnPreCommitTimeout enterNewRound", "h.roundState.Height", bs.Height, "h.roundState.Round", bs.Round)
	bs.enterNewRound(bs.Height, bs.Round+1)
}

func (bs *BftState) getNewRoundReqList() []common.Address {

	//Fixme
	if bs.Step != model2.RoundStepNewRound {
		return nil
	}

	byteArray := bs.NewRound.MissingAtRound(bs.Round)
	var reqAddresses []common.Address
	for i := range bs.CurVerifiers {
		if !byteArray.GetIndex(i) {
			reqAddresses = append(reqAddresses, bs.CurVerifiers[i])
		}
	}
	return reqAddresses
}

//Fixme can i receive proposal of other round?
func (bs *BftState) validProposal(p *model2.Proposal) bool {
	// check height and round
	// we accept bs.Round <= p.Round < bs.Round + 10
	if bs.Height != p.Height || bs.Round > p.Round || bs.Round+10 < p.Round {
		log.PBft.Warn("check height and round error", "height", bs.Height, "round", bs.Round, "proposal height", p.Height, "proposal round", p.Round)
		return false
	}
	// check already have proposal for this round
	if bs.Proposal.Have(p.Round) {
		log.PBft.Warn("BftState#OnNewProposal  already have proposal in this round", "round", p.Round, "cur h", bs.Height)
		return false
	}
	// valid witness
	if err := p.Witness.Valid(p.Hash().Bytes()); err != nil {
		log.PBft.Warn("BftState#OnNewProposal receive invalid proposal", "p.Hash().Bytes()", p.Hash().Bytes(), "err", err)
		return false
	}

	// valid proposer
	if !bs.proposerAtRound(p.Round).IsEqual(p.Witness.Address) {
		log.PBft.Warn("BftState#OnNewProposal  receive invalid proposal", "cur proposer", bs.curProposer(), "proposal addr", p.Witness.Address)
		return false
	}
	return true
}

// Checked
func (bs *BftState) tryEnterPropose() {
	log.PBft.Debug(fmt.Sprintf("[BftState-tryEnterPropose] (H: %v, R: %v, S:%v)", bs.Height, bs.Round, bs.Step))

	if !bs.BlockPoolNotEmpty {
		log.PBft.Debug(fmt.Sprintf("[BftState-tryEnterPropose] failed 1 (H: %v, R: %v, S:%v)", bs.Height, bs.Round, bs.Step))
		return
	}

	if bs.Step != model2.RoundStepNewRound {
		log.PBft.Debug(fmt.Sprintf("[BftState-tryEnterPropose] failed 2 (H: %v, R: %v, S:%v)", bs.Height, bs.Round, bs.Step))
		return
	}

	if !bs.NewRound.enoughAtRound(bs.Round) {
		log.PBft.Debug(fmt.Sprintf("[BftState-tryEnterPropose] failed 3 (H: %v, R: %v, S:%v)", bs.Height, bs.Round, bs.Step))
		return
	}
	bs.enterPropose(bs.Round)
}

//New functions
func (bs *BftState) makePrevote() (msg *model.VoteMsg) {
	if bs.LockedBlock != nil {
		return &model.VoteMsg{
			Height:    bs.Height,
			Round:     bs.Round,
			BlockID:   bs.LockedBlock.Hash(),
			VoteType:  model.PreVoteMessage,
			Timestamp: time.Now(),
		}
	}

	proposal := bs.Proposal.GetProposal(bs.Round)
	block := bs.ProposalBlock.GetBlock(bs.Round)
	if proposal == nil || block == nil {
		return nil
	}
	if !block.Hash().IsEqual(proposal.BlockID) {
		return nil
	}

	return &model.VoteMsg{
		Height:    bs.Height,
		Round:     bs.Round,
		BlockID:   block.Hash(),
		VoteType:  model.PreVoteMessage,
		Timestamp: time.Now(),
	}
}
func (bs *BftState) makeVote() (msg *model.VoteMsg) {
	if bs.LockedBlock == nil {
		return nil
	}
	return &model.VoteMsg{
		Height:    bs.Height,
		Round:     bs.Round,
		BlockID:   bs.LockedBlock.Hash(),
		VoteType:  model.VoteMessage,
		Timestamp: time.Now(),
	}
}
func (bs *BftState) curProposer() (result common.Address) {
	return bs.proposerAtRound(bs.Round)
}

func (bs *BftState) proposerAtRound(round uint64) common.Address {
	vLen := len(bs.CurVerifiers)
	if vLen == 0 {
		panic("No verifiers")
		return common.Address{}
	}
	index := int(round) % vLen
	return bs.CurVerifiers[index]
}
