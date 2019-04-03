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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConsoleHistory(t *testing.T) {
	defer os.RemoveAll("/tmp/test_console.his")
	his, err := NewConsoleHistory("/tmp/test_console.his")
	assert.NoError(t, err)
	assert.NotNil(t, his)
	his.AddHistoryItem("123")
	err = his.DoWriteHistory()
	assert.NoError(t, err)
	err = his.ReadHistory()
	assert.NoError(t, err)
	err = his.ClearHistory()
	assert.NoError(t, err)

	var his2 *ConsoleHistory
	assert.Error(t, his2.DoWriteHistory())
}

func TestConsoleHistory_AddHistoryItem(t *testing.T) {
	defer os.RemoveAll("/tmp/test_console.his")
	his, err := NewConsoleHistory("/tmp/test_console.his")

	assert.NoError(t, err)
	his.AddHistoryItem("123")
	assert.Equal(t, len(his.historyStrs), 1)

	his.AddHistoryItem("1234")
	assert.Equal(t, len(his.historyStrs), 2)
	his.AddHistoryItem("1234")
	assert.Equal(t, len(his.historyStrs), 2)
}

func TestConsoleHistory_ClearHistory(t *testing.T) {
	var his *ConsoleHistory
	assert.Error(t, his.ClearHistory())
}
