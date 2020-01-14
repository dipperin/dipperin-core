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
	"github.com/dipperin/dipperin-core/common/g-event"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/common/g-timer"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chaincommunication"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/event"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

const (
	notSynchronized = iota
	synchronized
)

const (
	notProcessEmptyBlock = iota
	processEmptyBlock
)

var (
	checkSynStatusDuration = 1 * time.Minute
	checkVerHaltDuration   = 5 * time.Minute
	//checkVerHaltDuration         = 30 * time.Second
	waitProposalResponseDuration = 1 * time.Minute
	waitVerifierVote             = 1 * time.Minute

	LogDuration = 30 * time.Second
)

type broadcastEmptyBlock func(block model.AbstractBlock)

//go:generate mockgen -destination=./peer_mock_test.go -package=verifiers_halt_check github.com/dipperin/dipperin-core/core/chaincommunication PmAbstractPeer

//go:generate mockgen -destination=./cs_protocol_mock_test.go -package=verifiers_halt_check github.com/dipperin/dipperin-core/core/verifiers-halt-check CsProtocolFunction
type CsProtocolFunction interface {
	GetVerifierBootNode() map[string]chaincommunication.PmAbstractPeer
	GetNextVerifierPeers() map[string]chaincommunication.PmAbstractPeer
	GetCurrentVerifierPeers() map[string]chaincommunication.PmAbstractPeer
}

type ProposalMsg struct {
	Round      uint64
	EmptyBlock model.Block
	VoteMsg    model.VoteMsg
}

type SystemHaltedCheck struct {
	//verifier boot node synchronization status
	SynStatus uint32

	//node type
	nodeType int

	//msg handlers
	handlers map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error

	//other verifier boot node current block number
	otherBootNodeHeight map[string]uint64
	//verifier current block number
	verifierHeight    map[string]uint64
	verifierMaxHeight uint64

	broadcaster broadcastEmptyBlock

	//handler halt status
	haltHandler *VBHaltHandler

	//handle block
	haltCheckStateHandle *StateHandler

	csProtocol CsProtocolFunction

	//start empty process flag
	startEmptyProcessFlag uint32

	proposalFail      chan bool
	aliveVerifierVote chan model.VoteMsg
	stopEmptyProcess  chan bool
	selectedProposal  chan ProposalMsg
	proposalInfoMsg   chan ProposalMsg
	heightInfo        chan heightResponseInfo
	quit              chan bool

	feed event.Feed
}

type getHeightResponse struct {
	Height uint64
}

type heightResponseInfo struct {
	Response getHeightResponse
	NodeName string
	NodeType uint64
	Address  common.Address
}

type HaltCheckConf struct {
	NodeType        int
	CsProtocol      CsProtocolFunction
	NeedChainReader NeedChainReaderFunction
	WalletSigner    NeedWalletSigner
	Broadcast       broadcastEmptyBlock
	EconomyModel    economy_model.EconomyModel
}

