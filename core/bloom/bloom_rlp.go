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
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

// The struct(s) end with -RLP are what actually transmitting
// on the networks. The conversion from original bloom to -RLP
// bloom is implemented, they are used in EncodeRLP and DecodeRLP
// functions which implicitly implements RLP encoder.
type BloomRLP struct {
	Bloom  []byte
	Config BloomConfig
}

type InvBloomRLP struct {
	Buckets []*BucketRLP
	Config  InvBloomConfig
	Salt    uint8
}

type BucketRLP struct {
	Idx     uint
	Count   uint
	KeySum  Data
	KeyHash DataHash
}

type GrapheneRLP struct {
	Bloom    *BloomRLP
	InvBloom *InvBloomRLP
}

type HybridEstimatorRLP struct {
	Strata  *StrataEstimatorRLP
	MinWise *HashPoolRLP
	Config  HybridEstimatorConfig
}

type StrataEstimatorRLP struct {
	Count  uint
	Strata []*IBloomStrataRLP
	Config EstimatorConfig
}

// This is different from InvBloomRLP, we have to drop
// the config field since it's unnecessary to carry the
// same config for different strata. They are saved under
// strata's config field.
type IBloomStrataRLP struct {
	Buckets []*BucketRLP
	Salt    uint8
}

type HashPoolRLP struct {
	Hashes SortDataHash
	Config HashPoolConfig
}

func (g Graphene) grapheneRLP() *GrapheneRLP {
	res := &GrapheneRLP{
		Bloom:    g.bloom.BloomRLP(),
		InvBloom: g.invBloom.invBloomRLP(),
	}
	return res
}

func (g *GrapheneRLP) graphene(res *Graphene) *Graphene {
	res.invBloom = NewInvBloom(g.InvBloom.Config)
	res.bloom = NewBloom(g.Bloom.Config)
	g.InvBloom.invBloom(res.invBloom)
	g.Bloom.CBloom(res.bloom)
	return res
}

func (e HybridEstimator) hybridRLP() *HybridEstimatorRLP {
	res := HybridEstimatorRLP{
		Strata:  e.strata.strataRLP(),
		MinWise: e.minWise.hashPoolRLP(),
	}
	return &res
}

func (e HybridEstimatorRLP) hybridEstimator(res *HybridEstimator) *HybridEstimator {
	res.minWise = NewHashPool(e.Config.MinWiseConfig)
	res.strata = NewEstimator(e.Config.StrataConfig)

	e.MinWise.hashPool(res.minWise)
	e.Strata.strataEstimator(res.strata)

	return res
}

func (e StrataEstimator) strataRLP() *StrataEstimatorRLP {
	res := StrataEstimatorRLP{
		Count:  e.count,
		Config: e.config,
	}

	res.Strata = make([]*IBloomStrataRLP, e.config.StrataNum)
	// deep copy each IBLT, this is non-trivial
	for i, s := range e.strata {
		res.Strata[i] = s.invBloomStrataRLP()
	}

	return &res
}

func (e StrataEstimatorRLP) strataEstimator(res *StrataEstimator) *StrataEstimator {
	res.config = e.Config
	res.count = e.Count

	res.strata = make([]*InvBloom, e.Config.StrataNum)
	for i, s := range e.Strata {
		res.strata[i] = NewInvBloom(e.Config.IBLTConfig)
		s.invBloom(res.strata[i])
	}

	return res
}

func (p HashPool) hashPoolRLP() *HashPoolRLP {
	res := HashPoolRLP{
		Config: p.config,
	}

	res.Hashes = make(SortDataHash, p.hashes.Len())
	// SortDataHash is a slice of byte slices
	// we have to deep copy each slice independently
	for i, h := range p.hashes {
		res.Hashes[i] = make(DataHash, len(h))
		copy(res.Hashes[i], h)
	}

	return &res
}

// note: must first `make` a slice of slices and then
// make each slice later
func (p *HashPoolRLP) hashPool(res *HashPool) *HashPool {
	res.hashes = make(SortDataHash, p.Hashes.Len())
	for i, h := range p.Hashes {
		res.hashes[i] = make(DataHash, len(h))
		copy(res.hashes[i], h)
	}
	res.config = p.Config

	return res
}

