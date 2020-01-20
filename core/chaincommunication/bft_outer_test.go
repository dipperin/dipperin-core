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

package chaincommunication

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	model2 "github.com/dipperin/dipperin-core/core/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestNewBftOuter(t *testing.T) {
	bo := NewBftOuter(&NewBftOuterConfig{})
	if assert.NotNil(t, bo) {
		assert.Equal(t, 3, len(bo.handlers))
	}
}

func TestBftOuter_SetBlockFetcher(t *testing.T) {
	bo := NewBftOuter(&NewBftOuterConfig{})
	if assert.Nil(t, bo.blockFetcher) {
		bo.SetBlockFetcher(&BlockFetcher{})
		assert.NotNil(t, bo.blockFetcher)
	}
}

func TestBftOuter_MsgHandlers(t *testing.T) {
	bo := NewBftOuter(&NewBftOuterConfig{})
	assert.Equal(t, len(bo.handlers), len(bo.MsgHandlers()))
}

func TestBftOuter_BroadcastVerifiedBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("testID").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("testPeer").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any())

	mockPm := NewMockPeerManager(ctrl)
	mockPm.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"testID": mockPeer}).AnyTimes()
	mockPm.EXPECT().GetPeer("testID").Return(mockPeer).AnyTimes()

	bftOuter := NewBftOuter(&NewBftOuterConfig{Pm: mockPm})

	header := model2.NewHeader(11, 101,
		common.HexToHash("ss"), common.HexToHash("fdfs"),
		common.StringToDiff("0x22"), big.NewInt(111),
		common.StringToAddress("fdsfds"), common.EncodeNonce(33))
	vr := &model.VerifyResult{
		Block:       model2.NewBlock(header, nil, nil),
		SeenCommits: nil,
	}

	bftOuter.BroadcastVerifiedBlock(vr) // TODO: Verification results are required. include sub-methods

	time.Sleep(600 * time.Microsecond)
}
