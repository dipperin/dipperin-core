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

package dipperin

import (
	"github.com/dipperin/dipperin-core/cmd/utils/debug"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/gmetrics"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/accounts/softwallet"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/cachedb"
	"github.com/dipperin/dipperin-core/core/chaincommunication"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	"github.com/dipperin/dipperin-core/core/csbft/csbftnode"
	"github.com/dipperin/dipperin-core/core/csbft/statemachine"
	"github.com/dipperin/dipperin-core/core/cschain"
	"github.com/dipperin/dipperin-core/core/cschain/chainstate"
	"github.com/dipperin/dipperin-core/core/cschain/chainwriter"
	"github.com/dipperin/dipperin-core/core/cschain/chainwriter/middleware"
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/core/mine/blockbuilder"
	"github.com/dipperin/dipperin-core/core/mine/minemaster"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/rpcinterface"
	"github.com/dipperin/dipperin-core/core/txpool"
	"github.com/dipperin/dipperin-core/core/verifiershaltcheck"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/dipperin/dipperin-core/third_party/p2p/nat"
	"github.com/dipperin/dipperin-core/third_party/p2p/netutil"
	"github.com/dipperin/dipperin-core/third_party/rpc"
	"github.com/dipperin/dipperin-core/third_party/vm-log-search"
	"go.uber.org/zap"
	"path/filepath"
	"strings"
	"sync/atomic"
)

type BlockValidator interface {
	Valid(block model.AbstractBlock) error
	FullValid(block model.AbstractBlock) error
}

type BaseComponent struct {
	nodeConfig           NodeConfig
	chainConfig          *chainconfig.ChainConfig
	DipperinConfig       *service.DipperinConfig
	csChainServiceConfig *cschain.CsChainServiceConfig
	bftConfig            *statemachine.BftConfig
	pmConf               *chaincommunication.CsProtocolManagerConfig
	txBConf              *chaincommunication.NewTxBroadcasterConfig
	verHaltCheckConfig   *verifiershaltcheck.HaltCheckConf

	prometheusServer *gmetrics.PrometheusMetricsServer
	//cacheDB                     *cachedb.CacheDB
	fullChain                   *cschain.CsChainService
	txPool                      *txpool.TxPool
	rpcService                  *rpcinterface.Service
	txSigner                    model.Signer
	defaultPriorityCalculator   model.PriofityCalculator
	coinbaseAddr                *atomic.Value
	blockDecoder                model.BlockDecoder
	consensusBeforeInsertBlocks BlockValidator
	defaultMsgDecoder           chaincommunication.P2PMsgDecoder
	verifiersReader             VerifiersReader
	chainService                *service.VenusFullChainService
	walletManager               *accounts.WalletManager
	msgSigner                   *accounts.WalletSigner
	bftNode                     *csbftnode.CsBft
	p2pServer                   *p2p.Server
	broadcastDelegate           *chaincommunication.BroadcastDelegate
	csPm                        *chaincommunication.CsProtocolManager
	minePm                      *chaincommunication.MineProtocolManager
	mineMaster                  minemaster.Master
	mineMasterServer            minemaster.MasterServer
	defaultAccountAddress       common.Address
	verHaltCheck                *verifiershaltcheck.SystemHaltedCheck
}

func NewBftNode(nodeConfig NodeConfig) (n Node) {
	// newBaseComponent
	baseComponent := newBaseComponent(nodeConfig)
	// init full chain
	baseComponent.initFullChain()
	// init tx pool
	baseComponent.initTxPool()
	// init Chain service
	baseComponent.initChainService()
	// init wallet manager
	baseComponent.initWalletManager()
	// init msg signer
	//baseComponent.initMsgSigner()
	// init bft node
	baseComponent.initBft()
	// init p2p service
	baseComponent.initP2PService()

	// set p2p info to bft
	baseComponent.setBftAfterP2PInit()
	// init rpc
	baseComponent.initRpc()
	// init mine master
	baseComponent.initMineMaster()
	// setup service config
	baseComponent.buildDipperinConfig()
	//init verifier halt check
	baseComponent.initVerHaltCheck()

	// wrap p2p protocols
	baseComponent.addP2PProtocols()
	// wrap all node service
	nodeServices := baseComponent.getNodeServices()

	// add extra services and rpc apis
	if baseComponent.nodeConfig.ExtraServiceFunc != nil {
		eApis, eServices := baseComponent.nodeConfig.ExtraServiceFunc(ExtraServiceFuncConfig{
			DipperinConfig: *baseComponent.DipperinConfig,
			ChainService:   baseComponent.chainService,
		})
		nodeServices = append(nodeServices, eServices...)
		baseComponent.rpcService.AddApis(eApis)
	}

	n = NewCsNode(NodeConfig{}, baseComponent, nodeServices)
	baseComponent.DipperinConfig.Node = n

	return
}

