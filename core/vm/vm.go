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

package vm

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/core/vm/resolver"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/life/exec"
	"github.com/dipperin/dipperin-core/third-party/life/mem-manage"
	"go.uber.org/zap"
	"math/big"
	"sync/atomic"
)

var EmptyCodeHash = cs_crypto.Keccak256Hash(nil)

var DEFAULT_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: mem_manage.DefaultPageSize,
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

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (vm *VM) Cancel() {
	atomic.StoreInt32(&vm.abort, 1)
}

func (vm *VM) Call(caller resolver.ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if vm.vmConfig.NoRecursion && vm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if vm.depth > int(model2.CallCreateDepth) {
		return nil, gas, g_error.ErrDepth
	}
	// Fail if we're trying to transfer more than the available balanceMap
	if !vm.Context.CanTransfer(vm.state, caller.Address(), value) {
		return nil, gas, g_error.ErrInsufficientBalance
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
	contract := NewContract(caller, to, value, gas, input)
	log.DLogger.Info("Call#NewContract", zap.Any("callerAddr", contract.CallerAddress), zap.Any("caller", contract.Caller().Address()), zap.Any("self", contract.Self().Address()))
	contract.SetCode(&addr, vm.state.GetCodeHash(addr), vm.state.GetCode(addr))
	contract.SetAbi(&addr, vm.state.GetAbiHash(addr), vm.state.GetAbi(addr))

	//start := time.Now()

	// Capture the tracer start/end events in debug mode
	/*	if evm.vmConfig.Debug && evm.depth == 0 {
		evm.vmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, gas, value)

		defer func() { // Lazy evaluation of the parameters
			evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
		}()
	}*/
	if to.Address().GetAddressType() == common.AddressTypeContractCall {
		ret, err = run(vm, contract, false)
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		vm.state.RevertToSnapshot(snapshot)
		if err != g_error.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
			log.DLogger.Info("callContract Use", zap.Uint64("gasUsed", contract.Gas), zap.Uint64("gasLeft", contract.Gas))
		}
	}
	return ret, contract.Gas, err
}

