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
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/dipperin/dipperin-core/third_party/p2p/enode"
	"net"
)

//go:generate mockgen -destination=./peer_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication PmAbstractPeer
//go:generate mockgen -destination=../mine/mineworker/peer_mock_test.go -package=mineworker github.com/dipperin/dipperin-core/core/chaincommunication PmAbstractPeer
// is responsible for sending and receiving messages
type PmAbstractPeer interface {
	// add node name
	NodeName() string
	// remote node type
	NodeType() uint64

	SendMsg(msgCode uint64, msg interface{}) error
	// remote node id
	ID() string
	// read peer msg
	ReadMsg() (p2p.Msg, error)

	GetHead() (common.Hash, uint64)

	SetHead(head common.Hash, height uint64)

	GetPeerRawUrl() string

	DisconnectPeer()

	RemoteVerifierAddress() (addr common.Address)
	// remote host and port
	RemoteAddress() net.Addr

	SetRemoteVerifierAddress(addr common.Address)
	SetNodeName(name string)
	SetNodeType(nt uint64)
	SetPeerRawUrl(rawUrl string)

	SetNotRunning()
	IsRunning() bool

	GetCsPeerInfo() *p2p.CsPeerInfo
}

//go:generate mockgen -destination=./peer_set_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication AbstractPeerSet
type AbstractPeerSet interface {
	BestPeer() PmAbstractPeer

	GetPeers() map[string]PmAbstractPeer

	AddPeer(p PmAbstractPeer) error
	RemovePeer(id string) error
	ReplacePeers(newPeers map[string]PmAbstractPeer)

	Peer(id string) PmAbstractPeer
	Len() int
	Close()

	GetPeersInfo() []*p2p.CsPeerInfo
}

// manage peer
type AbstractProtocolManager interface {
	Start() error
	Stop()
	Protocols() []p2p.Protocol
}

//go:generate mockgen -destination=./communication_service_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication CommunicationService
// specific function
type CommunicationService interface {
	MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
}

//go:generate mockgen -destination=./communication_executable_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication CommunicationExecutable
type CommunicationExecutable interface {
	Start() error
	Stop()
}

//go:generate mockgen -destination=./tx_pool_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication TxPool
type TxPool interface {
	AddLocal(tx model.AbstractTransaction) error
	AddRemote(tx model.AbstractTransaction) error
	AddLocals(txs []model.AbstractTransaction) []error
	AddRemotes(txs []model.AbstractTransaction) []error
	ConvertPoolToMap() map[common.Hash]model.AbstractTransaction
	Stats() (int, int)
	GetTxsEstimator(broadcastBloom *iblt.Bloom) *iblt.HybridEstimator
	Pending() (map[common.Address][]model.AbstractTransaction, error)
	Queueing() (map[common.Address][]model.AbstractTransaction, error)
}

//go:generate mockgen -destination=./pbft_node_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication PbftNode
type PbftNode interface {
	OnNewWaitVerifyBlock(block model.AbstractBlock, id string)
	OnNewMsg(msg interface{}) error
	ChangePrimary(primary string)

	OnNewP2PMsg(msg p2p.Msg, p PmAbstractPeer) error
	AddPeer(p PmAbstractPeer) error

	OnEnterNewHeight(h uint64)
}

//go:generate mockgen -destination=./pbft_decoder_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication P2PMsgDecoder
type P2PMsgDecoder interface {
	DecodeTxMsg(msg p2p.Msg) (model.AbstractTransaction, error)
	DecoderBlockMsg(msg p2p.Msg) (model.AbstractBlock, error)
	DecodeTxsMsg(msg p2p.Msg) (result []model.AbstractTransaction, err error)
}

//go:generate mockgen -destination=./p2p_server_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication P2PServer
type P2PServer interface {
	AddPeer(node *enode.Node)
	RemovePeer(node *enode.Node)
	Self() *enode.Node
}

//go:generate mockgen -destination=./chain_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication Chain
type Chain interface {
	CurrentBlock() model.AbstractBlock
	GetSlot(block model.AbstractBlock) *uint64
	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetBlockByHash(hash common.Hash) model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	GetSeenCommit(height uint64) []model.AbstractVerification
	SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error
}

//go:generate mockgen -destination=./pbft_signer_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication PbftSigner
type PbftSigner interface {
	GetAddress() common.Address
	SetBaseAddress(address common.Address)
	SignHash(hash []byte) ([]byte, error)
	PublicKey() *ecdsa.PublicKey
	ValidSign(hash []byte, pubKey []byte, sign []byte) error
	Evaluate(account accountsbase.Account, seed []byte) (index [32]byte, proof []byte, err error)
}

//go:generate mockgen -destination=./verifiers_reader_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication VerifiersReader
type VerifiersReader interface {
	CurrentVerifiers() []common.Address
	NextVerifiers() []common.Address
	ShouldChangeVerifier() bool
}

//go:generate mockgen -destination=./peer_manager_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication PeerManager
type PeerManager interface {
	GetPeers() map[string]PmAbstractPeer
	BestPeer() PmAbstractPeer
	IsSync() bool
	GetPeer(id string) PmAbstractPeer
	RemovePeer(id string)
}

type AbstractPbftProtocolManager interface {
	PeerManager
	GetCurrentConnectPeers() map[string]common.Address
	GetVerifierBootNode() map[string]PmAbstractPeer
	GetNextVerifierPeers() map[string]PmAbstractPeer
	SelfIsBootNode() bool
	GetSelfNode() *enode.Node
	MatchCurrentVerifiersToNext()
}

type ChainDownloader interface {
	Start() error
	Stop()
	MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
	//SetEiFetcher(f *EiBlockFetcher)
}

//go:generate mockgen -destination=./transaction_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/model AbstractTransaction

//go:generate mockgen -destination=./msgReadWriter_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/third-party/p2p MsgReadWriter

//go:generate mockgen -destination=./p2pPeer_mock_test.go -package=chaincommunication github.com/dipperin/dipperin-core/core/chaincommunication P2PPeer
