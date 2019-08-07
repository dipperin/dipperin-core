package model

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/crypto"
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
}

func bloom9t(b []byte) *big.Int {
	b = crypto.Keccak256(b)
	fmt.Println("crypto : ", common.Bytes2Hex(b))

	r := new(big.Int)
	fmt.Println("=====r", r)

	for i := 0; i < 6; i += 2 {
		t := big.NewInt(1)
		b := (uint(b[i+1]) + (uint(b[i]) << 8)) & 2047
		fmt.Println(i, "=====b", b)
		t = t.Lsh(t, b)
		fmt.Println(i, "=====t", t)
		r.Or(r, t)
		fmt.Println(i, "=====r", r)

	}

	return r
}

func TestBloom9(t *testing.T) {
	fmt.Println(bloom9t([]byte("t")))
}
