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
	"strings"
	"time"

	"errors"
)

// TimeLayout helps to parse a date string of the format YYYY-MM-DD
//   Intended to be used with the following function:
// 	 time.Parse(TimeLayout, date)
var TimeLayout = "2006-01-02" //this represents YYYY-MM-DD

// ParseDateRange parses a date range string of the format start:end
//   where the start and end date are of the format YYYY-MM-DD.
//   The parsed dates are time.Time and will return the zero time for
//   unbounded dates, ex:
//   unbounded start:	:2000-12-31
//	 unbounded end: 	2000-12-31:
func ParseDateRange(dateRange string) (startDate, endDate time.Time, err error) {
	dates := strings.Split(dateRange, ":")
	if len(dates) != 2 {
		err = errors.New("bad date range, must be in format date:date")
		return
	}
	parseDate := func(date string) (out time.Time, err error) {
		if len(date) == 0 {
			return
		}
		out, err = time.Parse(TimeLayout, date)
		return
	}
	startDate, err = parseDate(dates[0])
	if err != nil {
		return
	}
	endDate, err = parseDate(dates[1])
	if err != nil {
		return
	}
	return
}
