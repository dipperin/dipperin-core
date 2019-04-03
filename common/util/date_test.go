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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	date  = time.Date(2015, time.Month(12), 31, 0, 0, 0, 0, time.UTC)
	date2 = time.Date(2016, time.Month(12), 31, 0, 0, 0, 0, time.UTC)
	zero  time.Time
)

func TestParseDateRange(t *testing.T) {
	assert := assert.New(t)

	var testDates = []struct {
		dateStr string
		start   time.Time
		end     time.Time
		errNil  bool
	}{
		{"2015-12-31:2016-12-31", date, date2, true},
		{"2015-12-31:", date, zero, true},
		{":2016-12-31", zero, date2, true},
		{"2016-12-31", zero, zero, false},
		{"2016-31-12:", zero, zero, false},
		{":2016-31-12", zero, zero, false},
	}

	for _, test := range testDates {
		start, end, err := ParseDateRange(test.dateStr)
		if test.errNil {
			assert.Nil(err)
			testPtr := func(want, have time.Time) {
				assert.True(have.Equal(want))
			}
			testPtr(test.start, start)
			testPtr(test.end, end)
		} else {
			assert.NotNil(err)
		}
	}
}
