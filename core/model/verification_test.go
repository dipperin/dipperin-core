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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/secp256k1"
	"github.com/stretchr/testify/assert"
	"testing"
)

var as = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"

func TestWitMsg_Valid(t *testing.T) {
	key1, err1 := crypto.HexToECDSA(as)
	assert.NoError(t, err1)

	sign, err := crypto.Sign(common.HexToHash("123").Bytes(), key1)
	assert.NoError(t, err)

	msg := WitMsg{AliceAddr, sign}

	err = msg.Valid(common.HexToHash("123").Bytes())
	assert.NoError(t, err)

	err = msg.Valid(common.HexToHash("456").Bytes())
	assert.Equal(t, "signature not valid", err.Error())

	msg = WitMsg{AliceAddr, []byte{123}}
	err = msg.Valid(common.HexToHash("123").Bytes())
	assert.Equal(t, secp256k1.ErrInvalidSignatureLen, err)
}

func TestNewVoteMsg(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	assert.NotNil(t, voteMsg)
}

func TestNewVoteMsgWithSign(t *testing.T) {
	_, err := NewVoteMsgWithSign(10, 1, common.HexToHash("100"), 1, func(hash []byte) ([]byte, error) {
		return nil, nil
	}, AliceAddr)
	assert.NoError(t, err)

	_, err = NewVoteMsgWithSign(10, 1, common.HexToHash("100"), 1, func(hash []byte) ([]byte, error) {
		return nil, ErrInvalidSig
	}, AliceAddr)
	assert.Equal(t, ErrInvalidSig, err)
}

func TestVoteMsg_String(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.String()
	assert.NotNil(t, result)
}

func TestVoteMsg_GetViewID(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.GetViewID()
	assert.Equal(t, uint64(0), result)
}

func TestVoteMsg_GetHeight(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.GetHeight()
	assert.Equal(t, uint64(10), result)
}

func TestVoteMsg_GetRound(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.GetRound()
	assert.Equal(t, uint64(1), result)
}

func TestVoteMsg_GetType(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.GetType()
	assert.Equal(t, VoteMsgType(1), result)
}

func TestVoteMsg_GetBlockId(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.GetBlockId()
	assert.NotNil(t, result)
}

func TestVoteMsg_GetAddress(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.GetAddress()
	assert.NotNil(t, result)
}

func TestVoteMsg_GetBlockHash(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.GetBlockHash()
	assert.NotNil(t, result)
}

func TestVoteMsg_Valid(t *testing.T) {
	voteMsg := CreateSignedVote(10, 1, common.HexToHash(as), 1)
	result := voteMsg.Valid()
	assert.NoError(t, result)
}

func TestVoteMsg_HaltedVoteValid(t *testing.T) {
	// RecoverAddressFromSig error
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), VerBootNodeVoteMessage)
	voteMsg.Witness = &WitMsg{Address: AliceAddr, Sign: []byte{123}}
	err := voteMsg.HaltedVoteValid([]common.Address{AliceAddr})
	assert.Equal(t, secp256k1.ErrInvalidSignatureLen, err)

	// VoteRecoverAddrError
	voteMsg = CreateSignedVote(10, 1, common.HexToHash("100"), VerBootNodeVoteMessage)
	voteMsg.Witness.Address = common.HexToAddress("123")
	err = voteMsg.HaltedVoteValid([]common.Address{AliceAddr})
	assert.Equal(t, VoteRecoverAddrError, err)

	// wrong vote type
	voteMsg = CreateSignedVote(10, 1, common.HexToHash("100"), PreVoteMessage)
	err = voteMsg.HaltedVoteValid([]common.Address{AliceAddr})
	assert.Equal(t, WrongVoteType, err)

	voteMsg = CreateSignedVote(10, 1, common.HexToHash("100"), VoteMessage)
	err = voteMsg.HaltedVoteValid([]common.Address{AliceAddr})
	assert.Equal(t, WrongVoteType, err)

	// CheckAddressIsVerifierBootNode failed
	voteMsg = CreateSignedVote(10, 1, common.HexToHash("100"), VerBootNodeVoteMessage)
	err = voteMsg.HaltedVoteValid([]common.Address{AliceAddr})
	assert.Equal(t, AddressIsNotVerifierBootNode, err)

	// CheckAddressIsCurrentVerifier failed
	chain_config.VerBootNodeAddress = []common.Address{AliceAddr}
	voteMsg = CreateSignedVote(10, 1, common.HexToHash("100"), AliveVerifierVoteMessage)
	err = voteMsg.HaltedVoteValid([]common.Address{BobAddr})
	assert.Equal(t, AddressIsNotCurrentVerifier, err)

	// no error
	err = voteMsg.HaltedVoteValid([]common.Address{AliceAddr})
	assert.NoError(t, err)
}

func TestCheckAddressIsVerifierBootNode(t *testing.T) {
	result := CheckAddressIsVerifierBootNode(common.HexToAddress("123"))
	assert.Equal(t, false, result)

	chain_config.VerBootNodeAddress = []common.Address{common.HexToAddress("123")}
	result = CheckAddressIsVerifierBootNode(common.HexToAddress("123"))
	assert.Equal(t, true, result)
}

func TestCheckAddressIsCurrentVerifier(t *testing.T) {
	result := CheckAddressIsCurrentVerifier(common.HexToAddress("123"), []common.Address{common.HexToAddress("123"), common.HexToAddress("456")})
	assert.Equal(t, true, result)

	result = CheckAddressIsCurrentVerifier(common.HexToAddress("123"), []common.Address{common.HexToAddress("456"), common.HexToAddress("789")})
	assert.Equal(t, false, result)
}

func TestVoteMsg_Hash(t *testing.T) {
	voteMsg := NewVoteMsg(10, 1, common.HexToHash("100"), 1)
	result := voteMsg.Hash()
	assert.NotNil(t, result)
}
