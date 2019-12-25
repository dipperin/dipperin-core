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

package dipperin_prompts

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/soft-wallet"
	"github.com/manifoldco/promptui"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	RegPort = "^([0-9]|[1-9]\\d{1,3}|[1-5]\\d{4}|6[0-4]\\d{3}|65[0-4]\\d{2}|655[0-2]\\d|6553[0-5])$"
)

func NodeName() (string, error) {

	p := promptui.Prompt{
		Label:     "Node Name",
		Validate:  emptyValidate,
		Templates: PromptTemplate,
	}

	return p.Run()
}

func DataDir() (string, error) {

	defaultPath := filepath.Join(util.HomeDir(), "tmp/dipperin_apps/node")

	p := promptui.Prompt{
		Label:     fmt.Sprintf("Data Directory(default: %s)", defaultPath),
		Validate:  filepathValidate,
		Templates: PromptTemplate,
	}

	path, err := p.Run()

	if path == "" {
		path = defaultPath
	}

	return path, err
}

func P2PListener() (string, error) {

	validate := func(input string) error {
		if match, err := regexp.MatchString(RegPort, input); !match || err != nil {
			return fmt.Errorf("Not a valid p2p listener")
		}
		return nil
	}

	p := promptui.Prompt{
		Label:     "P2P Listener",
		Validate:  validate,
		Templates: PromptTemplate,
	}

	return p.Run()

}

func HTTPPort() (string, error) {
	p := promptui.Prompt{
		Label:     "HTTP Port",
		Validate:  portValidate,
		Templates: PromptTemplate,
	}

	return p.Run()
}

func WSPort() (string, error) {
	p := promptui.Prompt{
		Label:     "WebSocket Port",
		Validate:  portValidate,
		Templates: PromptTemplate,
	}

	return p.Run()
}

func WalletPassword() (string, error) {
	p := promptui.Prompt{
		Label:     "Wallet Password",
		Validate:  passwordValidate,
		Templates: PromptTemplate,
		Mask:      '*',
	}

	return p.Run()
}

func WalletPassPhrase() (string, error) {
	p := promptui.Prompt{
		Label:     "establish Wallet PassPhrase",
		Templates: PromptTemplate,
		Mask:      '*',
	}

	return p.Run()
}

func WalletPath(nodePath string) (string, error) {
	defaultPath := filepath.Join(util.HomeDir(), "tmp/dipperin_apps/node")
	p := promptui.Prompt{
		Label:     "open or establish Wallet path(default:"+defaultPath+"/CSWallet)",
		Validate: homePathValidate,
		Templates: PromptTemplate,
	}
	path, err := p.Run()
	if err != nil {
		return "", err
	}

	if path == "" {
		path = filepath.Join(nodePath, soft_wallet.WalletDefaultName)
	}

	return path, nil
}

func portValidate(port string) error {
	if match, err := regexp.MatchString(RegPort, port); !match || err != nil {
		return fmt.Errorf("Not a valid port")
	}
	return nil
}

func filepathValidate(path string) error {
	if path == "" {
		return nil
	}

	if strings.Contains(path,util.HomeDir()){
		return nil
	}
	//if !filepath.IsAbs(path) {
	//	return fmt.Errorf("Please enter an absolute path")
	//}

	return fmt.Errorf("choose or create a subdirectory in "+util.HomeDir())
}

func homePathValidate(path string) error {
	if path == "" {
		return nil
	}

	if strings.Contains(path,util.HomeDir()){
		return nil
	}

	return fmt.Errorf("choose or create a subdirectory in "+util.HomeDir())
}

func emptyValidate(input string) error {
	regNoCh := regexp.MustCompile("[\u4e00-\u9fa5]")
	chStr := regNoCh.FindAllString(input, -1)
	regBlank := regexp.MustCompile(" ")
	blankStr := regBlank.FindAllString(input, -1)

	if len(chStr) <= 0 && len(blankStr) <= 0{
		if len(input) == 0 {
			return fmt.Errorf("please provide a string input")
		}
		return nil
	}

	return errors.New("input illegal,cannot contain Chinese and spaces")
}

func passwordValidate(input string) error {
	reg := regexp.MustCompile(`[0-9]|[a-z]|[A-Z]|[~!@#$%^&*()_+<>?:"{},.\\/;'[\]` + "`]")
	regNoCh := regexp.MustCompile("[\u4e00-\u9fa5]")
	chStr := regNoCh.FindAllString(input, -1)
	regBlank := regexp.MustCompile(" ")
	blankStr := regBlank.FindAllString(input, -1)

	if reg.MatchString(input) && len(chStr) <= 0 && len(blankStr) <= 0 {
		if len(input) >= accounts.PasswordMin && len(input) <= accounts.PassWordMax {
			return nil
		}
	}
	return accounts.ErrPasswordOrPassPhraseIllegal
}
