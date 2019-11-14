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
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/dipperin/dipperin-core/cmd/dipperin-console"
	"github.com/dipperin/dipperin-core/cmd/dipperin-prompts"
	"github.com/dipperin/dipperin-core/cmd/dipperin/config"
	service2 "github.com/dipperin/dipperin-core/cmd/dipperin/service"
	"github.com/dipperin/dipperin-core/cmd/dipperincli/commands"
	config2 "github.com/dipperin/dipperin-core/cmd/dipperincli/config"
	"github.com/dipperin/dipperin-core/cmd/dipperincli/service"
	"github.com/dipperin/dipperin-core/cmd/utils/debug"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/dipperin"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/urfave/cli"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
)

var (
	app  *cli.App
	node dipperin.Node

	// running in backend, not start commandline
	BackendFName = "backend"
)

func main() {
	log.InitLogger(log.LvlInfo)
	//cslog.InitLogger(zap.InfoLevel, "", true)
	app = newApp()
	app.Run(os.Args)
}

func newApp() (nApp *cli.App) {
	nApp = cli.NewApp()
	nApp.Name = "DipperinCli"
	nApp.Version = chain_config.Version
	nApp.Author = "dipperin"
	nApp.Copyright = "(c) 2016-2019 dipperin."
	nApp.Usage = "Dipperin commandline tool for " + runtime.GOOS + "/" + runtime.GOARCH
	nApp.Description = ``

	nApp.Action = appAction
	nApp.Flags = append(config.Flags, debug.Flags...)
	nApp.Flags = append(nApp.Flags, cli.BoolFlag{Name: BackendFName, Usage: "set cli run without console"})
	nApp.Commands = commands.CliCommands

	sort.Sort(cli.FlagsByName(nApp.Flags))
	sort.Sort(cli.CommandsByName(nApp.Commands))
	return nApp
}

type startConf struct {
	NodeName    string `json:"node_name"`
	NodeType    int    `json:"node_type"`
	DataDir     string `json:"data_dir"`
	P2PListener string `json:"p2p_listener"`
	HTTPPort    string `json:"http_port"`
	WSPort      string `json:"ws_port"`
}

func initStartFlag() *startConf {
	//TODO
	startConfPath := filepath.Join(util.HomeDir(), ".dipperin", "start_conf.json")
	//startConfPath := filepath.Join(util.HomeDir(), ".dipperin", "start_conf2.json")
	fb, err := ioutil.ReadFile(startConfPath)
	var conf startConf
	if err != nil {
		doPrompts(&conf, startConfPath)
		return &conf
	}
	if err = util.ParseJsonFromBytes(fb, &conf); err != nil || len(conf.P2PListener) == 0 {
		doPrompts(&conf, startConfPath)
	}
	log.Info("load start flags file, you can open and change it, or rm it for reset. then must restart dipperincli", "conf_path", startConfPath)
	return &conf
}

func doPrompts(conf *startConf, saveTo string) {
	// do prompts
	conf.NodeName, _ = dipperin_prompts.NodeName()
	conf.NodeType, _ = dipperin_prompts.NodeType()
	conf.DataDir, _ = dipperin_prompts.DataDir()
	conf.P2PListener, _ = dipperin_prompts.P2PListener()
	conf.HTTPPort, _ = dipperin_prompts.HTTPPort()
	conf.WSPort, _ = dipperin_prompts.WSPort()

	// write to file
	exist, _ := soft_wallet.PathExists(saveTo)
	if !exist {
		os.MkdirAll(filepath.Dir(saveTo), 0766)
	}

	ioutil.WriteFile(saveTo, util.StringifyJsonToBytes(conf), 0644)
	log.Info("write start flags file, you can open and change it, or rm it for reset. then must restart dipperincli", "conf_path", saveTo)
}

