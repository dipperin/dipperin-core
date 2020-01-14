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
	"github.com/dipperin/dipperin-core/common"
	g_error "github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/core/model"
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
				pm := NewMockPeerManager(ctrl)
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
				pm := NewMockPeerManager(ctrl)
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

	blockMock := NewMockAbstractBlock(ctrl)
	blockMock.EXPECT().Header().Return(&model.Header{}).Times(1)
	blockMock.EXPECT().Body().Return(&model.Body{}).Times(1)

	cr := NewMockChainInterface(ctrl)
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

func TestDipperinVenusApi_GetBlockByNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *DipperinVenusApi
		expect error
	}{
		{
			name: "block not found",
			given: func() *DipperinVenusApi {
				cr := NewMockChainInterface(ctrl)
				cr.EXPECT().GetBlockByNumber(uint64(0)).Return(nil).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						ChainReader: cr,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: g_error.ErrBlockNotFound,
		},
		{
			name: "block found",
			given: func() *DipperinVenusApi {
				blockMock := NewMockAbstractBlock(ctrl)
				blockMock.EXPECT().Header().Return(&model.Header{}).Times(1)
				blockMock.EXPECT().Body().Return(&model.Body{}).Times(1)

				cr := NewMockChainInterface(ctrl)
				cr.EXPECT().GetBlockByNumber(uint64(0)).Return(blockMock).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						ChainReader: cr,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: nil,
		},
	}

	for _, tc := range testCases {
		_, err := tc.given().GetBlockByNumber(0)
		if !assert.Equal(t, tc.expect, err) {
			t.Errorf("case: %s, expect:%v, got:%v", tc.name, tc.expect, err)
		}
	}
}

func TestDipperinVenusApi_GetBlockByHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *DipperinVenusApi
		expect bool
	}{
		{
			name: "block not found",
			given: func() *DipperinVenusApi {
				cr := NewMockChainInterface(ctrl)
				cr.EXPECT().GetBlockByHash(common.Hash{}).Return(nil).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						ChainReader: cr,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: false,
		},
		{
			name: "block found",
			given: func() *DipperinVenusApi {
				blockMock := NewMockAbstractBlock(ctrl)
				blockMock.EXPECT().Header().Return(&model.Header{}).Times(1)
				blockMock.EXPECT().Body().Return(&model.Body{}).Times(1)

				cr := NewMockChainInterface(ctrl)
				cr.EXPECT().GetBlockByHash(common.Hash{}).Return(blockMock).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						ChainReader: cr,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		_, err := tc.given().GetBlockByHash(common.Hash{})
		if !assert.Equal(t, tc.expect, err == nil) {
			t.Errorf("case: %s, expect:%v, got:%v", tc.name, tc.expect, err)
		}
	}
}

func TestDipperinVenusApi_GetBlockNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var expect uint64 = 100

	cr := NewMockChainInterface(ctrl)
	cr.EXPECT().GetBlockNumber(common.Hash{}).Return(&expect).Times(1)
	api := DipperinVenusApi{service: &service.VenusFullChainService{
		DipperinConfig: &service.DipperinConfig{
			ChainReader: cr,
		},
		TxValidator: nil,
	}}

	assert.Equal(t, expect, *api.GetBlockNumber(common.Hash{}))
}

func TestDipperinVenusApi_GetGenesis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *DipperinVenusApi
		expect bool
	}{
		{
			name: "genesis not found",
			given: func() *DipperinVenusApi {
				cr := NewMockChainInterface(ctrl)
				cr.EXPECT().Genesis().Return(nil).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						ChainReader: cr,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: false,
		},
		{
			name: "genesis found",
			given: func() *DipperinVenusApi {
				blockMock := NewMockAbstractBlock(ctrl)
				blockMock.EXPECT().Header().Return(&model.Header{}).Times(1)
				blockMock.EXPECT().Body().Return(&model.Body{}).Times(1)
				cr := NewMockChainInterface(ctrl)
				cr.EXPECT().Genesis().Return(blockMock).Times(1)
				api := DipperinVenusApi{service: &service.VenusFullChainService{
					DipperinConfig: &service.DipperinConfig{
						ChainReader: cr,
					},
					TxValidator: nil,
				}}
				return &api
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		_, err := tc.given().GetGenesis()
		if !assert.Equal(t, tc.expect, err == nil) {
			t.Errorf("case: %s, expect:%v, got:%v", tc.name, tc.expect, err)
		}
	}
}

func TestDipperinVenusApi_GetBlockBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cr := NewMockChainInterface(ctrl)
	cr.EXPECT().GetBody(common.Hash{}).Return(&model.Body{}).Times(1)

	api := DipperinVenusApi{service: &service.VenusFullChainService{
		DipperinConfig: &service.DipperinConfig{
			ChainReader: cr,
		},
		TxValidator: nil,
	}}

	assert.NotNil(t, api.GetBlockBody(common.Hash{}))
}
