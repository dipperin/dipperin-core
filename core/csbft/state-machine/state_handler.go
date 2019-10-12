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
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"time"
)

type StateHandler struct {
	util.BaseService
	*BftConfig
	timeoutConfig Config

	bs        *BftState
	blockPool *components.BlockPool
	ticker    components.TimeoutTicker

	newHeightChan        chan uint64
	newRoundChan         chan *model2.NewRoundMsg
	poolNotEmptyChan     chan struct{}
	newProposalChan      chan *model2.Proposal
	preVoteChan          chan *model.VoteMsg
	voteChan             chan *model.VoteMsg
	getProposalBlockChan chan getProposalBlockMsg
}

type BftConfig struct {
	//FetcherConnAdaptCsBft components.FetcherConn
	ChainReader ChainReader
	Fetcher     Fetcher
	Signer      MsgSigner
	Sender      MsgSender
	Validator   Validator
}

type ReqRoundMsg struct {
	Height uint64
	Round  uint64
}

type getProposalBlockMsg struct {
	hash       common.Hash
	resultChan chan model.AbstractBlock
}

func NewStateHandler(bftConfig *BftConfig, timeConfig Config, blockPool *components.BlockPool) *StateHandler {
	h := &StateHandler{
		bs:                   &BftState{Height: 0, BlockPoolNotEmpty: false},
		timeoutConfig:        timeConfig,
		blockPool:            blockPool,
		newHeightChan:        make(chan uint64, 5),
		newRoundChan:         make(chan *model2.NewRoundMsg, 5),
		poolNotEmptyChan:     make(chan struct{}, 5),
		newProposalChan:      make(chan *model2.Proposal, 5),
		preVoteChan:          make(chan *model.VoteMsg, 5),
		voteChan:             make(chan *model.VoteMsg, 5),
		getProposalBlockChan: make(chan getProposalBlockMsg),
		ticker:               components.NewTimeoutTicker(),
	}
	h.BftConfig = bftConfig

	//Setup service
	h.BaseService = *util.NewBaseService(log.Root(), "cs_bft", h)

	//Update state
	curHeight := h.ChainReader.CurrentBlock().Number()
	h.OnNewHeight(curHeight + 1)

	return h
}

func (h *StateHandler) OnStart() error {
	log.PBft.Info("StateHandler OnStart~~~~~~~~~~~~~~~~~")
	h.ticker = components.NewTimeoutTicker()
	h.ticker.Start()
	go h.loop()
	return nil
}

func (h *StateHandler) OnStop() {
	h.ticker.Stop()
}

func (h *StateHandler) OnReset() error { return nil }

func (h *StateHandler) loop() {
	for {
		select {
		// outer events, each method of processing an event must determine whether it needs to be executed according to the current state.
		// For example, if you receive a 2/3 prevote and change the status, then you should not execute the same logic if you receive the prevote again.
		case height := <-h.newHeightChan:
			h.OnNewHeight(height)
			// trigger for start
		case <-h.poolNotEmptyChan:
			h.OnBlockPoolNotEmpty()
		case proposal := <-h.newProposalChan:
			h.OnNewProposal(proposal)
		case nRound := <-h.newRoundChan:
			h.OnNewRound(nRound)
		case pv := <-h.preVoteChan:
			h.OnPreVote(pv)
		case v := <-h.voteChan:
			h.OnVote(v)
			// timeout event
		case toutInfo := <-h.ticker.Chan():
			h.OnTimeout(toutInfo)
		case m := <-h.getProposalBlockChan:
			h.onGetProposalBlock(m)
		case <-h.Quit():
			log.PBft.Info("state handler stopped")
			return
		}
	}
}

