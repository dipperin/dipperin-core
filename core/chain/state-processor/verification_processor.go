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

package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"math"
)

var (
	performanceInitial = uint64(30)
	reward             = float64(1)
	// penalty            = float64(-10)
)

/*func (state *AccountStateDB) processCommitNum(v *model.VoteMsg) error {
	CommitNum, CommitErr := state.GetCommitNum(v.GetAddress())
	if CommitErr != nil {
		return CommitErr
	}
	state.SetCommitNum(v.GetAddress(), CommitNum+1)
	return nil
}*/

func (state *AccountStateDB) ProcessPerformance(address common.Address, amount float64) error {
	performance, performanceErr := state.GetPerformance(address)
	commitNum, _ := state.GetCommitNum(address)
	verifyNum, _ := state.GetVerifyNum(address)
	if performanceErr != nil {
		return performanceErr
	}

	currentPerf := float64(performance) + amount
	currentPerf = math.Max(currentPerf, 0)
	performanceCeil := float64(100)
	if verifyNum != 0 {
		performanceCeil = float64(commitNum) * 100 / float64(verifyNum)
	}
	currentPerf = math.Min(currentPerf, performanceCeil)
	state.SetPerformance(address, uint64(currentPerf))
	return nil
}

func (state *AccountStateDB) ProcessVerification(v model.AbstractVerification, index int) error {
	commitNum, commitErr := state.GetCommitNum(v.GetAddress())
	if commitErr != nil {
		return commitErr
	}
	state.SetCommitNum(v.GetAddress(), commitNum+1)
	return nil
}

func (state *AccountStateDB) ProcessVerifierNumber(address common.Address) error {
	verifyNum, verifyErr := state.GetVerifyNum(address)
	if verifyErr != nil {
		return verifyErr
	}
	state.SetVerifyNum(address, verifyNum+1)
	return nil
}
