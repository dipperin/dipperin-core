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
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/common/g-timer"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"path/filepath"
	"sync/atomic"
	"time"
)

const (
	// Normal and MineMaster are both base
	base = 1
	// Verifier
	verifier = 2
	// Boot
	boot = 3
)

type CsProtocolManagerConfig struct {
	ChainConfig     chain_config.ChainConfig
	Chain           Chain
	P2PServer       P2PServer
	NodeConf        NodeConf
	VerifiersReader VerifiersReader
	PbftNode        PbftNode
	MsgSigner       PbftSigner
}

/*

Peer management logic: According to their own roles, the peers are processed differently and added.

1. I am Normal/MineMaster: I will insert all the peers into the nor peer set.
1. i am a verifier
	（1）I am the current round and the next round: judging that the other party type . If boot then is added to the boot, if the other party is not the current and the next round then is added to the ordinary, if it is the next round, then join the cur and next set, otherwise belong Which round then join that set
	（2）I am the current round: I judge the other party type ,if it is boot and add it to the boot. If the other party is not the current one, it is added to the normal one. It is judged that the other party is the current round and joins the current round.
	（3）I am the next round: judge the other party is boot then add to the boot, the other party is not the next round, then join the ordinary, judge the other party is the next round to join the next round

1. Self is boot: if the other party is boot then is added to the boot, the other party is the current round then join the current round, the other party is the next round then join the next round, the others are added to the ordinary

switch verifiers logic：
*/
type CsProtocolManager struct {
	BaseProtocolManager
	*CsProtocolManagerConfig
	// Maximum number of peer connections
	maxPeers int
	// Synchronize transactions in the tx pool during handshake
	txSync *NewTxBroadcaster
	pmType atomic.Value

	// finder for verifiers
	vf *VFinder

	// verifier boot nodes list
	verifierBootNodes []*enode.Node

	peerSetManager *CsPmPeerSetManager

	stop chan struct{}
}

func (pm *CsProtocolManager) ShowPmInfo() *p2p.CsPmPeerInfo {
	return &p2p.CsPmPeerInfo{
		SelfNodeID:   pm.P2PServer.Self().ID().String(),
		SelfNodeName: pm.NodeConf.GetNodeName(),
		SelfType:     uint64(pm.NodeConf.GetNodeType()),
		Base:         pm.peerSetManager.basePeers.GetPeersInfo(),
		CurVerifier:  pm.peerSetManager.currentVerifierPeers.GetPeersInfo(),
		NextVerifier: pm.peerSetManager.nextVerifierPeers.GetPeersInfo(),
		VerifierBoot: pm.peerSetManager.verifierBootNode.GetPeersInfo(),
	}
}

//
func (pm *CsProtocolManager) BestPeer() PmAbstractPeer {
	return pm.peerSetManager.BestPeer()
}

func (pm *CsProtocolManager) RemovePeer(id string) {
	pm.peerSetManager.RemovePeer(id)
}

func (pm *CsProtocolManager) GetPeer(id string) PmAbstractPeer {
	if p := pm.peerSetManager.basePeers.Peer(id); p != nil {
		return p
	}
	if p := pm.peerSetManager.currentVerifierPeers.Peer(id); p != nil {
		return p
	}
	if p := pm.peerSetManager.nextVerifierPeers.Peer(id); p != nil {
		return p
	}
	if p := pm.peerSetManager.verifierBootNode.Peer(id); p != nil {
		return p
	}

	return nil
}

func (pm *CsProtocolManager) GetPeers() map[string]PmAbstractPeer {
	result := map[string]PmAbstractPeer{}
	norPs := pm.peerSetManager.basePeers.GetPeers()
	curPs := pm.peerSetManager.currentVerifierPeers.GetPeers()
	nextPs := pm.peerSetManager.nextVerifierPeers.GetPeers()
	vbPs := pm.peerSetManager.verifierBootNode.GetPeers()

	mergePeersTo(norPs, result)
	mergePeersTo(curPs, result)
	mergePeersTo(nextPs, result)
	mergePeersTo(vbPs, result)
	//log.Info("get pm peers", "total", len(result))
	return result
}

func mergePeersTo(from, to map[string]PmAbstractPeer) {
	for k, p := range from {
		to[k] = p
	}
}

func (pm *CsProtocolManager) Protocols() []p2p.Protocol {
	if len(pm.protocols) != 0 {
		return pm.protocols
	}
	pm.protocols = []p2p.Protocol{pm.getCsProtocol()}
	return pm.protocols
}

