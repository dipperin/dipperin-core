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

package rpcinterface

import (
	"errors"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDipperinP2PApi_AddPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := "AddPeer"
	expect := errors.New("AddPeerError")

	p2pAPIMock := NewMockP2PAPI(ctrl)
	p2pAPIMock.EXPECT().AddPeer(input).Return(expect).Times(1)

	api := DipperinP2PApi{service: p2pAPIMock}

	assert.Equal(t, expect, api.AddPeer(input))
}

func TestDipperinP2PApi_RemovePeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := "RemovePeer"
	expect := errors.New("RemovePeerError")

	p2pAPIMock := NewMockP2PAPI(ctrl)
	p2pAPIMock.EXPECT().RemovePeer(input).Return(expect).Times(1)

	api := DipperinP2PApi{service: p2pAPIMock}

	assert.Equal(t, expect, api.RemovePeer(input))
}

func TestDipperinP2PApi_AddTrustedPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := "AddTrustedPeer"
	expect := errors.New("AddTrustedPeerError")

	p2pAPIMock := NewMockP2PAPI(ctrl)
	p2pAPIMock.EXPECT().AddTrustedPeer(input).Return(expect).Times(1)

	api := DipperinP2PApi{service: p2pAPIMock}

	assert.Equal(t, expect, api.AddTrustedPeer(input))
}

func TestDipperinP2PApi_RemoveTrustedPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	input := "RemoveTrustedPeer"
	expect := errors.New("RemoveTrustedPeerError")

	p2pAPIMock := NewMockP2PAPI(ctrl)
	p2pAPIMock.EXPECT().RemoveTrustedPeer(input).Return(expect).Times(1)

	api := DipperinP2PApi{service: p2pAPIMock}

	assert.Equal(t, expect, api.RemoveTrustedPeer(input))
}

func TestDipperinP2PApi_Peers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expect1 := make([]*p2p.PeerInfo, 0)
	expect2 := errors.New("PeersError")

	p2pAPIMock := NewMockP2PAPI(ctrl)
	p2pAPIMock.EXPECT().Peers().Return(expect1, expect2).Times(1)

	api := DipperinP2PApi{service: p2pAPIMock}
	pInfos, err := api.Peers()

	assert.Equal(t, len(expect1), len(pInfos))
	assert.Equal(t, expect2, err)
}

func TestDipperinP2PApi_CsPmInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expect1 := &p2p.CsPmPeerInfo{}
	expect2 := errors.New("CsPmInfoError")

	p2pAPIMock := NewMockP2PAPI(ctrl)
	p2pAPIMock.EXPECT().CsPmInfo().Return(expect1, expect2).Times(1)

	api := DipperinP2PApi{service: p2pAPIMock}
	pInfo, err := api.CsPmInfo()

	assert.Equal(t, expect1, pInfo)
	assert.Equal(t, expect2, err)
}
