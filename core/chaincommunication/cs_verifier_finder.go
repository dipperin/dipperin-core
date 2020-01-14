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
	"fmt"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/gevent"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/dipperin/dipperin-core/third_party/p2p/enode"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

var (
	fetchConnTimeout = 15 * time.Second
)

//go:generate mockgen -destination=./abs_peer_manager_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication AbsPeerManager
type AbsPeerManager interface {
	BestPeer() PmAbstractPeer
	SelfIsCurrentVerifier() bool
	SelfIsNextVerifier() bool
	HaveEnoughVerifiers(withOrganizeVSet bool) (missCur uint, missNext uint)
	GetVerifierBootNode() map[string]PmAbstractPeer
	GetSelfNode() *enode.Node
	ConnectPeer(node *enode.Node)
	CurrentVerifierPeersSet() AbstractPeerSet
	NextVerifierPeersSet() AbstractPeerSet
}

type GetVerifiersReq struct {
	ID uint64

	CurMiss  uint
	NextMiss uint
	Slot     uint64
}

type GetVerifiersResp struct {
	ReqID uint64

	Cur     []string
	Next    []string
	ErrInfo string
}

type fetchReq struct {
	ReqID    uint64
	RespChan chan *GetVerifiersResp
}

func NewVfFetcher() *vfFetcher {
	return &vfFetcher{
		reqs:       make(map[uint64]chan *GetVerifiersResp),
		addReqChan: make(chan fetchReq),
		respChan:   make(chan GetVerifiersResp),
	}
}

// fetch GetVerifiersResp
type vfFetcher struct {
	reqs map[uint64]chan *GetVerifiersResp

	addReqChan chan fetchReq
	respChan   chan GetVerifiersResp
}

func (f *vfFetcher) loop() {
	for {
		select {
		case req := <-f.addReqChan:
			if f.reqs[req.ReqID] != nil {
				log.DLogger.Error("dup get conn req", zap.Any("req", req))
			}
			f.reqs[req.ReqID] = req.RespChan
		case resp := <-f.respChan:
			if f.reqs[resp.ReqID] != nil {
				select {
				case f.reqs[resp.ReqID] <- &resp:
				case <-time.After(100 * time.Millisecond):
				}
				delete(f.reqs, resp.ReqID)
			}
		}
	}
}

func (f *vfFetcher) OnGetVerifiersResp(msg p2p.Msg, p PmAbstractPeer) error {
	var resp GetVerifiersResp
	if err := msg.Decode(&resp); err != nil {
		return err
	}

	select {
	case f.respChan <- resp:
	case <-time.After(100 * time.Millisecond):
		log.DLogger.Warn("can't write to fetcher.respChan, maybe fetcher not started")
	}
	return nil
}

func (f *vfFetcher) getVerifiersFromBoot(req GetVerifiersReq, peer PmAbstractPeer) (resp *GetVerifiersResp) {
	req.ID = uint64(time.Now().UnixNano())
	respChan := make(chan *GetVerifiersResp)

	if err := peer.SendMsg(GetVerifiersConnFromBootNode, req); err != nil {
		log.DLogger.Warn("send get v conn msg to v boot failed", zap.Error(err))
		return nil
	}

	f.addReqChan <- fetchReq{ReqID: req.ID, RespChan: respChan}

	select {
	case resp = <-respChan:
		return resp
	case <-time.After(fetchConnTimeout):
		log.DLogger.Warn("fetch v conn from v boot timeout")
		return nil
	}
}

func NewVFinder(chain Chain, peerManager AbsPeerManager, chainCfg chainconfig.ChainConfig) *VFinder {
	return &VFinder{
		chain:       chain,
		peerManager: peerManager,
		chainCfg:    chainCfg,
		fetcher:     NewVfFetcher(),
	}
}

type VFinder struct {
	chain       Chain
	peerManager AbsPeerManager
	chainCfg    chainconfig.ChainConfig
	fetcher     *vfFetcher

	findingVerifiers uint32
	started          uint32
}

func (vf *VFinder) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	return map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error{
		BootNodeVerifiersConn: vf.fetcher.OnGetVerifiersResp,
	}
}

func (vf *VFinder) Start() error {
	if !atomic.CompareAndSwapUint32(&vf.started, 0, 1) {
		return gerror.ErrAlreadyStarted
	}

	go vf.fetcher.loop()

	return nil
}

func (vf *VFinder) Stop() {}

