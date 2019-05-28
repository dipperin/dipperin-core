package vmcommon

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/vmcommon"
	"testing"
)

func TestInt64ToBytes(t *testing.T) {
/*	a := []byte("hello")
	fmt.Println(a)
	b := BytesToInt64(a)
	fmt.Println(b)
	c := Int64ToBytes(b)
	fmt.Println(c)*/

	a := int64(2323)
	b := vmcommon.Int64ToBytes(a)
	fmt.Println(b)
	c := vmcommon.BytesToInt64(b)
	fmt.Println(c)
}