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

// code generated by "stringer -type=RoundStepType"; DO NOT EDIT.

package model

import "strconv"

//const _RoundStepType_name = "RoundWaitBlockRoundStepNewHeightRoundStepNewRoundRoundStepProposeRoundStepPrevoteRoundStepPrevoteWaitRoundStepPrecommitRoundStepPrecommitWaitRoundStepCommit"
const _RoundStepType_name = "RoundStepNewHeightRoundStepNewRoundRoundStepProposeRoundStepPreVoteRoundStepPreCommit"

//var _RoundStepType_index = [...]uint8{14, 32, 49, 65, 81, 101, 119, 141, 156}
var _RoundStepType_index = [...]uint8{0, 18, 35, 51, 67, 85}

func (i RoundStepType) String() string {
	if i < 0 || i >= RoundStepType(len(_RoundStepType_index)-1) {
		return "RoundStepType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _RoundStepType_name[_RoundStepType_index[i]:_RoundStepType_index[i+1]]
}
