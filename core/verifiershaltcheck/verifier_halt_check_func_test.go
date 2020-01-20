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
	"github.com/dipperin/dipperin-core/core/model"
)

func TestGenVoteMsg(t *testing.T) {
	type args struct {
		emptyBlock *model.Block
		signFunc   SignHashFunc
		addr       common.Address
		voteType   model.VoteMsgType
	}
	tests := []struct {
		name    string
		args    args
		want    *model.VoteMsg
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenVoteMsg(tt.args.emptyBlock, tt.args.signFunc, tt.args.addr, tt.args.voteType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenVoteMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenVoteMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkProposalValid(t *testing.T) {
	type args struct {
		proposal ProposalMsg
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkProposalValid(tt.args.proposal); (err != nil) != tt.wantErr {
				t.Errorf("checkProposalValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_selectEmptyProposal(t *testing.T) {
	type args struct {
		proposalA ProposalMsg
		proposalB ProposalMsg
	}
	tests := []struct {
		name string
		args args
		want ProposalMsg
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := selectEmptyProposal(tt.args.proposalA, tt.args.proposalB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectEmptyProposal() = %v, want %v", got, tt.want)
			}
		})
	}
}
