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
	"reflect"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
)

// TODO: need env of tests

func TestChainState_BuildRegisterProcessor(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		preBlockRegisterRoot common.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *registerdb.RegisterDB
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
			got, err := cs.BuildRegisterProcessor(tt.args.preBlockRegisterRoot)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChainState.BuildRegisterProcessor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.BuildRegisterProcessor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_CurrentSeed(t *testing.T) {
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
		want   common.Hash
		want1  uint64
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
			got, got1 := cs.CurrentSeed()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.CurrentSeed() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ChainState.CurrentSeed() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestChainState_IsChangePoint(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		block                 model.AbstractBlock
		isProcessPackageBlock bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
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
			if got := cs.IsChangePoint(tt.args.block, tt.args.isProcessPackageBlock); got != tt.want {
				t.Errorf("ChainState.IsChangePoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetLastChangePoint(t *testing.T) {
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
		name   string
		fields fields
		args   args
		want   *uint64
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
			if got := cs.GetLastChangePoint(tt.args.block); got != tt.want {
				t.Errorf("ChainState.GetLastChangePoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetSlotByNum(t *testing.T) {
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
		name   string
		fields fields
		args   args
		want   *uint64
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
			if got := cs.GetSlotByNum(tt.args.num); got != tt.want {
				t.Errorf("ChainState.GetSlotByNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetSlot(t *testing.T) {
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
		name   string
		fields fields
		args   args
		want   *uint64
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
			if got := cs.GetSlot(tt.args.block); got != tt.want {
				t.Errorf("ChainState.GetSlot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetCurrVerifiers(t *testing.T) {
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
		want   []common.Address
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
			if got := cs.GetCurrVerifiers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetCurrVerifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetVerifiers(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		slot uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []common.Address
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
			if got := cs.GetVerifiers(tt.args.slot); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetVerifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetNextVerifiers(t *testing.T) {
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
		want   []common.Address
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
			if got := cs.GetNextVerifiers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetNextVerifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_NumBeforeLastBySlot(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		slot uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *uint64
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
			if got := cs.NumBeforeLastBySlot(tt.args.slot); got != tt.want {
				t.Errorf("ChainState.NumBeforeLastBySlot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetNumBySlot(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		slot uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *uint64
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
			if got := cs.GetNumBySlot(tt.args.slot); got != tt.want {
				t.Errorf("ChainState.GetNumBySlot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_CalVerifiers(t *testing.T) {
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
		name   string
		fields fields
		args   args
		want   []common.Address
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
			if got := cs.CalVerifiers(tt.args.block); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.CalVerifiers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_getTopVerifiers(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		address     common.Address
		priority    uint64
		topAddress  []common.Address
		topPriority []uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []common.Address
		want1  []uint64
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
			got, got1 := cs.getTopVerifiers(tt.args.address, tt.args.priority, tt.args.topAddress, tt.args.topPriority)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.getTopVerifiers() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ChainState.getTopVerifiers() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestChainState_calPriority(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		addr     common.Address
		blockNum uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint64
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
			got, err := cs.calPriority(tt.args.addr, tt.args.blockNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChainState.calPriority() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ChainState.calPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_getLuck(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		addr     common.Address
		blockNum uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   common.Hash
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
			if got := cs.getLuck(tt.args.addr, tt.args.blockNum); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.getLuck() = %v, want %v", got, tt.want)
			}
		})
	}
}
