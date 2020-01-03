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

package economy_model_test

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"math/big"
	"testing"
)

type testService struct {
}

func (*testService) GetVerifiers(slotNum uint64) (addresses []common.Address) {
	return chain.VerifierAddress
}

func (*testService) GetSlot(block model.AbstractBlock) *uint64 {
	slotSize := chain_config.GetChainConfig().SlotSize
	slot := block.Number() / slotSize
	return &slot
}

var testEconomyService = &testService{}

func TestDipperinEconomyModel_MapMerge(t *testing.T) {
	src := map[common.Address]*big.Int{}
	des := map[common.Address]*big.Int{}

	src[common.HexToAddress("1234")] = big.NewInt(1)
	err := economy_model.MapMerge(des, src)
	assert.NoError(t, err)

	//existed error
	err = economy_model.MapMerge(des, src)
	assert.EqualError(t, err, economy_model.ErrAddressExist.Error())
}

// test the reward for a single block (macroeconomics)
func TestDipperinEconomyModel_GetOneBlockTotalDIPReward(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)

	//block num=0
	_, err := economyModel.GetOneBlockTotalDIPReward(0)
	assert.EqualError(t, err, economy_model.ErrBlockNumberIs0.Error())

	//block num <= HeightAfterTenYear
	v, _ := economyModel.GetOneBlockTotalDIPReward(economy_model.HeightAfterTenYear - 1)
	assert.Equal(t, v, economy_model.TotalReWardDIPOneBlock)

	v, _ = economyModel.GetOneBlockTotalDIPReward(economy_model.HeightAfterTenYear)
	assert.Equal(t, v, economy_model.TotalReWardDIPOneBlock)

	//block num > HeightAfterTenYear
	v, _ = economyModel.GetOneBlockTotalDIPReward(economy_model.HeightAfterTenYear + 1)
	assert.NotEqual(t, v, economy_model.TotalReWardDIPOneBlock)
}

// test the foundation
func TestDipperinEconomyModel_GetFoundation(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	f := economyModel.GetFoundation()
	assert.NotNil(t, f)
}

// test the minimum transaction fee
func TestDipperinEconomyModel_GetMinimumTxFee(t *testing.T) {
	fee := economy_model.GetMinimumTxFee(30)
	assert.NotNil(t, fee)
}

// test address type
func TestDipperinEconomyModel_CheckAddressType(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	infoInvest := economyModel.GetInvestorInitBalance()

	infoDev := economyModel.GetDeveloperInitBalance()

	for k := range infoInvest {
		assert.Equal(t, economyModel.CheckAddressType(k), economy_model.InvestorAddress)
	}

	for k := range infoDev {
		assert.Equal(t, economyModel.CheckAddressType(k), economy_model.DeveloperAddress)
	}
	assert.Equal(t, economyModel.CheckAddressType(common.HexToAddress("0x1234")), economy_model.NotEconomyModelAddress)
}

// test unlocking mechanism
func TestDipperinEconomyModel_GetAddressLockMoney(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)

	testData := make([]map[common.Address]*big.Int, 0)
	infoInvest := economyModel.GetInvestorInitBalance()
	testData = append(testData, infoInvest)
	infoDev := economyModel.GetDeveloperInitBalance()
	testData = append(testData, infoDev)
	infoMaintenance := economyModel.GetFoundationInfo(economy_model.Maintenance)
	testData = append(testData, infoMaintenance)
	infoEarlyToken := economyModel.GetFoundationInfo(economy_model.EarlyToken)
	testData = append(testData, infoEarlyToken)
	infoRemainReward := economyModel.GetFoundationInfo(economy_model.RemainReward)
	testData = append(testData, infoRemainReward)

	for i, info := range testData {
		for k := range info {
			if i < 2 {
				//test investor and developer
				_, err := economyModel.GetAddressLockMoney(k, 0)
				assert.Error(t, err, economy_model.ErrBlockNumberIs0.Error())
			}
			_, err := economyModel.GetAddressLockMoney(k, 1)
			assert.NoError(t, err)
		}
	}

	v, _ := economyModel.GetAddressLockMoney(common.HexToAddress("0x1234"), 1)
	assert.Equal(t, v, big.NewInt(0))
}

