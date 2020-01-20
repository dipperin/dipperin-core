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

package commands

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/rpcinterface"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"strings"
)

// getP2pRpcMethodByName get the rpc method name based on the method name
func getP2pRpcMethodByName(mName string) string {
	lm := strings.ToLower(string(mName[0])) + mName[1:]
	return "p2p_" + lm
}

func (caller *rpcCaller) AddPeer(c *cli.Context) {
	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	if len(cParams) != 1 {
		l.Error("url cannot be empty")
		return
	}

	var resp error
	if err = client.Call(&resp, getP2pRpcMethodByName(mName), cParams[0]); err != nil {
		l.Error("add peer error", zap.Error(err))
		return
	}

	if resp != nil {
		l.Error("add peer resp error", zap.Error(err))
		return
	}

	l.Info("add peer success")
}

func (caller *rpcCaller) Peers(c *cli.Context) {
	mName, _, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error")
		return
	}

	var resp []*p2p.PeerInfo
	if err = client.Call(&resp, getP2pRpcMethodByName(mName)); err != nil {
		l.Error("get peers error", zap.Error(err))
		return
	}

	for i := range resp {
		fmt.Println(util.StringifyJson(resp[i]))
	}
}

func (caller *rpcCaller) Debug(c *cli.Context) {
	var respBlock rpcinterface.BlockResp
	if err := client.Call(&respBlock, getDipperinRpcMethodByName("CurrentBlock")); err != nil {
		l.Error("look up for current block", zap.Error(err))
		return
	}

	printBlockInfo(respBlock)
	var resp p2p.CsPmPeerInfo
	if err := client.Call(&resp, getP2pRpcMethodByName("CsPmInfo")); err != nil {
		l.Error("add peer error", zap.Error(err))
		return
	}

	fmt.Println("base len:", len(resp.Base))
	fmt.Println("cur v len:", len(resp.CurVerifier))
	fmt.Println("next v len:", len(resp.NextVerifier))
	fmt.Println("v boot len:", len(resp.VerifierBoot))

	//fmt.Println(resp.String())
}
