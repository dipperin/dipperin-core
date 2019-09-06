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

package csbftnode

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/csbft/state-machine"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestCsBft_Start(t *testing.T) {
	node1 := NewTestNode()

	err := node1.Start()
	assert.NoError(t, err)

	err = node1.Start()
	assert.NoError(t, err)

	adderr := node1.AddPeer(nil)
	assert.NoError(t, adderr)

	block := &FakeBlock{
		Height: 1,
	}
	node1.OnNewWaitVerifyBlock(block, "")
	node1.Stop()
}

func TestCsBft_canStart(t *testing.T) {
	node1 := NewTestNode()
	can := node1.canStart()
	assert.Equal(t, can, true)
}

func TestCsBft_ChangePrimary(t *testing.T) {
	node1 := NewTestNode()
	node1.Stop()

	node1.ChangePrimary("")
	assert.Equal(t, true, node1.stateHandler.IsRunning())

	node1.stateHandler.ChainReader.SaveBlock(makeValidBlock(uint64(9), uint64(9)))
	node1.ChangePrimary("")
	assert.Equal(t, false, node1.stateHandler.IsRunning())

	node1.Start()
	assert.Equal(t, false, node1.stateHandler.IsRunning())
}

func TestCsBft_OnNewWaitVerifyBlock(t *testing.T) {
	node1 := NewTestNode()
	node1.Stop()

	block := &FakeBlock{uint64(3), common.HexToHash("0x123"), nil}
	node1.OnNewWaitVerifyBlock(block, "")
	assert.Equal(t, true, node1.blockPool.IsEmpty())

	node1.Start()
	node1.OnNewWaitVerifyBlock(block, "")
	assert.Equal(t, true, node1.blockPool.IsEmpty())

	block1 := &FakeBlock{uint64(1), common.HexToHash("0x123"), nil}
	node1.OnNewWaitVerifyBlock(block1, "")
	assert.Equal(t, false, node1.blockPool.IsEmpty())
}

func TestCsBft_isNextVerifier(t *testing.T) {
	node1 := NewTestNode()
	assert.Equal(t, true, node1.isNextVerifier())
}

func TestCsBft_OnNewP2PMsg(t *testing.T) {
	node1 := NewTestNode()
	msg := p2p.Msg{}

	err := node1.OnNewMsg(msg)
	assert.Equal(t, err, nil)
	address := common.HexToAddress("0x54bbe8ffddc")
	node1.OnNewP2PMsg(msg, &tPeer{0, "", "", address})
}

func TestCsBft_OnNewP2PMsg2(t *testing.T) {
	node1 := NewTestNode()
	node1.Start()
	address := common.HexToAddress("0x54bbe8ffddc")

	// NewRoundMsg
	msg := &model2.NewRoundMsg{
		Height: 1,
		Round:  1,
	}
	size, r, err := rlp.EncodeToReader(msg)
	assert.NoError(t, err)
	p2pMsg1 := p2p.Msg{
		Code:    uint64(model2.TypeOfNewRoundMsg),
		Size:    uint32(size),
		Payload: r,
	}

	// ProposalMsg
	p2pMsg2 := p2p.Msg{
		Code:    uint64(model2.TypeOfProposalMsg),
		Size:    uint32(size),
		Payload: r,
	}

	// PrevoteMsg
	p2pMsg3 := p2p.Msg{
		Code:    uint64(model2.TypeOfPreVoteMsg),
		Size:    uint32(size),
		Payload: r,
	}

	// VoteMsg
	p2pMsg4 := p2p.Msg{
		Code:    uint64(model2.TypeOfVoteMsg),
		Size:    uint32(size),
		Payload: r,
	}

	p2pMsg5 := p2p.Msg{
		Code:    uint64(model2.TypeOfFetchBlockReqMsg),
		Size:    uint32(size),
		Payload: r,
	}
	node1.OnNewP2PMsg(p2pMsg5, &tPeer{0, "", "", address})

	// FetchBlockResp
	p2pMsg6 := p2p.Msg{
		Code:    uint64(model2.TypeOfFetchBlockRespMsg),
		Size:    uint32(size),
		Payload: r,
	}

	p2pMsg7 := p2p.Msg{
		Code:    uint64(model2.TypeOfSyncBlockMsg),
		Size:    uint32(size),
		Payload: r,
	}

	p2pMsg8 := p2p.Msg{
		Code:    uint64(model2.TypeOfReqNewRoundMsg),
		Size:    uint32(size),
		Payload: r,
	}

	p2pMsg9 := p2p.Msg{
		Code:    uint64(model2.TypeOfReqNewRoundMsg + 1),
		Size:    uint32(size),
		Payload: r,
	}

	node1.OnNewP2PMsg(p2pMsg1, &tPeer{0, "", "", address})
	node1.OnNewP2PMsg(p2pMsg2, &tPeer{0, "", "", address})
	node1.OnNewP2PMsg(p2pMsg3, &tPeer{0, "", "", address})
	node1.OnNewP2PMsg(p2pMsg4, &tPeer{0, "", "", address})
	node1.OnNewP2PMsg(p2pMsg5, &tPeer{0, "", "", address})
	node1.OnNewP2PMsg(p2pMsg6, &tPeer{0, "", "", address})
	node1.OnNewP2PMsg(p2pMsg7, &tPeer{0, "", "", address})
	node1.OnNewP2PMsg(p2pMsg8, &tPeer{0, "", "", address})

	assert.Panics(t, func() {
		node1.OnNewP2PMsg(p2pMsg9, &tPeer{0, "", "", address})
	})
}

