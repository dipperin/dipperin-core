package vm

import (
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
)

var DEFAULT_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: exec.DefaultPageSize,
}

type VM struct {
	Context
	interpreter Interpreter
	vmconfig    exec.VMConfig
	resolver    exec.ImportResolver
	state       StateDB
}

func NewVM(context Context, state StateDB, config exec.VMConfig) *VM {
	interpreter := NewWASMInterpreter(state,context,config)
	vm := VM{context, interpreter, DEFAULT_VM_CONFIG, &Resolver{}, state}
	return &vm
}

func (vm *VM) Call(caller ContractRef, addr common.Address, input []byte, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	code := vm.state.GetState(addr,[]byte("code"))
	abi := vm.state.GetState(addr,[]byte("abi"))
	contract := &Contract{
		CallerAddress:caller.Address(),
		caller:caller,
		self:&Caller{addr:addr},
		ABI:abi,
		Code:code,
	}

	ret, err = run(vm, contract, input)
	return
}

func (vm *VM) Create(caller ContractRef, code []byte, abi []byte, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = common.HexToAddress("0x1223")
	return vm.create(caller, code, abi, value, contractAddr)
}

func (vm *VM) create(caller ContractRef, code []byte,abi []byte, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {


	contract := &Contract{
		CallerAddress:caller.Address(),
		caller:caller,
		self:&Caller{addr:address},
		ABI:abi,
		Code:code,
	}

	vm.state.SetState(contract.self.Address(),[]byte("code"),code)
	vm.state.SetState(contract.self.Address(),[]byte("abi"),abi)
	// call run
	run(vm, contract, nil)

	return nil, address, uint64(0), nil
}

func run(vm *VM, contract *Contract, input []byte) ([]byte, error) {

	// call interpreter.Run()
	vm.interpreter.Run(contract, input)
	return nil, nil
}

type Context struct {
	// Message information
	Origin common.Address // Provides information for ORIGIN

	// Block information
	Coinbase common.Address // Provides information for COINBASE
	//GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int // Provides information for NUMBER
	Time        *big.Int // Provides information for TIME
	Difficulty  *big.Int // Provides information for DIFFICULTY
	Log         log.Logger
}

func NewVMContext(tx model.AbstractTransaction) Context {
	sender, _ := tx.Sender(nil)
	return Context{
		Origin: sender,
	}
}

type ContractRef interface {
	Address() common.Address
}

type Caller struct {
	addr common.Address
}

func (c *Caller) Address() common.Address {
	return c.addr
}
