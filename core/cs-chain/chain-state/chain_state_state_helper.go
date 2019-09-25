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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
)

func (cs *ChainState) BuildStateProcessor(preAccountStateRoot common.Hash) (*state_processor.AccountStateDB, error) {
	return state_processor.NewAccountStateDB(preAccountStateRoot, cs.StateStorage)
}

func (cs *ChainState) GetStateStorage() state_processor.StateStorage {
	return cs.StateStorage
}

func (cs *ChainState) CurrentState() (*state_processor.AccountStateDB, error) {
	curBlock := cs.CurrentBlock()

	if curBlock == nil {
		return nil, errors.New("current block is nil")
	}

	stateRoot := cs.CurrentBlock().StateRoot()

	return cs.StateAtByStateRoot(stateRoot)
}

func (cs *ChainState) StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error) {
	block := cs.GetBlockByNumber(num)

	if block == nil {
		return nil, errors.New("block not found")
	}

	return cs.BuildStateProcessor(block.StateRoot())
}

func (cs *ChainState) StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error) {
	return cs.BuildStateProcessor(root)
}
