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
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"math/big"
)

// use outside mercury
// the proportion of each account while pre-mining
var NotMercuryDIPProportion = AddressDIPProportion{
	BaseNumber: 100,
	InvestorProportion: map[string]int{
		"0x000017cD9dB648440102E85371AB17871f62809FD183": 100,
	},
	DeveloperProportion: map[string]int{
		"0x00001305de8cfDa24E6aa2d5386D428E276A0785aF70": 100,
	},
	MaintenanceProportion: map[string]int{
		"0x00007737f9432D99fA9c5a052Dcb2D0b7217aA24cc17": 100,
	},
	EarlyTokenProportion: map[string]int{
		"0x0000D21AF412898F2F0A646d3686Adce8dAFBAfD045B": 100,
	},
	ReMainRewardProportion: map[string]int{
		"0x000016D54544B1eBCa3343847f6b0eE26509643F0Fd5": 100,
	},
}

// use inside mercury
var MercuryDIPProportion = AddressDIPProportion{
	BaseNumber: 100,
	InvestorProportion: map[string]int{
		"0x00000aa9bD540140d5c4E9c68685c360a091E46F4c27": 100,
	},
	DeveloperProportion: map[string]int{
		"0x0000125f454Cbf30F4354cD343f5cfb4C100D0ccCB34": 100,
	},
	MaintenanceProportion: map[string]int{
		"0x00005E3f1964Fad70Da7377BD53a76846E3bEFb734eA": 100,
	},
	EarlyTokenProportion: map[string]int{
		"0x0000960854D56506697b387463DB1C601eAb14be763E": 100,
	},
	ReMainRewardProportion: map[string]int{
		"0x0000c2FCE2Cf78953BAd21890215d4e2Ae823b56A73e": 100,
	},
}

// system parameters
var (
	TotalReWardDIPOneBlock = big.NewInt(0).Mul(big.NewInt(20), big.NewInt(consts.DIP))

	// the reward proportions of mineMaster and verifiers respectively for one block
	MineMasterRewardProportion = int64(87) //87%
	VerifierRewardProportion   = int64(13) //13%

	// the reward weight for each type of verifier
	MainVerifierRewardWeight     = 2
	CommitVerifierRewardWeight   = 4
	NoCommitVerifierRewardWeight = 1

	// to be adjusted in 10 years
	ChangeIssuingYear = uint64(10)

	// the issuing rate 10 years after,
	// the issuance amount during the nth year = the issuance amount * 3%
	IssuingRate = 3

	// block interval is 8s, namely 8 seconds one block
	// todo  there is a problem that the block generate interval is not equal the chainconfig.DefaultBlockGenerateInterval, should use the chainconfig.DefaultBlockGenerateInterval
	GenerateBlockDuration = 8

	// set the minimal deposit to register as a verifier
	MiniPledgeValue = big.NewInt(0).Mul(big.NewInt(1000), big.NewInt(consts.DIP))

	// the total supply of pre-mining
	PreMineDIP = big.NewInt(0).Mul(big.NewInt(525600000), big.NewInt(consts.DIP))

	// the supply of pre-mining for investors, developers and maintenance operators
	PreMineMineDIP = big.NewInt(0).Mul(big.NewInt(438000000), big.NewInt(consts.DIP))

	// the supply of pre-mining of each group
	InvestorProportion    = int64(60) //60%
	DevelopProportion     = int64(20) //20%
	MaintenanceProportion = int64(20) //20%

	// the supply of pre-mining for reward
	PreMineRewardDIP = big.NewInt(0).Mul(big.NewInt(87600000), big.NewInt(consts.DIP))

	// the supply of pre-mining for reward for each group, respectively 50%
	EarlyTokenReward  = int64(50)
	RemainOtherReward = int64(50)

	// 5 years
	EarlyTokenPeriod = uint64(5) //5

	// unlocking proportion each year for pre-mining investors
	InvestorUnlockInfo = map[int]int64{
		1: 10,
		2: 20,
		3: 25,
		4: 25,
		5: 20,
	}

	// unlocking proportion each year for pre-mining developers
	DeveloperUnlockInfo = map[int]int64{
		1: 10,
		2: 20,
		3: 20,
		4: 20,
		5: 20,
		6: 10,
	}

	// unlocking times each year for pre-mining, 4 means unlocking each quarter
	UnlockTimeOneYear = uint64(4)

	// proportion of pre-mining for each address
	DIPProportion AddressDIPProportion

	// the total early token supply measured by token unit
	EarlyTokenAmount = big.NewInt(1340280000)

	// the conversion rate of early token with respect to DIP
	EarlyTokenExchangeBase = int64(10000)
)

