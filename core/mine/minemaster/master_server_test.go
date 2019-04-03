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

package minemaster

import (
	"bytes"
	"errors"
	"github.com/dipperin/dipperin-core/tests/peer"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"sync/atomic"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
)

var fakeBlock *model.Block

func fakeGetCurWorkBlockFunc() model.AbstractBlock {
	return fakeBlock
}

func Test_server_RegisterWorker(t *testing.T) {
	m := testMasterBuilder(testMineConfig)
	dispatch := mockDispatch{curBlock: nil}
	m.setWorkDispatcher(&dispatch)
	wm := newDefaultWorkManager(testMineConfig)
	s := newServer(m, wm, fakeGetCurWorkBlockFunc)

	p := peer_spec.PeerBuilder()

	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")
	s.master.(*master).Start()
	s.RegisterWorker(worker)
	s.UnRegisterWorker("123")
}

func Test_server_ReceiveMsg(t *testing.T) {
	testMineConfig := MineConfig{
		CoinbaseAddress:  &atomic.Value{},
		BlockBroadcaster: fakeBlockBroadcaster{},
	}
	m := testMasterBuilder(testMineConfig)
	dispatch := mockDispatch{curBlock: nil}
	m.setWorkDispatcher(&dispatch)
	wm := newDefaultWorkManager(testMineConfig)
	s := newServer(m, wm, fakeGetCurWorkBlockFunc)
	s.master.(*master).Start()
	p := peer_spec.PeerBuilder()
	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")
	s.RegisterWorker(worker)

	s.ReceiveMsg("123", uint64(11), nil)

	s.ReceiveMsg("123", minemsg.SubmitDefaultWorkMsg, nil)

	work := mockWork{}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock = factory.CreateBlock2(diff, 1)

	s.ReceiveMsg("123", minemsg.SubmitDefaultWorkMsg, &work)

	model.CalNonce(fakeBlock)

	s.ReceiveMsg("123", minemsg.SubmitDefaultWorkMsg, &work)

	work2 := mockWork{fillSealResultError: errors.New("test")}

	s.ReceiveMsg("123", minemsg.SubmitDefaultWorkMsg, &work2)
}

func Test_server_SetMineMasterPeer(t *testing.T) {

	m := testMasterBuilder(testMineConfig)
	dispatch := mockDispatch{curBlock: nil}
	m.setWorkDispatcher(&dispatch)
	wm := newDefaultWorkManager(testMineConfig)
	s := newServer(m, wm, fakeGetCurWorkBlockFunc)

	p := peer_spec.PeerBuilder()

	s.SetMineMasterPeer(p)
}

func Test_server_OnNewMsg(t *testing.T) {
	testMineConfig := MineConfig{
		CoinbaseAddress:  &atomic.Value{},
		BlockBroadcaster: fakeBlockBroadcaster{},
	}
	m := testMasterBuilder(testMineConfig)
	dispatch := mockDispatch{curBlock: nil}
	m.setWorkDispatcher(&dispatch)
	wm := newDefaultWorkManager(testMineConfig)
	s := newServer(m, wm, fakeGetCurWorkBlockFunc)
	s.master.(*master).Start()
	p := peer_spec.PeerBuilder()
	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")
	s.RegisterWorker(worker)

	var register minemsg.Register

	payload, _ := rlp.EncodeToBytes(register)

	err := s.OnNewMsg(p2p.Msg{
		Code:    minemsg.RegisterMsg,
		Payload: bytes.NewReader(payload),
	}, p)

	assert.Error(t, err)

	var register2 = minemsg.Register{Coinbase: common.HexToAddress("0x213")}

	payload2, _ := rlp.EncodeToBytes(register2)

	err = s.OnNewMsg(p2p.Msg{
		Code:    minemsg.RegisterMsg,
		Payload: bytes.NewReader(payload2),
		Size:    uint32(len(payload2)),
	}, p)

	assert.NoError(t, err)

	err = s.OnNewMsg(p2p.Msg{
		Code:    minemsg.UnRegisterMsg,
		Payload: bytes.NewReader(payload),
		Size:    uint32(len(payload)),
	}, p)

	assert.NoError(t, err)

	err = s.OnNewMsg(p2p.Msg{
		Code:    minemsg.SetCurrentCoinbaseMsg,
		Payload: bytes.NewReader(payload),
		Size:    uint32(len(payload)),
	}, p)

	assert.NoError(t, err)

	err = s.OnNewMsg(p2p.Msg{
		Code:    minemsg.SubmitDefaultWorkMsg,
		Payload: bytes.NewReader(payload),
		Size:    uint32(len(payload)),
	}, p)

	assert.Error(t, err)

	work := minemsg.DefaultWork{
		WorkerCoinbaseAddress: common.HexToAddress("0x1234"),
		RlpPreCal:             []byte("123"),
		ResultNonce:           common.BlockNonce{0},
		BlockHeader:           *model.NewHeader(0, 0, common.HexToHash("0x123"), common.HexToHash("0x123"), common.HexToDiff("0x123"), big.NewInt(1), common.HexToAddress("0x123"), common.BlockNonce{0}),
	}

	payload3, _ := rlp.EncodeToBytes(work)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock = factory.CreateBlock2(diff, 1)

	err = s.OnNewMsg(p2p.Msg{
		Code:    minemsg.SubmitDefaultWorkMsg,
		Payload: bytes.NewReader(payload3),
		Size:    uint32(len(payload3)),
	}, p)

	assert.NoError(t, err)

	err = s.OnNewMsg(p2p.Msg{
		Code:    uint64(33),
		Payload: bytes.NewReader(payload),
		Size:    uint32(len(payload)),
	}, p)

	assert.Error(t, err)
}
