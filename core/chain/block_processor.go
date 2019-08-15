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

package chain

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"reflect"
)

var (
	//performanceInitial = uint64(30)
	reward  = float64(1)
	penalty = float64(-10)
)

type AccountDBChainReader interface {
	CurrentBlock() model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	GetVerifiers(round uint64) []common.Address
	StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error)

	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetLastChangePoint(block model.AbstractBlock) *uint64
	GetSlot(block model.AbstractBlock) *uint64
}

// process chain state before insert a block
type BlockProcessor struct {
	fullChain AccountDBChainReader
	*state_processor.AccountStateDB
	economyModel economy_model.EconomyModel
}

func NewBlockProcessor(fullChain AccountDBChainReader, preStateRoot common.Hash, db state_processor.StateStorage) (*BlockProcessor, error) {
	aDB, err := state_processor.NewAccountStateDB(preStateRoot, db)
	if err != nil {
		return nil, err
	}
	return &BlockProcessor{
		fullChain:      fullChain,
		AccountStateDB: aDB,
	}, nil
}

func (state *BlockProcessor) GetBlockHashByNumber(number uint64) common.Hash {
	if number > state.fullChain.CurrentBlock().Number() {
		log.Info("GetBlockHashByNumber failed, can't get future block")
		return common.Hash{}
	}
	return state.fullChain.GetBlockByNumber(number).Hash()
}

func (state *BlockProcessor) Process(block model.AbstractBlock, economyModel economy_model.EconomyModel) (err error) {
	log.Mpt.Debug("AccountStateDB Process begin~~~~~~~~~~~~~~", "pre state", state.PreStateRoot().Hex(), "blockId", block.Hash().Hex())

	state.economyModel = economyModel
	blockHeader := block.Header().(*model.Header)
	gasUsed := uint64(0)
	gasLimit := blockHeader.GasLimit
	// special block doesn't process txs
	if !block.IsSpecial() {
		if err = block.TxIterator(func(i int, tx model.AbstractTransaction) error {
			conf := state_processor.TxProcessConfig{
				Tx:       tx,
				Header:   blockHeader,
				GetHash:  state.GetBlockHashByNumber,
				GasUsed:  &gasUsed,
				GasLimit: &gasLimit,
			}
			innerError := state.ProcessTxNew(&conf)
			/*// unrecognized tx means no processing of the tx
			if innerError == g_error.UnknownTxTypeErr {
				log.Warn("unknown tx type", "type", tx.GetType())
				return g_error.UnknownTxTypeErr
			}*/
			if innerError != nil {
				return innerError
			}

			return nil
		}); err != nil {
			return err
		}
	}

	if err = state.ProcessExceptTxs(block, economyModel, false); err != nil {
		return
	}
	log.Mpt.Debug("AccountStateDB Process end~~~~~~~~~~~~~~~~~", "pre state", state.PreStateRoot().Hex())
	return
}

func (state *BlockProcessor) ProcessExceptTxs(block model.AbstractBlock, economyModel economy_model.EconomyModel, isProcessPackageBlock bool) (err error) {
	log.Mpt.Debug("ProcessExceptTxs begin", "pre state", state.PreStateRoot().Hex())
	state.economyModel = economyModel
	if block.Number() == 0 {
		log.Mpt.Debug("ProcessExceptTxs bug block num is 0")
		return nil
	}

	// do rewards
	if err = state.doRewards(block); err != nil {
		log.Mpt.Debug("ProcessExceptTxs doRewards failed", "storageErr", err)
		return
	}
	// process commits
	err = state.processCommitList(block, isProcessPackageBlock)
	log.Mpt.Debug("ProcessExceptTxs finished ---", "pre state", state.PreStateRoot().Hex())
	return
}