func (pm *CsProtocolManager) getCsProtocol() p2p.Protocol {
	var version uint = chain_config.CsProtocolVersion
	// Use a different protocol to make it unable to connect in the underlying layer
	protocolName := chain_config.AppName + "_cs_local"
	switch chain_config.GetCurBootsEnv() {
	case "mercury":
		log.Info("use mercury cs protocol")
		protocolName = chain_config.AppName + "_cs"
	case "test":
		log.Info("use test cs protocol")
		protocolName = chain_config.AppName + "_cs_test"
	case "venus":
		log.Info("use test cs protocol")
		protocolName = chain_config.AppName + "_vs"
	default:
		log.Info("use local cs protocol")
	}
	p := p2p.Protocol{Name: protocolName, Version: version, Length: 0x200}
	p.Run = func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
		g_metrics.Add(g_metrics.TotalHandledPeer, "", 1)
		log.Info("getPbftProtocol new pbft peer in", "protocol", protocolName)
		// format with communication peer
		tmpPmPeer := newPeer(int(version), peer, rw)
		pm.wg.Add(1)
		defer pm.wg.Done()
		// read msg loop in here
		return pm.handle(tmpPmPeer)
	}
	return p
}

func newCsProtocolManager(config *CsProtocolManagerConfig) *CsProtocolManager {
	pm := &CsProtocolManager{
		BaseProtocolManager: BaseProtocolManager{
			msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
		},
		CsProtocolManagerConfig: config,
		maxPeers:                P2PMaxPeerCount,
		verifierBootNodes:       chain_config.VerifierBootNodes,
		stop:                    make(chan struct{}),
	}

	// load msg code & method
	//pm.msgHandlers[GetCurrentVerifierFromBootNode] = pm.onGetBootNodeCurrentVerifier
	//pm.msgHandlers[BootNodeCurrentVerifier] = pm.onBootNodeCurrentVerifier
	//pm.msgHandlers[GetNextVerifierFromBootNode] = pm.onGetBootNodeNextVerifier
	//pm.msgHandlers[BootNodeNextVerifier] = pm.onBootNodeNextVerifier

	// check the maximum number of connections
	if pm.maxPeers < PbftMaxPeerCount {
		panic("max peer count should >= bft max peer count")
	}

	psManager := newCsPmPeerSetManager(pm.selfPmType(), pm.maxPeers, pm.SelfIsNextVerifier, pm.SelfIsCurrentVerifier,
		pm.isCurrentVerifierNode, pm.isNextVerifierNode, pm.isVerifierBootNode)

	pm.peerSetManager = psManager

	return pm
}

func (pm *CsProtocolManager) Start() error {
	if err := pm.BaseProtocolManager.Start(); err != nil {
		return err
	}

	// debug information, delete when not in use
	go pm.logCurPeersInfo()

	go func() {
		if util.IsTestEnv() {
			return
		}

		if err := pm.handleInsertEventForBft(); err != nil {
			log.Error("cs protocol manager handle insert event failed", "err", err)
			return
		}
	}()

	if !util.IsTestEnv() {
		go pm.bootVerifierConnCheck()
	}
	// if self is cur or next verifier, will add boot node

	return nil
}

func (pm *CsProtocolManager) logCurPeersInfo() {
	if util.IsTestEnv() {
		return
	}

	tick := g_timer.SetPeriodAndRun(pm.PrintPeerHealthCheck, 15*time.Second)
	defer g_timer.StopWork(tick)

	<-pm.stop
}

func (pm *CsProtocolManager) Stop() {
	pm.BaseProtocolManager.Stop()

	// need to stop the finder here
	if pm.stop == nil {
		return
	}

	close(pm.stop)
	pm.stop = nil
}

func (pm *CsProtocolManager) BroadcastMsg(msgCode uint64, msg interface{}) {
	vPeers := pm.peerSetManager.currentVerifierPeers.GetPeers()
	log.PBft.Info("broadcast msg to pbft nodes", "msg code", msgCode, "peer len", len(vPeers))
	for _, p := range vPeers {
		//log.PBft.Info("broadcast msg to pbft nodes", "msg code", msgCode, "to", p.NodeName())
		if err := p.SendMsg(msgCode, msg); err != nil {
			log.Warn("broadcast pbft msg failed", "to", p.NodeName(), "msg code", msgCode, "err", err)
		}
	}
}

func (pm *CsProtocolManager) BroadcastMsgToTargetVerifiers(msgCode uint64, from []common.Address, msg interface{}) {
	log.PBft.Debug("Broadcast msg to targets called", "len", len(from), "cur v len", pm.peerSetManager.currentVerifierPeers.Len())
	vPeers := pm.peerSetManager.currentVerifierPeers.GetPeers()
	for _, add := range from {
		for _, p := range vPeers {
			if p.RemoteVerifierAddress().IsEqual(add) {
				log.PBft.Debug("send fetch round msg", "to", p.NodeName())
				if err := p.SendMsg(msgCode, msg); err != nil {
					log.Warn("broadcast pbft msg to target verifier failed", "to", p.NodeName(),
						"msg code", msgCode, "err", err)
					log.PBft.Warn("broadcast pbft msg to target verifier failed", "to", p.NodeName(),
						"msg code", msgCode, "err", err)
				}
			}
		}
	}
}

