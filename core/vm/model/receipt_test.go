package model

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestReceipt(t *testing.T) {
	topics := []common.Hash{common.HexToHash("topic")}
	l := &Log{Topics: topics}
	receipt1 := NewReceipt([]byte{}, false, uint64(100), []*Log{l})
	receipt2 := NewReceipt([]byte{}, true, uint64(100), []*Log{l})
	receipt3 := &Receipt{Status: uint64(10)}
	assert.Equal(t, "Successful", receipt1.GetStatusStr())
	assert.Equal(t, "Failed", receipt2.GetStatusStr())
	assert.Equal(t, "UnKnown", receipt3.GetStatusStr())

	receipts := &Receipts{receipt1, receipt2, receipt3}
	rlpReceipt1, err := rlp.EncodeToBytes(receipt1)
	assert.NoError(t, err)
	assert.Equal(t, []byte(strconv.Itoa(0)), receipts.GetKey(0))
	assert.Equal(t, 3, receipts.Len())
	assert.Equal(t, rlpReceipt1, receipts.GetRlp(0))

	storageReceipts := []*ReceiptForStorage{(*ReceiptForStorage)(receipt1), (*ReceiptForStorage)(receipt2)}
	enc, err := rlp.EncodeToBytes(storageReceipts)
	assert.NoError(t, err)

	var resp []*ReceiptForStorage
	err = rlp.DecodeBytes(enc, &resp)
	assert.NoError(t, err)
	assert.Equal(t, len(receipt1.Logs), len(resp[0].Logs))
	for _, v := range resp {
		assert.NotNil(t, (*Receipt)(v).String())
	}
}