func (b Bloom) BloomRLP() *BloomRLP {
	res := BloomRLP{
		Bloom:  make([]byte, len(b.bloom)),
		Config: b.config,
	}
	copy(res.Bloom, b.bloom)

	return &res
}

func (b BloomRLP) CBloom(res *Bloom) *Bloom {
	res.bloom = make([]byte, len(b.Bloom))
	copy(res.bloom, b.Bloom)
	res.config = b.Config
	return res
}

func (b *InvBloom) invBloomRLP() *InvBloomRLP {
	res := InvBloomRLP{
		Config: b.config,
		Salt:   b.salt,
	}

	for i, bkt := range b.buckets {
		if bkt != nil {
			res.Buckets = append(res.Buckets, bkt.bucketRLP(uint(i)))
		}
	}

	return &res
}

func (b *InvBloom) invBloomStrataRLP() *IBloomStrataRLP {
	res := IBloomStrataRLP{
		Salt: b.salt,
	}

	for i, bkt := range b.buckets {
		if bkt != nil {
			res.Buckets = append(res.Buckets, bkt.bucketRLP(uint(i)))
		}
	}

	return &res
}

// observes the argument and returns after modifications
// we have to do it in this way because the conversion from -RLP
// struct to original struct needs to take in the function receiver
// returning a new data struct to replace the function receiver
// seems not applicable and impossible.
// Receiver could only be altered.
func (b *InvBloomRLP) invBloom(res *InvBloom) *InvBloom {
	res.config = b.Config
	res.salt = b.Salt
	res.buckets = make([]*Bucket, b.Config.BucketNum)

	for _, bkt := range b.Buckets {
		i := bkt.Idx
		if bkt != nil {
			res.buckets[i] = bkt.bucket()
		}
	}

	return res
}

func (b *IBloomStrataRLP) invBloom(res *InvBloom) *InvBloom {
	res.salt = b.Salt
	res.buckets = make([]*Bucket, res.config.BucketNum)

	for _, bkt := range b.Buckets {
		i := bkt.Idx
		if bkt != nil {
			res.buckets[i] = bkt.bucket()
		}
	}

	return res
}

// since Bucket is an intermediate data structure, and we didn't
// implement its corresponding -RLP struct.
func (b Bucket) bucketRLP(idx uint) *BucketRLP {
	return &BucketRLP{
		Idx:     idx,
		Count:   uint(b.Count),
		KeySum:  b.DataSum,
		KeyHash: b.DataHash,
	}
}

func (b BucketRLP) bucket() *Bucket {
	return &Bucket{
		Count:    int32(b.Count),
		DataSum:  b.KeySum,
		DataHash: b.KeyHash,
	}
}

func (e *HybridEstimator) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, e.hybridRLP())
}

func (e *HybridEstimator) DecodeRLP(s *rlp.Stream) error {
	var estimator HybridEstimatorRLP
	err := s.Decode(&estimator)
	estimator.hybridEstimator(e)
	return err
}

func (g *Graphene) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, g.grapheneRLP())
}

func (g *Graphene) DecodeRLP(s *rlp.Stream) error {
	var graphene GrapheneRLP
	err := s.Decode(&graphene)
	graphene.graphene(g)
	return err
}

func (e *StrataEstimator) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, e.strataRLP())
}

func (e *StrataEstimator) DecodeRLP(s *rlp.Stream) error {
	var estimator StrataEstimatorRLP
	err := s.Decode(&estimator)
	estimator.strataEstimator(e)
	return err
}

func (p *HashPool) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, p.hashPoolRLP())
}

func (p *HashPool) DecodeRLP(s *rlp.Stream) error {
	var pool HashPoolRLP
	err := s.Decode(&pool)
	pool.hashPool(p)
	return err
}

func (b *Bloom) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, b.BloomRLP())
}

func (b *Bloom) DecodeRLP(s *rlp.Stream) error {
	var bloom BloomRLP
	err := s.Decode(&bloom)
	bloom.CBloom(b)
	return err
}

//EncodeRLP implements rlp.Encoder
func (b *InvBloom) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, b.invBloomRLP())
}

//DecodeRLP implements rlp.Decoder
func (b *InvBloom) DecodeRLP(s *rlp.Stream) error {
	var bloom InvBloomRLP
	err := s.Decode(&bloom)
	bloom.invBloom(b)
	return err
}
