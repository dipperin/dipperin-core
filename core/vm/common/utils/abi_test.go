package utils

import (
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWasmAbi_FromJson(t *testing.T) {
	abiByte := new(WasmAbi)
	err := abiByte.FromJson(nil)
	assert.Equal(t, errEmptyInput, err)

	WASMPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	_, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)

	err = abiByte.FromJson(abi)
	assert.NoError(t, err)
}