// newBaseComponent configs and base components
func newBaseComponent(nodeConfig NodeConfig) *BaseComponent {
	promeS := gmetrics.NewPrometheusMetricsServer(nodeConfig.GetPMetricsPort())
	gmetrics.InitCSMetrics()
	b := &BaseComponent{
		prometheusServer:          promeS,
		chainConfig:               chainconfig.GetChainConfig(),
		DipperinConfig:            &service.DipperinConfig{},
		csChainServiceConfig:      &cschain.CsChainServiceConfig{},
		defaultPriorityCalculator: model.DefaultPriorityCalculator,
		defaultMsgDecoder:         chaincommunication.MakeDefaultMsgDecoder(),
		coinbaseAddr:              &atomic.Value{},
		nodeConfig:                nodeConfig,
	}
	b.txSigner = model.NewSigner(b.chainConfig.ChainId)

	// init block decoder
	b.blockDecoder = model.MakeDefaultBlockDecoder()

	// load boot nodes from datadir file
	chainconfig.InitBootNodes(nodeConfig.DataDir)

	// init data decoder
	model.SetBlockRlpHandler(&model.PBFTBlockRlpHandler{})
	model.SetBlockJsonHandler(&model.PBFTBlockJsonHandler{})

	return b
}

func (b *BaseComponent) setNodeSignerInfo() error {
	account, err := b.walletManager.GetMainAccount()
	if err != nil {
		return err
	}

	b.coinbaseAddr.Store(account.Address)
	b.defaultAccountAddress = account.Address
	b.msgSigner = accounts.MakeWalletSigner(b.defaultAccountAddress, b.walletManager)

	b.DipperinConfig.DefaultAccount = b.defaultAccountAddress
	b.DipperinConfig.MsgSigner = b.msgSigner

	//protocol manager
	b.pmConf.MsgSigner = b.msgSigner
	b.csPm.MsgSigner = b.pmConf.MsgSigner

	//bft and verifier halt check
	if b.nodeConfig.NodeType == chainconfig.NodeTypeOfVerifier || b.nodeConfig.NodeType == chainconfig.NodeTypeOfVerifierBoot {
		b.bftConfig.Signer = b.msgSigner
		b.verHaltCheckConfig.WalletSigner = b.msgSigner
		b.verHaltCheck.SetMsgSigner(b.msgSigner)
	}

	//mineMaster
	if b.nodeConfig.NodeType == chainconfig.NodeTypeOfMineMaster {
		b.mineMaster.SetMsgSigner(b.msgSigner)
		b.mineMaster.SetCoinbaseAddress(b.defaultAccountAddress)
	}

	return nil
}

func (b *BaseComponent) buildDipperinConfig() {
	b.DipperinConfig.PbftPm = b.csPm
	b.DipperinConfig.Broadcaster = b.broadcastDelegate
	b.DipperinConfig.ChainReader = b.fullChain
	b.DipperinConfig.TxPool = b.txPool
	b.DipperinConfig.NodeConf = b.nodeConfig
	b.DipperinConfig.ChainConfig = *b.chainConfig
	b.DipperinConfig.PriorityCalculator = b.defaultPriorityCalculator
	b.DipperinConfig.P2PServer = b.p2pServer
	b.DipperinConfig.NormalPm = b.csPm
	b.DipperinConfig.NormalPm = b.csPm

	b.DipperinConfig.WalletManager = b.walletManager
	b.DipperinConfig.MineMaster = b.mineMaster
	b.DipperinConfig.MineMasterServer = b.mineMasterServer
	//b.DipperinConfig.DefaultAccount = b.defaultAccountAddress
	//b.DipperinConfig.MsgSigner = b.msgSigner
	b.DipperinConfig.ChainIndex = vm_log_search.NewBloomIndexer(b.DipperinConfig.ChainReader, b.fullChain.CacheChainState.ChainState.GetDB(), vm_log_search.BloomBitsBlocks, vm_log_search.BloomConfirms)
}

