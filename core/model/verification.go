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


package model

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log/ver_halt_check_log"
	"time"
)

//go:generate stringer -type=VoteMsgType
type VoteMsgType byte

const (
	PreVoteMessage           VoteMsgType = 0
	VoteMessage              VoteMsgType = 1
	VerBootNodeVoteMessage   VoteMsgType = 2
	AliveVerifierVoteMessage VoteMsgType = 3
)

var (
	AddressIsNotVerifierBootNode = errors.New("the Address isn't verifier boot node")
	VoteRecoverAddrError         = errors.New("the aliveVerifierVote recover Address is invalid")
	AddressIsNotCurrentVerifier  = errors.New("the recover Address isn't current verifier")
	WrongVoteType                = errors.New("wrong vote type")
)

type SignHashFunc func(hash []byte) ([]byte, error)

type WitMsg struct {
	Address common.Address `json:"address"`
	Sign    []byte         `json:"sign"`
}

func (witMsg *WitMsg) Valid(dataHash []byte) error {
	pubKey, err := crypto.SigToPub(dataHash, witMsg.Sign)
	if err != nil {
		return err
	}
	address := cs_crypto.GetNormalAddress(*pubKey)
	if address.IsEqual(witMsg.Address) {
		return nil
	} else {
		return errors.New("signature not valid")
	}
}

//Check whether the signature is correct and match the corresponding address to exclude the case that one person votes twice
type BftBlockBody struct {
	Txs    []*Transaction `json:"transactions"`
	Vers   []*VoteMsg     `json:"commit_msg"`
	Inters InterLink      `json:"interlinks"`
}

type VoteMsg struct {
	Height    uint64
	Round     uint64
	BlockID   common.Hash
	VoteType  VoteMsgType
	Timestamp time.Time
	Witness   *WitMsg
}

func NewVoteMsg(height, round uint64, blockID common.Hash, voteType VoteMsgType) *VoteMsg {
	return &VoteMsg{
		Height:    height,
		Round:     round,
		BlockID:   blockID,
		VoteType:  voteType,
		Timestamp: time.Now(),
		Witness:   &WitMsg{},
	}
}

// Create a new signature VoteMsg
func NewVoteMsgWithSign(height, round uint64, blockID common.Hash, voteType VoteMsgType, signFunc SignHashFunc, signAddress common.Address) (*VoteMsg, error) {

	vote := &VoteMsg{
		Height:    height,
		Round:     round,
		BlockID:   blockID,
		VoteType:  voteType,
		Timestamp: time.Now(),
	}

	sign, err := signFunc(vote.Hash().Bytes())
	if err != nil {
		return nil, err
	}
	vote.Witness = &WitMsg {
		Address: signAddress,
		Sign:    sign,
	}

	return vote, nil
}

func (v VoteMsg) String() string {
	return util.StringifyJson(v)
}

func (v VoteMsg) GetViewID() uint64 {
	return 0
}

func (v VoteMsg) GetHeight() uint64 {
	return v.Height
}

func (v VoteMsg) GetRound() uint64 {
	return v.Round
}

func (v VoteMsg) GetType() VoteMsgType {
	return v.VoteType
}

func (v VoteMsg) GetBlockId() common.Hash {
	return v.BlockID
}

func (v VoteMsg) GetAddress() common.Address {
	return v.Witness.Address
}

func (v VoteMsg) GetBlockHash() string {
	return v.BlockID.Hex()
}

func (v VoteMsg) Valid() error {

	// deep copy
	cp := v
	return v.Witness.Valid(cp.Hash().Bytes())
}

func (v VoteMsg) HaltedVoteValid(verifiers []common.Address) error {

	if v.GetType() == VoteMessage || v.GetType() == PreVoteMessage {
		return WrongVoteType
	}

	recoverAddress, err := cs_crypto.RecoverAddressFromSig(v.Hash(), v.Witness.Sign)
	if err != nil {
		ver_halt_check_log.Error("recover Address error from witness")
		return err
	}

	if recoverAddress != v.Witness.Address {
		return VoteRecoverAddrError
	}

	if v.GetType() == VerBootNodeVoteMessage {
		checkResult := CheckAddressIsVerifierBootNode(recoverAddress)
		if !checkResult {
			ver_halt_check_log.Warn("the Address isn't verifier boot node")
			return AddressIsNotVerifierBootNode
		}
	} else if v.GetType() == AliveVerifierVoteMessage {
		if !CheckAddressIsCurrentVerifier(recoverAddress, verifiers) {
			ver_halt_check_log.Warn("the Address isn't current verifier")
			return AddressIsNotCurrentVerifier
		}
	}

	return nil
}

func CheckAddressIsVerifierBootNode(address common.Address) bool {

	for _, addr := range chain_config.VerBootNodeAddress {
		if address == addr {
			return true
		}
	}
	return false
}

func CheckAddressIsCurrentVerifier(address common.Address, verifiers []common.Address) bool {

	for _, addr := range verifiers {
		if address == addr {
			return true
		}
	}
	return false
}

func (v VoteMsg) Hash() common.Hash {
	v.Witness = nil
	return common.RlpHashKeccak256(v)
}