// test of getting block number by the number of year
func TestDipperinEconomyModel_GetBlockYear(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	year, _ := economyModel.GetBlockYear(0)
	assert.Equal(t, year, uint64(0))
}

// test the reward for miners
func TestDipperinEconomyModel_GetMineMasterDIPReward(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)

	controller := gomock.NewController(t)
	defer controller.Finish()
	testBlock := economy_model.NewMockAbstractBlock(controller)
	testBlock.EXPECT().Number().Return(uint64(30))

	reward, err := economyModel.GetMineMasterDIPReward(testBlock)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0).Mul(big.NewInt(17400000000), big.NewInt(consts.GDIPUNIT)), reward)
	testBlock.EXPECT().Number().Return(uint64(economy_model.HeightAfterTenYear + 1))
	reward, err = economyModel.GetMineMasterDIPReward(testBlock)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0).Mul(big.NewInt(8700000000), big.NewInt(consts.GDIPUNIT)), reward)

	testBlock.EXPECT().Number().Return(uint64(11*economy_model.HeightAfterOneYear) + 1)
	reward, err = economyModel.GetMineMasterDIPReward(testBlock)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(0).Mul(big.NewInt(12441000000), big.NewInt(consts.GDIPUNIT)), reward)

	testBlock.EXPECT().Number().Return(uint64(0))
	_, err = economyModel.GetMineMasterDIPReward(testBlock)
	assert.Equal(t, economy_model.ErrBlockNumberIs0, err)
}

// test the reward for verifiers
func TestDipperinEconomyModel_GetVerifierDIPReward(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockBlock1 := economy_model.NewMockAbstractBlock(controller)
	mockBlock1.EXPECT().Number().Return(uint64(30)).AnyTimes()
	reward, err := economyModel.GetVerifierDIPReward(mockBlock1)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(75362318840579710), reward[economy_model.MasterVerifier])
	assert.EqualValues(t, big.NewInt(150724637681159420), reward[economy_model.CommitVerifier])
	assert.EqualValues(t, big.NewInt(37681159420289855), reward[economy_model.NotCommitVerifier])

	mockBlock2 := economy_model.NewMockAbstractBlock(controller)
	mockBlock2.EXPECT().Number().Return(uint64(economy_model.HeightAfterTenYear + 1)).AnyTimes()
	reward, err = economyModel.GetVerifierDIPReward(mockBlock2)
	assert.NoError(t, err)
	assert.EqualValues(t, big.NewInt(37681159420289855), reward[economy_model.MasterVerifier])
	assert.EqualValues(t, big.NewInt(75362318840579710), reward[economy_model.CommitVerifier])
	assert.EqualValues(t, big.NewInt(18840579710144927), reward[economy_model.NotCommitVerifier])

	mockBlock3 := economy_model.NewMockAbstractBlock(controller)
	mockBlock3.EXPECT().Number().Return(uint64(0)).AnyTimes()
	_, err = economyModel.GetVerifierDIPReward(mockBlock3)
	assert.Equal(t, economy_model.ErrBlockNumberIs0, err)
}

// test investors
func TestDipperinEconomyModel_GetInvestorInitBalance(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	info := economyModel.GetInvestorInitBalance()
	assert.EqualValues(t, big.NewInt(0).Mul(big.NewInt(262800000000000000), big.NewInt(consts.GDIPUNIT)), info[economy_model.InvestorAddresses[0]])
}

// test developers
func TestDipperinEconomyModel_GetDeveloperInitBalance(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	info := economyModel.GetDeveloperInitBalance()
	assert.EqualValues(t, big.NewInt(0).Mul(big.NewInt(87600000000000000), big.NewInt(consts.GDIPUNIT)), info[economy_model.DeveloperAddresses[0]])
}

// test investor unlocking mechanism
func TestDipperinEconomyModel_GetInvestorLockDIP(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	investorTotalDIP := big.NewInt(0).Mul(big.NewInt(262800000000000000), big.NewInt(consts.GDIPUNIT))

	// unlocking by quarters each year
	for i := 1; i <= 4; i++ {
		blockNumber := economy_model.HeightAfterOneYear / 4 * uint64(i)
		lockDIP, err := economyModel.GetInvestorLockDIP(economy_model.InvestorAddresses[0], blockNumber)
		assert.NoError(t, err)

		unlockDIP := big.NewInt(0).Div(investorTotalDIP, big.NewInt(10))
		unlockDIP.Mul(unlockDIP, big.NewInt(int64(i)))
		unlockDIP.Div(unlockDIP, big.NewInt(4))
		lockValue := big.NewInt(0).Sub(investorTotalDIP, unlockDIP)

		assert.EqualValues(t, lockValue, lockDIP)
	}

	_, err := economyModel.GetInvestorLockDIP(common.Address{}, 2)
	assert.Equal(t, economy_model.ErrAddress, err)
}

