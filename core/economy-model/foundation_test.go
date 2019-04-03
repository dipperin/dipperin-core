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
	"testing"
	"github.com/stretchr/testify/assert"
	"math/big"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/third-party/log"
)

var testPortion = map[int]int64{
	1: 50,
	2: 46,
	3: 38,
	4: 26,
	5: 10,
}

// test foundation

func TestDipperinFoundation_GetMineMasterEDIPReward(t *testing.T) {
	testFoundation := MakeDipperinFoundation(DIPProportion)

	DIPReward := big.NewInt(1740000000)
	tokenDecimal := 3
	decimalBase := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(tokenDecimal)), nil)
	// year 1~5
	for i := 1; i <= 5; i++ {
		blockNumber := uint64(i) * HeightAfterOneYear
		reward, err := testFoundation.GetMineMasterEDIPReward(DIPReward, blockNumber, tokenDecimal)
		assert.NoError(t, err)

		testReward := big.NewInt(0)
		testReward.Mul(big.NewInt(testPortion[i]), DIPReward)
		testReward.Mul(testReward, decimalBase)
		testReward.Div(testReward, big.NewInt(10))
		testReward.Div(testReward, big.NewInt(consts.DIP))

		log.Info("the reward is:", "reward", reward)
		assert.EqualValues(t, testReward, reward)
	}

}

func TestDipperinEconomyModel_GetVerifierEDIPReward(t *testing.T) {
	testFoundation := MakeDipperinFoundation(DIPProportion)

	DIPReward := map[VerifierType]*big.Int{
		MasterVerifier:    big.NewInt(7536231),
		CommitVerifier:    big.NewInt(15072463),
		NotCommitVerifier: big.NewInt(3768115),
	}

	tokenDecimal := 3
	decimalBase := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(tokenDecimal)), nil)
	// year 1~5
	for i := 1; i <= 1; i++ {
		blockNumber := uint64(i) * HeightAfterOneYear
		reward, err := testFoundation.GetVerifierEDIPReward(DIPReward, blockNumber, tokenDecimal)
		assert.NoError(t, err)

		for key, value := range reward {
			testReward := big.NewInt(0)
			testReward.Mul(big.NewInt(testPortion[i]), DIPReward[key])
			testReward.Mul(testReward, decimalBase)
			testReward.Div(testReward, big.NewInt(10))
			testReward.Div(testReward, big.NewInt(consts.DIP))

			log.Info("the reward is:", "type",key,"value", value)
			assert.EqualValues(t, testReward, value)
		}

	}
}
