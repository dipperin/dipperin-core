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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/log"
)

// Create a new InvBloom Config according two the config of buckets
// in the context of transmission of transactions
func NewInvBloomConfig(bkt uint, used uint) InvBloomConfig {
	config := InvBloomConfig{
		BucketNum:  bkt,
		BucketUsed: used,
	}

	config.KeyLen = 4
	config.SerNumLen = 2
	config.ValLen = 120

	// The total length of value on bytes
	dataLen := config.KeyLen + config.ValLen + config.SerNumLen
	// TODO: adjust parameters
	config.BktConfig = NewBucketConfig(dataLen, 4)

	return config
}

// Configuration of IBLT
type InvBloomConfig struct {
	// The configuration of buckets of IBLT
	BktConfig BucketConfig

	// BucketNum represents the number of total buckets
	BucketNum uint

	// Number of used buckets in BucketNum of buckets
	BucketUsed uint

	// Number of bytes reserved to record the key
	KeyLen    uint

	// Number of bytes reserved to record the value
	ValLen    uint

	// Number of bytes reserved to record the section number
	SerNumLen uint
}

// The structure of IBLT
type InvBloom struct {
	// todo Each struct with a config will take up a lot of memory
	config  InvBloomConfig
	buckets []*Bucket
	salt    uint8
}

// Create a new IBLT
func NewInvBloom(config InvBloomConfig) *InvBloom {
	return &InvBloom{
		config:  config,
		buckets: make([]*Bucket, config.BucketNum),
	}
}

// deprecated
//func (b *InvBloom) ListEntries() map[common.Hash]Data {
//	return b.listEntries()
//}

//func (config InvBloomConfig) String() string {
//	return fmt.Sprintf(
//		"%d out of %d buckets\n[Key(%d), SerNum(%d), Val(%d)]\nBucket config:\n%v",
//		config.BucketUsed,
//		config.BucketNum,
//		config.KeyLen,
//		config.SerNumLen,
//		config.ValLen,
//		config.BktConfig,
//	)
//}

//func (config InvBloomConfig) String() string {
//	return fmt.Sprintf(`{
//	BucketUsed: %v,
//	BucketNum: %v,
//	KeyLen: %v,
//	SerNumLen: %v,
//	ValLen: %v,
//	BktConfig: %v
//}`, 	config.BucketUsed,
//		config.BucketNum,
//		config.KeyLen,
//		config.SerNumLen,
//		config.ValLen,
//		config.BktConfig,
//	)
//}

func (config InvBloomConfig) String() string {
	return util.StringifyJson(config)
}

// Call the insert function
func (b *InvBloom) Insert(x Data) {
	b.insert(x)
}

// Insert inserts the data into certain buckets in the InvBloom.
func (b *InvBloom) insert(x Data) {
	idx := hashIndex(x, b.config.BucketNum, b.config.BucketUsed)
	bkt := b.NewBucket()
	bkt.Put(x)

	b.operateBucket(idx, bkt, 1)
}

// For each element call the Insert function
func (b *InvBloom) InsertMap(m map[common.Hash]Data) *InvBloom {
	for _, v := range m {
		b.Insert(v)
	}
	return b
}

// For each element call the insert function
func (b *InvBloom) insertMap(m map[common.Hash]Data) *InvBloom {
	for _, v := range m {
		b.insert(v)
	}
	return b
}

// Delete a data from the corresponding buckets
func (b *InvBloom) Delete(k Data) bool {
	idx := hashIndex(k, b.config.BucketNum, b.config.BucketUsed)

	for _, i := range idx {
		if b.buckets[i] == nil || b.buckets[i].IsEmpty() {
			return false
		}
	}

	b.delete(k)

	return true
}

//func (b InvBloom) String() string {
//	s := string("InvBloom:\n")
//	for i := 0; i < int(b.config.BucketNum); i++ {
//		if b.buckets[i] != nil && !b.buckets[i].IsEmpty() {
//			s += fmt.Sprintf("At [%d] \n%v", i, b.buckets[i])
//		}
//	}
//	return s
//}

func (b InvBloom) String() string {
	return fmt.Sprintf(`{
config: %v, 
buckets: %v, 
salt: %v
}`, b.config, b.buckets, b.salt)
}

// Since this style of operator should support syntax like
// a.Subtract(a, b), attention should be paid to the value
// change when then are assigned, consider all the situations
//
// with the help of Linxin Liu at 2018-08-03
//
// subtract sets z to the difference of a-b and returns z
func (z *InvBloom) Subtract(a, b *InvBloom) *InvBloom {
	for i := range a.buckets {
		if a.buckets[i] == nil && b.buckets[i] == nil {
			// both nil, nil to z
			z.buckets[i] = nil
		} else if a.buckets[i] != nil && b.buckets[i] == nil {
			// a has value, z = a
			z.buckets[i] = a.buckets[i]
		} else if a.buckets[i] == nil && b.buckets[i] != nil {
			// b has value, z = -b
			if z.buckets[i] == nil {
				z.buckets[i] = z.NewBucket()
			}
			z.buckets[i].Subtract(z.NewBucket(), b.buckets[i])
		} else {
			// z = a - b
			if z.buckets[i] == nil {
				z.buckets[i] = z.NewBucket()
			}
			z.buckets[i].Subtract(a.buckets[i], b.buckets[i])
		}
	}
	return z
}