// unlocking proportion each year for pre-mining investors
var (
	// the number of blocks to be generated during 1 year
	HeightAfterOneYear = uint64(365 * 24 * 3600 / GenerateBlockDuration)
	// the number of blocks to be generated during 10 years, with bias of course due to the non fixed interval for block generation
	HeightAfterTenYear = 10 * HeightAfterOneYear

	// total amount possessed by investors
	InvestorDIP = big.NewInt(0).Div(big.NewInt(0).Mul(PreMineMineDIP, big.NewInt(InvestorProportion)), big.NewInt(100))
	// total amount possessed by developers
	DeveloperDIP = big.NewInt(0).Div(big.NewInt(0).Mul(PreMineMineDIP, big.NewInt(DevelopProportion)), big.NewInt(100))
	// total amount possessed by maintenance operators
	MaintenanceDIP = big.NewInt(0).Div(big.NewInt(0).Mul(PreMineMineDIP, big.NewInt(MaintenanceProportion)), big.NewInt(100))
	// total amount possessed by earlyToken
	EarlyTokenDIP = big.NewInt(0).Div(big.NewInt(0).Mul(PreMineRewardDIP, big.NewInt(EarlyTokenReward)), big.NewInt(100))
	// remaining coin rewards reserved
	RemainRewardDIP = big.NewInt(0).Div(big.NewInt(0).Mul(PreMineRewardDIP, big.NewInt(RemainOtherReward)), big.NewInt(100))

	InvestorAddresses     = make([]common.Address, 0)
	DeveloperAddresses    = make([]common.Address, 0)
	MaintenanceAddresses  = make([]common.Address, 0)
	EarlyTokenAddresses   = make([]common.Address, 0)
	RemainRewardAddresses = make([]common.Address, 0)
)

var InitExchangeRate int64

func init() {

	if chainconfig.GetCurBootsEnv() != chainconfig.BootEnvMercury {
		DIPProportion = NotMercuryDIPProportion
	} else {
		DIPProportion = MercuryDIPProportion
	}

	//log.DLogger.Info("the EarlyTokenDIP is:", "EarlyTokenDIP", EarlyTokenDIP)
	exchangeRate := big.NewInt(0).Mul(EarlyTokenDIP, big.NewInt(EarlyTokenExchangeBase))
	exchangeRate.Div(exchangeRate, big.NewInt(consts.DIP))
	exchangeRate.Div(exchangeRate, EarlyTokenAmount)
	InitExchangeRate = exchangeRate.Int64()

	for address := range DIPProportion.InvestorProportion {
		InvestorAddresses = append(InvestorAddresses, common.HexToAddress(address))
	}

	for address := range DIPProportion.DeveloperProportion {
		DeveloperAddresses = append(DeveloperAddresses, common.HexToAddress(address))
	}

	for address := range DIPProportion.MaintenanceProportion {
		MaintenanceAddresses = append(MaintenanceAddresses, common.HexToAddress(address))
	}

	for address := range DIPProportion.EarlyTokenProportion {
		EarlyTokenAddresses = append(EarlyTokenAddresses, common.HexToAddress(address))
	}

	for address := range DIPProportion.ReMainRewardProportion {
		RemainRewardAddresses = append(RemainRewardAddresses, common.HexToAddress(address))
	}
}

type PreMineMainType int

const (
	Investor PreMineMainType = iota
	Developer
)

type EconomyModelAddress int

const (
	InvestorAddress EconomyModelAddress = iota
	DeveloperAddress
	NotEconomyModelAddress
)

//go:generate mockgen -destination=./economy_need_service_mock.go -package=economymodel github.com/dipperin/dipperin-core/core/economymodel EconomyNeedService
type EconomyNeedService interface {
	GetVerifiers(slotNum uint64) (addresses []common.Address)
	GetSlot(block model.AbstractBlock) *uint64
}

type DipperinEconomyModel struct {
	Service EconomyNeedService

	Foundation
	// the corresponding initial value for each pre-mining address
	investInitBalance    map[common.Address]*big.Int
	developerInitBalance map[common.Address]*big.Int
}

// configuration of DIP proportion of each address
type AddressDIPProportion struct {
	BaseNumber          int64          `json:"base_number"`
	InvestorProportion  map[string]int `json:"investor_proportion"`
	DeveloperProportion map[string]int `json:"developer_proportion"`

	// maintenance operator proportion
	MaintenanceProportion map[string]int `json:"maintenance_proportion"`

	EarlyTokenProportion map[string]int `json:"early_token_proportion"`

	// other user reward proportion preserved
	ReMainRewardProportion map[string]int `json:"remain_reward_proportion"`
}

