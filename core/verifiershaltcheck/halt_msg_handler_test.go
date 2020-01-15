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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHaltHandler(t *testing.T) {
	assert.NotNil(t, MakeTestHaltHandle(t, 0))
}

func TestVBHaltHandler_ProposeEmptyBlock(t *testing.T) {
	_, err := ProposeEmptyBlockForTest(t, 0)
	assert.NoError(t, err)
}

func TestVBHaltHandler_OnNewProposalMsg(t *testing.T) {
	testMsg, err := MakeTestProposalMsg(t, 0)
	assert.NoError(t, err)
	
	testHandle, err := ProposeEmptyBlockForTest(t, 0)
	assert.NoError(t, err)
	
	// TODO: the aliveVerifierVote recover Address is invalid
	err = testHandle.OnNewProposalMsg(*testMsg)
	// assert.NoError(t, err)
}

func TestVBHaltHandler_GetProposalMsg(t *testing.T) {
	handle, err := ProposeEmptyBlockForTest(t, 0)
	assert.NoError(t, err)
	assert.NotNil(t, handle.GetProposalMsg())
}

func TestVBHaltHandler_GetOtherProposalMessages(t *testing.T) {
	handle, err := ProposeEmptyBlockForTest(t, 1)
	assert.NoError(t, err)
	assert.NotNil(t, handle.GetOtherProposalMessages())
}

func TestVBHaltHandler_GetMinProposalMsg(t *testing.T) {
	handle, err := ProposeEmptyBlockForTest(t, 0)
	assert.NoError(t, err)
	assert.NotNil(t, handle.GetMinProposalMsg())
}

func TestVBHaltHandler_HandlerProposalMessages(t *testing.T) {
	// TODO: goroutine
}

func TestVBHaltHandler_HandlerAliveVerVotes(t *testing.T) {
	handle, err := ProposeEmptyBlockForTest(t, 0)
	assert.NoError(t, err)
	assert.NotNil(t, handle)
	// TODO: New a Account for test
	// verVote, err := GenVoteMsg(handle.GetProposalMsg().EmptyBlock, signHash, )
}

func TestVBHaltHandler_VotesLen(t *testing.T) {
	type fields struct {
		pgConfig                 ProposalGeneratorConfig
		proposalMsg              *ProposalMsg
		proposalMessagesByOthers []ProposalMsg
		minProposalMsg           ProposalMsg
		aliveVerVotes            map[common.Address]model.VoteMsg
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &VBHaltHandler{
				pgConfig:                 tt.fields.pgConfig,
				proposalMsg:              tt.fields.proposalMsg,
				proposalMessagesByOthers: tt.fields.proposalMessagesByOthers,
				minProposalMsg:           tt.fields.minProposalMsg,
				aliveVerVotes:            tt.fields.aliveVerVotes,
			}
			if got := handler.VotesLen(); got != tt.want {
				t.Errorf("VBHaltHandler.VotesLen() = %v, want %v", got, tt.want)
			}
		})
	}
}
//
// func TestVBHaltHandler_MinProposalMsg(t *testing.T) {
// 	type fields struct {
// 		pgConfig                 ProposalGeneratorConfig
// 		proposalMsg              *ProposalMsg
// 		proposalMessagesByOthers []ProposalMsg
// 		minProposalMsg           ProposalMsg
// 		aliveVerVotes            map[common.Address]model.VoteMsg
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		want   ProposalMsg
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			handler := &VBHaltHandler{
// 				pgConfig:                 tt.fields.pgConfig,
// 				proposalMsg:              tt.fields.proposalMsg,
// 				proposalMessagesByOthers: tt.fields.proposalMessagesByOthers,
// 				minProposalMsg:           tt.fields.minProposalMsg,
// 				aliveVerVotes:            tt.fields.aliveVerVotes,
// 			}
// 			if got := handler.MinProposalMsg(); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("VBHaltHandler.MinProposalMsg() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestNewAliveVerHaltHandler(t *testing.T) {
// 	type args struct {
// 		signFunc SignHashFunc
// 		addr     common.Address
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *AliveVerHaltHandler
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewAliveVerHaltHandler(tt.args.signFunc, tt.args.addr); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewAliveVerHaltHandler() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestAliveVerHaltHandler_OnMinimalHashBlock(t *testing.T) {
// 	type fields struct {
// 		signHashFunc     SignHashFunc
// 		ownAddress       common.Address
// 		receivedProposal ProposalMsg
// 		ownVote          model.VoteMsg
// 	}
// 	type args struct {
// 		selectedProposal ProposalMsg
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *model.VoteMsg
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			handler := &AliveVerHaltHandler{
// 				signHashFunc:     tt.fields.signHashFunc,
// 				ownAddress:       tt.fields.ownAddress,
// 				receivedProposal: tt.fields.receivedProposal,
// 				ownVote:          tt.fields.ownVote,
// 			}
// 			got, err := handler.OnMinimalHashBlock(tt.args.selectedProposal)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("AliveVerHaltHandler.OnMinimalHashBlock() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("AliveVerHaltHandler.OnMinimalHashBlock() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
