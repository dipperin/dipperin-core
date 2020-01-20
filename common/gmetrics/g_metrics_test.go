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

package gmetrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPrometheusMetricsServer(t *testing.T) {
	//port must be 0 for this error test
	svr := NewPrometheusMetricsServer(0)
	assert.NotNil(t, svr)

	assert.Nil(t, svr.Start())

	svr = NewPrometheusMetricsServer(35231)
	svr.Start()
	time.Sleep(100 * time.Millisecond)
	svr.Start()
	svr.Stop()
}

func TestCreateCounter(t *testing.T) {
	//test enable switch close
	CreateCounter("fail", "fail", nil)

	EnableMeter()
	//create counter
	CreateCounter("count", "count", nil)

	//create counter vector
	CreateCounter("countVector", "countVector", []string{"countvVectorS"})
}

func TestCreateGauge(t *testing.T) {
	EnableMeter()
	//create gauge
	CreateGauge("gauge", "gauge", nil)

	//create gauge vector
	CreateGauge("gaugeVector", "gaugeVector", []string{"gaugeVectorS"})
}

type testMeter struct {
	name  string
	help  string
	label []string
}

var meters []testMeter

func init() {
	meters = makeMeters()
}

func makeMeters() []testMeter {
	meters := []testMeter{
		{"cnt", "cnt", nil},
		{"cntV", "cntV", []string{"cv"}},
		{"gg", "gg", nil},
		{"ggV", "ggV", []string{"gv"}},
	}
	EnableMeter()
	CreateCounter(meters[0].name, meters[0].help, meters[0].label)
	CreateCounter(meters[1].name, meters[1].help, meters[1].label)
	CreateGauge(meters[2].name, meters[2].help, meters[2].label)
	CreateGauge(meters[3].name, meters[3].help, meters[3].label)

	return meters
}

func TestAdd(t *testing.T) {
	//test non-register
	Add("notExist", "", 1)

	for _, m := range meters {
		if len(m.label) > 0 {
			//for vector, test no label
			Add(m.name, "", 1)
			//with lable
			Add(m.name, m.label[0], 1)
		} else {
			Add(m.name, "", 1)
		}
	}
}

func TestSet(t *testing.T) {
	//test non-register
	Set("notExist", "", 1)

	for _, m := range meters {
		if len(m.label) > 0 {
			//for vector, test no label
			Set(m.name, "", 1)
			//with lable
			Set(m.name, m.label[0], 1)
		} else {
			Set(m.name, "", 1)
		}
	}
}

func TestSub(t *testing.T) {
	//test non-register
	Sub("notExist", "", 1)

	for _, m := range meters {
		if len(m.label) > 0 {
			//for vector, test no label
			Sub(m.name, "", 1)
			//with lable
			Sub(m.name, m.label[0], 1)
		} else {
			Sub(m.name, "", 1)
		}
	}
}

func TestDisable(t *testing.T) {
	enable = false
	CreateCounter("", "", nil)
	CreateGauge("", "", nil)
	Set("", "", 1)
	Add("", "", 1)
	Sub("", "", 1)
}
