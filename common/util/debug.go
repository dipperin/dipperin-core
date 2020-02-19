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
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

func WriteDebugStack(toFile string) error {
	// prepare dir
	toDir := filepath.Dir(toFile)
	if !FileExist(toDir) {
		if err := os.MkdirAll(toDir, 0755); err != nil {
			return err
		}
	}

	// write stack
	buf := make([]byte, 5*1024*1024)
	buf = buf[:runtime.Stack(buf, true)]
	return ioutil.WriteFile(toFile, buf, 0755)
}
