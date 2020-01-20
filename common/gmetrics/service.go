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
	"fmt"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
)

// This method should be placed at the forefront to ensure that it can be call before other services are registered.
func NewPrometheusMetricsServer(port int) *PrometheusMetricsServer {
	pms := &PrometheusMetricsServer{
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: promhttp.Handler(),
		},
	}
	if port != 0 {
		EnableMeter()
	}
	return pms
}

type PrometheusMetricsServer struct {
	port   int
	server *http.Server
}

func (p *PrometheusMetricsServer) Start() error {
	if p.port == 0 {
		log.DLogger.Info("port is 0, do not start prometheus metrics server")
		return nil
	}

	log.DLogger.Info("start prometheus metrics", zap.String("addr", p.server.Addr))
	go func() {
		if err := p.server.ListenAndServe(); err != nil {
			log.DLogger.Error("pMetrics serve failed: " + err.Error())
		}
	}()
	return nil
}

func (p *PrometheusMetricsServer) Stop() {}
