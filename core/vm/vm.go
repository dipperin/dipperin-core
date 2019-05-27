package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	cs_crypto "github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/core/vm/resolver"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
)

var emptyCodeHash = cs_crypto.Keccak256Hash(nil)
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

	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

func NewVM(context Context, state StateDB, config exec.VMConfig) *VM {
	interpreter := NewWASMInterpreter(state, context, config)
	vm := VM{
		Context:context,
		interpreter:interpreter,
		vmconfig:DEFAULT_VM_CONFIG,
		resolver:&resolver.Resolver{},
		state:state,
	}
	return &vm
}

func (vm *VM)GetCallGasTemp() uint64{
	return vm.callGasTemp
}

func (vm *VM) GasPrice() int64 {
	return vm.Context.GasPrice.Int64()
}

func (vm *VM) BlockHash(num uint64) common.Hash {
	return vm.Context.GetHash(num)
}

func (vm *VM) BlockNumber() *big.Int {
	return vm.Context.BlockNumber
}

func (vm *VM) GasLimit() uint64 {
	return vm.Context.GasLimit
}

func (vm *VM) Time() *big.Int {
	return vm.Context.Time
}

func (vm *VM) CoinBase() common.Address {
	return vm.Context.Coinbase
}

func (vm *VM) Origin() common.Address {
	return vm.Context.Origin
}

func (vm *VM) Call(caller ContractRef, addr common.Address, input []byte,gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	code := vm.state.GetState(addr,[]byte("code"))
	abi := vm.state.GetState(addr,[]byte("abi"))
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

func (vm *VM)DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error){
	return nil ,0,nil
}
func (vm *VM) Create(caller ContractRef, code []byte, abi []byte, value []byte) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = cs_crypto.CreateContractAddress(caller.Address(), vm.state.GetNonce(caller.Address()))
	return vm.create(caller, code, abi, value, contractAddr)
}

func (vm *VM) create(caller ContractRef, code []byte,abi []byte, input []byte, address common.Address) ([]byte, common.Address, uint64, error) {
	// Caller nonce ++
	vm.state.AddNonce(caller.Address(), uint64(1))

	// Ensure there's no existing contract already at the designated address
	contractHash := vm.state.GetCodeHash(address)
	if vm.state.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}

	// Create a new account on the state
	// snapshot := vm.state.Snapshot()

	// vm.state.CreateAccount(address)
	// vm.Transfer(evm.StateDB, caller.Address(), address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, AccountRef(address), code, abi)

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

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, common.Address, common.Address, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

type Context struct {
	// Message information
	Origin common.Address // Provides information for ORIGIN

	GetHash GetHashFunc

	// Block information
	Coinbase common.Address // Provides information for COINBASE

	GasPrice *big.Int       // Provides information for GASPRICE
	GasLimit    uint64         // Provides information for GASLIMIT
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

func AccountRef(addr common.Address) ContractRef{
	return &Caller{addr}
}

func NewContract(caller ContractRef, object ContractRef, code []byte, abi []byte) *Contract{
	return &Contract{
		CallerAddress: caller.Address(),
		caller:        caller,
		self:          object,
		ABI:           abi,
		Code:          code,
	}
}



