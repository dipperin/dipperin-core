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
	"github.com/dipperin/dipperin-core/common"
	chain_config "github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestMakeCsProtocolManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// txpool
	mockTxPool := NewMockTxPool(ctrl)

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier).AnyTimes()

	// peer manager
	mockPeerManager := NewMockPeerManager(ctrl)

	// chain
	mockChain := NewMockChain(ctrl)

	// p2p server
	mockP2pServer := NewMockP2PServer(ctrl)

	// verifier
	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	// PbftSigner
	mockPbftSigner := NewMockPbftSigner(ctrl)

	cfg := &CsProtocolManagerConfig{
		ChainConfig:     *chain_config.GetChainConfig(),
		Chain:           mockChain,
		P2PServer:       mockP2pServer,
		NodeConf:        mockNodeConf,
		VerifiersReader: mockVerifiersReader,
		PbftNode:        mockPbftNode,
		MsgSigner:       mockPbftSigner,
	}

	mockPbftDecoder := NewMockP2PMsgDecoder(ctrl)

	tbCfg := &NewTxBroadcasterConfig{
		P2PMsgDecoder: mockPbftDecoder,
		TxPool:        mockTxPool,
		NodeConf:      mockNodeConf,
		Pm:            mockPeerManager,
	}

	csPm, bd := MakeCsProtocolManager(cfg, tbCfg)

	assert.NotNil(t, csPm)
	assert.NotNil(t, bd)
}

func TestNewBroadcastDelegate(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// txpool
	mockTxPool := NewMockTxPool(ctrl)

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)

	// peer manager
	mockPeerManager := NewMockPeerManager(ctrl)

	// chain
	mockChain := NewMockChain(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	bd := NewBroadcastDelegate(mockTxPool, mockNodeConf, mockPeerManager, mockChain, mockPbftNode)

	assert.NotNil(t, bd)
}

func TestBroadcastDelegate_BroadcastMinedBlock(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// txpool
	mockTxPool := NewMockTxPool(ctrl)

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)

	// peer manager
	mockPeerManager := NewMockPeerManager(ctrl)

	// chain
	mockChain := NewMockChain(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	bd := NewBroadcastDelegate(mockTxPool, mockNodeConf, mockPeerManager, mockChain, mockPbftNode)

	assert.NotNil(t, bd)

	mockPeerManager.EXPECT().GetPeers().Return(nil)

	bd.BroadcastMinedBlock(model.NewBlock(model.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil))
}

func TestBroadcastDelegate_BroadcastTx(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// txpool
	mockTxPool := NewMockTxPool(ctrl)

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)

	// peer manager
	mockPeerManager := NewMockPeerManager(ctrl)

	// chain
	mockChain := NewMockChain(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	bd := NewBroadcastDelegate(mockTxPool, mockNodeConf, mockPeerManager, mockChain, mockPbftNode)

	assert.NotNil(t, bd)

	mockPeerManager.EXPECT().GetPeers().Return(nil)

	gasPrice := big.NewInt(1)
	gasLimit := 2 * model.TxGas

	tx := model.NewTransaction(11, common.StringToAddress("dsad"), big.NewInt(11),
		gasPrice, gasLimit, nil)

	bd.BroadcastTx([]model.AbstractTransaction{tx})

}

func TestBroadcastDelegate_BroadcastEiBlock(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// txpool
	mockTxPool := NewMockTxPool(ctrl)

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)

	// peer manager
	mockPeerManager := NewMockPeerManager(ctrl)

	// chain
	mockChain := NewMockChain(ctrl)
	mockChain.EXPECT().GetSeenCommit(gomock.Any()).Return([]model.AbstractVerification{model.NewVoteMsg(11, 1, common.HexToHash("asdd"), model.AliveVerifierVoteMessage)})

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	bd := NewBroadcastDelegate(mockTxPool, mockNodeConf, mockPeerManager, mockChain, mockPbftNode)

	assert.NotNil(t, bd)

	mockPeerManager.EXPECT().GetPeers().Return(nil)

	block := model.NewBlock(model.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	bd.BroadcastEiBlock(block)
}
