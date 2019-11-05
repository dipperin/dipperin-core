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

package chain_state

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
	"testing"
)

func TestChainStateWriter(t *testing.T) {
	//check.TestingT(t)
}

var _ = check.Suite(&chainWriterSuite{})

type chainWriterSuite struct {
	BaseChainSuite
}

func (suite *chainWriterSuite) TestChainState_SaveBlock(c *check.C) {
	block := suite.blockBuilder.Build()
	err := suite.chainState.SaveBlock(block)
	assert.NoError(c, err)

	fmt.Println(block.Number())
	assert.NoError(c, suite.chainState.Rollback(1))
	assert.Error(c, suite.chainState.SaveBlockWithoutVotes(block))
}
