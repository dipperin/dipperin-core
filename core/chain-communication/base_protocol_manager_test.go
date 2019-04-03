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
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"net"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/p2p"
)

type testCm struct{}

func (tcm *testCm) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	a := func(msg p2p.Msg, p PmAbstractPeer) error {
		return nil
	}

	m := make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)
	m[100] = a

	return m
}

type testCe struct{}

func (tce *testCe) Start() error {
	return nil
}

func (tce *testCe) Stop() {}

func TestBaseProtocolManager_registerCommunicationService(t *testing.T) {
	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	testCmService := &testCm{}

	bpm.registerCommunicationService(testCmService, nil)

	// panic
	assert.Panics(t, func() {
		bpm.registerCommunicationService(testCmService, nil)
	})

	testCeService := &testCe{}

	bpm.registerCommunicationService(nil, testCeService)

	assert.Equal(t, 1, len(bpm.executables))
}

func TestBaseProtocolManager_RemovePeer(t *testing.T) {
	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	assert.Panics(t, func() {
		bpm.RemovePeer("sss")
	})
}

type mockHandleMsgPeer struct{}

func (mp *mockHandleMsgPeer) NodeName() string {
	return "mock_handle_msg_peer"
}

func (mp *mockHandleMsgPeer) NodeType() uint64 {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) SendMsg(msgCode uint64, msg interface{}) error {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) ID() string {
	panic("implement me")
}

type mockMsg struct {
	I string
}

func (mp *mockHandleMsgPeer) ReadMsg() (p2p.Msg, error) {
	data := &mockMsg{I: "hhhh"}
	size, r, err := rlp.EncodeToReader(data)

	if err != nil {
		return p2p.Msg{}, err
	}

	msg := p2p.Msg{Code: 101, Size: uint32(size), Payload: r}

	return msg, nil
}

func (mp *mockHandleMsgPeer) GetHead() (common.Hash, uint64) {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) SetHead(head common.Hash, height uint64) {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) GetPeerRawUrl() string {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) DisconnectPeer() {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) RemoteVerifierAddress() (addr common.Address) {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) RemoteAddress() net.Addr {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) SetRemoteVerifierAddress(addr common.Address) {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) SetNodeName(name string) {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) SetNodeType(nt uint64) {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) SetPeerRawUrl(rawUrl string) {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) SetNotRunning() {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) IsRunning() bool {
	panic("implement me")
}

func (mp *mockHandleMsgPeer) GetCsPeerInfo() *p2p.CsPeerInfo {
	panic("implement me")
}

func TestBaseProtocolManager_handleMsg(t *testing.T) {
	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	p := &mockHandleMsgPeer{}

	//assert.NoError(t, bpm.handleMsg(p))

	assert.EqualError(t, bpm.handleMsg(p), msgHandleFuncNotFoundErr.Error())

}

func TestBaseProtocolManager_handleMsg1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	mPeer := NewMockPmAbstractPeer(ctrl)

	mPeer.EXPECT().NodeName().Return("aaaa")
	mPeer.EXPECT().ReadMsg().Return(p2p.Msg{}, errors.New("dddd"))
	mPeer.EXPECT().NodeName().Return("aaaa")
	mPeer.EXPECT().NodeName().Return("aaaa")
	assert.EqualError(t, bpm.handleMsg(mPeer), "dddd")
}

func TestBaseProtocolManager_handleMsg2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	mPeer := NewMockPmAbstractPeer(ctrl)

	mPeer.EXPECT().NodeName().Return("aaaa")

	// blockHashMsg
	data := &blockHashMsg{BlockHash: common.HexToHash("vfd"), BlockNumber: 11}
	_, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: ProtocolMaxMsgSize*10, Payload: r}

	mPeer.EXPECT().ReadMsg().Return(msg, nil)
	assert.EqualError(t, bpm.handleMsg(mPeer), msgTooLargeErr.Error())
}

func TestBaseProtocolManager_handleMsg3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	bpm.msgHandlers[VerifyBlockHashResultMsg] = func(msg p2p.Msg, p PmAbstractPeer) error {
		return errors.New("dddd")
	}

	mPeer := NewMockPmAbstractPeer(ctrl)

	mPeer.EXPECT().NodeName().Return("aaaa")

	// blockHashMsg
	data := &blockHashMsg{BlockHash: common.HexToHash("vfd"), BlockNumber: 11}
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: uint32(size), Payload: r}

	mPeer.EXPECT().ReadMsg().Return(msg, nil)
	mPeer.EXPECT().SetNotRunning()

	assert.EqualError(t, bpm.handleMsg(mPeer), "dddd")
}

