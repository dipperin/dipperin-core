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
	"github.com/dipperin/dipperin-core/core/vm/common"
	"github.com/dipperin/dipperin-core/tests/util"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getBaseVirtualMachine(t *testing.T)(*gomock.Controller, *exec.ImportResolver, *exec.VirtualMachine)  {
	ctrl := gomock.NewController(t)

	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contract, state)

	code, _ := test_util.GetTestData("demo")
	vm, err := exec.NewVirtualMachine(code, common.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)
	return ctrl, &solver,vm
}


func Test_Instructions(t *testing.T) {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	entryID, ok := vm.GetFunctionExport("can_payable")
	assert.Equal(t, true, ok)
	vm.GasLimit = model.TestGasLimit * 100

	res, err := vm.Run(entryID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res)
}


func Test_envEmitEvent(t *testing.T)  {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	topic := int64(0)
	topicLen:= int64(10)
	dataSrc := int64(20)
	dataLen := int64(20)

	//data := []byte("hello")
	//lenPos := int64(len(data))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals:[]int64{topic, topicLen, dataSrc, dataLen},
	}
	vm.CallStack[0] = frame

	//dest:= solver.(Resolver).envEmitEvent(vm)
	//assert.NoError(t, err)
	//assert.Equal(t, int64(131072), dest)
}


func Test_envMalloc(t *testing.T)  {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	data := []byte("hello")
	lenPos := int64(len(data))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals:[]int64{lenPos},
	}
	vm.CallStack[0] = frame

	dest := envMalloc(vm)
	assert.Equal(t, int64(131072), dest)
}



func Test_envMemcpy(t *testing.T)  {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	destPos := int64(10)
	srcPos  := int64(0)
	data := []byte("hello")
	lenPos := int64(len(data))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals:[]int64{destPos,srcPos,lenPos},
	}
	vm.CallStack[0] = frame
	vm.Memory.Memory = append(data, vm.Memory.Memory[:lenPos]...)

	dest := envMemcpy(vm)
	assert.Equal(t, data, vm.Memory.Memory[dest:dest+int64(len(data))])
}

func Test_envMemmove(t *testing.T)  {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	destPos := int64(10)
	srcPos  := int64(0)
	data := []byte("hello")
	lenPos := int64(len(data))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals:[]int64{destPos,srcPos,lenPos},
	}
	vm.CallStack[0] = frame
	vm.Memory.Memory = append(data, vm.Memory.Memory[:lenPos]...)

	dest := envMemmove(vm)
	assert.Equal(t, data, vm.Memory.Memory[dest:dest+int64(len(data))])
}

func Test_MallocString(t *testing.T) {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()

	testCases := []struct{
		name string
		given string
		expect string
	}{
		{
			name:"normalString",
			given:"DIPP",
			expect:"DIPP",
		},
		{
			name:"emptyString",
			given:"",
			expect:"",
		},
		{
			name:"addrString",
			given:model.AliceAddr.String(),
			expect:model.AliceAddr.String(),
		},
	}

	for _,tc := range testCases{
		pos := MallocString(vm, tc.given)
		assert.Equal(t, tc.given, string(vm.Memory.Memory[pos:pos+int64(len(tc.expect))]) )
	}
}




