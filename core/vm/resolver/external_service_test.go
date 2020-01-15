// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the Dipperin-core library.
//
// The Dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The Dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package resolver

import (
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestResolverNeedExternalService_Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	service := &resolverNeedExternalService{
		contract,
		vmValue,
		state,
	}
	transferValue := big.NewInt(100)

	contract.EXPECT().Self().Return(common.AccountRef(model.AliceAddr)).AnyTimes()
	vmValue.EXPECT().Call(common.AccountRef(model.AliceAddr),model.AliceAddr,nil, uint64(0), transferValue)


	resp, gasLeft, err := service.Transfer(model.AliceAddr, transferValue)

	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), resp)
	assert.Equal(t, uint64(0), gasLeft)
}

func TestResolverNeedExternalService_ResolverCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)

	service := &resolverNeedExternalService{
		contract,
		vmValue,
		state,
	}

	contract.EXPECT().Self().Return(common.AccountRef(model.AliceAddr)).AnyTimes()
	contract.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
	contract.EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
	paramsInit, err := rlp.EncodeToBytes([]interface{}{"init"})
	assert.NoError(t, err)

	params, err := rlp.EncodeToBytes([]interface{}{"name"})
	assert.NoError(t, err)
	vmValue.EXPECT().Call(contract,model.AliceAddr,params,uint64(0), big.NewInt(100) ).Return([]byte(nil),uint64(0),nil).Times(1)


	testCases := []struct{
		name string
		given func() error
		expect error
	}{
		{
			name:"ErrFunctionInitCanNotCalled",
			given: func() error {
				resp, err := service.ResolverCall(model.AliceAddr.Bytes(), paramsInit)
				assert.Equal(t, []byte(nil), resp)
				return err
			},
			expect:gerror.ErrFunctionInitCanNotCalled,
		},
		{
			name:"ResolverCall",
			given: func() error {
				resp, err := service.ResolverCall(model.AliceAddr.Bytes(), params)
				assert.Equal(t, []byte(nil), resp)
				return err
			},
			expect:nil,
		},
	}


	for _,tc := range testCases{
		err := tc.given()
		if err != nil {
			assert.Equal(t, tc.expect,err)
		} else {
			assert.NoError(t, err)
		}
	}



}


func TestResolverNeedExternalService_ResolverDelegateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()


	testCases := []struct{
		name string
		given func() []byte
		expect []byte
	}{
		{
			name: "ResolverDelegateCall",
			given: func() []byte {
				state := NewMockStateDBService(ctrl)
				vmValue := NewMockVmContextService(ctrl)
				contract := NewMockContractService(ctrl)

				service := &resolverNeedExternalService{
					contract,
					vmValue,
					state,
				}

				params, err := rlp.EncodeToBytes([]interface{}{"name"})
				assert.NoError(t, err)

				contract.EXPECT().Self().Return(common.AccountRef(model.AliceAddr)).AnyTimes()
				contract.EXPECT().GetGas().Return(uint64(0)).AnyTimes()
				contract.EXPECT().CallValue().Return(big.NewInt(100)).AnyTimes()
				vmValue.EXPECT().Call(contract,model.AliceAddr,params,uint64(0), big.NewInt(100) ).Return([]byte(nil),uint64(0),nil).AnyTimes()
				vmValue.EXPECT().DelegateCall(contract, model.AliceAddr,params,uint64(0)).Return([]byte(nil),uint64(0),nil)
				resp, err := service.ResolverDelegateCall(model.AliceAddr.Bytes(), params)
				assert.NoError(t, err)
				return resp
			},
			expect:[]byte(nil),
		},
	}

	for _,tc := range testCases{
		assert.Equal(t, tc.expect, tc.given())
	}
}

