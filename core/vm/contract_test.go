package vm

import (
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewContract(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)

	inputs := genInput(t, "hello", nil)
	contract := getContract(WASMPath, AbiPath, inputs)

	value := g_testData.TestValue
	gasLimit := g_testData.TestGasLimit
	assert.Equal(t, aliceAddr, contract.Caller().Address())
	assert.Equal(t, value, contract.CallValue())
	assert.Equal(t, gasLimit, contract.GetGas())
	assert.Equal(t, false, contract.UseGas(uint64(gasLimit*2)))
	assert.Equal(t, true, contract.UseGas(uint64(gasLimit/2)))
	assert.Equal(t, gasLimit/2, contract.GetGas())
}

func TestContract_AsDelegate(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)

	inputs := genInput(t, "hello", nil)
	callerContract := getContract(WASMPath, AbiPath, inputs)
	contract := &Contract{caller: callerContract}

	deContract := contract.AsDelegate()
	assert.Equal(t, true, deContract.DelegateCall)
	assert.Equal(t, deContract.Caller().Address(), callerContract.Address())
	assert.Equal(t, deContract.CallValue(), callerContract.CallValue())
}