func appAction(c *cli.Context) {
	// init log
	//lvStr := c.String(config.LogLevelFlagName)
	//lv, _ := cslog.LvlFromString(lvStr)
	//cslog.InitLogger(lv, "", true)
	//commands.InitLog(lv)

	startFlagsConf := initStartFlag()
	log.Debug("set loaded conf flags")
	c.Set(config.NodeNameFlagName, startFlagsConf.NodeName)
	c.Set(config.NodeTypeFlagName, fmt.Sprintf("%v", startFlagsConf.NodeType))
	c.Set(config.DataDirFlagName, startFlagsConf.DataDir)
	c.Set(config.P2PListenerFlagName, startFlagsConf.P2PListener)
	c.Set(config.HttpPortFlagName, startFlagsConf.HTTPPort)
	c.Set(config.WsPortFlagName, startFlagsConf.WSPort)

	path, _ := dipperin_prompts.WalletPath(startFlagsConf.DataDir)
	c.Set(config.SoftWalletPath, path)

	pwd, _ := dipperin_prompts.WalletPassword()
	c.Set(config.SoftWalletPasswordFlagName, pwd)

	passPhrase, _ := dipperin_prompts.WalletPassPhrase()
	c.Set(config.SoftWalletPassPhraseFlagName, passPhrase)

	log.Debug("the c.Args number is:", "number", c.NArg())
	log.Debug("the c.Args is:", "args", c.Args())

	nodeType := getNodeType(c.Int("node_type"))
	if nodeType == "" {
		return
	}

	log.Info("node info", "name", c.String(config.NodeNameFlagName), "type", nodeType)

	if err := debug.Setup(c); err != nil {
		log.Error("debug setup failed", "err", err)
	}

	log.Info("network info", "name", os.Getenv("boots_env"))

	if os.Getenv("boots_env") == "mercury" {
		log.Error("the Mercury testnet is not stopped forever, please try set boots_env = venus.")
		// return
	}

	err := startNode(c)
	if err != nil {
		panicInfo := "start node error err: " + err.Error()
		panic(panicInfo)
	}

	commands.InitRpcClient(node.GetNodeInfo())

	commands.InitAccountInfo(c.Int("node_type"), path, pwd, passPhrase)

	// Non-command line background startup
	if c.Bool(BackendFName) {
		// Receive stop command
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		select {
		case <-c:
		}
		return
	}

	csConsole := dipperin_console.NewConsole(Executor(c), config2.DipperinCliCompleterNew)

	defer func() {
		if node != nil {
			node.Stop()
		}
		err = csConsole.History.DoWriteHistory()
		if err != nil {
			log.Error("do write history error", "err", err)
		}
		closeApp()
	}()

	commands.PrintCommandsModuleName()
	commands.PrintDefaultAccountStake()

	go commands.CheckDownloaderSyncStatus()

	//Use the format of block subscription instead of timing printing
	/*	if nodeType == "verifier" {
		commands.AsyncLogElectionTx()
	}*/

	csConsole.Prompt.Run()
}

func getNodeType(nodeType int) (nodeTypeStr string) {
	switch nodeType {
	case 0:
		if commands.CheckRegistration() {
			log.Warn("You are registered and are not allowed to open the normal node. Please switch the verifier node to send the cancellation transaction and link up before switching the normal node")
			return
		}
		nodeTypeStr = "normal"
	case 1:
		if commands.CheckRegistration() {
			log.Warn("You are registered and are not allowed to open the miner node. Please switch the verifier node to send the cancellation transaction and link up before switching the miner node")
			return
		}
		nodeTypeStr = "mine master"
	case 2:
		nodeTypeStr = "verifier"
	}
	return
}

func haveCmd(cmd string) bool {
	cmdStr := strings.Split(cmd, " ")
	for _, c := range commands.CliCommands {
		if c.Name == cmdStr[0] {
			return true
		}
	}
	return false
}

func Executor(c *cli.Context) prompt.Executor {
	return func(command string) {
		if command == "" {
			return
		} else if command == "exit" {
			closeApp()
		} else if !haveCmd(command) {
			fmt.Println("unknown command: " + command)
			return
		}

		cmdArgs := strings.Split(strings.TrimSpace(command), " ")
		if len(cmdArgs) == 0 {
			return
		} else if len(cmdArgs) == 1 && cmdArgs[0] != "-h" && cmdArgs[0] != "--help" {
			fmt.Println("Please assign the method you want to call!")
			return
		}
		s := []string{os.Args[0]}
		s = append(s, cmdArgs...)
		//fmt.Println("s final:",s)
		//s = []string{"dipperincli", "tx", "SendTransactionContract", "--abi"}
		//s = []string{"dipperincli", "tx", "SendTransactionContract"}
		if len(cmdArgs) >= 2 {
			if config2.CheckModuleMethodIsRight(cmdArgs[0], cmdArgs[1]) {
				err := c.App.Run(s)
				log.Info("Executor", "err", err)
			} else {
				fmt.Println("module", cmdArgs[0], "has not method", cmdArgs[1])
			}
		} else {
			err := c.App.Run(s)
			log.Info("Executor", "err", err)
		}
	}
}

func startNode(c *cli.Context) error {
	if node == nil {
		var err error
		if node, err = service2.StartNode(c, true, false, true); err != nil {
			return err
		}
	} else {
		return errors.New("node is running")
	}
	return nil
}

func closeApp() {
	service.Exit()
}
