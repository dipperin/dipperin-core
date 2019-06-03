package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/resolver"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"math/big"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
)

var emptyCodeHash = cs_crypto.Keccak256Hash(nil)

var DEFAULT_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: exec.DefaultPageSize,
}

type VM struct {
	Context
	Interpreter Interpreter
	vmConfig    exec.VMConfig
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
		vmConfig:    config,
		state:       state,
	}
	return &vm
}

func (vm *VM) GetStateDB() StateDB {
	return vm.state
}

func (vm *VM) Call(caller resolver.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if vm.vmConfig.NoRecursion && vm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if vm.depth > int(model2.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if !vm.Context.CanTransfer(vm.state, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}

	var (
		to       = AccountRef(addr)
		snapshot = vm.state.Snapshot() // - snapshot.
	)

	if !vm.state.Exist(addr) {
		/*precompiles := PrecompiledContractsHomestead
		if evm.ChainConfig().IsByzantium(evm.BlockNumber) {
			precompiles = PrecompiledContractsByzantium
		}
		if precompiles[addr] == nil && PrecompiledContractsPpos[addr] == nil && evm.ChainConfig().IsEIP158(evm.BlockNumber) && value.Sign() == 0 {
			// Calling a non existing account, don't do anything, but ping the tracer
			if evm.vmConfig.Debug && evm.depth == 0 {
				evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)
				evm.vmConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, gas, nil
		}*/
		vm.state.CreateAccount(addr)
	}
	vm.Transfer(vm.state, caller.Address(), to.Address(), value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, vm.state.GetCodeHash(addr), vm.state.GetCode(addr))

	//start := time.Now()

	// Capture the tracer start/end events in debug mode
	/*	if evm.vmConfig.Debug && evm.depth == 0 {
			evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)

			defer func() { // Lazy evaluation of the parameters
				evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
			}()
		}*/
	ret, err = run(vm, contract, input, false)
	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		vm.state.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
			log.Info("callContract Use", "gasUsed", contract.Gas, "gasLeft", contract.Gas)
		}
	}
	log.Info("lifeVm run successful", "gasLeft", contract.Gas)
	return ret, contract.Gas, err
}

func (vm *VM) DelegateCall(caller resolver.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if vm.vmConfig.NoRecursion && vm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if vm.depth > int(model2.CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	var (
		snapshot = vm.GetStateDB().Snapshot()
		to       = AccountRef(caller.Address())
	)

	// Initialise a new contract and make initialise the delegate values
	contract := NewContract(caller, to, nil, gas).AsDelegate()
	contract.SetCallCode(&addr, vm.GetStateDB().GetCodeHash(addr), vm.GetStateDB().GetCode(addr))

	ret, err = run(vm, contract, input, false)
	if err != nil {
		vm.GetStateDB().RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

func (vm *VM) Create(caller resolver.ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = cs_crypto.CreateContractAddress(caller.Address(), vm.state.GetNonce(caller.Address()))
	return vm.create(caller, code, gas, value, contractAddr)
}

func (vm *VM) create(caller resolver.ContractRef, code []byte, gas uint64, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if vm.depth > int(model2.CallCreateDepth) {
		return nil, common.Address{}, gas, ErrDepth
	}

	if !vm.CanTransfer(vm.state, caller.Address(), value) {
		return nil, common.Address{}, gas, ErrInsufficientBalance
	}
	vm.state.AddNonce(caller.Address(), uint64(1))

	// Ensure there's no existing contract already at the designated address
	contractHash := vm.state.GetCodeHash(address)
	if vm.state.GetNonce(address) != 0 || (contractHash != common.Hash{} && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}

	// Create a new account on the state
	snapshot := vm.state.Snapshot()
	vm.state.CreateAccount(address)
	vm.Transfer(vm.state, caller.Address(), address, value)

	// initialise a new contract and set the code that is to be used by the
	// EVM. The contract is a scoped environment for this execution context
	// only.
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCallCode(&address, cs_crypto.Keccak256Hash(code), code)

	if vm.vmConfig.NoRecursion && vm.depth > 0 {
		return nil, address, gas, nil
	}

	/*	if vm.vmConfig.Debug && vm.depth == 0 {
			vm.vmConfig.Tracer.CaptureStart(caller.Address(), address, true, code, gas, value)
		}*/
	//start := time.Now()

	ret, err := run(vm, contract, nil, true)
	log.Info("lifeVm run successful", "gasLeft", contract.Gas)
	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > model2.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * model2.CreateDataGas
		if contract.UseGas(createDataGas) {
			vm.state.SetCode(address, ret)
			log.Info("createDataGas Use", "CodeLen", len(ret), "gasUsed", createDataGas, "gasLeft", contract.Gas)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && err != ErrCodeStoreOutOfGas) {
		log.Info("Run lifeVm failed", "err", err)
		vm.state.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = ErrMaxCodeSizeExceeded
	}

	/*	if vm.vmConfig.Debug && vm.depth == 0 {
			vm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
		}*/
	return ret, address, contract.Gas, err
}

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

	// Block information
	Coinbase common.Address // Provides information for COINBASE

	GasPrice    *big.Int // Provides information for GASPRICE
	GasLimit    uint64   // Provides information for GASLIMIT
	BlockNumber *big.Int // Provides information for NUMBER
	Time        *big.Int // Provides information for TIME
	Difficulty  *big.Int // Provides information for DIFFICULTY
	BlockHash   common.Hash
	TxHash      common.Hash
	TxIndex     uint64

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

func (context *Context) GetTxHash() common.Hash {
	return context.TxHash
}

func (context *Context) GetTxIdx() uint64 {
	return context.TxIndex
}

func (context *Context) GetCallGasTemp() uint64 {
	return context.callGasTemp
}

func (context *Context) GetGasPrice() int64 {
	return context.GasPrice.Int64()
}

func (context *Context) GetBlockHash(num uint64) common.Hash {
	return context.BlockHash
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
	txIndex, err := tx.GetTxIndex()
	if err != nil {
		panic("GetTxIndex failed")
	}
	return Context{
		Origin:      sender,
		GasPrice:    tx.GetGasPrice(),
		GasLimit:    tx.Fee().Uint64(),
		BlockNumber: new(big.Int).SetUint64(block.Number()),
		Time:        block.Timestamp(),
		Coinbase:    block.CoinBaseAddress(),
		Difficulty:  block.Difficulty().Big(),
		callGasTemp: tx.Fee().Uint64(),
		BlockHash:   block.Hash(),
		TxHash:      tx.CalTxId(),
		TxIndex:     uint64(txIndex),
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
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
