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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	model2 "github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestNewBftOuter(t *testing.T) {
	out := NewBftOuter(&NewBftOuterConfig{})
	assert.NotNil(t, out)
}

func TestBftOuter_SetBlockFetcher(t *testing.T) {
	out := NewBftOuter(&NewBftOuterConfig{})
	assert.NotNil(t, out)

	assert.Nil(t, out.blockFetcher)
	out.SetBlockFetcher(&BlockFetcher{})
	assert.NotNil(t, out.blockFetcher)
}

func TestBftOuter_MsgHandlers(t *testing.T) {
	out := NewBftOuter(&NewBftOuterConfig{})
	assert.NotNil(t, out)

	m := out.MsgHandlers()

	assert.Equal(t, 3, len(m))
}

func TestBftOuter_BroadcastVerifiedBlock(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("testID").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("testPeer").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Times(1)

	// pm
	mockPm := NewMockPeerManager(ctrl)
	mockPm.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"testID": mockPeer}).AnyTimes()
	mockPm.EXPECT().GetPeer("testID").Return(mockPeer).AnyTimes()

	// out
	out := NewBftOuter(&NewBftOuterConfig{Pm: mockPm})

	// vr
	vr := &model.VerifyResult{
		Block:       model2.NewBlock(model2.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil),
		SeenCommits: nil,
	}

	out.BroadcastVerifiedBlock(vr)

	time.Sleep(600 * time.Microsecond)
}

func TestBftOuter_onVerifiedResultBlockHash(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("testPeer").AnyTimes()
	mockPeer.EXPECT().ID().Return("test1").AnyTimes()
	mockPeer.EXPECT().SetHead(gomock.Any(), gomock.Any()).Times(1)

	// chain
	mockChain := NewMockChain(ctrl)
	mockChain.EXPECT().CurrentBlock().Return(model2.NewBlock(model2.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)).AnyTimes()
	mockChain.EXPECT().GetBlockByNumber(gomock.Any()).Return(nil)

	// out
	out := NewBftOuter(&NewBftOuterConfig{Chain: mockChain})
	out.blockFetcher = NewBlockFetcher(nil, nil, nil, nil)

	// blockHashMsg
	data := &blockHashMsg{BlockHash: common.HexToHash("vfd"), BlockNumber: 11}
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: uint32(size), Payload: r}

	assert.NoError(t, out.onVerifiedResultBlockHash(msg, mockPeer))
}

func TestBftOuter_onVerifiedResultBlockHash2(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("testPeer").AnyTimes()
	mockPeer.EXPECT().ID().Return("test1").AnyTimes()
	//mockPeer.EXPECT().SetHead(gomock.Any(), gomock.Any()).Times(1)

	// chain
	mockChain := NewMockChain(ctrl)
	mockChain.EXPECT().CurrentBlock().Return(model2.NewBlock(model2.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)).AnyTimes()
	//mockChain.EXPECT().GetBlockByNumber(gomock.Any()).Return(nil)

	// out
	out := NewBftOuter(&NewBftOuterConfig{Chain: mockChain})
	out.blockFetcher = NewBlockFetcher(nil, nil, nil, nil)

	// blockHashMsg
	data := &blockHashMsg{BlockHash: common.HexToHash("vfd"), BlockNumber: 9}
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: uint32(size), Payload: r}

	// error
	assert.Nil(t, out.onVerifiedResultBlockHash(msg, mockPeer))
}

func TestBftOuter_onVerifiedResultBlockHash3(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("testPeer").AnyTimes()
	mockPeer.EXPECT().ID().Return("test1").AnyTimes()
	//mockPeer.EXPECT().SetHead(gomock.Any(), gomock.Any()).Times(0)

	// chain
	mockChain := NewMockChain(ctrl)
	mockChain.EXPECT().CurrentBlock().Return(model2.NewBlock(model2.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)).AnyTimes()
	//mockChain.EXPECT().GetBlockByNumber(gomock.Any()).Return(nil)

	// out
	out := NewBftOuter(&NewBftOuterConfig{Chain: mockChain})
	out.blockFetcher = NewBlockFetcher(nil, nil, nil, nil)

	// blockHashMsg
	//data := &blockHashMsg{BlockHash:common.HexToHash("vfd"), BlockNumber:11}
	//size, r, err := rlp.EncodeToReader(data)
	//assert.NoError(t, err)
	//
	//msg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: uint32(size), Payload: r}

	errSize, errR, err := rlp.EncodeToReader("aaaa")
	assert.NoError(t, err)
	errorMsg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: uint32(errSize), Payload: errR}

	// error
	assert.Error(t, out.onVerifiedResultBlockHash(errorMsg, mockPeer))
}

func TestBftOuter_onVerifiedResultBlockHash4(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("testPeer").AnyTimes()
	mockPeer.EXPECT().ID().Return("test1").AnyTimes()
	mockPeer.EXPECT().SetHead(gomock.Any(), gomock.Any()).Times(1)

	// chain
	mockChain := NewMockChain(ctrl)
	mockChain.EXPECT().CurrentBlock().Return(model2.NewBlock(model2.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)).AnyTimes()
	mockChain.EXPECT().GetBlockByNumber(gomock.Any()).Return(model2.NewBlock(model2.NewHeader(11, 11, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil))

	// out
	out := NewBftOuter(&NewBftOuterConfig{Chain: mockChain})
	out.blockFetcher = NewBlockFetcher(nil, nil, nil, nil)

	// blockHashMsg
	data := &blockHashMsg{BlockHash: common.HexToHash("vfd"), BlockNumber: 11}
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: uint32(size), Payload: r}

	// error
	assert.Nil(t, out.onVerifiedResultBlockHash(msg, mockPeer))
}

