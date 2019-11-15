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

package g_metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var metrics map[string]interface{}
var enable bool = false

func init() {
	metrics = make(map[string]interface{})
}

func CreateCounter(name string, help string, label []string) {
	if !enable {
		return
	}
	if label == nil {
		counter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: name,
			Help: help,
		})
		metrics[name] = counter
		prometheus.MustRegister(counter)
	} else {
		counter := prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: name,
			Help: help,
		}, label)
		metrics[name] = counter
		prometheus.MustRegister(counter)
	}
}

func CreateGauge(name string, help string, label []string) {
	if !enable {
		return
	}
	if label == nil {
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		})
		metrics[name] = gauge
		prometheus.MustRegister(gauge)
	} else {
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, label)
		metrics[name] = gauge
		prometheus.MustRegister(gauge)
	}
}

func CreateHistogram(name string, help string){
	if !enable {
		return
	}

	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: name,
		Help: help,
	})
	metrics[name] = histogram
	prometheus.MustRegister(histogram)

}

func NewTimer(name string) *prometheus.Timer{
	if !enable {
		return nil
	}
	if metrics[name] == nil {
		return nil
	}

	switch metrics[name].(type){
	case prometheus.Histogram:
		return prometheus.NewTimer(metrics[name].(prometheus.Histogram))
	}

	return nil
}

func EnableMeter() {
	enable = true
}

func Set(name string, label string, value float64) {
	if !enable {
		return
	}
	if metrics[name] == nil {
		return
	}

	switch g := metrics[name].(type) {
	case *prometheus.GaugeVec:
		if label == "" {
			return
		}
		g.WithLabelValues(label).Set(value)
	case prometheus.Gauge:
		g.Set(value)
	default:
	}
}

func Add(name string, label string, value float64) {
	if !enable {
		return
	}
	if metrics[name] == nil {
		return
	}

	switch meter := metrics[name].(type) {
	case prometheus.Gauge:
		meter.Add(value)
	case prometheus.Counter:
		meter.Add(value)
	case *prometheus.GaugeVec:
		if label == "" {
			return
		}
		meter.WithLabelValues(label).Add(value)
	case *prometheus.CounterVec:
		if label == "" {
			return
		}
		meter.WithLabelValues(label).Add(value)
	}
}

func Sub(name string, label string, value float64) {
	if !enable {
		return
	}
	if metrics[name] == nil {
		return
	}

	switch meter := metrics[name].(type) {
	case prometheus.Gauge:
		meter.Sub(value)

	case *prometheus.GaugeVec:
		if label == "" {
			return
		}
		meter.WithLabelValues(label).Sub(value)

	default:
	}
}
