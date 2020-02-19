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

package chain_communication

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/dipperin/dipperin-core/third-party/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-event"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCsProtocolManager_cvFinderLoop(t *testing.T) {
	pm := &CsProtocolManager{}

	pm.pmType.Store(base)
}

func TestCsProtocolManager_handleInsertEvent(t *testing.T) {
	pm := &CsProtocolManager{}
	pm.pmType.Store(base)

	assert.NoError(t, pm.handleInsertEventForBft())
}

func TestCsProtocolManager_handleInsertEvent1(t *testing.T) {
	pm := &CsProtocolManager{stop: make(chan struct{})}
	pm.pmType.Store(base)

	go func() {
		time.Sleep(20 * time.Millisecond)
		pm.stop <- struct{}{}
	}()
	assert.NoError(t, pm.handleInsertEventForBft())
}

func TestCsProtocolManager_handleInsertEvent2(t *testing.T) {
	g_event.Add(g_event.NewBlockInsertEvent)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChain := NewMockChain(ctrl)
	mVReader := NewMockVerifiersReader(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mPbftNode := NewMockPbftNode(ctrl)

	pm := &CsProtocolManager{
		stop: make(chan struct{}),
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			Chain:           mChain,
			VerifiersReader: mVReader,
			MsgSigner:       mSigner,
			PbftNode:        mPbftNode,
		},
	}

	pm.pmType.Store(verifier)

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	go func() {
		time.Sleep(100 * time.Millisecond)
		//time.Sleep(500 * time.Millisecond)
		g_event.Send(g_event.NewBlockInsertEvent, *block)

		time.Sleep(100 * time.Millisecond)
		//time.Sleep(2 * time.Second)
		pm.stop <- struct{}{}
	}()

	addr := common.StringToAddress("aaa")

	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false)

	// SelfIsCurrentVerifier
	mVReader.EXPECT().ShouldChangeVerifier().Return(false)
	mVReader.EXPECT().CurrentVerifiers().Return([]common.Address{addr})
	mVReader.EXPECT().NextVerifiers().Return([]common.Address{})
	mSigner.EXPECT().GetAddress().Return(addr)

	mVReader.EXPECT().NextVerifiers().Return([]common.Address{addr})
	mVReader.EXPECT().ShouldChangeVerifier().Return(false)
	mSigner.EXPECT().GetAddress().Return(addr)

	solt := uint64(11)
	mChain.EXPECT().GetSlot(gomock.Any()).Return(&solt)

	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false)

	// SelfIsCurrentVerifier
	mVReader.EXPECT().ShouldChangeVerifier().Return(false)
	mVReader.EXPECT().CurrentVerifiers().Return([]common.Address{addr})
	mVReader.EXPECT().NextVerifiers().Return([]common.Address{})
	mSigner.EXPECT().GetAddress().Return(addr)

	mPbftNode.EXPECT().OnEnterNewHeight(gomock.Any())

	assert.NoError(t, pm.handleInsertEventForBft())

	time.Sleep(200 * time.Millisecond)
	//time.Sleep(3 * time.Second)
}

func TestNewVfFetcher(t *testing.T) {
	assert.NotNil(t, NewVfFetcher())
}

func Test_vfFetcher_loop(t *testing.T) {
	vff := NewVfFetcher()

	req := fetchReq{ReqID: 8, RespChan: make(chan *GetVerifiersResp)}

	go vff.loop()
	time.Sleep(time.Microsecond)

	vff.addReqChan <- req

	time.Sleep(2 * time.Microsecond)

	resp := GetVerifiersResp{ReqID: req.ReqID}
	vff.respChan <- resp
	time.Sleep(2 * time.Microsecond)
}

func Test_vfFetcher_OnGetVerifiersResp(t *testing.T) {
	vff := NewVfFetcher()

	req := fetchReq{ReqID: 8, RespChan: make(chan *GetVerifiersResp)}

	go vff.loop()
	time.Sleep(time.Microsecond)

	data := &GetVerifiersResp{
		ReqID: req.ReqID,
	}

	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: BootNodeVerifiersConn, Size: uint32(size), Payload: r}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)

	assert.NoError(t, vff.OnGetVerifiersResp(msg, mPeer))
}

