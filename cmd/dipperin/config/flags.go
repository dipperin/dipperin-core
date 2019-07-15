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

package config

import (
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/urfave/cli"
)

// define flag names
const (
	LogLevelFlagName = "log_level"
	//LogTypeFlagName = "log_type"

	DataDirFlagName  = "data_dir"
	NodeTypeFlagName = "node_type"

	P2PListenerFlagName = "p2p_listener"
	HttpHostFlagName    = "http_host"
	HttpPortFlagName    = "http_port"
	WsHostFlagName      = "ws_host"
	WsPortFlagName      = "ws_port"
	IPCPathFlagName     = "ipc_path"
	DebugModeFlagName   = "debug_mode"

	UseStaticNodesFlagName       = "use_static_nodes"
	SoftWalletPasswordFlagName   = "soft_wallet_pwd"
	SoftWalletPassPhraseFlagName = "soft_wallet_pass_phrase"
	SoftWalletPath               = "soft_wallet_path"

	IsScannerFlagName = "is_scanner"

	IsUploadNodeData = "is_upload_node_data"

	IsPerformanceFlagName = "is_performance"

	UploadURL = "upload_url"

	NodeNameFlagName = "node_name"

	IsStartMine = "is_start_mine"
	NoDiscovery = "no_discovery"
	Nat         = "nat"

	AllowHostsFlagName = "allow_hosts"

	MetricsPortFlagName = "m_port"

	PProfBoolFlagName ="pprof"
	PProfPortFlagName = "pprofport"
)

var (
	Flags = []cli.Flag{
		MetricsPortFlag,
		//LogTypeFlag,
		LogLevelFlag,
		DataDirFlag,
		NodeTypeFlag,
		P2PListenerFlag,
		HttpHostFlag,
		HttpPortFlag,
		WsHostFlag,
		WsPortFlag,
		IPCPathFlag,
		UseStaticNodesFlag,
		NodeNameFlag,
		DebugModeFlag,
		SoftWalletPasswordFlag,
		SoftWalletPassPhraseFlag,
		SoftWalletPathFlag,
		IsScannerFlag,
		IsUploadNodeDataFlag,
		//IsPerformanceFlag,
		UploadURLFlag,

		MetricsEnabledFlag,

		IsStartMineFlag,
		NoDiscoveryFlag,
		NatFlag,
		AllowHostsFlag,
	}
)

var (
	MetricsPortFlag = cli.IntFlag{
		Name:  MetricsPortFlagName,
		Usage: "set metrics port, not start metrics server if =0",
		Value: 0,
	}
	AllowHostsFlag = cli.StringSliceFlag{
		Name:  AllowHostsFlagName,
		Usage: "set rpc client allow hosts",
		Value: &cli.StringSlice{"localhost", "127.0.0.1"},
	}
	UploadURLFlag = cli.StringFlag{
		Name:  UploadURL,
		Usage: "set uploading data url",
		Value: "http://localhost:8887/api/Dipperin_nodes",
	}
	IsUploadNodeDataFlag = cli.IntFlag{
		Name:  IsUploadNodeData,
		Usage: "set whether uploading data,0 yes,1 no",
		Value: 0,
	}
	IsScannerFlag = cli.IntFlag{
		Name:  IsScannerFlagName,
		Usage: "set whether including browser，0 no,1 yes",
		Value: 0,
	}
	//IsPerformanceFlag = cli.IntFlag{
	//	Name:  IsPerformanceFlagName,
	//	Usage: "set whether including TPS monitoring,0 no,1 yes",
	//	Value: 0,
	//}
	SoftWalletPasswordFlag = cli.StringFlag{
		Name:  SoftWalletPasswordFlagName,
		Usage: "set whether needing password of creating or openning wallet",
		Value: "",
	}
	SoftWalletPassPhraseFlag = cli.StringFlag{
		Name:  SoftWalletPassPhraseFlagName,
		Usage: "set whether needing salt string for creating wallet",
		Value: "",
	}
	SoftWalletPathFlag = cli.StringFlag{
		Name:  SoftWalletPath,
		Usage: "set whether needing path for creating or openning wallet",
		Value: "",
	}
	HttpHostFlag = cli.StringFlag{
		Name:  HttpHostFlagName,
		Usage: "set http host",
		Value: "127.0.0.1",
	}
	WsHostFlag = cli.StringFlag{
		Name:  WsHostFlagName,
		Usage: "set web socket host",
		Value: "127.0.0.1",
	}
	UseStaticNodesFlag = cli.BoolFlag{
		Name:  UseStaticNodesFlagName,
		Usage: "set whether using static node configuration",
	}
	IPCPathFlag = cli.StringFlag{
		Name:  IPCPathFlagName,
		Usage: "set ipc directory",
		Value: "/tmp/dipperin.ipc",
	}
	//LogTypeFlag = cli.IntFlag{
	//	Name: LogTypeFlagName,
	//	Usage: "Set the log type, 0 colored log, 1 colorless log. Default 0",
	//	Value: 0,
	//}
	HttpPortFlag = cli.IntFlag{
		Name:  HttpPortFlagName,
		Usage: "set http port",
		Value: 7001,
	}
	WsPortFlag = cli.IntFlag{
		Name:  WsPortFlagName,
		Usage: "set web socket port",
		Value: 7002,
	}
	P2PListenerFlag = cli.StringFlag{
		Name:  P2PListenerFlagName,
		Usage: "set p2p port",
		Value: ":22222",
	}
	NodeNameFlag = cli.StringFlag{
		Name:  NodeNameFlagName,
		Usage: "set node alias",
	}
	DataDirFlag = cli.StringFlag{
		Name:  DataDirFlagName,
		Value: chain_config.DefaultDataDir(),
		Usage: "set node data saving directory",
	}
	LogLevelFlag = cli.StringFlag{
		Name:  LogLevelFlagName,
		Value: "info",
		Usage: "set log level: debug info warn error",
	}
	NodeTypeFlag = cli.IntFlag{
		Name:  NodeTypeFlagName,
		Value: 0,
		Usage: "set node type, normal: 0 mine master:1 verifier:2",
	}
	DebugModeFlag = cli.IntFlag{
		Name:  DebugModeFlagName,
		Value: 2,
		Usage: "set debug mode，0 single node without broadcasting，1 multiple nodes with PBFT，2 multiple nodes with PBFT and election",
	}

	MetricsEnabledFlag = cli.BoolFlag{
		Name:  metrics.MetricsEnabledFlag,
		Usage: "Enable metrics collection and reporting",
	}

	IsStartMineFlag = cli.IntFlag{
		Name:  IsStartMine,
		Value: 0,
		Usage: "set whether mining，0 no，1 yes",
	}

	NoDiscoveryFlag = cli.IntFlag{
		Name:  NoDiscovery,
		Value: 0,
		Usage: "whether closing node discovery",
	}

	NatFlag = cli.StringFlag{
		Name:  Nat,
		Value: "",
		Usage: "nat mode",
	}
)
