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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	//"time"
	//"fmt"
)

func Test_makeNewTxBroadcaster(t *testing.T) {
	config := &NewTxBroadcasterConfig{}
	assert.NotNil(t, makeNewTxBroadcaster(config))
}

func Test_MsgHandlers(t *testing.T) {
	config := &NewTxBroadcasterConfig{}
	ntb := makeNewTxBroadcaster(config)
	assert.NotNil(t, ntb.MsgHandlers()[TxV1Msg])
}

func Test_BroadcastTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPM := NewMockPeerManager(ctrl)

	config := &NewTxBroadcasterConfig{
		Pm: mockPM,
	}
	ntb := makeNewTxBroadcaster(config)

	tx := NewMockAbstractTransaction(ctrl)
	tx.EXPECT().CalTxId().Return(common.HexToHash("0x123")).AnyTimes()

	txs := []model.AbstractTransaction{tx}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	peers := make(map[string]PmAbstractPeer)
	peers[mockPeer.ID()] = mockPeer

	mockPM.EXPECT().GetPeers().Return(peers).AnyTimes()
	mockPM.EXPECT().GetPeer(mockPeer.ID()).Return(mockPeer)

	ntb.BroadcastTx(txs)

	time.Sleep(100 * time.Millisecond)

}

func Test_onNewTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxPool := NewMockTxPool(ctrl)
	mockNodeConf := NewMockNodeConf(ctrl)
	mockDecoder := NewMockP2PMsgDecoder(ctrl)
	mockPM := NewMockPeerManager(ctrl)

	config := &NewTxBroadcasterConfig{
		TxPool:        mockTxPool,
		P2PMsgDecoder: mockDecoder,
		NodeConf: mockNodeConf,
		Pm: mockPM,
	}

	ntb := makeNewTxBroadcaster(config)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockDecoder.EXPECT().DecodeTxsMsg(gomock.Any()).Return(nil, errors.New("test"))
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier)

	peers := make(map[string]PmAbstractPeer)
	peers[mockPeer.ID()] = mockPeer

	mockPM.EXPECT().GetPeers().Return(peers)

	err := ntb.onNewTx(p2p.Msg{}, mockPeer)

	assert.Error(t, err)

	txs := []model.AbstractTransaction{nil}

	mockDecoder.EXPECT().DecodeTxsMsg(gomock.Any()).Return(txs, nil)

	err = ntb.onNewTx(p2p.Msg{}, mockPeer)

	assert.Error(t, err)

	tx := NewMockAbstractTransaction(ctrl)
	tx.EXPECT().CalTxId().Return(common.HexToHash("0x123")).AnyTimes()

	txs = []model.AbstractTransaction{tx}

	mockDecoder.EXPECT().DecodeTxsMsg(gomock.Any()).Return(txs, nil)
	mockTxPool.EXPECT().AddRemotes(gomock.Any()).Return([]error{errors.New("test")})

	err = ntb.onNewTx(p2p.Msg{}, mockPeer)

	assert.NoError(t, err)

	mockDecoder.EXPECT().DecodeTxsMsg(gomock.Any()).Return(txs, nil)
	mockTxPool.EXPECT().AddRemotes(gomock.Any()).Return([]error{})

	err = ntb.onNewTx(p2p.Msg{}, mockPeer)

	assert.NoError(t, err)

	time.Sleep(1 * time.Millisecond)
}

func Test_send2MinerMaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeConf := NewMockNodeConf(ctrl)
	mockPM := NewMockPeerManager(ctrl)

	config := &NewTxBroadcasterConfig{
		NodeConf: mockNodeConf,
		Pm:       mockPM,
	}

	ntb := makeNewTxBroadcaster(config)

	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfMineMaster)

	ntb.send2MinerMaster([]model.AbstractTransaction{})

	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier)

	tx := NewMockAbstractTransaction(ctrl)
	tx.EXPECT().CalTxId().Return(common.HexToHash("0x123")).AnyTimes()

	txs := []model.AbstractTransaction{tx}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockPeer.EXPECT().NodeType().Return(uint64(chain_config.NodeTypeOfVerifier))

	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeer2.EXPECT().ID().Return("2").AnyTimes()
	mockPeer2.EXPECT().NodeName().Return("tes2").AnyTimes()
	mockPeer2.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockPeer2.EXPECT().NodeType().Return(uint64(chain_config.NodeTypeOfMineMaster))

	peers := make(map[string]PmAbstractPeer)
	peers[mockPeer.ID()] = mockPeer
	peers[mockPeer2.ID()] = mockPeer2

	mockPM.EXPECT().GetPeers().Return(peers)
	mockPM.EXPECT().GetPeer(mockPeer.ID()).Return(mockPeer).AnyTimes()
	mockPM.EXPECT().GetPeer(mockPeer2.ID()).Return(mockPeer2).AnyTimes()

	ntb.send2MinerMaster(txs)

	time.Sleep(100 * time.Millisecond)
}

func Test_getReceiver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config := &NewTxBroadcasterConfig{}

	ntb := makeNewTxBroadcaster(config)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	assert.NotNil(t, ntb.getReceiver(mockPeer))
	assert.NotNil(t, ntb.getReceiver(mockPeer))
}

func Test_syncTxs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxPool := NewMockTxPool(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	config := &NewTxBroadcasterConfig{
		TxPool: mockTxPool,
		Pm: mockPM,
	}

	ntb := makeNewTxBroadcaster(config)

	txs := make(map[common.Address][]model.AbstractTransaction)

	mockTxPool.EXPECT().Pending().Return(txs, nil)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	ntb.syncTxs(mockPeer)

	tx := NewMockAbstractTransaction(ctrl)
	tx.EXPECT().Size().Return(common.StorageSize(10))
	tx.EXPECT().CalTxId().Return(common.HexToHash("0x123"))

	txs[common.HexToAddress("0x123")] = []model.AbstractTransaction{tx}

	mockTxPool.EXPECT().Pending().Return(txs, nil)
	mockPM.EXPECT().GetPeer(mockPeer.ID()).Return(mockPeer)

	ntb.syncTxs(mockPeer)

	time.Sleep(100 * time.Millisecond)
}

func Test_txSyncLoop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxPool := NewMockTxPool(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	config := &NewTxBroadcasterConfig{
		TxPool: mockTxPool,
		Pm: mockPM,
	}

	ntb := makeNewTxBroadcaster(config)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	txs := make(map[common.Address][]model.AbstractTransaction)

	tx := NewMockAbstractTransaction(ctrl)
	tx.EXPECT().Size().Return(common.StorageSize(10)).AnyTimes()
	tx.EXPECT().CalTxId().Return(common.HexToHash("0x123")).AnyTimes()

	txs[common.HexToAddress("0x123")] = []model.AbstractTransaction{tx}

	mockTxPool.EXPECT().Pending().Return(txs, nil)
	mockPM.EXPECT().GetPeer(mockPeer.ID()).Return(nil)

	ntb.syncTxs(mockPeer)

	time.Sleep(100 * time.Millisecond)

	txs[common.HexToAddress("0x122")] = []model.AbstractTransaction{tx}

	mockTxPool.EXPECT().Pending().Return(txs, nil).Times(2)
	mockPM.EXPECT().GetPeer(mockPeer.ID()).Return(mockPeer).Times(2)

	ntb.syncTxs(mockPeer)
	ntb.syncTxs(mockPeer)

	time.Sleep(100 * time.Millisecond)
}