func Test_vfFetcher_getVerifiersFromBoot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)

	vff := NewVfFetcher()
	go vff.loop()
	time.Sleep(time.Microsecond)

	req := GetVerifiersReq{}

	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("aaa"))
	assert.Nil(t, vff.getVerifiersFromBoot(req, mPeer))

	time.Sleep(2 * time.Millisecond)

	fetchConnTimeout = 2 * time.Microsecond

	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)
	assert.Nil(t, vff.getVerifiersFromBoot(req, mPeer))
}

func TestNewVFinder(t *testing.T) {
	assert.NotNil(t, NewVFinder(nil, nil,
		*chain_config.GetChainConfig()))
}

func TestVFinder_MsgHandlers(t *testing.T) {
	vf := NewVFinder(nil, nil, *chain_config.GetChainConfig())
	vf.MsgHandlers()
}

func TestVFinder_Start(t *testing.T) {
	vf := NewVFinder(nil, nil, *chain_config.GetChainConfig())
	vf.started = 1
	assert.EqualError(t, vf.Start(), g_error.ErrAlreadyStarted.Error())

	vf.started = 0

	assert.NoError(t, vf.Start())
	time.Sleep(time.Millisecond)
}

func TestVFinder_Stop(t *testing.T) {
	vf := NewVFinder(nil, nil, *chain_config.GetChainConfig())
	vf.Stop()
}

func TestVFinder_findVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	mPeer := NewMockPmAbstractPeer(ctrl)
	mChain := NewMockChain(ctrl)
	//mPs := NewMockAbstractPeerSet(ctrl)
	vf := &VFinder{
		peerManager: mPm,
		chain:       mChain,
		fetcher:     NewVfFetcher(),
	}
	mPm.EXPECT().BestPeer().Return(nil)
	vf.findingVerifiers = 1
	vf.findVerifiers()
	vf.findingVerifiers = 0
	vf.findVerifiers()

	mPm.EXPECT().BestPeer().Return(mPeer).AnyTimes()
	mPeer.EXPECT().GetHead().Return(common.Hash{}, uint64(1)).AnyTimes()
	mChain.EXPECT().CurrentBlock().Return(model.NewBlock(&model.Header{Number: 1}, nil, nil)).AnyTimes()
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	mPm.EXPECT().SelfIsCurrentVerifier().Return(true).AnyTimes()
	mPm.EXPECT().HaveEnoughVerifiers(gomock.Any()).Return(uint(0), uint(0))
	vf.findVerifiers()

	mPm.EXPECT().HaveEnoughVerifiers(gomock.Any()).Return(uint(1), uint(0))
	s := uint64(1)
	mChain.EXPECT().GetSlot(gomock.Any()).Return(&s)
	mPm.EXPECT().GetVerifierBootNode().Return(map[string]PmAbstractPeer{
		"1": mPeer,
	})
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	go func() {
		req := <-vf.fetcher.addReqChan
		assert.NotEqual(t, uint64(0), req.ReqID)
		assert.NoError(t, vf.Start())
		vf.fetcher.addReqChan <- req
		vf.fetcher.addReqChan <- req
		vf.fetcher.respChan <- GetVerifiersResp{ReqID: req.ReqID, ErrInfo: "failed"}
	}()
	vf.findVerifiers()
}

func TestVFinder_shouldFindVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	mPeer := NewMockPmAbstractPeer(ctrl)
	mChain := NewMockChain(ctrl)
	//mPs := NewMockAbstractPeerSet(ctrl)
	vf := &VFinder{
		peerManager: mPm,
		chain:       mChain,
		fetcher:     NewVfFetcher(),
	}

	mPm.EXPECT().BestPeer().Return(mPeer).AnyTimes()
	mPeer.EXPECT().GetHead().Return(common.Hash{}, uint64(1)).AnyTimes()
	mChain.EXPECT().CurrentBlock().Return(model.NewBlock(&model.Header{Number: 1}, nil, nil)).AnyTimes()
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	mPm.EXPECT().SelfIsCurrentVerifier().Return(false).AnyTimes()
	mPm.EXPECT().SelfIsNextVerifier().Return(false).AnyTimes()
	s := uint64(1)
	mChain.EXPECT().GetSlot(gomock.Any()).Return(&s).AnyTimes()

	assert.Error(t, vf.shouldFindVerifiers())
}

