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
	"github.com/dipperin/dipperin-core/third_party/crypto"
	cs_crypto "github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type mockCommunicationService struct {
	handlers map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
}

func (m mockCommunicationService) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	return m.handlers
}

type mockCommunicationExecutable struct {
	startErr error
	stopErr  error
}

func (m mockCommunicationExecutable) Start() error {
	return m.startErr
}

func (m mockCommunicationExecutable) Stop() {
	if m.stopErr != nil {
		panic(m.stopErr)
	}
}

func TestBaseProtocolManager_registerCommunicationService(t *testing.T) {
	expect1 := map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error{
		15: func(msg p2p.Msg, p PmAbstractPeer) error {
			return nil
		},
	}

	cs := mockCommunicationService{handlers: expect1}
	ex := mockCommunicationExecutable{}

	bpm := BaseProtocolManager{msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error)}
	bpm.registerCommunicationService(cs, ex)

	assert.Equal(t, len(expect1), len(bpm.msgHandlers))
	assert.Equal(t, 1, len(bpm.executables))
}

func TestStatusData_Sender(t *testing.T) {
	testCases := []struct {
		name           string
		givenAndExpect func() (*StatusData, common.Address)
	}{
		{
			name: "verify signature failed",
			givenAndExpect: func() (*StatusData, common.Address) {
				return &StatusData{
					HandShakeData: HandShakeData{},
					PubKey:        []byte{},
					Sign:          []byte{},
				}, common.Address{}
			},
		},
		{
			name: "get sender success",
			givenAndExpect: func() (*StatusData, common.Address) {
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
				pk, _ := crypto.GenerateKey()
				statusData.Sign, _ = crypto.Sign(statusData.DataHash().Bytes(), pk)
				statusData.PubKey = crypto.CompressPubkey(&pk.PublicKey)
				addr := cs_crypto.GetNormalAddress(pk.PublicKey)
				return statusData, addr
			},
		},
	}

	for _, tc := range testCases {
		given, expect := tc.givenAndExpect()
		addr := given.Sender()
		if !assert.True(t, expect.IsEqual(addr)) {
			t.Errorf("case:%s, expect:%+v, got:%+v", tc.name, expect, addr)
		}
	}
}

func TestStatusData_DataHash(t *testing.T) {
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

	assert.False(t, hash.IsEmpty())
}
