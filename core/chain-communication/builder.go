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

package chain_communication

import (
	"github.com/dipperin/dipperin-core/cmd/utils/debug"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chain-config"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
)

const (
	P2PMaxPeerCount    = 100
	NormalMaxPeerCount = 45
)

func NewBroadcastDelegate(txPool TxPool, nodeConfig NodeConf, peerManager PeerManager, chain Chain, pbftNode PbftNode) *BroadcastDelegate {
	txConfig := NewTxBroadcasterConfig{
		P2PMsgDecoder: MakeDefaultMsgDecoder(),
		TxPool:        txPool,
		NodeConf:      nodeConfig,
		Pm:            peerManager,
	}

	broadConfig := NewBlockBroadcasterConfig{
		Chain:    chain,
		Pm:       peerManager,
		PbftNode: pbftNode,
	}

	bftConfig := NewBftOuterConfig{
		Chain: chain,
		Pm:    peerManager,
	}
	return &BroadcastDelegate{
		newTxBroadcaster: makeNewTxBroadcaster(&txConfig),
		blockBroadcaster: makeNewBlockBroadcaster(&broadConfig),
		bftOut:           NewBftOuter(&bftConfig),
	}
}

type BroadcastDelegate struct {
	newTxBroadcaster *NewTxBroadcaster
	//eiBlockBroadcaster           *EiBlockBroadcaster
	//eiWaitVerifyBlockBroadcaster *EiWaitVerifyBlockBroadcaster

	blockBroadcaster *NewBlockBroadcaster
	bftOut           *BftOuter
}

func (delegate *BroadcastDelegate) BroadcastMinedBlock(block model.AbstractBlock) {
	log.DLogger.Debug("BroadcastBlock", zap.String("block id", block.Hash().Hex()), zap.Int("txs", block.TxCount()))
	delegate.blockBroadcaster.BroadcastBlock(block) //IBLT-method
}

func (delegate *BroadcastDelegate) BroadcastTx(txs []model.AbstractTransaction) {
	delegate.newTxBroadcaster.BroadcastTx(txs)
}

func (delegate *BroadcastDelegate) BroadcastEiBlock(block model.AbstractBlock) {
	vr := &model2.VerifyResult{
		Block:       block,
		SeenCommits: delegate.bftOut.Chain.GetSeenCommit(block.Number()),
	}
	delegate.bftOut.BroadcastVerifiedBlock(vr)
}

func MakeCsProtocolManager(pmConfig *CsProtocolManagerConfig, txBConf *NewTxBroadcasterConfig) (*CsProtocolManager, *BroadcastDelegate) {
	pm := newCsProtocolManager(pmConfig)
	txBConf.Pm = pm
	newTxBroadcaster := makeNewTxBroadcaster(txBConf)
	pm.registerCommunicationService(newTxBroadcaster, nil)
	pm.txSync = newTxBroadcaster

	blockBroadcaster := makeNewBlockBroadcaster(&NewBlockBroadcasterConfig{
		Pm:       pm,
		Chain:    pmConfig.Chain,
		PbftNode: pmConfig.PbftNode,
	})
	pm.registerCommunicationService(blockBroadcaster, nil)

	bftOut := NewBftOuter(&NewBftOuterConfig{
		Pm:    pm,
		Chain: pmConfig.Chain,
	})
	pm.registerCommunicationService(bftOut, nil)

	//eiBlockBroadcaster := makeEiBlockBroadcaster(&EiBlockBroadcasterConfig{
	//	Pm: pm,
	//	Chain: pmConfig.Chain,
	//	TxPool: txBConf.TxPool,
	//})
	//pm.registerCommunicationService(eiBlockBroadcaster, nil)
	//
	//eiWaitVerifyBlockBroadcaster := makeEiWaitVerifyBlockBroadcaster(&EiWaitVerifyBlockBroadcasterConfig{
	//	Pm: pm,
	//	NodeConf: pmConfig.NodeConf,
	//	TxPool: txBConf.TxPool,
	//	Chain: pmConfig.Chain,
	//})
	//pm.registerCommunicationService(eiWaitVerifyBlockBroadcaster, nil)

	//eiBlockFetcher := NewEiBlockFetcher(pmConfig.Chain.CurrentBlock, pmConfig.Chain.GetBlockByHash, pmConfig.Chain.SaveBlock, txBConf.TxPool.ConvertPoolToMap, eiBlockBroadcaster.BroadcastBlock)
	//pm.registerCommunicationService(nil, eiBlockFetcher)

	blockFetcher := NewBlockFetcher(pmConfig.Chain.CurrentBlock, pmConfig.Chain.GetBlockByHash, pmConfig.Chain.SaveBlock, bftOut.BroadcastVerifiedBlock)
	pm.registerCommunicationService(nil, blockFetcher)

	//wvEiBlockFetcher := NewWvEiBlockFetcher(&WvEiBlockFetcherConfig{
	//	PbftNode: pmConfig.PbftNode,
	//}, pmConfig.Chain.CurrentBlock, txBConf.TxPool.ConvertPoolToMap, eiWaitVerifyBlockBroadcaster.BroadcastBlock)
	//pm.registerCommunicationService(nil, wvEiBlockFetcher)
	//
	//eiWaitVerifyBlockBroadcaster.fetcher = wvEiBlockFetcher
	//
	//eiBlockBroadcaster.fetcher = eiBlockFetcher

	bftOut.SetBlockFetcher(blockFetcher)

	// have diff downloader
	downloader := MakeNewPbftDownloader(&NewPbftDownloaderConfig{
		Chain:    pmConfig.Chain,
		Pm:       pm,
		PbftNode: pmConfig.PbftNode,
		fetcher:  blockFetcher,
	})

	//downloader.SetFetcher(bftOuterFetcher)
	pm.registerCommunicationService(downloader, downloader)

	broadcastDelegate := &BroadcastDelegate{
		newTxBroadcaster: newTxBroadcaster,
		blockBroadcaster: blockBroadcaster,
		bftOut:           bftOut,
	}

	switch pmConfig.NodeConf.GetNodeType() {
	case chain_config.NodeTypeOfVerifier:
		pm.vf = NewVFinder(pmConfig.Chain, pm, pmConfig.ChainConfig)
		pm.registerCommunicationService(pm.vf, pm.vf)
	case chain_config.NodeTypeOfVerifierBoot:
		pm.registerCommunicationService(NewVFinderBoot(pm, pmConfig.Chain), nil)
	}

	// add mem size info
	if chain_config.GetCurBootsEnv() != chain_config.BootEnvMercury {
		debug.Memsize.Add("cs_protocol", pm)
		debug.Memsize.Add("newTxBroadcaster", newTxBroadcaster)
		debug.Memsize.Add("blockBroadcaster", blockBroadcaster)
		debug.Memsize.Add("bftOut", bftOut)
		debug.Memsize.Add("blockFetcher", blockFetcher)
		debug.Memsize.Add("downloader", downloader)
	}

	return pm, broadcastDelegate
}