func (pm *CsProtocolManager) SendFetchBlockMsg(msgCode uint64, from common.Address, msg *model.FetchBlockReqDecodeMsg) error {
	vPeers := pm.peerSetManager.currentVerifierPeers.GetPeers()
	for _, p := range vPeers {
		if p.RemoteVerifierAddress().IsEqual(from) {
			log.PBft.Info("send fetch block msg", "to", p.NodeName())
			return p.SendMsg(msgCode, msg)
		}
	}
	return errors.New("no verifier peer for fetcher")
}

// change verifier，
//This method is only triggered when the change is made, so if there is a problem when in the peer handling process, it is difficult to correct it.
func (pm *CsProtocolManager) ChangeVerifiers() {
	log.PBft.Info("Change verifiers", "is new slot verifier", pm.SelfIsCurrentVerifier())
	vReader := pm.VerifiersReader
	nextVerifiers := vReader.NextVerifiers()
	if pm.NodeConf.GetNodeType() == verifier {
		pm.PbftNode.ChangePrimary(nextVerifiers[0].Hex())
	}
	// todo All the change points to get Next are going to get next next
	pm.peerSetManager.OrganizeVerifiersSet(nextVerifiers, []common.Address{})
}

// only for finder
func (pm *CsProtocolManager) HaveEnoughVerifiers(withOrganizeVSet bool) (mc uint, mn uint) {
	// exclude self
	vShouldLen := pm.ChainConfig.VerifierNumber - 1
	cLen := pm.peerSetManager.currentVerifierPeers.Len()
	nLen := pm.peerSetManager.nextVerifierPeers.Len()

	if withOrganizeVSet && (cLen < vShouldLen || nLen < vShouldLen) {
		vReader := pm.VerifiersReader
		pm.peerSetManager.OrganizeVerifiersSet(vReader.CurrentVerifiers(), vReader.NextVerifiers())
		cLen = pm.peerSetManager.currentVerifierPeers.Len()
		nLen = pm.peerSetManager.nextVerifierPeers.Len()
	}

	missCur := vShouldLen - cLen
	missNext := vShouldLen - nLen

	if missCur <= 0 {
		mc = 0
		log.Warn("too many v peers", "cur len", cLen)
	} else {
		mc = uint(missCur)
	}

	if missNext <= 0 {
		mn = 0
		log.Warn("too many v peers", "next len", nLen)
	} else {
		mn = uint(missNext)
	}

	return
}

// returns the map key of the current verifier -- > node id, value -- > address
func (pm *CsProtocolManager) GetCurrentConnectPeers() map[string]common.Address {
	peerInfo := make(map[string]common.Address)

	if pm.selfPmType() != verifier {
		return peerInfo
	}

	currentConnectPeers := pm.peerSetManager.currentVerifierPeers.GetPeers()

	for nodeId, peer := range currentConnectPeers {
		peerInfo[nodeId] = peer.RemoteVerifierAddress()
	}

	return peerInfo
}

func (pm *CsProtocolManager) GetVerifierBootNode() map[string]PmAbstractPeer {
	return pm.peerSetManager.verifierBootNode.GetPeers()
}

func (pm *CsProtocolManager) GetNextVerifierPeers() map[string]PmAbstractPeer {
	return pm.peerSetManager.nextVerifierPeers.GetPeers()
}

func (pm *CsProtocolManager) CurrentVerifierPeersSet() AbstractPeerSet {
	return pm.peerSetManager.currentVerifierPeers
}

func (pm *CsProtocolManager) NextVerifierPeersSet() AbstractPeerSet {
	return pm.peerSetManager.nextVerifierPeers
}

func (pm *CsProtocolManager) SelfIsBootNode() bool {
	if pm.selfPmType() == boot {
		return true
	}

	return false
}

func (pm *CsProtocolManager) GetSelfNode() *enode.Node {
	return pm.P2PServer.Self()
}

func (pm *CsProtocolManager) ConnectPeer(node *enode.Node) {
	pm.P2PServer.AddPeer(node)
}