func MakeSystemHaltedCheck(conf *HaltCheckConf) *SystemHaltedCheck {
	systemHaltedCheck := &SystemHaltedCheck{
		handlers:             make(map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error),
		haltCheckStateHandle: MakeHaltCheckStateHandler(conf.NeedChainReader, conf.WalletSigner, conf.EconomyModel),
		csProtocol:           conf.CsProtocol,
		nodeType:             conf.NodeType,
	}

	chainConfig := chain_config.GetChainConfig()
	//verifier only need onCurrentBlockNumberRequest and onSendMinimalHashBlock
	//verifier don't need onSendMinimalHashBlock
	systemHaltedCheck.handlers[chaincommunication.CurrentBlockNumberRequest] = systemHaltedCheck.onCurrentBlockNumberRequest
	if conf.NodeType == chain_config.NodeTypeOfVerifierBoot {
		atomic.StoreUint32(&systemHaltedCheck.SynStatus, notSynchronized)
		systemHaltedCheck.otherBootNodeHeight = make(map[string]uint64)
		systemHaltedCheck.verifierHeight = make(map[string]uint64, 0)

		// channel buffer should be set as VerifierNumber+VerifierBootNodeNumber-1,
		// and abandon the message when receive a message simultaneously
		systemHaltedCheck.heightInfo = make(chan heightResponseInfo, chainConfig.VerifierNumber+chainConfig.VerifierBootNodeNumber-1)
		// systemHaltedCheck.emptyBlocks = make(map[common.Address]model.Block, 0)
		// channel buffer should be set as VerifierBootNodeNumber-1,
		// prevent all other verifier boot nodes from losing messages when broadcasting messages
		systemHaltedCheck.proposalInfoMsg = make(chan ProposalMsg, chainConfig.VerifierBootNodeNumber-1)
		systemHaltedCheck.quit = make(chan bool, 0)
		systemHaltedCheck.selectedProposal = make(chan ProposalMsg, 0)
		systemHaltedCheck.stopEmptyProcess = make(chan bool, 0)
		// channel buffer should be set as the number of Verifiers,
		// in case messages are lost when all other alive verifiers reply vote simultaneously
		systemHaltedCheck.aliveVerifierVote = make(chan model.VoteMsg, chainConfig.VerifierNumber)
		// be used to broadcast failure message when propose empty block
		systemHaltedCheck.proposalFail = make(chan bool, 0)

		atomic.StoreUint32(&systemHaltedCheck.startEmptyProcessFlag, notProcessEmptyBlock)
		systemHaltedCheck.broadcaster = conf.Broadcast

		systemHaltedCheck.handlers[chaincommunication.CurrentBlockNumberResponse] = systemHaltedCheck.onCurrentBlockNumberResponse
		systemHaltedCheck.handlers[chaincommunication.ProposeEmptyBlockMsg] = systemHaltedCheck.onProposeEmptyBlockMsg
		systemHaltedCheck.handlers[chaincommunication.SendMinimalHashBlockResponse] = systemHaltedCheck.onSendMinimalHashBlockResponse
	} else if conf.NodeType == chain_config.NodeTypeOfVerifier {
		systemHaltedCheck.handlers[chaincommunication.SendMinimalHashBlock] = systemHaltedCheck.onSendMinimalHashBlock
	}

	return systemHaltedCheck
}

func (systemHaltedCheck *SystemHaltedCheck) SetMsgSigner(walletSigner NeedWalletSigner) {
	systemHaltedCheck.haltCheckStateHandle.walletSigner = walletSigner
}

func (systemHaltedCheck *SystemHaltedCheck) MsgHandlers() map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {
	return systemHaltedCheck.handlers
}