func TestVFinder_getVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	//mPeer := NewMockPmAbstractPeer(ctrl)
	mChain := NewMockChain(ctrl)
	//mPs := NewMockAbstractPeerSet(ctrl)
	vf := &VFinder{
		peerManager: mPm,
		chain:       mChain,
		fetcher:     NewVfFetcher(),
	}

	mChain.EXPECT().CurrentBlock().Return(model.NewBlock(&model.Header{Number: 1}, nil, nil)).AnyTimes()
	mChain.EXPECT().GetSlot(gomock.Any()).Return(nil)
	assert.Panics(t, func() {
		vf.getVerifiers(1, 1)
	})
}

func TestVFinder_getVerifiersFromBoot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	mPeer := NewMockPmAbstractPeer(ctrl)
	mChain := NewMockChain(ctrl)
	mPs := NewMockAbstractPeerSet(ctrl)
	vf := &VFinder{
		peerManager: mPm,
		chain:       mChain,
		fetcher:     NewVfFetcher(),
	}

	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)
	go func() {
		req := <-vf.fetcher.addReqChan
		assert.NotEqual(t, uint64(0), req.ReqID)
		assert.NoError(t, vf.Start())
		vf.fetcher.addReqChan <- req
		vf.fetcher.respChan <- GetVerifiersResp{
			ReqID:   req.ReqID,
			Cur:     []string{"enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@127.0.0.1:10003"},
			Next:    []string{"enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@127.0.0.1:10003"},
			ErrInfo: "",
		}
	}()
	n, _ := enode.ParseV4(fmt.Sprintf("enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@%v:%v", "127.0.0.1", 10003))
	assert.NotNil(t, n)
	mPm.EXPECT().GetSelfNode().Return(n).AnyTimes()
	mPm.EXPECT().CurrentVerifierPeersSet().Return(mPs).AnyTimes()
	mPm.EXPECT().NextVerifierPeersSet().Return(mPs).AnyTimes()
	vf.getVerifiersFromBoot(GetVerifiersReq{ID: uint64(1), CurMiss: 1}, mPeer)
}

func TestVFinder_checkAndConnectNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPs := NewMockAbstractPeerSet(ctrl)
	mPeer := NewMockPmAbstractPeer(ctrl)

	vf := &VFinder{}

	assert.Panics(t, func() {
		vf.checkAndConnectNode("a", "fdsf", mPs)
	})

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")
	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	assert.Equal(t, false, vf.checkAndConnectNode(n.ID().String(), n.String(), mPs))

	mPs.EXPECT().Peer(gomock.Any()).Return(mPeer)
	assert.Equal(t, false, vf.checkAndConnectNode("a", n.String(), mPs))

	mApm := NewMockAbsPeerManager(ctrl)
	vf.peerManager = mApm
	mApm.EXPECT().ConnectPeer(gomock.Any())
	mPs.EXPECT().Peer(gomock.Any()).Return(nil)
	assert.Equal(t, true, vf.checkAndConnectNode("a", n.String(), mPs))
}

func TestNewVFinderBoot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	mChain := NewMockChain(ctrl)

	assert.NotNil(t, NewVFinderBoot(mPm, mChain))
}

func TestVFinderBoot_MsgHandlers(t *testing.T) {
	vfb := &VFinderBoot{}
	vfb.MsgHandlers()
}

func TestVFinderBoot_onGetVerifiersReq(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	mChain := NewMockChain(ctrl)
	mPeer := NewMockPmAbstractPeer(ctrl)
	mBestPerr := NewMockPmAbstractPeer(ctrl)

	solt := uint64(10)
	req := &GetVerifiersReq{ID: 11}
	currentBlock := model.NewBlock(model.NewHeader(11, 39, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	vfb := &VFinderBoot{
		peerManager: mPm,
		chain:       mChain,
	}

	mPm.EXPECT().BestPeer().Return(mBestPerr)
	mBestPerr.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	block := model.NewBlock(model.NewHeader(11, 30, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	mChain.EXPECT().CurrentBlock().Return(block)
	vfb.onGetVerifiersReq(req, mPeer)

	mPm.EXPECT().BestPeer().Return(mBestPerr)

	mBestPerr.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false)
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().GetSlot(gomock.Any()).Return(nil)
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	assert.Panics(t, func() {
		vfb.onGetVerifiersReq(req, mPeer)
	})

	mPm.EXPECT().BestPeer().Return(mBestPerr)
	mBestPerr.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false)
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().GetSlot(gomock.Any()).Return(&solt)
	req.Slot = 9
	mPeer.EXPECT().NodeName().Return("test")
	vfb.onGetVerifiersReq(req, mPeer)

	mPm.EXPECT().BestPeer().Return(mBestPerr)
	mBestPerr.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false)
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().GetSlot(gomock.Any()).Return(&solt)
	req.Slot = 10
	mPeer.EXPECT().ID().Return("testID")
	req.CurMiss = 1
	req.NextMiss = 1
	mGpPeer1 := NewMockPmAbstractPeer(ctrl)
	mGpPeer2 := NewMockPmAbstractPeer(ctrl)
	mPs1 := NewMockAbstractPeerSet(ctrl)
	mPs2 := NewMockAbstractPeerSet(ctrl)
	mPm.EXPECT().CurrentVerifierPeersSet().Return(mPs1)
	mPs1.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"a": mGpPeer1})
	mGpPeer1.EXPECT().GetPeerRawUrl().Return("aaa@127.0.0.1:3333")
	mPm.EXPECT().NextVerifierPeersSet().Return(mPs2)
	mPs2.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"b": mGpPeer2})
	mGpPeer2.EXPECT().GetPeerRawUrl().Return("bbb@127.0.0.1:3333")
	vfb.onGetVerifiersReq(req, mPeer)

}

