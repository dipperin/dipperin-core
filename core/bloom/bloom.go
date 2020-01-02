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
	"github.com/dipperin/dipperin-core/common/hexutil"
	"math/big"
)

// Struct of bloom comprises a bloom of byte slice and a config
type Bloom struct {
	// The data of bloom stored in byte slice
	bloom []byte
	// Configuration of bloom
	config BloomConfig
}

// Configuration of bloom
type BloomConfig struct {
	// Length of bloom on byte (not en bit)
	BloomByteLength uint
	// Number of Hash functions
	BloomBits uint
}

// Create a new BloomConfig where the input len is the log2 of the
// length of bloom on bit
func NewBloomConfig(len, bits uint) BloomConfig {
	return BloomConfig{
		BloomBits: bits,
		// Conversion from bits to bytes
		BloomByteLength: 1 << (len - 3),
	}
}

// Proposed config of Bloom according to the paper
// given the number of elements.
func DeriveBloomConfig(count int) BloomConfig {
	l := uint(0)
	count *= 30
	for count != 0 {
		l++
		count >>= 1
	}
	return NewBloomConfig(l, 4)
}

func NewBloom(config BloomConfig) *Bloom {
	return &Bloom{
		bloom:  make([]byte, config.BloomByteLength),
		config: config,
	}
}

// SetBytes sets the content of b to the given bytes.
func (b *Bloom) SetBytes(d []byte) {
	if len(b.bloom) < len(d) {
		d = d[uint(len(d))-b.config.BloomByteLength:]
	}
	copy(b.bloom[b.config.BloomByteLength-uint(len(d)):], d)
}

// IsEqual tests whether the caller Bloom object is identical
// to the parameter B.
func (b Bloom) IsEqual(B *Bloom) bool {
	return bytes.Equal(b.bloom[:], B.bloom[:])
}

// bitwise operation Or.
func (b *Bloom) Or(x, y *Bloom) *Bloom {
	for i := range b.bloom {
		b.bloom[i] = x.bloom[i] | y.bloom[i]
	}
	return b
}

func (b Bloom) Hex() string {
	return hexutil.Encode(b.bloom[:])
}

// For certain value k, execute the required hashes, and set
// the value of the corresponding bloom bit to 1
func (b *Bloom) Digest(k []byte) *Bloom {
	//bloom_log.DLogger.Info("bloom Digest")
	hashes := kHash(k[:], b.config.BloomBits)

	for _, hash := range hashes {
		// in case the value of idx is greater than the size of the bloom
		idx := hash & uint32(b.config.BloomByteLength*8-1)

		b.SetAt(idx)
	}

	return b
}

// SetAt sets certain bit to 1 in big-endian
func (b *Bloom) SetAt(idx uint32) {
	op := uint8(1)

	if idx < uint32(b.config.BloomByteLength*8) {
		byteIdx := uint32(b.config.BloomByteLength) - (idx / 8) - 1
		bitIdx := idx % 8
		b.bloom[byteIdx] |= op << bitIdx
	}
}

// Return the value of the corresponding bloom bit
func (b *Bloom) LookAt(idx uint32) bool {
	op := uint8(1)
	if idx < uint32(b.config.BloomByteLength*8) {
		byteIdx := uint32(b.config.BloomByteLength) - (idx / 8) - 1
		bitIdx := idx % 8

		op <<= bitIdx

		return op == (b.bloom[byteIdx] & op)
	}

	return false
}

// Big converts b to a big integer.
func (b Bloom) Big() *big.Int {
	return new(big.Int).SetBytes(b.bloom[:])
}

// Find out whether an element k is in the set according to the bloom.
// Notice that there are possibilities of false positive.
func (b Bloom) LookUp(k []byte) bool {
	hashes := kHash(k[:], b.config.BloomBits)
	//startAt := time.Now()
	for _, hash := range hashes {

		idx := hash & uint32(b.config.BloomByteLength*8-1)

		if !b.LookAt(idx) {
			return false
		}
	}
	//bloom_log.DLogger.Info("bloom look up", "use", time.Now().Sub(startAt), "hash len", len(hashes))

	return true
}