// prepare the next round of verifier
func (pm *CsProtocolManager) MatchCurrentVerifiersToNext() {
	// If it is a new block to do this, next should be the real next
	if !pm.SelfIsNextVerifier() {
		return
	}

	vReader := pm.VerifiersReader
	nextVPeersLen := pm.peerSetManager.nextVerifierPeers.Len()

	if nextVPeersLen == (totalVerifier - 1) {
		return
	}
	log.Health.Info("MatchCurrentVerifiersToNext", "next p len", nextVPeersLen, "total", totalVerifier)

	nextVs := vReader.NextVerifiers()
	pm.pickNextVerifierFromPs(nextVs)
	pm.pickNextVerifierFromCps(nextVs)

}

// Check if there is a next round in the normal peer set
func (pm *CsProtocolManager) pickNextVerifierFromPs(nextVs []common.Address) {
	// check if there are any normal nodes
	for _, cp := range pm.peerSetManager.basePeers.GetPeers() {
		for _, n := range nextVs {
			if cp.RemoteVerifierAddress().IsEqual(n) {

				// If it already exists in next, it should be removed from the normal
				if pm.peerSetManager.nextVerifierPeers.Peer(cp.ID()) != nil {
					if err := pm.peerSetManager.basePeers.RemovePeer(cp.ID()); err != nil {
						log.Error("remove peer from peer set failed", "err", err)
					}
					break
				}

				if err := pm.peerSetManager.nextVerifierPeers.AddPeer(cp); err != nil {
					log.Error("add peer to next verifier set failed", "err", err)
					break
				}

				// need to be removed from normal
				if err := pm.peerSetManager.basePeers.RemovePeer(cp.ID()); err != nil {
					log.Error("remove peer from peer set failed", "err", err)
				}

				break
			}
		}
	}
}

// check is there is the next round verifier in the current verifier set
func (pm *CsProtocolManager) pickNextVerifierFromCps(nextVs []common.Address) {
	// check if there is in the normal peers
	for _, cp := range pm.peerSetManager.currentVerifierPeers.GetPeers() {
		for _, n := range nextVs {
			if cp.RemoteVerifierAddress().IsEqual(n) {

				// If it already exists in next, it should jump out of the loop
				if pm.peerSetManager.nextVerifierPeers.Peer(cp.ID()) != nil {
					break
				}

				if err := pm.peerSetManager.nextVerifierPeers.AddPeer(cp); err != nil {
					log.Error("add peer to next verifier set failed", "err", err)
				}

				break
			}
		}
	}
}

// get the type of protocol run
func (pm *CsProtocolManager) selfPmType() int {
	if tp := pm.pmType.Load(); tp != nil {
		return tp.(int)
	}

	// check node type
	nodeType := pm.NodeConf.GetNodeType()

	if nodeType == chain_config.NodeTypeOfNormal || nodeType == chain_config.NodeTypeOfMineMaster {
		pm.pmType.Store(base)
		return base
	}

	if nodeType == chain_config.NodeTypeOfVerifier {
		pm.pmType.Store(verifier)
		return verifier
	}

	if nodeType == chain_config.NodeTypeOfVerifierBoot {

		curNodeID := pm.P2PServer.Self().ID().String()
		//log.Info("the node id is:","curNodeId",curNodeID)

		for _, n := range pm.verifierBootNodes {
			if curNodeID == n.ID().String() {
				pm.pmType.Store(boot)
				return boot
			}
		}
	}

	panic(fmt.Sprintf("illegal node type: %v. nodekey is wrong if is v boot", nodeType))
}

// determine whether the peer is a verifier boot node
func (pm *CsProtocolManager) isVerifierBootNode(p PmAbstractPeer) bool {
	for _, bn := range pm.verifierBootNodes {
		//log.Info("-----------------check remote peer is boot node", "saved b", bn.ID.String(), "p id", p.ID())
		if p.ID() == bn.ID().String() {
			return true
		}
	}
	return false
}

// Determine if the node is the next round of verifier
func (pm *CsProtocolManager) isNextVerifierNode(p PmAbstractPeer) bool {
	vReader := pm.VerifiersReader

	result := vReader.ShouldChangeVerifier()
	nextNodes := make([]common.Address, 0)
	if result {
		//When it is detected that the verifiers should be changed, the next verifier in pbft should be the next verifier after the change.
		//The next verifier interface is not provided. Here, it returns true first, that is, not connect the next verifier of success change in the change point.
		return false
	} else {
		nextNodes = vReader.NextVerifiers()
	}

	for _, n := range nextNodes {
		if n.IsEqual(p.RemoteVerifierAddress()) {
			return true
		}
	}
	return false
}

