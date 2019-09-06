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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"sort"
	"testing"
	"time"
)

var defaultInvBloomConfig = NewInvBloomConfig(1<<12, 4)

func TestInvBloom_Insert(t *testing.T) {
	b := NewInvBloom(defaultInvBloomConfig)
	testKey := b.NewData()
	testKey.SetBytes([]byte{1, 2, 3, 4})

	keyHash := NewDataHash(b.config.BktConfig.HashLen)
	keyHash.SetBytes(hash(testKey.Bytes()))

	b.Insert(testKey)

	// Indexes are generated with respect to the choice of hash function.
	indexes := hashIndex(testKey, b.config.BucketNum, b.config.BucketUsed)

	for _, idx := range indexes {
		assert.EqualValues(t, 1, b.buckets[idx].Count)

		assert.EqualValues(t, testKey, b.buckets[idx].DataSum)

		assert.EqualValues(t, keyHash, b.buckets[idx].DataHash)
	}
}

func TestInvBloom_Subtract(t *testing.T) {
	a := NewInvBloom(defaultInvBloomConfig)
	b := NewInvBloom(defaultInvBloomConfig)

	a.buckets[0] = a.NewBucket()
	a.buckets[0].Count = 2

	b.buckets[1] = a.NewBucket()
	b.buckets[1].Count = -1

	a.buckets[2] = a.NewBucket()
	a.buckets[2].Count = 3

	b.buckets[2] = b.NewBucket()
	b.buckets[2].Count = -2

	c := NewInvBloom(defaultInvBloomConfig).Subtract(a, b)

	assert.EqualValues(t, 2, c.buckets[0].Count)
	assert.EqualValues(t, 1, c.buckets[1].Count)
	assert.EqualValues(t, 5, c.buckets[2].Count)
}

func invBloomStats(b *InvBloom, t *testing.T) *map[int]uint {
	stats := make(map[int]uint)
	sum := uint(0)

	for _, bucket := range b.buckets {
		// Iterate over the array to count the `Count` field
		switch {
		case bucket == nil:
			stats[0]++
		default:
			stats[int(bucket.Count)]++
		}
	}

	// maintaining sorted key map
	// see URL: https://stackoverflow.com/questions/23330781/sort-go-map-values-by-keys

	var keys []int
	for i, count := range stats {
		sum += count
		keys = append(keys, int(i))
	}
	sort.Ints(keys)

	for _, k := range keys {
		fmt.Printf("Count: %d, Counts: %d\n", k, stats[k])
	}
	fmt.Printf("Sum: %d\n", sum)

	assert.EqualValues(t, b.config.BucketNum, sum)

	return &stats
}

func TestInvBloom_Delete(t *testing.T) {
	b := NewInvBloom(defaultInvBloomConfig)

	testKey := NewData(defaultBucketConfig.DataLen)
	testKey.SetBytes([]byte{1, 2, 3, 4})

	b.Insert(testKey)

	tmp := NewData(defaultBucketConfig.DataLen)
	tmp.SetBytes([]byte{10, 1, 1, 2, 3, 4})
	b.Insert(tmp)

	tmp = NewData(defaultBucketConfig.DataLen)
	tmp.SetBytes([]byte{2, 2, 2, 3, 4})
	b.Insert(tmp)

	tmp = NewData(defaultBucketConfig.DataLen)
	tmp.SetBytes([]byte{1, 2, 3, 4})
	assert.True(t, b.Delete(tmp))

	tmp = NewData(defaultBucketConfig.DataLen)
	tmp.SetBytes([]byte{1, 1, 1})

	assert.False(t, b.Delete(tmp))
	// Assumes safe deletion

	indexes := hashIndex(testKey, b.config.BucketNum, b.config.BucketUsed)

	for _, idx := range indexes {
		assert.True(t, b.buckets[idx].Count == 0)

		assert.True(t, b.buckets[idx].DataSum.IsEmpty())
	}

}

func generateRandomKV(size int) map[common.Hash]Data {
	rand.Seed(time.Now().Unix() + rand.Int63())

	KV := make(map[common.Hash]Data)

	p := NewData(defaultInvBloomConfig.BktConfig.DataLen)

	// generate random test bytes
	for i := 0; i < size; i++ {
		rand.Read(p)

		k := NewData(defaultInvBloomConfig.BktConfig.DataLen)
		k.SetBytes(p)

		KV[k.Hash()] = k
	}
	return KV
}
