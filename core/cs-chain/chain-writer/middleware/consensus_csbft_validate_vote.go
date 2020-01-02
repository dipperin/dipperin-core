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

package middleware

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
)

// SaveBlock use
func ValidateVotes(c *BftBlockContext) Middleware {
	return func() error {
		log.DLogger.Info("ValidateVotes start")
		// verify the number of seen commits meets the need
		block := c.Block
		slot := c.Chain.GetSlot(block)
		verifiers := c.Chain.GetVerifiers(*slot)
		if err := validVotesForBlock(c.Votes, block, verifiers); err != nil {
			return err
		}

		if err := validBlockHash(c.Votes, block); err != nil {
			return err
		}

		err := validateVotes(c.Block, c.Chain)
		if err != nil {
			return err
		}
		log.DLogger.Info("ValidateVotes success")
		return c.Next()
	}
}

// BFT use
func ValidateVotesForBFT(c *BlockContext) Middleware {
	return func() error {
		err := validateVotes(c.Block, c.Chain)
		if err != nil {
			return err
		}
		return c.Next()
	}
}

func validateVotes(block model.AbstractBlock, chain ChainInterface) error {

	// The first block has no votes
	if block.Number() == 1 {
		if !block.VerificationRoot().IsEqual(model.EmptyVerfRoot) || len(block.GetVerifications()) != 0 {
			return g_error.ErrFirstBlockHaveVerifications
		}
		return nil
	}

	// the v root of the current block in PBFT is the merkle root of the previous block's verifications, and the body's verifications are also that of the previous block
	if err := validVerificationRoot(block.GetVerifications(), block.VerificationRoot()); err != nil {
		return err
	}

	// 2.1 2.2 Verify votes is on previous block
	preBlockHeight := block.Number() - 1
	preBlock := chain.GetBlockByNumber(preBlockHeight)
	preBlockSlot := chain.GetSlot(preBlock)
	preBlockVerifiers := chain.GetVerifiers(*preBlockSlot)
	if err := validVotesForBlock(block.GetVerifications(), preBlock, preBlockVerifiers); err != nil {
		return err
	}

	if err := validBlockHash(block.GetVerifications(), preBlock); err != nil {
		return err
	}
	return nil
}

func validVotesForBlock(votes []model.AbstractVerification, block model.AbstractBlock, verifiers []common.Address) error {
	if len(votes) == 0 {
		return g_error.ErrEmptyVoteList
	}

	// valid different votes
	if sameVote(votes) {
		return g_error.ErrSameVoteSingerInVotes
	}

	// valid numbers of verifiers
	if block.IsSpecial() {
		v0 := votes[0]
		// check first vote
		if v0.GetType() != model.VerBootNodeVoteMessage {
			fmt.Println(v0.GetAddress())
			return g_error.ErrInvalidFirstVoteInSpecialBlock
		}

		for _, ver := range votes {
			if err := ver.HaltedVoteValid(verifiers); err != nil {
				return err
			}
		}
	} else {

		// check verification lens
		totalNeed := len(verifiers)*2/3 + 1
		if len(votes) < totalNeed {
			return g_error.ErrBlockVotesNotEnough
		}

		for _, ver := range votes {
			if err := ver.Valid(); err != nil {
				return err
			}

			if !verificationSignerInVerifiers(ver.GetAddress(), verifiers) {
				fmt.Println("invalid addr", ver.GetAddress(), "should in", verifiers)
				return g_error.ErrNotCurrentVerifier
			}
		}
	}
	return nil
}

func validVerificationRoot(verifications []model.AbstractVerification, vRoot common.Hash) error {

	targetRoot := model.DeriveSha(model.Verifications(verifications))
	//fmt.Println("===t root", targetRoot)
	if !targetRoot.IsEqual(vRoot) {
		log.DLogger.Error("verification root not match", zap.String("targetRoot", targetRoot.Hex()), zap.String("blockRoot", vRoot.Hex()), zap.Int("verificationLen", len(verifications)))
		return g_error.ErrVerificationRootNotMatch
	}

	return nil
}

// valid signer is verifiers
func verificationSignerInVerifiers(signer common.Address, verifiers []common.Address) bool {
	for _, v := range verifiers {
		if v.IsEqual(signer) {
			return true
		}
	}
	return false
}

func sameVote(votes []model.AbstractVerification) bool {
	var addressList []common.Address
	for _, vote := range votes {
		addressList = append(addressList, vote.GetAddress())
	}

	for i := 0; i < len(votes); i++ {
		address := addressList[0]
		addressList = addressList[1:]
		//fmt.Println(address, i)
		//fmt.Println(addressList, i)
		if address.InSlice(addressList) {
			return true
		}
	}
	return false
}

func validBlockHash(verifications []model.AbstractVerification, block model.AbstractBlock) error {
	for _, ver := range verifications {

		// valid block hash
		vBlockHashStr := block.Hash().Hex()
		//log.DLogger.Info("the block hash is:","ver",ver.GetBlockHash(),"block",vBlockHashStr)
		if ver.GetBlockHash() != vBlockHashStr {
			return g_error.ErrInvalidBlockHashInVotes
		}
	}
	return nil
}
