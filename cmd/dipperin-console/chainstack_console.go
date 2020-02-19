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

package dipperin_console

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	// EnvConfigDir configuration path
	EnvConfigDir = "CS_CI_EX_GO_CONFIG_DIR"
	// ConfigName configuration file name
	ConfigName = "cs_config.json"
	// app name
	AppName = "DipperinCIEx"
)

var (
	historyFilePath = filepath.Join(GetConfigDir(), "cs_command_history.txt")
)

// GetConfigDir get the configuration path
func GetConfigDir() string {
	// from environment variables
	configDir, ok := os.LookupEnv(EnvConfigDir)
	if ok {
		if filepath.IsAbs(configDir) {
			return configDir
		}
		// if not absolute path, loop up in directory
		return util.ExecutablePathJoin(configDir)
	}

	// use the old version
	// If the old version of the configuration file exists, use the old version
	oldConfigDir := util.ExecutablePath()
	_, err := os.Stat(filepath.Join(oldConfigDir, ConfigName))
	if err == nil {
		return oldConfigDir
	}

	switch runtime.GOOS {
	case "windows":
		return getWinConfigDir(oldConfigDir)
	default:
		dataPath, ok := os.LookupEnv("HOME")
		if !ok {
			log.Warn("Environment HOME not set")
			return oldConfigDir
		}
		//configDir = filepath.Join(dataPath, ".config", AppName)
		configDir = filepath.Join(dataPath, "tmp", AppName)

		// check if it is writable
		err = os.MkdirAll(configDir, 0700)
		if err != nil {
			log.Warn("check config dir error", "err", err)
			return oldConfigDir
		}
		return configDir
	}
}

func getWinConfigDir(oldConfigDir string) string {
	dataPath, ok := os.LookupEnv("APPDATA")
	if !ok {
		log.Warn("Environment APPDATA not set")
		return oldConfigDir
	}
	return filepath.Join(dataPath, AppName)
}

type Console struct {
	Prompt  *prompt.Prompt
	History *ConsoleHistory

	paused bool
}

// NewConsole return *Console , completer Completer
func NewConsole(executor prompt.Executor, completer prompt.Completer) *Console {
	var err error
	c := &Console{}

	log.Info("historyFilePath", "historyFilePath", historyFilePath)
	if !common.FileExist(historyFilePath) {
		_ = os.MkdirAll(filepath.Dir(historyFilePath), 0644)
	}
	c.History, err = NewConsoleHistory(historyFilePath)
	if err != nil {
		panic(err)
	}
	err = c.History.ReadHistory()
	if err != nil {
		fmt.Printf("warning reading history command file error, %s\n", err)
	}

	// New prompt and load history input
	p := prompt.New(
		c.WrapExecutor(executor),
		completer,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("Dipperin"),
		prompt.OptionHistory(c.History.historyStrs),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)

	c.Prompt = p

	return c
}

// Wrap the executor function for add the input to history
func (c *Console) WrapExecutor(executor prompt.Executor) prompt.Executor {
	return func(s string) {
		//log.Info("add history", "s", s)
		c.History.AddHistoryItem(s)
		executor(s)
	}
}
