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

package economymodel

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"math/big"
)

type Foundation interface {
	GetAddressType(address common.Address) FoundationAddressType
	GetFoundationInfo(usage FoundationDIPUsage) map[common.Address]*big.Int
	GetMaintenanceLockDIP(address common.Address, blockNumber uint64) (*big.Int, error)
	GetReMainRewardLockDIP(address common.Address, blockNumber uint64) (*big.Int, error)
	GetEarlyTokenLockDIP(address common.Address, blockNumber uint64) (*big.Int, error)
	// set the conversion rate of EDIP w.r.t DIP
	GetMineMasterEDIPReward(DIPReward *big.Int, blockNumber uint64, tokenDecimals int) (*big.Int, error)
	GetVerifierEDIPReward(DIPReward map[VerifierType]*big.Int, blockNumber uint64, tokenDecimals int) (map[VerifierType]*big.Int, error)
}

type FoundationDIPUsage int

const (
	EarlyToken FoundationDIPUsage = iota
	Maintenance
	RemainReward
)

type FoundationAddressType int

const (
	EarlyTokenAddress FoundationAddressType = iota
	MaintenanceAddress
	RemainRewardAddress
	NotFoundationAddress
)

type DipperinFoundation struct {
	//foundation address initial DIP
	foundationInfo map[FoundationDIPUsage]map[common.Address]*big.Int
}

func MakeDipperinFoundation(proportion AddressDIPProportion) *DipperinFoundation {
	foundation := &DipperinFoundation{
		foundationInfo: make(map[FoundationDIPUsage]map[common.Address]*big.Int, 0),
	}
	foundation.foundationInfo[EarlyToken] = paddingAddressDIP(EarlyTokenDIP, proportion.EarlyTokenProportion, proportion.BaseNumber)
	foundation.foundationInfo[Maintenance] = paddingAddressDIP(MaintenanceDIP, proportion.MaintenanceProportion, proportion.BaseNumber)
	foundation.foundationInfo[RemainReward] = paddingAddressDIP(RemainRewardDIP, proportion.ReMainRewardProportion, proportion.BaseNumber)
	return foundation
}

//y = -0.0046x^2 + 0.0505x - 0.0009
/*func calcReMainRewardLock(totalValue *big.Int, blockNumber uint64) (*big.Int, error) {
	if blockNumber == 0 {
		return big.NewInt(0), ErrBlockNumberIs0
	}
	year := (blockNumber + HeightAfterOneYear - 1) / HeightAfterOneYear
	if year > ChangeIssuingYear {
		return big.NewInt(0), nil
	}

	UnlockProportion := big.NewInt(0)
	UnlockMoney := big.NewInt(0)
	for i := 1; i <= int(year); i++ {
		x := big.NewInt(int64(i))
		tmp := big.NewInt(0)
		UnlockProportion.Mul(big.NewInt(505), x)
		UnlockProportion.Sub(UnlockProportion, big.NewInt(9))
		UnlockProportion.Sub(UnlockProportion, big.NewInt(0).Mul(big.NewInt(46), x.Exp(x, big.NewInt(2), nil)))

		tmp.Mul(totalValue, UnlockProportion)
		tmp.Div(UnlockMoney, big.NewInt(10000))

		UnlockMoney.Add(UnlockMoney, tmp)
	}
	return big.NewInt(0).Sub(totalValue, UnlockMoney), nil
}*/

// calculate early token reward
// Reward(n) = 5-0.4*∑(t-1) t in 1~n
func calcEarlyTokenValue(DIPReward *big.Int, blockNumber uint64, tokenDecimals int) (*big.Int, error) {
	if blockNumber == 0 {
		return big.NewInt(0), ErrBlockNumberIs0
	}
	year := (blockNumber + HeightAfterOneYear - 1) / HeightAfterOneYear
	if year > EarlyTokenPeriod {
		return big.NewInt(0), nil
	}

	sum := 0
	for i := 1; i <= int(year); i++ {
		sum += i - 1
	}
	reward := big.NewInt(0).Mul(big.NewInt(int64(50-4*sum)), DIPReward)
	reward.Mul(reward, big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(tokenDecimals)), nil))
	reward.Div(reward, big.NewInt(10))
	reward.Div(reward, big.NewInt(consts.DIP))

	return reward, nil
}

// get type of address
func (foundation *DipperinFoundation) GetAddressType(address common.Address) FoundationAddressType {
	if _, ok := foundation.foundationInfo[EarlyToken][address]; ok {
		return EarlyTokenAddress
	} else if _, ok := foundation.foundationInfo[Maintenance][address]; ok {
		return MaintenanceAddress
	} else if _, ok := foundation.foundationInfo[RemainReward][address]; ok {
		return RemainRewardAddress
	} else {
		return NotFoundationAddress
	}
}

// get foundation address initial DIP
func (foundation *DipperinFoundation) GetFoundationInfo(usage FoundationDIPUsage) map[common.Address]*big.Int {
	return foundation.foundationInfo[usage]
}

// get foundation Maintenance each address Lock DIP, maintenance not locked
func (foundation *DipperinFoundation) GetMaintenanceLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return big.NewInt(0), nil
}

//get foundation Maintenance each address Lock DIP, RemainReward not locked
func (foundation *DipperinFoundation) GetReMainRewardLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return big.NewInt(0), nil
	//return calcReMainRewardLock(foundation.foundationInfo[RemainReward][address],blockNumber)
}

//get foundation EarlyToken each address unlock DIP
func (foundation *DipperinFoundation) GetEarlyTokenLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return big.NewInt(0), nil
}

// get mineMaster EDIPReward
func (foundation *DipperinFoundation) GetMineMasterEDIPReward(DIPReward *big.Int, blockNumber uint64, tokenDecimals int) (*big.Int, error) {
	return calcEarlyTokenValue(DIPReward, blockNumber, tokenDecimals)
}

// get verifier EDIPReward
func (foundation *DipperinFoundation) GetVerifierEDIPReward(DIPReward map[VerifierType]*big.Int, blockNumber uint64, tokenDecimals int) (map[VerifierType]*big.Int, error) {
	result := make(map[VerifierType]*big.Int, 0)
	for verifierType, value := range DIPReward {
		tokenValue, err := calcEarlyTokenValue(value, blockNumber, tokenDecimals)
		if err != nil {
			return map[VerifierType]*big.Int{}, err
		}
		result[verifierType] = tokenValue
	}
	return result, nil
}