func MapMerge(des, src map[common.Address]*big.Int) error {
	for address, value := range src {
		if _, ok := des[address]; ok {
			log.DLogger.Error("address already in des", zap.String("addr", address.Hex()))
			return gerror.ErrAddressExist
		}
		des[address] = value
	}

	return nil
}

func paddingAddressDIP(totalDIP *big.Int, proportion map[string]int, baseNumber int64) map[common.Address]*big.Int {
	result := make(map[common.Address]*big.Int, 0)
	for address, value := range proportion {
		tmp := big.NewInt(0).Mul(totalDIP, big.NewInt(int64(value)))
		tmp.Div(tmp, big.NewInt(baseNumber))
		result[common.HexToAddress(address)] = tmp
	}
	return result
}

func MakeDipperinEconomyModel(service EconomyNeedService, proportion AddressDIPProportion) *DipperinEconomyModel {
	economyModel := &DipperinEconomyModel{
		Service: service,
	}

	economyModel.investInitBalance = paddingAddressDIP(InvestorDIP, proportion.InvestorProportion, proportion.BaseNumber)
	economyModel.developerInitBalance = paddingAddressDIP(DeveloperDIP, proportion.DeveloperProportion, proportion.BaseNumber)
	economyModel.Foundation = MakeDipperinFoundation(proportion)

	return economyModel
}

// calculate the total supply during the first n years, pre-mining not included
func CalcDIPTotalCirculation(year uint64) (value *big.Int) {
	if year == 1 {
		return big.NewInt(0).Mul(big.NewInt(int64(HeightAfterOneYear)), TotalReWardDIPOneBlock)
	} else if year > 1 && year <= ChangeIssuingYear {
		return big.NewInt(0).Mul(CalcDIPTotalCirculation(1), big.NewInt(int64(year)))
	} else if year > ChangeIssuingYear {
		tmp := CalcDIPTotalCirculation(year - 1)
		tmp.Add(tmp, PreMineDIP)
		currentYearDIP := big.NewInt(0).Div(big.NewInt(0).Mul(tmp, big.NewInt(int64(IssuingRate))), big.NewInt(100))
		return big.NewInt(0).Add(tmp, currentYearDIP)
	} else {
		return big.NewInt(0)
	}
}

func (economyModel *DipperinEconomyModel) GetOneBlockTotalDIPReward(blockNumber uint64) (*big.Int, error) {
	rewardOneBlock := big.NewInt(0)
	if blockNumber == 0 {
		return big.NewInt(0), gerror.ErrBlockNumberIs0
	} else if blockNumber <= HeightAfterTenYear {
		rewardOneBlock = TotalReWardDIPOneBlock
	} else {
		tmp := CalcDIPTotalCirculation((blockNumber+HeightAfterOneYear-1)/HeightAfterOneYear - 1)
		tmp.Add(tmp, PreMineDIP)
		currentYearDIP := big.NewInt(0).Div(big.NewInt(0).Mul(tmp, big.NewInt(int64(IssuingRate))), big.NewInt(100))
		rewardOneBlock = rewardOneBlock.Div(currentYearDIP, big.NewInt(int64(HeightAfterOneYear)))
	}
	return rewardOneBlock, nil
}

//calculate different verifier reward
func (economyModel *DipperinEconomyModel) calcDifferentVerifierReward(totalReward *big.Int) map[VerifierType]*big.Int {
	conf := chainconfig.GetChainConfig()
	commitNumber := conf.VerifierNumber*2/3 + 1
	notCommitNumber := conf.VerifierNumber - commitNumber
	mineVerifierNumber := 1

	totalWeight := mineVerifierNumber*MainVerifierRewardWeight + commitNumber*CommitVerifierRewardWeight + notCommitNumber*NoCommitVerifierRewardWeight
	verifierReward := make(map[VerifierType]*big.Int, 0)
	verifierReward[MasterVerifier] = big.NewInt(0).Div(big.NewInt(0).Mul(totalReward, big.NewInt(int64(MainVerifierRewardWeight))), big.NewInt(int64(totalWeight)))
	verifierReward[CommitVerifier] = big.NewInt(0).Div(big.NewInt(0).Mul(totalReward, big.NewInt(int64(CommitVerifierRewardWeight))), big.NewInt(int64(totalWeight)))
	verifierReward[NotCommitVerifier] = big.NewInt(0).Div(big.NewInt(0).Mul(totalReward, big.NewInt(int64(NoCommitVerifierRewardWeight))), big.NewInt(int64(totalWeight)))

	return verifierReward
}


