// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package chain_state

import (
	"errors"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
)

var (
	// errSectionOutOfBounds is returned if the user tried to add more bloom filters
	// to the batch than available space, or if tries to retrieve above the capacity.
	errSectionOutOfBounds = errors.New("section out of bounds")

	// errBloomBitOutOfBounds is returned if the user tried to retrieve specified
	// bit bloom above the capacity.
	errBloomBitOutOfBounds = errors.New("bloom bit out of bounds")
)

// Generator takes a number of bloom filters and generates the rotated bloom bits
// to be used for batched filtering.
type Generator struct {
	blooms   [model2.BloomBitLength][]byte // Rotated blooms for per-bit matching
	sections uint                          // Number of sections to batch together
	nextSec  uint                          // Next section to set when adding a bloom
}

// NewGenerator creates a rotated bloom generator that can iteratively fill a
// batched bloom filter's bits.
func NewGenerator(sections uint) (*Generator, error) {
	if sections%8 != 0 {
		return nil, errors.New("section count not multiple of 8")
	}
	b := &Generator{sections: sections}
	for i := 0; i < model2.BloomBitLength; i++ {
		b.blooms[i] = make([]byte, sections/8)
	}
	return b, nil
}

// AddBloom takes a single bloom filter and sets the corresponding bit column
// in memory accordingly.
// 整个方法的作用是判断bloom中该位是否为0，如果不为0则与generator中blooms中该位在bloom中对应的序号中对应的
func (b *Generator) AddBloom(index uint, bloom model2.Bloom) error {
	// Make sure we're not adding more bloom filters than our capacity
	if b.nextSec >= b.sections {
		return errSectionOutOfBounds
	}
	if b.nextSec != index {
		return errors.New("bloom filter with unexpected index")
	}
	// Rotate the bloom and insert into our collection
	byteIndex := b.nextSec / 8
	bitMask := byte(1) << byte(7-b.nextSec%8)
	//fmt.Println("AddBloom ============", bitMask, byte(1), byte(7-b.nextSec%8))

	for i := 0; i < model2.BloomBitLength; i++ {
		bloomByteIndex := model2.BloomByteLength - 1 - i/8
		bloomBitMask := byte(1) << byte(i%8)
		//fmt.Println("Adoom =====","i", i,",bloomByteIndex",byte(bloomByteIndex),",bloomBitMask",byte(bloomBitMask),
		//	"bloom[bloomByteIndex]",byte(bloom[bloomByteIndex]),
		//	"(bloom[bloomByteIndex] & bloomBitMask)",byte(bloom[bloomByteIndex] & bloomBitMask), "(b.blooms[i][byteIndex] | bitMask)", byte(b.blooms[i][byteIndex] | bitMask), "bteIndex", byte(byteIndex), "bitMask", byte(bitMask))

		if (bloom[bloomByteIndex] & bloomBitMask) != 0 {
			b.blooms[i][byteIndex] |= bitMask
		}
	}
	b.nextSec++

	return nil
}

// Bitset returns the bit vector belonging to the given bit index after all
// blooms have been added.
func (b *Generator) Bitset(idx uint) ([]byte, error) {
	if b.nextSec != b.sections {
		return nil, errors.New("bloom not fully generated yet")
	}
	if idx >= model2.BloomBitLength {
		return nil, errBloomBitOutOfBounds
	}
	return b.blooms[idx], nil
}
