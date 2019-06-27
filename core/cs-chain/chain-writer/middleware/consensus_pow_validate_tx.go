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

package middleware

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strings"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/third-party/log"
)

// special tx validators
var txValidators = map[common.TxType]func(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error{
	// normal tx have no special validation
	common.TxType(common.AddressTypeNormal): func(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
		return nil
	},
	common.TxType(common.AddressTypeStake):       validRegisterTx,
	common.TxType(common.AddressTypeCancel):      validCancelTx,
	common.TxType(common.AddressTypeUnStake):     validUnStakeTx,
	common.TxType(common.AddressTypeEvidence):    validEvidenceTx,
	common.TxType(common.AddressTypeERC20):       validContractTx,
	common.TxType(common.AddressTypeEarlyReward): validEarlyTokenTx,
	common.TxType(common.AddressTypeContractCall): func(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
		return nil
	},
	common.TxType(common.AddressTypeContractCreate): func(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
		return nil
	},
}

//type TxContext struct {
//	MiddlewareContext
//
//	Tx model.AbstractTransaction
//	Chain ChainInterface
//	BlockHeight uint64
//}

// NewValidatorTx create a validator for transactions
func NewTxValidatorForRpcService(chain ChainInterface) *TxValidatorForRpcService {
	return &TxValidatorForRpcService{Chain: chain}
}

type TxValidatorForRpcService struct {
	Chain ChainInterface
}

// Valid valid transactions
func (v *TxValidatorForRpcService) Valid(tx model.AbstractTransaction) error {
	return validTx(tx, v.Chain, 0)
}

func ValidateBlockTxs(c *BlockContext) Middleware {
	return func() error {
		txs := c.Block.GetAbsTransactions()
		targetRoot := model.DeriveSha(model.AbsTransactions(txs))
		pbft_log.Info("the header tx root is:","root",c.Block.TxRoot().Hex())
		pbft_log.Info("the calculated tx root is:","root",targetRoot.Hex())
		pbft_log.Info("the block txs is:","len",len(txs))
		for _,tx := range txs{
			pbft_log.Info("the tx is:","tx",tx)
		}
		if !targetRoot.IsEqual(c.Block.TxRoot()) {
			return errors.New(fmt.Sprintf("tx root not match, target: %v, root in block: %v", targetRoot.Hex(), c.Block.TxRoot().Hex()))
		}

		if c.Block.IsSpecial() {
			if txs != nil {
				return errors.New("special block should not have transactions")
			} else {
				return c.Next()
			}
		}

		// start:=time.Now()
		for _, tx := range txs {
			if err := validTx(tx, c.Chain, c.Block.Number()); err != nil {
				return err
			}
		}

		return c.Next()
	}
}

// TODO The size of transaction will influence the transaction fee?
//func ValidTxSizeM(c *TxContext) Middleware {
//	return func() error {
//		return ValidTxSize(c.Tx)
//	}
//}
//
//// check whether the sender and the amount are correct
//func ValidTxSenderM(c *TxContext) Middleware {
//	return func() error {
//		return ValidTxSender(c.Tx, c.Chain, c.BlockHeight)
//	}
//}
//
//// do checking for different types of transactions
//func ValidTxByTypeM(c *TxContext) Middleware {
//	return func() error {
//		validator := txValidators[c.Tx.GetType()]
//		if validator == nil {
//			return g_error.ErrInvalidTxType
//		}
//		if err := validator(c.Tx, c.Chain, c.BlockHeight); err != nil {
//			return err
//		}
//		return nil
//	}
//}

func ValidTxSize(tx model.AbstractTransaction) error {
	if tx.Size() > chain_config.MaxTxSize {
		return g_error.ErrTxOverSize
	}
	return nil
}

// valid sender and amount
func ValidTxSender(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	economy := chain.GetEconomyModel()
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}

	if tx.GetType() == common.AddressTypeContractCreate || tx.GetType() == common.AddressTypeContractCall {
		gas, err := model.IntrinsicGas(tx.ExtraData(), tx.GetType() == common.AddressTypeContractCreate , true)
		if err !=nil{
			return err
		}

		if gas > tx.GetGasLimit() {
			return fmt.Errorf("gas limit is to low, need:%v got:%v",gas,tx.GetGasLimit())
		}

	}else{
		// valid tx fee
		if tx.Fee().Cmp(economy_model.GetMinimumTxFee(tx.Size())) == -1 {
			log.Error("the tx fee is:", "fee", tx.Fee(),"needFee",economy_model.GetMinimumTxFee(tx.Size()))
			return g_error.ErrTxFeeTooLow
		}
	}

	// valid tx fee
/*	if tx.Fee().Cmp(economy_model.GetMinimumTxFee(tx.Size())) == -1 {
		log.Info("the tx fee is:", "fee", tx.Fee(),"needFee",economy_model.GetMinimumTxFee(tx.Size()))
		return g_error.ErrTxFeeTooLow
	}*/

	// log.Info("ValidTxSender the blockHeight is:","blockHeight",blockHeight)
	state, err := getPreStateForHeight(blockHeight, chain)
	if err != nil {
		return err
	}
	credit, err := state.GetBalance(sender)
	log.Info("ValidTxSender#credit", "credit",  credit)
	if err != nil {
		return err
	}

	// get locked money
	lockValue, err := economy.GetAddressLockMoney(sender, chain.CurrentBlock().Number())
	if err != nil {
		return err
	}
	usage := big.NewInt(0).Add(tx.Amount(), tx.Fee())
	usage.Add(usage, lockValue)


	log.Info("the credit and the usage is:","credit",credit,"usage",usage)
	if credit.Cmp(usage) < 0 {
		return state_processor.NotEnoughBalanceError
	}
	return nil
}

