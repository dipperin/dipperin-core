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
	"github.com/kardianos/osext"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// IsIPhoneOS whether it is an apple mobile device
func IsIPhoneOS() bool {
	if runtime.GOOS == "darwin" && (runtime.GOARCH == "arm" || runtime.GOARCH == "arm64") {
		_, err := os.Stat("Info.plist")
		return err == nil
	}
	return false
}

// ChWorkDir switch back to working directory
func ChWorkDir() {
	if !IsIPhoneOS() {
		return
	}

	dir, err := filepath.Abs("")
	if err != nil {
		return
	}

	subPath := filepath.Dir(os.Args[0])
	os.Chdir(strings.TrimSuffix(dir, subPath))
}

// Executable get the real directory or real relative path where the program is located
func Executable() string {
	executablePath, err := osext.Executable()
	if err != nil {
		//cslog.Debug().Err(err).Msg("osext.Executable")
		executablePath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			//cslog.Debug().Err(err).Msg("filepath.Abs")
			executablePath = filepath.Dir(os.Args[0])
		}
	}

	if IsIPhoneOS() {
		executablePath = filepath.Join(strings.TrimSuffix(executablePath, os.Args[0]), filepath.Base(os.Args[0]))
	}

	// read link
	linkedExecutablePath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		//cslog.Debug().Err(err).Msg("filepath.EvalSymlinks")
		return executablePath
	}
	return linkedExecutablePath
}

// ExecutablePath get the directory where the program is located
func ExecutablePath() string {
	return filepath.Dir(Executable())
}

// ExecutablePathJoin returns a subdirectory of the directory where the program is located
func ExecutablePathJoin(subPath string) string {
	return filepath.Join(ExecutablePath(), subPath)
}

// WalkDir get all the files in the specified directory and all subdirectories, and match the suffix filtering.
func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix) // Ignore the case of the suffix matching

	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //  traversing the directory
		if err != nil {
			return err
		}
		if fi.IsDir() { // ignore the directory
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

// ConvertToUnixPathSeparator convert windows directory separator to unix
func ConvertToUnixPathSeparator(p string) string {
	return strings.Replace(p, "\\", "/", -1)
}
