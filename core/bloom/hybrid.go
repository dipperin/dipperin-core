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

// HybridEstimator combines strata estimator and min-hash estimator
// to estimate the set difference between two parties.
type HybridEstimator struct {
	strata  *StrataEstimator
	minWise *HashPool
}

// The configuration of both the Strata Estimator and the HashPool
type HybridEstimatorConfig struct {
	StrataConfig  EstimatorConfig
	MinWiseConfig HashPoolConfig
}

// Create a new hybrid estimatorConfig with the suggested configuration
func NewHybridEstimatorConfig() HybridEstimatorConfig {
	config := HybridEstimatorConfig{
		StrataConfig:  NewEstimatorConfig(7),
		MinWiseConfig: NewHashPoolConfig(4, 500),
	}

	// TODO adjust parameters here
	// Set the length of data in IBLT of Strata en bytes
	config.StrataConfig.IBLTConfig.BktConfig.DataLen = 16

	return config
}

// Create a new hybrid estimator using the hybridEstimator config
func NewHybridEstimator(config HybridEstimatorConfig) *HybridEstimator {
	return &HybridEstimator{
		strata:  NewEstimator(config.StrataConfig),
		minWise: NewHashPool(config.MinWiseConfig),
	}
}

// Return the configuration of the two estimators
func (e *HybridEstimator) Config() HybridEstimatorConfig {
	return HybridEstimatorConfig{
		MinWiseConfig: e.minWise.config,
		StrataConfig:  e.strata.config,
	}
}

// Create a new byte slice with length equal to the length of data
func (e HybridEstimator) NewData() Data {
	return NewData(e.strata.config.IBLTConfig.BktConfig.DataLen)
}

func (e HybridEstimator) NewDataHash() DataHash {
	return NewDataHash(e.strata.config.IBLTConfig.BktConfig.HashLen)
}

// Put a data in the hybridEstimator
// Whether we should deal with it by strata or by minHash
// depends on its number of layer
func (e *HybridEstimator) Encode(d Data) {
	data := e.NewData()
	data.SetBytes(d.Bytes())
	h := hash(data)
	i := e.strata.TrailingZeros(h)

	// conditional encode data
	if i < e.strata.config.StrataNum {
		// TODO: should separate to avoid repetitive operations
		// lower layer, insert to specific layer in strata estimator
		e.strata.EncodeData(data)
	} else {
		// high layer, insert to min hash pool
		e.minWise.Encode(data)
	}
}

// Equivalent to Encode
func (e *HybridEstimator) EncodeByte(b []byte) {
	d := e.NewData()
	d.SetBytes(b)
	e.Encode(d)
}

// According to the estimator of the sender, the receiver determines
// the configuration of the IBLT used for later transmission
func (e HybridEstimator) DeriveConfig(r *HybridEstimator) InvBloomConfig {
	diff := e.Decode(r)
	// TODO adjust coefficient here
	config := NewInvBloomConfig(uint(float32(diff)*5*2), 4)

	return config
}

// Calculate the difference of the two hybridEstimators
func (e *HybridEstimator) Decode(r *HybridEstimator) uint {
	// The number of elements in the sender
	count := e.strata.count + uint(e.minWise.hashes.Len())
	// Get the total number of elements in the sender and the receiver
	count += r.strata.count + uint(r.minWise.hashes.Len())

	// Calculate the difference between the two strata
	s := e.strata.DecodeData(r.strata)

	countFactor := float32(1)
	if s > 100 {
		countFactor = 1.45
	}
	s = uint(float32(s) * countFactor)
	d := uint(0)
	if s == 0 {
		r := e.minWise.similarity(r.minWise)
		d = uint((1 - r) / (1 + r) * float32(count))
	} else {
		d = e.minWise.Decode(r.minWise)
	}

	// To sum the two differences to get the total difference
	res := s + d
	//log.DLogger.Info("strata estimator estimates difference", "difference", s)
	//log.DLogger.Info("minWise estimator estimates difference", "difference", d)

	// Set the sup limit and sub limit of res
	if count < res {
		res = count
	}

	if res < 20 {
		res = 20
	}

	//log.DLogger.Info("set difference estimated", "difference", res)
	return res
}

func (e *HybridEstimator) String() string {
	return fmt.Sprintf("Hybrid Estimator:\nStrata:\n%vMinhash:\n%vStrata Config:%v\nMinwise Config:%v\n",
		e.strata, e.minWise, e.strata.config, e.minWise.config)
}
