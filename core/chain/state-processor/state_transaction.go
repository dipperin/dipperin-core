package state_processor

import (
	"github.com/dipperin/dipperin-core/core/model"
	"math/big"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vm"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/third-party/log"
)

type StateTransition struct {
	gp         *uint64
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      vm.StateDB
	lifeVm     *vm.VM
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(vm *vm.VM, msg Message, gp *uint64) *StateTransition {
	return &StateTransition{
		gp:       gp,
		lifeVm:   vm,
		msg:      msg,
		gasPrice: msg.GasPrice(),
		value:    msg.Value(),
		data:     msg.Data(),
		state:    vm.GetStateDB(),
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(vm *vm.VM, msg Message, gp *uint64) ([]byte, uint64, bool, error) {
	return NewStateTransition(vm, msg, gp).TransitionDb()
}

// to returns the recipient of the message.
func (st *StateTransition) to() common.Address {
	if st.msg == nil || st.msg.To() == nil {
		return common.Address{}
	}
	return *st.msg.To()
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gas < amount {
		return g_error.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	msgVal := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), st.gasPrice)
	if st.lifeVm.GetStateDB().GetBalance(st.msg.From()).Cmp(msgVal) < 0 {
		return g_error.ErrInsufficientBalanceForGas
	}

	log.Info("Call buyGas", "gasPool", *st.gp, "balance", st.lifeVm.GetStateDB().GetBalance(st.msg.From()), "value", msgVal, "st.msg.Gas()", st.msg.Gas())


	if *st.gp < st.msg.Gas() {
		return g_error.ErrGasLimitReached
	}
	*st.gp -= st.msg.Gas()
	st.gas += st.msg.Gas()

	st.initialGas = st.msg.Gas()
	st.state.SubBalance(st.msg.From(), msgVal)

	log.Info("BuyGas successful", "gasPool", *st.gp, "gas", st.gas, "initialGas", st.initialGas)
	log.Info("BuyGas successful", "balance", st.lifeVm.GetStateDB().GetBalance(st.msg.From()))
	return nil
}

func (st *StateTransition) preCheck() error {
	// Make sure this transaction's nonce is correct.
	if st.msg.CheckNonce() {
		nonce := st.lifeVm.GetStateDB().GetNonce(st.msg.From())
		//log.Info("StateTransition#preCheck", "nonce", nonce, "st.msg.Nonce()", st.msg.Nonce())
		if nonce < st.msg.Nonce() {
			return g_error.ErrNonceTooHigh
		} else if nonce > st.msg.Nonce() {
			return g_error.ErrNonceTooLow
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the result including the used gas. It returns an error if failed.
// An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, usedGas uint64, failed bool, err error) {
	// init initialGas value = txMsg.gas
	if err = st.preCheck(); err != nil {
		return
	}
	msg := st.msg
	sender := vm.AccountRef(msg.From())
	contractCreation := msg.To().GetAddressType() == common.AddressTypeContractCreate

	// Pay intrinsic gas
	gas, err := model.IntrinsicGas(st.data, contractCreation, true)
	if err != nil {
		return nil, 0, false, err
	}

	log.Info("TransitionDb st.gas", "st.gas", st.gas, "gas", gas)
	if err = st.useGas(gas); err != nil {
		log.Error("TransitionDb#IntrinsicGas", "err", err)
		return nil, 0, false, err
	}

	log.Info("IntrinsicGas Used", "used", gas, "left", st.gas)
	var (
		lifeVm = st.lifeVm
		// lifeVm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
	)
	if contractCreation {
		// data = rlp
		ret, _, st.gas, vmerr = lifeVm.Create(sender, st.data, st.gas, st.value)
	} else {
		// Increment the nonce for the next transaction
		st.lifeVm.GetStateDB().AddNonce(msg.From(), uint64(1))
		//log.Debug("Nonce tracking: SetNonce", "from", msg.From(), "nonce", st.state.GetNonce(sender.Address()))
		ret, st.gas, vmerr = lifeVm.Call(sender, st.to(), st.data, st.gas, st.value)
	}
	if vmerr != nil {
		log.Info("VM returned with error", "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == g_error.ErrInsufficientBalance {
			return nil, 0, false, vmerr
		}
	}
	st.refundGas()
	//todo handle reward in ProcessExceptTxs
	//st.state.AddBalance(st.lifeVm.Coinbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))
	return ret, st.gasUsed(), vmerr != nil, err
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	/*refund := st.gasUsed() / 2
	if refund > st.lifeVm.GetStateDB().GetRefund() {
		refund = st.lifeVm.GetStateDB().GetRefund()
	}*/
	log.Info("Call refundGas", "gasPool", *st.gp, "balance", st.lifeVm.GetStateDB().GetBalance(st.msg.From()))
	// Return ETH for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.state.AddBalance(st.msg.From(), remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	if *st.gp > math.MaxUint64-st.gas {
		panic("gas pool pushed above uint64")
	}
	*st.gp += st.gas
	log.Info("Gas Refund", "gasPool", *st.gp, "balance", st.lifeVm.GetStateDB().GetBalance(st.msg.From()))
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}
