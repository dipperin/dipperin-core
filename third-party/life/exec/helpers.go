package exec

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/perlin-network/life/compiler"
	"github.com/perlin-network/life/utils"
)

var _ ImportResolver = (*NopResolver)(nil)

// NopResolver is a nil WebAssembly module import resolver.
type NopResolver struct{}

func (r *NopResolver) ResolveFunc(module, field string) *FunctionImport {
	panic("implement me")
}

func (r *NopResolver) ResolveGlobal(module, field string) int64 {
	panic("global import not allowed")
}

// RunWithGasLimit runs a WebAssembly modules function denoted by its ID with a specified set
// of parameters for a specified amount of instructions (also known as gas) denoted by `limit`.
// Panics on logical errors.
func (vm *VirtualMachine) RunWithGasLimit(entryID, limit int, params ...int64) (int64, error) {
	count := 0

	vm.Ignite(entryID, params...)
	for !vm.Exited {
		vm.Execute()
		if vm.Delegate != nil {
			vm.Delegate()
			vm.Delegate = nil
		}
		count++
		if count == limit {
			return -1, errors.New("gas limit exceeded")
		}
	}

	if vm.ExitError != nil {
		return -1, utils.UnifyError(vm.ExitError)
	}
	return vm.ReturnValue, nil
}

// Run runs a WebAssembly modules function denoted by its ID with a specified set
// of parameters.
// Panics on logical errors.
func (vm *VirtualMachine) Run(entryID int, params ...int64) (retVal int64, retErr error) {
	vm.Ignite(entryID, params...) // call Ignite() to perform necessary checks even if we are using the AOT mode.
	// vmcommon.AOTService is nil
	if vm.AOTService != nil {
		recoveryFunc := func() {
			if err := recover(); err != nil {
				if err, ok := err.(error); ok {
					retErr = err
				} else {
					panic(err)
				}
			} else {
				vm.CurrentFrame = -1
			}
		}
		targetName := fmt.Sprintf("%s%d", compiler.NGEN_FUNCTION_PREFIX, entryID)
		switch len(params) {
		case 0:
			defer recoveryFunc()
			return int64(vm.AOTService.UnsafeInvokeFunction_0(vm, targetName)), nil
		case 1:
			defer recoveryFunc()
			return int64(vm.AOTService.UnsafeInvokeFunction_1(vm, targetName, uint64(params[0]))), nil
		case 2:
			defer recoveryFunc()
			return int64(vm.AOTService.UnsafeInvokeFunction_2(vm, targetName, uint64(params[0]), uint64(params[1]))), nil
		default:
		}
	}
	//vmcommon.Exited = false when vmcommon initializes the call frame in vmcommon.Ignite
	for !vm.Exited {
		vm.Execute()
		if vm.Delegate != nil {
			vm.Delegate()
			vm.Delegate = nil
		}
	}

	if vm.ExitError != nil {
		log.Debug("VirtualMachine#Run ", "err", vm.ExitError)
		return -1, utils.UnifyError(vm.ExitError)
	}
	return vm.ReturnValue, nil
}

func (vm *VirtualMachine) Stop() (err error) {
	vm.Memory.Put(vm.Memory.Memory)
	for _, pos := range vm.ExternalParams {
		err = vm.Memory.Free(int(pos))

	}
	vm.Memory.PutTree(vm.Memory.Tree)
	return err
}
