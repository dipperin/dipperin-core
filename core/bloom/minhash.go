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
	"sort"
)

// The configuration of HashPool
type HashPoolConfig struct {
	// length of Hash by bytes
	HashLen uint

	// K is the number of elements with smallest hash values
	// chosen for the estimation of difference
	K uint
}

// Create a new hashpool
func NewHashPoolConfig(l, k uint) HashPoolConfig {
	return HashPoolConfig{
		HashLen: l,
		K:       k,
	}
}

// HashPool stores sorted hashes, the first K hashes will be
// selected for comparison in order to estimate the set difference
type HashPool struct {
	hashes SortDataHash
	config HashPoolConfig
}

// Create an empty HashPool
func NewHashPool(config HashPoolConfig) *HashPool {
	return &HashPool{
		hashes: make(SortDataHash, 0),
		config: config,
	}
}

// Add the hash of a new data in the hashpool
func (p *HashPool) Encode(d Data) {
	h := hash(d)
	dh := NewDataHash(p.config.HashLen)
	dh.SetBytes(h)

	p.hashes = append(p.hashes, dh)
}

// MinHash chooses and returns K smallest hashes
func (p HashPool) MinHash(k uint) SortDataHash {
	sort.Sort(p.hashes)

	if k > uint(p.hashes.Len()) {
		k = uint(p.hashes.Len())
	}

	res := make(SortDataHash, k)
	copy(res, p.hashes[:k])

	return res
}

// Get the minimum between the number of hashes
// in the pool and K
func (p HashPool) GetK() int {
	k := int(p.config.K)
	if k < p.hashes.Len() {
		return k
	} else {
		return p.hashes.Len()
	}
}

// Comp returns the number of common hashes among K hashes
func (p HashPool) Comp(pp *HashPool) int {
	a := p.MinHash(p.config.K)
	b := pp.MinHash(p.config.K)

	cnt := 0

	// TODO: two loops could be optimized
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i].IsEqual(b[j]) {
				cnt++
			}
		}
	}

	return cnt
}

// This function returns the proportion of common elements
// with respect to the total elements of two hashPools
func (p HashPool) similarity(pp *HashPool) float32 {
	m := p.Comp(pp)
	// similarity between parties a and b
	r := float32(m) / float32(p.GetK()+pp.GetK()-m)

	if m == 0 {
		r = 0
	}
	return r
}

// Decode calculates the set difference
// based on the value of similarity
func (p HashPool) Decode(pp *HashPool) uint {
	r := p.similarity(pp)

	d := (1 - r) / (1 + r) * float32(p.hashes.Len()+pp.hashes.Len())

	return uint(d)
}

func (p *HashPool) String() string {
	return fmt.Sprintf("HashPool:\n%v\nConfig:%v\n", p.hashes, p.config)
}