// todo  this method seems to be wrong
// calculate the amount of locked DIP for investors and developers for different height
func (economyModel *DipperinEconomyModel) calcLockDIP(unlockType PreMineMainType, address common.Address, blockNumber uint64) (*big.Int, error) {
	unlockTotalDIP := make(map[common.Address]*big.Int, 0)
	unlockInfo := make(map[int]int64, 0)
	if unlockType == Investor {
		unlockTotalDIP = economyModel.investInitBalance
		unlockInfo = InvestorUnlockInfo
	} else if unlockType == Developer {
		unlockTotalDIP = economyModel.developerInitBalance
		unlockInfo = DeveloperUnlockInfo
	} else {
		return big.NewInt(0), gerror.ErrLockTypeError
	}

	if totalDIP, ok := unlockTotalDIP[address]; ok {
		if blockNumber == 0 {
			return big.NewInt(0), gerror.ErrBlockNumberIs0
		}

		// locking period is bypassed
		year := (blockNumber + HeightAfterOneYear - 1) / HeightAfterOneYear
		if _, ok = unlockInfo[int(year)]; !ok {
			return big.NewInt(0), nil
		}
		unlockMoney := big.NewInt(0)
		for i := 1; i <= int(year); i++ {
			tmp := big.NewInt(0).Div(big.NewInt(0).Mul(totalDIP, big.NewInt(unlockInfo[i])), big.NewInt(100))
			if i == int(year) {
				tmpNumber := blockNumber - (year-1)*HeightAfterOneYear
				blockNumberOneUnlock := HeightAfterOneYear / UnlockTimeOneYear
				unlockProportion := (tmpNumber + blockNumberOneUnlock - 1) / blockNumberOneUnlock
				tmp.Div(tmp.Mul(tmp, big.NewInt(int64(unlockProportion))), big.NewInt(int64(UnlockTimeOneYear)))
			}
			unlockMoney.Add(unlockMoney, tmp)
		}

		return big.NewInt(0).Sub(totalDIP, unlockMoney), nil
	} else {
		return big.NewInt(0), gerror.ErrAddress
	}
}

//get foundation
func (economyModel *DipperinEconomyModel) GetFoundation() Foundation {
	return economyModel.Foundation
}

// get coin reward of mineMaster for each block
func (economyModel *DipperinEconomyModel) GetMineMasterDIPReward(block model.AbstractBlock) (*big.Int, error) {
	totalReward, err := economyModel.GetOneBlockTotalDIPReward(block.Number())
	if err != nil {
		return big.NewInt(0), err
	}
	//log.DLogger.Info("onBlock reward is:", "cur block total reward", totalReward)
	return big.NewInt(0).Div(big.NewInt(0).Mul(totalReward, big.NewInt(MineMasterRewardProportion)), big.NewInt(100)), nil
}

// get corresponding reward for different type of verifiers for each block
func (economyModel *DipperinEconomyModel) GetVerifierDIPReward(block model.AbstractBlock) (map[VerifierType]*big.Int, error) {
	// no reward for genesis
	if block.Number() == 0 {
		return map[VerifierType]*big.Int{}, gerror.ErrBlockNumberIs0
	}

	totalReward, err := economyModel.GetOneBlockTotalDIPReward(block.Number())
	if err != nil {
		return map[VerifierType]*big.Int{}, err
	}

	verifierReward := big.NewInt(0).Div(big.NewInt(0).Mul(totalReward, big.NewInt(VerifierRewardProportion)), big.NewInt(100))
	return economyModel.calcDifferentVerifierReward(verifierReward), nil
}

