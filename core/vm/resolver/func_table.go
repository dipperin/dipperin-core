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
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
)

// Gas costs
const (
	GasQuickStep   uint64 = 2
	GasFastestSetp uint64 = 3

	// ...
)

type Resolver struct {
	Service resolverNeedExternalService
}

func NewResolver(vmValue VmContextService, contract ContractService, state StateDBService) exec.ImportResolver {
	return &Resolver{
		Service: resolverNeedExternalService{
			ContractService:  contract,
			VmContextService: vmValue,
			StateDBService:   state,
		},
	}
}

func (r *Resolver) ResolveFunc(module, field string) *exec.FunctionImport {
	sysFunc := newSystemFuncSet(r)
	df := &exec.FunctionImport{
		Execute: func(vm *exec.VirtualMachine) int64 {
			panic(fmt.Sprintf("unsupport func module:%s field:%s", module, field))
		},
		GasCost: func(vm *exec.VirtualMachine) (uint64, error) {
			panic(fmt.Sprintf("unsupport gas cost module:%s field:%s", module, field))
		},
	}

	if m, exist := sysFunc[module]; exist == true {
		if f, exist := m[field]; exist == true {
			return f
		} else {
			return df
		}
	} else {
		return df
	}
}

func (r *Resolver) ResolveGlobal(module, field string) int64 {
	if m, exist := newGlobalSet()[module]; exist == true {
		if g, exist := m[field]; exist == true {
			return g
		} else {
			return 0
			//panic("unknown field " + field)

		}
	} else {
		return 0
		//panic("unknown module " + module)
	}
}

func newGlobalSet() map[string]map[string]int64 {
	return map[string]map[string]int64{
		"env": {
			"stderr": 0,
			"stdin":  0,
			"stdout": 0,
		},
	}
}

