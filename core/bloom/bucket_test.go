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
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultBucketConfig = NewBucketConfig(8, 4)

func TestBucket_IsEmpty(t *testing.T) {
	b := NewBucket(defaultBucketConfig)
	assert.True(t, b.IsEmpty())
}

func TestBucket_subtract(t *testing.T) {
	a := NewBucket(defaultBucketConfig)
	b := NewBucket(defaultBucketConfig)

	a.Count = 1

	a.DataSum = NewData(defaultBucketConfig.DataLen)

	b.Subtract(b, a)

	assert.Equal(t, -a.Count, b.Count)
	assert.EqualValues(t, a.DataSum, b.DataSum)
}

func TestBucket_Empty(t *testing.T) {
	a := NewData(defaultBucketConfig.DataLen)
	assert.True(t, a.IsEmpty())

	a = NewData(defaultBucketConfig.DataLen)
	a.SetBytes([]byte{1, 2, 3})

	assert.False(t, a.IsEmpty())
}

func TestInvBloom_String(t *testing.T) {
	config := NewInvBloomConfig(1000, 6)
	fmt.Println(config)
}

func TestBucket_IsPure(t *testing.T) {
	bkt := NewBucket(defaultInvBloomConfig.BktConfig)
	k := NewData(defaultInvBloomConfig.BktConfig.DataLen)
	k.SetBytes([]byte{1, 2, 3})

	bkt.Put(k)

	assert.True(t, bkt.isPure())
}

func TestBucket_Put(t *testing.T) {
	bucketConfig := NewBucketConfig(16, 16)
	bucket := NewBucket(bucketConfig)
	data := make(Data, 8)
	data[5] = 86
	bucket.Put(data)
}

func TestBucket_Subtract(t *testing.T) {
	bucketConfig := NewBucketConfig(32, 16)
	bucket1 := NewBucket(bucketConfig)
	bucket2 := NewBucket(bucketConfig)
	data1 := make(Data, 8)
	data1[1] = 117
	data1[2] = 66
	data2 := make(Data, 8)
	data2[5] = 0xf7
	bucket1.Put(data1)
	bucket2.Put(data2)
	bucket1.Subtract(bucket1, bucket2)
	data3 := make(Data, 4)
	assert.EqualValues(t, data3, data3.Bytes())
}

func TestBucketConfig_String(t *testing.T) {
	bucketConfig := NewBucketConfig(32, 16)
	assert.Equal(t, bucketConfig.String(), "Data: 32\nHash: 16")
}

func TestBucket_String(t *testing.T) {
	bucketConfig := NewBucketConfig(16, 16)
	bucket := NewBucket(bucketConfig)

	assert.Equal(t, bucket.String(), "Count: \t0 \nDataSum: \t0x00000000000000000000000000000000 \nDatahash: \t0x00000000000000000000000000000000\n")
}