func (h *StateHandler) OnNewHeight(height uint64) {
	log.PBft.Info("[**********************start new Block************************]")
	log.PBft.Info("[StateHandler-OnNewHeight]", "height", height)
	round := uint64(0)
	chainHeight := h.ChainReader.CurrentBlock().Number()
	if height != chainHeight+1 {
		return
	}

	if height > 1 {
		seenCommit := h.ChainReader.GetSeenCommit(chainHeight)
		if len(seenCommit) == 0 {
			log.Error(g_error.ErrCannotLoadSeenCommit.Error(), "chainHeight", chainHeight)
			return
		}
		round = h.ChainReader.GetSeenCommit(chainHeight)[0].GetRound()
	}

	h.blockPool.NewHeight(height)

	Block := h.ChainReader.CurrentBlock()
	log.PBft.Debug("New Height Called", "height", height, "chain height", Block.Number())
	// check where it is a change point, add verifiers list and set round as 0
	if h.ChainReader.IsChangePoint(Block, false) {
		verifiers := h.ChainReader.GetNextVerifiers()
		h.bs.OnNewHeight(height, 0, verifiers)
		return
	}

	h.bs.OnNewHeight(height, round+1, h.ChainReader.GetCurrVerifiers())
	log.PBft.Debug(fmt.Sprintf("EnterNewHeight (H: %v, R: %v, S: %v)", h.bs.Height, h.bs.Round, h.bs.Step))
}

func (h *StateHandler) OnNewRound(nRound *model2.NewRoundMsg) {
	log.PBft.Info("[StateHandler-OnNewRound]", "address", nRound.Witness.Address.Hex(), "round", nRound.Round, "Height", nRound.Height)
	_, preRound, preStep := h.RecordCurState()
	h.bs.OnNewRound(nRound)
	_, curRound, curStep := h.RecordCurState()
	fmt.Println(fmt.Sprintf("New round msg, state change from (R:%v, S:%s) to (R:%v, S:%s)",preRound, preStep,curRound, curStep))

	// A new round message, lead to state machine change state not more than 3 times.
	// 1. Catch up with new round.
	if preRound < curRound {
		h.onEnterNewRound()
		switch curStep {
		case model2.RoundStepPropose:
			h.onEnterPropose()
		case model2.RoundStepPreVote:
			h.onEnterPrevote() // ignore enter propose, because state machine already get the right proposal
		default:
			panic("State change error")
		}
	} else { // 2. preRound = curRound
		switch {
		case preStep == curStep: // state not change, do not need send out message
		case curStep == model2.RoundStepPropose: // preStep != curStep
			h.onEnterPropose()
		case curStep == model2.RoundStepPreVote: // preStep != curStep
			h.onEnterPrevote() // ignore enter propose, because state machine already get the right proposal
		default:
			panic("State change error")
		}
	}
	/*switch {
	case preStep == model2.RoundStepNewRound && curStep == model2.RoundStepPropose:
		//fmt.Println("enter","id",reflect.ValueOf(h.bs).Pointer())
		if preRound != curRound {
			log.PBft.Info("[StateHandler-OnNewRound]:onEnterNewRound catch up round")
		}
		log.PBft.Info("[StateHandler-OnNewRound]:onEnterPropose", "pre", preStep, "new", curStep)
		h.onEnterPropose()
	case preStep == model2.RoundStepNewRound && curStep == model2.RoundStepPreVote:
		log.PBft.Info("[StateHandler-OnNewRound]: onEnterPrevote", "pre", preStep, "new", curStep)
		h.onEnterPrevote()
	default:
		log.PBft.Info("[StateHandler-OnNewRound]:on new round", "pre", preStep, "new", curStep)
	}*/
}