// do checking for different types of transactions
func ValidTxByType(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) Middleware {
	return func() error {
		validator := txValidators[tx.GetType()]
		if validator == nil {
			return g_error.ErrInvalidTxType
		}
		if err := validator(tx, chain, blockHeight); err != nil {
			return err
		}
		return nil
	}
}

/*

1. valid transactions signature
2. valid transactions according to account balance (balance is always positive)
3. valid transactions type is logical for safety requirements

*/
func validTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	// start := time.Now()
	if err := ValidTxSender(tx, chain, blockHeight); err != nil {
		return err
	}

	if err := ValidTxSize(tx); err != nil {
		return err
	}

	validator := txValidators[tx.GetType()]
	if validator == nil {
		return errors.New(fmt.Sprintf("no validator for tx, type: %v", tx.GetType()))
	}

	// start = time.Now()
	if err := validator(tx, chain, blockHeight); err != nil {
		return err
	}
	return nil
}

func validRegisterTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	if tx == nil || chain == nil{
		log.Error("the tx and chain:","tx",tx,"chain",chain)
		return errors.New("the tx or chain is nil")
	}
	if tx.Amount().Cmp(economy_model.MiniPledgeValue) == -1{
		return errors.New("the register tx delegate is lower than MiniPledgeValue")
	}
	return nil
}

func validUnStakeTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	if err := haveStack(tx, chain, blockHeight); err != nil {
		return err
	}
	if err := validUnStakeTime(tx, chain, blockHeight); err != nil {
		return err
	}
	return nil
}

func validCancelTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}
	state, err := getPreStateForHeight(blockHeight, chain)
	if err != nil {
		return err
	}

	// whether sent register tx
	stake, err := state.GetStake(sender)
	if err != nil {
		return err
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return state_processor.SendRegisterTxFirst
	}

	// whether sent cancel tx
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock != 0 {
		return state_processor.SendRegisterTxFirst
	}
	return nil
}

func validContractTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	curState, err := chain.CurrentState()
	if err != nil {
		return err
	}

	return contract.NewProcessor(curState, blockHeight).Process(tx)
}

func validEarlyTokenTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	return nil
}

func validEvidenceTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	if err := conflictVote(tx, chain, blockHeight); err != nil {
		return err
	}
	if err := validEvidenceTime(tx, chain, blockHeight); err != nil {
		return err
	}
	if err := validTargetStake(tx, chain, blockHeight); err != nil {
		return err
	}
	return nil
}

