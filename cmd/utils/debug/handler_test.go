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


package debug

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestHandlerT_MemStats(t *testing.T) {
	var h *HandlerT
	h.MemStats()
}

func TestHandlerT_GcStats(t *testing.T) {
	var h *HandlerT
	h.GcStats()
}

func TestHandlerT_MemProfile(t *testing.T) {
	var h *HandlerT
	tmpF := "/tmp/xxx"
	defer os.RemoveAll(tmpF)
	assert.Nil(t, h.MemProfile(tmpF))
}

func TestHandlerT_CpuProfile(t *testing.T) {
	h := &HandlerT{}
	assert.Error(t, h.StopCPUProfile())
	tmpF := "/tmp/xxx"
	defer os.RemoveAll(tmpF)
	go h.CpuProfile(tmpF, 999999)
	time.Sleep(100 * time.Millisecond)
	assert.Error(t, h.CpuProfile(tmpF, 999999))
	h.StopCPUProfile()
}

func TestHandlerT_GoTrace(t *testing.T) {
	h := &HandlerT{}
	assert.Error(t, h.StopGoTrace())

	defer os.RemoveAll("/tmp/go_trace_test")
	go h.GoTrace("/tmp/go_trace_test", 99999)
	time.Sleep(100 * time.Millisecond)
	assert.Error(t, h.GoTrace("/tmp/go_trace_test", 99999))
	assert.NoError(t, h.StopGoTrace())
}

func TestHandlerT_BlockProfile(t *testing.T) {
	h := &HandlerT{}

	defer os.RemoveAll("/tmp/block_p_test")
	defer os.RemoveAll("/tmp/block_p")
	go h.BlockProfile("/tmp/block_p_test", 99999)
	time.Sleep(100 * time.Millisecond)
	h.SetBlockProfileRate(1)
	assert.NoError(t, h.WriteBlockProfile("/tmp/block_p"))
}

func TestHandlerT_MutexProfile(t *testing.T) {
	h := &HandlerT{}

	defer os.RemoveAll("/tmp/mutex_test")
	defer os.RemoveAll("/tmp/mutex")
	defer os.RemoveAll("/tmp/mem")
	go h.MutexProfile("/tmp/mutex_test", 9999)
	time.Sleep(100 * time.Millisecond)
	assert.NoError(t, h.WriteMutexProfile("/tmp/mutex"))
	assert.NoError(t, h.WriteMemProfile("/tmp/mem"))
	h.SetMutexProfileFraction(1)

	h.Stacks()
	h.FreeOSMemory()
	h.SetGCPercent(50)
}

func Test_expandHome(t *testing.T) {
	assert.NoError(t, os.Setenv("HOME", ""))
	expandHome("~/tmp/aaa")
}
