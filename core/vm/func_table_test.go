package vm

import (
	"testing"
	"fmt"
)

func TestResolver_ResolveFunc(t *testing.T) {
	resolver := newResolver(nil, nil, nil)
	resolver.ResolveFunc("env","getState")
}

func TestSystemFuncSet(t *testing.T) {
	cfgSet := newSystemFuncSet(nil)
	for k, v := range cfgSet {
		fmt.Println("key:", k)
		for k1, v1 := range v {
			fmt.Printf("key1: %v, v1 type: %T \n", k1, v1)
		}
	}
}