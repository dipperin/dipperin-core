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

var (
	ErrNotCurrentVerifier                   = errors.New("vote signer is not current verifier")
	ErrSameVoteSingerInVotes                = errors.New("same vote signer in votes")
	ErrBlockVotesNotEnough                  = errors.New("block votes not enough")
	ErrInvalidBlockHashInVotes              = errors.New("invalid block hash in votes")
	ErrFirstBlockShouldNotHaveVerifications = errors.New("first block shouldn't have verifications")
	ErrInvalidFirstVoteInSpecialBlock       = errors.New("first vote in special block should be boot node's vote")
	ErrInvalidTxType                        = errors.New("invalid type, no validator for tx")
	ErrTxOverSize                           = errors.New("tx over size")
	ErrEmptyVoteList                        = errors.New("empty vote list")
	ErrTxNonceNotMatch                      = errors.New("tx nonce not match")
	ErrTxGasUsedIsOverGasLimit				= errors.New("the tx gasUsed is over the gasLimit")
)
