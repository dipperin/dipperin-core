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
	"errors"
	"github.com/dipperin/dipperin-core/common/g-event"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/cachedb"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
)

var testVerBootAccounts []tests.Account

func init() {
	var err error
	testVerBootAccounts, err = tests.ChangeVerBootNodeAddress()
	if err != nil {
		panic("change verifier boot node address error for test")
	}
}

func TestMakeSystemHaltedCheck(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mChainR := NewMockNeedChainReaderFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfVerifierBoot,
		NeedChainReader: mChainR,
	})
	assert.NotEmpty(t, haltedCheck.MsgHandlers())

	mp1 := NewMockPmAbstractPeer(c)
	mChainR.EXPECT().CurrentBlock().Return(model.NewBlock(&model.Header{Number: 1}, nil, nil))
	mp1.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("failed"))
	assert.NoError(t, haltedCheck.onCurrentBlockNumberRequest(p2p.Msg{}, mp1))

	mp2 := NewMockPmAbstractPeer(c)
	assert.Error(t, haltedCheck.onCurrentBlockNumberResponse(getMsg(1, struct{}{}), mp2))
	mp2.EXPECT().NodeType().Return(uint64(chain_config.NodeTypeOfNormal)).AnyTimes()
	assert.Error(t, haltedCheck.onCurrentBlockNumberResponse(getMsg(1, getHeightResponse{Height: 1}), mp2))

	mp3 := NewMockPmAbstractPeer(c)
	mp3.EXPECT().NodeType().Return(uint64(chain_config.NodeTypeOfVerifierBoot)).AnyTimes()
	mp3.EXPECT().NodeName().Return("test").AnyTimes()
	mp3.EXPECT().RemoteVerifierAddress().Return(common.Address{0x12}).AnyTimes()
	assert.NoError(t, haltedCheck.onCurrentBlockNumberResponse(getMsg(1, getHeightResponse{Height: 1}), mp3))

	haltedCheck = MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType: chain_config.NodeTypeOfVerifier,
	})
	assert.NotEmpty(t, haltedCheck.MsgHandlers())
	haltedCheck.Stop()
}

func TestSystemHaltedCheck_checkPeerHeight(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mChainR := NewMockNeedChainReaderFunction(c)
	mProtocol := NewMockCsProtocolFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfNormal,
		NeedChainReader: mChainR,
		CsProtocol:      mProtocol,
	})
	assert.NoError(t, haltedCheck.checkPeerHeight())

	haltedCheck = MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfVerifierBoot,
		NeedChainReader: mChainR,
		CsProtocol:      mProtocol,
	})
	mProtocol.EXPECT().GetVerifierBootNode().Return(map[string]chain_communication.PmAbstractPeer{})
	assert.NoError(t, haltedCheck.checkPeerHeight())
	mProtocol.EXPECT().GetVerifierBootNode().Return(map[string]chain_communication.PmAbstractPeer{
		"b1": newMockPeer(c),
		"b2": newMockPeer(c),
		"b3": newMockPeer(c),
	})
	mProtocol.EXPECT().GetCurrentVerifierPeers().Return(map[string]chain_communication.PmAbstractPeer{
		"v1": newMockPeer(c),
		"v2": newMockPeer(c),
		"v3": newMockPeer(c),
	})
	checkSynStatusDuration = time.Millisecond
	go haltedCheck.checkPeerHeight()
	time.Sleep(2 * time.Millisecond)

	haltedCheck.heightInfo <- heightResponseInfo{NodeType: chain_config.NodeTypeOfVerifierBoot}
	haltedCheck.heightInfo <- heightResponseInfo{NodeType: chain_config.NodeTypeOfVerifier}
	time.Sleep(2 * time.Millisecond)
	haltedCheck.quit <- true
	time.Sleep(time.Millisecond)
	haltedCheck.Stop()
}

