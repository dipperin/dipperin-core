package model

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateBloom(t *testing.T) {
	topics := []common.Hash{common.HexToHash("topic")}
	l := &Log{Topics: topics}
	receipt1 := NewReceipt([]byte{}, false, uint64(100), []*Log{l})
	bloom := CreateBloom(Receipts{receipt1})
	assert.Equal(t, true, BloomLookup(bloom, topics[0]))
	assert.Equal(t, false, BloomLookup(bloom, common.HexToHash("bloom")))

	bloom.Add(big.NewInt(100))
	assert.Equal(t, true, bloom.Test(big.NewInt(100)))
	assert.Equal(t, true, bloom.TestBytes(big.NewInt(100).Bytes()))
	assert.Equal(t, bloom.Hex(), hexutil.Encode(bloom.Bytes()))
	assert.Panics(t, func() {
		d := make([]byte, BloomByteLength+1)
		bloom.SetBytes(d)
	})
}

func TestBloom_MarshalText(t *testing.T) {
	topics := []common.Hash{common.HexToHash("topic")}
	l := &Log{Topics: topics}
	receipt1 := NewReceipt([]byte{}, false, uint64(100), []*Log{l})
	bloom := CreateBloom(Receipts{receipt1})

	// MarshalJSON
	enc, err1 := bloom.MarshalText()
	assert.NoError(t, err1)

	// UnmarshalJSON
	b1get := &Bloom{}
	err2 := b1get.UnmarshalText(enc)
	assert.NoError(t, err2)
	assert.EqualValues(t, &bloom, b1get)
}
