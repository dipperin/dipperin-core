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
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMakeNewPbftDownloader(t *testing.T) {
	pbftDownloaderConfig := &NewPbftDownloaderConfig{}
	assert.NotNil(t, MakeNewPbftDownloader(pbftDownloaderConfig))
}

func TestNewPbftDownloader_MsgHandlers(t *testing.T) {
	pbftDownloaderConfig := &NewPbftDownloaderConfig{}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)

	handles := pbftDownloader.MsgHandlers()

	assert.NotNil(t, handles[GetBlocksMsg])
	assert.NotNil(t, handles[BlocksMsg])
}

func TestNewPbftDownloader_onGetBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)

	fetcher := getFetcher()

	pbftDownloaderConfig := &NewPbftDownloaderConfig{
		Chain:   mockChain,
		fetcher: fetcher,
	}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	err := pbftDownloader.onGetBlocks(p2p.Msg{
		Payload: bytes.NewReader([]byte{}),
	}, mockPeer)

	assert.Error(t, err)

	fakeBlock1 := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)

	mockChain.EXPECT().GetBlockByNumber(uint64(1)).Return(fakeBlock1).Times(1)
	mockChain.EXPECT().GetBlockByNumber(uint64(2)).Return(nil).Times(1)
	mockChain.EXPECT().GetSeenCommit(uint64(1)).Return(nil).Times(1)
	mockChain.EXPECT().GetSeenCommit(uint64(2)).Return(nil).Times(1)

	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockGetBlockHeaders, _ := rlp.EncodeToBytes(&getBlockHeaders{
		Amount:       2,
		OriginHeight: 1,
	})

	msg := p2p.Msg{
		Payload: bytes.NewBuffer(mockGetBlockHeaders),
	}

	err = pbftDownloader.onGetBlocks(msg, mockPeer)

	assert.NoError(t, err)

}

func TestNewPbftDownloader_onBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)

	fetcher := getFetcher()

	pbftDownloaderConfig := &NewPbftDownloaderConfig{
		Chain:   mockChain,
		fetcher: fetcher,
	}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().ID().Return("1").AnyTimes()

	err := pbftDownloader.onBlocks(p2p.Msg{
		Payload: bytes.NewReader([]byte{}),
	}, mockPeer)

	assert.Error(t, err)

	mockCatchRlp, _ := rlp.EncodeToBytes(&[]catchupRlp{})

	msg := p2p.Msg{
		Payload: bytes.NewReader(mockCatchRlp),
	}

	err = pbftDownloader.onBlocks(msg, mockPeer)

	assert.NoError(t, err)

	fakeBlock1 := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)

	mockCatchRlp, _ = rlp.EncodeToBytes(&[]catchupRlp{{Block: fakeBlock1}})

	msg = p2p.Msg{
		Payload: bytes.NewReader(mockCatchRlp),
	}

	go func() {
		err := pbftDownloader.onBlocks(msg, mockPeer)
		assert.Error(t, err)
	}()

	fetcher.Start()

	time.Sleep(1 * time.Millisecond)

	pbftDownloader.Stop()

	mockChain2 := NewMockChain(ctrl)

	fetcher2 := getFetcher()

	pbftDownloaderConfig2 := &NewPbftDownloaderConfig{
		Chain:   mockChain2,
		fetcher: fetcher2,
	}

	pbftDownloader2 := MakeNewPbftDownloader(pbftDownloaderConfig2)

	msg2 := p2p.Msg{
		Payload: bytes.NewReader(mockCatchRlp),
	}

	go func() {
		err := pbftDownloader2.onBlocks(msg2, mockPeer)
		assert.NoError(t, err)
	}()

	fetcher2.Start()

	time.Sleep(1 * time.Millisecond)
	fmt.Println(reflect.ValueOf(pbftDownloader2).Pointer(), reflect.ValueOf(pbftDownloader).Pointer())
	assert.NotNil(t, <-pbftDownloader2.blockC)

}

func TestNewPbftDownloader_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPM := NewMockPeerManager(ctrl)
	mockChain := NewMockChain(ctrl)

	fetcher := getFetcher()

	pbftDownloaderConfig := &NewPbftDownloaderConfig{
		Chain:   mockChain,
		Pm:      mockPM,
		fetcher: fetcher,
	}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)

	mockPM.EXPECT().BestPeer().Return(nil).AnyTimes()

	pollingInterval = 1 * time.Millisecond

	assert.NoError(t, pbftDownloader.Start())

	time.Sleep(5 * time.Millisecond)

	pbftDownloader.Stop()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x1234"), uint64(3)).AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	mockPM.EXPECT().BestPeer().Return(mockPeer).AnyTimes()

	fakeBlock1 := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)

	mockChain.EXPECT().CurrentBlock().Return(fakeBlock1).AnyTimes()

	pbftDownloader2 := MakeNewPbftDownloader(pbftDownloaderConfig)

	assert.NoError(t, pbftDownloader2.Start())

	time.Sleep(5 * time.Millisecond)

	pbftDownloader2.Stop()

}

func TestNewPbftDownloader_getBestPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPM := NewMockPeerManager(ctrl)
	mockChain := NewMockChain(ctrl)

	fetcher := getFetcher()

	pbftDownloaderConfig := &NewPbftDownloaderConfig{
		Chain:   mockChain,
		Pm:      mockPM,
		fetcher: fetcher,
	}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)

	mockPM.EXPECT().BestPeer().Return(nil).Times(1)

	assert.Equal(t, pbftDownloader.getBestPeer(), nil)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	mockPM.EXPECT().BestPeer().Return(mockPeer).AnyTimes()

	mockChain.EXPECT().CurrentBlock().Return(nil).Times(1)

	assert.Equal(t, pbftDownloader.getBestPeer(), nil)

	fakeBlock1 := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)

	mockChain.EXPECT().CurrentBlock().Return(fakeBlock1).AnyTimes()

	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x123"), uint64(1)).Times(1)

	assert.Equal(t, pbftDownloader.getBestPeer(), nil)

	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x123"), uint64(2)).Times(1)

	assert.Equal(t, pbftDownloader.getBestPeer(), mockPeer)

}