func TestSystemHaltedCheck_onProposeEmptyBlockMsg(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType: chain_config.NodeTypeOfVerifierBoot,
		//NeedChainReader: mChainR,
		//CsProtocol: mProtocol,
	})
	mp0 := NewMockPmAbstractPeer(c)
	mp0.EXPECT().RemoteVerifierAddress().Return(common.Address{})
	assert.NoError(t, haltedCheck.onProposeEmptyBlockMsg(p2p.Msg{}, mp0))

	mp1 := newMockPeer(c)
	mp1.EXPECT().RemoteVerifierAddress().Return(chain_config.VerBootNodeAddress[0]).AnyTimes()
	assert.Error(t, haltedCheck.onProposeEmptyBlockMsg(getMsg(1, struct{}{}), mp1))

	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	v, err := model.NewVoteMsgWithSign(1, 0, eB.Hash(), model.VerBootNodeVoteMessage, testVerBootAccounts[0].SignHash, testVerBootAccounts[0].Address())
	assert.NoError(t, err)
	assert.NoError(t, haltedCheck.onProposeEmptyBlockMsg(getMsg(1, &ProposalMsg{
		EmptyBlock: *eB,
		VoteMsg:    *v,
	}), mp1))
}

func TestSystemHaltedCheck_onSendMinimalHashBlock(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mp1 := newMockPeer(c)
	mw := NewMockNeedWalletSigner(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:     chain_config.NodeTypeOfVerifierBoot,
		WalletSigner: mw,
		//NeedChainReader: mChainR,
		//CsProtocol: mProtocol,
	})
	assert.Error(t, haltedCheck.onSendMinimalHashBlock(getMsg(1, struct{}{}), mp1))

	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	v, err := model.NewVoteMsgWithSign(1, 0, eB.Hash(), model.VerBootNodeVoteMessage, testVerBootAccounts[0].SignHash, testVerBootAccounts[0].Address())
	assert.NoError(t, err)
	mw.EXPECT().GetAddress().Return(testVerBootAccounts[0].Address())
	mw.EXPECT().SignHash(gomock.Any()).Return([]byte{0x1}, nil)
	assert.NoError(t, haltedCheck.onSendMinimalHashBlock(getMsg(1, &ProposalMsg{
		EmptyBlock: *eB,
		VoteMsg:    *v,
	}), mp1))
}

func TestSystemHaltedCheck_onSendMinimalHashBlockResponse(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mp1 := newMockPeer(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType: chain_config.NodeTypeOfVerifierBoot,
		//WalletSigner: mw,
		//NeedChainReader: mChainR,
		//CsProtocol: mProtocol,
	})
	assert.Error(t, haltedCheck.onSendMinimalHashBlockResponse(getMsg(1, struct{}{}), mp1))
	assert.NoError(t, haltedCheck.onSendMinimalHashBlockResponse(getMsg(1, model.NewVoteMsg(1, 1, common.Hash{0x1}, model.VerBootNodeVoteMessage)), mp1))
}