func (h *StateHandler) OnBlockPoolNotEmpty() {
	log.PBft.Info("[StateHandler-OnBlockPoolNotEmpty]")
	_, preRound, preStep := h.RecordCurState()
	h.bs.OnBlockPoolNotEmpty()
	_, curRound, curStep := h.RecordCurState()
	fmt.Println(fmt.Sprintf("Block pool not empty, state change from (R:%v, S:%s) to (R:%v, S:%s)",preRound, preStep,curRound, curStep))


	//Fixme (R:3, S:NewHeight) if get 2/3 vote on round 4 before block pool not empty; state machine jumped into NewRound step, and ignore this msg
	switch {
	case preStep == model2.RoundStepNewHeight && curStep == model2.RoundStepNewRound:
		log.PBft.Info("[StateHandler-OnBlockPoolNotEmpty]:onEnterNewRound")
		h.onEnterNewRound()
	default:
		log.PBft.Info("[StateHandler-OnBlockPoolNotEmpty]:block pool not empty", "pre", preStep, "new", curStep)
	}
}

func (h *StateHandler) OnNewProposal(proposal *model2.Proposal) {
	log.PBft.Info("[StateHandler-OnNewProposal]", "block", proposal.BlockID.Hex())
	if !h.bs.validProposal(proposal) {
		return
	}
	log.PBft.Info("[StateHandler-OnNewProposal] proposal accepted, try fetching block", "block", proposal.BlockID.Hex())
	block := h.fetchProposalBlock(proposal.BlockID, proposal.Witness.Address)
	if block == nil || block.IsSpecial() {
		log.PBft.Info("[StateHandler-OnNewProposal] fetch block failed", "block", proposal.BlockID.Hex())
		return
	}

	if err := h.Validator.FullValid(block); err != nil {
		log.PBft.Info("[StateHandler-OnNewProposal] proposed block not valide", "block", proposal.BlockID.Hex())
		return
	}

	if h.blockPool.IsEmpty() == true {
		h.blockPool.AddBlock(block)
	}

	_, preRound, preStep := h.RecordCurState()
	h.bs.OnNewProposal(proposal, block)
	_, curRound, curStep := h.RecordCurState()
	fmt.Println(fmt.Sprintf("Receive proposal, state change from (R:%v, S:%s) to (R:%v, S:%s)",preRound, preStep,curRound, curStep))

	switch {
	case preRound != curRound:
		panic("receive a proposal should not lead round change")
	case preStep == curStep:
	case preStep == model2.RoundStepPropose && curStep == model2.RoundStepPreVote: // preRound == curRound
		h.onEnterPrevote()
	case preStep == model2.RoundStepNewHeight && curStep == model2.RoundStepPreVote:
		h.onEnterPrevote()
	case preStep == model2.RoundStepNewHeight && curStep == model2.RoundStepPropose:
		h.onEnterPropose()
	default:
		fmt.Println(fmt.Sprintf("State change from (R:%s, S:%s) to (R:%s, S:%s)",preRound, preStep,curRound, curStep))
		panic("unexpected state change")
	}

	//fixme Add other cases
	/*switch {
	case preStep == model2.RoundStepPropose && curStep == model2.RoundStepPreVote:
		log.PBft.Info("[StateHandler-OnNewProposal]:onEnterPrevote")
		h.onEnterPrevote()
	default:
		log.PBft.Info("[StateHandler-OnNewProposal]:block pool not empty", "pre", preStep, "new", curStep)
	}*/
}

func (h *StateHandler) OnPreVote(pv *model.VoteMsg) {
	log.PBft.Info("[StateHandler-OnPreVote]")
	_, preRound, preStep := h.RecordCurState()
	h.bs.OnPreVote(pv)
	_, curRound, curStep := h.RecordCurState()
	fmt.Println(fmt.Sprintf("Prevote, state change from (R:%v, S:%s) to (R:%v, S:%s)",preRound, preStep,curRound, curStep))

	switch {
	case preRound == curRound && preStep == model2.RoundStepPreVote && curStep == model2.RoundStepPreCommit:
		h.onEnterPrecommit()
	default:
		log.PBft.Info("[StateHandler-OnPreVote]:on prevote", "pre", preStep, "new", curStep)
	}
}

