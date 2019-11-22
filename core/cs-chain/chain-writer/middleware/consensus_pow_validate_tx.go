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
	"fmt"
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
var txValidators = map[common.TxType]func(tx model.AbstractTransaction, conf *validTxNeedConfig) error{
	// normal tx have no special validation
	common.TxType(common.AddressTypeNormal): func(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
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

type validTxNeedConfig struct {
	economyModel       economy_model.EconomyModel
	preState           *state_processor.AccountStateDB
	curState           *state_processor.AccountStateDB
	currentBlockNumber uint64
	validBlockNumber   uint64
}

func newValidTxSenderNeedConfig(chain ChainInterface, blockNumber uint64) *validTxNeedConfig {
	var tmpConfig validTxNeedConfig

	tmpConfig.economyModel = chain.GetEconomyModel()
	tmpConfig.validBlockNumber = blockNumber
	tmpConfig.currentBlockNumber = chain.CurrentHeader().GetNumber()
	preState, err := getPreStateForHeight(blockNumber, chain)
	if err != nil {
		panic(fmt.Sprintf("newValidTxSenderNeedConfig get preState error,blockNumber:%v", tmpConfig.validBlockNumber))
	}
	tmpConfig.preState = preState

	curState,err := chain.CurrentState()
	if err != nil {
		panic(fmt.Sprintf("newValidTxSenderNeedConfig get curState error,blockNumber:%v", tmpConfig.validBlockNumber))
	}
	tmpConfig.curState = curState
	return &tmpConfig
}

// NewValidatorTx create a validator for transactions
func NewTxValidatorForRpcService(chain ChainInterface) *TxValidatorForRpcService {
	return &TxValidatorForRpcService{Chain: chain}
}

type TxValidatorForRpcService struct {
	Chain ChainInterface
}

// Valid valid transactions
func (v *TxValidatorForRpcService) Valid(tx model.AbstractTransaction) error {
	conf := newValidTxSenderNeedConfig(v.Chain, 0)
	return validTx(tx, conf)
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

		conf := newValidTxSenderNeedConfig(c.Chain, c.Block.Number())
		// start:=time.Now()
		for _, tx := range txs {
			if err := validTx(tx, conf); err != nil {
				return err
			}
		}
		log.Middleware.Info("ValidateBlockTxs success")
		return c.Next()
	}
}

// valid sender and amount
func ValidTxSender(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	economy := conf.economyModel
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
	credit, err := conf.preState.GetBalance(sender)
	log.Info("ValidTxSender#credit", "credit", credit)
	if err != nil {
		return err
	}

	// get locked money
	lockValue, err := economy.GetAddressLockMoney(sender, conf.currentBlockNumber)
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
		conf := newValidTxSenderNeedConfig(chain, blockHeight)
		if err := validator(tx,conf); err != nil {
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
func validTx(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	// start := time.Now()
	if err := ValidTxSender(tx, conf); err != nil {
		return err
	}

	validator := txValidators[tx.GetType()]
	if validator == nil {
		return g_error.ErrInvalidTxType
	}

	// start = time.Now()
	if err := validator(tx, conf); err != nil {
		return err
	}
	return nil
}

func validRegisterTx(tx model.AbstractTransaction,conf *validTxNeedConfig) error {
	if tx.Amount().Cmp(economy_model.MiniPledgeValue) == -1 {
		return g_error.ErrTxDelegatesNotEnough
	}
	return nil
}

func validUnStakeTx(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	if err := haveStack(tx, conf); err != nil {
		return err
	}
	if err := validUnStakeTime(tx, conf); err != nil {
		return err
	}
	return nil
}

func validCancelTx(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}

	// whether sent register tx
	stake, err := conf.preState.GetStake(sender)
	if err != nil {
		return err
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return g_error.ValidateSendRegisterTxFirst
	}

	// whether sent cancel tx
	lastBlock, err := conf.preState.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock != 0 {
		return g_error.ValidateSendRegisterTxFirst
	}
	return nil
}

func validContractTx(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	return contract.NewProcessor(conf.curState, conf.currentBlockNumber).Process(tx)
}

func validEarlyTokenTx(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	return nil
}

func validContractCreateTx(tx model.AbstractTransaction,conf *validTxNeedConfig) error {
	return nil
}

func validContractCallTx(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	return nil
}

func validEvidenceTx(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	if err := conflictVote(tx, conf); err != nil {
		return err
	}
	if err := validEvidenceTime(tx, conf); err != nil {
		return err
	}
	if err := validTargetStake(tx,conf); err != nil {
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

func validEvidenceTime(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	//config := nodeContext.ChainConfig()
	config := chain_config.GetChainConfig()

	target := tx.To()

	targetNormal := cs_crypto.GetNormalAddressFromEvidence(*target)
	lastBlock, err := conf.preState.GetLastElect(targetNormal)
	if err != nil {
		return err
	}

	if lastBlock != 0 {
		// lastBlock < (current +1) < (lastBlock/SlotSize + StakeLockSlot)*SlotSize
		current := conf.currentBlockNumber
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

func conflictVote(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
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

func validUnStakeTime(tx model.AbstractTransaction,conf *validTxNeedConfig) error {
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}
	//config := nodeContext.ChainConfig()
	config := chain_config.GetChainConfig()
	// whether sent register tx
	stake, err := conf.preState.GetStake(sender)
	if err != nil {
		return err
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return g_error.ValidateSendRegisterTxFirst
	}

	// whether sent cancel tx
	lastBlock, err := conf.preState.GetLastElect(sender)
	if err != nil {
		return err
	}
	if lastBlock == 0 {
		return g_error.ValidateSendCancelTxFirst
	}

	// whether in lockup period
	current := conf.currentBlockNumber
	slotSpace := (current+1)/config.SlotSize - lastBlock/config.SlotSize
	if slotSpace < config.StakeLockSlot {
		return g_error.ErrInvalidUnStakeTime
	}
	return nil
}

func haveStack(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	singer := tx.GetSigner()
	sender, err := tx.Sender(singer)
	if err != nil {
		return err
	}

	stake, err := conf.curState.GetStake(sender)
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
func validTargetStake(tx model.AbstractTransaction, conf *validTxNeedConfig) error {
	to := tx.To()
	target := cs_crypto.GetNormalAddressFromEvidence(*to)
	stake, err := conf.curState.GetStake(target)
	if err != nil {
		return err
	}
	if stake.Cmp(big.NewInt(0)) == 0 {
		return g_error.ErrTxTargetStakeNotEnough
	}
	return nil
}
