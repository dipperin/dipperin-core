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

package util

import (
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestGetCurPcIp(t *testing.T) {
	ip := GetCurPcIp("10.200.0")
	fmt.Println(ip)
}

func TestStringifyJsonToBytes(t *testing.T) {
	m := map[string]int{"123": 1}
	var rm map[string]int
	ParseJsonFromBytes(StringifyJsonToBytes(m), &rm)
	assert.Equal(t, m["123"], rm["123"])

	hb, err := hexutil.Decode("0x307830303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030")
	assert.NoError(t, err)
	fmt.Println(hb)
	// 0 is 48
	fmt.Println([]byte("0x0000000000000000000000000000000000000000000000"))
	fmt.Println(string(hb))
	fmt.Println(hexutil.Encode(hb))
}

func TestTxFeeCmp(t *testing.T) {
	//maxAmount, success := big.NewInt(0).SetString("100000000000000000000", 10)
	//assert.True(t, success)
	sumAmount := big.NewInt(3500000000)

	fmt.Println(sumAmount.Cmp(big.NewInt(consts.MinAmount)), sumAmount.Cmp(consts.MaxAmount))
	if sumAmount.Cmp(big.NewInt(consts.MinAmount)) == -1 || sumAmount.Cmp(consts.MaxAmount) == 1 {
		fmt.Println("false")
	} else {
		fmt.Println(true)
	}
}

type testRunner struct {
	stopChan chan struct{}
}

func TestStopChanClosed(t *testing.T) {
	runner := &testRunner{}
	assert.True(t, StopChanClosed(runner.stopChan))
	runner.stopChan = make(chan struct{})
	assert.False(t, StopChanClosed(runner.stopChan))
	close(runner.stopChan)
	assert.True(t, StopChanClosed(runner.stopChan))
}

func TestSetTimeout(t *testing.T) {
	stopTimerFunc := SetTimeout(func() {
		fmt.Println("timeout")
	}, time.Second)
	time.Sleep(500 * time.Millisecond)
	stopTimerFunc()

	time.Sleep(time.Second)
}

func TestExecuteFuncWithTimeout(t *testing.T) {
	ExecuteFuncWithTimeout(func() {
		time.Sleep(500 * time.Millisecond)
	}, time.Second)
}

type absHand interface {
	Hit()
}

type hand1 struct {
	X uint
}

func (h *hand1) Hit() {
	fmt.Println(h.X)
}

// 1000000	      1163 ns/op
func BenchmarkInterfaceSliceCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		from := []absHand{&hand1{X: 1}}
		to := make([]*hand1, len(from))
		InterfaceSliceCopy(to, from)
		assert.Equal(b, uint(1), to[0].X)
	}
}

// 1000000	      1004 ns/op
func BenchmarkNormalSliceCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		from := []absHand{&hand1{X: 1}}
		fLen := len(from)
		to := make([]*hand1, fLen)
		for j, f := range from {
			to[j] = f.(*hand1)
		}
		assert.Equal(b, uint(1), to[0].X)
	}
}

func TestGetUniqueId(t *testing.T) {
	id := GetUniqueId()

	assert.Len(t, id, 25)
}

func TestParseJson(t *testing.T) {
	type test struct {
		Test string `json:"test"`
	}

	var r test

	err := ParseJson(`{"test": "1"}`, &r)

	assert.NoError(t, err)

	assert.Equal(t, r.Test, "1")
}

func TestStringifyJson(t *testing.T) {
	s := struct {
		Test string `json:"test"`
	}{
		Test: "1",
	}

	str := StringifyJson(s)

	assert.Equal(t, str, `{"test":"1"}`)

}

func TestParseJsonFromBytes(t *testing.T) {
	type test struct {
		Test string `json:"test"`
	}

	var r test

	err := ParseJsonFromBytes([]byte(`{"test": "1"}`), &r)

	assert.NoError(t, err)

	assert.Equal(t, r.Test, "1")
}

func TestStringifyJsonToBytesWithErr(t *testing.T) {
	s := struct {
		Test string `json:"test"`
	}{
		Test: "1",
	}

	str, err := StringifyJsonToBytesWithErr(s)

	assert.NoError(t, err)

	assert.Equal(t, str, []byte(`{"test":"1"}`))
}

func TestFileExist(t *testing.T) {
	tf, tfClean := testTempJSONFile(t)
	defer tfClean()

	assert.True(t, FileExist(tf))
	assert.False(t, FileExist("test.log"))
}

func TestPathExists(t *testing.T) {
	assert.True(t, PathExists("/"))
	assert.False(t, PathExists("/test"))
}

func TestInterfaceIsNil(t *testing.T) {
	assert.True(t, InterfaceIsNil(nil))
	x := struct{ Test string }{Test: "1"}
	assert.False(t, InterfaceIsNil(&x))
}

func TestInterfaceSliceCopy(t *testing.T) {
	type Test struct {
		Test string
	}

	x := []Test{
		{
			Test: "3",
		},
	}

	y := make([]Test, 1)

	InterfaceSliceCopy(y, x)

	assert.Equal(t, x[0].Test, y[0].Test)
}

func TestHomeDir(t *testing.T) {
	x := HomeDir()
	assert.NotEmpty(t, x)
	assert.NoError(t, os.Setenv("HOME", ""))
	assert.NotEmpty(t, HomeDir())

	assert.True(t, IsTestEnv())
}

func TestIsTestEnv(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTestEnv(); got != tt.want {
				t.Errorf("IsTestEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