func TestNewPbftDownloader_runSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPM := NewMockPeerManager(ctrl)
	mockChain := NewMockChain(ctrl)

	fetcher := getFetcher()

	pbftDownloaderConfig := &NewPbftDownloaderConfig{
		Chain:   mockChain,
		Pm:      mockPM,
		fetcher: fetcher,
	}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)

	pbftDownloader.synchronising = int32(1)

	pbftDownloader.runSync()

	pbftDownloader.synchronising = int32(0)

	fakeBlock1 := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)

	mockCatchupRlp := &catchupRlp{
		Block:      fakeBlock1,
		SeenCommit: nil,
	}

	mockNpbPack := &npbPack{
		peerID: "string",
		blocks: []*catchupRlp{mockCatchupRlp},
	}

	go func() { pbftDownloader.blockC <- mockNpbPack }()

	time.Sleep(1 * time.Millisecond)

	mockPM.EXPECT().BestPeer().Return(nil).Times(1)

	pbftDownloader.runSync()
}

func TestNewPbftDownloader_fetchBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPM := NewMockPeerManager(ctrl)
	mockChain := NewMockChain(ctrl)

	fetcher := getFetcher()

	pbftDownloaderConfig := &NewPbftDownloaderConfig{
		Chain:   mockChain,
		Pm:      mockPM,
		fetcher: fetcher,
	}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)

	mockPeer := NewMockPmAbstractPeer(ctrl)

	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x123"), uint64(2)).Times(2)
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockPeer.EXPECT().IsRunning().Return(false).AnyTimes()

	fakeBlock1 := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)
	mockChain.EXPECT().CurrentBlock().Return(fakeBlock1).AnyTimes()

	fakeBlock2 := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)
	mockChain.EXPECT().GetBlockByNumber(uint64(1)).Return(fakeBlock1).AnyTimes()

	mockCatchupRlp := &catchupRlp{
		Block:      fakeBlock2,
		SeenCommit: nil,
	}

	mockNpbPack := &npbPack{
		peerID: "2",
		blocks: []*catchupRlp{mockCatchupRlp},
	}

	mockNpbPack2 := &npbPack{
		peerID: "1",
		blocks: []*catchupRlp{},
	}

	go func() {
		pbftDownloader.blockC <- mockNpbPack
		pbftDownloader.blockC <- mockNpbPack2
	}()

	time.Sleep(1 * time.Millisecond)

	pbftDownloader.fetchBlocks(mockPeer)

	mockChain.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)

	mockNpbPack3 := &npbPack{
		peerID: "1",
		blocks: []*catchupRlp{mockCatchupRlp},
	}

	go func() {
		pbftDownloader.blockC <- mockNpbPack3
	}()

	time.Sleep(1 * time.Millisecond)

	pbftDownloader.fetchBlocks(mockPeer)

	mockChain.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x123"), uint64(1)).Times(1)

	go func() {
		pbftDownloader.blockC <- mockNpbPack3
	}()

	time.Sleep(1 * time.Millisecond)

	pbftDownloader.fetchBlocks(mockPeer)

	mockChain.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x123"), uint64(2)).Times(1)

	go func() {
		pbftDownloader.blockC <- mockNpbPack3
	}()

	time.Sleep(1 * time.Millisecond)

	go func() {
		pbftDownloader.fetchBlocks(mockPeer)
	}()

	time.Sleep(1 * time.Millisecond)

	pbftDownloader.Stop()

	pbftDownloader2 := MakeNewPbftDownloader(pbftDownloaderConfig)

	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x123"), uint64(2)).Times(1)

	fetchBlockTimeout = 1 * time.Millisecond

	go func() {
		pbftDownloader2.fetchBlocks(mockPeer)
	}()

	time.Sleep(10 * time.Millisecond)
}

func TestNewPbftDownloader_importBlockResults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPM := NewMockPeerManager(ctrl)
	mockChain := NewMockChain(ctrl)

	fetcher := getFetcher()

	pbftDownloaderConfig := &NewPbftDownloaderConfig{
		Chain:   mockChain,
		Pm:      mockPM,
		fetcher: fetcher,
	}

	pbftDownloader := MakeNewPbftDownloader(pbftDownloaderConfig)

	mockChain.EXPECT().SaveBlock(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	fakeBlock := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)

	mockCatchupRlp := &catchupRlp{
		Block:      fakeBlock,
		SeenCommit: []*model.VoteMsg{model.NewVoteMsg(uint64(1), uint64(1), common.HexToHash("0x123"), 0)},
	}

	assert.NoError(t, pbftDownloader.importBlockResults([]*catchupRlp{mockCatchupRlp}))

}

func Test_catchup_DecodeRLP(t *testing.T) {
	c := &catchup{}

	size, reader, _ := rlp.EncodeToReader([]byte{})

	stream := rlp.NewStream(reader, uint64(size))

	err := c.DecodeRLP(stream)

	assert.Error(t, err)

	fakeBlock := factory.CreateBlock2(common.HexToDiff("0x1effffff"), 1)

	mockCatchup := &catchup{
		Block:      fakeBlock,
		SeenCommit: nil,
	}

	size, reader, _ = rlp.EncodeToReader(mockCatchup)

	stream = rlp.NewStream(reader, uint64(size))

	err = c.DecodeRLP(stream)

	assert.NoError(t, err)
}
