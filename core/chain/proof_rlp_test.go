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


package chain

//import (
//	"github.com/ethereum/go-ethereum/rlp"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestProofRLP(t *testing.T) {
//	setUpInterlinkTest(t)
//
//	p, err := GetInfixProof(reader, 300, 3, 3, 20)
//	assert.NoError(t, err)
//	err = p.Valid(reader)
//	assert.NoError(t, err)
//	b, err := rlp.EncodeToBytes(p)
//	assert.NoError(t, err)
//
//	var proof *Proof
//	err = rlp.DecodeBytes(b, &proof)
//	assert.NoError(t, err)
//
//	bb, err := rlp.EncodeToBytes(proof)
//	assert.NoError(t, err)
//	assert.Equal(t, b, bb)
//
//	err = proof.Valid(reader)
//	assert.NoError(t, err)
//}
//
//func TestProofsRLP(t *testing.T) {
//	setUpInterlinkTest(t)
//
//	p, err := GetInfixProof(reader, 300, 3, 3, 20)
//	assert.NoError(t, err)
//	err = p.Valid(reader)
//	assert.NoError(t, err)
//	b, err := rlp.EncodeToBytes(Proofs{*p})
//	assert.NoError(t, err)
//
//	var proofs Proofs
//	err = rlp.DecodeBytes(b, &proofs)
//	assert.NoError(t, err)
//
//	bb, err := rlp.EncodeToBytes(proofs)
//	assert.NoError(t, err)
//	assert.Equal(t, b, bb)
//
//	err = proofs[0].Valid(reader)
//	assert.NoError(t, err)
//}
