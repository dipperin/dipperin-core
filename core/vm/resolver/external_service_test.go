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

/*func TestResolverNeedExternalService_Transfer(t *testing.T) {
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

	resp, gasLeft, err := service.Transfer(model.AliceAddr, big.NewInt(100))
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

	resp, err := service.ResolverCall(model.AliceAddr.Bytes(), []byte{123})
	assert.Error(t, err)
	assert.Equal(t, []byte(nil), resp)

	params, err := rlp.EncodeToBytes([]interface{}{"init"})
	assert.NoError(t, err)
	resp, err = service.ResolverCall(model.AliceAddr.Bytes(), params)
	assert.Equal(t, g_error.ErrFunctionInitCanNotCalled, err)
	assert.Equal(t, []byte(nil), resp)

	params, err = rlp.EncodeToBytes([]interface{}{"name"})
	assert.NoError(t, err)
	resp, err = service.ResolverCall(model.AliceAddr.Bytes(), params)
	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), resp)
}
*/
/*
func TestResolverNeedExternalService_ResolverDelegateCall(t *testing.T) {
	vmValue := &fakeVmContextService{}
	contract := &fakeContractService{}
	state := NewFakeStateDBService()
	service := &resolverNeedExternalService{
		contract,
		vmValue,
		state,
	}

	resp, err := service.ResolverDelegateCall(aliceAddr.Bytes(), []byte{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), resp)
}
*/
