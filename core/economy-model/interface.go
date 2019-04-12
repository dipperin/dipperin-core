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


package economy_model

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"math/big"
)

type VerifierType uint

const (
	MasterVerifier    VerifierType = iota
	CommitVerifier
	NotCommitVerifier
)

//go:generate mockgen -destination=./../verifiers-halt-check/economy_model_mock_test.go -package=verifiers_halt_check github.com/dipperin/dipperin-core/core/economy-model EconomyModel
type EconomyModel interface {
	GetMineMasterDIPReward(block model.AbstractBlock) (*big.Int, error)
	GetVerifierDIPReward(block model.AbstractBlock) (map[VerifierType]*big.Int, error)
	GetInvestorInitBalance() map[common.Address]*big.Int
	GetDeveloperInitBalance() map[common.Address]*big.Int

	GetInvestorLockDIP(address common.Address, blockNumber uint64) (*big.Int, error)
	GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*big.Int, error)

	GetFoundation() Foundation
	CheckAddressType(address common.Address) EconomyModelAddress

	GetDiffVerifierAddress(preBlock,block model.AbstractBlock) (map[VerifierType][]common.Address, error)

	GetAddressLockMoney(address common.Address,blockNumber uint64) (*big.Int,error)

	GetBlockYear(blockNumber uint64) (uint64,error)
	GetOneBlockTotalDIPReward(blockNumber uint64) (*big.Int, error)
}
