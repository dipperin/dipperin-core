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

package verifiershaltcheck

import (
	"reflect"
	"testing"
	
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/economymodel"
	"github.com/dipperin/dipperin-core/core/model"
)

func TestMakeHaltCheckStateHandler(t *testing.T) {
	type args struct {
		needChainReader NeedChainReaderFunction
		walletSigner    NeedWalletSigner
		economyModel    economymodel.EconomyModel
	}
	tests := []struct {
		name string
		args args
		want *StateHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeHaltCheckStateHandler(tt.args.needChainReader, tt.args.walletSigner, tt.args.economyModel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeHaltCheckStateHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateHandler_GenProposalConfig(t *testing.T) {
	type fields struct {
		chainReader  NeedChainReaderFunction
		walletSigner NeedWalletSigner
		economyModel economymodel.EconomyModel
	}
	type args struct {
		voteType model.VoteMsgType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    ProposalGeneratorConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haltCheckStateHandle := &StateHandler{
				chainReader:  tt.fields.chainReader,
				walletSigner: tt.fields.walletSigner,
				economyModel: tt.fields.economyModel,
			}
			got, err := haltCheckStateHandle.GenProposalConfig(tt.args.voteType)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateHandler.GenProposalConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StateHandler.GenProposalConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStateHandler_ProcessAccountAndRegisterState(t *testing.T) {
	type fields struct {
		chainReader  NeedChainReaderFunction
		walletSigner NeedWalletSigner
		economyModel economymodel.EconomyModel
	}
	type args struct {
		block           model.AbstractBlock
		preStateRoot    common.Hash
		preRegisterRoot common.Hash
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		wantStateRoot    common.Hash
		wantRegisterRoot common.Hash
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			haltCheckStateHandle := &StateHandler{
				chainReader:  tt.fields.chainReader,
				walletSigner: tt.fields.walletSigner,
				economyModel: tt.fields.economyModel,
			}
			gotStateRoot, gotRegisterRoot, err := haltCheckStateHandle.ProcessAccountAndRegisterState(tt.args.block, tt.args.preStateRoot, tt.args.preRegisterRoot)
			if (err != nil) != tt.wantErr {
				t.Errorf("StateHandler.ProcessAccountAndRegisterState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStateRoot, tt.wantStateRoot) {
				t.Errorf("StateHandler.ProcessAccountAndRegisterState() gotStateRoot = %v, want %v", gotStateRoot, tt.wantStateRoot)
			}
			if !reflect.DeepEqual(gotRegisterRoot, tt.wantRegisterRoot) {
				t.Errorf("StateHandler.ProcessAccountAndRegisterState() gotRegisterRoot = %v, want %v", gotRegisterRoot, tt.wantRegisterRoot)
			}
		})
	}
}

func TestStateHandler_SaveFinalEmptyBlock(t *testing.T) {
	type fields struct {
		chainReader  NeedChainReaderFunction
		walletSigner NeedWalletSigner
		economyModel economymodel.EconomyModel
	}
	type args struct {
		proposal ProposalMsg
		votes    map[common.Address]model.VoteMsg
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
			haltCheckStateHandle := &StateHandler{
				chainReader:  tt.fields.chainReader,
				walletSigner: tt.fields.walletSigner,
				economyModel: tt.fields.economyModel,
			}
			if err := haltCheckStateHandle.SaveFinalEmptyBlock(tt.args.proposal, tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("StateHandler.SaveFinalEmptyBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