func TestSystemHaltedCheck_proposeEmptyBlock(t *testing.T) {

	c := gomock.NewController(t)
	defer c.Finish()
	//mChainR := NewMockNeedChainReaderFunction(c)
	cfg := chain_config.GetChainConfig()
	cs := chain_state.NewChainState(&chain_state.ChainStateConfig{
		ChainConfig:   cfg,
		DataDir:       "",
		WriterFactory: chain_writer.NewChainWriterFactory(),
	})
	cs_chain.GenesisSetUp = true
	mChainR := cs_chain.NewCsChainService(&cs_chain.CsChainServiceConfig{
		CacheDB: cachedb.NewCacheDB(ethdb.NewMemDatabase()),
	}, cs)
	mw := NewMockNeedWalletSigner(c)
	mProtocol := NewMockCsProtocolFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfVerifierBoot,
		WalletSigner:    mw,
		NeedChainReader: mChainR,
		EconomyModel:    mChainR.GetEconomyModel(),
		CsProtocol:      mProtocol,
	})
	//mChainR.EXPECT().CurrentBlock().Return(model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)).AnyTimes()
	mw.EXPECT().GetAddress().Return(testVerBootAccounts[0].Address()).AnyTimes()
	mw.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return([32]byte{}, []byte{}, errors.New("failed"))
	assert.Error(t, haltedCheck.proposeEmptyBlock())

	mw.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return([32]byte{0x12}, []byte{0x11}, nil)
	//mChainR.EXPECT().GetSeenCommit(gomock.Any()).Return([]model.AbstractVerification{&model.VoteMsg{}})
	mw.EXPECT().PublicKey().Return(&testVerBootAccounts[0].Pk.PublicKey).AnyTimes()
	mw.EXPECT().SignHash(gomock.Any()).Return([]byte{0x12}, nil)
	//mChainR.EXPECT().BlockProcessor(gomock.Any()).Return()
	mProtocol.EXPECT().GetVerifierBootNode().Return(map[string]chain_communication.PmAbstractPeer{
		"1": newSuccessSendMockPeer(c),
	}).AnyTimes()
	go haltedCheck.proposeEmptyBlock()
	time.Sleep(50 * time.Millisecond)

	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	v, err := model.NewVoteMsgWithSign(1, 0, eB.Hash(), model.VerBootNodeVoteMessage, testVerBootAccounts[0].SignHash, testVerBootAccounts[0].Address())
	assert.NoError(t, err)
	haltedCheck.proposalInfoMsg <- ProposalMsg{EmptyBlock: *eB, VoteMsg: *v}
	time.Sleep(50 * time.Millisecond)
}

func TestSystemHaltedCheck_proposeEmptyBlock1(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	//mChainR := NewMockNeedChainReaderFunction(c)
	cfg := chain_config.GetChainConfig()
	cs := chain_state.NewChainState(&chain_state.ChainStateConfig{
		ChainConfig:   cfg,
		DataDir:       "",
		WriterFactory: chain_writer.NewChainWriterFactory(),
	})
	cs_chain.GenesisSetUp = true
	mChainR := cs_chain.NewCsChainService(&cs_chain.CsChainServiceConfig{
		CacheDB: cachedb.NewCacheDB(ethdb.NewMemDatabase()),
	}, cs)
	mw := NewMockNeedWalletSigner(c)
	mProtocol := NewMockCsProtocolFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfVerifierBoot,
		WalletSigner:    mw,
		NeedChainReader: mChainR,
		EconomyModel:    mChainR.GetEconomyModel(),
		CsProtocol:      mProtocol,
	})
	//mChainR.EXPECT().CurrentBlock().Return(model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)).AnyTimes()
	mw.EXPECT().GetAddress().Return(testVerBootAccounts[0].Address()).AnyTimes()
	mw.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return([32]byte{}, []byte{}, errors.New("failed"))
	assert.Error(t, haltedCheck.proposeEmptyBlock())

	mw.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return([32]byte{0x12}, []byte{0x11}, nil)
	//mChainR.EXPECT().GetSeenCommit(gomock.Any()).Return([]model.AbstractVerification{&model.VoteMsg{}})
	mw.EXPECT().PublicKey().Return(&testVerBootAccounts[0].Pk.PublicKey).AnyTimes()
	mw.EXPECT().SignHash(gomock.Any()).Return([]byte{0x12}, nil)
	//mChainR.EXPECT().BlockProcessor(gomock.Any()).Return()
	mProtocol.EXPECT().GetVerifierBootNode().Return(map[string]chain_communication.PmAbstractPeer{
		"1": newMockPeer(c),
	}).AnyTimes()
	go haltedCheck.proposeEmptyBlock()
	time.Sleep(50 * time.Millisecond)

	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	v, err := model.NewVoteMsgWithSign(1, 0, eB.Hash(), model.VerBootNodeVoteMessage, testVerBootAccounts[0].SignHash, testVerBootAccounts[0].Address())
	assert.NoError(t, err)
	haltedCheck.proposalInfoMsg <- ProposalMsg{EmptyBlock: *eB, VoteMsg: *v}
	time.Sleep(50 * time.Millisecond)
}

