package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	cs_crypto "github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"math/big"
	"time"
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
	interpreter := NewWASMInterpreter(state, context, config)
	vm := VM{context, interpreter, DEFAULT_VM_CONFIG, &Resolver{}, state}
	return &vm
}


func (vm *VM) Call(caller ContractRef, addr common.Address, input []byte, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	code := vm.state.GetState(addr, []byte("code"))
	abi := vm.state.GetState(addr, []byte("abi"))
	contract := &Contract{
		CallerAddress: caller.Address(),
		caller:        caller,
		self:          &Caller{addr: addr},
		ABI:           abi,
		Code:          code,
	}

	ret, err = run(vm, contract, input,false)
	return
}

func (vm *VM) Create(caller ContractRef, code []byte, abi []byte, value []byte) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = cs_crypto.CreateContractAddress(caller.Address(), vm.state.GetNonce(caller.Address()))
	return vm.create(caller, code, abi, value, contractAddr)
}

func (vm *VM) create(caller ContractRef, code []byte,abi []byte, input []byte, address common.Address) ([]byte, common.Address, uint64, error) {
	contract := &Contract{
		CallerAddress: caller.Address(),
		caller:        caller,
		self:          &Caller{addr: address},
		ABI:           abi,
		Code:          code,
	}

	// Caller nonce ++
	vm.state.AddNonce(caller.Address(), uint64(1))

	// Ensure there's no existing contract already at the designated address
	contractHash := vm.state.GetCodeHash(address)
	if evm.StateDB.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(address)
	if evm.ChainConfig().IsEIP158(evm.BlockNumber) {
		evm.StateDB.SetNonce(address, 1)
	}
	evm.Transfer(evm.StateDB, caller.Address(), address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, address, gas, nil
	}

	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), address, true, codeAndHash.code, gas, value)
	}
	start := time.Now()


	vm.state.SetState(contract.self.Address(), []byte("code"), code)
	vm.state.SetState(contract.self.Address(), []byte("abi"), abi)
	// call run
	run(vm, contract, input,true)

	return nil, address, uint64(0), nil
}

func run(vm *VM, contract *Contract, input []byte, create bool) ([]byte, error) {

	// call interpreter.Run()
	vm.interpreter.Run(contract, input,create)
	return nil, nil
}

type Context struct {
	// Message information
	Origin common.Address // Provides information for ORIGIN

	// Block information
	Coinbase common.Address // Provides information for COINBASE
	//GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int    // Provides information for NUMBER
	BlockHash   common.Hash // Provides information for Hash
	Time        *big.Int    // Provides information for TIME
	Difficulty  *big.Int    // Provides information for DIFFICULTY
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
