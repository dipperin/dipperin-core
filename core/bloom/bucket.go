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
)

// Configuration of bucket in IBLT.
type BucketConfig struct {
	//length of data by bytes
	DataLen uint
	//length of hash by bytes
	HashLen uint
}

// IBLT are made of buckets. The construction of bucket is slightly
// different from what the paper says, as there is no key here.
type Bucket struct {
	Count    int32
	DataSum  Data
	DataHash DataHash
}

// Create a new bucket configuration
func NewBucketConfig(dataLen, hashLen uint) BucketConfig {
	return BucketConfig{
		DataLen: dataLen,
		HashLen: hashLen,
	}
}

// Create a new bucket according to its configuraion.
func NewBucket(config BucketConfig) (b *Bucket) {
	return &Bucket{
		Count:    0,
		DataSum:  NewData(config.DataLen),
		DataHash: NewDataHash(config.HashLen),
	}
}

func (config BucketConfig) String() string {
	return fmt.Sprintf(
		"Data: %d\nHash: %d",
		config.DataLen,
		config.HashLen,
	)
}

// To add an element d in an empty bucket.
func (b *Bucket) Put(d Data) (*Bucket) {
	b.DataSum.SetBytes(d)

	h := hash(b.DataSum.Bytes())
	b.DataHash.SetBytes(h)
	b.Count = 1

	return b
}

func (b Bucket) String() string {
	return fmt.Sprintf("Count: \t%v \nDataSum: \t%v \nDatahash: \t%v\n",
		b.Count,
		b.DataSum.Hex(),
		b.DataHash.Hex(),
	)
}

// Call method subtract
func (z *Bucket) Subtract(a, b *Bucket) *Bucket {
	if a != nil && b != nil {
		return z.subtract(a, b)
	} else {
		return nil
	}
}

// Subtract set the value of z to the difference of a-b, and returns z
func (z *Bucket) subtract(a, b *Bucket) *Bucket {
	z.Count = a.Count - b.Count

	z.DataSum.Xor(a.DataSum, b.DataSum)

	z.DataHash.Xor(a.DataHash, b.DataHash)

	return z
}

// Bitwise XOR operation between two buckets.
func (z *Bucket) Xor(a, b *Bucket) *Bucket {
	if a != nil && b != nil {
		z.DataSum.Xor(a.DataSum, b.DataSum)

		z.DataHash.Xor(a.DataHash, b.DataHash)
	}
	return z
}

// To find out whether a bucket is pure. This is very useful in
// decoding a subtraction of two IBLTs.
func (b Bucket) isPure() bool {
	if b.Count == 1 || b.Count == -1 {
		h := hash(b.DataSum.Bytes())

		t := NewDataHash(uint(len(b.DataHash)))
		t.SetBytes(h)

		return t.IsEqual(b.DataHash)
	}
	return false
}

// Return whether a bucket is empty.
func (b Bucket) IsEmpty() bool {
	if b.Count == 0 &&
		b.DataSum.IsEmpty() &&
		b.DataHash.IsEmpty() {
		return true
	}

	return false
}

// Return whether two buckets are equal.
//func (b Bucket) IsEqual(a Bucket) bool {
//	return false
//}
