// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// bootnode runs a bootstrap node for the Ethereum Discovery Protocol.
package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/dipperin/dipperin-core/cmd/utils"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p/discover"
	"github.com/dipperin/dipperin-core/third-party/p2p/discv5"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/dipperin/dipperin-core/third-party/p2p/nat"
	"github.com/dipperin/dipperin-core/third-party/p2p/netutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"strconv"
)

var (
	listenAddr  = flag.String("addr", ":30301", "listen address")
	genKey      = flag.String("genkey", "", "generate a node key")
	writeAddr   = flag.Bool("writeaddress", false, "write out the node's public key and quit")
	nodeKeyFile = flag.String("nodekey", "", "private key filename")
	nodeKeyHex  = flag.String("nodekeyhex", "", "private key as hex (for testing)")
	natdesc     = flag.String("nat", "none", "port mapping mechanism (any|none|upnp|pmp|extip:<IP>)")
	netrestrict = flag.String("netrestrict", "", "restrict network communication to the given IP networks (CIDR masks)")
	runv5       = flag.Bool("v5", false, "run a v5 topic discovery bootnode")
)

func printBootNodeFlag() {
	log.DLogger.Info("the listenAddr is:", zap.String("listenAddr", *listenAddr))
	log.DLogger.Info("the genKey is:", zap.String("genKey", *genKey))
	log.DLogger.Info("the writeAddr is:", zap.Bool("writeAddr", *writeAddr))
	log.DLogger.Info("the nodeKeyFile is:", zap.String("nodeKeyFile", *nodeKeyFile))
	log.DLogger.Info("the nodeKeyHex is:", zap.String("nodeKeyHex", *nodeKeyHex))
	log.DLogger.Info("the natdesc is:", zap.String("natdesc", *natdesc))
	log.DLogger.Info("the netrestrict is:", zap.String("netrestrict", *netrestrict))
	log.DLogger.Info("the runv5 is:", zap.Bool("runv5", *runv5))
}

func main() {
	flag.Parse()

	var (
		nodeKey *ecdsa.PrivateKey
		err     error
	)
	cnf := log.LoggerConfig{
		Lvl:         zapcore.DebugLevel,
		FilePath:    "",
		Filename:    "",
		WithConsole: true,
		WithFile:    false,
	}
	switch chain_config.GetCurBootsEnv() {
	case chain_config.BootEnvVenus, chain_config.BootEnvMercury:
		cnf.Lvl = zapcore.InfoLevel
	}
	log.InitLogger(cnf)

	natm, err := nat.Parse(*natdesc)
	if err != nil {
		utils.Fatalf("-nat: %v", err)
	}
	switch {
	case *genKey != "":
		nodeKey, _ = crypto.GenerateKey()
		if err = crypto.SaveECDSA(*genKey, nodeKey); err != nil {
			utils.Fatalf("%v", err)
		}
		return
	case *nodeKeyFile == "" && *nodeKeyHex == "":
		utils.Fatalf("Use -nodekey or -nodekeyhex to specify a private key")
	case *nodeKeyFile != "" && *nodeKeyHex != "":
		utils.Fatalf("Options -nodekey and -nodekeyhex are mutually exclusive")
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA(*nodeKeyFile); err != nil {
			utils.Fatalf("-nodekey: %v", err)
		}
	case *nodeKeyHex != "":
		if nodeKey, err = crypto.HexToECDSA(*nodeKeyHex); err != nil {
			utils.Fatalf("-nodekeyhex: %v", err)
		}
	}

	if *writeAddr {
		fmt.Printf("%x\n", crypto.FromECDSAPub(&nodeKey.PublicKey)[1:])
		return
	}

	var restrictList *netutil.Netlist
	switch chain_config.GetCurBootsEnv() {
	case chain_config.BootEnvMercury:
		if *netrestrict != "" {
			restrictList, err = netutil.ParseNetlist(*netrestrict)
			if err != nil {
				utils.Fatalf("-netrestrict: %v", err)
			}
		}
		/*	case chain_config.BootEnvVenus:
			if *netrestrict != "" {
				restrictList, err = netutil.ParseNetlist(*netrestrict)
				if err != nil {
					utils.Fatalf("-netrestrict: %v", err)
				}
			}*/
	case "test":
		restrictList, _ = netutil.ParseNetlist(chain_config.TestIPWhiteList)
	}

	addr, err := net.ResolveUDPAddr("udp", *listenAddr)
	if err != nil {
		utils.Fatalf("-ResolveUDPAddr: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		utils.Fatalf("-ListenUDP: %v", err)
	}

	realaddr := conn.LocalAddr().(*net.UDPAddr)
	if natm != nil {
		if !realaddr.IP.IsLoopback() {
			go nat.Map(natm, nil, "udp", realaddr.Port, realaddr.Port, "ethereum discovery")
		}
		// TODO: react to external IP changes over time.
		if ext, err := natm.ExternalIP(); err == nil {
			realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
		}
	}

	udpPort, _ := strconv.ParseInt((*listenAddr)[1:], 10, 64)

	n := enode.NewV4(&nodeKey.PublicKey, net.ParseIP("127.0.0.1"), int(udpPort), int(udpPort))
	fmt.Println("bootnode conn:", n.String())

	if *runv5 {
		if _, err := discv5.ListenUDP(nodeKey, conn, "", restrictList); err != nil {
			utils.Fatalf("%v", err)
		}
	} else {
		db, _ := enode.OpenDB("")
		ln := enode.NewLocalNode(db, nodeKey)
		cfg := discover.Config{
			PrivateKey:  nodeKey,
			NetRestrict: restrictList,
		}
		if _, err := discover.ListenUDP(conn, ln, cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	select {}
}
