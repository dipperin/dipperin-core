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
	"github.com/dipperin/dipperin-core/core/economymodel"
	"reflect"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/ethereum/go-ethereum/ethdb"
)

// TODO: need env of tests

func TestChainState_BuildStateProcessor(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		preAccountStateRoot common.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *stateprocessor.AccountStateDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChainState{
				ChainStateConfig: tt.fields.ChainStateConfig,
				ethDB:            tt.fields.ethDB,
				ChainDB:          tt.fields.ChainDB,
				StateStorage:     tt.fields.StateStorage,
				EconomyModel:     tt.fields.EconomyModel,
			}
			got, err := cs.BuildStateProcessor(tt.args.preAccountStateRoot)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChainState.BuildStateProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.BuildStateProcessor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetStateStorage(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	tests := []struct {
		name   string
		fields fields
		want   stateprocessor.StateStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChainState{
				ChainStateConfig: tt.fields.ChainStateConfig,
				ethDB:            tt.fields.ethDB,
				ChainDB:          tt.fields.ChainDB,
				StateStorage:     tt.fields.StateStorage,
				EconomyModel:     tt.fields.EconomyModel,
			}
			if got := cs.GetStateStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetStateStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_CurrentState(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	tests := []struct {
		name    string
		fields  fields
		want    *stateprocessor.AccountStateDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChainState{
				ChainStateConfig: tt.fields.ChainStateConfig,
				ethDB:            tt.fields.ethDB,
				ChainDB:          tt.fields.ChainDB,
				StateStorage:     tt.fields.StateStorage,
				EconomyModel:     tt.fields.EconomyModel,
			}
			got, err := cs.CurrentState()
			if (err != nil) != tt.wantErr {
				t.Errorf("ChainState.CurrentState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.CurrentState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_StateAtByBlockNumber(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		num uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *stateprocessor.AccountStateDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChainState{
				ChainStateConfig: tt.fields.ChainStateConfig,
				ethDB:            tt.fields.ethDB,
				ChainDB:          tt.fields.ChainDB,
				StateStorage:     tt.fields.StateStorage,
				EconomyModel:     tt.fields.EconomyModel,
			}
			got, err := cs.StateAtByBlockNumber(tt.args.num)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChainState.StateAtByBlockNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.StateAtByBlockNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_StateAtByStateRoot(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		root common.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *stateprocessor.AccountStateDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &ChainState{
				ChainStateConfig: tt.fields.ChainStateConfig,
				ethDB:            tt.fields.ethDB,
				ChainDB:          tt.fields.ChainDB,
				StateStorage:     tt.fields.StateStorage,
				EconomyModel:     tt.fields.EconomyModel,
			}
			got, err := cs.StateAtByStateRoot(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChainState.StateAtByStateRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.StateAtByStateRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
