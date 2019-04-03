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

// Test the performance of bloom and IBLT when the data sizes are different.

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func benchmarkBloomInsertDeleteN(n int) bool {
	samples := generateRandomKV(n)

	bloom := NewInvBloom(defaultInvBloomConfig)

	bloom.InsertMap(samples)

	for _, v := range samples {
		bloom.Delete(v)
	}

	for _, bucket := range bloom.buckets {
		if bucket != nil && !bucket.IsEmpty() {
			// after deletion, IBLT is not empty
			return false
		}
	}

	return true
}

// inserts n KV pairs, and checks whether it decodes successfully
//func benchmarkInvBloomListEntriesN(n int) bool {
//	samples := generateRandomKV(n)
//
//	bloom := NewInvBloom(defaultInvBloomConfig).InsertMap(samples)
//
//	results := bloom.ListEntries()
//
//	count := 0
//	for k := range samples {
//		if samples[k].IsEqual(results[k]) {
//			count++
//		}
//	}
//
//	return count == len(samples)
//}

// benchmarks the efficiency for insert and delete KV pairs
func BenchmarkInvbloom_InsertDelete1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBloomInsertDeleteN(1)
	}
}

func BenchmarkInvbloom_InsertDelete10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBloomInsertDeleteN(10)
	}
}

func BenchmarkInvbloom_InsertDelete100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBloomInsertDeleteN(100)
	}
}

func BenchmarkInvbloom_InsertDelete1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBloomInsertDeleteN(1000)
	}
}

func BenchmarkInvbloom_InsertDelete10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBloomInsertDeleteN(10000)
	}
}

func BenchmarkInvbloom_InsertDelete100000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBloomInsertDeleteN(100000)
	}
}

const TestSize = 5000

func TestInvBloom_diff(t *testing.T) {
	const (
		intersect = 1000
		aSize     = 1000
		bSize     = aSize
	)

	a := NewInvBloom(defaultInvBloomConfig)
	b := NewInvBloom(defaultInvBloomConfig)

	whole := generateRandomKV(intersect + aSize + bSize)

	aMap := make(map[common.Hash]Data)
	bMap := make(map[common.Hash]Data)

	cnt := 0
	for k, v := range whole {
		if cnt < intersect {
			a.Insert(v)
			b.Insert(v)
		} else if cnt < intersect+aSize {
			a.Insert(v)
			aMap[k] = v
		} else if cnt < intersect+aSize+bSize {
			b.Insert(v)
			bMap[k] = v
		}
		cnt++
	}

	//invBloomStats(a, t)
	//invBloomStats(b, t)

	a.Subtract(a, b)
	//invBloomStats(a, t)

	aAlice := make(map[common.Hash]Data)
	aBob := make(map[common.Hash]Data)

	success := a.Decode(aAlice, aBob)

	for k, v := range aMap {
		assert.EqualValues(t, v, aAlice[k])
	}

	for k, v := range bMap {
		assert.EqualValues(t, v, aBob[k])
	}

	assert.True(t, success)
}

func TestInvBloom_InsertDelete(t *testing.T) {
	b := NewInvBloom(defaultInvBloomConfig)

	m := generateRandomKV(TestSize)
	b.InsertMap(m)
	//invBloomStats(b, t)

	// delete these bytes
	for _, v := range m {
		b.Delete(v)
	}

	//invBloomStats(b, t)
	for _, bucket := range b.buckets {
		if bucket != nil {
			// asserts IsEmpty
			assert.True(t, bucket.DataSum.IsEmpty())
		}
	}
}

func TestInvBloom_FilterDecode(t *testing.T) {
	var aKey []Data

	const (
		intersect = 100
		aSize     = 100
		bSize     = 100
	)

	a := NewGraphene(defaultInvBloomConfig, defaultConfig)

	whole := generateRandomKV(intersect + aSize + bSize)

	// construct distinctive KV pairs
	bMap := make(map[common.Hash]Data)

	cnt := 0
	for k, v := range whole {
		if cnt < intersect {
			a.Insert(v)
			//aKey = append(aKey, k)
			bMap[k] = v
			//aValue = append(aValue, v)
		} else if cnt < intersect+aSize {
			a.Insert(v)
			aKey = append(aKey, v)
		} else if cnt < intersect+aSize+bSize {
			bMap[k] = v
		}
		cnt++
	}

	//invBloomStats(&a, t)

	//bloomStats(a.bloom, t)

	// guess what element in b might in `a`
	// so aAlice is the distinct element `a` only has
	aAlice := a.FilterDecode(bMap)

	assert.Equal(t, aSize, len(aAlice))
	for _, k := range aKey {
		assert.NotNil(t, aAlice[k.Hash()])
	}
	//bloomStats(a.bloom, t)
}

// inserts n KV pairs nn times, to extract the threshold of success decode
// larger nn ensures more precise result
//func benchmarkInvBloomListEntriesNN(n int, nn int) float32 {
//	numerator := 0
//	denominator := 0
//
//	for i := 0; i < nn; i++ {
//		if benchmarkInvBloomListEntriesN(n) {
//			numerator++
//		}
//		denominator++
//	}
//
//	return float32(numerator) / float32(denominator)
//}

// This test takes a long time about 20s, not recommended to run every build
//func TestInvBloom_ListEntries(t *testing.T) {
//	assert.True(t, benchmarkInvBloomListEntriesN(100))
//	assert.True(t, benchmarkInvBloomListEntriesN(1000))
//	fmt.Println(benchmarkInvBloomListEntriesNN(1575, 1000))
//}
