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

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRand2(t *testing.T) {
	nr := NewRand()
	assert.NotNil(t, nr)
}

func TestSeed2(t *testing.T) {
	Seed(1)
}

func TestRandStr2(t *testing.T) {
	RandStr(1)
}

func TestAllFunctions(t *testing.T) {
	RandUint16()
	RandUint32()
	RandUint16()
	RandUint()
	RandInt()
	RandInt16()
	RandInt32()
	RandInt64()
	RandInt31()
	RandInt31n(2)
	RandInt63()
	RandInt63n(2)
	RandBool()
	RandFloat32()
	RandFloat64()
	RandTime()
	RandBytes(2)
	RandIntn(2)
	RandPerm(2)

}

/*

func TestNewRand(t *testing.T) {
	tests := []struct {
		name string
		want *Rand
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_init(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.init()
		})
	}
}

func TestRand_reset(t *testing.T) {
	type args struct {
		seed int64
	}
	tests := []struct {
		name string
		r    *Rand
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.reset(tt.args.seed)
		})
	}
}

func TestSeed(t *testing.T) {
	type args struct {
		seed int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Seed(tt.args.seed)
		})
	}
}

func TestRandStr(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandStr(tt.args.length); got != tt.want {
				t.Errorf("RandStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandUint16(t *testing.T) {
	tests := []struct {
		name string
		want uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandUint16(); got != tt.want {
				t.Errorf("RandUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandUint32(t *testing.T) {
	tests := []struct {
		name string
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandUint32(); got != tt.want {
				t.Errorf("RandUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandUint64(t *testing.T) {
	tests := []struct {
		name string
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandUint64(); got != tt.want {
				t.Errorf("RandUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandUint(t *testing.T) {
	tests := []struct {
		name string
		want uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandUint(); got != tt.want {
				t.Errorf("RandUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt16(t *testing.T) {
	tests := []struct {
		name string
		want int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt16(); got != tt.want {
				t.Errorf("RandInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt32(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt32(); got != tt.want {
				t.Errorf("RandInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt64(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt64(); got != tt.want {
				t.Errorf("RandInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt(); got != tt.want {
				t.Errorf("RandInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt31(t *testing.T) {
	tests := []struct {
		name string
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt31(); got != tt.want {
				t.Errorf("RandInt31() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt31n(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		args args
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt31n(tt.args.n); got != tt.want {
				t.Errorf("RandInt31n() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt63(t *testing.T) {
	tests := []struct {
		name string
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt63(); got != tt.want {
				t.Errorf("RandInt63() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandInt63n(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandInt63n(tt.args.n); got != tt.want {
				t.Errorf("RandInt63n() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandBool(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandBool(); got != tt.want {
				t.Errorf("RandBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandFloat32(t *testing.T) {
	tests := []struct {
		name string
		want float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandFloat32(); got != tt.want {
				t.Errorf("RandFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandFloat64(t *testing.T) {
	tests := []struct {
		name string
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandFloat64(); got != tt.want {
				t.Errorf("RandFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandTime(t *testing.T) {
	tests := []struct {
		name string
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RandTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandBytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandBytes(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RandBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandIntn(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandIntn(tt.args.n); got != tt.want {
				t.Errorf("RandIntn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandPerm(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandPerm(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RandPerm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Seed(t *testing.T) {
	type args struct {
		seed int64
	}
	tests := []struct {
		name string
		r    *Rand
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Seed(tt.args.seed)
		})
	}
}

func TestRand_Str(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		r    *Rand
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Str(tt.args.length); got != tt.want {
				t.Errorf("Rand.Str() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Uint16(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Uint16(); got != tt.want {
				t.Errorf("Rand.Uint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Uint32(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Uint32(); got != tt.want {
				t.Errorf("Rand.Uint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Uint64(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Uint64(); got != tt.want {
				t.Errorf("Rand.Uint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Uint(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want uint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Uint(); got != tt.want {
				t.Errorf("Rand.Uint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int16(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int16(); got != tt.want {
				t.Errorf("Rand.Int16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int32(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int32(); got != tt.want {
				t.Errorf("Rand.Int32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int64(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int64(); got != tt.want {
				t.Errorf("Rand.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int(); got != tt.want {
				t.Errorf("Rand.Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int31(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int31(); got != tt.want {
				t.Errorf("Rand.Int31() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int31n(t *testing.T) {
	type args struct {
		n int32
	}
	tests := []struct {
		name string
		r    *Rand
		args args
		want int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int31n(tt.args.n); got != tt.want {
				t.Errorf("Rand.Int31n() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int63(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int63(); got != tt.want {
				t.Errorf("Rand.Int63() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Int63n(t *testing.T) {
	type args struct {
		n int64
	}
	tests := []struct {
		name string
		r    *Rand
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Int63n(tt.args.n); got != tt.want {
				t.Errorf("Rand.Int63n() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Float32(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Float32(); got != tt.want {
				t.Errorf("Rand.Float32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Float64(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Float64(); got != tt.want {
				t.Errorf("Rand.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Time(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Time(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rand.Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Bytes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		r    *Rand
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Bytes(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rand.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Intn(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		r    *Rand
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Intn(tt.args.n); got != tt.want {
				t.Errorf("Rand.Intn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Bool(t *testing.T) {
	tests := []struct {
		name string
		r    *Rand
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Bool(); got != tt.want {
				t.Errorf("Rand.Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRand_Perm(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		r    *Rand
		args args
		want []int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Perm(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rand.Perm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cRandBytes(t *testing.T) {
	type args struct {
		numBytes int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cRandBytes(tt.args.numBytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cRandBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
