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
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
)

//Reward Verifiers, current block reward previous block
func (state *BlockProcessor) RewardByzantiumVerifier(Block model.AbstractBlock, earlyContract *contract.EarlyRewardContract) (err error) {
	//use economy model calculate the verifier reward
	//should use preBlock as the parameter
	preBlock := state.fullChain.GetBlockByNumber(Block.Number() - 1)
	rewards, err := state.economyModel.GetVerifierDIPReward(preBlock)
	if err != nil {
		return err
	}
	rewardAddress, err := state.economyModel.GetDiffVerifierAddress(preBlock, Block)
	if err != nil {
		return err
	}

	for addressType, addresses := range rewardAddress {
		rewardValue := rewards[addressType]
		for _, address := range addresses {
			if empty := state.IsEmptyAccount(address); empty {
				if err = state.NewAccountState(address); err != nil {
					return err
				}
			}

			if err = state.AddBalance(address, rewardValue); err != nil {
				return err
			}
		}
	}

	//reward earlyToken
	err = earlyContract.RewardVerifier(rewards, preBlock.Number(), rewardAddress)
	if err != nil {
		return err
	}

	/*	storageErr = Block.VersIterator(func(i int, verification model.AbstractVerification, block model.AbstractBlock) (error) {
		innerErr := state.rewardVerifier(verification)
		if innerErr != nil{
			return innerErr
		}
		return nil
	})*/

	return
}

func (state *BlockProcessor) RewardCoinBase(block model.AbstractBlock, earlyContract *contract.EarlyRewardContract) (err error) {
	err = state.rewardCoinBase(block, earlyContract)
	return
}

//Reward Miner
func (state *BlockProcessor) rewardCoinBase(block model.AbstractBlock, earlyContract *contract.EarlyRewardContract) error {
	//Miner's reward has two parts, transaction fees and coin base reward
	transactionFees := block.GetTransactionFees()

	//use economy model calculate the mineMaster reward
	//coinBase := chain_config.FrontierBlockReward
	coinBase, err := state.economyModel.GetMineMasterDIPReward(block)
	if err != nil {
		return err
	}

	total := transactionFees.Add(transactionFees, coinBase)

	coinBaseAddress := block.CoinBaseAddress()
	if coinBaseAddress.IsEqual(common.Address{}) {
		return g_error.InvalidCoinBaseAddressErr
	}

	empty := state.IsEmptyAccount(coinBaseAddress)
	if empty {
		err := state.NewAccountState(coinBaseAddress)
		if err != nil {
			return err
		}
	}

	//reward early token
	err = earlyContract.RewardMineMaster(coinBase, block.Number(), coinBaseAddress)
	if err != nil {
		return err
	}

	log.Info("reward to coinBase", "total", total, "address", block.CoinBaseAddress(), "num", block.Number())
	err = state.AddBalance(coinBaseAddress, total)
	if err != nil {
		return err
	}
	return nil
}
