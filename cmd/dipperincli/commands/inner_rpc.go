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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/urfave/cli"
	"strconv"
)

var InnerRpcForbid = false

func (caller *rpcCaller) GetBlockDiffVerifierInfo(c *cli.Context) {
	if InnerRpcForbid {
		l.Error("the rpc function isn't external")
		return
	}

	if checkSync() {
		return
	}

	mName, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if len(cParams) != 1 {
		l.Error("Sending a transaction requires parameter at least：blockNumber")
		return
	}

	blockNumber, err := strconv.Atoi(cParams[0])
	if err != nil {
		l.Error("blockNumber convert to int error")
		return
	}

	var resp map[economy_model.VerifierType][]common.Address
	if err = client.Call(&resp, getDipperinRpcMethodByName(mName), blockNumber); err != nil {
		l.Error("call failed", "err", err)
		return
	}

	fmt.Println("","the MasterVerifier address is:")
	printAddress(resp[economy_model.MasterVerifier])
	fmt.Println("","the CommitVerifier address is:")
	printAddress(resp[economy_model.CommitVerifier])
	fmt.Println("","the NotCommitVerifier address is:")
	printAddress(resp[economy_model.NotCommitVerifier])
}

func printAddress(addresses []common.Address) {
	for _, address := range addresses {
		fmt.Println("\t", "address:", address.Hex())
	}
}

func (caller *rpcCaller) CheckVerifierType(c *cli.Context) {
	if InnerRpcForbid {
		l.Error("the rpc function isn't external")
		return
	}

	if checkSync() {
		return
	}

	config := chain_config.GetChainConfig()
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if len(cParams) != 2 {
		l.Error("Sending a transaction requires Parameter at least：slot,address")
		return
	}

	slot, err := strconv.Atoi(cParams[0])
	if err != nil {
		l.Error("slot convert to int error")
		return
	}
	numberStart := uint64(slot) * config.SlotSize
	numberEnd := numberStart + config.SlotSize - 1

	// get currentBlockNumber
	var respBlock rpc_interface.BlockResp
	if err = client.Call(&respBlock, getDipperinRpcMethodByName("CurrentBlock")); err != nil {
		l.Error("look up for current block", "err", err)
		return
	}

	if numberStart > respBlock.Header.Number {
		l.Error("the slot more than current block number")
		return
	}

	if numberEnd > respBlock.Header.Number {
		numberEnd = respBlock.Header.Number
	}

	addr, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	for i := numberStart; i <= numberEnd; i++ {
		findType := false
		var resp map[economy_model.VerifierType][]common.Address
		if err = client.Call(&resp, getDipperinRpcMethodByName("GetBlockDiffVerifierInfo"), i); err != nil {
			l.Error("call failed", "err", err)
			return
		}

		for key, value := range resp {
			for _, address := range value {
				if addr == address {
					addressType := []string{"MasterVerifier", "CommitVerifier", "NotCommitVerifier"}
					findType = true
					l.Info("the address type is: ", addressType[key], "blockNumber", i)
				}
			}
		}

		if findType == false {
			l.Info("the address isn't verifier in the block", "blockNumber", i)
		}
	}
}

/*
func (caller *rpcCaller)SpecialCharDebug(c *cli.Context){
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}
	l.Info("the cParams is:","cParams",cParams)
}*/