func TestSystemHaltedCheck_sendMinimalHashBlock(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mProtocol := NewMockCsProtocolFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:   chain_config.NodeTypeOfVerifierBoot,
		CsProtocol: mProtocol,
		//WalletSigner: mw,
		//NeedChainReader: mChainR,
		//EconomyModel: mChainR.GetEconomyModel(),
	})
	mProtocol.EXPECT().GetCurrentVerifierPeers().Return(map[string]chain_communication.PmAbstractPeer{
		"1": newMockPeer(c),
	}).AnyTimes()
	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	go haltedCheck.sendMinimalHashBlock(ProposalMsg{EmptyBlock: *eB})
	time.Sleep(time.Millisecond)
}

func TestSystemHaltedCheck_sendMinimalHashBlock1(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mProtocol := NewMockCsProtocolFunction(c)
	mChainR := NewMockNeedChainReaderFunction(c)
	mw := NewMockNeedWalletSigner(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfVerifierBoot,
		CsProtocol:      mProtocol,
		NeedChainReader: mChainR,
		WalletSigner:    mw,
		Broadcast:       func(block model.AbstractBlock) {},
	})
	mw.EXPECT().GetAddress().Return(common.Address{0x11})
	mw.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return([32]byte{0x11}, []byte{0x11}, nil)
	mChainR.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mChainR.EXPECT().CurrentBlock().Return(model.NewBlock(&model.Header{Number: 0}, nil, nil))
	mChainR.EXPECT().GetSeenCommit(gomock.Any()).Return([]model.AbstractVerification{&model.VoteMsg{}})
	mw.EXPECT().PublicKey().Return(&testVerBootAccounts[0].Pk.PublicKey)
	hhCfg, err := haltedCheck.haltCheckStateHandle.GenProposalConfig(model.VerBootNodeVoteMessage)
	assert.NoError(t, err)
	haltedCheck.haltHandler = NewHaltHandler(hhCfg)
	mProtocol.EXPECT().GetCurrentVerifierPeers().Return(map[string]chain_communication.PmAbstractPeer{
		"1": newSuccessSendMockPeer(c),
	}).AnyTimes()
	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	go haltedCheck.sendMinimalHashBlock(ProposalMsg{EmptyBlock: *eB})
	waitVerifierVote = time.Millisecond
	time.Sleep(10 * time.Millisecond)
}

func TestSystemHaltedCheck_handleFinalEmptyBlock(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mChainR := NewMockNeedChainReaderFunction(c)
	//mProtocol := NewMockCsProtocolFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfVerifierBoot,
		NeedChainReader: mChainR,
		Broadcast:       func(block model.AbstractBlock) {},
		//CsProtocol: mProtocol,
		//WalletSigner: mw,
		//EconomyModel: mChainR.GetEconomyModel(),
	})
	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	mChainR.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).Return(errors.New("failed"))
	assert.Error(t, haltedCheck.handleFinalEmptyBlock(ProposalMsg{EmptyBlock: *eB}, map[common.Address]model.VoteMsg{}))

	mChainR.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).Return(nil)
	assert.NoError(t, haltedCheck.handleFinalEmptyBlock(ProposalMsg{EmptyBlock: *eB}, map[common.Address]model.VoteMsg{}))
}

func TestSystemHaltedCheck_checkVerClusterStatus(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	g_event.Add(g_event.NewBlockInsertEvent)
	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType: chain_config.NodeTypeOfVerifierBoot,
	})

	go haltedCheck.checkVerClusterStatus()
	time.Sleep(time.Millisecond)
	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	g_event.Send(g_event.NewBlockInsertEvent, *eB)
	time.Sleep(time.Millisecond)
	haltedCheck.quit <- true
	time.Sleep(time.Millisecond)
}

