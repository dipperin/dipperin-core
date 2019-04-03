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

package rpc_interface

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDipperinP2PApi_AddPeer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mp := NewMockP2PAPI(controller)
	mp.EXPECT().AddPeer(gomock.Any()).Return(nil).AnyTimes()
	mp.EXPECT().RemovePeer(gomock.Any()).Return(nil).AnyTimes()
	mp.EXPECT().AddTrustedPeer(gomock.Any()).Return(nil).AnyTimes()
	mp.EXPECT().RemoveTrustedPeer(gomock.Any()).Return(nil).AnyTimes()
	mp.EXPECT().Peers().Return(nil, nil).AnyTimes()
	mp.EXPECT().CsPmInfo().Return(nil, nil).AnyTimes()

	s := &DipperinP2PApi{service: mp}
	assert.NoError(t, s.AddPeer(""))
	assert.NoError(t, s.RemovePeer(""))
	assert.NoError(t, s.AddTrustedPeer(""))
	assert.NoError(t, s.RemoveTrustedPeer(""))
	_, err := s.Peers()
	assert.NoError(t, err)
	_, err = s.CsPmInfo()
	assert.NoError(t, err)
}