// Determine if the node is the current round verifier
func (pm *CsProtocolManager) isCurrentVerifierNode(p PmAbstractPeer) bool {
	vReader := pm.VerifiersReader

	result := vReader.ShouldChangeVerifier()
	currentVerifiers := make([]common.Address, 0)
	if result {
		//When it is detected that the verifiers should be changed, the current verifier in pbft should be the verifier after the change. Otherwise the node that was disconnected at change will be reconnected
		currentVerifiers = vReader.NextVerifiers()
	} else {
		currentVerifiers = vReader.CurrentVerifiers()
	}

	for _, n := range currentVerifiers {
		if n.IsEqual(p.RemoteVerifierAddress()) {
			return true
		}
	}
	return false
}

// determine if you are current verifier
func (pm *CsProtocolManager) SelfIsCurrentVerifier() bool {
	vReader := pm.VerifiersReader
	pbftSigner := pm.MsgSigner

	shouldChange := vReader.ShouldChangeVerifier()

	curs := vReader.CurrentVerifiers()
	ns := vReader.NextVerifiers()

	//If it is a change block, the fetched next is actually the real current.
	if shouldChange {
		curs = ns
		ns = []common.Address{}
		log.Info("check self is cur verifier, cur block should change verifier, so next trans to cur")
	}

	baseAddr := pbftSigner.GetAddress()

	//log.Pm.Info("check self is current verifier","selfAddr",baseAddr.Hex())
	//log.Pm.Info("check self is current verifier","currentVer",curs)

	if baseAddr.InSlice(curs) {
		return true
	}

	return false
}

//determine if you are a verifier boot node
func (pm *CsProtocolManager) selfIsVerifierBootNode() bool {
	baseAddr := pm.MsgSigner.GetAddress()
	if baseAddr.InSlice(chain_config.VerBootNodeAddress) {
		return true
	}

	return false
}

// determine if you are the next round of verifier
// It is important to consider how this next concept should be defined after the block is inserted.
//isRemoteNext is used to mark whether you want to take the real next round when at the change point (because the taken next    is actually the upcoming verifier, which is the real current)
func (pm *CsProtocolManager) SelfIsNextVerifier() bool {
	vReader := pm.VerifiersReader
	pbftSigner := pm.MsgSigner

	ns := vReader.NextVerifiers()
	shouldChange := vReader.ShouldChangeVerifier()
	if shouldChange {
		ns = []common.Address{}
		log.Info("check self is next verifier, cur block should change verifier, so next is remote next")
	}

	baseAddr := pbftSigner.GetAddress()
	if baseAddr.InSlice(ns) {
		return true
	}
	return false
}

// check the number of connections
func (pm *CsProtocolManager) checkConnCount() bool {
	switch pm.selfPmType() {
	case base:
		if pm.peerSetManager.basePeers.Len() >= pm.maxPeers {
			return true
		}
	case verifier:
		if pm.peerSetManager.currentVerifierPeers.Len() >= PbftMaxPeerCount && pm.peerSetManager.basePeers.Len() >= (pm.maxPeers-PbftMaxPeerCount) {
			return true
		}
	case boot:
		if pm.peerSetManager.currentVerifierPeers.Len() >= PbftMaxPeerCount && pm.peerSetManager.basePeers.Len() >= (pm.maxPeers-PbftMaxPeerCount) {
			return true
		}
	}

	return false
}

