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

package chainstate

import (
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/cschain/chainwriter"
	"gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner
func TestChainStateChainHelper(t *testing.T) { check.TestingT(t) }

type chainHelperSuite struct {
	cs *ChainState
}

var _ = check.Suite(&chainHelperSuite{})

func (suite *chainHelperSuite) SetUpSuite(c *check.C) {
	writerF := chainwriter.NewChainWriterFactory()
	suite.cs = NewChainState(&ChainStateConfig{
		DataDir:       "",
		WriterFactory: writerF,
		ChainConfig:   chainconfig.GetChainConfig(),
	})
	
	c.Check(suite.cs, check.NotNil)
}

func (suite *chainHelperSuite) Test_GetChainConfig(c *check.C) {
	chainConfig := suite.cs.GetChainConfig()
	
	c.Check(chainConfig, check.NotNil)
}

func (suite *chainHelperSuite) Test_GetEconomyModel(c *check.C) {
	economyModel := suite.cs.GetEconomyModel()
	
	c.Check(economyModel, check.NotNil)
}

func (suite *chainHelperSuite) Test_GetChainDB(c *check.C) {
	chainDB := suite.cs.GetChainDB()
	
	c.Check(chainDB, check.NotNil)
}