func (b *InvBloom) NewBucket() *Bucket {
	return NewBucket(b.config.BktConfig)
}

// insert and delete are similar functions, the only difference
// is the increment/decrement on Count.
// Attention, these two functions are private functions,
// safety checks are not included.

// delete operates similar to insert, however, when classical bloom
// is used, more deletions worse filter effect at the receiver side.
func (b *InvBloom) delete(x Data) {
	idx := hashIndex(x, b.config.BucketNum, b.config.BucketUsed)

	bkt := b.NewBucket()
	bkt.Put(x)

	b.operateBucket(idx, bkt, -1)
}

// operateBucket modifies the buckets in b
// buckets are specified at the location idx,
// bkt is the Bucket pointer to be Xor with
// c is the constant updating the Count field.
func (b *InvBloom) operateBucket(idx []uint, bkt *Bucket, c int32) *InvBloom {
	for _, i := range idx {
		if b.buckets[i] == nil {
			b.buckets[i] = b.NewBucket()
		}

		b.buckets[i].Count += c
		b.buckets[i].Xor(b.buckets[i], bkt)
	}
	return b
}

// deprecated, if subtraction of InvBloom is introduced, return contains
// kV pairs with Count 1.
// TODO optimize to use queue
//func (b *InvBloom) listEntries() map[common.Hash]Data {
//	var pureList []Bucket
//	//output := map[common.Hash]Data{}{}
//	output := make(map[common.Hash]Data)
//	for {
//		// go style do ... while
//		// iterate through the buckets until all the pure
//		// buckets are extracted, time complexity O(n*long(n))
//		for _, bucket := range b.buckets {
//			if bucket != nil && bucket.isPure() {
//				pureList = append(pureList, *bucket)
//
//				output[bucket.DataSum.Hash()] = bucket.DataSum
//				b.delete(bucket.DataSum)
//			}
//		}
//
//		if len(pureList) <= 0 {
//			break
//		}
//		// IsEmpty list
//		pureList = []Bucket{}
//	}
//	return output
//}

// Decode should be called after InvBloom subtraction
// alice stores the KV pairs that b, the caller has,
// bob stores the KV pairs that b, the caller has not.
//
// Example:
// A: {"foo":"a1", "bar":"a2", "foo1":"common1", "bar1":"common2"}
// B: {"baz":"b1", "qux":"b2", "foo1":"common1", "bar1":"common2"}
//
// C := A - B
//
// C.Decode(alice, bob)
//
// alice: {"foo":"a1", "bar":"a2"}
// bob: {"baz":"b1", "qux":"b2"}
//
// Their common KV pairs are canceled, only when it returns true
// could their distinct pairs be Decode-ed. If it returns false,
// the partial pairs are not guaranteed to be correct.
//
// Attention: Decode modifies the callee object. If one needs a
// non-corrupted version of decode, call Decode by a copy of that.
// TODO, optimize to use queue
func (b *InvBloom) Decode(alice, bob map[common.Hash]Data) bool {
	var pureList []Bucket
	loopCount := 0
	for {
		// todo only for debug
		loopCount++
		if loopCount > 200000 {
			log.Error("infinite loop in InvBloom Decode")
			//panic("infinite loop in InvBloom Decode")
			return false
		}
		for i := 0; i < len(b.buckets); i++ {
			bucket := b.buckets[i]
			if bucket != nil && bucket.isPure() {
				pureList = append(pureList, *bucket)
				k := b.NewData()
				k.SetBytes(bucket.DataSum)

				c := bucket.Count
				if bucket.Count > 0 {
					if alice != nil {
						alice[k.Hash()] = k
					}
				} else {
					if bob != nil {
						bob[k.Hash()] = k
					}
				}

				idx := hashIndex(k, b.config.BucketNum, b.config.BucketUsed)
				// deleting
				bkt := b.NewBucket()
				bkt.Put(k)

				b.operateBucket(idx, bkt, -c)
			}
		}
		if len(pureList) <= 0 {
			break
		}
		pureList = []Bucket{}
	}

	for i := range b.buckets {
		if b.buckets[i] != nil && !b.buckets[i].IsEmpty() {
			log.Error("invertible bloom lookup table decode failed", "failed bucket index", i)
			return false
		}
	}
	return true
}

// Create a new data with the length of value in buckets
func (b InvBloom) NewData() Data {
	return NewData(b.config.BktConfig.DataLen)
}

// Return the configuration of IBLT
func (b InvBloom) Config() InvBloomConfig {
	return b.config
}
