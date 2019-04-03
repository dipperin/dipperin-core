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

import "testing"

func TestMakeName(t *testing.T) {
	t.Skip()
	type args struct {
		name    string
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				name: "test",
				version: "1.0",
			},
			want: "test/v1.0/linux/go1.11.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeName(tt.args.name, tt.args.version); got != tt.want {
				t.Errorf("MakeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileExist(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "File not exists",
			args: args{
				filePath: "test",
			},
			want: false,
		},
		{
			name: "File is exists",
			args: args{
				filePath: "./path.go",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExist(tt.args.filePath); got != tt.want {
				t.Errorf("FileExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAbsolutePath(t *testing.T) {
	type args struct {
		Datadir  string
		filename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Is absolute path",
			args: args{
				Datadir: "/test1",
				filename: "/test/",
			},
			want: "/test/",
		},
		{
			name: "Not a absolute path",
			args: args{
				Datadir: "/test1",
				filename: "test",
			},
			want: "/test1/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsolutePath(tt.args.Datadir, tt.args.filename); got != tt.want {
				t.Errorf("AbsolutePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
