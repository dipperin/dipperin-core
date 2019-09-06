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
	ErrBlockNumberNotContinuous = errors.New("block number not continuous")
	ErrBlockPreHashNotMatch     = errors.New("block pre hash not match")
	ErrNoGenesis                = errors.New("genesis not found in chain")
	ErrFutureBlock              = errors.New("future block") //not an error
	ErrFutureBlockTooFarAway    = errors.New("future block too far away")
	ErrStateRootNotMatch        = errors.New("state root not match")
	ErrRegisterRootNotMatch     = errors.New("register root not match")
	// ErrUnknownAncestor is returned when validating a block requires an ancestor
	// that is unknown.
	ErrUnknownAncestor = errors.New("unknown ancestor")
	//ErrUnexpectedReward is returned when validating a block and reward for miner is more than expected
	ErrUnexpectedReward = errors.New("unexpected reward")
	//ErrHashNotMatchAncestor is returned when block hash not match ancestor head hash
	ErrHashNotMatchAncestor = errors.New("preHash not match ancestor")
	//ErrBadMerkleRoot is returned when merkle root not match with the merkle hash in block header.
	ErrBadMerkleRoot = errors.New("bad merkle root")

	/*Transaction errors*/
	//ErrDuplicateTx is returned when a transaction is already in the tx pool.
	ErrDuplicateTx = errors.New("duplicate transaction")
	//ErrFeeNotEnough is returned when in transaction fee is not enough for a transaction.
	ErrFeeNotEnough = errors.New("tx fee not enough")
	//ErrBadSignature is returned when the signature not match
	ErrBadSignature = errors.New("signature not match")
	//ErrEmptyTxouts is returned when a transaction has no txouts
	ErrEmptyTxouts = errors.New("txouts is empty")
	//ErrBadIns is returned when one Ins has a nil outpoint
	ErrBadIns = errors.New("input has a nil outpoint")
	//ErrWitNotMatch is returned when on transaction witness is not mach
	ErrWitNotMatch = errors.New("witness not match")
	//ErrNotEnoughCredit is returned when a transaction try to spent more than coin
	ErrNotEnoughCredit = errors.New("credit smaller than spent")
	//ErrAlreadyHaveThisBlock is returned when
	ErrAlreadyHaveThisBlock = errors.New("already have this block")

	ErrBlockHeightTooLow                   = errors.New("block height too low")
	ErrBlockHeightIsCurrentAndIsNotSpecial = errors.New("block height is the same as current block height and isn't empty block")
	ErrBlockSizeTooLarge                   = errors.New("block size too large")

	ErrBlockNotFound     = errors.New("block not found")
	ErrCurrentBlockIsNil = errors.New("current block is nil")

	ErrPreBlockIsNil        = errors.New("pre block cannot be null")
	ErrPreBlockHashNotMatch = errors.New("pre block hash not match")

	ErrSpecialInvalidCoinBase = errors.New("invalid special block CoinBase address")
	ErrInvalidDiff            = errors.New("invalid difficulty for this block")
	ErrWrongHashDiff          = errors.New("block hash not valid for difficulty")

	ErrNotGetPk        = errors.New("can not get pk from header")
	ErrSeedNotMatch    = errors.New("block seed not match")
	ErrPkNotIsCoinBase = errors.New("pk not belongs to CoinBase")

	ErrBlockVer       = errors.New("block version not accept")
	ErrBlockTimeStamp = errors.New("the block time stamp is invalid")

	ErrReceiptIsNil    = errors.New("the transaction receipt is nil")
	ErrReceiptNotFound = errors.New("the transaction receipt not found")
)