//Process verifiers and commit list of previous block
func (state *BlockProcessor) processCommitList(block model.AbstractBlock, isProcessPackageBlock bool) (err error) {
	if block.Number() != 0 {
		// fixme cur block v list is pre,
		previous := block.Number() - 1
		preBlock := state.fullChain.GetBlockByNumber(previous)
		if preBlock == nil {
			return g_error.NotHavePreBlockErr
		}

		preBlockSlot := state.fullChain.GetSlot(preBlock)
		verifiers := state.fullChain.GetVerifiers(*preBlockSlot)
		for index, ver := range verifiers {
			innerErr := state.ProcessVerifierNumber(ver)
			if innerErr != nil {
				log.Error("process block verifiers error", "storageErr", innerErr, "verifier", ver, "index", index)
				return innerErr
			}
		}

		// boot node verifier does't process verification
		verifications := block.GetVerifications()
		if preBlock.IsSpecial() {
			verifications = verifications[1:]
		}

		for _, ver := range verifications {
			innerErr := state.ProcessVerification(ver, 0)
			if innerErr != nil {
				log.Error("process block verifications error", "storageErr", innerErr, "verifier", ver)
				return innerErr
			}
		}

		// If the current block is at the transition point, penalize the non-working verifier, reward the voter with the vote
		// Calculated by comparing whether each prover has a change in the commit num in one round
		config := chain_config.GetChainConfig()
		slot := state.fullChain.GetSlot(block)
		if state.fullChain.IsChangePoint(block, isProcessPackageBlock) && *slot >= config.SlotMargin {

			//　firstStateBySlot is the state of the first block in a round
			var firstStateBySlot *state_processor.AccountStateDB
			lastPoint := state.fullChain.GetLastChangePoint(block)

			//If the previous block is also a change point, then there is only one block in this round.
			if state.fullChain.IsChangePoint(preBlock, false) {
				firstStateBySlot = state.AccountStateDB
			} else {
				firstStateBySlot, _ = state.fullChain.StateAtByBlockNumber(*lastPoint + 1)
			}

			log.Info("process performance", "slot", slot, "current num", block.Number(), "len(vers)", verifiers)
			for _, ver := range verifiers {
				commitNum, _ := state.GetCommitNum(ver)
				firstCommitNumBySlot, err := firstStateBySlot.GetCommitNum(ver)
				if err != nil {
					log.Error("process performance error")
					return err
				}
				amount := reward
				if commitNum == firstCommitNumBySlot {
					amount = penalty
				}
				state.ProcessPerformance(ver, amount)
			}
		}
	}
	return
}

func (state *BlockProcessor) doRewards(block model.AbstractBlock) (err error) {
	//get earlyToken contract
	earlyContractV, err := state.GetContract(contract.EarlyContractAddress, reflect.TypeOf(contract.EarlyRewardContract{}))
	if err != nil {
		return
	}

	earlyContract := earlyContractV.Interface().(*contract.EarlyRewardContract)

	earlyContract.AccountDB = state
	earlyContract.Early = state.economyModel.GetFoundation()

	// reward coinBase, special block doesn't reward CoinBase
	if !block.IsSpecial() {
		err = state.RewardCoinBase(block, earlyContract)
		if err != nil {
			log.Error("process block reward coin base error", "num", block.Number(), "storageErr", err)
			return err
		}
	}

	// reward verifier
	if block.Number() != 0 && block.Number() != 1 {
		// fixme reward to cur block verifiers list
		previous := block.Number() - 1
		//	log.Info("the chainReader is:","chainReader",state.fullChain)
		//	log.Info("the previous is:","previous",previous)
		preBlock := state.fullChain.GetBlockByNumber(previous)

		//log.Info("the preBlock info is:","blockNumber",preBlock.Number(),"ver",preBlock.GetVerifications())
		if preBlock == nil {
			return g_error.NotHavePreBlockErr
		}

		err = state.RewardByzantiumVerifier(block, earlyContract)
		if err != nil {
			log.Error("process block reward pre block verifiers error", "num", block.Number(), "storageErr", err)
			return
		}
	}

	err = state.PutContract(contract.EarlyContractAddress, earlyContractV)
	if err != nil {
		return
	}
	return
}
