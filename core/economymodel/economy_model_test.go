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
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)


func TestDipperinEconomyModel_MapMerge(t *testing.T) {

	testCases := []struct{
		name string
		given func() (map[common.Address]*big.Int,map[common.Address]*big.Int)
		expect error
	} {
		{
			name:"MapMerge",
			given: func() (map[common.Address]*big.Int,map[common.Address]*big.Int) {
				src := map[common.Address]*big.Int{}
				des := map[common.Address]*big.Int{}

				src[common.HexToAddress("1234")] = big.NewInt(1)
				return des,src
			},
			expect:nil,
		},
		{
			name:"ErrAddressExist",
			given: func() (map[common.Address]*big.Int,map[common.Address]*big.Int) {
				src := map[common.Address]*big.Int{}
				des := map[common.Address]*big.Int{}

				src[common.HexToAddress("1234")] = big.NewInt(1)
				err := MapMerge(des, src)
				assert.NoError(t, err)
				return des, src
			},
			expect:gerror.ErrAddressExist,
		},
	}

	for _, tc := range testCases{
		t.Log("TestDipperinEconomyModel_MapMerge", tc.name)
		des,src := tc.given()
		err := MapMerge(des, src)
		assert.Equal(t, tc.expect, err)
	}
}

// test the reward for a single block (macroeconomics)
func TestDipperinEconomyModel_GetOneBlockTotalDIPReward(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)

	economyModel := MakeDipperinEconomyModel(service, DIPProportion)

	type result struct {
		reword *big.Int
		err error
	}

	testCases := []struct{
		name string
		given uint64
		expect result
	}{
		{
			name:"BlockNumberIs0",
			given: uint64(0),
			expect:result{big.NewInt(0), gerror.ErrBlockNumberIs0},
		},
		{
			name:"HeightBeforeTenYear",
			given: uint64(HeightAfterTenYear - 1),
			expect:result{TotalReWardDIPOneBlock, nil},
		},
		{
			name:"HeightInTenYear",
			given: uint64(HeightAfterTenYear),
			expect:result{TotalReWardDIPOneBlock, nil},
		},
		{
			name:"HeightAfterTenYear",
			given: uint64(HeightAfterTenYear),
			expect:result{TotalReWardDIPOneBlock, nil},
		},
	}

	for _,tc := range testCases{
		t.Log("TestDipperinEconomyModel_GetOneBlockTotalDIPReward", tc.name)
		num := tc.given
		reword, err := economyModel.GetOneBlockTotalDIPReward(num)
		assert.Equal(t, tc.expect.reword, reword)
		assert.Equal(t, tc.expect.err, err)
	}
}

// test the foundation
func TestDipperinEconomyModel_GetFoundation(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)

	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	f := economyModel.GetFoundation()
	assert.NotNil(t, f)
}

// test the minimum transaction fee
//func TestDipperinEconomyModel_GetMinimumTxFee(t *testing.T) {
//	fee := GetMinimumTxFee(30)
//	assert.NotNil(t, fee)
//}

// test address type
func TestDipperinEconomyModel_CheckAddressType(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)

	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	infoInvest := economyModel.GetInvestorInitBalance()

	infoDev := economyModel.GetDeveloperInitBalance()

	for k := range infoInvest {
		assert.Equal(t, economyModel.CheckAddressType(k), InvestorAddress)
	}

	for k := range infoDev {
		assert.Equal(t, economyModel.CheckAddressType(k), DeveloperAddress)
	}
	assert.Equal(t, economyModel.CheckAddressType(common.HexToAddress("0x1234")), NotEconomyModelAddress)
}

