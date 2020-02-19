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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strings"
)

// special tx validators
var txValidators = map[common.TxType]func(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error{
	// normal tx have no special validation
	common.TxType(common.AddressTypeNormal): func(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
		return nil
	},
	common.TxType(common.AddressTypeStake):          validRegisterTx,
	common.TxType(common.AddressTypeCancel):         validCancelTx,
	common.TxType(common.AddressTypeUnStake):        validUnStakeTx,
	common.TxType(common.AddressTypeEvidence):       validEvidenceTx,
	common.TxType(common.AddressTypeERC20):          validContractTx,
	common.TxType(common.AddressTypeEarlyReward):    validEarlyTokenTx,
	common.TxType(common.AddressTypeContractCall):   validContractCallTx,
	common.TxType(common.AddressTypeContractCreate): validContractCreateTx,
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
		log.Middleware.Info("ValidateBlockTxs start")
		txs := c.Block.GetAbsTransactions()
		targetRoot := model.DeriveSha(model.AbsTransactions(txs))
		if !targetRoot.IsEqual(c.Block.TxRoot()) {
			log.Error("tx root not match", "targetRoot", targetRoot.Hex(), "blockRoot", c.Block.TxRoot().Hex())
			return g_error.ErrTxRootNotMatch
		}

		if c.Block.IsSpecial() {
			if txs != nil {
				return g_error.ErrTxInSpecialBlock
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
		log.Middleware.Info("ValidateBlockTxs success")
		return c.Next()
	}
}

// valid sender and amount
func ValidTxSender(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	economy := chain.GetEconomyModel()
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}

	//check minimal gasUsed
	gas, err := model.IntrinsicGas(tx.ExtraData(), tx.GetType() == common.AddressTypeContractCreate, true)
	if err != nil {
		return err
	}

	if gas > tx.GetGasLimit() {
		log.Error("tx gas limit is too low", "need", gas, "got", tx.GetGasLimit())
		return g_error.ErrTxGasLimitNotEnough
	}

	// log.Info("ValidTxSender the blockHeight is:","blockHeight",blockHeight)
	state, err := getPreStateForHeight(blockHeight, chain)
	if err != nil {
		return err
	}
	credit, err := state.GetBalance(sender)
	log.Info("ValidTxSender#credit", "credit", credit)
	if err != nil {
		return err
	}

	// get locked money
	lockValue, err := economy.GetAddressLockMoney(sender, chain.CurrentBlock().Number())
	if err != nil {
		return err
	}

	gasFee := big.NewInt(0).Mul(big.NewInt(int64(tx.GetGasLimit())), tx.GetGasPrice())
	usage := big.NewInt(0).Add(tx.Amount(), gasFee)
	usage.Add(usage, lockValue)

	log.Info("the credit and the usage is:", "credit", credit, "usage", usage)
	if credit.Cmp(usage) < 0 {
		return g_error.ErrTxSenderBalanceNotEnough
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

	validator := txValidators[tx.GetType()]
	if validator == nil {
		return g_error.ErrInvalidTxType
	}

	// start = time.Now()
	if err := validator(tx, chain, blockHeight); err != nil {
		return err
	}
	return nil
}

func validRegisterTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	if tx.Amount().Cmp(economy_model.MiniPledgeValue) == -1 {
		return g_error.ErrTxDelegatesNotEnough
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
		return g_error.ValidateSendRegisterTxFirst
	}

	// whether sent cancel tx
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock != 0 {
		return g_error.ValidateSendRegisterTxFirst
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

func validContractCreateTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	return nil
}

func validContractCallTx(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
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
			return g_error.ErrInvalidEvidenceTime
		}
		if current < lastBlock {
			return g_error.ErrInvalidEvidenceTime
		}
	}
	// == 0 no err?
	return nil
}

func conflictVote(tx model.AbstractTransaction, chain ChainInterface, blockHeight uint64) error {
	extraData := tx.ExtraData()
	proofData := model.Proofs{}
	if err := rlp.DecodeBytes(extraData, &proofData); err != nil {
		return err
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
		return g_error.ErrEvidenceVoteNotConflict
	}

	// Test target match voter
	if !voteA.GetAddress().IsEqual(cs_crypto.GetNormalAddressFromEvidence(*tx.To())) {
		return g_error.ErrTxTargetAddressNotMatch
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
		return g_error.ValidateSendRegisterTxFirst
	}

	// whether sent cancel tx
	lastBlock, err := state.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock == 0 {
		return g_error.ValidateSendCancelTxFirst
	}

	// whether in lockup period
	current := chainReader.CurrentBlock().Number()
	slotSpace := (current+1)/config.SlotSize - lastBlock/config.SlotSize
	if slotSpace < config.StakeLockSlot {
		return g_error.ErrInvalidUnStakeTime
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
		return g_error.ErrTxSenderStakeNotEnough
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
		return g_error.ErrTxTargetStakeNotEnough
	}
	return nil
}
