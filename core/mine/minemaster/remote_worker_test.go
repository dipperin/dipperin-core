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
	"github.com/dipperin/dipperin-core/tests/peer"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
)

type mockWork struct {
	fillSealResultError error
}

func (w *mockWork) GetWorkerCoinbaseAddress() common.Address {
	return common.Address{}
}

func (w *mockWork) SetWorkerCoinbaseAddress(address common.Address) {
	panic("implement me")
}

func (w *mockWork) FillSealResult(curBlock model.AbstractBlock) error {
	return w.fillSealResultError
}

func Test_newRemoteWorker(t *testing.T) {
	assert.Nil(t, newRemoteWorker(nil, common.Address{}, ""))

	p := peer_spec.PeerBuilder()

	assert.NotNil(t, newRemoteWorker(p, common.HexToAddress("0x1234"), "123"))
}

func Test_remoteWorker_SetCoinbase(t *testing.T) {
	p := peer_spec.PeerBuilder()

	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")

	assert.Equal(t, worker.CurrentCoinbaseAddress().Hex(), "0x00000000000000000000000000000000000000001234")

	worker.SetCoinbase(common.Address{})

	assert.Equal(t, worker.CurrentCoinbaseAddress().Hex(), "0x00000000000000000000000000000000000000000000")

	worker2 := remoteWorker{}

	assert.Equal(t, worker2.CurrentCoinbaseAddress().Hex(), "0x00000000000000000000000000000000000000000000")
}

func Test_remoteWorker_Start(t *testing.T) {
	p := peer_spec.PeerBuilder()

	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")

	worker.Start()

	assert.Equal(t, p.(*peer_spec.FakePeer).TestMsg, uint64(minemsg.StartMineMsg))
}

func Test_remoteWorker_Stop(t *testing.T) {
	p := peer_spec.PeerBuilder()

	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")

	worker.Stop()

	assert.Equal(t, p.(*peer_spec.FakePeer).TestMsg, uint64(minemsg.StopMineMsg))
}

func Test_remoteWorker_GetId(t *testing.T) {
	p := peer_spec.PeerBuilder()

	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")

	assert.Equal(t, string(worker.GetId()), "123")
}

func Test_remoteWorker_WaitForCommit(t *testing.T) {
	p := peer_spec.PeerBuilder()

	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")

	worker.WaitForCommit()

	assert.Equal(t, p.(*peer_spec.FakePeer).TestMsg, uint64(minemsg.WaitForCommitMsg))
}

func Test_remoteWorker_SendNewWork(t *testing.T) {
	p := peer_spec.PeerBuilder()
	w := mockWork{}
	worker := newRemoteWorker(p, common.HexToAddress("0x1234"), "123")

	worker.SendNewWork(13, &w)

	assert.Equal(t, p.(*peer_spec.FakePeer).TestMsg, uint64(13))
}