// peer handle msg
func (pm *CsProtocolManager) handle(p PmAbstractPeer) error {
	g_metrics.Add(g_metrics.CurHandelPeer, "", 1)
	defer g_metrics.Sub(g_metrics.CurHandelPeer, "", 1)

	// check the number of connections
	if pm.checkConnCount() {
		log.Warn("too many peers, can't add new peer")
		g_metrics.Add(g_metrics.TotalFailedHandle, "", 1)
		return p2p.DiscTooManyPeers
	}

	if err := pm.HandShake(p); err != nil {
		g_metrics.Add(g_metrics.TotalFailedHandle, "", 1)
		log.Warn("CsProtocolManager hand shake failed", "err", err, "remote host", p.RemoteAddress())
		return err
	}

	// determine the same address repeated connection
	if pm.isCurrentVerifierNode(p) {
		for _, peer := range pm.peerSetManager.currentVerifierPeers.GetPeers() {
			if peer.RemoteVerifierAddress().IsEqual(p.RemoteVerifierAddress()) {
				g_metrics.Add(g_metrics.TotalFailedHandle, "", 1)
				return errors.New("current verifier address already in peer set")
			}
		}
	}

	// determine the same address repeated connection
	if pm.isNextVerifierNode(p) {
		for _, peer := range pm.peerSetManager.nextVerifierPeers.GetPeers() {
			if peer.RemoteVerifierAddress().IsEqual(p.RemoteVerifierAddress()) {
				g_metrics.Add(g_metrics.TotalFailedHandle, "", 1)
				return errors.New("next verifier address already in peer set")
			}
		}
	}

	// The add peer that is repeated will report an error, so there is no need to check if there is a duplicate peer.
	if err := pm.peerSetManager.AddPeer(p); err != nil {
		log.Warn("add peer to peer set failed", "err", err)

		g_metrics.Add(g_metrics.TotalFailedHandle, "", 1)
		return err
	}

	// add the condition after add succeeds
	defer func() {
		// rm peer && disconnect
		pm.peerSetManager.RemovePeer(p.ID())
	}()

	// Propagate existing transactions. new transactions appearing
	//pm.txSync.syncTxs(p)

	g_metrics.Add(g_metrics.TotalSuccessHandle, "", 1)

	for {
		if err := pm.handleMsg(p); err != nil {
			log.Error("handle peer msg failed", "err", err, "p name", p.NodeName())

			if InPmBrokenError(err) {
				p.SetNotRunning()
				return err
			} else {
				log.Info("handleMsg err is not broken err, do not disconnect", "err", err)
				// todo This is not very good, but can avoid the for engage completely CPU
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

}

func (pm *CsProtocolManager) handleMsg(p PmAbstractPeer) error {

	//log.Info("the protocolManager connected peer number is:","number",len(pm.peers.GetPeers()))

	msg, err := p.ReadMsg()

	if err != nil {
		log.Info("base protocol read msg from peer failed", "err", err, "peer name", p.NodeName())
		log.Pm.Info("base protocol read msg from peer failed", "node", p.NodeName(), "err", err)
		return err
	}

	// todo only for debug
	finishedChan := make(chan struct{})
	go func() {
		select {
		case <-finishedChan:
			// It is likely that the trading pool is too full and it takes a lot of time to do IBLT parsing.
		case <-time.After(10 * time.Second):
			// log debug stack
			debugStackFile := filepath.Join(util.HomeDir(), "tmp", "cs_debug", "stack", pm.NodeConf.GetNodeName(), time.Now().Format("2006-1-2 15:04:05")+".log")
			if err := util.WriteDebugStack(debugStackFile); err != nil {
				log.Error("write debug stack failed", "err", err)
			}
			//panic(fmt.Sprintf("handle msg use more than 10s, msg code: 0x%x", msg.Code))
			log.Error(fmt.Sprintf("handle msg use more than 10s, msg code: 0x%x, remote node: %v. disconnect this peer, and write debug stack to: %v", msg.Code, p.NodeName(), debugStackFile))
			// not disconnect, the handle may be indeed time-consuming
			//pm.removePeerFromAllSet(p.ID())
		}
	}()

	defer func() {
		_ = msg.Discard()
		timer := time.NewTimer(2 * time.Second)
		select {
		case finishedChan <- struct{}{}:
		case <-timer.C:
			log.Error("can't write to finishedChan. means handle msg finished, but use more than 10s")
		}
		timer.Stop()
	}()

	if msg.Size > ProtocolMaxMsgSize {
		return msgTooLargeErr
	}

	// msg to bft node
	if pm.selfPmType() != base && uint64(msg.Code) > 0x100 {
		// handle this msg
		if err = pm.PbftNode.OnNewP2PMsg(msg, p); err != nil {
			log.PBft.Error("handle pbft msg failed", "err", err, "msg code", fmt.Sprintf("%x", msg.Code))
			return err
		}

		// if the message is processed then jump out
		return nil
	}

	// I am the current round verifier, but the node type that is started is not verifier
	if uint64(msg.Code) > 0x100 {
		log.Warn("self in current verifiers, but node not start with verifier type", "cur node type", pm.NodeConf.GetNodeType())
		return nil
	}
	// find handler for this msg
	tmpHandler := pm.msgHandlers[uint64(msg.Code)]
	if tmpHandler == nil {
		log.Error("Get message processing error", "msg code", fmt.Sprintf("0x%x", msg.Code))
		return msgHandleFuncNotFoundErr
	}

	// handle this msg
	if err = tmpHandler(msg, p); err != nil {
		log.Warn("handle msg failed", "err", err, "msg code", fmt.Sprintf("%x", msg.Code))
		return err
	}

	return nil
}

// handle handshake
func (pm *CsProtocolManager) HandShake(p PmAbstractPeer) error {
	chainReader := pm.Chain
	chainConf := pm.ChainConfig
	genesisBlock := chainReader.GetBlockByNumber(0)
	nodeConf := pm.NodeConf
	pbftSigner := pm.MsgSigner

	statusDataChan := make(chan *StatusData)

	go func() {
		curB := chainReader.CurrentBlock()
		//log.Info("send hand shake msg", "cur block", curB.Number())
		sData := StatusData{
			HandShakeData: HandShakeData{
				ChainID:            chainConf.ChainId,
				NetworkId:          chainConf.NetworkID,
				ProtocolVersion:    chain_config.CsProtocolVersion,
				NodeType:           uint64(nodeConf.GetNodeType()),
				NodeName:           nodeConf.GetNodeName(),
				CurrentBlockHeight: curB.Number(),
				CurrentBlock:       curB.Hash(),
				GenesisBlock:       genesisBlock.Hash(),
				RawUrl:             pm.P2PServer.Self().String(),
			},
			//NodeType:
		}
		log.Debug("before sign hand shake msg", "data hash", sData.DataHash().Hex())
		log.Debug("the node type is:", "type", nodeConf.GetNodeType())
		if nodeConf.GetNodeType() != chain_config.NodeTypeOfNormal {
			// sign
			if signB, err := pbftSigner.SignHash(sData.DataHash().Bytes()); err != nil {
				// send even if there is an error
				log.Error("sign status data hash failed", "err", err)
			} else {
				sData.Sign = signB
				sData.PubKey = crypto.CompressPubkey(pbftSigner.PublicKey())
			}
		}

		if err := p.SendMsg(StatusMsg, sData); err != nil {
			log.Error("send status msg error", "err", err)
		}
		log.Debug("send hand shake message success")
	}()

	go func() {
		msg, err := p.ReadMsg()
		if err != nil {
			log.Error("read handshake response error", "err", err)
			statusDataChan <- nil
			return
		}
		defer func() {
			if err := msg.Discard(); err != nil {
				log.Error("discard status msg err")
			}
		}()

		var tmpStatus StatusData
		if err = msg.Decode(&tmpStatus); err != nil {
			log.Warn("decode hand shake msg failed", "err", err)
		}

		log.Debug("read  handshake data is:", "tmpStatus", tmpStatus)
		statusDataChan <- &tmpStatus
	}()

	select {
	case remoteStatus := <-statusDataChan:
		if remoteStatus == nil {
			return errors.New("can't read hand shake msg from remote")
		}
		if remoteStatus.NetworkId != chainConf.NetworkID {
			log.Error(fmt.Sprintf("network id not match, remote: %v local: %v", remoteStatus.NetworkId, chainConf.NetworkID))
			return errors.New("network id not match")
		}
		if !genesisBlock.Hash().IsEqual(remoteStatus.GenesisBlock) {

			log.Error(fmt.Sprintf("genesis block not match, local: %v remote: %v", genesisBlock.Hash(), remoteStatus.GenesisBlock))
			return errors.New("genesis block not match")
		}
		if remoteStatus.ProtocolVersion == 0 {
			return errors.New("can't read hand shake msg")
		}

		if remoteStatus.ProtocolVersion != chain_config.CsProtocolVersion {
			return errors.New("cs protocol version not match")
		}

		// If the other party does not sign, they will get an empty address.
		verifierAddress := remoteStatus.Sender()
		p.SetRemoteVerifierAddress(verifierAddress)

		p.SetNodeType(remoteStatus.NodeType)
		p.SetNodeName(remoteStatus.NodeName)
		p.SetHead(remoteStatus.CurrentBlock, remoteStatus.CurrentBlockHeight)
		remoteStatus.RawUrl = getRealRawUrl(remoteStatus.RawUrl, p.RemoteAddress().String())
		p.SetPeerRawUrl(remoteStatus.RawUrl)

		log.Info("cs protocol hand shake success", "remote", remoteStatus.NodeName, "remote bh", remoteStatus.CurrentBlockHeight, "remote nt", remoteStatus.NodeType, "raw url", remoteStatus.RawUrl)
	}

	return nil

}

//check and connect verifier boot nodes
//If it is the current or next round verifier, then connect the unconnected verifier boot nodes
//If it is a verifier boot node, connect to other unconnected verifier boot nodes
//If you are a different type of node, disconnect your own connected verifier boot nodes
func (pm *CsProtocolManager) checkAndHandleVerBootNodes() {
	if pm.selfPmType() == base {
		return
	}

	if pm.chainHeightTooLow() {
		//log.Info("chain height too low, disconnect v boots")
		//pm.disconnectVBoots()
		return
	}

	if pm.SelfIsCurrentVerifier() || pm.SelfIsNextVerifier() || pm.SelfIsBootNode() {
		pm.connectVBoots()
	} else {
		//pm.disconnectVBoots()
	}

	return
}

//func (pm *CsProtocolManager) disconnectVBoots() {
//	if pm.peerSetManager.verifierBootNode.Len() == 0 {
//		return
//	}
//
//	log.Info("do disconnectVBoots")
//
//	for _, vbNode := range chain_config.VerifierBootNodes {
//		pm.P2PServer.RemovePeer(vbNode)
//		pm.peerSetManager.removeVerifierBootNodePeers(vbNode.ID().String())
//	}
//}

func (pm *CsProtocolManager) connectVBoots() {
	chainConfig := chain_config.GetChainConfig()

	if pm.peerSetManager.verifierBootNode.Len() == chainConfig.VerifierBootNodeNumber-1 {
		return
	}

	log.Info("do connectVBoots", "chain_config.VerifierBootNodes len", len(chain_config.VerifierBootNodes))

	selfID := pm.P2PServer.Self().ID().String()
	for _, vbNode := range chain_config.VerifierBootNodes {
		vbID := vbNode.ID().String()
		if vbID == selfID {
			log.Info("v boot is cur node", "id", selfID)
			continue
		}

		if pm.peerSetManager.verifierBootNode.Peer(vbID) != nil {
			log.Info("v boot in verifierBootSet", "vbID", vbID, "self id", selfID)
			continue
		}
		log.Info("add v boot", "conn", vbNode.String())
		pm.P2PServer.AddPeer(vbNode)
	}
}

func (pm *CsProtocolManager) chainHeightTooLow() bool {
	bp := pm.BestPeer()
	if bp == nil {
		log.Error("can't get best peer, for chainHeightTooLow check")
		return true
	}
	_, h := bp.GetHead()
	curB := pm.Chain.CurrentBlock()
	if curB.Number()+2 < h {
		return true
	}
	return false
}

// check the connection boot node periodically
func (pm *CsProtocolManager) bootVerifierConnCheck() {
	if pm.NodeConf.GetNodeType() != chain_config.NodeTypeOfVerifier && pm.NodeConf.GetNodeType() != chain_config.NodeTypeOfVerifierBoot {
		log.Info("cur node isn't verifier or verifier boot node, do not check v boot conn")
		return
	}

	tw := g_timer.SetPeriodAndRun(pm.checkAndHandleVerBootNodes, 8*time.Second)
	defer g_timer.StopWork(tw)

	<-pm.stop
}

//provide register communication service for external package
func (pm *CsProtocolManager) RegisterCommunicationService(cService CommunicationService, executable CommunicationExecutable) {
	pm.registerCommunicationService(cService, executable)
	return
}

//provide get current verifier peers function for external package
func (pm *CsProtocolManager) GetCurrentVerifierPeers() map[string]PmAbstractPeer {
	return pm.peerSetManager.currentVerifierPeers.GetPeers()
}

func (pm *CsProtocolManager) IsSync() bool {
	currentBlock := pm.Chain.CurrentBlock()

	if currentBlock == nil {
		return true
	}

	bestPeer := pm.BestPeer()
	// if best peer is nil, node is no sync
	if bestPeer == nil {
		return false
	}

	_, bestPeerHeight := bestPeer.GetHead()

	//log.Info("the currentBlock.Number is:","number",currentBlock.Number())
	//log.Info("the bestPeerHeight is:","bestPeerHeight",bestPeerHeight)
	// if peer current block number + 10 > best peer height , node is sync, but return false
	if currentBlock.Number()+10 >= bestPeerHeight {
		return false
	}

	return true
}

func (pm *CsProtocolManager) PrintPeerHealthCheck() {
	basePeers := pm.peerSetManager.basePeers.GetPeers()
	curPeers := pm.peerSetManager.currentVerifierPeers.GetPeers()
	nextPeers := pm.peerSetManager.nextVerifierPeers.GetPeers()
	vBootPeers := pm.peerSetManager.verifierBootNode.GetPeers()

	if log.OutputHealthLog() {
		printPeerInfo("base", basePeers)
		printPeerInfo("cur", curPeers)
		printPeerInfo("next", nextPeers)
		printPeerInfo("vboot", vBootPeers)
		log.Health.Debug("======================")
	}

	norLen := len(basePeers)
	curLen := len(curPeers)
	nextLen := len(nextPeers)
	vBootLen := len(vBootPeers)
	log.Health.Info("pm print cur peers info", "normal", norLen, "cur vers", curLen, "next vers", nextLen, "v boots", vBootLen)

	g_metrics.Set(g_metrics.NorPeerSetGauge, "", float64(norLen))
	g_metrics.Set(g_metrics.CurPeerSetGauge, "", float64(curLen))
	g_metrics.Set(g_metrics.NextPeerSetGauge, "", float64(nextLen))
	g_metrics.Set(g_metrics.VBootPeerSetGauge, "", float64(vBootLen))
}

func printPeerInfo(pSet string, ps map[string]PmAbstractPeer) {
	for _, p := range ps {
		log.Health.Debug("peer conn info", "node", p.NodeName(), "is running", p.IsRunning(), "remote addr", p.RemoteVerifierAddress(), "in set", pSet)
	}
}
