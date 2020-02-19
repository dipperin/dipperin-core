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

package main

import (
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"time"
)

func main() {
	g_metrics.NewPrometheusMetricsServer(9200).Start()

	timeoutCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "time_out_count",
		Help: "mem info",
	}, []string{"count"})
	prometheus.MustRegister(timeoutCount)
	for {
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(100)
		if x%2 == 0 {
			timeoutCount.WithLabelValues("prepare").Add(1)
		} else {
			timeoutCount.WithLabelValues("prevote").Add(1)
		}
		time.Sleep(5 * time.Second)
	}
}
