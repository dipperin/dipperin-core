package components

import (
	"testing"
	"github.com/issue9/assert"
)

func TestNewBlockPool(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	t.Log("www:", rsp)
}

func TestBlockPool_AddBlock(t *testing.T) {
	NewBlockPool(0, nil).SetNodeConfig(nil)
}

func TestBlockPool_SetPoolEventNotifier(t *testing.T) {
	NewBlockPool(0, nil).SetPoolEventNotifier(nil)

}

func TestBlockPool_Start(t *testing.T) {
	NewBlockPool(0, nil).Start()
}

func TestBlockPool_Stop(t *testing.T) {
	NewBlockPool(0, nil).Stop()
}

func TestBlockPool_IsEmpty(t *testing.T) {
	sign := NewBlockPool(0, nil).IsEmpty()
	assert.Equal(t, sign, true)
}

func TestBlockPool_IsRunning(t *testing.T) {
	sign := NewBlockPool(0, nil).IsRunning()
	assert.Equal(t, sign, false)
}
