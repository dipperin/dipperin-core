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

package middleware

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateBlockVerifier(t *testing.T) {
	_, _, _, passChain := getTxTestEnv(t)
	ivc := &BlockContext{
		Block: &fakeBlock{},
		Chain: passChain,
	}
	assert.Error(t, UpdateBlockVerifier(ivc)())
	assert.Error(t, ValidBlockVerifier(ivc)())

	vc := &BlockContext{
		Block: &fakeBlock{ registerRoot: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421") },
		Chain: passChain,
	}
	assert.NoError(t, UpdateBlockVerifier(vc)())
	assert.NoError(t, ValidBlockVerifier(vc)())

	passChain.block.registerRoot = common.Hash{0x12}
	assert.Error(t, ValidBlockVerifier(vc)())

	passChain.slot = 1
	assert.NoError(t, NextRoundVerifier(vc)())
}

