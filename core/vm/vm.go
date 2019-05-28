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
	state StateDB
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

func (vm *VM) GetStateDB() StateDB {
	return vm.state
}

func (vm *VM) PreCheck() error {

	return nil
}

func (vm *VM) Call(caller resolver.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	code := vm.state.GetState(addr, []byte("code"))
	abi := vm.state.GetState(addr, []byte("abi"))

	contract := &Contract{
		CallerAddress: caller.Address(),
		caller:        caller,
		self:          &Caller{addr: addr},
		ABI:           abi,
		Code:          code,
		value:         value,
		Gas:           gas,
	}

	ret, err = run(vm, contract, input, false)
	return
}

func (vm *VM) DelegateCall(caller resolver.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	return nil, 0, nil
}

func (vm *VM) Create(caller resolver.ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = cs_crypto.CreateContractAddress(caller.Address(), vm.state.GetNonce(caller.Address()))
	return
}

/*func (vm *VM) create(caller resolver.ContractRef, code []byte, gas uint64, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Caller nonce ++
	vm.state.AddNonce(caller.Address(), uint64(1))

	// Ensure there's no existing contract already at the designated address
	contractHash := vm.state.GetCodeHash(address)
	if vm.state.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}

	// Create a new account on the state
	// snapshot := vmcommon.state.Snapshot()

	vm.state.CreateAccount(address)
	// vmcommon.Transfer(evm.StateDB, caller.Address(), address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, AccountRef(address), code, nil, gas)
	vm.state.SetState(contract.self.Address(), []byte("code"), code)
	vm.state.SetState(contract.self.Address(), []byte("abi"), abi)
	// call run
	run(vm, contract, input, true)

	return nil, address, uint64(0), nil

	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if evm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gas, ErrDepth
	}
	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, common.Address{}, gas, ErrInsufficientBalance
	}
	nonce := evm.StateDB.GetNonce(caller.Address())
	evm.StateDB.SetNonce(caller.Address(), nonce+1)

	// Ensure there's no existing contract already at the designated address
	contractHash := evm.StateDB.GetCodeHash(address)
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

	// initialise a new contract and set the code that is to be used by the
	// EVM. The contract is a scoped environment for this execution context
	// only.
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCallCode(&address, crypto.Keccak256Hash(code), code)

	if evm.vmConfig.NoRecursion && evm.depth > 0 {
		return nil, address, gas, nil
	}

	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), address, true, code, gas, value)
	}
	start := time.Now()

	ret, err := run(evm, contract, nil, false)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := evm.ChainConfig().IsEIP158(evm.BlockNumber) && len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			evm.StateDB.SetCode(address, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (evm.ChainConfig().IsHomestead(evm.BlockNumber) || err != ErrCodeStoreOutOfGas)) {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}
	return ret, address, contract.Gas, err
}
*/
func run(vm *VM, contract *Contract, input []byte, create bool) ([]byte, error) {
	// call Interpreter.Run()
	return vm.Interpreter.Run(vm, contract, input, create)
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

	GasPrice    *big.Int // Provides information for GASPRICE
	GasLimit    uint64   // Provides information for GASLIMIT
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

func (context *Context) GetCallGasTemp() uint64 {
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
		Origin:      sender,
		GasPrice:    tx.GetGasPrice(),
		GasLimit:    tx.Fee().Uint64(),
		BlockNumber: new(big.Int).SetUint64(block.Number()),
		Time:        block.Timestamp(),
		Coinbase:    block.CoinBaseAddress(),
		Difficulty:  block.Difficulty().Big(),
		callGasTemp: tx.Fee().Uint64(),
	}
}

type Caller struct {
	addr common.Address
}

func (c *Caller) Address() common.Address {
	return c.addr
}

func AccountRef(addr common.Address) resolver.ContractRef {
	return &Caller{addr}
}