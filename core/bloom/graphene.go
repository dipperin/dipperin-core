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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
)

// Graphene are made of a bloom and an invBloom
type Graphene struct {
	invBloom *InvBloom
	bloom    *Bloom
}

// Create a new Graphene with the configurations of bloom and invBloom.
func NewGraphene(i InvBloomConfig, b BloomConfig) *Graphene {
	return &Graphene{
		invBloom: NewInvBloom(i),
		bloom:    NewBloom(b),
	}
}

// Filter filters out the keys in bob that are impossible in Alice by the
// filtering of bloom in Graphene. Remember there is still false positive.
func (g Graphene) Filter(bob map[common.Hash]Data) {
	// this is safe deleting map element when iterating
	// through itself
	for k, v := range bob {
		if !g.bloom.LookUp(v) {
			delete(bob, k)
		}
	}
}

// Put the elements in the IBLT of Graphene.
func (g *Graphene) InsertMap(m map[common.Hash]Data) *Graphene {
	g.invBloom.InsertMap(m)
	return g
}

// Insert inserts the key-value pair into certain bucket in the bloom.
// This operation always succeeds, assuming that all keys are distinct.
func (g *Graphene) Insert(x Data) {
	g.invBloom.insert(x)
	g.bloom.Digest(x[:])
}

// call the InsertRLP method of invBloom
func (g *Graphene) InsertRLP(k, v interface{}) {
	g.invBloom.InsertRLP(k, v)

	data, err := rlp.EncodeToBytes(k)

	if err != nil {
		log.DLogger.Error("key RLP encode error")
		return
	}

	g.bloom.Digest(data)
}

// Return the Bloom of the Graphene
func (g Graphene) Bloom() *Bloom {
	return g.bloom
}

// Return the InvBloom of the Graphene
func (g Graphene) InvBloom() *InvBloom {
	return g.invBloom
}

// Check the existence of data in the bloom.
// Note that there is still false positive
func (g Graphene) LookUp(x Data) bool {
	return g.bloom.LookUp(x)
}

// FilterDecode returns the unique elements of Alice not existing in Bob.
// It modified the callee.
func (g *Graphene) FilterDecode(bob map[common.Hash]Data) map[common.Hash]Data {
	// Reserve the elements of Bob possible to be in Alice(including false positive).
	g.Filter(bob)

	alice := make(map[common.Hash]Data)

	// Generate a new InvBloom by the elements of b to subtract Alice
	b := NewInvBloom(g.invBloom.config).InsertMap(bob)

	g.invBloom.Subtract(g.invBloom, b)

	if g.invBloom.Decode(alice, nil) {
		return alice
	} else {
		return nil
	}
}

// FilterRLP filters out the keys that are impossible in the map's keys
// it will not modify the input interface.
// it panics if the input is not an valid map.
func (g Graphene) FilterRLP(i interface{}) map[interface{}]interface{} {
	m := interfaceToInterfaceMap(i)

	res := make(map[interface{}]interface{})

	for k, v := range m {
		if k == nil {
			log.DLogger.Error("insert key is nil")
			return nil
		}

		kBytes, kError := rlp.EncodeToBytes(k)

		if kError != nil {
			log.DLogger.Error("key RLP encode error", zap.Error(kError), zap.Any("key", k))
			return nil
		}

		// The key value are registered only if kBytes is possible is be in Alice.
		if g.bloom.LookUp(kBytes) {
			res[k] = v
		}
	}

	return res
}

// FilterListRLP returns the slice of RLP coded slices that are uniquely
// in caller's bloom lookup table.
func (g *Graphene) FilterListRLP(i interface{}) (Alice [][]byte, err error) {
	m := g.FilterRLP(i)
	b := NewInvBloom(g.invBloom.config)
	for k, v := range m {
		b.InsertRLP(k, v)
	}

	g.invBloom.Subtract(g.invBloom, b)

	Alice, _, err = g.invBloom.ListRLP()

	return
}

// Recover recovers the slice of RLP coded slices that are in caller's bloom
func (g *Graphene) Recover(i interface{}) (alice [][]byte, err error) {
	mayOverlap := g.FilterRLP(i)

	b := NewInvBloom(g.invBloom.config)
	for k, v := range mayOverlap {
		b.InsertRLP(k, v)
	}

	g.invBloom.Subtract(g.invBloom, b)

	// This function gets the unique elements of Alice and of Bob
	alice, bob, err := g.invBloom.ListRLP()

	if err != nil {
		return nil, err
	}

	// a hash-map
	// since decoded elements are RLP coded.
	// the only way to test existence without applying detail types is
	// to test by using RLP code.
	bobSet := make(map[common.Hash]struct{})

	for _, b := range bob {
		bobSet[common.BytesToHash(hash(b))] = struct{}{}
	}

	for _, v := range mayOverlap {
		coded, err := rlp.EncodeToBytes(v)
		if err != nil {
			return nil, err
		}

		// append these in mayOverlap but exclude bob's
		if _, exist := bobSet[common.BytesToHash(hash(coded))]; !exist {
			alice = append(alice, coded)
		}
	}

	return alice, nil
}

// Call the ListRLP method of the IBLT in the Graphene
func (g Graphene) ListRLP() (Alice, Bob [][]byte, err error) {
	return g.invBloom.ListRLP()
}

// Return the config of the IBLT in the Graphene
func (g Graphene) InvBloomConfig() InvBloomConfig {
	return g.invBloom.config
}

// Return the config of the Bloom in the Graphene
func (g Graphene) BloomConfig() BloomConfig {
	return g.bloom.config
}
