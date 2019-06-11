package model

import (
	"testing"
	"fmt"
)

func TestReceipt_String(t *testing.T) {
	receipt := NewReceipt([]byte{}, false, uint64(100))
	fmt.Println(receipt.String())
}