func TestBftOuter_onGetVerifiedResult(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("testPeer").AnyTimes()
	mockPeer.EXPECT().ID().Return("test1").AnyTimes()

	// chain
	mockChain := NewMockChain(ctrl)
	mockChain.EXPECT().GetBlockByNumber(gomock.Any()).Return(model2.NewBlock(model2.NewHeader(11, 11, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)).AnyTimes()

	mockChain.EXPECT().GetSeenCommit(gomock.Any()).Return(nil).AnyTimes()

	// out
	out := NewBftOuter(&NewBftOuterConfig{Chain: mockChain})

	data := uint64(11)
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: GetVerifyResultMsg, Size: uint32(size), Payload: r}

	assert.NoError(t, out.onGetVerifiedResult(msg, mockPeer))

	data1 := struct {
		s string
	}{
		s: "ssss",
	}

	size1, r1, err := rlp.EncodeToReader(data1)
	assert.NoError(t, err)
	msgE := p2p.Msg{Code: GetVerifyResultMsg, Size: uint32(size1), Payload: r1}

	assert.NotNil(t, out.onGetVerifiedResult(msgE, mockPeer))

	mockChain.EXPECT().GetSeenCommit(gomock.Any()).Return([]model2.AbstractVerification{model2.NewVoteMsg(11, 1, common.HexToHash("asdd"), model2.AliveVerifierVoteMessage)}).AnyTimes()

	data2 := uint64(11)
	size2, r2, err := rlp.EncodeToReader(data2)
	assert.NoError(t, err)
	msg1 := p2p.Msg{Code: GetVerifyResultMsg, Size: uint32(size2), Payload: r2}
	assert.NoError(t, out.onGetVerifiedResult(msg1, mockPeer))
}

func TestBftOuter_onVerifiedResultBlock(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("test1").AnyTimes()

	// data
	data := &model.VerifyResult{
		Block: model2.NewBlock(model2.NewHeader(11, 11, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil),

		SeenCommits: []model2.AbstractVerification{model2.NewVoteMsg(11, 1, common.HexToHash("asdd"), model2.AliveVerifierVoteMessage)},
	}

	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: VerifyBlockResultMsg, Size: uint32(size), Payload: r}

	// out
	out := NewBftOuter(&NewBftOuterConfig{})

	// cause out.fetcher is nil
	assert.Panics(t, func() {
		_ = out.onVerifiedResultBlock(msg, mockPeer)
	})

	// err
	data1 := struct {
		s string
	}{
		s: "ssss",
	}

	size1, r1, err := rlp.EncodeToReader(data1)
	assert.NoError(t, err)
	msgE := p2p.Msg{Code: VerifyBlockResultMsg, Size: uint32(size1), Payload: r1}

	assert.NotNil(t, out.onVerifiedResultBlock(msgE, mockPeer))
}

func TestBftOuter_getReceiver(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	out := NewBftOuter(&NewBftOuterConfig{})
	assert.NotNil(t, out)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("testID1").Times(3)
	mockPeer.EXPECT().NodeName().Return("test1")

	assert.NotNil(t, out.getReceiver(mockPeer))

	// use cache
	mockPeer1 := NewMockPmAbstractPeer(ctrl)
	mockPeer1.EXPECT().ID().Return("testID1").Times(1)

	assert.NotNil(t, out.getReceiver(mockPeer1))
}

func TestBftOuter_newBlockReceiver(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPm := NewMockPeerManager(ctrl)

	mockPeer.EXPECT().ID().Return("testID1")
	mockPeer.EXPECT().NodeName().Return("test1")

	out := NewBftOuter(&NewBftOuterConfig{
		Pm: mockPm,
	})
	assert.NotNil(t, out)

	assert.NotNil(t, out.newBlockReceiver(mockPeer))
}

func TestBftOuter_getPeersWithoutBlock(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// init mock obj
	returnMap := make(map[string]PmAbstractPeer)

	mockPm := NewMockPeerManager(ctrl)
	mockPm.EXPECT().GetPeers().Return(returnMap).AnyTimes()

	out := NewBftOuter(&NewBftOuterConfig{
		Pm: mockPm,
	})
	assert.NotNil(t, out)

	assert.Equal(t, 0, len(out.getPeersWithoutBlock(common.HexToHash("ssss"))))

	//
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("testID1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test1")

	receiver := out.newBlockReceiver(mockPeer)
	assert.NotNil(t, receiver)

	receiver.knownBlocks.Add(common.HexToHash("123"), 1)
	out.vResultBroadcast.Store("testID1", receiver)

	returnMap["testID1"] = mockPeer

	assert.Equal(t, 0, len(out.getPeersWithoutBlock(common.HexToHash("123"))))

	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeer2.EXPECT().ID().Return("testID2").AnyTimes()
	mockPeer2.EXPECT().NodeName().Return("test2")

	returnMap["testID2"] = mockPeer2

	assert.Equal(t, 1, len(out.getPeersWithoutBlock(common.HexToHash("123"))))
}
