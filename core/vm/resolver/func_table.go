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

var (
	cfc  = newCfcSet()
	cgbl = newGlobalSet()
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

func newCfcSet() map[string]map[string]*exec.FunctionImport {
	return map[string]map[string]*exec.FunctionImport{}
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
	return sysFunc[module][field]
}

func (r *Resolver) ResolveGlobal(module, field string) int64 {
	return 0
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
			"printsf":    &exec.FunctionImport{Execute: envPrintsf, GasCost: envPrintsfGasCost},
			"printdf":    &exec.FunctionImport{Execute: envPrintdf, GasCost: envPrintdfGasCost},
			"printqf":    &exec.FunctionImport{Execute: envPrintqf, GasCost: envPrintqfGasCost},
			"printn":     &exec.FunctionImport{Execute: envPrintn, GasCost: envPrintnGasCost},
			"printhex":   &exec.FunctionImport{Execute: envPrinthex, GasCost: envPrinthexGasCost},

			"abort": &exec.FunctionImport{Execute: envAbort, GasCost: envAbortGasCost},

			// compiler builtins
			// arithmetic long double
			"__ashlti3": &exec.FunctionImport{Execute: env__ashlti3, GasCost: env__ashlti3GasCost},
			"__ashrti3": &exec.FunctionImport{Execute: env__ashrti3, GasCost: env__ashrti3GasCost},
			"__lshlti3": &exec.FunctionImport{Execute: env__lshlti3, GasCost: env__lshlti3GasCost},
			"__lshrti3": &exec.FunctionImport{Execute: env__lshrti3, GasCost: env__lshrti3GasCost},
			"__divti3":  &exec.FunctionImport{Execute: env__divti3, GasCost: env__divti3GasCost},
			"__udivti3": &exec.FunctionImport{Execute: env__udivti3, GasCost: env__udivti3GasCost},
			"__modti3":  &exec.FunctionImport{Execute: env__modti3, GasCost: env__modti3GasCost},
			"__umodti3": &exec.FunctionImport{Execute: env__umodti3, GasCost: env__umodti3GasCost},
			"__multi3":  &exec.FunctionImport{Execute: env__multi3, GasCost: env__multi3GasCost},
			/*			"__addtf3":  &exec.FunctionImport{Execute: env__addtf3, GasCost: env__addtf3GasCost},
						"__subtf3":  &exec.FunctionImport{Execute: env__subtf3, GasCost: env__subtf3GasCost},
						"__multf3":  &exec.FunctionImport{Execute: env__multf3, GasCost: env__multf3GasCost},
						"__divtf3":  &exec.FunctionImport{Execute: env__divtf3, GasCost: env__divtf3GasCost},

						// conversion long double
						"__floatsitf":   &exec.FunctionImport{Execute: env__floatsitf, GasCost: env__floatsitfGasCost},
						"__floatunsitf": &exec.FunctionImport{Execute: env__floatunsitf, GasCost: env__floatunsitfGasCost},
						"__floatditf":   &exec.FunctionImport{Execute: env__floatditf, GasCost: env__floatditfGasCost},
						"__floatunditf": &exec.FunctionImport{Execute: env__floatunditf, GasCost: env__floatunditfGasCost},
						"__floattidf":   &exec.FunctionImport{Execute: env__floattidf, GasCost: env__floattidfGasCost},
						"__floatuntidf": &exec.FunctionImport{Execute: env__floatuntidf, GasCost: env__floatuntidfGasCost},
						"__floatsidf":   &exec.FunctionImport{Execute: env__floatsidf, GasCost: env__floatsidfGasCost},
						"__extendsftf2": &exec.FunctionImport{Execute: env__extendsftf2, GasCost: env__extendsftf2GasCost},
						"__extenddftf2": &exec.FunctionImport{Execute: env__extenddftf2, GasCost: env__extenddftf2GasCost},
						"__fixtfti":     &exec.FunctionImport{Execute: env__fixtfti, GasCost: env__fixtftiGasCost},
						"__fixtfdi":     &exec.FunctionImport{Execute: env__fixtfdi, GasCost: env__fixtfdiGasCost},
						"__fixtfsi":     &exec.FunctionImport{Execute: env__fixtfsi, GasCost: env__fixtfsiGasCost},
						"__fixunstfti":  &exec.FunctionImport{Execute: env__fixunstfti, GasCost: env__fixunstftiGasCost},
						"__fixunstfdi":  &exec.FunctionImport{Execute: env__fixunstfdi, GasCost: env__fixunstfdiGasCost},
						"__fixunstfsi":  &exec.FunctionImport{Execute: env__fixunstfsi, GasCost: env__fixunstfsiGasCost},
						"__fixsfti":     &exec.FunctionImport{Execute: env__fixsfti, GasCost: env__fixsftiGasCost},
						"__fixdfti":     &exec.FunctionImport{Execute: env__fixdfti, GasCost: env__fixdftiGasCost},
						"__trunctfdf2":  &exec.FunctionImport{Execute: env__trunctfdf2, GasCost: env__trunctfdf2GasCost},
						"__trunctfsf2":  &exec.FunctionImport{Execute: env__trunctfsf2, GasCost: env__trunctfsf2GasCost},

						"__eqtf2":    &exec.FunctionImport{Execute: env__eqtf2, GasCost: env__eqtf2GasCost},
						"__netf2":    &exec.FunctionImport{Execute: env__netf2, GasCost: env__netf2GasCost},
						"__getf2":    &exec.FunctionImport{Execute: env__getf2, GasCost: env__getf2GasCost},
						"__gttf2":    &exec.FunctionImport{Execute: env__gttf2, GasCost: env__gttf2GasCost},
						"__lttf2":    &exec.FunctionImport{Execute: env__lttf2, GasCost: env__lttf2GasCost},
						"__letf2":    &exec.FunctionImport{Execute: env__letf2, GasCost: env__letf2GasCost},
						"__cmptf2":   &exec.FunctionImport{Execute: env__cmptf2, GasCost: env__cmptf2GasCost},
						"__unordtf2": &exec.FunctionImport{Execute: env__unordtf2, GasCost: env__unordtf2GasCost},
						"__negtf2":   &exec.FunctionImport{Execute: env__negtf2, GasCost: env__negtf2GasCost},*/

			// for blockchain function
			"value":        &exec.FunctionImport{Execute: r.envValue, GasCost: constGasFunc(GasQuickStep)},
			"gasPrice":     &exec.FunctionImport{Execute: r.envGasPrice, GasCost: constGasFunc(GasQuickStep)},
			"blockHash":    &exec.FunctionImport{Execute: r.envBlockHash, GasCost: constGasFunc(GasQuickStep)},
			"number":       &exec.FunctionImport{Execute: r.envNumber, GasCost: constGasFunc(GasQuickStep)},
			"gasLimit":     &exec.FunctionImport{Execute: r.envGasLimit, GasCost: constGasFunc(GasQuickStep)},
			"timestamp":    &exec.FunctionImport{Execute: r.envTimestamp, GasCost: constGasFunc(GasQuickStep)},
			"coinbase":     &exec.FunctionImport{Execute: r.envCoinbase, GasCost: constGasFunc(GasQuickStep)},
			"balance":      &exec.FunctionImport{Execute: r.envBalance, GasCost: constGasFunc(GasQuickStep)},
			"origin":       &exec.FunctionImport{Execute: r.envOrigin, GasCost: constGasFunc(GasQuickStep)},
			"caller":       &exec.FunctionImport{Execute: r.envCaller, GasCost: constGasFunc(GasQuickStep)},
			"callValue":    &exec.FunctionImport{Execute: r.envCallValue, GasCost: constGasFunc(GasQuickStep)},
			"address":      &exec.FunctionImport{Execute: r.envAddress, GasCost: constGasFunc(GasQuickStep)},
			"sha3":         &exec.FunctionImport{Execute: envSha3, GasCost: envSha3GasCost},
			"emitEvent":    &exec.FunctionImport{Execute: r.envEmitEvent, GasCost: envEmitEventGasCost},
			"setState":     &exec.FunctionImport{Execute: r.envSetState, GasCost: envSetStateGasCost},
			"getState":     &exec.FunctionImport{Execute: r.envGetState, GasCost: envGetStateGasCost},
			"getStateSize": &exec.FunctionImport{Execute: r.envGetStateSize, GasCost: envGetStateSizeGasCost},

			// supplement
			"getCallerNonce": &exec.FunctionImport{Execute: r.envGetCallerNonce, GasCost: constGasFunc(GasQuickStep)},
			// "currentTime": &exec.FunctionImport{Execute: r.envCurrentTime, GasCost: constGasFunc(GasQuickStep)},
			/*			"callTransfer":   &exec.FunctionImport{Execute: r.envCallTransfer, GasCost: constGasFunc(GasQuickStep)},

						"dipcCall":               &exec.FunctionImport{Execute: r.envDipperCall, GasCost: envDipperCallGasCost},
						"dipcCallInt64":          &exec.FunctionImport{Execute: r.envDipperCallInt64, GasCost: envDipperCallInt64GasCost},
						"dipcCallString":         &exec.FunctionImport{Execute: r.envDipperCallString, GasCost: envDipperCallStringGasCost},
						"dipcDelegateCall":       &exec.FunctionImport{Execute: r.envDipperDelegateCall, GasCost: envDipperCallStringGasCost},
						"dipcDelegateCallInt64":  &exec.FunctionImport{Execute: r.envDipperDelegateCallInt64, GasCost: envDipperCallStringGasCost},
						"dipcDelegateCallString": &exec.FunctionImport{Execute: r.envDipperDelegateCallString, GasCost: envDipperCallStringGasCost},*/
		},
	}
}
