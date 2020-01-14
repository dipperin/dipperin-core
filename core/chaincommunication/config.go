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

package chaincommunication

import (
	"errors"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"strings"
)

const (
	StatusMsg          = 0x00
	NewBlockHashesMsg  = 0x01
	TxMsg              = 0x02
	GetBlocksMsg       = 0x03
	BlocksMsg          = 0x04
	NewBlockMsg        = 0x07
	NewBlockByBloomMsg = 0x08

	// finder verifier
	GetVerifiersConnFromBootNode = 0x60
	BootNodeVerifiersConn        = 0x61

	// Estimator & InvBloom msg type
	EiNewBlockHashMsg    = 0x80
	EiEstimatorMsg       = 0x81
	EiNewBlockByBloomMsg = 0x82
	// for wait verify blocks
	EiWaitVerifyBlockHashMsg    = 0x83
	EiWaitVerifyEstimatorMsg    = 0x84
	EiWaitVerifyBlockByBloomMsg = 0x85

	TxV1Msg        = 0x71
	BlockHashesMsg = 0x72
	NewBlockV1Msg  = 0x73
	// verify result
	VerifyBlockHashResultMsg = 0x74
	GetVerifyResultMsg       = 0x75
	VerifyBlockResultMsg     = 0x76

	//verifier halt check protocol
	CurrentBlockNumberRequest    = 0x90
	CurrentBlockNumberResponse   = 0x91
	ProposeEmptyBlockMsg         = 0x92
	SendMinimalHashBlock         = 0x93
	SendMinimalHashBlockResponse = 0x94
)

const ProtocolMaxMsgSize = 10 * 1024 * 1024

var (
	msgTooLargeErr           = errors.New("msg too large")
	msgHandleFuncNotFoundErr = errors.New("msg handle func not found")
	quitErr                  = errors.New("download canceled")
)

const (
	MaxBlockFetch = 16
)

var totalVerifierBootNode int
var totalVerifier int
var PbftMaxPeerCount int

func init() {
	chinConfig := chain_config.GetChainConfig()
	totalVerifier = chinConfig.VerifierNumber
	totalVerifierBootNode = chinConfig.VerifierBootNodeNumber
	PbftMaxPeerCount = totalVerifier*2 + totalVerifierBootNode
}

func getRealRawUrl(remoteRawUrl string, remoteAddr string) string {
	i1 := strings.Index(remoteRawUrl, "@")
	i2 := strings.LastIndex(remoteRawUrl, ":")
	i3 := strings.LastIndex(remoteAddr, ":")

	return remoteRawUrl[:i1+1] + remoteAddr[:i3] + remoteRawUrl[i2:]
}
