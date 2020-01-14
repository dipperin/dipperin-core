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

package gerror

import "errors"

var (
	/*Validate votes errors*/
	ErrFirstBlockHaveVerifications    = errors.New("first block shouldn't have verifications")
	ErrEmptyVoteList                  = errors.New("empty vote list")
	ErrSameVoteSingerInVotes          = errors.New("same vote signer in votes")
	ErrInvalidFirstVoteInSpecialBlock = errors.New("invalid first vote in special block")
	ErrBlockVotesNotEnough            = errors.New("block votes not enough")
	ErrNotCurrentVerifier             = errors.New("vote signer is not current verifier")
	ErrVerificationRootNotMatch       = errors.New("verification root not match")
	ErrInvalidBlockHashInVotes        = errors.New("invalid block hash in votes")
	ErrRegisterRootNotMatch           = errors.New("register root not match")

	/*Validate block errors*/
	ErrChainOrBlockIsNil       = errors.New("chain or block is nil")
	ErrBlockHeightTooLow       = errors.New("block height is too low")
	ErrNormalBlockHeightTooLow = errors.New("normal block height is too low")
	ErrFutureBlock             = errors.New("future block") //not an error
	ErrFutureBlockTooFarAway   = errors.New("future block too far away")
	ErrPreBlockIsNil           = errors.New("pre block is nil")
	ErrPreBlockHashNotMatch    = errors.New("pre block hash not match")
	ErrInvalidDiff             = errors.New("invalid difficulty for this block")
	ErrInvalidHashDiff         = errors.New("invalid block hash for difficulty")
	ErrInvalidCoinBase         = errors.New("invalid special block CoinBase")
	ErrPkIsNil                 = errors.New("header pk is nil")
	ErrSeedNotMatch            = errors.New("block seed not match")
	ErrCoinBaseNotMatch        = errors.New("pk doesn't belongs to CoinBase")
	ErrInvalidBlockVersion     = errors.New("invalid block version")
	ErrInvalidBlockTimeStamp   = errors.New("invalid block time stamp")
	ErrInvliadHeaderGasLimit   = errors.New("invalid header gas limit")
	ErrHeaderGasLimitNotEnough = errors.New("header gas limit not enough compare parent block")

	/*Validate tx errors*/
	ErrTxRootNotMatch           = errors.New("transaction root not match")
	ErrTxInSpecialBlock         = errors.New("special block have transactions")
	ErrTxGasLimitNotEnough      = errors.New("tx gas limit not enough")
	ErrTxSenderBalanceNotEnough = errors.New("tx sender balance not enough")
	ErrTxSenderStakeNotEnough   = errors.New("tx sender stake not enough")
	ErrTxTargetStakeNotEnough   = errors.New("tx target stake not enough")
	ErrInvalidTxType            = errors.New("invalid tx type, no validator for tx")
	ErrTxDelegatesNotEnough     = errors.New("register tx delegate is not enough")
	ValidateSendRegisterTxFirst = errors.New("validate: need to send register tx first")
	ValidateSendCancelTxFirst   = errors.New("validate: need to send cancel tx first")
	ErrInvalidEvidenceTime      = errors.New("invalid evidence time")
	ErrEvidenceVoteNotConflict  = errors.New("evidence vote not conflict")
	ErrTxTargetAddressNotMatch  = errors.New("tx target address not match")
	ErrInvalidUnStakeTime       = errors.New("invalid unStake time")

	/*Insert receipts errors*/
	ErrReceiptHashNotMatch      = errors.New("receipt hash not match")
	ErrInvalidHeaderGasUsed     = errors.New("invalid header gas used")
	ErrHeaderGasUsedOverRanging = errors.New("header gas used is over-ranging ")
	ErrTxReceiptIsNil           = errors.New("tx receipt is nil")

	/*Insert block errors*/
	ErrInvalidBlockNum = errors.New("invalid block number")
)
