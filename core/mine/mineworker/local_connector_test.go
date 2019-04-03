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
	"github.com/dipperin/dipperin-core/tests/mine"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/stretchr/testify/assert"
)

func TestLocalWorkMsg_MsgCode(t *testing.T) {
	msg := localWorkMsg{code: minemsg.NewDefaultWorkMsg}
	assert.Equal(t, minemsg.NewDefaultWorkMsg, msg.MsgCode())
}

func TestLocalWorkMsg_Decode(t *testing.T) {
	dWork := &minemsg.DefaultWork{WorkerCoinbaseAddress: common.HexToAddress("0xaaabb123123")}
	msg := &localWorkMsg{code: minemsg.NewDefaultWorkMsg, work: dWork}
	var tmpWork minemsg.DefaultWork
	assert.EqualError(t, msg.Decode(tmpWork), WorkMsgDecodeResultShouldBePtrErr.Error())
	assert.NoError(t, msg.Decode(&tmpWork))
	assert.Equal(t, common.HexToAddress("0xaaabb123123"), tmpWork.WorkerCoinbaseAddress)

	var fw fakeWork
	assert.EqualError(t, msg.Decode(&fw), WorkMsgDecodeShouldBeSameTypeErr.Error())
}

func Test_newLocalConnector(t *testing.T) {
	assert.NotNil(t, newLocalConnector("123", mine_spec.MasterServerBuilder()))
}

func Test_localConnector_Register(t *testing.T) {
	ms := mine_spec.MasterServerBuilder()
	lc := newLocalConnector("123", ms)

	lc.Register()

	assert.NotNil(t, ms.(*mine_spec.FakeMasterServer).Workers["123"])

	lc.UnRegister()

	assert.Nil(t, ms.(*mine_spec.FakeMasterServer).Workers["123"])
}

func Test_localConnector_SendMsg(t *testing.T) {
	ms := mine_spec.MasterServerBuilder()
	lc := newLocalConnector("123", ms)
	assert.NoError(t, lc.SendMsg(uint64(10), "test"))
}

func Test_localConnector_WaitForCommit(t *testing.T) {
	ms := mine_spec.MasterServerBuilder()
	lc := newLocalConnector("123", ms)
	lc.WaitForCommit()
}

func Test_localConnector_SetCoinbase(t *testing.T) {
	ms := mine_spec.MasterServerBuilder()
	lc := newLocalConnector("123", ms)
	lc.SetCoinbase(common.HexToAddress("0x123"))
}

func Test_localConnector_Start(t *testing.T) {
	ms := mine_spec.MasterServerBuilder()
	lc := newLocalConnector("123", ms)
	w := newWorker(common.HexToAddress("0x123"), 1, lc)
	lc.worker = w

	lc.Start()
	lc.Stop()
}

func Test_localConnector_SendNewWork(t *testing.T) {
	ms := mine_spec.MasterServerBuilder()
	lc := newLocalConnector("123", ms)
	dWork := &minemsg.DefaultWork{WorkerCoinbaseAddress: common.HexToAddress("0xaaabb123123")}
	w := newWorker(common.HexToAddress("0x123"), 1, lc)
	manager := newWorkManager(lc, w.Miners, w.CurrentCoinbaseAddress)
	lc.receiver = manager
	lc.SendNewWork(1, dWork)
}

func Test_localConnector_CurrentCoinbaseAddress(t *testing.T) {
	ms := mine_spec.MasterServerBuilder()
	lc := newLocalConnector("123", ms)
	w := newWorker(common.HexToAddress("0x123"), 1, lc)
	lc.worker = w

	assert.Equal(t, lc.CurrentCoinbaseAddress().Hex(), "0x00000000000000000000000000000000000000000123")
}
