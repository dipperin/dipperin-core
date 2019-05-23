package vm

import (
	"github.com/dipperin/dipperin-core/third-party/life/exec"
)

type Resolver struct {
	context  *Context
	contract *Contract
	state    StateDB
}

func (r *Resolver) ResolveFunc(module, field string) exec.FunctionImport {
	sysFunc := newSystemFuncSet(r)
	return sysFunc[module][field]
}

func (r *Resolver) ResolveGlobal(module, field string) int64 {
	return 0
}

func newResolver(context *Context, contract *Contract, state StateDB) exec.ImportResolver {
	return &Resolver{context, contract, state}
}

func newSystemFuncSet(r *Resolver) map[string]map[string]exec.FunctionImport {
	return map[string]map[string]exec.FunctionImport{
		"env": {
			// for blockchain function
			"malloc":       r.envMalloc,
			"free":         r.envFree,
			"setState":     r.envSetState,
			"getState":     r.envGetState,
			"getStateSize": r.envGetStateSize,

			"prints":   r.envPrints,
			"printi":   r.envPrinti,
			"printui":  r.envPrintui,
			"prints_l": r.envPrints_l,

			"coinBase":  r.envCoinBase,
			"blockNum":  r.envBlockNum,
			"difficult": r.envDifficulty,

			"abort": r.envAbort,

			"memcpy":  r.envMemcpy,
			"memmove":  r.envMemmove,
		},
	}
}