func (h *StateHandler) OnVote(v *model.VoteMsg) {
	log.PBft.Debug("[StateHandler-OnVote]: handle new vote")
	blockId, commits := h.bs.OnVote(v)
	if commits != nil {
		log.PBft.Info("the commit 0 is:", "round", commits[0].GetRound(), "height", commits[0].GetHeight(), "blockId", commits[0].GetBlockId().Hex(), "address", commits[0].GetAddress().Hex())
		block := h.bs.ProposalBlock.GetBlock(commits[0].GetRound())

		if block == nil {
			log.PBft.Info("not find the block in the state proposalBlockSet ")
			block = h.fetchProposalBlock(blockId, h.bs.proposerAtRound(commits[0].GetRound()))
		}

		if block != nil {
			log.PBft.Info("[StateHandler-OnVote]:finalBlock", "blockNumber", block.Number())
			h.finalBlock(block, commits)
			log.PBft.Info("==================================pbft save block end=======================================")
			log.PBft.Info("")
		}
	}
}

func (h *StateHandler) OnTimeout(toutInfo components.TimeoutInfo) {
	log.PBft.Info("[StateHandler-OnTimeout]", "height", toutInfo.Height, "round", toutInfo.Round, "step", toutInfo.Step, "duration", toutInfo.Duration)
	if toutInfo.Height != h.bs.Height || toutInfo.Round != h.bs.Round {
		return
	}
	log.Health.Info("pbft state handler timeout", "height", toutInfo.Height, "round", toutInfo.Round, "step", toutInfo.Step, "duration", toutInfo.Duration)
	//log.PBft.Info("Get timeout", "timeout info", toutInfo)
	switch toutInfo.Step {
	case model2.RoundStepPropose:
		h.addTimeoutCount("RoundStepPropose")
		h.onProposeTimeout()

	case model2.RoundStepPreVote:
		h.addTimeoutCount("RoundStepPreVote")
		h.onPreVoteTimeout()

	case model2.RoundStepPreCommit:
		h.addTimeoutCount("RoundStepPreCommit")
		h.onPreCommitTimeout()

	case model2.RoundStepNewRound:
		h.addTimeoutCount("RoundStepNewRound")
		h.onNewRoundTimeout()

	default:
		panic(fmt.Sprintf("unknown timeout type: %v code: %v", toutInfo.Step.String(), int(toutInfo.Step)))
	}
}

// broadcast new round msg
func (h *StateHandler) broadcastNewRoundMsg() {
	msg := &model2.NewRoundMsg{
		Height: h.bs.Height,
		Round:  h.bs.Round,
	}

	sign, err := h.BftConfig.Signer.SignHash(msg.Hash().Bytes())
	if err != nil {
		log.Warn("sign new round msg failed", "err", err)
		return
	}
	msg.Witness = &model.WitMsg{
		Address: h.BftConfig.Signer.GetAddress(),
		Sign:    sign,
	}

	log.PBft.Info("StateHandler#broadcastNewRoundMsg")
	h.Sender.BroadcastMsg(uint64(model2.TypeOfNewRoundMsg), msg)
	h.OnNewRound(msg)
}

// todo mv to BftState
func (h *StateHandler) curProposer() (result common.Address) {
	vLen := len(h.bs.CurVerifiers)
	if vLen == 0 {
		return
	}
	index := int(h.bs.Round) % vLen
	return h.bs.CurVerifiers[index]
}

// broadcast new round msg
func (h *StateHandler) broadcastReqRoundMsg(adds []common.Address) {
	msg := &ReqRoundMsg{
		Height: h.bs.Height,
		Round:  h.bs.Round,
	}
	h.Sender.SendReqRoundMsg(uint64(model2.TypeOfReqNewRoundMsg), adds, msg)
}

func (h *StateHandler) fetchProposalBlock(blockId common.Hash, from common.Address) (block model.AbstractBlock) {

	// get proposal block in pool
	block = h.blockPool.GetBlockByHash(blockId)
	if block != nil {
		return
	}

	block = h.Fetcher.FetchBlock(from, blockId)
	return block
}

