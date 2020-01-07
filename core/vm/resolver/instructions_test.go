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
	"github.com/dipperin/dipperin-core/tests/factory/g-testData"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstructions(t *testing.T) {
	vmValue := &fakeVmContextService{}
	contract := &fakeContractService{}
	state := NewFakeStateDBService()
	solver := NewResolver(vmValue, contract, state)

	WASMPath := g_testData.GetWASMPath("dipclib_test", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("dipclib_test", g_testData.CoreVmTestData)
	code, _ := g_testData.GetCodeAbi(WASMPath, AbiPath)
	vm, err := exec.NewVirtualMachine(code, TEST_VM_CONFIG, solver, nil)
	assert.NoError(t, err)
	entryID, ok := vm.GetFunctionExport("libTest")
	assert.Equal(t, true, ok)
	vm.GasLimit = g_testData.TestGasLimit * 100

	res, err := vm.Run(entryID)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), res)
}

func TestMallocString(t *testing.T) {
	vmValue := &fakeVmContextService{}
	contract := &fakeContractService{}
	state := NewFakeStateDBService()
	solver := NewResolver(vmValue, contract, state)

	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	code, _ := g_testData.GetCodeAbi(WASMPath, AbiPath)
	vm, err := exec.NewVirtualMachine(code, TEST_VM_CONFIG, solver, nil)
	assert.NoError(t, err)

	param := "DIPP"
	pos := MallocString(vm, param)
	assert.Equal(t, param, string(vm.Memory.Memory[pos:pos+int64(len(param))]))

	param = ""
	pos = MallocString(vm, param)
	assert.Equal(t, param, string(vm.Memory.Memory[pos:pos+int64(len(param))]))

	param = aliceAddr.String()
	pos = MallocString(vm, param)
	assert.Equal(t, param, string(vm.Memory.Memory[pos:pos+int64(len(param))]))
}