// test unlocking mechanism
func TestDipperinEconomyModel_GetAddressLockMoney(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)

	economyModel := MakeDipperinEconomyModel(service, DIPProportion)

	testData := make([]map[common.Address]*big.Int, 0)
	infoInvest := economyModel.GetInvestorInitBalance()
	testData = append(testData, infoInvest)
	infoDev := economyModel.GetDeveloperInitBalance()
	testData = append(testData, infoDev)
	infoMaintenance := economyModel.GetFoundationInfo(Maintenance)
	testData = append(testData, infoMaintenance)
	infoEarlyToken := economyModel.GetFoundationInfo(EarlyToken)
	testData = append(testData, infoEarlyToken)
	infoRemainReward := economyModel.GetFoundationInfo(RemainReward)
	testData = append(testData, infoRemainReward)

	for i, info := range testData {
		for k := range info {
			if i < 2 {
				//test investor and developer
				_, err := economyModel.GetAddressLockMoney(k, 0)
				assert.Error(t, err, gerror.ErrBlockNumberIs0.Error())
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
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)


	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	year, _ := economyModel.GetBlockYear(uint64(365 * 24 * 3600 / GenerateBlockDuration) )
	assert.Equal(t, year, uint64(1))
}

// test the reward for miners
func TestDipperinEconomyModel_GetMineMasterDIPReward(t *testing.T) {

	controller := gomock.NewController(t)
	defer controller.Finish()
	service := NewMockEconomyNeedService(controller)
	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	testBlock := model.NewMockAbstractBlock(controller)


	type result struct {
		reword *big.Int
		err error
	}


	testCases := []struct{
		name string
		given func() *model.MockAbstractBlock
		expect result
	}{
		{
			name:"HeightInTenYear",
			given: func() *model.MockAbstractBlock {
				testBlock.EXPECT().Number().Return(uint64(30))
				return testBlock
			},
			expect:result{big.NewInt(0).Mul(big.NewInt(17400000000), big.NewInt(consts.GDIPUNIT)), nil},
		},
		{
			name:"HeightAfterTenYear",
			given: func() *model.MockAbstractBlock {
				testBlock.EXPECT().Number().Return(uint64(HeightAfterTenYear + 1))
				return testBlock
			},
			expect:result{big.NewInt(0).Mul(big.NewInt(8700000000), big.NewInt(consts.GDIPUNIT)), nil},
		},
		{
			name:"HeightAfterElevenYear",
			given: func() *model.MockAbstractBlock {
				testBlock.EXPECT().Number().Return(uint64(11*HeightAfterOneYear) + 1)
				return testBlock
			},
			expect:result{big.NewInt(0).Mul(big.NewInt(12441000000), big.NewInt(consts.GDIPUNIT)), nil},
		},
		{
			name:"genesisBlock",
			given: func() *model.MockAbstractBlock {
				testBlock.EXPECT().Number().Return(uint64(0))
				return testBlock
			},
			expect:result{big.NewInt(0), gerror.ErrBlockNumberIs0},
		},
	}

	for _,tc := range testCases{
		t.Log("TestDipperinEconomyModel_GetMineMasterDIPReward", tc.name)
		block := tc.given
		reword, err := economyModel.GetMineMasterDIPReward(block())
		assert.Equal(t, tc.expect.reword, reword)
		assert.Equal(t, tc.expect.err, err)
	}

}

// test the reward for verifiers
func TestDipperinEconomyModel_GetVerifierDIPReward(t *testing.T) {
	controller := gomock.NewController(t)
	service := NewMockEconomyNeedService(controller)
	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	defer controller.Finish()

	testCases := []struct{
		name string
		given func() *model.MockAbstractBlock
		expect error
		expectResult map[VerifierType]*big.Int
	}{
		{
			name:"HeightInTenYear",
			given: func()  *model.MockAbstractBlock {
				mockBlock1 := model.NewMockAbstractBlock(controller)
				mockBlock1.EXPECT().Number().Return(uint64(30)).AnyTimes()
				return mockBlock1
			},
			expect:nil,
			expectResult:map[VerifierType]*big.Int{MasterVerifier:big.NewInt(75362318840579710),CommitVerifier:big.NewInt(150724637681159420), NotCommitVerifier:big.NewInt(37681159420289855)},
		},
		{
			name:"HeightInTenYear",
			given: func()  *model.MockAbstractBlock {
				mockBlock1 := model.NewMockAbstractBlock(controller)
				mockBlock1.EXPECT().Number().Return(uint64(HeightAfterTenYear + 1)).AnyTimes()
				return mockBlock1
			},
			expect:nil,
			expectResult:map[VerifierType]*big.Int{MasterVerifier:big.NewInt(37681159420289855),CommitVerifier:big.NewInt(75362318840579710), NotCommitVerifier:big.NewInt(18840579710144927)},
		},
		{
			name:"genesisBlock",
			given: func()  *model.MockAbstractBlock {
				mockBlock1 := model.NewMockAbstractBlock(controller)
				mockBlock1.EXPECT().Number().Return(uint64(0)).AnyTimes()
				return mockBlock1
			},
			expect:gerror.ErrBlockNumberIs0,
		},
	}

	for _,tc := range testCases{
		t.Log("TestDipperinEconomyModel_GetVerifierDIPReward", tc.name)
		block  := tc.given()
		result, err := economyModel.GetVerifierDIPReward(block)
		if err != nil {
			assert.Equal(t, tc.expect, err)
		}else {
			assert.EqualValues(t, tc.expectResult[MasterVerifier], result[MasterVerifier])
			assert.EqualValues(t, tc.expectResult[CommitVerifier], result[CommitVerifier])
			assert.EqualValues(t, tc.expectResult[NotCommitVerifier], result[NotCommitVerifier])
		}
		assert.Equal(t, tc.expect, err)
	}
}

// test investors
func TestDipperinEconomyModel_GetInvestorInitBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)
	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	info := economyModel.GetInvestorInitBalance()
	assert.EqualValues(t, big.NewInt(0).Mul(big.NewInt(262800000000000000), big.NewInt(consts.GDIPUNIT)), info[InvestorAddresses[0]])
}