// test developer unlocking mechanism
func TestDipperinEconomyModel_GetDeveloperLockDIP(t *testing.T) {
	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)
	investorTotalDIP := big.NewInt(0).Mul(big.NewInt(87600000000000000), big.NewInt(consts.GDIPUNIT))

	// unlocking by quarters each year
	for i := 1; i <= 4; i++ {
		blockNumber := economy_model.HeightAfterOneYear / 4 * uint64(i)
		lockDIP, err := economyModel.GetDeveloperLockDIP(economy_model.DeveloperAddresses[0], blockNumber)
		assert.NoError(t, err)

		unlockDIP := big.NewInt(0).Div(investorTotalDIP, big.NewInt(10))
		unlockDIP.Mul(unlockDIP, big.NewInt(int64(i)))
		unlockDIP.Div(unlockDIP, big.NewInt(4))
		lockValue := big.NewInt(0).Sub(investorTotalDIP, unlockDIP)

		assert.EqualValues(t, lockValue, lockDIP)
	}

	_, err := economyModel.GetDeveloperLockDIP(common.Address{}, 2)
	assert.Equal(t, economy_model.ErrAddress, err)
}

func TestGenerateAddress(t *testing.T) {
	addressNumber := 5
	for i := 0; i < addressNumber; i++ {
		sk, err := crypto.GenerateKey()
		assert.NoError(t, err)
		address := cs_crypto.GetNormalAddress(sk.PublicKey)
		log.DLogger.Info("the address is:", zap.String("address", address.Hex()))
	}
}

func TestCalcDIPTotalCirculation(t *testing.T) {
	value := economy_model.CalcDIPTotalCirculation(0)
	assert.Equal(t, big.NewInt(0), value)

	value = economy_model.CalcDIPTotalCirculation(1)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(78840000), big.NewInt(consts.DIP)), value)

	value = economy_model.CalcDIPTotalCirculation(5)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(394200000), big.NewInt(consts.DIP)), value)

	value = economy_model.CalcDIPTotalCirculation(15)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(3788167915366200000), big.NewInt(consts.GDIPUNIT)), value)
}

func TestDipperinEconomyModel_GetDiffVerifierAddress(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockPreBlock := economy_model.NewMockAbstractBlock(controller)
	mockBlock := economy_model.NewMockAbstractBlock(controller)
	mockBlock.EXPECT().Number().Return(uint64(1))

	economyModel := economy_model.MakeDipperinEconomyModel(testEconomyService, economy_model.DIPProportion)

	_, err := economyModel.GetDiffVerifierAddress(mockPreBlock, mockBlock)
	assert.Equal(t, economy_model.ErrBlockNumberIs0Ore1, err)

	mockPreBlock.EXPECT().Number().Return(uint64(2))
	mockBlock.EXPECT().Number().Return(uint64(3))

	mockVerifier := economy_model.NewMockAbstractVerification(controller)
	mockVerifier.EXPECT().GetRound().Return(uint64(0))
	mockVerifier.EXPECT().GetAddress().Return(common.HexToAddress("0x000078b33598Be2b405206F44B018557e6F851FD230C")).AnyTimes()

	mockBlock.EXPECT().GetVerifications().Return(model.Verifications{mockVerifier})

	verAddr, err := economyModel.GetDiffVerifierAddress(mockPreBlock, mockBlock)
	assert.NoError(t, err)
	log.DLogger.Info("the verAddr is:", zap.Any("verAddr", verAddr))
	//assert.Equal(t,map[economy_model.VerifierType][]common.Address{},verAddr)
}

func TestDIP(t *testing.T) {
	tmpValue := int64(^uint64(1) >> 1)
	fmt.Printf("the tmpValue is:%x\r\n", tmpValue)
}