// return current state if height == 0
func getPreStateForHeight(height uint64, reader ChainInterface) (s *state_processor.AccountStateDB, err error) {
	if height == 0 {
		s, err = reader.CurrentState()
	} else {
		s, err = reader.StateAtByBlockNumber(height - 1)
	}
	return
}

func validEvidenceTime(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	chainReader := chain
	//config := nodeContext.ChainConfig()
	config := chain_config.GetChainConfig()

	target := tx.To()
	state, err := getPreStateForHeight(blockHeight, chainReader)
	if err != nil {
		return err
	}
	targetNormal := cs_crypto.GetNormalAddressFromEvidence(*target)
	lastBlock, err := state.GetLastElect(targetNormal)
	if err != nil {
		return err
	}

	if lastBlock != 0 {
		// lastBlock < (current +1) < (lastBlock/SlotSize + StakeLockSlot)*SlotSize
		current := chainReader.CurrentBlock().Number()
		slotSpace := (current+1)/config.SlotSize - lastBlock/config.SlotSize
		if slotSpace > config.StakeLockSlot {
			return errors.New("invalid evidence time")
		}
		if current < lastBlock {
			return errors.New("invalid evidence time")
		}
	}
	// == 0 no err?
	return nil
}

func conflictVote(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	extraData := tx.ExtraData()
	proofData := model.Proofs{}
	if err := rlp.DecodeBytes(extraData, &proofData); err != nil {
		return errors.New(fmt.Sprintf("decode proof data failed: %v", err))
	}
	voteA := proofData.VoteA
	voteB := proofData.VoteB

	// valid vote
	err := voteA.Valid()
	if err != nil {
		return err
	}
	err = voteB.Valid()
	if err != nil {
		return err
	}

	// Two vote conflict check
	if voteA.GetType() != voteB.GetType() || voteA.GetViewID() != voteB.GetViewID() || voteA.GetHeight() != voteB.GetHeight() || strings.Compare(voteA.GetBlockHash(), voteB.GetBlockHash()) == 0 || !voteA.GetAddress().IsEqual(voteB.GetAddress()) {
		return errors.New("vote not conflict")
	}

	// Test target match voter
	if !voteA.GetAddress().IsEqual(cs_crypto.GetNormalAddressFromEvidence(*tx.To())) {
		return errors.New("invalid to address")
	}
	return nil
}

func validUnStakeTime(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}
	chainReader := chain
	//config := nodeContext.ChainConfig()
	config := chain_config.GetChainConfig()
	state, err := getPreStateForHeight(blockHeight, chainReader)
	if err != nil {
		return err
	}

	// whether sent register tx
	stake, err := state.GetStake(sender)
	if err != nil {
		return err
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return state_processor.SendRegisterTxFirst
	}

	// whether sent cancel tx
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock == 0 {
		return state_processor.SendCancelTxFirst
	}

	// whether in lockup period
	current := chainReader.CurrentBlock().Number()
	slotSpace := (current+1)/config.SlotSize - lastBlock/config.SlotSize
	if slotSpace < config.StakeLockSlot {
		return errors.New("invalid unStake time")
	}
	return nil
}

func haveStack(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}

	// get current stake
	currentStake, err := chain.CurrentState()
	if err != nil {
		return err
	}

	stake, err := currentStake.GetStake(sender)
	if err != nil {
		return err
	}

	if stake.Cmp(big.NewInt(0)) == 0 {
		return state_processor.NotEnoughStakeErr
	}

	return nil
}

/*
TxTargetStakeValidator is to validate the target has stake or is a validator candidate.
It consider to be valid, the target's stake more than 0.
It implemented TransactionValidator interface.
 */
func validTargetStake(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	to := tx.To()
	target := cs_crypto.GetNormalAddressFromEvidence(*to)
	currentStake, err := chain.CurrentState()
	stake, err := currentStake.GetStake(target)
	if err != nil {
		return err
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return errors.New("not enough stake")
	}
	return nil
}
