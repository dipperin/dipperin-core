// Copyright 2014 The go-ethereum Authors
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

import (
	"testing"
)


func TestStorageSize_String(t *testing.T) {
	tests := []struct {
		name string
		s    StorageSize
		want string
	}{
		{name: "MB", s: 2381273, want: "2.38 mB"},
		{name: "kB", s: 2192, want: "2.19 kB"},
		{name: "B", s: 12, want: "12.00 B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("StorageSize.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorageSize_TerminalString(t *testing.T) {
	tests := []struct {
		name string
		s    StorageSize
		want string
	}{
		{name: "MB", s: 2381273, want: "2.38mB"},
		{name: "kB", s: 2192, want: "2.19kB"},
		{name: "B", s: 12, want: "12.00B"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.TerminalString(); got != tt.want {
				t.Errorf("StorageSize.TerminalString() = %v, want %v", got, tt.want)
			}
		})
	}
}