// Takes a long time
func TestCsBft_onSyncBlockMsg(t *testing.T) {
	node1 := NewTestNode()
	node1.Start()

	address := common.HexToAddress("0x54bbe8ffddc")
	hash := common.HexToHash("0x54bbe8ffddc")
	node1.onSyncBlockMsg(address, common.Hash{})
	node1.onSyncBlockMsg(common.Address{}, common.Hash{})
	node1.onSyncBlockMsg(common.Address{}, hash)
	node1.onSyncBlockMsg(address, hash)
	node1.Stop()
}

func TestCsBft_onSyncBlockMsg2(t *testing.T) {
	node1 := NewTestNode()

	address := common.HexToAddress("0x54bbe8ffddc42")
	node1.onSyncBlockMsg(address, common.Hash{})

}

// Test fetch block req msg
func TestCsBft_OnNewP2PMsg3(t *testing.T) {
	node1 := NewTestNode()
	node1.stateHandler.Start()
	node1.blockPool.Start()
	node1.blockPool.AddBlock(&FakeBlock{1, common.HexToHash("0x22"), nil})
	address := common.HexToAddress("0x54bbe8ffddc42")
	node1.onSyncBlockMsg(address, common.HexToHash("0x22"))

	msg := model2.FetchBlockReqDecodeMsg{uint64(1), common.HexToHash("22")}
	size, r, _ := rlp.EncodeToReader(msg)
	node1.OnNewP2PMsg(p2p.Msg{Code: uint64(model2.TypeOfFetchBlockReqMsg), Size: uint32(size), Payload: r}, &tPeer{0, "", "", common.HexToAddress("0x54bbe8ffddc")})

}

// Test msg TypeOfReqNewRoundMsg
func TestCsBft_OnNewP2PMsg4(t *testing.T) {
	node1 := NewTestNode()
	node1.stateHandler.Start()

	msg := model2.ReqRoundMsg{uint64(1), uint64(1)}
	size, r, _ := rlp.EncodeToReader(msg)
	node1.OnNewP2PMsg(p2p.Msg{Code: uint64(model2.TypeOfReqNewRoundMsg), Size: uint32(size), Payload: r}, &tPeer{0, "", "", common.HexToAddress("0x54bbe8ffddc")})
}

// Test NewRoundMsg
func TestCsBft_OnNewP2PMsg5(t *testing.T) {
	node1 := NewTestNode()
	node1.Start()
	address := common.HexToAddress("0x54bbe8ffddc")

	// NewRoundMsg
	msg := &model2.NewRoundMsg{
		Height: 1,
		Round:  1,
	}
	size, r, err := rlp.EncodeToReader(msg)
	assert.NoError(t, err)
	newRound := p2p.Msg{
		Code:    uint64(model2.TypeOfNewRoundMsg),
		Size:    uint32(size),
		Payload: r,
	}
	node1.OnNewP2PMsg(newRound, &tPeer{0, "", "", address})
}

// Test Proposal
func TestCsBft_OnNewP2PMsg6(t *testing.T) {
	node1 := NewTestNode()
	node1.Start()
	address := common.HexToAddress("0x54bbe8ffddc")

	p := MakeNewProposal(1, 1, &FakeBlock{}, 1)
	size, r, _ := rlp.EncodeToReader(p)
	proposalMsg := p2p.Msg{
		Code:    uint64(model2.TypeOfProposalMsg),
		Size:    uint32(size),
		Payload: r,
	}
	node1.OnNewP2PMsg(proposalMsg, &tPeer{0, "", "", address})
}

// Test PrevoteMsg
func TestCsBft_OnNewP2PMsg7(t *testing.T) {
	node1 := NewTestNode()
	node1.Start()
	address := common.HexToAddress("0x54bbe8ffddc")

	preVote := MakeNewProVote(1, 1, &FakeBlock{}, 1)
	size, r, _ := rlp.EncodeToReader(preVote)
	p2pMsg3 := p2p.Msg{
		Code:    uint64(model2.TypeOfPreVoteMsg),
		Size:    uint32(size),
		Payload: r,
	}
	node1.OnNewP2PMsg(p2pMsg3, &tPeer{0, "", "", address})
}

