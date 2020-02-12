package base

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/life/exec"
	"github.com/dipperin/dipperin-core/third_party/life/mem-manage"
	"math/big"
)

var DEFAULT_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: mem_manage.DefaultPageSize,
}

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

type AccountRef common.Address

func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, common.Address, common.Address, *big.Int) error
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

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
		return gerror.ErrVMTransfer
	}
	return nil
}

