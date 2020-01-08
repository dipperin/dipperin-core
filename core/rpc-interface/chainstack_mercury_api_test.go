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
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/mock/chain-communication-mock"
	middleware_mock "github.com/dipperin/dipperin-core/tests/mock/cs-chain-mock/chain-writer-mock/middleware-mock"
	model_mock "github.com/dipperin/dipperin-core/tests/mock/model-mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDipperinVenusApi_GetSyncStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *DipperinVenusApi
		expect bool
	}{
		{
			name: "sync status is false",
			given: func() *DipperinVenusApi {
				pm := chain_communication_mock.NewMockPeerManager(ctrl)
				pm.EXPECT().IsSync().Return(false).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						NormalPm: pm,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: false,
		},
		{
			name: "sync status is true",
			given: func() *DipperinVenusApi {
				pm := chain_communication_mock.NewMockPeerManager(ctrl)
				pm.EXPECT().IsSync().Return(true).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						NormalPm: pm,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().GetSyncStatus()) {
			t.Errorf("case: %s, expect:%v", tc.name, tc.expect)
		}
	}
}

func TestDipperinVenusApi_CurrentBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	blockMock := model_mock.NewMockAbstractBlock(ctrl)
	blockMock.EXPECT().Header().Return(&model.Header{}).Times(1)
	blockMock.EXPECT().Body().Return(&model.Body{}).Times(1)

	cr := middleware_mock.NewMockChainInterface(ctrl)
	cr.EXPECT().CurrentBlock().Return(blockMock).Times(1)

	api := DipperinVenusApi{service: &service.VenusFullChainService{
		DipperinConfig: &service.DipperinConfig{
			ChainReader: cr,
		},
		TxValidator: nil,
	}}

	resp, err := api.CurrentBlock()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
