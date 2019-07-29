package resolver

import (
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewResolver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vmValue := NewMockVmContextService(ctrl)
	contract := NewMockContractService(ctrl)
	state := NewMockStateDBService(ctrl)
	solver := NewResolver(vmValue, contract, state)

	// test resolve function
	resolverFunc := solver.ResolveFunc("module", "field")
	assert.Panics(t, func() {
		resolverFunc.Execute(&exec.VirtualMachine{})
	})
	assert.Panics(t, func() {
		resolverFunc.GasCost(&exec.VirtualMachine{})
	})

	resolverFunc = solver.ResolveFunc("env", "field")
	assert.Panics(t, func() {
		resolverFunc.Execute(&exec.VirtualMachine{})
	})
	assert.Panics(t, func() {
		resolverFunc.GasCost(&exec.VirtualMachine{})
	})

	vmValue.EXPECT().GetGasPrice().Return(g_testData.TestGasPrice).AnyTimes()
	resolverFunc = solver.ResolveFunc("env", "gasPrice")
	gasPrice := resolverFunc.Execute(&exec.VirtualMachine{})
	cost, err := resolverFunc.GasCost(&exec.VirtualMachine{})
	assert.NoError(t, err)
	assert.Equal(t, g_testData.TestGasPrice.Int64(), gasPrice)
	assert.Equal(t, GasQuickStep, cost)

	// test resolve global
	assert.Equal(t, int64(0), solver.ResolveGlobal("module", "field"))
	assert.Equal(t, int64(0), solver.ResolveGlobal("env", "field"))
	assert.Equal(t, int64(0), solver.ResolveGlobal("env", "stderr"))
}
