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

package vm_log_search

import (
	"bytes"
	"github.com/dipperin/dipperin-core/core/model"
	"math/rand"
	"testing"
)

// Tests that batched bloom bits are correctly rotated from the input bloom
// filters.
func TestGenerator(t *testing.T) {
	// Generate the input and the rotated output
	var input, output [model.BloomBitLength][model.BloomByteLength]byte

	for i := 0; i < model.BloomBitLength; i++ {
		for j := 0; j < model.BloomBitLength; j++ {
			bit := byte(rand.Int() % 2)

			input[i][j/8] |= bit << byte(7-j%8)
			output[model.BloomBitLength-1-j][i/8] |= bit << byte(7-i%8)
		}
	}
	// Crunch the input through the generator and verify the result
	gen, err := NewGenerator(model.BloomBitLength)
	//fmt.Println("fmt println   ", gen)
	if err != nil {
		t.Fatalf("failed to create bloombit generator: %v", err)
	}
	for i, bloom := range input {
		if err := gen.AddBloom(uint(i), bloom); err != nil {
			t.Fatalf("bloom %d: failed to add: %v", i, err)
		}
	}
	for i, want := range output {
		have, err := gen.Bitset(uint(i))
		if err != nil {
			t.Fatalf("output %d: failed to retrieve bits: %v", i, err)
		}
		if !bytes.Equal(have, want[:]) {
			t.Errorf("output %d: bit vector mismatch have %x, want %x", i, have, want)
		}
	}
}