func (b *BaseComponent) buildBftConfig() {
	b.bftConfig = &statemachine.BftConfig{
		//FetcherConnAdaptCsBft:csPm,
		ChainReader: b.fullChain,
		//Fetcher:components.NewFetcher(csPm),
		Signer: b.msgSigner,
		//Sender:MsgSender,
		Validator: b.consensusBeforeInsertBlocks,
	}
}

func (b *BaseComponent) buildHaltCheckConfig() {
	b.verHaltCheckConfig = &verifiershaltcheck.HaltCheckConf{
		NodeType:        b.nodeConfig.NodeType,
		CsProtocol:      b.csPm,
		NeedChainReader: b.fullChain,
		WalletSigner:    b.msgSigner,
		Broadcast:       b.broadcastDelegate.BroadcastEiBlock,
		EconomyModel:    b.fullChain.EconomyModel,
	}
}

func (b *BaseComponent) buildCommunicationConfig() {
	b.pmConf = &chaincommunication.CsProtocolManagerConfig{
		ChainConfig:     *b.chainConfig,
		Chain:           b.fullChain,
		P2PServer:       b.p2pServer,
		NodeConf:        b.nodeConfig,
		VerifiersReader: b.verifiersReader,
		PbftNode:        b.bftNode,
		MsgSigner:       b.msgSigner,
	}
	b.txBConf = &chaincommunication.NewTxBroadcasterConfig{
		P2PMsgDecoder: b.defaultMsgDecoder,
		TxPool:        b.txPool,
		NodeConf:      b.nodeConfig,
	}
}

func (b *BaseComponent) buildMineConfig(modelConfig blockbuilder.ModelConfig) minemaster.MineConfig {
	return minemaster.MineConfig{
		GasFloor:         &atomic.Value{},
		GasCeil:          &atomic.Value{},
		CoinbaseAddress:  b.coinbaseAddr,
		BlockBuilder:     blockbuilder.MakeBftBlockBuilder(modelConfig),
		BlockBroadcaster: b.broadcastDelegate,
	}
}

func (b *BaseComponent) builderModelConfig() blockbuilder.ModelConfig {
	return blockbuilder.ModelConfig{
		ChainReader:        b.fullChain,
		TxPool:             b.txPool,
		PriorityCalculator: b.defaultPriorityCalculator,
		MsgSigner:          b.msgSigner,
		ChainConfig:        *b.chainConfig,
	}
}

func (b *BaseComponent) initFullChain() {
	// init full chain
	b.fullChain = cschain.NewCsChainService(b.csChainServiceConfig, chainstate.NewChainState(&chainstate.ChainStateConfig{
		ChainConfig:   b.chainConfig,
		DataDir:       b.nodeConfig.DataDir,
		WriterFactory: chainwriter.NewChainWriterFactory(),
	}))
	b.csChainServiceConfig.CacheDB = cachedb.NewCacheDB(b.fullChain.GetDB())
	cachedb.SetCacheDataDecoder(&cachedb.BFTCacheDataDecoder{})

	b.verifiersReader = chain.MakeVerifiersReader(b.fullChain)
	b.consensusBeforeInsertBlocks = middleware.NewBftBlockValidator(b.fullChain)

	// Add Venus Testnet
	if chainconfig.GetCurBootsEnv() != chainconfig.BootEnvMercury && chainconfig.GetCurBootsEnv() != chainconfig.BootEnvVenus {
		debug.Memsize.Add("fullChain", b.fullChain)
		// TODo confirm if you need
		//debug.Memsize.Add("consensusBeforeInsertBlocks", consensusBeforeInsertBlocks)
		//debug.Memsize.Add("txValidator", c.txValidator)
		//debug.Memsize.Add("cacheDB", cacheDB)
		//debug.Memsize.Add("verifiersReader", verifiersReader)
	}
}