func TestVFinderBoot_OnGetVerifiersReq(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	mChain := NewMockChain(ctrl)
	mPeer := NewMockPmAbstractPeer(ctrl)
	mBestPerr := NewMockPmAbstractPeer(ctrl)

	solt := uint64(10)
	req := &GetVerifiersReq{ID: 11}
	currentBlock := model.NewBlock(model.NewHeader(11, 39, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	vfb := &VFinderBoot{
		peerManager: mPm,
		chain:       mChain,
	}

	req.CurMiss = 1
	req.NextMiss = 1
	req.Slot = 10
	mGpPeer1 := NewMockPmAbstractPeer(ctrl)
	mGpPeer2 := NewMockPmAbstractPeer(ctrl)
	mPs1 := NewMockAbstractPeerSet(ctrl)
	mPs2 := NewMockAbstractPeerSet(ctrl)

	mPm.EXPECT().BestPeer().Return(mBestPerr)
	mBestPerr.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false)
	mChain.EXPECT().CurrentBlock().Return(currentBlock)
	mChain.EXPECT().GetSlot(gomock.Any()).Return(&solt)
	mPeer.EXPECT().ID().Return("testID")
	mPm.EXPECT().CurrentVerifierPeersSet().Return(mPs1)
	mPs1.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"a": mGpPeer1})
	mGpPeer1.EXPECT().GetPeerRawUrl().Return("aaa@127.0.0.1:3333")
	mPm.EXPECT().NextVerifierPeersSet().Return(mPs2)
	mPs2.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"b": mGpPeer2})
	mGpPeer2.EXPECT().GetPeerRawUrl().Return("bbb@127.0.0.1:3333")
	size, r, err := rlp.EncodeToReader(req)
	assert.NoError(t, err)
	msg := p2p.Msg{Code: BootNodeVerifiersConn, Size: uint32(size), Payload: r}
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("aaa"))
	mPeer.EXPECT().NodeName().Return("aa")
	assert.Nil(t, vfb.OnGetVerifiersReq(msg, mPeer), "aaa")
}

func Test_canFind(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPm := NewMockAbsPeerManager(ctrl)
	mChain := NewMockChain(ctrl)
	mPeer := NewMockPmAbstractPeer(ctrl)

	mPm.EXPECT().BestPeer().Return(nil)
	assert.EqualError(t, canFind(mPm, mChain), g_error.ErrNoBestPeerFound.Error())

	mPm.EXPECT().BestPeer().Return(mPeer)
	mPeer.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	block := model.NewBlock(model.NewHeader(11, 30, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	mChain.EXPECT().CurrentBlock().Return(block)
	assert.EqualError(t, canFind(mPm, mChain), g_error.ErrCurHeightTooLow.Error())

	mPm.EXPECT().BestPeer().Return(mPeer)
	mPeer.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	block1 := model.NewBlock(model.NewHeader(11, 39, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	mChain.EXPECT().CurrentBlock().Return(block1)
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(true)
	assert.EqualError(t, canFind(mPm, mChain), g_error.ErrIsChangePointDoNotFind.Error())

	mPm.EXPECT().BestPeer().Return(mPeer)
	mPeer.EXPECT().GetHead().Return(common.HexToHash("a"), uint64(40))
	mChain.EXPECT().CurrentBlock().Return(block1)
	mChain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(false)
	assert.NoError(t, canFind(mPm, mChain))
}
