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

package mineworker

import (
	"bytes"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/tests/peer"
	"github.com/stretchr/testify/assert"
)

type fakeReceiver struct {
	Msg workMsg
}

func (r *fakeReceiver) OnNewWork(msg workMsg) {
	r.Msg = msg
}

func Test_newRemoteConnector(t *testing.T) {
	rc := newRemoteConnector()
	assert.NotNil(t, rc)
}

func TestRemoteConnector_SetMineMasterPeer(t *testing.T) {
	p := peer_spec.PeerBuilder()
	rc := newRemoteConnector()
	w, _ := MakeRemoteWorker(common.HexToAddress("0x123"), 1)
	w.(*worker).connector = rc
	rc.worker = w
	rc.SetMineMasterPeer(p)
	time.Sleep(1 * time.Millisecond)
}

func TestRemoteConnector_Register(t *testing.T) {
	p := peer_spec.PeerBuilder()
	rc := newRemoteConnector()
	w, _ := MakeRemoteWorker(common.HexToAddress("0x123"), 1)
	w.(*worker).connector = rc
	rc.worker = w
	rc.SetMineMasterPeer(p)
	err := rc.Register()
	assert.NoError(t, err)
	assert.Equal(t, p.(*peer_spec.FakePeer).TestMsg, uint64(minemsg.RegisterMsg))
	rc.UnRegister()
	assert.Equal(t, p.(*peer_spec.FakePeer).TestMsg, uint64(minemsg.UnRegisterMsg))
}

func TestRemoteConnector_SendMsg(t *testing.T) {
	p := peer_spec.PeerBuilder()
	rc := newRemoteConnector()
	w, _ := MakeRemoteWorker(common.HexToAddress("0x123"), 1)
	w.(*worker).connector = rc
	rc.worker = w
	rc.SetMineMasterPeer(p)
	rc.SendMsg(uint64(1), "")
	assert.Equal(t, p.(*peer_spec.FakePeer).TestMsg, uint64(1))
}

func TestRemoteConnector_OnNewMsg(t *testing.T) {
	p := peer_spec.PeerBuilder()
	rc := newRemoteConnector()
	w, _ := MakeRemoteWorker(common.HexToAddress("0x123"), 1)
	w.(*worker).connector = rc
	rc.worker = w
	receiver := fakeReceiver{}
	rc.receiver = &receiver
	rc.SetMineMasterPeer(p)

	err := rc.OnNewMsg(p2p.Msg{Code: minemsg.StartMineMsg}, p)
	assert.NoError(t, err)

	err = rc.OnNewMsg(p2p.Msg{Code: minemsg.StopMineMsg}, p)
	assert.NoError(t, err)

	err = rc.OnNewMsg(p2p.Msg{Code: minemsg.WaitForCommitMsg}, p)
	assert.NoError(t, err)

	var register minemsg.Register

	payload, _ := rlp.EncodeToBytes(register)

	err = rc.OnNewMsg(p2p.Msg{
		Code:    minemsg.NewDefaultWorkMsg,
		Payload: bytes.NewReader(payload),
	}, p)
	assert.Error(t, err)

	var dWork2 = minemsg.DefaultWork{
		WorkerCoinbaseAddress: common.HexToAddress("0x1234"),
		RlpPreCal:             []byte("123"),
		ResultNonce:           common.BlockNonce{0},
		BlockHeader:           *model.NewHeader(0, 0, common.HexToHash("0x123"), common.HexToHash("0x123"), common.HexToDiff("0x123"), big.NewInt(1), common.HexToAddress("0x123"), common.BlockNonce{0}),
	}

	payload2, _ := rlp.EncodeToBytes(dWork2)

	err = rc.OnNewMsg(p2p.Msg{
		Code:    minemsg.NewDefaultWorkMsg,
		Payload: bytes.NewReader(payload2),
	}, p)
	assert.NoError(t, err)

	time.Sleep(1 * time.Millisecond)

	assert.Equal(t, receiver.Msg.MsgCode(), minemsg.NewDefaultWorkMsg)
}