func (h *StateHandler) onEnterNewRound() {
	h.ticker.ScheduleTimeout(components.TimeoutInfo{Duration: h.timeoutConfig.WaitNewRound, Height: h.bs.Height, Round: h.bs.Round, Step: model2.RoundStepNewRound})

	log.PBft.Debug(fmt.Sprintf("EnterNewRound (H: %v, R: %v, S: %v)", h.bs.Height, h.bs.Round, h.bs.Step))
	h.broadcastNewRoundMsg()
}

func (h *StateHandler) onEnterPropose() {
	h.ticker.ScheduleTimeout(components.TimeoutInfo{Duration: h.timeoutConfig.ProposalTimeout, Height: h.bs.Height, Round: h.bs.Round, Step: model2.RoundStepPropose})

	log.PBft.Debug(fmt.Sprintf("EnterPropose (H: %v, R: %v, S: %v)", h.bs.Height, h.bs.Round, h.bs.Step))

	//fmt.Println("iscu","id",reflect.ValueOf(h.bs).Pointer(),"round",h.bs.Round,"iscu",h.isCurProposer())
	if !h.isCurProposer() {
		return
	}

	//fmt.Println("on enter proposal","id",reflect.ValueOf(h.bs).Pointer())
	//Pick a valid block
	block := h.bs.LockedBlock
	//log.PBft.Info("[onEnterPropose] the block is:","block",block)
	if block == nil {

		block = h.blockPool.GetProposalBlock()
		for block != nil && !block.IsSpecial() {
			err := h.Validator.FullValid(block)
			if err == nil {
				log.PBft.Info(fmt.Sprintf("StateHandler#broadcastProposal  Get a good block from pool. CurMiss (H: %v, R: %v, S:%s)", h.bs.Height, h.bs.Round, h.bs.Step))
				break
			} else {
				log.PBft.Error("StateHandler#broadcastProposal  valid propose block failed", "block hash", block.Hash().Hex(), "result", err)
			}
			block = h.blockPool.GetProposalBlock()
		}

		//No valid block in pool
		if block == nil {
			log.PBft.Info(fmt.Sprintf("StateHandler#broadcastProposal  No valid block in pool, stop propose. CurMiss (H: %v, R: %v, S:%s)", h.bs.Height, h.bs.Round, h.bs.Step))
			return
		}
	}

	//log.PBft.Info("[onEnterPropose] get the proposal block is:","block",block)

	msg := model2.Proposal{
		Height:    h.bs.Height,
		Round:     h.bs.Round,
		BlockID:   block.Hash(),
		Timestamp: time.Now(),
	}
	sign, err := h.BftConfig.Signer.SignHash(msg.Hash().Bytes())
	if err != nil {
		log.PBft.Warn("sign new round msg failed", "err", err)
		return
	}
	msg.Witness = &model.WitMsg{
		Address: h.BftConfig.Signer.GetAddress(),
		Sign:    sign,
	}

	//Send proposal to other verifiers
	h.Sender.BroadcastMsg(uint64(model2.TypeOfProposalMsg), msg)

	// Send to myself the proposal
	h.blockPool.AddBlock(block)
	h.OnNewProposal(&msg)
	/*preStep := h.bs.Step
	h.bs.OnNewProposal(&msg, block)
	curStep := h.bs.Step

	switch {
	case preStep == model2.RoundStepPropose && curStep == model2.RoundStepPreVote:
		h.onEnterPrevote()
	default:
		log.PBft.Info("propose time out", "pre", preStep, "new", curStep)
	}*/
}

//New functions
func (h *StateHandler) finalBlock(block model.AbstractBlock, commits []model.AbstractVerification) {
	log.Health.Info("enter final block", "num", block.Number())
	err := h.ChainReader.SaveBlock(block, commits)
	if err != nil {
		log.Health.Warn("pbft save block failed", "err", err)
		if err.Error() != g_error.ErrAlreadyHaveThisBlock.Error() {
			return
		}
	}
	log.Health.Info("pbft save block success, broadcast it", "block", block.Number())
	// broadcast result
	h.Sender.BroadcastEiBlock(block)
	// change to new height, clear block pool
	// h.OnNewHeight(h.ChainReader.CurrentBlock().Number() + 1)
}

