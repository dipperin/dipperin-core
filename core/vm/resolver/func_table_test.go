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
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewResolver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contract, state)

	vmValue.EXPECT().GetGasPrice().Return(model.TestGasPrice).AnyTimes()

	testCases := []struct{
		name string
		given func() (string,string,bool)
		expect error
	}{
		{
			name:"wrongModule",
			given: func() (string,string,bool) {
				return "module", "field",false
			},
			expect: nil,
		},
		{
			name:"wrongField",
			given: func() (string,string,bool) {
				return "env", "field",false
			},
			expect:nil,
		},
		{
		    name:"ResolveFuncRight",
			given: func() (string,string,bool) {
				return "env", "gasPrice",true
			},
		    expect:nil,
		},
	}

	for _,tc := range testCases{
		module, field, canExecute := tc.given()
		resolverFunc := solver.ResolveFunc(module, field)
		if(canExecute){
			gasPrice := resolverFunc.Execute(&exec.VirtualMachine{})
			cost, err := resolverFunc.GasCost(&exec.VirtualMachine{})
			assert.Equal(t, model.TestGasPrice.Int64(), gasPrice)
			assert.Equal(t, GasQuickStep, cost)
			assert.NoError(t, err)
		}else {
			assert.Panics(t, func() {
				resolverFunc.Execute(&exec.VirtualMachine{})
			})
			assert.Panics(t, func() {
				resolverFunc.GasCost(&exec.VirtualMachine{})
			})
		}
	}


	// test resolve global
	assert.Equal(t, int64(0), solver.ResolveGlobal("module", "field"))
	assert.Equal(t, int64(0), solver.ResolveGlobal("env", "field"))
	assert.Equal(t, int64(0), solver.ResolveGlobal("env", "stderr"))
}