func TestSystemHaltedCheck_logCurrentVerifier(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mChainR := NewMockNeedChainReaderFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:        chain_config.NodeTypeOfVerifierBoot,
		NeedChainReader: mChainR,
	})
	LogDuration = time.Millisecond
	mChainR.EXPECT().GetCurrVerifiers().Return([]common.Address{{0x12}}).AnyTimes()
	go haltedCheck.LogCurrentVerifier()
	time.Sleep(2 * time.Millisecond)
}

func TestSystemHaltedCheck_logConnectedCurrentVerifier(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mProtocol := NewMockCsProtocolFunction(c)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:   chain_config.NodeTypeOfVerifierBoot,
		CsProtocol: mProtocol,
	})

	LogDuration = time.Millisecond
	mProtocol.EXPECT().GetCurrentVerifierPeers().Return(map[string]chain_communication.PmAbstractPeer{
		"1": newMockPeer(c),
	}).AnyTimes()
	mProtocol.EXPECT().GetNextVerifierPeers().Return(map[string]chain_communication.PmAbstractPeer{
		"1": newMockPeer(c),
	}).AnyTimes()

	go haltedCheck.LogConnectedCurrentVerifier()
	time.Sleep(2 * time.Millisecond)
}

func TestSystemHaltedCheck_loop(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()
	mProtocol := NewMockCsProtocolFunction(c)

	g_event.Add(g_event.NewBlockInsertEvent)

	haltedCheck := MakeSystemHaltedCheck(&HaltCheckConf{
		NodeType:   chain_config.NodeTypeOfVerifierBoot,
		CsProtocol: mProtocol,
	})

	mProtocol.EXPECT().GetVerifierBootNode().Return(map[string]chain_communication.PmAbstractPeer{}).AnyTimes()

	assert.NoError(t, haltedCheck.Start())
	haltedCheck.nodeType = chain_config.NodeTypeOfNormal
	assert.NoError(t, haltedCheck.Start())
	haltedCheck.loop()
}

func TestCheckProposalValid(t *testing.T) {
	eB := model.NewBlock(&model.Header{Number: 1, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	v, err := model.NewVoteMsgWithSign(1, 1, eB.Hash(), model.VerBootNodeVoteMessage, testVerBootAccounts[0].SignHash, testVerBootAccounts[0].Address())
	assert.NoError(t, err)
	assert.NoError(t, checkProposalValid(ProposalMsg{
		EmptyBlock: *eB,
		VoteMsg:    *v,
	}))

	fakeB := model.NewBlock(&model.Header{Number: 2, Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}, nil, nil)
	assert.Error(t, checkProposalValid(ProposalMsg{
		EmptyBlock: *fakeB,
		VoteMsg:    *v,
	}))

	fakeV := model.NewVoteMsg(1, 1, eB.Hash(), model.VerBootNodeVoteMessage)
	assert.Error(t, checkProposalValid(ProposalMsg{
		EmptyBlock: *eB,
		VoteMsg:    *fakeV,
	}))
}

func newSuccessSendMockPeer(c *gomock.Controller) *MockPmAbstractPeer {
	mp := NewMockPmAbstractPeer(c)
	mp.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mp.EXPECT().NodeName().Return("test").AnyTimes()
	mp.EXPECT().RemoteVerifierAddress().Return(chain_config.VerifierBootNodeAddress[0]).AnyTimes()
	return mp
}

func newMockPeer(c *gomock.Controller) *MockPmAbstractPeer {
	mp := NewMockPmAbstractPeer(c)
	mp.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("failed")).AnyTimes()
	mp.EXPECT().NodeName().Return("test").AnyTimes()
	mp.EXPECT().RemoteVerifierAddress().Return(chain_config.VerifierBootNodeAddress[0]).AnyTimes()
	return mp
}

func getMsg(mCode uint64, data interface{}) p2p.Msg {
	s, r, err := rlp.EncodeToReader(data)
	if err != nil {
		panic(err)
	}

	return p2p.Msg{
		Code:       mCode,
		Size:       uint32(s),
		Payload:    r,
		ReceivedAt: time.Now(),
	}
}
