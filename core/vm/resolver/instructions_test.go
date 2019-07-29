package resolver

import (
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"testing"
)

var TEST_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: exec.DefaultPageSize,
}

func TestInstructions(t *testing.T) {
	/*	ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		vmValue := NewMockVmContextService(ctrl)
		contract := NewMockContractService(ctrl)
		state := NewMockStateDBService(ctrl)
		solver := NewResolver(vmValue, contract, state)
		vm, err := exec.NewVirtualMachine([]byte{}, TEST_VM_CONFIG, solver, nil)
		assert.NoError(t, err)

		solverFunc := solver.ResolveFunc("env", "setState")
		gasCost, err := solverFunc.GasCost(vm)
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), gasCost)
		assert.Equal(t, int64(0), solverFunc.Execute(vm))*/
}
