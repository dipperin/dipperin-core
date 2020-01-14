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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
)

func (cs *ChainState) BuildStateProcessor(preAccountStateRoot common.Hash) (*stateprocessor.AccountStateDB, error) {
	return stateprocessor.NewAccountStateDB(preAccountStateRoot, cs.StateStorage)
}

func (cs *ChainState) GetStateStorage() stateprocessor.StateStorage {
	return cs.StateStorage
}

func (cs *ChainState) CurrentState() (*stateprocessor.AccountStateDB, error) {
	curHeader := cs.CurrentHeader()

	if curHeader == nil {
		return nil, errors.New("current header is nil")
	}

	stateRoot := curHeader.GetStateRoot()

	return cs.StateAtByStateRoot(stateRoot)
}

func (cs *ChainState) StateAtByBlockNumber(num uint64) (*stateprocessor.AccountStateDB, error) {
	header := cs.GetHeaderByNumber(num)
	if header == nil {
		return nil, errors.New("header not found")
	}

	return cs.BuildStateProcessor(header.GetStateRoot())
}

func (cs *ChainState) StateAtByStateRoot(root common.Hash) (*stateprocessor.AccountStateDB, error) {
	return cs.BuildStateProcessor(root)
}
