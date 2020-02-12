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
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/economymodel"
	"testing"

	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
)

// TODO: need env of tests

func TestChainState_SaveBlock(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		block model.AbstractBlock
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := cs.SaveBlock(tt.args.block); (err != nil) != tt.wantErr {
				t.Errorf("ChainState.SaveBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainState_SaveBlockWithoutVotes(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		block model.AbstractBlock
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := cs.SaveBlockWithoutVotes(tt.args.block); (err != nil) != tt.wantErr {
				t.Errorf("ChainState.SaveBlockWithoutVotes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChainState_Rollback(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		target uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
			if err := cs.Rollback(tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("ChainState.Rollback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