// test developers
func TestDipperinEconomyModel_GetDeveloperInitBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)
	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	info := economyModel.GetDeveloperInitBalance()
	assert.EqualValues(t, big.NewInt(0).Mul(big.NewInt(87600000000000000), big.NewInt(consts.GDIPUNIT)), info[DeveloperAddresses[0]])
}

// test investor unlocking mechanism
// todo  need to check the economy model
func TestDipperinEconomyModel_GetInvestorLockDIP(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)
	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	investorTotalDIP := big.NewInt(0).Mul(big.NewInt(262800000000000000), big.NewInt(consts.GDIPUNIT))

	// unlocking by quarters each year
	for i := 1; i <= 4; i++ {
		blockNumber := HeightAfterOneYear / 4 * uint64(i)
		lockDIP, err := economyModel.GetInvestorLockDIP(InvestorAddresses[0], blockNumber)
		assert.NoError(t, err)

		unlockDIP := big.NewInt(0).Div(investorTotalDIP, big.NewInt(10))
		unlockDIP.Mul(unlockDIP, big.NewInt(int64(i)))
		unlockDIP.Div(unlockDIP, big.NewInt(4))
		lockValue := big.NewInt(0).Sub(investorTotalDIP, unlockDIP)

		assert.EqualValues(t, lockValue, lockDIP)
	}

	_, err := economyModel.GetInvestorLockDIP(common.Address{}, 2)
	assert.Equal(t, gerror.ErrAddress, err)
}

