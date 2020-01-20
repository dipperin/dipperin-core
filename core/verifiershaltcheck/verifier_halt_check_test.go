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
	"github.com/dipperin/dipperin-core/core/chaincommunication"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/ethereum/go-ethereum/event"
)

func TestMakeSystemHaltedCheck(t *testing.T) {
	type args struct {
		conf *HaltCheckConf
	}
	tests := []struct {
		name string
		args args
		want *SystemHaltedCheck
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeSystemHaltedCheck(tt.args.conf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeSystemHaltedCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemHaltedCheck_SetMsgSigner(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	type args struct {
		walletSigner NeedWalletSigner
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			systemHaltedCheck.SetMsgSigner(tt.args.walletSigner)
		})
	}
}

func TestSystemHaltedCheck_MsgHandlers(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name   string
		fields fields
		want   map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if got := systemHaltedCheck.MsgHandlers(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SystemHaltedCheck.MsgHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemHaltedCheck_onCurrentBlockNumberRequest(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	type args struct {
		msg p2p.Msg
		p   chaincommunication.PmAbstractPeer
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
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.onCurrentBlockNumberRequest(tt.args.msg, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.onCurrentBlockNumberRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_onCurrentBlockNumberResponse(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	type args struct {
		msg p2p.Msg
		p   chaincommunication.PmAbstractPeer
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
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.onCurrentBlockNumberResponse(tt.args.msg, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.onCurrentBlockNumberResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_checkPeerHeight(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.checkPeerHeight(); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.checkPeerHeight() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_onProposeEmptyBlockMsg(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	type args struct {
		msg p2p.Msg
		p   chaincommunication.PmAbstractPeer
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
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.onProposeEmptyBlockMsg(tt.args.msg, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.onProposeEmptyBlockMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_onSendMinimalHashBlock(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	type args struct {
		msg p2p.Msg
		p   chaincommunication.PmAbstractPeer
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
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.onSendMinimalHashBlock(tt.args.msg, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.onSendMinimalHashBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_onSendMinimalHashBlockResponse(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	type args struct {
		msg p2p.Msg
		p   chaincommunication.PmAbstractPeer
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
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.onSendMinimalHashBlockResponse(tt.args.msg, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.onSendMinimalHashBlockResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_proposeEmptyBlock(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.proposeEmptyBlock(); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.proposeEmptyBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_sendMinimalHashBlock(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	type args struct {
		proposal ProposalMsg
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
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.sendMinimalHashBlock(tt.args.proposal); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.sendMinimalHashBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_handleFinalEmptyBlock(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
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
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.handleFinalEmptyBlock(tt.args.proposal, tt.args.votes); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.handleFinalEmptyBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_checkVerClusterStatus(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.checkVerClusterStatus(); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.checkVerClusterStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_LogCurrentVerifier(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			systemHaltedCheck.LogCurrentVerifier()
		})
	}
}

func TestSystemHaltedCheck_LogConnectedCurrentVerifier(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			systemHaltedCheck.LogConnectedCurrentVerifier()
		})
	}
}

func TestSystemHaltedCheck_loop(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			systemHaltedCheck.loop()
		})
	}
}

func TestSystemHaltedCheck_Start(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			if err := systemHaltedCheck.Start(); (err != nil) != tt.wantErr {
				t.Errorf("SystemHaltedCheck.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystemHaltedCheck_Stop(t *testing.T) {
	type fields struct {
		SynStatus             uint32
		nodeType              int
		handlers              map[uint64]func(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error
		otherBootNodeHeight   map[string]uint64
		verifierHeight        map[string]uint64
		verifierMaxHeight     uint64
		broadcaster           broadcastEmptyBlock
		haltHandler           *VBHaltHandler
		haltCheckStateHandle  *StateHandler
		csProtocol            CsProtocolFunction
		startEmptyProcessFlag uint32
		proposalFail          chan bool
		aliveVerifierVote     chan model.VoteMsg
		stopEmptyProcess      chan bool
		selectedProposal      chan ProposalMsg
		proposalInfoMsg       chan ProposalMsg
		heightInfo            chan heightResponseInfo
		quit                  chan bool
		feed                  event.Feed
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			systemHaltedCheck := &SystemHaltedCheck{
				SynStatus:             tt.fields.SynStatus,
				nodeType:              tt.fields.nodeType,
				handlers:              tt.fields.handlers,
				otherBootNodeHeight:   tt.fields.otherBootNodeHeight,
				verifierHeight:        tt.fields.verifierHeight,
				verifierMaxHeight:     tt.fields.verifierMaxHeight,
				broadcaster:           tt.fields.broadcaster,
				haltHandler:           tt.fields.haltHandler,
				haltCheckStateHandle:  tt.fields.haltCheckStateHandle,
				csProtocol:            tt.fields.csProtocol,
				startEmptyProcessFlag: tt.fields.startEmptyProcessFlag,
				proposalFail:          tt.fields.proposalFail,
				aliveVerifierVote:     tt.fields.aliveVerifierVote,
				stopEmptyProcess:      tt.fields.stopEmptyProcess,
				selectedProposal:      tt.fields.selectedProposal,
				proposalInfoMsg:       tt.fields.proposalInfoMsg,
				heightInfo:            tt.fields.heightInfo,
				quit:                  tt.fields.quit,
				feed:                  tt.fields.feed,
			}
			systemHaltedCheck.Stop()
		})
	}
}
