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

package g_error

import "errors"

// peer type error
var (
	PeerTypeError                 = errors.New("peer isn't verifier or bootNodeVerifier")
	AddressIsNotVerifierBootNode  = errors.New("the Address isn't verifier boot node")
	ProposalMsgDecodeError        = errors.New("proposal msg decode error")
	MinimalBlockDecodeError       = errors.New("decode minimal block msg error")
	EmptyBlockNumberNotMatchError = errors.New("the empty block number not match")
	VoteMsgBlockHashNotMatchError = errors.New("vote msg block hash not match")
	WaitEmptyBlockExpireError     = errors.New("wait empty block expire")
	VoteMsgDecodeError            = errors.New("decode aliveVerifierVote message error")
	AlreadyHaveVoteMsgError       = errors.New("already have this vote msg")
	GenProposalConfigError        = errors.New("generate proposal config error")
	AliveVoteBlockHashError       = errors.New("the alive verifier vote block hash error")
	ProposeNotEnough              = errors.New("the propose isn't enough")
)
