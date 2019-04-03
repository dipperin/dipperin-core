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


package rpc_interface

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"net"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	apis []rpc.API

	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	ipcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	ipcListener net.Listener // IPC RPC listener socket to serve API requests
	ipcHandler  *rpc.Server  // IPC RPC request handler to process the API requests

	httpEndpoint  string       // HTTP endpoint (interface + port) to listen at (empty = HTTP disabled)
	//httpWhitelist []string     // HTTP RPC modules to allow through this endpoint
	httpListener  net.Listener // HTTP RPC listener socket to server API requests
	httpHandler   *rpc.Server  // HTTP RPC request handler to process the API requests

	wsEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)
	wsListener net.Listener // Websocket RPC listener socket to server API requests
	wsHandler  *rpc.Server  // Websocket RPC request handler to process the API requests

	allowHosts []string
}

// add extra api
func (service *Service) AddApis(apis []rpc.API) {
	service.apis = append(service.apis, apis...)
}

func (service *Service) Start() error {
	log.Info("start rpc service")
	if err := service.startInProc(service.apis); err != nil {
		return err
	}
	if err := service.startHTTP(service.httpEndpoint, service.apis, []string{}, service.allowHosts, service.allowHosts); err != nil {
		service.Stop()
		return err
	}
	log.Info("start websocket", "allow hosts", service.allowHosts)
	if err := service.startWS(service.wsEndpoint, service.apis, []string{}, service.allowHosts, true); err != nil {
		service.Stop()
		return err
	}
	return nil
}

func (service *Service) Stop() {
	service.stopInProc()
	service.stopHTTP()
	service.stopWS()
}

// startInProc initializes an in-process RPC endpoint.
func (service *Service) startInProc(apis []rpc.API) error {
	log.Debug("startInProc")
	// AddPeerSet all the APIs exposed by the services
	handler := rpc.NewServer()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			return err
		}
	}
	service.inprocHandler = handler
	return nil
}

// stopInProc terminates the in-process RPC endpoint.
func (service *Service) stopInProc() {
	if service.inprocHandler != nil {
		service.inprocHandler.Stop()
		service.inprocHandler = nil
	}
}

func (service *Service) startHTTP(endpoint string, apis []rpc.API, modules []string, cors []string, vhosts []string) error {
	log.Debug("start rpc http", "endpoint", endpoint)
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts, rpc.HTTPTimeouts{
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout: 50 * time.Second,
	})
	if err != nil {
		return err
	}
	log.Info("HTTP endpoint opened", "url", fmt.Sprintf("http://%s", endpoint), "cors", strings.Join(cors, ","), "vhosts", strings.Join(vhosts, ","))
	service.httpListener = listener
	service.httpHandler = handler
	return nil
}

func (service *Service) stopHTTP() {
	if service.httpListener != nil {
		service.httpListener.Close()
		service.httpListener = nil
	}
	if service.httpHandler != nil {
		service.httpHandler.Stop()
		service.httpHandler = nil
	}
}

// startWS initializes and starts the websocket RPC endpoint.
func (service *Service) startWS(endpoint string, apis []rpc.API, modules []string, wsOrigins []string, exposeAll bool) error {
	log.Debug("startWS", "endpoint", endpoint)
	// Short circuit if the WS endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartWSEndpoint(endpoint, apis, modules, wsOrigins, exposeAll)
	if err != nil {
		return err
	}
	log.Info("WebSocket endpoint opened", "url", fmt.Sprintf("ws://%v", endpoint))
	// All listeners booted successfully
	service.wsListener = listener
	service.wsHandler = handler
	return nil
}

// stopWS terminates the websocket RPC endpoint.
func (service *Service) stopWS() {
	if service.wsListener != nil {
		service.wsListener.Close()
		service.wsListener = nil
	}
	if service.wsHandler != nil {
		service.wsHandler.Stop()
		service.wsHandler = nil
	}
}