func (vm *VM) DelegateCall(caller resolver.ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if vm.vmConfig.NoRecursion && vm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if vm.depth > int(model2.CallCreateDepth) {
		return nil, gas, g_error.ErrDepth
	}

	var (
		snapshot = vm.state.Snapshot()
		to       = AccountRef(caller.Address())
	)

	// Initialise a new contract and make initialise the delegate values
	contract := NewContract(caller, to, nil, gas, input).AsDelegate()
	log.DLogger.Info("DelegateCall#NewContract", zap.Any("callerAddr", contract.CallerAddress), zap.Any("caller", contract.Caller().Address()), zap.Any("self", contract.Self().Address()))
	contract.SetCode(&addr, vm.state.GetCodeHash(addr), vm.state.GetCode(addr))
	contract.SetAbi(&addr, vm.state.GetAbiHash(addr), vm.state.GetAbi(addr))

	if to.Address().GetAddressType() == common.AddressTypeContractCall {
		ret, err = run(vm, contract, false)
	}
	if err != nil {
		vm.GetStateDB().RevertToSnapshot(snapshot)
		if err != g_error.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, contract.Gas, err
}

func (vm *VM) Create(caller resolver.ContractRef, data []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	nonce, err := vm.state.GetNonce(caller.Address())
	if err != nil {
		log.DLogger.Error("can't get the caller nonce")
		return nil, common.Address{}, 0, err
	}
	contractAddr = cs_crypto.CreateContractAddress(caller.Address(), nonce)
	return vm.create(caller, data, gas, value, contractAddr)
}

func (vm *VM) create(caller resolver.ContractRef, data []byte, gas uint64, value *big.Int, address common.Address) (rest []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	defer func() {
		if er := recover(); er != nil {
			log.DLogger.Error("VM#create err  ", zap.Error(er.(error)))
			rest, contractAddr, leftOverGas, err = nil, common.Address{}, gas, er.(error)
		}
	}()
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if vm.depth > int(model2.CallCreateDepth) {
		return nil, common.Address{}, gas, g_error.ErrDepth
	}

	if !vm.CanTransfer(vm.state, caller.Address(), value) {
		return nil, common.Address{}, gas, g_error.ErrInsufficientBalance
	}
	vm.state.AddNonce(caller.Address(), uint64(1))

	// Ensure there's no existing contract already at the designated address
	contractHash := vm.state.GetCodeHash(address)
	nonce, _ := vm.state.GetNonce(address)
	if nonce != uint64(0) || (contractHash != common.Hash{} && contractHash != EmptyCodeHash) {
		return nil, common.Address{}, 0, g_error.ErrContractAddressCollision
	}

	// Create a new account on the state
	snapshot := vm.state.Snapshot()
	err = vm.state.CreateAccount(address)
	if err != nil {
		return nil, common.Address{}, 0, g_error.ErrContractAddressCreate
	}

	if err = vm.Transfer(vm.state, caller.Address(), address, value); err != nil {
		return nil, common.Address{}, 0, err
	}

	// initialise a new contract and set the data that is to be used by the
	// EVM. The contract is a scoped environment for this execution context
	// only.
	code, abi, rlpInit, err := ParseCreateExtraData(data)
	if err != nil {
		log.DLogger.Error("ParseCreateExtraData failed", zap.Error(err))
		return nil, common.Address{}, 0, err
	}
	contract := NewContract(caller, AccountRef(address), value, gas, rlpInit)
	contract.SetCode(&address, cs_crypto.Keccak256Hash(code), code)
	contract.SetAbi(&address, cs_crypto.Keccak256Hash(abi), abi)

	if vm.vmConfig.NoRecursion && vm.depth > 0 {
		return nil, address, gas, nil
	}

	/*	if api.vmConfig.Debug && api.depth == 0 {
		api.vmConfig.Tracer.CaptureStart(caller.Address(), address, true, code, gas, value)
	}*/
	//start := time.Now()

	ret, err := run(vm, contract, true)
	// check whether the max data size has been exceeded
	maxCodeSizeExceeded := len(ret) > model2.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the data. If the data could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		log.DLogger.Info("LifeVm run successful", zap.Uint64("gasLeft", contract.Gas))
		createDataGas := uint64(len(ret)+len(abi)) * model2.CreateDataGas
		if contract.UseGas(createDataGas) {
			vm.state.SetCode(address, ret)
			vm.state.SetAbi(address, abi)
			log.DLogger.Info("CreateDataGas Use", zap.Uint64("gasUsed", createDataGas), zap.Uint64("gasLeft", contract.Gas))
		} else {
			err = g_error.ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation data
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for data storage gas errors.
	if maxCodeSizeExceeded || (err != nil && err != g_error.ErrCodeStoreOutOfGas) {
		log.DLogger.Info("Run lifeVm failed", zap.Error(err))
		vm.state.RevertToSnapshot(snapshot)
		if err != g_error.ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract data size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = g_error.ErrMaxCodeSizeExceeded
	}

	/*	if api.vmConfig.Debug && api.depth == 0 {
		api.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}*/
	return ret, address, contract.Gas, err
}

func run(vm *VM, contract *Contract, create bool) ([]byte, error) {
	// call Interpreter.Run()
	return vm.Interpreter.Run(vm, contract, create)
}

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, common.Address, common.Address, *big.Int) error
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

type Context struct {
	// Message information
	Origin common.Address // Provides information for ORIGIN

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasPrice    *big.Int       // Provides information for GASPRICE
	GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
	TxHash      common.Hash

	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	//callGasTemp uint64

	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc

	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc
}

func (context *Context) GetTxHash() common.Hash {
	return context.TxHash
}

/*func (context *Context) GetCallGasTemp() uint64 {
	return context.callGasTemp
}
*/

func (context *Context) GetGasPrice() *big.Int {
	return context.GasPrice
}

func (context *Context) GetGasLimit() uint64 {
	return context.GasLimit
}

func (context *Context) GetBlockHash(num uint64) common.Hash {
	return context.GetHash(num)
}

func (context *Context) GetBlockNumber() *big.Int {
	return context.BlockNumber
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
func NewVMContext(tx model.AbstractTransaction, header model.AbstractHeader, GetHash GetHashFunc) Context {
	sender, _ := tx.Sender(tx.GetSigner())
	return Context{
		Origin:      sender,
		GasPrice:    tx.GetGasPrice(),
		GasLimit:    tx.GetGasLimit(),
		BlockNumber: new(big.Int).SetUint64(header.GetNumber()),
		//callGasTemp:  tx.Fee().Uint64(),
		TxHash:      tx.CalTxId(),
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		Coinbase:    header.CoinBaseAddress(),
		Time:        header.GetTimeStamp(),
		GetHash:     GetHash,
	}
}

// CanTransfer checks whether there are enough funds in the address' account to make a transfer.
// This does not take the necessary gas in to account to make the transfer valid.
func CanTransfer(db StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

// Transfer subtracts amount from sender and adds amount to recipient using the given Db
func Transfer(db StateDB, sender, recipient common.Address, amount *big.Int) error {
	err1 := db.SubBalance(sender, amount)
	err2 := db.AddBalance(recipient, amount)
	if err1 != nil || err2 != nil {
		return g_error.ErrVMTransfer
	}
	return nil
}

/*type Caller struct {
	addr common.Address
}

func (c *Caller) Address() common.Address {
	return c.addr
}

func AccountRef(addr common.Address) resolver.ContractRef {
	return &Caller{addr}
}*/

type AccountRef common.Address

func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }
