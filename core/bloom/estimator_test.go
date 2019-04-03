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


package iblt

import (
	"github.com/dipperin/dipperin-core/common"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEstimator_IBLT(t *testing.T) {
	const (
		intersect = 1
		aSize     = 1
		bSize     = aSize
	)

	aliceEstimator := NewHybridEstimator(defaultHybridEstimatorConfig)
	bobEstimator := NewHybridEstimator(defaultHybridEstimatorConfig)

	whole := generateRandomKV(intersect + aSize + bSize)
	// save their data set for future use
	aliceSet := make(map[common.Hash]Data)
	bobSet := make(map[common.Hash]Data)

	// construct distinctive KV pairs
	cnt := 0
	for k, v := range whole {
		if cnt < intersect {
			aliceEstimator.Encode(v)
			bobEstimator.Encode(v)
			aliceSet[k] = v
			bobSet[k] = v
		} else if cnt < intersect+aSize {
			aliceEstimator.Encode(v)
			aliceSet[k] = v
		} else if cnt < intersect+aSize+bSize {
			bobEstimator.Encode(v)
			bobSet[k] = v
		}
		cnt++
	}

	// bobEstimator sends to aliceEstimator, aliceEstimator has the result
	estimatedDiff := aliceEstimator.Decode(bobEstimator)

	// coefficient here is 1.6
	c := NewInvBloomConfig(uint(float32(estimatedDiff)*1.6), 4)
	fmt.Println("estimatedDiff:", estimatedDiff)
	fmt.Println("IBLT size:", c.BucketNum)

	// alice sets up the IBLT based on the estimated buckets number
	aliceIBLT := NewInvBloom(c)

	for _, v := range aliceSet {
		aliceIBLT.Insert(v)
	}

	// alice sends aliceIBLT to bob
	// bob sets up an aliceIBLT based on the parameters given in
	// the IBLT sent by alice
	aliceRLP, err := rlp.EncodeToBytes(aliceIBLT)
	assert.NoError(t, err)
	fmt.Println("alice IBLT length", len(aliceRLP))
	fmt.Println("There are", uint(len(aliceSet))*aliceIBLT.config.BktConfig.DataLen, "bytes in total")

	bobIBLT := NewInvBloom(aliceIBLT.config)

	// bob constructs and inserts its aliceIBLT
	for _, v := range bobSet {
		bobIBLT.Insert(v)
	}

	bobHas, aliceHas := make(map[common.Hash]Data), make(map[common.Hash]Data)

	// bob tries to decode
	temp := NewInvBloom(bobIBLT.config)
	temp.Subtract(bobIBLT, aliceIBLT)
	success := temp.Decode(bobHas, aliceHas)
	assert.True(t, success)

	assert.Equal(t, bSize, len(bobHas))
	assert.Equal(t, aSize, len(aliceHas))

	// aliceHas stores the data in alice side only
	// bobHas stores the data in bob side only
	for k, v := range aliceHas {
		assert.True(t, aliceSet[k].IsEqual(v))
	}
}