// test developer unlocking mechanism
// todo  need to check the economy model
func TestDipperinEconomyModel_GetDeveloperLockDIP(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockEconomyNeedService(ctrl)
	economyModel := MakeDipperinEconomyModel(service, DIPProportion)
	investorTotalDIP := big.NewInt(0).Mul(big.NewInt(87600000000000000), big.NewInt(consts.GDIPUNIT))

	// unlocking by quarters each year
	for i := 1; i <= 4; i++ {
		blockNumber := HeightAfterOneYear / 4 * uint64(i)
		lockDIP, err := economyModel.GetDeveloperLockDIP(DeveloperAddresses[0], blockNumber)
		assert.NoError(t, err)

		unlockDIP := big.NewInt(0).Div(investorTotalDIP, big.NewInt(10))
		unlockDIP.Mul(unlockDIP, big.NewInt(int64(i)))
		unlockDIP.Div(unlockDIP, big.NewInt(4))
		lockValue := big.NewInt(0).Sub(investorTotalDIP, unlockDIP)

		assert.EqualValues(t, lockValue, lockDIP)
	}

	_, err := economyModel.GetDeveloperLockDIP(common.Address{}, 2)
	assert.Equal(t, gerror.ErrAddress, err)
}

//func TestGenerateAddress(t *testing.T) {
//	addressNumber := 5
//	for i := 0; i < addressNumber; i++ {
//		sk, err := crypto.GenerateKey()
//		assert.NoError(t, err)
//		address := cs_crypto.GetNormalAddress(sk.PublicKey)
//		log.DLogger.Info("the address is:", zap.String("address", address.Hex()))
//	}
//}

func TestCalcDIPTotalCirculation(t *testing.T) {
	value := CalcDIPTotalCirculation(0)
	assert.Equal(t, big.NewInt(0), value)

	value = CalcDIPTotalCirculation(1)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(78840000), big.NewInt(consts.DIP)), value)

	value = CalcDIPTotalCirculation(5)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(394200000), big.NewInt(consts.DIP)), value)

	value = CalcDIPTotalCirculation(15)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(3788167915366200000), big.NewInt(consts.GDIPUNIT)), value)
}

func TestDipperinEconomyModel_GetDiffVerifierAddress(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockPreBlock := model.NewMockAbstractBlock(controller)
	mockBlock := model.NewMockAbstractBlock(controller)
	service := NewMockEconomyNeedService(controller)
	slot := uint64(0)
	service.EXPECT().GetSlot(mockPreBlock).Return(&slot).AnyTimes()
	service.EXPECT().GetVerifiers(slot).Return([]common.Address{common.HexToAddress("0x000078b33598Be2b405206F44B018557e6F851FD230C"),}).AnyTimes()
	economyModel := MakeDipperinEconomyModel(service, DIPProportion)

	testCases := []struct{
		name string
		given func() (*model.MockAbstractBlock, *model.MockAbstractBlock)
		expect error
	}{
		{
			name:"ErrBlockNumberIs0Ore1",
			given: func() (*model.MockAbstractBlock, *model.MockAbstractBlock) {
				mockBlock.EXPECT().Number().Return(uint64(1))
				return mockPreBlock, mockBlock
			},
			expect:gerror.ErrBlockNumberIs0Ore1,
		},
		{
			name:"GetDiffVerifierAddressRight",
			given: func() (*model.MockAbstractBlock, *model.MockAbstractBlock) {
				mockPreBlock.EXPECT().Number().Return(uint64(2)).AnyTimes()
				mockBlock.EXPECT().Number().Return(uint64(3)).AnyTimes()

				mockVerifier := model.NewMockAbstractVerification(controller)
				mockVerifier.EXPECT().GetRound().Return(uint64(0))
				mockVerifier.EXPECT().GetAddress().Return(common.HexToAddress("0x000078b33598Be2b405206F44B018557e6F851FD230C")).AnyTimes()

				mockBlock.EXPECT().GetVerifications().Return(model.Verifications{mockVerifier})
				return mockPreBlock,mockBlock
			},
			expect:nil,
		},

	}
	for _,tc := range testCases{
		t.Log("TestDipperinEconomyModel_GetDiffVerifierAddress", tc.name)
		mockPreBlock, mockBlock := tc.given()
		_, err := economyModel.GetDiffVerifierAddress(mockPreBlock, mockBlock)
		assert.Equal(t, tc.expect, err)
	}


}

