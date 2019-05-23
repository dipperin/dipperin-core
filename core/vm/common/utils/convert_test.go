package utils

import (
	"testing"
	"fmt"
)

func TestInt64ToBytes(t *testing.T) {
/*	a := []byte("hello")
	fmt.Println(a)
	b := BytesToInt64(a)
	fmt.Println(b)
	c := Int64ToBytes(b)
	fmt.Println(c)*/

	a := int64(2323)
	b := Int64ToBytes(a)
	fmt.Println(b)
	c := BytesToInt64(b)
	fmt.Println(c)
}