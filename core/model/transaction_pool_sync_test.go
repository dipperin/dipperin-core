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
	"fmt"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Alice sends block to Bob using IBLT
func TestTxPool_InvBloom(t *testing.T) {
	const (
		aliceNumTx  = 0
		bobNumTx    = 5
		commonNumTx = 5
	)
	//log.InitLogger(log.LvlError)

	tBloom := setUpBloomTest(aliceNumTx, bobNumTx, commonNumTx)

	// There are aliceNumTx + commonNumTx in Alice's block
	aliceMap := tBloom.aliceMap

	// There are bobNumTx + commonNumTx in Bob's transactions pool
	bobMap := tBloom.bobMap

	// Alice first constructs a bloom filter to initiate the transmission process
	aliceBloom := iblt.NewBloom(iblt.DeriveBloomConfig(len(aliceMap)))
	for k := range aliceMap {
		aliceBloom.Digest(k.Bytes())
	}

	aliceBloomRLP, err := rlp.EncodeToBytes(aliceBloom)
	assert.NoError(t, err)
	fmt.Printf("sending classical bloom using: %d bytes\n", len(aliceBloomRLP))
	// ###########################################
	// Alice sends its bloom filter to Bob to initiate estimator syncing
	// ###########################################

	// Bob receives Alice's bloom, decodes and reconstructs
	var bobReceivedBloom iblt.Bloom
	err = rlp.DecodeBytes(aliceBloomRLP, &bobReceivedBloom)
	assert.NoError(t, err)

	c := iblt.NewHybridEstimatorConfig()
	c.StrataConfig.IBLTConfig.BktConfig.DataLen = 16
	bobEstimator := iblt.NewHybridEstimator(c)

	// Bob constructs its estimator locally
	for k := range bobMap {
		if bobReceivedBloom.LookUp(k.Bytes()) {
			// construct the estimator only using the items that
			// are possible in Alice's
			bobEstimator.EncodeByte(k.Bytes())
		}
	}

	// Bob sends its estimator to Alice
	bEstimatorRLP, err := rlp.EncodeToBytes(bobEstimator)
	assert.NoError(t, err)

	// ###################################################
	// Network transport byte slice packet `bEstimatorRLP`
	// ###################################################

	// Alice receives the RLP bytes sent by Bob
	// Alice recovers the estimator
	var aReceivedEstimator iblt.HybridEstimator
	err = rlp.DecodeBytes(bEstimatorRLP, &aReceivedEstimator)
	assert.NoError(t, err)

	// Alice then constructs its estimator and decodes the estimator sent by Bob
	aliceEstimator := iblt.NewHybridEstimator(aReceivedEstimator.Config())
	for k := range aliceMap {
		aliceEstimator.EncodeByte(k.Bytes())
	}

	// Alice decodes and estimate the difference
	estimatedDiff := aliceEstimator.Decode(&aReceivedEstimator)
	fmt.Printf("estimated difference using: %d bytes\n", len(bEstimatorRLP))
	fmt.Println("estimated difference:", estimatedDiff)

	// Estimation ends here, start to sync data using IBLT, and IBLT is sized
	// by the estimated difference

	// Alice knows the set difference, then she constructs corresponding IBLT
	estimatedConfig := aliceEstimator.DeriveConfig(&aReceivedEstimator)
	fmt.Println("number of buckets", estimatedConfig.BucketNum)
	//estimatedConfig.BloomConfig = iblt.NewBloomConfig(0, 0)
	BloomConfig := iblt.DeriveBloomConfig(len(aliceMap))
	aliceInvBloom := iblt.NewGraphene(estimatedConfig, BloomConfig)
	//fmt.Println("-------------------", aliceInvBloom.InvBloom())
	var arr []*Transaction
	for k, v := range aliceMap {
		aliceInvBloom.InsertRLP(k, v)
		aliceInvBloom.Bloom().Digest(k.Bytes())
		arr = append(arr, v)
	}
	aliceBlockRLP, err := rlp.EncodeToBytes(arr)
	assert.NoError(t, err)

	// Alice then sends its IBLT to Bob
	aInvBloomRLP, err := rlp.EncodeToBytes(aliceInvBloom)
	assert.NoError(t, err)

	// ###################################################
	// Network transport byte slice packet `aInvBloomRLP`
	// ###################################################

	// Bob receives IBLT's RLP byte slice, and tries to decode
	var bReceivedInvBLoom iblt.Graphene
	err = rlp.DecodeBytes(aInvBloomRLP, &bReceivedInvBLoom)
	assert.NoError(t, err)

	// Bob constructs its IBLT
	bobInvBloom := iblt.NewGraphene(bReceivedInvBLoom.InvBloomConfig(), bReceivedInvBLoom.BloomConfig())

	for k, v := range bobMap {
		// TODO: should use classical bloom filter here
		if bReceivedInvBLoom.LookUp(k.Bytes()) {
			bobInvBloom.InsertRLP(k, v)
			// possible = map[hash]Txs
		} else {
			//
		}
	}

	tempGraphene := iblt.NewGraphene(bReceivedInvBLoom.InvBloomConfig(), bReceivedInvBLoom.BloomConfig())
	fmt.Println("==========tempGraphene.InvBloom()==========", bReceivedInvBLoom.InvBloom(), bobInvBloom.InvBloom())
	alice, bob, err := tempGraphene.InvBloom().Subtract(bReceivedInvBLoom.InvBloom(), bobInvBloom.InvBloom()).ListRLP()
	//res, err := aliceBloom.FilterListRLP(bobMap)
	fmt.Println("==========alice, bob==========", alice, bob)
	assert.NoError(t, err)

	assertTransactionInMap(t, bob, bobMap)

	assert.True(t, len(bob) <= bobNumTx)
	// for b in bob delete(possible)
	// alice + possible = block
	assert.Equal(t, aliceNumTx, len(alice))

	assertTransactionInMap(t, alice, aliceMap)

	fmt.Printf("sync %d entries using: %d bytes\n", aliceNumTx, len(aInvBloomRLP))
	fmt.Printf("naively sending the whole block's transaction uses: %d bytes\n", len(aliceBlockRLP))
	//fmt.Println(bobInvBloom.Bloom())
	fmt.Println(bobInvBloom.InvBloom().Config())
}
