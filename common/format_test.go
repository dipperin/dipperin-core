// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package common

import "testing"

func TestPrettyDuration_String(t *testing.T) {
	tests := []struct {
		name string
		d    PrettyDuration
		want string
	}{
		{
			name: "lengt > 4",
			d:    1,
			want: "1ns",
		},
		{
			name: "ns",
			d:    123,
			want: "123ns",
		},
		{
			name: "µs",
			d:    1200,
			want: "1.2µs",
		},
		{
			name: "ms",
			d:    1234567,
			want: "1.234ms",
		},
		{
			name: "s",
			d:    1200000000,
			want: "1.2s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("PrettyDuration.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
