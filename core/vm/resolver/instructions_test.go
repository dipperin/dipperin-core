package resolver

import (
	"github.com/dipperin/dipperin-core/tests/g-testData"
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
