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
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

// TODO: need env of tests

func TestChainState_Genesis(t *testing.T) {
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
		want   model.AbstractBlock
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
			if got := cs.Genesis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.Genesis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_CurrentBlock(t *testing.T) {
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
		want   model.AbstractBlock
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
			if got := cs.CurrentBlock(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.CurrentBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_CurrentHeader(t *testing.T) {
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
		want   model.AbstractHeader
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
			if got := cs.CurrentHeader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.CurrentHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBlock(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash   common.Hash
		number uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractBlock
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
			if got := cs.GetBlock(tt.args.hash, tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBlockByHash(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractBlock
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
			if got := cs.GetBlockByHash(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetBlockByHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBlockByNumber(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		number uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractBlock
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
			if got := cs.GetBlockByNumber(tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetBlockByNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_HasBlock(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash   common.Hash
		number uint64
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
			if got := cs.HasBlock(tt.args.hash, tt.args.number); got != tt.want {
				t.Errorf("ChainState.HasBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBody(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractBody
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
			if got := cs.GetBody(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBodyRLP(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   rlp.RawValue
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
			if got := cs.GetBodyRLP(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetBodyRLP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetHeader(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash   common.Hash
		number uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractHeader
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
			if got := cs.GetHeader(tt.args.hash, tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetHeaderByHash(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractHeader
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
			if got := cs.GetHeaderByHash(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetHeaderByHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetHeaderByNumber(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		number uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractHeader
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
			if got := cs.GetHeaderByNumber(tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetHeaderByNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetHeaderRLP(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   rlp.RawValue
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
			if got := cs.GetHeaderRLP(tt.args.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetHeaderRLP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_HasHeader(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash   common.Hash
		number uint64
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
			if got := cs.HasHeader(tt.args.hash, tt.args.number); got != tt.want {
				t.Errorf("ChainState.HasHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBlockNumber(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash common.Hash
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
			if got := cs.GetBlockNumber(tt.args.hash); got != tt.want {
				t.Errorf("ChainState.GetBlockNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetTransaction(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		txHash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AbstractTransaction
		want1  common.Hash
		want2  uint64
		want3  uint64
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
			got, got1, got2, got3 := cs.GetTransaction(tt.args.txHash)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetTransaction() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ChainState.GetTransaction() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ChainState.GetTransaction() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("ChainState.GetTransaction() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func TestChainState_GetReceipts(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash   common.Hash
		number uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.Receipts
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
			if got := cs.GetReceipts(tt.args.hash, tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetReceipts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBloomLog(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		hash   common.Hash
		number uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.Bloom
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
			if got := cs.GetBloomLog(tt.args.hash, tt.args.number); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetBloomLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetBloomBits(t *testing.T) {
	type fields struct {
		ChainStateConfig *ChainStateConfig
		ethDB            ethdb.Database
		ChainDB          chaindb.Database
		StateStorage     stateprocessor.StateStorage
		EconomyModel     economymodel.EconomyModel
	}
	type args struct {
		head    common.Hash
		bit     uint
		section uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
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
			if got := cs.GetBloomBits(tt.args.head, tt.args.bit, tt.args.section); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetBloomBits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainState_GetLatestNormalBlock(t *testing.T) {
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
		want   model.AbstractBlock
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
			if got := cs.GetLatestNormalBlock(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChainState.GetLatestNormalBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