func (h *StateHandler) onEnterPrevote() {
	h.ticker.ScheduleTimeout(components.TimeoutInfo{Duration: h.timeoutConfig.WaitNewRound, Height: h.bs.Height, Round: h.bs.Round, Step: model2.RoundStepPreVote})
	voteMsg := h.bs.makePrevote()

	log.PBft.Debug(fmt.Sprintf("EnterPrevote (H: %v, R: %v, S: %v)", h.bs.Height, h.bs.Round, h.bs.Step))

	if voteMsg != nil {
		h.signAndPrevote(voteMsg)
	}
}

func (h *StateHandler) onNewRoundTimeout() {
	log.PBft.Info("[StateHandler-onNewRoundTimeout]")
	//Ignore this timeout when already entered other steps.
	if h.bs.Step != model2.RoundStepNewRound {
		return
	}

	reqAddresses := h.bs.getNewRoundReqList()
	h.broadcastReqRoundMsg(reqAddresses)

	if !h.bs.NewRound.EnoughAtRound(h.bs.Round) {
		h.ticker.ScheduleTimeout(components.TimeoutInfo{Duration: h.timeoutConfig.WaitNewRound, Height: h.bs.Height, Round: h.bs.Round, Step: model2.RoundStepNewRound})
	}
}

func (h *StateHandler) isCurProposer() bool {
	return h.curProposer().IsEqual(h.BftConfig.Signer.GetAddress())
}

func (h *StateHandler) onEnterPrecommit() {
	h.ticker.ScheduleTimeout(components.TimeoutInfo{Duration: h.timeoutConfig.WaitNewRound, Height: h.bs.Height, Round: h.bs.Round, Step: model2.RoundStepPreCommit})
	voteMsg := h.bs.makeVote()

	log.PBft.Debug(fmt.Sprintf("EnterPrecommit (H: %v, R: %v, S: %v)", h.bs.Height, h.bs.Round, h.bs.Step))
	if voteMsg != nil {
		h.signAndVote(voteMsg)
	}
}

func (h *StateHandler) onProposeTimeout() {
	log.PBft.Info("[StateHandler-onProposeTimeout]")
	preStep := h.bs.Step
	h.bs.OnProposeTimeout()
	curStep := h.bs.Step

	switch {
	case preStep == model2.RoundStepPropose && curStep == model2.RoundStepNewRound:
		log.PBft.Info("[StateHandler-onProposeTimeout]:onEnterNewRound")
		h.onEnterNewRound()
	default:
		log.PBft.Info("[StateHandler-onProposeTimeout]:block pool not empty", "pre", preStep, "new", curStep)
	}
}

func (h *StateHandler) onPreVoteTimeout() {
	log.PBft.Info("[StateHandler-onPreVoteTimeout]")
	preStep := h.bs.Step
	h.bs.OnPreVoteTimeout()
	curStep := h.bs.Step

	switch {
	case preStep == model2.RoundStepPreVote && curStep == model2.RoundStepNewRound:
		log.PBft.Info("[StateHandler-onPreVoteTimeout]:onEnterNewRound")
		h.onEnterNewRound()
	default:
		log.PBft.Info("[StateHandler-onPreVoteTimeout]: prevote time out", "pre", preStep, "new", curStep)
	}
}

func (h *StateHandler) onPreCommitTimeout() {
	log.PBft.Info("[StateHandler-onPreCommitTimeout]")
	preStep := h.bs.Step
	h.bs.OnPreCommitTimeout()
	curStep := h.bs.Step

	switch {
	case preStep == model2.RoundStepPreCommit && curStep == model2.RoundStepNewRound:
		log.PBft.Info("[StateHandler-onPreCommitTimeout]:onEnterNewRound")
		h.onEnterNewRound()
	default:
		log.PBft.Info("[StateHandler-onPreCommitTimeout]:precommit time out", "pre", preStep, "new", curStep)
	}
}

