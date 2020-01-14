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

package rpcinterface

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

type Service struct {
	apis []rpc.API

	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	ipcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	ipcListener net.Listener // IPC RPC listener socket to serve API requests
	ipcHandler  *rpc.Server  // IPC RPC request handler to process the API requests

	httpEndpoint string // HTTP endpoint (interface + port) to listen at (empty = HTTP disabled)
	//httpWhitelist []string     // HTTP RPC modules to allow through this endpoint
	httpListener net.Listener // HTTP RPC listener socket to server API requests
	httpHandler  *rpc.Server  // HTTP RPC request handler to process the API requests

	wsEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)
	wsListener net.Listener // Websocket RPC listener socket to server API requests
	wsHandler  *rpc.Server  // Websocket RPC request handler to process the API requests

	allowHosts []string
}

func (service *Service) GetInProcHandler() *rpc.Server {
	return service.inprocHandler
}

// add extra api
func (service *Service) AddApis(apis []rpc.API) {
	service.apis = append(service.apis, apis...)
}

func (service *Service) Start() error {
	log.DLogger.Info("start rpc service")
	if err := service.startInProc(service.apis); err != nil {
		return err
	}
	//start ipc
	log.DLogger.Info("start ipc service")
	if err := service.startIPC(service.apis); err != nil {
		service.stopInProc()
		return err
	}
	log.DLogger.Info("start http service")
	if err := service.startHTTP(service.httpEndpoint, service.apis, []string{}, service.allowHosts, service.allowHosts); err != nil {
		service.stopInProc()
		service.stopIPC()
		return err
	}
	log.DLogger.Info("start websocket", zap.Strings("allow hosts", service.allowHosts))
	if err := service.startWS(service.wsEndpoint, service.apis, []string{}, service.allowHosts, false); err != nil {
		service.stopInProc()
		service.stopIPC()
		service.stopHTTP()
		return err
	}

	return nil
}

func (service *Service) Stop() {
	service.stopInProc()
	service.stopIPC()
	service.stopHTTP()
	service.stopWS()
}

// startInProc initializes an in-process RPC endpoint.
func (service *Service) startInProc(apis []rpc.API) error {
	log.DLogger.Debug("startInProc")
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
	log.DLogger.Debug("start rpc http", zap.String("endpoint", endpoint))
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartHTTPEndpoint(endpoint, apis, modules, cors, vhosts, rpc.HTTPTimeouts{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  50 * time.Second,
	})
	if err != nil {
		return err
	}
	log.DLogger.Info("HTTP endpoint opened", zap.String("url", fmt.Sprintf("http://%s", endpoint)), zap.String("cors", strings.Join(cors, ",")), zap.String("vhosts", strings.Join(vhosts, ",")))
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
	log.DLogger.Debug("startWS", zap.String("endpoint", endpoint))
	// Short circuit if the WS endpoint isn't being exposed
	if endpoint == "" {
		return nil
	}
	listener, handler, err := rpc.StartWSEndpoint(endpoint, apis, modules, wsOrigins, exposeAll)
	if err != nil {
		return err
	}
	log.DLogger.Info("WebSocket endpoint opened", zap.String("url", fmt.Sprintf("ws://%v", endpoint)))
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

// startIPC initializes and starts the IPC RPC endpoint.
func (service *Service) startIPC(apis []rpc.API) error {
	if service.ipcEndpoint == "" {
		return nil // IPC disabled.
	}
	listener, handler, err := rpc.StartIPCEndpoint(service.ipcEndpoint, apis)
	if err != nil {
		return err
	}
	service.ipcListener = listener
	service.ipcHandler = handler
	log.DLogger.Info("IPC endpoint opened", zap.String("url", service.ipcEndpoint))
	return nil
}

// stopIPC terminates the IPC RPC endpoint.
func (service *Service) stopIPC() {
	if service.ipcListener != nil {
		service.ipcListener.Close()
		service.ipcListener = nil

		log.DLogger.Info("IPC endpoint closed", zap.String("url", service.ipcEndpoint))
	}
	if service.ipcHandler != nil {
		service.ipcHandler.Stop()
		service.ipcHandler = nil
	}
}