// Test VoteMsg
func TestCsBft_OnNewP2PMsg8(t *testing.T) {
	node1 := NewTestNode()
	node1.Start()
	address := common.HexToAddress("0x54bbe8ffddc")

	vote := MakeNewVote(1, 1, &FakeBlock{}, 1)
	size, r, _ := rlp.EncodeToReader(vote)
	p2pMsg4 := p2p.Msg{
		Code:    uint64(model2.TypeOfVoteMsg),
		Size:    uint32(size),
		Payload: r,
	}
	node1.OnNewP2PMsg(p2pMsg4, &tPeer{0, "", "", address})
}

// Test FetchBlockResponse
func TestCsBft_OnNewP2PMsg9(t *testing.T) {
	node1 := NewTestNode()
	node1.Start()
	address := common.HexToAddress("0x54bbe8ffddc")
	var block model.Block

	fetch := model2.FetchBlockRespDecodeMsg{1, &block}
	size, r, _ := rlp.EncodeToReader(fetch)
	p2pMsg4 := p2p.Msg{
		Code:    uint64(model2.TypeOfFetchBlockRespMsg),
		Size:    uint32(size),
		Payload: r,
	}
	node1.OnNewP2PMsg(p2pMsg4, &tPeer{0, "", "", address})
}

func TestCsBft_Remind(t *testing.T) {
	node1 := NewTestNode()
	curB := node1.ChainReader.CurrentBlock()
	node1.ChainReader.IsChangePoint(curB, false)
	err := node1.Start()
	assert.Equal(t, err, nil)
}

func makeValidBlock(height uint64, round uint64) (fakeblock *FakeBlock, commits []model.AbstractVerification) {
	fakeblock = &FakeBlock{height, common.HexToHash(strconv.Itoa(int(height))), nil}
	commits = append(commits, MakeNewVote(height, round, fakeblock, 0))
	commits = append(commits, MakeNewVote(height, round, fakeblock, 1))
	commits = append(commits, MakeNewVote(height, round, fakeblock, 2))
	return
}

func MakeNewVote(height uint64, round uint64, block model.AbstractBlock, i int) *model.VoteMsg {
	sks, _ := CreateKey()
	signer := newFackSigner(sks[i])
	voteMsg := model.VoteMsg{
		Height:    height,
		Round:     round,
		BlockID:   block.Hash(),
		VoteType:  model.VoteMessage,
		Timestamp: time.Now(),
	}
	voteMsg.Witness = Sign(voteMsg, signer)
	return &voteMsg
}
func Sign(b H, signer *fakeSigner) *model.WitMsg {
	sign, err := signer.SignHash(b.Hash().Bytes())
	if err != nil {
		return nil
	}
	witness := &model.WitMsg{
		Address: signer.GetAddress(),
		Sign:    sign,
	}
	return witness
}

type H interface {
	Hash() common.Hash
}

func NewTestNode() *CsBft {
	fc := NewFakeFullChain()
	sks, _ := CreateKey()
	fs := newFackSigner(sks[1])
	fcn := &FC{}
	fetcher := components.NewFetcher(fcn)
	config := &state_machine.BftConfig{fc, fetcher, fs, &FackMsgSender{}, &FakeValidtor{}}
	csbft := NewCsBft(config)
	fc.SetNewHeightNotifier(csbft.OnEnterNewHeight)
	csbft.SetFetcher(fetcher)
	return csbft
}

type FC struct{}

func (fc *FC) SendFetchBlockMsg(msgCode uint64, from common.Address, msg *model2.FetchBlockReqDecodeMsg) error {
	return nil
}

func MakeNewProposal(height uint64, round uint64, block model.AbstractBlock, key int) *model2.Proposal {
	sks, _ := CreateKey()
	signer := newFackSigner(sks[key])
	proposal := model2.Proposal{
		Height:    height,
		Round:     round,
		BlockID:   block.Hash(),
		Timestamp: time.Now(),
	}
	proposal.Witness = Sign(proposal, signer)
	return &proposal
}

func MakeNewProVote(height uint64, round uint64, block model.AbstractBlock, i int) *model.VoteMsg {
	sks, _ := CreateKey()
	signer := newFackSigner(sks[i])
	voteMsg := model.VoteMsg{
		Height:    height,
		Round:     round,
		BlockID:   block.Hash(),
		VoteType:  model.PreVoteMessage,
		Timestamp: time.Now(),
	}
	voteMsg.Witness = Sign(voteMsg, signer)
	return &voteMsg
}