func (h *StateHandler) signAndPrevote(msg *model.VoteMsg) {
	// sign msg
	sign, err := h.BftConfig.Signer.SignHash(msg.Hash().Bytes())
	if err != nil {
		log.Warn("sign vote msg failed", "err", err)
		return
	}
	msg.Witness = &model.WitMsg{
		Address: h.BftConfig.Signer.GetAddress(),
		Sign:    sign,
	}

	h.Sender.BroadcastMsg(uint64(model2.TypeOfPreVoteMsg), msg)

	h.OnPreVote(msg)
}

func (h *StateHandler) signAndVote(msg *model.VoteMsg) {
	// sign msg
	sign, err := h.BftConfig.Signer.SignHash(msg.Hash().Bytes())
	if err != nil {
		log.Warn("sign vote msg failed", "err", err)
		return
	}
	msg.Witness = &model.WitMsg{
		Address: h.BftConfig.Signer.GetAddress(),
		Sign:    sign,
	}

	h.Sender.BroadcastMsg(uint64(model2.TypeOfVoteMsg), msg)

	h.OnVote(msg)
}

func (h *StateHandler) addTimeoutCount(label string) {
	g_metrics.Add(g_metrics.BftTimeoutCount, label, 1)
}

//Receive msgs
func (h *StateHandler) NewHeight(height uint64) {
	if h.IsRunning() {
		h.newHeightChan <- height
	}
}
func (h *StateHandler) NewRound(r *model2.NewRoundMsg) {
	if h.IsRunning() {
		h.newRoundChan <- r
	}
}
func (h *StateHandler) BlockPoolNotEmpty() {
	if h.IsRunning() {
		h.poolNotEmptyChan <- struct{}{}
	}
}
func (h *StateHandler) NewProposal(p *model2.Proposal) {
	if h.IsRunning() {
		h.newProposalChan <- p
	}
}
func (h *StateHandler) PreVote(pv *model.VoteMsg) {
	if h.IsRunning() {
		h.preVoteChan <- pv
	}
}
func (h *StateHandler) Vote(v *model.VoteMsg) {
	if h.IsRunning() {
		h.voteChan <- v
	}
}

func (h *StateHandler) GetProposalBlock(hash common.Hash) model.AbstractBlock {
	if !h.IsRunning() {
		return nil
	}
	result := make(chan model.AbstractBlock)
	h.getProposalBlockChan <- getProposalBlockMsg{
		hash:       hash,
		resultChan: result,
	}
	return <-result
}

func (h *StateHandler) onGetProposalBlock(msg getProposalBlockMsg) {
	msg.resultChan <- h.bs.ProposalBlock.GetBlockByHash(msg.hash)
}

//When peer request your round msg, return that
func (h *StateHandler) GetRoundMsg(height, round uint64) *model2.NewRoundMsg {
	if height != h.bs.Height {
		return nil
	}
	if round > h.bs.Round {
		return nil
	}
	msg := &model2.NewRoundMsg{
		Height: h.bs.Height,
		Round:  round,
	}
	sign, err := h.Signer.SignHash(msg.Hash().Bytes())
	if err != nil {
		log.Warn("sign new round msg failed", "err", err)
		return nil
	}
	msg.Witness = &model.WitMsg{
		Address: h.Signer.GetAddress(),
		Sign:    sign,
	}
	return msg
}

func (h *StateHandler) SetFetcher(fetcher components.Fetcher) {
	h.Fetcher = fetcher
}

func (h *StateHandler) RecordCurState() (uint64, uint64, model2.RoundStepType) {
	return h.bs.Height, h.bs.Round, h.bs.Step
}
