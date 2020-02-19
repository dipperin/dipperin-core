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

package dipperin_console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"unicode/utf8"
)

const (
	HistoryLimit = 1000
)

type ConsoleHistory struct {
	historyFilePath string
	historyFile     *os.File
	historyMutex    sync.RWMutex
	historyStrs     []string
}

// NewLineHistory setting history
func NewConsoleHistory(filePath string) (ch *ConsoleHistory, err error) {
	ch = &ConsoleHistory{
		historyFilePath: filePath,
	}

	ch.historyFile, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (ch *ConsoleHistory) AddHistoryItem(s string) {
	sLen := len(ch.historyStrs)
	if sLen != 0 {
		if s == ch.historyStrs[sLen-1] {
			return
		}
	}
	ch.historyStrs = append(ch.historyStrs, s)
}

// DoWriteHistory execute write history file
func (ch *ConsoleHistory) DoWriteHistory() (err error) {
	if ch == nil {
		return fmt.Errorf("history not set")
	}
	ch.historyMutex.RLock()
	defer func() {
		if ch != nil {
			ch.historyMutex.RUnlock()
		}
	}()

	ch.historyFile, err = os.Create(ch.historyFilePath)
	if err != nil {
		return fmt.Errorf("write to history file error, %s", err)
	}

	for _, item := range ch.historyStrs {
		_, err := fmt.Fprintln(ch.historyFile, item)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return fmt.Errorf("write to history file error: %s", err)
	}

	return nil
}

// ReadHistory reading from history file
func (ch *ConsoleHistory) ReadHistory() (err error) {
	if ch == nil {
		return fmt.Errorf("history not set")
	}

	ch.historyMutex.Lock()
	defer func() {
		if ch != nil {
			ch.historyMutex.Unlock()
		}
	}()

	in := bufio.NewReader(ch.historyFile)
	num := 0
	for {
		line, part, err := in.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if part {
			return fmt.Errorf("line %d is too long", num+1)
		}
		if !utf8.Valid(line) {
			return fmt.Errorf("invalid string at line %d", num+1)
		}
		num++
		ch.historyStrs = append(ch.historyStrs, string(line))
		if len(ch.historyStrs) > HistoryLimit {
			ch.historyStrs = ch.historyStrs[1:]
		}
	}

	return err
}

func (ch *ConsoleHistory) ClearHistory() (err error) {
	if ch == nil {
		return fmt.Errorf("history not set")
	}

	ch.historyStrs = []string{}
	ch.DoWriteHistory()

	return nil
}
