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

package model

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"go.uber.org/zap"
	"math"
	"math/big"
)

var (
	DefaultPriorityCalculator = getCalculator()

	epsilon11 = 10000.0
	epsilon12 = 0.001
	epsilon13 = 1000.0

	epsilon21 = 100.0
	epsilon22 = 0.008
	epsilon23 = 500.0

	epsilon31 = 1.0
	epsilon32 = 5.0
	epsilon33 = 0.1

	scaleElem = float64(10000)

	w1 = 0.2
	w2 = 0.1
	w3 = 0.7

	StakeValMin = float64(100) // todo: this value is too low, maybe should mul dip(1000000000000000000)
)

func getCalculator() PriofityCalculator {
	return TestCalculator{}
}

func CalPriority(hash common.Hash, reputation uint64) (uint64, error) {
	priority := uint64(float64(reputation) * math.Pow(float64(hash[31])/256, 3))
	return priority, nil
}

func CalReputation(nonce uint64, stake *big.Int, performance uint64) (uint64, error) {
	stakeVal := float64(stake.Int64())
	if stakeVal < StakeValMin {
		log.DLogger.Info("the stakeVal is:", zap.Float64("stakeVal", stakeVal))
		return 0, errors.New("stake not sufficient")
	}

	nonceVal := float64(nonce)

	R1, err := Elem(stakeVal, epsilon11, epsilon12, epsilon13)
	if err != nil {
		return uint64(0), err
	}
	//fmt.Println("R1=", R1)

	R2, err := Elem(nonceVal, epsilon21, epsilon22, epsilon23)
	if err != nil {
		return uint64(0), err
	}

	var R3 float64
	performanceFloat := float64(performance) / 100
	R3, err = Elem(performanceFloat, epsilon31, epsilon32, epsilon33)
	if err != nil {
		return uint64(0), err
	}

	return uint64(w1*R1 + w2*R2 + w3*R3), nil
}

func Elem(x float64, epsilon1 float64, epsilon2 float64, epsilon3 float64) (float64, error) {

	if x < 0 {
		return 0, errors.New("invalid number")
	} else if x >= epsilon1 {
		return scaleElem / (1 + math.Exp(-epsilon2*(x-epsilon3))), nil
	} else {
		return scaleElem * math.Sqrt(x) / (math.Sqrt(epsilon1) * (1 + math.Exp(epsilon2*(epsilon3-epsilon1)))), nil
	}
}

//func Elembig (x *big.Float) (*big.Float, error){
//	y,_:=x.Float64()
//	z,_:=Elem(y)
//	return big.NewFloat(z),nil
//}

type TestCalculator struct {
}

// hash means luck which calculate by block seed and address
func (tc TestCalculator) GetElectPriority(hash common.Hash, nonce uint64, stake *big.Int, performance uint64) (uint64, error) {
	//log.DLogger.Info("get elect priority", "hash", hash.Hex(), "nonce", nonce, "stake", stake.String(), "performance", performance)
	reputation, err := CalReputation(nonce, stake, performance)
	//log.DLogger.Info("TestCalculator#GetElectPriority", "reputation", reputation, "seed", hash.Hex(), "nonce", nonce, "stake", stake.String(), "performance", performance)
	if err != nil {
		return uint64(0), err
	}
	priority, err := CalPriority(hash, reputation)
	if err != nil {
		return uint64(0), err
	} else {
		return priority, nil
	}
}

// acquire reputation value
func (tc TestCalculator) GetReputation(nonce uint64, stake *big.Int, performance uint64) (uint64, error) {
	//log.DLogger.Info("get reputation", "nonce", nonce, "stake", stake.String(), "performance", performance)
	return CalReputation(nonce, stake, performance)
}
