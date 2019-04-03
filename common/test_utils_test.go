// Copyright 2015 The go-ethereum Authors
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
	"io/ioutil"
	"os"
	"testing"
)

func testTempJSONFile(t *testing.T, content []byte) (string, func()) {
	t.Helper()
	tf, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatalf(":err: %s", err)
	}

	tf.Write(content)

	tf.Close()
	return tf.Name(), func() { os.Remove(tf.Name()) }
}

func Test_LoadJSON(t *testing.T) {

	type args struct {
		val  interface{}
	}

	var res []string

	tests := []struct {
		name    string
		content string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			content: "[\"test\"]",
			args: args{
				val: &res,
			},
			wantErr: false,
		},
		{
			name: "syntax error",
			content: "",
			args: args{
				val: &res,
			},
			wantErr: true,
		},
		{
			name: "unmarshal error",
			content: "{\"test\": 1}",
			args: args{
				val: &res,
			},
			wantErr: true,
		},
		{
			name: "read file error",
			content: "1",
			args: args{
				val: &res,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tf, tfClean := testTempJSONFile(t, []byte(tt.content))
			defer tfClean()
			if tt.content == "1" {
				tf = ""
			}

			if err := LoadJSON(tf, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("LoadJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findLine(t *testing.T) {
	type args struct {
		data   []byte
		offset int64
	}
	tests := []struct {
		name     string
		args     args
		wantLine int
	}{
		{
			name: "One line",
			args: args{
				data: []byte("213121"),
				offset: 0,
			},
			wantLine: 1,
		},
		{
			name: "two line",
			args: args{
				data: []byte("213121\n34234"),
				offset: 16,
			},
			wantLine: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLine := findLine(tt.args.data, tt.args.offset); gotLine != tt.wantLine {
				t.Errorf("findLine() = %v, want %v", gotLine, tt.wantLine)
			}
		})
	}
}