// do find verifiers
func (vf *VFinder) findVerifiers() {
	if !atomic.CompareAndSwapUint32(&vf.findingVerifiers, 0, 1) {
		log.DLogger.Warn("call find verifiers, but last finding not finished")
		return
	}
	defer atomic.CompareAndSwapUint32(&vf.findingVerifiers, 1, 0)

	if err := vf.shouldFindVerifiers(); err != nil {
		log.DLogger.Info("don't find verifiers, because of " + err.Error())
		// if ErrNotCurrentOrNextVerifier, disable to connect v boot
		return
	}

	// Organize connection and check have enough conn
	missC, missN := vf.peerManager.HaveEnoughVerifiers(true)
	if missC <= 0 && missN <= 0 {
		log.DLogger.Debug("have enough verifiers, don't get from v boot")
		return
	}
	log.DLogger.Info("don't have enough verifiers", zap.Uint("missC", missC), zap.Uint("missN", missN))

	vf.getVerifiers(missC, missN)
}

func (vf *VFinder) shouldFindVerifiers() error {
	// common check
	if err := canFind(vf.peerManager, vf.chain); err != nil {
		return err
	}

	// is cur or next verifiers
	if !vf.peerManager.SelfIsCurrentVerifier() && !vf.peerManager.SelfIsNextVerifier() {
		return gerror.ErrNotCurrentOrNextVerifier
	}

	// check conn enough outside
	return nil
}

// get verifiers from v boot node until got enough peers or got from all v boot
func (vf *VFinder) getVerifiers(missCur uint, missNext uint) {
	slot := vf.chain.GetSlot(vf.chain.CurrentBlock())
	if slot == nil {
		log.DLogger.Error("can't get slot for current block", zap.Uint64("cur b", vf.chain.CurrentBlock().Number()))
		panic("can't get slot for current block")
	}
	req := GetVerifiersReq{
		CurMiss:  missCur,
		NextMiss: missNext,
		Slot:     *slot,
	}

	bPeers := vf.peerManager.GetVerifierBootNode()
	for _, vb := range bPeers {
		if ok := vf.getVerifiersFromBoot(req, vb); ok {
			break
		}
	}
}

// get conn info, return true if got enough conn
func (vf *VFinder) getVerifiersFromBoot(req GetVerifiersReq, peer PmAbstractPeer) (ok bool) {
	resp := vf.fetcher.getVerifiersFromBoot(req, peer)
	if resp == nil {
		return
	}

	if resp.ErrInfo != "" {
		log.DLogger.Warn("err from v boot", zap.String("info", resp.ErrInfo))
		return
	}

	selfID := vf.peerManager.GetSelfNode().ID().String()

	// connect cur
	missC := int(req.CurMiss)
	curSet := vf.peerManager.CurrentVerifierPeersSet()
	for _, vStr := range resp.Cur {
		if vf.checkAndConnectNode(selfID, vStr, curSet) {
			missC--
		}
	}

	// connect next
	missN := int(req.NextMiss)
	nextSet := vf.peerManager.NextVerifierPeersSet()
	for _, vStr := range resp.Next {
		if vf.checkAndConnectNode(selfID, vStr, nextSet) {
			missN--
		}
	}

	return missC <= 0 && missN <= 0
}

func (vf *VFinder) checkAndConnectNode(selfID string, connStr string, ps AbstractPeerSet) (ok bool) {
	node, err := enode.ParseV4(connStr)
	if err != nil {
		panic("parse node conn from v boot failed: " + err.Error())
	}

	nodeID := node.ID().String()
	if selfID == nodeID {
		return
	}

	if ps.Peer(nodeID) != nil {
		return
	}

	vf.peerManager.ConnectPeer(node)

	return true
}

func NewVFinderBoot(peerManager AbsPeerManager, chain Chain) *VFinderBoot {
	return &VFinderBoot{
		peerManager: peerManager,
		chain:       chain,
	}
}

type VFinderBoot struct {
	peerManager AbsPeerManager
	chain       Chain
}

func (vfb *VFinderBoot) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	return map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error{
		GetVerifiersConnFromBootNode: vfb.OnGetVerifiersReq,
	}
}

