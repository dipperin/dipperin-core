// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"math/big"
	"errors"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
)

var (
	TransactionTypeError  = errors.New("transaction type miss match with processor function")
	NotEnoughBalanceError = errors.New("not enough balance error")
	NotEnoughStakeErr     = errors.New("target account stake not enough")
	SenderNotExistErr     = errors.New("sender account does not exist")
	ReceiverNotExistErr   = errors.New("receiver account does not exist")
	SendRegisterTxFirst   = errors.New("target need to send register transaction first")
	SendCancelTxFirst     = errors.New("target need to send cancel transaction first")
)

/*
Basic operations
Stake money from balance
*/
func (state *AccountStateDB) Stake(addr common.Address, amount *big.Int) error {
	balance, err := state.GetBalance(addr)
	if err != nil || balance.Cmp(amount) < 0 {
		log.Debug("stake failed", "addr", addr.Hex(), "err", NotEnoughBalanceError, "balance", balance, "amount", amount)
		return NotEnoughBalanceError
	}
	err = state.SubBalance(addr, amount)
	if err != nil {
		return err
	}
	err = state.AddStake(addr, amount)
	if err != nil {
		return err
	}
	err = state.SetLastElect(addr, uint64(0))
	if err != nil {
		return err
	}
	pbft_log.Info("stake money", "address", addr.Hex(), "amount", amount)
	return nil
}

/*
Retrieval Stake
Sub stake and add balance
*/
func (state *AccountStateDB) UnStake(addr common.Address) error {
	amount, err := state.GetStake(addr)
	if err != nil {
		return err
	}
	if amount.Cmp(big.NewInt(0)) == 0 {
		log.Warn("unStake value is zero", "address", addr.Hex())
		return NotEnoughStakeErr
	}
	err = state.AddBalance(addr, amount)
	if err != nil {
		return err
	}
	err = state.SubStake(addr, amount)
	if err != nil {
		return err
	}
	pbft_log.Info("unStake", "address", addr.Hex(), "amount", amount)
	return nil
}

/*Move stake to some address*/
func (state *AccountStateDB) MoveStakeToAddress(fromAdd common.Address, toAdd common.Address) error {
	amount, err := state.GetStake(fromAdd)
	if err != nil {
		return err
	}
	if amount.Cmp(big.NewInt(0)) == 0 {
		log.Warn("unStake value is zero", "address", fromAdd.Hex())
		return NotEnoughStakeErr
	}

	empty := state.IsEmptyAccount(toAdd)
	if empty {
		err := state.NewAccountState(toAdd)
		if err != nil {
			return err
		}
	}

	err = state.SubStake(fromAdd, amount)
	if err != nil {
		return err
	}
	err = state.AddBalance(toAdd, amount)
	if err != nil {
		return err
	}
	pbft_log.Debug("move stake", "from", fromAdd.Hex(), "to", toAdd.Hex(), "amount", amount)
	return nil
}

/*
* Process verifier related transactions
* Include AddPeerSet(Stake), Evidence, Cancel, UnStake

* Process register Tx
* Stake some money
  */
func (state *AccountStateDB) processStakeTx(tx model.AbstractTransaction) (err error) {

	//Check
	sender, _ := tx.Sender(nil)
	receiver := *(tx.To())
	if receiver.GetAddressType() != common.AddressTypeStake {
		return TransactionTypeError
	}

	//judging the balance of the deposit
	stake, err := state.GetStake(sender)
	minStake := int64(model.StakeValMin)
	if err != nil {
		return
	}
	if stake.Cmp(big.NewInt(0)) == 0 && tx.Amount().Cmp(big.NewInt(minStake)) == -1 {
		log.Debug("process register transaction failed", "err", NotEnoughStakeErr)
		return NotEnoughStakeErr
	}

	//Process
	err = state.Stake(sender, tx.Amount())
	if err != nil {
		return
	}
	pbft_log.Info("success process a register transaction", "Tx hash", tx.CalTxId().Hex())

	//TODO add receipt?
	return
}

/*
Process cancel Tx, num is processing block num
 */
func (state *AccountStateDB) processCancelTx(tx model.AbstractTransaction, num uint64) (err error) {
	//Check
	sender, _ := tx.Sender(nil)
	receiver := *(tx.To())
	if receiver.GetAddressType() != common.AddressTypeCancel {
		return TransactionTypeError
	}

	//have you sent a registered transaction
	stake, err := state.GetStake(sender)
	if err != nil {
		return
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return SendRegisterTxFirst
	}

	//have you sent a cancellation transaction
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock != 0 {
		return SendRegisterTxFirst
	}

	//Process
	err = state.SetLastElect(sender, num)
	if err != nil {
		return
	}
	pbft_log.Info("success process a cancel transaction", "Tx hash", tx.CalTxId().Hex())

	//TODO add receipt return?
	return
}

/*
Process UnStake Tx
Un stake money
 */
func (state *AccountStateDB) processUnStakeTx(tx model.AbstractTransaction) (err error) {

	//Check
	sender, _ := tx.Sender(nil)
	receiver := *(tx.To())
	if receiver.GetAddressType() != common.AddressTypeUnStake {
		return TransactionTypeError
	}

	//have you sent a registered transaction
	stake, err := state.GetStake(sender)
	if err != nil {
		return
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return SendRegisterTxFirst
	}

	//have you sent a cancellation transaction
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return
	}
	if lastBlock == 0 {
		return SendCancelTxFirst
	}

	//Process
	err = state.UnStake(sender)
	if err != nil {
		return
	}
	pbft_log.Info("success process a unStake transaction", "Tx hash", tx.CalTxId().Hex())

	//TODO add receipt return?
	return
}

/*
Process Evidence Tx
Punish target account
Move all target account stake to the sender of this transaction
*/
func (state *AccountStateDB) processEvidenceTx(tx model.AbstractTransaction) (err error) {

	//Check
	sender, _ := tx.Sender(nil)
	receiver := *(tx.To())
	originalReceiver := common.Address{}
	if receiver.GetAddressType() != common.AddressTypeEvidence {
		return TransactionTypeError
	}
	if sender.GetAddressType() != common.AddressTypeNormal {
		return errors.New("sender must be a normal account")
	}
	originalReceiver = cs_crypto.GetNormalAddressFromEvidence(receiver)
	if empty := state.IsEmptyAccount(originalReceiver); empty {
		return ReceiverNotExistErr
	}

	//Process
	err = state.MoveStakeToAddress(originalReceiver, sender)
	if err != nil {
		return err
	}

	//TODO add receipt return
	return nil
}
