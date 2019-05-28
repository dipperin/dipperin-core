package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/resolver"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"math/big"
)

var emptyCodeHash = cs_crypto.Keccak256Hash(nil)

var DEFAULT_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: exec.DefaultPageSize,
}

type VM struct {
	Context
	Interpreter Interpreter
	vmconfig    exec.VMConfig
	// state gives access to the underlying state
	state       StateDB
	// Depth is the current call stack
	depth int
	// abort is used to abort the VM calling operations
	// NOTE: must be set atomically
	abort int32
}

func NewVM(context Context, state StateDB, config exec.VMConfig) *VM {
	interpreter := NewWASMInterpreter(state, context, config)
	vm := VM{
		Context:     context,
		Interpreter: interpreter,
		vmconfig:    config,
		state:       state,
	}
	return &vm
}

func (vm *VM) PreCheck() error {


	return nil
}


func (vm *VM) Call(caller resolver.ContractRef, addr common.Address, input []byte,gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
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

func (vm *VM)DelegateCall(caller resolver.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error){
	return nil ,0,nil
}
func (vm *VM) Create(caller resolver.ContractRef, code []byte, abi []byte, value []byte) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = cs_crypto.CreateContractAddress(caller.Address(), vm.state.GetNonce(caller.Address()))
	return vm.create(caller, code, abi, value, contractAddr)
}

func (vm *VM) create(caller resolver.ContractRef, code []byte,abi []byte, input []byte, address common.Address) ([]byte, common.Address, uint64, error) {
	// Caller nonce ++
	vm.state.AddNonce(caller.Address(), uint64(1))

	// Ensure there's no existing contract already at the designated address
	contractHash := vm.state.GetCodeHash(address)
	if vm.state.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}

	// Create a new account on the state
	// snapshot := vm.state.Snapshot()

	vm.state.CreateAccount(address)
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

	// call Interpreter.Run()
	vm.Interpreter.Run(vm,contract, input,create)
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

	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64


	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
}

func (context *Context)GetCallGasTemp() uint64{
	return context.callGasTemp
}

func (context *Context) GetGasPrice() int64 {
	return context.GasPrice.Int64()
}

func (context *Context) BlockHash(num uint64) common.Hash {
	return context.GetHash(num)
}

func (context *Context) GetBlockNumber() *big.Int {
	return context.BlockNumber
}

func (context *Context) GetGasLimit() uint64 {
	return context.GasLimit
}

func (context *Context) GetTime() *big.Int {
	return context.Time
}

func (context *Context) GetCoinBase() common.Address {
	return context.Coinbase
}

func (context *Context) GetOrigin() common.Address {
	return context.Origin
}

// NewVMContext creates a new context for use in the VM.
func NewVMContext(tx model.AbstractTransaction, block model.AbstractBlock) Context {
	sender, _ := tx.Sender(tx.GetSigner())
	return Context{
		Origin: sender,
		GasPrice:tx.GetGasPrice(),
		GasLimit: tx.Fee().Uint64(),
		BlockNumber:new(big.Int).SetUint64(block.Number()),
		Time:block.Timestamp(),
		Coinbase:block.CoinBaseAddress(),
		Difficulty:block.Difficulty().Big(),
		callGasTemp:tx.Fee().Uint64(),
	}
}

type Caller struct {
	addr common.Address
}

func (c *Caller) Address() common.Address {
	return c.addr
}

func AccountRef(addr common.Address) resolver.ContractRef{
	return &Caller{addr}
}

func NewContract(caller resolver.ContractRef, object resolver.ContractRef, code []byte, abi []byte) *Contract{
	return &Contract{
		CallerAddress: caller.Address(),
		caller:        caller,
		self:          object,
		ABI:           abi,
		Code:          code,
	}
}