func (b *BaseComponent) initTxPool() {
	txPoolConfig := txpool.DefaultTxPoolConfig
	txPoolConfig.Journal = filepath.Join(b.nodeConfig.DataDir, "transaction.rlp")
	// no need to replace with context
	b.txPool = txpool.NewTxPool(txPoolConfig, *b.chainConfig, b.fullChain)
	b.csChainServiceConfig.TxPool = b.txPool

	if chainconfig.GetCurBootsEnv() != chainconfig.BootEnvMercury {
		debug.Memsize.Add("tx pool", b.txPool)
	}
}

func (b *BaseComponent) initChainService() {
	b.DipperinConfig.ChainReader = b.fullChain
	//b.DipperinConfig.ChainIndex =
	// init service
	b.chainService = service.MakeFullChainService(b.DipperinConfig)
}

func (b *BaseComponent) initWalletManager() {
	var err error
	if b.nodeConfig.NoWalletStart {
		if b.walletManager, err = accounts.NewWalletManager(b.chainService); err != nil {
			panic("init wallet manager failed: " + err.Error())
		}
		return
	}

	// load wallet manager
	defaultWallet, wErr := softwallet.NewSoftWallet()
	if wErr != nil {
		panic("new soft wallet failed: " + wErr.Error())
	}

	var mnemonic string
	exit, _ := softwallet.PathExists(b.nodeConfig.SoftWalletFile())
	log.DLogger.Info("initWalletManager the wallet exit", zap.Bool("exit", exit))
	if exit {
		err = defaultWallet.Open(b.nodeConfig.SoftWalletFile(), b.nodeConfig.SoftWalletName(), b.nodeConfig.SoftWalletPassword)
	} else {
		mnemonic, err = defaultWallet.Establish(b.nodeConfig.SoftWalletFile(), b.nodeConfig.SoftWalletName(), b.nodeConfig.SoftWalletPassword, b.nodeConfig.SoftWalletPassPhrase)
		mnemonic = strings.Replace(mnemonic, " ", ",", -1)
		log.DLogger.Info("EstablishWallet mnemonic is:", zap.String("mnemonic", mnemonic))
	}

	if err != nil {
		log.DLogger.Info("open or establish wallet error ", zap.Error(err))
		panic("initWalletManager open or establish wallet error: " + err.Error())
	}
	if b.walletManager, err = accounts.NewWalletManager(b.chainService, defaultWallet); err != nil {
		log.DLogger.Info("init wallet manager failed:", zap.Any("walletManager", b.walletManager), zap.Error(err))
		panic("init wallet manager failed: " + err.Error())
	}

	log.DLogger.Info("the wallet number is:", zap.Int("number", len(b.walletManager.Wallets)))
	return
}

/*func (b *BaseComponent) initWalletManager() {
	tmpLog := log.New()
	tmpLog.SetHandler(log.StdoutHandler)
	var err error
	log.DLogger.Info("the nodeType is:", "nodeType", b.nodeConfig.NodeType)
	// No need to create or open a default wallet when the normal node starts
	if b.nodeConfig.NodeType == chain_config.NodeTypeOfNormal {
		if b.walletManager, err = accounts.NewWalletManager(b.chainService); err != nil {
			panic("init wallet manager failed: " + err.Error())
		}
		return
	}
	// load wallet manager
	defaultWallet, wErr := soft_wallet.NewSoftWallet()
	if wErr != nil {
		panic("new soft wallet failed: " + wErr.Error())
	}

	var mnemonic string
	exit, _ := soft_wallet.PathExists(b.nodeConfig.SoftWalletFile())
	if exit {
		err = defaultWallet.Open(b.nodeConfig.SoftWalletFile(), b.nodeConfig.SoftWalletName(), b.nodeConfig.SoftWalletPassword)
	} else {
		mnemonic, err = defaultWallet.Establish(b.nodeConfig.SoftWalletFile(), b.nodeConfig.SoftWalletName(), b.nodeConfig.SoftWalletPassword, b.nodeConfig.SoftWalletPassPhrase)
		mnemonic = strings.Replace(mnemonic, " ", ",", -1)
		tmpLog.Info("EstablishWallet mnemonic is:", "mnemonic", mnemonic)
	}

	if err != nil {
		tmpLog.Info("open or establish wallet error ", zap.Error(err))
		os.Exit(1)
	}

	if b.walletManager, err = accounts.NewWalletManager(b.chainService, defaultWallet); err != nil {
		tmpLog.Info("init wallet manager failed: ", zap.Error(err))
		log.DLogger.Info("init wallet manager failed:", "walletManager", b.walletManager, zap.Error(err))
		os.Exit(1)
	}
	var defaultAccounts []accounts.Account
	if defaultAccounts, err = defaultWallet.Accounts(); err != nil {
		tmpLog.Info("get default accounts failed: ", zap.Error(err))
		os.Exit(1)
	}
	b.coinbaseAddr.Store(defaultAccounts[0].Address)
	b.defaultAccountAddress = defaultAccounts[0].Address
	log.DLogger.Info("open wallet success", "b.defaultAccountAddress", b.defaultAccountAddress)
}*/

