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

package vm

import (
	"github.com/dipperin/dipperin-core/tests/factory/g-testData"
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