func (systemHaltedCheck *SystemHaltedCheck) onCurrentBlockNumberRequest(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {

	blockNumber := systemHaltedCheck.haltCheckStateHandle.chainReader.CurrentBlock().Number()
	response := getHeightResponse{
		Height: blockNumber,
	}

	err := p.SendMsg(chaincommunication.CurrentBlockNumberResponse, &response)
	if err != nil {
		log.DLogger.Warn("send msg error", zap.Int("msgCode", chaincommunication.CurrentBlockNumberResponse), zap.Error(err))
	}

	return nil
}

func (systemHaltedCheck *SystemHaltedCheck) onCurrentBlockNumberResponse(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {

	var blockNumberResponse getHeightResponse
	err := msg.Decode(&blockNumberResponse)
	if err != nil {
		log.DLogger.Error("blockNumberResponse decode error", zap.Error(err))
		return err
	}

	if p.NodeType() != chain_config.NodeTypeOfVerifier && p.NodeType() != chain_config.NodeTypeOfVerifierBoot {
		return g_error.PeerTypeError
	}

	systemHaltedCheck.heightInfo <- heightResponseInfo{
		Response: blockNumberResponse,
		NodeName: p.NodeName(),
		NodeType: p.NodeType(),
		Address:  p.RemoteVerifierAddress(),
	}
	return nil
}

func (systemHaltedCheck *SystemHaltedCheck) checkPeerHeight() error {
	if systemHaltedCheck.nodeType != chain_config.NodeTypeOfVerifierBoot {
		return nil
	}

	chainConfig := chain_config.GetChainConfig()
	bootNodes := systemHaltedCheck.csProtocol.GetVerifierBootNode()
	log.DLogger.Info("the bootNodes number is:", zap.Int("number", len(bootNodes)))
	if len(bootNodes) != chainConfig.VerifierBootNodeNumber-1 {
		atomic.StoreUint32(&systemHaltedCheck.SynStatus, notSynchronized)
		return nil
	}

	currentVerifier := systemHaltedCheck.csProtocol.GetCurrentVerifierPeers()
	log.DLogger.Info("connect current verifier number is:", zap.Int("number", len(currentVerifier)))

	tickHandler := func() {
		for _, peer := range bootNodes {
			err := peer.SendMsg(chaincommunication.CurrentBlockNumberRequest, "")
			if err != nil {
				log.DLogger.Warn("send CurrentBlockNumberRequest error", zap.Error(err), zap.String("ToNodeName", peer.NodeName()))
			}
		}

		for _, peer := range currentVerifier {
			err := peer.SendMsg(chaincommunication.CurrentBlockNumberRequest, "")
			if err != nil {
				log.DLogger.Warn("send CurrentBlockNumberRequest error", zap.Error(err), zap.String("ToNodeName", peer.NodeName()))
			}
		}
	}
	ticker := g_timer.SetPeriodAndRun(tickHandler, checkSynStatusDuration)
	defer g_timer.StopWork(ticker)

	for {
		select {
		case heightResponse := <-systemHaltedCheck.heightInfo:
			if heightResponse.NodeType == chain_config.NodeTypeOfVerifier {
				systemHaltedCheck.verifierHeight[heightResponse.NodeName] = heightResponse.Response.Height
			} else if heightResponse.NodeType == chain_config.NodeTypeOfVerifierBoot {
				systemHaltedCheck.otherBootNodeHeight[heightResponse.NodeName] = heightResponse.Response.Height
			}
			log.DLogger.Info("the connected verBootNode peer block height is:", zap.Any("verBootNode", systemHaltedCheck.otherBootNodeHeight))
			log.DLogger.Info("the connected verifier peer block height is:", zap.Any("verifier", systemHaltedCheck.verifierHeight))
		case <-systemHaltedCheck.quit:
			return nil
		}
	}
}

func (systemHaltedCheck *SystemHaltedCheck) onProposeEmptyBlockMsg(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {
	// check msg from verifier boot node
	address := p.RemoteVerifierAddress()
	if !model.CheckAddressIsVerifierBootNode(address) {
		log.DLogger.Warn("the msg sender Address isn't verifier boot node")
		return nil
	}

	var proposal ProposalMsg
	err := msg.Decode(&proposal)
	if err != nil {
		log.DLogger.Warn("decode proposal msg error")
		return g_error.ProposalMsgDecodeError
	}
	log.DLogger.Info("receive an empty block proposal", zap.String("proposal", proposal.EmptyBlock.Hash().Hex()), zap.Uint64("height", proposal.EmptyBlock.Number()), zap.String("nodeName", p.NodeName()))

	// in case no coroutine reads emptyBlockInfo so that it is blocked
	select {
	case systemHaltedCheck.proposalInfoMsg <- proposal:
		log.DLogger.Info("received proposal~~~~", zap.String("proposal", proposal.EmptyBlock.Hash().Hex()), zap.String("fromNodeName", p.NodeName()), zap.Uint64("height", proposal.EmptyBlock.Number()))
	case <-time.After(100 * time.Millisecond):
	}

	return nil
}

// alive verifier send the vote of emptyBlock after receiving the minimalHashBlock from boot node verifier
func (systemHaltedCheck *SystemHaltedCheck) onSendMinimalHashBlock(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {
	var selectedProposal ProposalMsg

	err := msg.Decode(&selectedProposal)
	if err != nil {
		log.DLogger.Warn("decode minimal block msg error")
		return g_error.MinimalBlockDecodeError
	}

	log.DLogger.Info("received minimal hash block", zap.String("blockHash", selectedProposal.EmptyBlock.Hash().Hex()), zap.String("nodeName", p.NodeName()))
	// new AliveVerHaltHandler to valid and response the minimal hash block
	aliveVerHandler := NewAliveVerHaltHandler(systemHaltedCheck.haltCheckStateHandle.walletSigner.SignHash, systemHaltedCheck.haltCheckStateHandle.walletSigner.GetAddress())
	vote, err := aliveVerHandler.OnMinimalHashBlock(selectedProposal)
	if err != nil {
		log.DLogger.Warn("generateEmptyVoteMsg failed", zap.Error(err))
		return nil
	}

	log.DLogger.Info("send the vote for minimal hash empty Block", zap.Any("vote", vote))
	err = p.SendMsg(chaincommunication.SendMinimalHashBlockResponse, &vote)
	if err != nil {
		log.DLogger.Warn("SendMinimalHashBlockResponse error ", zap.String("nodeName", p.NodeName()))
	}

	return nil
}

// verifier bootNode receive the vote from the alive verifiers
func (systemHaltedCheck *SystemHaltedCheck) onSendMinimalHashBlockResponse(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {
	var vote model.VoteMsg
	err := msg.Decode(&vote)
	if err != nil {
		log.DLogger.Warn("decode aliveVerifierVote msg error")
		return g_error.VoteMsgDecodeError
	}

	log.DLogger.Info("receive the vote from alive verifier", zap.Any("vote", vote), zap.String("nodeName", p.NodeName()))
	// in case no coroutine reads systemHaltedCheck.vote so that it is blocked
	select {
	case systemHaltedCheck.aliveVerifierVote <- vote:
	case <-time.After(100 * time.Millisecond):
	}

	return nil
}

/*

Generate empty block and empty vote msg
Send the proposal to all v boot
Collect all the proposals and insert the block into the chain and broadcast it again.

*/
// verifier boot node propose an empty block when check the verifiers is halted
func (systemHaltedCheck *SystemHaltedCheck) proposeEmptyBlock() (err error) {

	log.DLogger.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~proposeEmptyBlock start")
	proposalConfig, err := systemHaltedCheck.haltCheckStateHandle.GenProposalConfig(model.VerBootNodeVoteMessage)
	if err != nil {
		log.DLogger.Info("proposeEmptyBlock generate proposal config error", zap.Error(err))
		return err
	}

	// New haltHandler to send and receive proposal msg
	// If the last propose height is the same as the current propose height,
	// use the same handler to prevent different propose block hashes for the same height empty block.
	// When all the boot nodes are out of sync, the empty block with the smallest hash cannot be picked out normally.
	var proposal ProposalMsg
	if systemHaltedCheck.haltHandler != nil && proposalConfig.CurBlock.Number() == systemHaltedCheck.haltHandler.pgConfig.CurBlock.Number() {
		proposal = *systemHaltedCheck.haltHandler.proposalMsg
	} else {
		systemHaltedCheck.haltHandler = NewHaltHandler(proposalConfig)
		proposal, err = systemHaltedCheck.haltHandler.ProposeEmptyBlock()
		if err != nil {
			log.DLogger.Info("proposeEmptyBlock error", zap.Error(err))
			return err
		}
	}

	log.DLogger.Info("the proposalMsg is:", zap.Any("proposalMsg", systemHaltedCheck.haltHandler.proposalMsg))
	log.DLogger.Info("propose empty block header is:", zap.String("header", systemHaltedCheck.haltHandler.proposalMsg.EmptyBlock.Header().Hash().Hex()))

	atomic.StoreUint32(&systemHaltedCheck.startEmptyProcessFlag, processEmptyBlock)

	errChan := make(chan error, 0)
	go func() {
		log.DLogger.Info("the connect ver boot node number is:", zap.Int("number", len(systemHaltedCheck.csProtocol.GetVerifierBootNode())))
		//propose empty block to other verifier boot node
		for _, node := range systemHaltedCheck.csProtocol.GetVerifierBootNode() {
			if !model.CheckAddressIsVerifierBootNode(node.RemoteVerifierAddress()) {
				errChan <- g_error.AddressIsNotVerifierBootNode
				return
			}

			log.DLogger.Info("send empty block propose is:", zap.String("propose", proposal.EmptyBlock.Hash().Hex()), zap.Uint64("height", proposal.EmptyBlock.Number()), zap.String("ToNodeName", node.NodeName()))
			err := node.SendMsg(chaincommunication.ProposeEmptyBlockMsg, &proposal)
			if err != nil {
				log.DLogger.Error("send propose empty block msg error", zap.Error(err), zap.String("nodeName", node.NodeName()))
				errChan <- err
				return
			}
		}
	}()

	waitOtherPropose := time.NewTimer(waitProposalResponseDuration)
	sub := systemHaltedCheck.feed.Subscribe(systemHaltedCheck.stopEmptyProcess)
	defer sub.Unsubscribe()
	for {
		select {
		case proposal := <-systemHaltedCheck.proposalInfoMsg:
			err := systemHaltedCheck.haltHandler.HandlerProposalMessages(proposal, systemHaltedCheck.selectedProposal)
			if err != nil {
				log.DLogger.Info("handler proposal messages err is:", zap.Error(err))
				// The received ms ms received error received by the processing received,
				// the error of handling proposal msg received,
				// should be returned until waitOtherPropose expires
			} else {
				log.DLogger.Info("handle proposal empty block end")
				return nil
			}

		case <-waitOtherPropose.C:
			//proposal empty block fail
			log.DLogger.Info("proposeEmptyBlock over time")
			systemHaltedCheck.proposalFail <- true
			return g_error.WaitEmptyBlockExpireError
		case <-systemHaltedCheck.stopEmptyProcess:
			log.DLogger.Info("stop proposeEmptyBlock")
			return nil
		case <-systemHaltedCheck.quit:
			return nil
		case readErr := <-errChan:
			log.DLogger.Info("proposeEmptyBlock err", zap.Error(readErr))
			//proposal empty block fail
			systemHaltedCheck.proposalFail <- true
			return readErr
		}
	}
}

//send proposal with minimal hash to verifiers
func (systemHaltedCheck *SystemHaltedCheck) sendMinimalHashBlock(proposal ProposalMsg) error {
	errChan := make(chan error, 0)
	//send selected block to verifiers
	go func() {
		log.DLogger.Info("sendMinimalHashBlock the currentVerifier peer number is:", zap.Int("number", len(systemHaltedCheck.csProtocol.GetCurrentVerifierPeers())))
		for _, node := range systemHaltedCheck.csProtocol.GetCurrentVerifierPeers() {
			//todo: check whether the peer in peerSet is currentVerifier. If the finder ensures the peerSet is normal ,then there is no need to check
			log.DLogger.Info("sendMinimalHashBlock the block hash is:", zap.Any("hash", proposal.EmptyBlock.Hash()), zap.String("toNodeName", node.NodeName()))
			err := node.SendMsg(chaincommunication.SendMinimalHashBlock, &proposal)
			if err != nil {
				log.DLogger.Error("send propose empty block msg error", zap.Error(err), zap.String("nodeName", node.NodeName()))
				//treat it as bad node when send error
				//errChan <- err
				return
			}
		}
	}()

	//chainConfig := chain_config.GetChainConfig()
	//collect verifier votes
	waitVerifierVote := time.NewTimer(waitVerifierVote)
	sub := systemHaltedCheck.feed.Subscribe(systemHaltedCheck.stopEmptyProcess)
	defer sub.Unsubscribe()
	for {
		select {
		case vote := <-systemHaltedCheck.aliveVerifierVote:
			// check aliveVerifierVote from verifier
			currentVerifiers := systemHaltedCheck.haltCheckStateHandle.chainReader.GetCurrVerifiers()
			systemHaltedCheck.haltHandler.HandlerAliveVerVotes(vote, currentVerifiers)

		case <-waitVerifierVote.C:
			//Broadcast empty block and the verifications
			err := systemHaltedCheck.handleFinalEmptyBlock(proposal, systemHaltedCheck.haltHandler.aliveVerVotes)
			if err != nil {
				//repeat proposal empty block if save error
				systemHaltedCheck.proposalFail <- true
			}
			return nil
		case <-systemHaltedCheck.stopEmptyProcess:
			log.DLogger.Info("stop sendMinimalHashBlock")
			return nil
		case <-systemHaltedCheck.quit:
			return nil
		case readErr := <-errChan:
			//repeat proposal empty block if sendMinimalHashBlock error
			systemHaltedCheck.proposalFail <- true
			return readErr

		}
	}
}

func (systemHaltedCheck *SystemHaltedCheck) handleFinalEmptyBlock(proposal ProposalMsg, votes map[common.Address]model.VoteMsg) error {
	//log.DLogger.Info("handleFinalEmptyBlock the votes is:","votes",votes)
	err := systemHaltedCheck.haltCheckStateHandle.SaveFinalEmptyBlock(proposal, votes)
	if err != nil {
		log.DLogger.Info("verifier boot node save final empty block failed", zap.Error(err))
		return err
	}

	//Broadcast empty block
	systemHaltedCheck.broadcaster(&proposal.EmptyBlock)
	return nil
}

/*

Check if the chain runs normally
1. if it expires, propose empty block
2. upon reception of proposed block, pick the empty block that should be inserted into the chain

*/
func (systemHaltedCheck *SystemHaltedCheck) checkVerClusterStatus() error {

	log.DLogger.Info("the systemHaltedCheck NodeType is:", zap.Int("NodeType", systemHaltedCheck.nodeType))
	if systemHaltedCheck.nodeType != chain_config.NodeTypeOfVerifierBoot {
		return nil
	}

	log.DLogger.Info("checkVerClusterStatus start~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	newBlockChan := make(chan model.Block, 0)
	//blockSub := systemHaltedCheck.haltCheckStateHandle.chainReader.SubscribeBlockEvent(newBlockChan)
	blockSub := g_event.Subscribe(g_event.NewBlockInsertEvent, newBlockChan)
	timer := time.NewTimer(checkVerHaltDuration)

	for {
		select {
		case <-timer.C:
			log.DLogger.Info("check system halt")
			go systemHaltedCheck.proposeEmptyBlock()
		case newBlock := <-newBlockChan:
			log.DLogger.Info("receive new block ", zap.Uint64("blockNumber", newBlock.Number()))
			g_metrics.Set(g_metrics.CurBlockNumberGauge, "", float64(newBlock.Number()))
			//receive new block from verifiers before haltDuration expire
			if atomic.LoadUint32(&systemHaltedCheck.startEmptyProcessFlag) == processEmptyBlock {
				// prevent the blockage caused due to the retirement of listening coroutine by timeout when writing
				systemHaltedCheck.feed.Send(true)
			}
			timer.Reset(checkVerHaltDuration)
		case minimalBlockProposal := <-systemHaltedCheck.selectedProposal:
			go systemHaltedCheck.sendMinimalHashBlock(minimalBlockProposal)
		case <-systemHaltedCheck.proposalFail:
			//repeat proposal when propose empty block fail
			log.DLogger.Info("repeat proposal empty block")
			go systemHaltedCheck.proposeEmptyBlock()
		case <-systemHaltedCheck.quit:
			blockSub.Unsubscribe()
			close(newBlockChan)
			return nil
		}
	}
}

func (systemHaltedCheck *SystemHaltedCheck) LogCurrentVerifier() {

	tickHandler := func() {
		currentVerifiers := systemHaltedCheck.haltCheckStateHandle.chainReader.GetCurrVerifiers()
		log.DLogger.Info("the current verifiers is:", zap.Any("ver", currentVerifiers))
	}
	tick := g_timer.SetPeriodAndRun(tickHandler, LogDuration)
	defer g_timer.StopWork(tick)

	<-systemHaltedCheck.quit
}

func (systemHaltedCheck *SystemHaltedCheck) LogConnectedCurrentVerifier() {

	tickHandler := func() {
		currentVerifierPeers := systemHaltedCheck.csProtocol.GetCurrentVerifierPeers()
		log.DLogger.Info("the connected current and nex verifier peer is:")
		for _, peer := range currentVerifierPeers {
			log.DLogger.Info("connected current verifier is:", zap.String("nodeName", peer.NodeName()))
		}

		nextVerifierPeers := systemHaltedCheck.csProtocol.GetNextVerifierPeers()
		for _, peer := range nextVerifierPeers {
			log.DLogger.Info("connected next verifier is:", zap.String("nodeName", peer.NodeName()))
		}
	}
	tick := g_timer.SetPeriodAndRun(tickHandler, LogDuration)
	defer g_timer.StopWork(tick)

	<-systemHaltedCheck.quit
}

func (systemHaltedCheck *SystemHaltedCheck) loop() {
	log.DLogger.Info("systemHaltedCheck loop start~~~~~~~~~~~~~~~~~")
	//go systemHaltedCheck.log.HaltConnectedCurrentVerifier()
	//go systemHaltedCheck.log.HaltCurrentVerifier()
	go systemHaltedCheck.checkPeerHeight()
	go systemHaltedCheck.checkVerClusterStatus()
	log.DLogger.Info("systemHaltedCheck loop end~~~~~~~~~~~~~~~~~")
}

func (systemHaltedCheck *SystemHaltedCheck) Start() error {
	log.DLogger.Info("SystemHaltedCheck start~~~~~~~~~~~~~~~", zap.Int("NodeType", systemHaltedCheck.nodeType))
	if systemHaltedCheck.nodeType != chain_config.NodeTypeOfVerifierBoot {
		return nil
	}
	systemHaltedCheck.loop()
	return nil
}

func (systemHaltedCheck *SystemHaltedCheck) Stop() {
	if systemHaltedCheck.nodeType != chain_config.NodeTypeOfVerifierBoot {
		return
	}
	close(systemHaltedCheck.quit)
	return
}
