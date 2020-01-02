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
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"go.uber.org/zap"
	"math/big"
)

/*
Basic operations
Stake money from balance
*/
func (state *AccountStateDB) Stake(addr common.Address, amount *big.Int) error {
	balance, err := state.GetBalance(addr)
	if err != nil || balance.Cmp(amount) < 0 {
		log.DLogger.Debug("stake failed", zap.String("addr", addr.Hex()), zap.Error(g_error.ErrBalanceNotEnough), zap.Any("balance", balance), zap.Any("amount", amount))
		return g_error.ErrBalanceNotEnough
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
	log.DLogger.Info("stake money", zap.String("address", addr.Hex()), zap.Any("amount", amount))
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
		log.DLogger.Warn("unStake value is zero", zap.String("address", addr.Hex()))
		return g_error.ErrStakeNotEnough
	}
	err = state.AddBalance(addr, amount)
	if err != nil {
		return err
	}
	err = state.SubStake(addr, amount)
	if err != nil {
		return err
	}
	log.DLogger.Info("unStake", zap.String("address", addr.Hex()), zap.Any("amount", amount))
	return nil
}

/*Move stake to some address*/
func (state *AccountStateDB) MoveStakeToAddress(fromAdd common.Address, toAdd common.Address) error {
	amount, err := state.GetStake(fromAdd)
	if err != nil {
		return err
	}
	if amount.Cmp(big.NewInt(0)) == 0 {
		log.DLogger.Warn("unStake value is zero", zap.String("address", fromAdd.Hex()))
		return g_error.ErrStakeNotEnough
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
	log.DLogger.Debug("move stake", zap.String("from", fromAdd.Hex()), zap.String("to", toAdd.Hex()), zap.Any("amount", amount))
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
		return g_error.ErrTxTypeNotMatch
	}

	//judging the balance of the deposit
	stake, err := state.GetStake(sender)
	minStake := int64(economy_model.MiniPledgeValue.Int64())
	if err != nil {
		return
	}

	if stake.Cmp(big.NewInt(0)) == 0 && tx.Amount().Cmp(big.NewInt(minStake)) == -1 {
		log.DLogger.Debug("process register transaction failed", zap.Error(g_error.ErrStakeNotEnough))
		return g_error.ErrStakeNotEnough
	}

	//Process
	err = state.Stake(sender, tx.Amount())
	if err != nil {
		return
	}
	log.DLogger.Info("success process a register transaction", zap.String("Tx hash", tx.CalTxId().Hex()))

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
		return g_error.ErrTxTypeNotMatch
	}

	//have you sent a registered transaction
	stake, err := state.GetStake(sender)
	if err != nil {
		return
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return g_error.StateSendRegisterTxFirst
	}

	//have you sent a cancellation transaction
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock != 0 {
		return g_error.StateSendRegisterTxFirst
	}

	//Process
	err = state.SetLastElect(sender, num)
	if err != nil {
		return
	}
	log.DLogger.Info("success process a cancel transaction", zap.String("Tx hash", tx.CalTxId().Hex()))

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
		return g_error.ErrTxTypeNotMatch
	}

	//have you sent a registered transaction
	stake, err := state.GetStake(sender)
	if err != nil {
		return
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return g_error.StateSendRegisterTxFirst
	}

	//have you sent a cancellation transaction
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return
	}
	if lastBlock == 0 {
		return g_error.StateSendCancelTxFirst
	}

	//Process
	err = state.UnStake(sender)
	if err != nil {
		return
	}
	log.DLogger.Info("success process a unStake transaction", zap.String("Tx hash", tx.CalTxId().Hex()))

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
		return g_error.ErrTxTypeNotMatch
	}
	if sender.GetAddressType() != common.AddressTypeNormal {
		return g_error.ErrAddressTypeNotMatch
	}
	originalReceiver = cs_crypto.GetNormalAddressFromEvidence(receiver)
	if empty := state.IsEmptyAccount(originalReceiver); empty {
		return g_error.ErrReceiverNotExist
	}

	//Process
	err = state.MoveStakeToAddress(originalReceiver, sender)
	if err != nil {
		return err
	}

	//TODO add receipt return
	return nil
}
