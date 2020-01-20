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
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestNewVfFetcher(t *testing.T) {
	assert.NotNil(t, NewVfFetcher())
}

func TestNewVFinder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockChain(ctrl)

	pm := NewMockAbsPeerManager(ctrl)

	cfg := chainconfig.ChainConfig{}

	assert.NotNil(t, NewVFinder(c, pm, cfg))
}

func TestVFinder_MsgHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := NewMockChain(ctrl)

	pm := NewMockAbsPeerManager(ctrl)

	cfg := chainconfig.ChainConfig{}

	vf := NewVFinder(c, pm, cfg)

	if assert.NotNil(t, vf) {
		assert.NotNil(t, vf.MsgHandlers()[BootNodeVerifiersConn])
	}
}

func TestVFinder_Start(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() *VFinder
		expect bool // exist error
	}{
		{
			name: "already started",
			given: func() *VFinder {
				vf := &VFinder{}
				atomic.CompareAndSwapUint32(&vf.started, 0, 1)
				return vf
			},
			expect: true,
		},
		{
			name: "already started",
			given: func() *VFinder {
				return &VFinder{
					fetcher: NewVfFetcher(),
				}
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().Start() != nil) {
			t.Errorf("case:%s, expect:%v", tc.name, tc.expect)
		}
	}
}
