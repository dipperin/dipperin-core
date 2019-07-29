package model

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateBloom(t *testing.T) {
	topics := []common.Hash{common.HexToHash("topic")}
	l := &Log{Topics: topics}
	receipt1 := NewReceipt([]byte{}, false, uint64(100), []*Log{l})
	bloom := CreateBloom(Receipts{receipt1})
	assert.Equal(t, true, BloomLookup(bloom, topics[0]))
	assert.Equal(t, false, BloomLookup(bloom, common.HexToHash("bloom")))
}
