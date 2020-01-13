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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVoteMsgType_String(t *testing.T) {
	result := VoteMsgType(byte(0)).String()
	assert.Equal(t, PreVoteMessage.String(), result)

	result = VoteMsgType(byte(1)).String()
	assert.Equal(t, VoteMessage.String(), result)

	result = VoteMsgType(byte(2)).String()
	assert.Equal(t, VerBootNodeVoteMessage.String(), result)

	result = VoteMsgType(byte(3)).String()
	assert.Equal(t, AliveVerifierVoteMessage.String(), result)

	result = VoteMsgType(byte(4)).String()
	assert.NotEqual(t, AliveVerifierVoteMessage.String(), result)

	result = VoteMsgType(byte(5)).String()
	assert.NotEqual(t, AliveVerifierVoteMessage.String(), result)
}
