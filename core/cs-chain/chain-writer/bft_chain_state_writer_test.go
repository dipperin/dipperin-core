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


package chain_writer

import (
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/tests/g-mockFile"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBftChainWriter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mc := g_mockFile.NewMockChainInterface(controller)
	mb := NewMockAbstractBlock(controller)
	mb.EXPECT().IsSpecial().Return(true)
	mb.EXPECT().Version().Return(uint64(100))
	mc.EXPECT().GetChainConfig().Return(chain_config.GetChainConfig()).AnyTimes()

	assert.Error(t, NewBftChainWriter(&middleware.BftBlockContext{
		BlockContext: middleware.BlockContext{ Block: mb, Chain: mc},
	}, mc).SaveBlock())

}

func TestBftChainWriter_SaveBlock(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mc := g_mockFile.NewMockChainInterface(controller)
	mb := NewMockAbstractBlock(controller)
	mb.EXPECT().IsSpecial().Return(true)
	mb.EXPECT().Version().Return(uint64(100))
	mc.EXPECT().GetChainConfig().Return(chain_config.GetChainConfig()).AnyTimes()

	assert.Error(t, NewBftChainWriterWithoutVotes(&middleware.BftBlockContextWithoutVotes{
		BlockContext: middleware.BlockContext{ Block: mb, Chain: mc },
	}, mc).SaveBlock())
}
