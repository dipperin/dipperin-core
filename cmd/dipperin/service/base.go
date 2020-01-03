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

package service

import (
	"github.com/dipperin/dipperin-core/cmd/dipperin/config"
	"github.com/dipperin/dipperin-core/cmd/utils/debug"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/dipperin"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func StartNode(c *cli.Context, async bool, logToConsole bool, logToFile bool) (dipperin.Node, error) {
	// set log level and others
	extraBeforeStart(c, logToConsole, logToFile)
	// make a node
	nodeConf := getNodeConf(c)
	if err := nodeConf.NodeConfigCheck(); err != nil {
		return nil, err
	}

	node := dipperin.NewBftNode(nodeConf)
	debug.Memsize.Add("node", node)
	// start node
	if err := node.Start(); err != nil {
		return nil, err
	}
	// listen stop
	go signalListen(node)

	// wait stop
	if async {
		go node.Wait()
	} else {
		node.Wait()
	}

	return node, nil
}

func extraBeforeStart(c *cli.Context, logToConsole bool, logToFile bool) {
	logLevel := c.String(config.LogLevelFlagName)
	lv, err := log.LvlFromString(logLevel)
	if err != nil {
		log.DLogger.Error("parse log level", zap.Error(err))
	}

	if !logToFile && os.Getenv("cslog") == "enable" {
		logToFile = true
		logToConsole = false
	}

	dataDir := c.String(config.DataDirFlagName)
	log.DLogger.Info("init logger", zap.String("lv", lv.String()), zap.Bool("log to file", logToFile))
	log.InitLogger(log.LoggerConfig{
		Lvl:         lv,
		FilePath:    dataDir,
		Filename:    "",
		WithConsole: logToConsole,
		WithFile:    logToFile,
	})

}

func getNodeConf(c *cli.Context) dipperin.NodeConfig {
	nodeConf := dipperin.DefaultNodeConf()
	nodeConf.Name = c.String(config.NodeNameFlagName)
	nodeConf.HTTPHost = c.String(config.HttpHostFlagName)
	nodeConf.HTTPPort = c.Int(config.HttpPortFlagName)
	nodeConf.WSHost = c.String(config.WsHostFlagName)
	nodeConf.WSPort = c.Int(config.WsPortFlagName)
	nodeConf.IPCPath = c.String(config.IPCPathFlagName)
	nodeConf.DataDir = c.String(config.DataDirFlagName)
	nodeConf.NodeType = c.Int(config.NodeTypeFlagName)
	nodeConf.DebugMode = c.Int(config.DebugModeFlagName)
	nodeConf.P2PListener = c.String(config.P2PListenerFlagName)
	if nodeConf.P2PListener[0] != ':' {
		nodeConf.P2PListener = ":" + nodeConf.P2PListener
	}
	nodeConf.NoWalletStart = c.Bool(config.NoWalletStartFlagName)
	nodeConf.SoftWalletPassword = c.String(config.SoftWalletPasswordFlagName)
	nodeConf.SoftWalletPassPhrase = c.String(config.SoftWalletPassPhraseFlagName)
	nodeConf.SoftWalletPath = c.String(config.SoftWalletPath)
	nodeConf.IsScanner = c.Int(config.IsScannerFlagName)
	nodeConf.IsUploadNodeData = c.Int(config.IsUploadNodeData)
	nodeConf.UploadURL = c.String(config.UploadURL)
	nodeConf.NoDiscovery = c.Int(config.NoDiscovery)
	nodeConf.Nat = c.String(config.Nat)
	nodeConf.AllowHosts = c.StringSlice(config.AllowHostsFlagName)
	nodeConf.PMetricsPort = c.Int(config.MetricsPortFlagName)

	if c.Int(config.IsStartMine) == 0 {
		nodeConf.IsStartMine = false
	} else {
		nodeConf.IsStartMine = true
	}

	log.DLogger.Info("getNodeConf the node type is:", zap.Int("nodeType", nodeConf.NodeType))

	return nodeConf
}

func signalListen(n dipperin.Node) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	s := <-c
	log.DLogger.Info("got system signal", zap.Any("signal", s))
	debug.Exit()
	n.Stop()
}