// get the address of different verifier type for each block
func (economyModel *DipperinEconomyModel) GetDiffVerifierAddress(preBlock, block model.AbstractBlock) (map[VerifierType][]common.Address, error) {

	// genesis and block 1 has no commit information because the commit information for current block is stored in the subsequent block
	if block.Number() < 2 {
		return map[VerifierType][]common.Address{}, gerror.ErrBlockNumberIs0Ore1
	}

	config := chainconfig.GetChainConfig()
	//log.DLogger.Info("get address", "economyModel", economyModel)
	slot := economyModel.Service.GetSlot(preBlock)
	log.DLogger.Debug("the slot is:", zap.Uint64("slot", *slot))
	verifiers := economyModel.Service.GetVerifiers(*slot)
	log.DLogger.Debug("get verifiers is:", zap.Any("verifier", verifiers))

	verifierAddress := make(map[VerifierType][]common.Address, 0)
	commitVerifier := make([]common.Address, 0)
	notCommitVerifier := make([]common.Address, len(verifiers))
	copy(notCommitVerifier, verifiers)

	verifications := block.GetVerifications()
	//log.DLogger.Info("the verifications number is:","number",len(verifications))
	for _, verification := range verifications {
		//log.DLogger.Info("the verification address is:","address",verification.GetAddress().Hex())
		commitVerifier = append(commitVerifier, verification.GetAddress())
		for i, tmpAddress := range notCommitVerifier {
			if verification.GetAddress() == tmpAddress {
				notCommitVerifier = append(notCommitVerifier[:i], notCommitVerifier[i+1:]...)
			}
		}
	}

	masterVerifierIndex := int(verifications[0].GetRound()) % config.VerifierNumber
	log.DLogger.Info("the masterVerifierIndex is:", zap.Int("index", masterVerifierIndex))
	log.DLogger.Info("the verifierAddress is:", zap.Int("number", len(verifiers)))
	verifierAddress[MasterVerifier] = []common.Address{verifiers[masterVerifierIndex]}
	verifierAddress[CommitVerifier] = commitVerifier
	verifierAddress[NotCommitVerifier] = notCommitVerifier

	return verifierAddress, nil
}

// get the minimum transaction fee according to the size of the transaction
// minimumTxFee = txSize * 0.0000001 * const.DIP
//func GetMinimumTxFee(txSize common.StorageSize) *big.Int {
//	return big.NewInt(0).Mul(big.NewInt(int64(txSize)), big.NewInt(100))
//}

// get the pre-mining DIP amount for each address of investors
func (economyModel *DipperinEconomyModel) GetInvestorInitBalance() map[common.Address]*big.Int {
	return economyModel.investInitBalance
}

// get the pre-mining DIP amount for each address of developers
func (economyModel *DipperinEconomyModel) GetDeveloperInitBalance() map[common.Address]*big.Int {
	return economyModel.developerInitBalance
}

// lock the pre-mining DIP amount for each address of investors
func (economyModel *DipperinEconomyModel) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return economyModel.calcLockDIP(Investor, address, blockNumber)
}

// lock the pre-mining DIP amount for each address of developers
func (economyModel *DipperinEconomyModel) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*big.Int, error) {
	return economyModel.calcLockDIP(Developer, address, blockNumber)
}

// identify the type of address
func (economyModel *DipperinEconomyModel) CheckAddressType(address common.Address) EconomyModelAddress {
	if _, ok := economyModel.investInitBalance[address]; ok {
		return InvestorAddress
	}

	if _, ok := economyModel.developerInitBalance[address]; ok {
		return DeveloperAddress
	}

	return NotEconomyModelAddress
}

// identify the type of address and return the value of locked coins
func (economyModel *DipperinEconomyModel) GetAddressLockMoney(address common.Address, blockNumber uint64) (*big.Int, error) {
	lockValue := big.NewInt(0)
	economyType := economyModel.CheckAddressType(address)

	switch economyType {
	case InvestorAddress:
		result, err := economyModel.GetInvestorLockDIP(address, blockNumber)
		if err != nil {
			return big.NewInt(0), err
		}
		lockValue.Add(lockValue, result)
	case DeveloperAddress:
		result, err := economyModel.GetDeveloperLockDIP(address, blockNumber)
		if err != nil {
			return big.NewInt(0), err
		}
		lockValue.Add(lockValue, result)
	}

	foundationType := economyModel.Foundation.GetAddressType(address)
	switch foundationType {
	case MaintenanceAddress:
		result, err := economyModel.Foundation.GetMaintenanceLockDIP(address, blockNumber)
		if err != nil {
			return big.NewInt(0), err
		}
		lockValue.Add(lockValue, result)
	case EarlyTokenAddress:
		result, err := economyModel.Foundation.GetEarlyTokenLockDIP(address, blockNumber)
		if err != nil {
			return big.NewInt(0), err
		}
		lockValue.Add(lockValue, result)
	case RemainRewardAddress:
		result, err := economyModel.Foundation.GetReMainRewardLockDIP(address, blockNumber)
		if err != nil {
			return big.NewInt(0), err
		}
		lockValue.Add(lockValue, result)
	}

	return lockValue, nil
}

func (economyModel *DipperinEconomyModel) GetBlockYear(blockNumber uint64) (uint64, error) {
	year := (blockNumber + HeightAfterOneYear - 1) / HeightAfterOneYear
	return year, nil
}
