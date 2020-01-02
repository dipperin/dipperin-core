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
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"math/big"
)

// Predefined constants

// Constants left here for reference
const (
	// BucketNum represents the number of total buckets
	BucketNumBits = 12
	BucketNum     = 1 << BucketNumBits

	// BucketKeyHashLength is the number of bytes taken from
	// the hash of inserted RLP key
	BucketKeyHashLength = 6

	// BucketSerNumLength is the number of bytes used as index
	// if modified, change uint16 to corresponding types
	BucketSerNumLength = 2

	// BucketKeyLength represents the length of bytes used by keys
	BucketKeyLength = BucketKeyHashLength + BucketSerNumLength

	// BucketValueLength represents the length fo bytes used by keys
	BucketValueLength = 8
	DataLength        = BucketKeyLength + BucketValueLength

	// Number of used buckets in BucketNum of buckets
	BucketUsed = 4

	// Number of byte used in key hash
	KeyHashLength = 4

	// BloomBitLength is the number of bits in a classical bloom filter
	BloomLengthBits = 18
	BloomBitLength  = 1 << BloomLengthBits

	// BloomByteLength is the number of bytes in a classical bloom filter
	BloomByteLength = BloomBitLength / 8

	// number of bits per Digest-ed element
	BloomUsedBits = 7
)

// Errors
var (
	ErrKeyDuplicate = errors.New("insert key duplicated")
	ErrPtrEmpty     = errors.New("input nil pointer")
	ErrDecodeFailed = errors.New("decode IBLT failed")
)

func hash(b []byte) []byte {
	return crypto.Keccak256(b[:])
}

// Defined data types used in IBLT
type Data []byte

type DataHash []byte

func NewData(size uint) Data {
	s := make(Data, size)
	return s
}

func (k Data) SetBytes(b []byte) {
	copy(k, b)
}

func (k Data) Bytes() []byte {
	return k[:]
}

func (k Data) IsEqual(K Data) bool {
	return bytes.Equal(k[:], K[:])
}

func (k Data) IsEmpty() bool {
	return k.IsEqual(NewData(uint(len(k))))
}

func (k Data) Hex() string {
	return hexutil.Encode(k[:])
}

func (k Data) Big() *big.Int {
	return new(big.Int).SetBytes(k[:])
}

func (k Data) Hash() common.Hash {
	hh := common.Hash{}

	hh.SetBytes(hash(k))
	return hh
}

// Xor performs exclusive Or operation on a and b
// c = a ^ b, equivalent to c = b ^ a
// a, b, c are assumed to be the same length
func (k Data) Xor(a, b Data) Data {
	for i := 0; i < len(a); i++ {
		k[i] = a[i] ^ b[i]
	}
	return k
}

func NewDataHash(size uint) DataHash {
	return make(DataHash, size)
}

func (d DataHash) SetBytes(b []byte) {
	copy(d, b)
}

func (d DataHash) Bytes() []byte {
	return d[:]
}

func (d DataHash) IsEqual(K DataHash) bool {
	return bytes.Equal(d[:], K[:])
}

func (d DataHash) IsEmpty() bool {
	return d.IsEqual(NewDataHash(uint(cap(d))))
}

func (d DataHash) Hex() string {
	return hexutil.Encode(d[:])
}

func (d DataHash) Big() *big.Int {
	return new(big.Int).SetBytes(d[:])
}

func (d DataHash) Xor(a, b DataHash) DataHash {
	for i := 0; i < len(a); i++ {
		d[i] = a[i] ^ b[i]
	}
	return d
}

func (d DataHash) Lsh(n int) DataHash {
	for i := 0; i < n; i++ {
		d.lsh()
	}

	return d
}

func (d DataHash) Lsb() bool {
	var b byte
	if len(d) > 0 {
		b = d[0]
	}

	return b >= byte(128)
}

func (d DataHash) lsh() DataHash {
	for i := range d {
		d[i] <<= 1
		if i+1 < len(d) {
			car := d[i+1] & byte(128)
			car >>= 7
			d[i] |= car
		}
	}

	return d
}

// hashIndex is a helper function that returns the bloom bucket
// indexes according to the hash of input data.
func hashIndex(d Data, n uint, k uint) []uint {
	indexes := make([]uint, 0)
	//bloom_log.DLogger.Info("hashIndex")
	hashes := kHash(d[:], k)

	for _, hash := range hashes {
		//idx := hash & uint32(n-1)
		idx := uint(hash % uint32(n-1))

		indexes = append(indexes, idx)
	}

	distinct(indexes, n)
	return indexes
}

// Assuming hash function is 32 bytes length temporarily. It executes
// one hash instead of k, but it realises the effect of k hashes.
func kHash(bytes []byte, k uint) []uint32 {
	h := hash(bytes[:])

	var res []uint32

	//startAt := time.Now()
	i := uint(0)
	hLen := len(h)
	for j := uint(0); j < uint(hLen) && i < k; j += 4 {
		t := binary.BigEndian.Uint32(h[j : j+4])

		res = append(res, t)
		i++
	}
	//bloom_log.DLogger.Info("kHash", "use", time.Now().Sub(startAt), "h len", hLen)
	return res
}

// distinct modifies the input array such that the array has distinct values
// in the range from 0 to (BucketNum - 1), see below comments
func distinct(s []uint, n uint) {
	set := make(map[uint]bool)

	// impossible parameters
	if uint(len(s)) > n {
		return
	}

	for i := 0; i < len(s); i++ {
		idx := s[i]
		if !set[s[i]] {
			// not exists before
			// insert new k to the set
			set[s[i]] = true
		} else {
			// exists
			// while the given key exists
			for set[idx] {
				// increase the key
				idx++
				idx %= n
			}
			set[idx] = true
			s[i] = idx
		}
	}
}

type SortData []Data

func (k SortData) Len() int           { return len(k) }
func (k SortData) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k SortData) Less(i, j int) bool { return k[i].Less(k[j]) }

// Less returns true if caller is less then argument
// compare is byte wise
func (k Data) Less(comp Data) bool {
	for i, v := range k {
		if v != comp[i] {
			if v > comp[i] {
				return false
			} else {
				return true
			}
		}
	}
	return true
}

type SortDataHash []DataHash

func (k SortDataHash) Len() int           { return len(k) }
func (k SortDataHash) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k SortDataHash) Less(i, j int) bool { return k[i].Less(k[j]) }

// Less returns true if callee is less then argument
// compare is byte wise
func (d DataHash) Less(comp DataHash) bool {
	for i, v := range d {
		if v != comp[i] {
			if v > comp[i] {
				return false
			} else {
				return true
			}
		}
	}
	return true
}