func (b *BaseComponent) initP2PService() {
	// load p2p
	p2pConf := DefaultP2PConf()

	if b.nodeConfig.NoDiscovery == 1 {
		p2pConf.NoDiscovery = true
	}
	// boot wait for connect, do not connect others by it self
	if b.nodeConfig.NodeType == chainconfig.NodeTypeOfVerifierBoot {
		p2pConf.NoDiscovery = true
	}

	// init nat
	if b.nodeConfig.Nat != "" {
		p2pNat, err := nat.Parse(b.nodeConfig.Nat)

		if err != nil {
			panic(err)
		}

		p2pConf.NAT = p2pNat
	}

	if chainconfig.GetCurBootsEnv() == chainconfig.BootEnvTest {
		restrictList, err := netutil.ParseNetlist(chainconfig.TestIPWhiteList)
		if err != nil {
			panic(err)
		}
		p2pConf.NetRestrict = restrictList
	}
	p2pConf.ListenAddr = b.nodeConfig.P2PListener
	p2pConf.BootstrapNodes = chainconfig.KBucketNodes
	p2pConf.PrivateKey = loadNodeKeyFromFile(b.nodeConfig.DataDir)
	p2pConf.StaticNodes = getNodeList(filepath.Join(b.nodeConfig.DataDir, staticNodes))
	p2pConf.TrustedNodes = getNodeList(filepath.Join(b.nodeConfig.DataDir, trustedNodes))

	p2pServer := &p2p.Server{Config: p2pConf}
	b.p2pServer = p2pServer

	b.buildCommunicationConfig()
	csPm, broadcastDelegate := chaincommunication.MakeCsProtocolManager(b.pmConf, b.txBConf)

	b.csPm = csPm
	b.broadcastDelegate = broadcastDelegate

	if chainconfig.GetCurBootsEnv() != chainconfig.BootEnvMercury && chainconfig.GetCurBootsEnv() != chainconfig.BootEnvVenus {
		debug.Memsize.Add("p2p server", p2pServer)
	}
}

func (b *BaseComponent) initRpc() {
	// load rpc service todo chainService not init
	rpcApi := rpcinterface.MakeDipperinVenusApi(b.chainService)
	debugApi := rpcinterface.MakeDipperinDebugApi(b.chainService)
	p2pApi := rpcinterface.MakeDipperinP2PApi(b.chainService)
	externalApi := rpcinterface.MakeDipperExternalApi(rpcApi)

	b.rpcService = rpcinterface.MakeRpcService(b.nodeConfig, []rpc.API{
		{
			Namespace: "dipperin",
			Version:   chainconfig.Version,
			Service:   rpcApi,
			Public:    false,
		},
		{
			Namespace: "dipperin",
			Version:   chainconfig.Version,
			Service:   externalApi,
			Public:    true,
		},
		{
			Namespace: "debug",
			Version:   "1.0",
			Service:   debug.Handler,
			Public:    true,
		},
		{
			Namespace: "debug",
			Version:   "1.0",
			Service:   debugApi,
			Public:    true,
		},
		{
			Namespace: "p2p",
			Version:   "1.0",
			Service:   p2pApi,
			Public:    false,
		},
	}, b.nodeConfig.GetAllowHosts())

	if chainconfig.GetCurBootsEnv() != chainconfig.BootEnvMercury {
		debug.Memsize.Add("rpc server", b.rpcService)
	}

	return
}