func (vfb *VFinderBoot) onGetVerifiersReq(req *GetVerifiersReq, from PmAbstractPeer) *GetVerifiersResp {
	resp := &GetVerifiersResp{ReqID: req.ID}

	// valid find
	if err := canFind(vfb.peerManager, vfb.chain); err != nil {
		resp.ErrInfo = err.Error()
		return resp
	}

	slot := vfb.chain.GetSlot(vfb.chain.CurrentBlock())
	if slot == nil {
		log.DLogger.Error("can't get slot for cur block", zap.Uint64("cur h", vfb.chain.CurrentBlock().Number()))
		panic("can't get slot for cur block")
	}
	if *slot != req.Slot {
		resp.ErrInfo = fmt.Sprintf("slot not match, req.slot: %v, boot.slot: %v", req.Slot, *slot)
		log.DLogger.Warn("slot not match", zap.Uint64("req.slot", req.Slot), zap.Uint64("boot.slot", *slot), zap.String("from", from.NodeName()))
		return resp
	}

	fID := from.ID()
	if req.CurMiss > 0 {
		curPeers := vfb.peerManager.CurrentVerifierPeersSet().GetPeers()
		for id, p := range curPeers {
			if id != fID {
				resp.Cur = append(resp.Cur, p.GetPeerRawUrl())
			}
		}
	}
	if req.NextMiss > 0 {
		nextPeers := vfb.peerManager.NextVerifierPeersSet().GetPeers()
		for id, p := range nextPeers {
			if id != fID {
				resp.Next = append(resp.Next, p.GetPeerRawUrl())
			}
		}
	}
	return resp
}

func (vfb *VFinderBoot) OnGetVerifiersReq(msg p2p.Msg, p PmAbstractPeer) error {
	var req GetVerifiersReq
	if err := msg.Decode(&req); err != nil {
		return err
	}

	resp := vfb.onGetVerifiersReq(&req, p)

	if err := p.SendMsg(BootNodeVerifiersConn, resp); err != nil {
		log.DLogger.Warn("send conn to verifier failed", zap.Error(err), zap.String("remote node", p.NodeName()))
	}
	return nil
}

// check can find for both verifier and v boot
func canFind(pm AbsPeerManager, chain Chain) error {
	bestPeer := pm.BestPeer()
	if bestPeer == nil {
		return gerror.ErrNoBestPeerFound
	}
	_, rHeight := bestPeer.GetHead()
	curB := chain.CurrentBlock()

	// valid height
	if curB.Number()+2 < rHeight {
		log.DLogger.Info("height too low, do not find verifiers", zap.Uint64("cur h", curB.Number()), zap.Uint64("remote h", rHeight))
		return gerror.ErrCurHeightTooLow
	}

	// is change point, do not find
	if chain.IsChangePoint(curB, false) {
		return gerror.ErrIsChangePointDoNotFind
	}

	return nil
}

// listen to new blocks
func (pm *CsProtocolManager) handleInsertEventForBft() error {
	// if it is a normal type go back directly
	if pm.selfPmType() == base {
		return nil
	}

	go func() {
		newBlockChan := make(chan model.Block, 0)
		//sub := pm.nodeContext.Chain().SubscribeBlockEvent(newBlockChan)
		sub := gevent.Subscribe(gevent.NewBlockInsertEvent, newBlockChan)
		defer sub.Unsubscribe()

		for {
			select {
			case newBlock := <-newBlockChan:

				log.DLogger.Debug("[Insert Event]", zap.Uint64("blockNumber", newBlock.Number()), zap.Bool("is change point", pm.Chain.IsChangePoint(&newBlock, false)), zap.Bool("is current", pm.SelfIsCurrentVerifier()), zap.Bool("is next", pm.SelfIsNextVerifier()))

				slot := pm.Chain.GetSlot(&newBlock)
				log.DLogger.Debug("the current slot is:", zap.Uint64("slot", *slot))

				// should change verifier
				if pm.Chain.IsChangePoint(&newBlock, false) && *slot > 0 {
					log.DLogger.Debug("[Insert Event] IsChangePoint CallChangeVerifier")
					pm.ChangeVerifiers()
				} else {
					log.DLogger.Debug("[Insert Event] NotChangePoint NotCurrentVerifier")

					if pm.vf != nil {
						pm.vf.findVerifiers()
					}
				}
				// goes to bft only the current round verifier
				if pm.SelfIsCurrentVerifier() {
					log.DLogger.Debug("[Insert Event] NotChangePoint IsCurrentVerifier")
					pm.PbftNode.OnEnterNewHeight(newBlock.Number() + 1)
				}

			case <-pm.stop:
				return
			}
		}
	}()

	return nil
}
