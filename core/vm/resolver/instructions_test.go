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
	common2 "github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/base"
	"github.com/dipperin/dipperin-core/tests/factory/vminfo"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func getBaseVirtualMachine(t *testing.T) (*gomock.Controller, exec.ImportResolver, *exec.VirtualMachine) {
	ctrl := gomock.NewController(t)

	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contract, state)

	code, _ := vminfo.GetTestData("demo")
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)
	return ctrl, solver, vm
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



func Test_envGetState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contract, state)

	code, _ := vminfo.GetTestData("demo")
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

	key := int64(0)
	keyLen := int64(10)
	value := int64(20)
	valueLen := int64(5)
	copyData := []byte{123}

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{key, keyLen, value, valueLen},
	}
	vm.CallStack[0] = frame

	contract.EXPECT().Address().Return(model.ContractAddr).AnyTimes()
	state.EXPECT().GetState(model.ContractAddr, vm.Memory.Memory[key:key+keyLen]).Return(copyData)

	appendData := make([]byte, int(valueLen)-len(copyData))
	result := append(copyData, appendData...)

	solver.(*Resolver).envGetState(vm)
	assert.Equal(t, result, vm.Memory.Memory[value:value+valueLen])

}

func Test_envGetStateSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contract, state)

	code, _ := vminfo.GetTestData("demo")
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

	key := int64(0)
	keyLen := int64(10)

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{key, keyLen},
	}
	vm.CallStack[0] = frame

	contract.EXPECT().Address().Return(model.ContractAddr).AnyTimes()
	state.EXPECT().GetState(model.ContractAddr, vm.Memory.Memory[key:key+keyLen]).Return([]byte{123})

	result := solver.(*Resolver).envGetStateSize(vm)
	assert.Equal(t, len([]byte{123}), int(result))

}

func prepare_envEmitEvent(t *testing.T) (*exec.VirtualMachine, exec.ImportResolver) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	state := NewMockStateDBService(ctrl)
	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	solver := NewResolver(vmValue, contract, state)

	code, _ := vminfo.GetTestData("demo")
	vm, err := exec.NewVirtualMachine(code, base.DEFAULT_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

	topic := int64(0)
	topicLen := int64(10)
	dataSrc := int64(20)
	dataLen := int64(20)

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{topic, topicLen, dataSrc, dataLen},
	}

	tData := make([]byte, topicLen)
	dData := make([]byte, dataLen)
	copy(tData, vm.Memory.Memory[topic:topic+topicLen])
	copy(dData, vm.Memory.Memory[dataSrc:dataSrc+dataLen])
	log := &model.Log{
		Address:   model.ContractAddr,
		Topics:    []common2.Hash{common2.BytesToHash(crypto.Keccak256(tData))},
		Data:      dData,
		TopicName: string(tData),
		TxHash:    common2.Hash{},
	}

	t.Log("vm info ", string(tData))
	t.Log("vm info hash ", common2.HexToHash(string(tData)), tData)

	vmValue.EXPECT().GetBlockNumber().Return(big.NewInt(0)).AnyTimes()
	contract.EXPECT().Address().Return(model.ContractAddr).AnyTimes()
	vmValue.EXPECT().GetTxHash().Return(common2.Hash{}).AnyTimes()
	state.EXPECT().AddLog(log).AnyTimes()

	vm.CallStack[0] = frame
	return vm, solver
}

func Test_envEmitEvent(t *testing.T) {
	vm, solver := prepare_envEmitEvent(t)
	solver.(*Resolver).envEmitEvent(vm)
}

func Test_envMalloc(t *testing.T) {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	data := []byte("hello")
	lenPos := int64(len(data))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{lenPos},
	}
	vm.CallStack[0] = frame

	dest := envMalloc(vm)
	assert.Equal(t, int64(131072), dest)
}

func Test_envMemcpy(t *testing.T) {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	destPos := int64(10)
	srcPos := int64(0)
	data := []byte("hello")
	lenPos := int64(len(data))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{destPos, srcPos, lenPos},
	}
	vm.CallStack[0] = frame
	vm.Memory.Memory = append(data, vm.Memory.Memory[:lenPos]...)

	dest := envMemcpy(vm)
	assert.Equal(t, data, vm.Memory.Memory[dest:dest+int64(len(data))])
}

func Test_envMemmove(t *testing.T) {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()
	destPos := int64(10)
	srcPos := int64(0)
	data := []byte("hello")
	lenPos := int64(len(data))

	t.Log("currentFrame", vm.CurrentFrame)
	vm.CurrentFrame++
	frame := exec.Frame{
		Locals: []int64{destPos, srcPos, lenPos},
	}
	vm.CallStack[0] = frame
	vm.Memory.Memory = append(data, vm.Memory.Memory[:lenPos]...)

	dest := envMemmove(vm)
	assert.Equal(t, data, vm.Memory.Memory[dest:dest+int64(len(data))])
}

func Test_MallocString(t *testing.T) {
	ctrl, _, vm := getBaseVirtualMachine(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  string
		expect string
	}{
		{
			name:   "normalString",
			given:  "DIPP",
			expect: "DIPP",
		},
		{
			name:   "emptyString",
			given:  "",
			expect: "",
		},
		{
			name:   "addrString",
			given:  model.AliceAddr.String(),
			expect: model.AliceAddr.String(),
		},
	}

	for _, tc := range testCases {
		pos := MallocString(vm, tc.given)
		assert.Equal(t, tc.given, string(vm.Memory.Memory[pos:pos+int64(len(tc.expect))]))
	}
}
