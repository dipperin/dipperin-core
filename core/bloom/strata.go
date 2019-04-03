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
)

// StrataEstimator is a strata estimator consists of an array of IBLT,
// each element is inserted to an IBLT corresponding to its level of key
type StrataEstimator struct {
	//number of data
	count  uint

	// slice of IBLTs whose number is equal to the number of layers
	strata []*InvBloom
	config EstimatorConfig
}

type EstimatorConfig struct {
	// number of layers of strata
	StrataNum  uint
	// Configuration of IBLT for each layer
	IBLTConfig InvBloomConfig
}

func (config EstimatorConfig) String() string {
	return fmt.Sprintf(`{ StrataNum: %v, IBLTConfig: %v }`, config.StrataNum, config.IBLTConfig)
}

// Create a new estimator config according to the default setting of IBLT
func NewEstimatorConfig(n uint) EstimatorConfig {
	config := NewInvBloomConfig(80, 4)
	return EstimatorConfig{
		StrataNum:  n,
		IBLTConfig: config,
	}
}

// Create a new estimator according to estimatorConfig
func NewEstimator(config EstimatorConfig) *StrataEstimator {
	e := &StrataEstimator{
		strata: make([]*InvBloom, config.StrataNum),
		config: config,
	}

	// Generate a new IBLT for each layer
	for i := range e.strata {
		e.strata[i] = NewInvBloom(config.IBLTConfig)
	}

	return e
}

// Add a new data into the strata estimator
func (e *StrataEstimator) EncodeData(d Data) {
	dh := e.NewDataHash()
	h := hash(d)
	dh.SetBytes(h)

	// determines the level of the given data
	i := e.TrailingZeros(dh)

	// update the total number of elements
	e.count++
	if i >= e.config.StrataNum {
		// if there are more zeros than our layer
		// put the data into up-most layer stratum
		i = e.config.StrataNum - 1
	}

	// insert into the IBLT according to the rule of IBLT
	e.strata[i].insert(d)
}

// Create an empty byte slice with the setting of hash length
func (e StrataEstimator) NewDataHash() DataHash {
	return NewDataHash(e.config.IBLTConfig.BktConfig.HashLen)
}

// Calculate the difference of two StrataEstimator
func (e *StrataEstimator) DecodeData(r *StrataEstimator) uint {
	count := uint(0)
	for i := int(e.config.StrataNum - 1); i >= 0; i -- {
		t := NewInvBloom(e.config.IBLTConfig)
		t.Subtract(e.strata[i], r.strata[i])
		a, b := make(map[common.Hash]Data), make(map[common.Hash]Data)

		// For loop until the decode of the subtraction of certain layer fails
		// and estimate the total difference. If it succeeds until the end
		// of the for loop, then the difference is directly the sum.
		if ! t.Decode(a, b) {
			return (1 << uint(i+1)) * count
		}

		count = uint(len(a)+len(b)) + count
	}

	return count
}

// Determine the layer of IBLT the data should be inserted into
// by calculating the number of preceding zeros
func (e StrataEstimator) TrailingZeros(h DataHash) uint {
	t := e.NewDataHash()
	copy(t, h)

	res := uint(0)
	length := uint(len(t)<<3)
	for ; !t.Lsb() && res < length; res ++ {
		t.lsh()
	}

	return res
}

func (e StrataEstimator) String() string {
	//for i := range e.strata {
	//	res += fmt.Sprintf("%d%v\n", i, e.strata[i])
	//}
	return fmt.Sprintf(`{
count: %v,
strata: %v,
config: %v
}`, e.count, e.strata, e.config)
}
