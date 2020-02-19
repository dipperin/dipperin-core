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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockProcessor_RewardByzantiumVerifier_Error(t *testing.T) {
	var earlyTokenContract contract.EarlyRewardContract
	block := createBlock(20)

	// GetDiffVerifierAddress error
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{addrErr: EconomyErr}
	err = processor.RewardByzantiumVerifier(block, &earlyTokenContract)
	assert.Equal(t, EconomyErr, err)

	// GetVerifierDIPReward error
	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{DIPErr: EconomyErr}
	err = processor.RewardByzantiumVerifier(block, &earlyTokenContract)
	assert.Equal(t, EconomyErr, err)

	// AddBalance error
	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{setErr: TrieError})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{}
	err = processor.RewardByzantiumVerifier(block, &earlyTokenContract)
	assert.Equal(t, TrieError, err)

	// NewAccountState error
	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{
		getErr: TrieError,
		setErr: TrieError,
	})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{}
	err = processor.RewardByzantiumVerifier(block, &earlyTokenContract)
	assert.Equal(t, TrieError, err)

	// No error
	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{}
	err = processor.RewardByzantiumVerifier(block, &earlyTokenContract)
	assert.NoError(t, err)
}

func TestBlockProcessor_RewardCoinBase_Error(t *testing.T) {
	var earlyTokenContract contract.EarlyRewardContract
	block := createBlock(20)

	// GetMineMasterCSKReward error
	processor, err := NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{DIPErr: EconomyErr}
	err = processor.RewardCoinBase(block, &earlyTokenContract)
	assert.Equal(t, EconomyErr, err)

	// AddBalance error
	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{setErr: TrieError})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{}
	err = processor.RewardCoinBase(block, &earlyTokenContract)
	assert.Equal(t, TrieError, err)

	// NewAccountState error
	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{
		getErr: TrieError,
		setErr: TrieError,
	})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{}
	err = processor.RewardCoinBase(block, &earlyTokenContract)
	assert.Equal(t, TrieError, err)

	// No error
	processor, err = NewBlockProcessor(fakeAccountDBChain{}, common.Hash{}, fakeStateStorage{})
	assert.NoError(t, err)
	processor.economyModel = fakeEconomyModel{}
	err = processor.RewardCoinBase(block, &earlyTokenContract)
	assert.NoError(t, err)

	// NoCoinBase error
	block = createBlockWithoutCoinBase()
	err = processor.RewardCoinBase(block, &earlyTokenContract)
	assert.Equal(t, g_error.InvalidCoinBaseAddressErr, err)
}
