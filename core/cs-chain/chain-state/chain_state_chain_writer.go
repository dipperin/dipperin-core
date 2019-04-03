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
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"errors"
	"github.com/dipperin/dipperin-core/third-party/log"
)

func (cs *ChainState) SaveBlock(block model.AbstractBlock) error {
	log.Info("chain state save block")
	return cs.WriterFactory.NewWriter(middleware.NewBlockContext(block, cs)).SaveBlock()
}

func (cs *ChainState) SaveBlockWithoutVotes(block model.AbstractBlock) error {
	log.Info("chain state SaveBlockWithoutVotes")
	return cs.WriterFactory.NewWriter(middleware.NewBftBlockContextWithoutVotes(block, cs)).SaveBlock()
}

func (cs *ChainState) Rollback(target uint64) error {
	return errors.New("Rollback do nothings")
}