func (b *BaseComponent) initMineMaster() {
	if b.nodeConfig.NodeType != chainconfig.NodeTypeOfMineMaster {
		return
	}

	mineConfig := b.buildMineConfig(b.builderModelConfig())
	mineConfig.GasFloor.Store(uint64(chainconfig.BlockGasLimit))
	mineConfig.GasCeil.Store(uint64(chainconfig.BlockGasLimit))
	// chain service not init here
	mineMaster, mineMasterServer := minemaster.MakeMineMaster(mineConfig)
	minePm := chaincommunication.NewMineProtocolManager(mineMasterServer)

	// p2p server not init
	b.minePm = minePm
	b.mineMaster = mineMaster
	b.mineMasterServer = mineMasterServer
}

// must have init wallet manager
func (b *BaseComponent) initMsgSigner() {
	if b.nodeConfig.NodeType == chainconfig.NodeTypeOfNormal {
		b.msgSigner = nil
	} else {
		log.DLogger.Info("setup default sign address", zap.String("addr", b.defaultAccountAddress.Hex()))
		b.msgSigner = accounts.MakeWalletSigner(b.defaultAccountAddress, b.walletManager)
	}
}

func (b *BaseComponent) initBft() {
	if b.nodeConfig.NodeType != chainconfig.NodeTypeOfVerifier && b.nodeConfig.NodeType != chainconfig.NodeTypeOfVerifierBoot {
		return
	}
	b.buildBftConfig()
	b.bftNode = csbftnode.NewCsBft(b.bftConfig)
}

func (b *BaseComponent) setBftAfterP2PInit() {
	if b.nodeConfig.NodeType == chainconfig.NodeTypeOfVerifier || b.nodeConfig.NodeType == chainconfig.NodeTypeOfVerifierBoot {
		msgSender := MsgSender{
			csPm:              b.csPm,
			broadcastDelegate: b.broadcastDelegate,
		}
		fetcher := components.NewFetcher(b.csPm)
		b.bftNode.SetFetcher(fetcher)
		b.bftNode.Sender = &msgSender
		b.bftNode.Fetcher = fetcher
	}
}

func (b *BaseComponent) initVerHaltCheck() {
	if b.nodeConfig.NodeType != chainconfig.NodeTypeOfVerifier && b.nodeConfig.NodeType != chainconfig.NodeTypeOfVerifierBoot {
		return
	}
	b.buildHaltCheckConfig()
	b.verHaltCheck = verifiershaltcheck.MakeSystemHaltedCheck(b.verHaltCheckConfig)
	b.csPm.RegisterCommunicationService(b.verHaltCheck, b.verHaltCheck)
}

func (b *BaseComponent) getNodeServices() []NodeService {
	// these services may have nil
	return filterNilService([]NodeService{
		b.chainService, b.bftNode, b.walletManager, b.csPm,
		b.p2pServer, b.rpcService, b.txPool, b.prometheusServer, b.DipperinConfig.ChainIndex,
	})
}

func filterNilService(ns []NodeService) (result []NodeService) {
	for _, s := range ns {
		if !util.InterfaceIsNil(s) {
			result = append(result, s)
		}
	}
	return
}

func (b *BaseComponent) addP2PProtocols() {
	if b.csPm != nil {
		b.p2pServer.Protocols = append(b.p2pServer.Protocols, b.csPm.Protocols()...)
	}
	if b.minePm != nil {
		b.p2pServer.Protocols = append(b.p2pServer.Protocols, b.minePm.GetProtocol())
	}
}

type MsgSender struct {
	broadcastDelegate *chaincommunication.BroadcastDelegate
	csPm              *chaincommunication.CsProtocolManager
}

func (m *MsgSender) BroadcastMsg(msgCode uint64, msg interface{}) {
	m.csPm.BroadcastMsg(msgCode, msg)

}

func (m *MsgSender) SendReqRoundMsg(msgCode uint64, from []common.Address, msg interface{}) {
	m.csPm.BroadcastMsgToTargetVerifiers(msgCode, from, msg)
}

func (m *MsgSender) BroadcastEiBlock(block model.AbstractBlock) {
	m.broadcastDelegate.BroadcastEiBlock(block)
}