func TestBaseProtocolManager_handleMsg4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	bpm.msgHandlers[VerifyBlockHashResultMsg] = func(msg p2p.Msg, p PmAbstractPeer) error {
		return nil
	}

	mPeer := NewMockPmAbstractPeer(ctrl)

	mPeer.EXPECT().NodeName().Return("aaaa")

	// blockHashMsg
	data := &blockHashMsg{BlockHash: common.HexToHash("vfd"), BlockNumber: 11}
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: VerifyBlockHashResultMsg, Size: uint32(size), Payload: r}

	mPeer.EXPECT().ReadMsg().Return(msg, nil)

	assert.Nil(t, bpm.handleMsg(mPeer))
}

type testCe1 struct{}

func (tce *testCe1) Start() error {
	return errors.New("ssss")
}

func (tce *testCe1) Stop() {}

func TestBaseProtocolManager_Start(t *testing.T) {
	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	testCmService := &testCm{}

	bpm.registerCommunicationService(testCmService, nil)

	// panic
	assert.Panics(t, func() {
		bpm.registerCommunicationService(testCmService, nil)
	})

	testCeService := &testCe{}

	bpm.registerCommunicationService(nil, testCeService)

	assert.Equal(t, 1, len(bpm.executables))

	testCeService1 := &testCe1{}

	bpm.registerCommunicationService(nil, testCeService1)
	assert.Equal(t, 2, len(bpm.executables))

	assert.EqualError(t, bpm.Start(), "ssss")
}

func TestBaseProtocolManager_Stop(t *testing.T) {
	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	testCmService := &testCm{}

	bpm.registerCommunicationService(testCmService, nil)

	// panic
	assert.Panics(t, func() {
		bpm.registerCommunicationService(testCmService, nil)
	})

	testCeService := &testCe{}

	bpm.registerCommunicationService(nil, testCeService)

	assert.Equal(t, 1, len(bpm.executables))

	testCeService1 := &testCe1{}

	bpm.registerCommunicationService(nil, testCeService1)
	assert.Equal(t, 2, len(bpm.executables))

	assert.EqualError(t, bpm.Start(), "ssss")

	bpm.Stop()
}

func TestBaseProtocolManager_validStatus(t *testing.T) {
	bpm := &BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}

	assert.Nil(t, bpm.validStatus(StatusData{}))
}

func Test_statusData_Sender(t *testing.T) {
	hsData := HandShakeData{
		ProtocolVersion:    1,
		ChainID:            big.NewInt(2),
		NetworkId:          1,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       common.HexToHash("sss"),
		NodeType:           2,
		NodeName:           "test",
		RawUrl:             "127.0.0.1",
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	assert.Equal(t, true, account.Address().IsEqual(statusData.Sender()))
}

func Test_validSign(t *testing.T) {

	hsData := HandShakeData{
		ProtocolVersion:    1,
		ChainID:            big.NewInt(2),
		NetworkId:          1,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       common.HexToHash("sss"),
		NodeType:           2,
		NodeName:           "test",
		RawUrl:             "1127.0.0.1",
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	assert.NoError(t, validSign(hash.Bytes(), statusData.PubKey, statusData.Sign))

	// error
	account2 := tests.AccFactory.GenAccount()
	assert.EqualError(t, validSign(hash.Bytes(), crypto.CompressPubkey(&account2.Pk.PublicKey), statusData.Sign), "verify signature fail")

	assert.EqualError(t, validSign(hash.Bytes(), crypto.CompressPubkey(&account2.Pk.PublicKey), []byte{}), "empty sign")
}

func Test_statusData_dataHash(t *testing.T) {
	hsData := HandShakeData{
		ProtocolVersion:    1,
		ChainID:            big.NewInt(2),
		NetworkId:          1,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       common.HexToHash("sss"),
		NodeType:           2,
		NodeName:           "test",
		RawUrl:             "127.0.0.1",
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())
}