func newSystemFuncSet(r *Resolver) map[string]map[string]*exec.FunctionImport {
	return map[string]map[string]*exec.FunctionImport{
		"env": {
			"malloc":  &exec.FunctionImport{Execute: envMalloc, GasCost: envMallocGasCost},
			"free":    &exec.FunctionImport{Execute: envFree, GasCost: envFreeGasCost},
			"calloc":  &exec.FunctionImport{Execute: envCalloc, GasCost: envCallocGasCost},
			"realloc": &exec.FunctionImport{Execute: envRealloc, GasCost: envReallocGasCost},

			"memcpy":  &exec.FunctionImport{Execute: envMemcpy, GasCost: envMemcpyGasCost},
			"memmove": &exec.FunctionImport{Execute: envMemmove, GasCost: envMemmoveGasCost},
			"memcmp":  &exec.FunctionImport{Execute: envMemcpy, GasCost: envMemmoveGasCost},
			"memset":  &exec.FunctionImport{Execute: envMemset, GasCost: envMemsetGasCost},

			"prints":     &exec.FunctionImport{Execute: r.envPrints, GasCost: envPrintsGasCost},
			"prints_l":   &exec.FunctionImport{Execute: envPrintsl, GasCost: envPrintslGasCost},
			"printi":     &exec.FunctionImport{Execute: envPrinti, GasCost: envPrintiGasCost},
			"printui":    &exec.FunctionImport{Execute: envPrintui, GasCost: envPrintuiGasCost},
			"printi128":  &exec.FunctionImport{Execute: envPrinti128, GasCost: envPrinti128GasCost},
			"printui128": &exec.FunctionImport{Execute: envPrintui128, GasCost: envPrintui128GasCost},
			/*			"printsf":    &exec.FunctionImport{Execute: envPrintsf, GasCost: envPrintsfGasCost},
						"printdf":    &exec.FunctionImport{Execute: envPrintdf, GasCost: envPrintdfGasCost},
						"printqf":    &exec.FunctionImport{Execute: envPrintqf, GasCost: envPrintqfGasCost},*/
			"printn":   &exec.FunctionImport{Execute: envPrintn, GasCost: envPrintnGasCost},
			"printhex": &exec.FunctionImport{Execute: envPrinthex, GasCost: envPrinthexGasCost},

			"abort": &exec.FunctionImport{Execute: envAbort, GasCost: envAbortGasCost},
			// for blockchain function
			"gasPrice":            &exec.FunctionImport{Execute: r.envGasPrice, GasCost: constGasFunc(GasQuickStep)},
			"blockHash":           &exec.FunctionImport{Execute: r.envBlockHash, GasCost: constGasFunc(GasQuickStep)},
			"number":              &exec.FunctionImport{Execute: r.envNumber, GasCost: constGasFunc(GasQuickStep)},
			"gasLimit":            &exec.FunctionImport{Execute: r.envGasLimit, GasCost: constGasFunc(GasQuickStep)},
			"timestamp":           &exec.FunctionImport{Execute: r.envTimestamp, GasCost: constGasFunc(GasQuickStep)},
			"coinbase":            &exec.FunctionImport{Execute: r.envCoinbase, GasCost: constGasFunc(GasQuickStep)},
			"balance":             &exec.FunctionImport{Execute: r.envBalance, GasCost: constGasFunc(GasQuickStep)},
			"origin":              &exec.FunctionImport{Execute: r.envOrigin, GasCost: constGasFunc(GasQuickStep)},
			"caller":              &exec.FunctionImport{Execute: r.envCaller, GasCost: constGasFunc(GasQuickStep)},
			"callValue":           &exec.FunctionImport{Execute: r.envCallValue, GasCost: constGasFunc(GasQuickStep)},
			"address":             &exec.FunctionImport{Execute: r.envAddress, GasCost: constGasFunc(GasQuickStep)},
			"sha3":                &exec.FunctionImport{Execute: envSha3, GasCost: envSha3GasCost},
			"hexStringSameWithVM": &exec.FunctionImport{Execute: envHexStringSameWithVM, GasCost: envHexStringSameWithVMGasCost},
			"emitEvent":           &exec.FunctionImport{Execute: r.envEmitEvent, GasCost: envEmitEventGasCost},
			"setState":            &exec.FunctionImport{Execute: r.envSetState, GasCost: envSetStateGasCost},
			"getState":            &exec.FunctionImport{Execute: r.envGetState, GasCost: envGetStateGasCost},
			"getStateSize":        &exec.FunctionImport{Execute: r.envGetStateSize, GasCost: envGetStateSizeGasCost},

			// supplement
			"getCallerNonce": &exec.FunctionImport{Execute: r.envGetCallerNonce, GasCost: constGasFunc(GasQuickStep)},
			// "currentTime": &exec.FunctionImport{Execute: r.envCurrentTime, GasCost: constGasFunc(GasQuickStep)},
			"callTransfer": &exec.FunctionImport{Execute: r.envCallTransfer, GasCost: constGasFunc(GasQuickStep)},

			"dipcCall":               &exec.FunctionImport{Execute: r.envDipperCall, GasCost: envDipperCallGasCost},
			"dipcCallInt64":          &exec.FunctionImport{Execute: r.envDipperCallInt64, GasCost: envDipperCallInt64GasCost},
			"dipcCallString":         &exec.FunctionImport{Execute: r.envDipperCallString, GasCost: envDipperCallStringGasCost},
			"dipcDelegateCall":       &exec.FunctionImport{Execute: r.envDipperDelegateCall, GasCost: envDipperCallStringGasCost},
			"dipcDelegateCallInt64":  &exec.FunctionImport{Execute: r.envDipperDelegateCallInt64, GasCost: envDipperCallStringGasCost},
			"dipcDelegateCallString": &exec.FunctionImport{Execute: r.envDipperDelegateCallString, GasCost: envDipperCallStringGasCost},
		},
	}
}
