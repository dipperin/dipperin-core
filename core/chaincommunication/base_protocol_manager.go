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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	crypto2 "github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"go.uber.org/zap"
	"math/big"
	"sync"
)

//go:generate mockgen -destination=./node_conf_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication NodeConf
type NodeConf interface {
	GetNodeType() int
	GetNodeName() string
}

type BaseProtocolManager struct {
	protocols []p2p.Protocol
	// dispatch msg to communicationServices
	msgHandlers map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
	// start when call pm start
	executables []CommunicationExecutable

	wg sync.WaitGroup
}

// cService executable will be same object
func (pm *BaseProtocolManager) registerCommunicationService(cService CommunicationService, executable CommunicationExecutable) {
	if cService != nil {
		newHandlers := cService.MsgHandlers()
		for k, v := range newHandlers {
			if pm.msgHandlers[k] != nil {
				panic(fmt.Sprintf("already have handler:%v", k))
			}
			pm.msgHandlers[k] = v
		}
	}

	if executable != nil {
		pm.executables = append(pm.executables, executable)
	}
}

func (pm *BaseProtocolManager) RemovePeer(id string) { panic("impl me") }

// handle msg for GetPeers,
func (pm *BaseProtocolManager) handleMsg(p PmAbstractPeer) error {

	log.DLogger.Info("base protocol handle msg", zap.String("remote node", p.NodeName()))

	msg, err := p.ReadMsg()

	if err != nil {
		log.DLogger.Info("base protocol read msg from peer failed", zap.Error(err), zap.String("peer name", p.NodeName()))
		return err
	}

	defer msg.Discard()
	if msg.Size > ProtocolMaxMsgSize {
		return msgTooLargeErr
	}

	// find handler for this msg
	tmpHandler := pm.msgHandlers[uint64(msg.Code)]
	if tmpHandler == nil {
		log.DLogger.Error("Get message processing error", zap.Uint64("msg code", msg.Code))
		return msgHandleFuncNotFoundErr
	}

	// handle this msg
	if err = tmpHandler(msg, p); err != nil {
		p.SetNotRunning()
		return err
	}

	return nil
}

func (pm *BaseProtocolManager) Start() error {
	for _, r := range pm.executables {
		if err := r.Start(); err != nil {
			pm.Stop()
			return err
		}
	}
	return nil
}

func (pm *BaseProtocolManager) Stop() {
	for _, r := range pm.executables {
		r.Stop()
	}
}

func (pm *BaseProtocolManager) validStatus(status StatusData) error {
	return nil
}

type HandShakeData struct {
	ProtocolVersion uint32
	ChainID         *big.Int
	NetworkId       uint64
	//TD              *big.Int
	CurrentBlock       common.Hash
	CurrentBlockHeight uint64
	GenesisBlock       common.Hash
	NodeType           uint64
	NodeName           string
	// for pbft
	RawUrl string
}

// for hand shake
type StatusData struct {
	HandShakeData

	PubKey []byte
	Sign   []byte
}

func (status *StatusData) Sender() (result common.Address) {
	if err := validSign(status.DataHash().Bytes(), status.PubKey, status.Sign); err != nil {
		log.DLogger.Error("verifier hand shake verify signature failed", zap.Error(err))
		return
	}
	pubKey, err := crypto2.DecompressPubkey(status.PubKey)
	if err != nil {
		log.DLogger.Error("can't decode pub key from status data", zap.Error(err))
		return
	}
	// pass check sign, then get address from pubkey
	result = cs_crypto.GetNormalAddress(*pubKey)
	return
}

func validSign(hash []byte, pubKey []byte, sign []byte) error {
	if len(sign) == 0 {
		return errors.New("empty sign")
	}
	if !crypto2.VerifySignature(pubKey, hash, sign[:len(sign)-1]) {
		return errors.New("verify signature fail")
	}
	return nil
}

func (status *StatusData) DataHash() common.Hash {
	v := common.RlpHashKeccak256(status.HandShakeData)
	return v
}
