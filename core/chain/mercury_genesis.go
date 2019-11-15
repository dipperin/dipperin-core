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

package chain

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

//the verifier address to generate the first twenty blocks
var VerifierAddress []common.Address

//delete angle verifier csWallet cipher in the Dipperin-core source code
func init() {
	env := os.Getenv("boots_env")
	log.Info("start env: " + env)
	if env == "mercury" {
		VerifierAddress = chain_config.MercuryVerifierAddress
		log.Debug("default mercury verifier", "count", len(VerifierAddress))
	} else {
		VerifierAddress = chain_config.LocalVerifierAddress
	}

	chain_config.VerBootNodeAddress = chain_config.VerifierBootNodeAddress
}

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config    *chain_config.ChainConfig `json:"config"`
	Nonce     uint64                    `json:"nonce"`
	Timestamp *big.Int                  `json:"timestamp"`
	ExtraData []byte                    `json:"extraData"`
	GasLimit  uint64                    `json:"gasLimit"   gencodec:"required"`
	//Difficulty *big.Int            `json:"difficulty" gencodec:"required"`
	Difficulty common.Difficulty `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash       `json:"mixHash"`
	Coinbase   common.Address    `json:"coinbase"`
	Alloc      GenesisAlloc      `json:"alloc"      gencodec:"required"`

	//add verifiers
	Verifiers []common.Address

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number uint64 `json:"number"`
	//GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`

	ChainDB               chaindb.Database
	AccountStateProcessor state_processor.AccountStateProcessor
	RegisterProcessor     registerdb.RegisterProcessor

	merged bool
}

func (g *Genesis) SetChainDB(db chaindb.Database) {
	g.ChainDB = db
}

func (g *Genesis) SetVerifiers(v []common.Address) {
	g.Verifiers = v
}

func (g *Genesis) SetAccountStateProcessor(p state_processor.AccountStateProcessor) {
	g.AccountStateProcessor = p
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]*big.Int

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  main-net default  |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(genesis *Genesis) (*chain_config.ChainConfig, common.Hash, error) {
	if genesis == nil {
		log.Info("Writing default main-net genesis block")
		return nil, common.Hash{}, errors.New("genesis can't be nil")
	}

	if genesis.Config == nil {
		return nil, common.Hash{}, errGenesisNoConfig
	}

	chainDB := genesis.ChainDB
	// Just commit the new block if there is no stored genesis block.
	stored := chainDB.GetBlockHashByNumber(0)
	log.Info("get stored genesis hash", "hash", stored.Hex())

	// prepare genesis block
	// Every time you start, you must compare the configured Genesis block with the Genesis block on the chain. If it is not the same, you need to report an error.
	block, err := genesis.Prepare()
	if err != nil {
		log.Error("prepare genesis failed", "storageErr", err)
		return nil, common.Hash{}, err
	}
	if (stored == common.Hash{}) {
		log.Info("write genesis block", "hash", block.Hash().Hex())
		err = genesis.Commit(block)
		// check contract committed
		genesis.checkEarlyContractExist()
		return genesis.Config, block.Hash(), err
	}

	// Check whether the genesis block is already written.
	if genesis != nil {
		hash := block.Hash()
		if hash != stored {
			// todo need reset ?
			//genesis.resetDataDirIfMercury(dataDir)
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
		// check early contract exist
		genesis.checkEarlyContractExist()
	}
	// write chain config to db?

	return genesis.Config, stored, nil
}

func (g *Genesis) checkEarlyContractExist() {
	earlyCV, err := g.AccountStateProcessor.(*state_processor.AccountStateDB).Copy().GetContract(contract.EarlyContractAddress, reflect.TypeOf(contract.EarlyRewardContract{}))
	if err != nil {
		panic("check early contract failed, storageErr: " + err.Error())
	}
	earlyC := earlyCV.Interface().(*contract.EarlyRewardContract)

	var originEarlyC contract.EarlyRewardContract
	if err = util.ParseJson(contract.EarlyRewardContractStr, &originEarlyC); err != nil {
		panic("parse origin early contract failed, storageErr: " + err.Error())
	}
	if !earlyC.Owner.IsEqual(originEarlyC.Owner) {
		panic(fmt.Sprintf("early contract owner not match, is: %v, should be: %v", earlyC.Owner.Hex(), originEarlyC.Owner.Hex()))
	}
	//log.Debug("originEarlyC", "addr", "0x000095Cfdd141b0aF2Bb92F0074d5Dbc9b5F554fF807", "eDIP balance", originEarlyC.Balances["0x000095Cfdd141b0aF2Bb92F0074d5Dbc9b5F554fF807"])

	// check if the alloc value is correct
	for addr, a := range g.Alloc {
		b, err := g.AccountStateProcessor.GetBalance(addr)
		if err != nil {
			panic("alloc check failed, storageErr: %v" + err.Error())
		}
		shouldBe := big.NewInt(0).Add(a, big.NewInt(0))
		if addr.IsEqual(earlyC.Owner) {
			shouldBe = big.NewInt(0).Sub(a, originEarlyC.NeedDIP)
		}
		if b.Cmp(shouldBe) != 0 {
			panic(fmt.Sprintf("alloc value not right, addr: %v, got: %v, should be: %v", addr, b, shouldBe))
		}
	}
	log.Info("genesis early contract check success")
}

/*func (g *Genesis) resetDataDirIfMercury(dataDir string) {
	if os.Getenv("boots_env") != "mercury" {
		return
	}
	if !common.FileExist(dataDir) {
		log.Error("can't reset mercury datadir", "datadir", dataDir)
		return
	}
	bakDir := filepath.Join(util.HomeDir(), "dipperin_mercury_latest_datadir_bak")
	os.RemoveAll(bakDir)
	if storageErr := os.Rename(dataDir, bakDir); storageErr != nil {
		panic(fmt.Sprintf("bak datadir failed, storageErr: %v", storageErr))
	}
	log.Info("genesis not match reset datadir, you should restart dipperin", "bak to", bakDir)
	os.Exit(1)
}*/

/*func (g *Genesis) configOrDefault(ghash common.Hash) *chain_config.ChainConfig {
	switch {
	//case g != nil:
	//	return g.Config
	//case ghash == params.MainnetGenesisHash:
	//	return params.MainnetChainConfig
	//case ghash == params.TestnetGenesisHash:
	//	return params.TestnetChainConfig
	default:
		return chain_config.GetChainConfig()
	}
}*/

func (g *Genesis) Valid() bool {
	var totalAmount = big.NewInt(0)
	for _, amount := range g.Alloc {
		totalAmount.Add(totalAmount, amount)
	}
	if totalAmount.Cmp(consts.MaxAmount) > 0 {
		fmt.Println("Genesis config not correct. Total allocation exceed maximum limit")
		return false
	}
	return true
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock() *model.Block {
	head := &model.Header{
		Version:   0,
		Number:    0,
		Nonce:     common.EncodeNonce(g.Nonce),
		TimeStamp: g.Timestamp,
		PreHash:   g.ParentHash,
		//FIXME generate a random seed.
		Seed:        g.ParentHash,
		Proof:       make([]byte, 0),
		MinerPubKey: make([]byte, 0),
		Diff:        g.Difficulty,
		Bloom:       iblt.NewBloom(model.DefaultBlockBloomConfig),
		GasLimit:    g.GasLimit,
	}

	if g.Difficulty.Equal(common.Difficulty{}) {
		head.Diff = chain_config.GenesisDifficulty
	}

	//padding economy model balance
	economyModel := economy_model.MakeDipperinEconomyModel(nil, economy_model.DIPProportion)
	err := g.paddingEconomyInfo(economyModel)
	if err != nil {
		panic("padding economy info error: " + err.Error())
	}

	block := model.NewBlock(head, nil, nil)

	return block
}

func (g *Genesis) Commit(block model.AbstractBlock) error {
	// commit states
	if _, err := g.AccountStateProcessor.Commit(); err != nil {
		log.Error("init accountStateProcessor failed", "storageErr", err)
		return err
	}

	// commit register
	if _, err := g.RegisterProcessor.Commit(); err != nil {
		log.Error("init registerProcessor failed", "storageErr", err)
		return err
	}

	// write block
	g.ChainDB.InsertBlock(block)
	return nil
}

func (g *Genesis) paddingEconomyInfo(economyModel economy_model.EconomyModel) (err error) {
	if g.merged {
		return nil
	}
	g.merged = true

	//set economy model addresses balance
	investorInfo := economyModel.GetInvestorInitBalance()

	developerInfo := economyModel.GetDeveloperInitBalance()

	maintenanceInfo := economyModel.GetFoundation().GetFoundationInfo(economy_model.Maintenance)

	remainRewardInfo := economyModel.GetFoundation().GetFoundationInfo(economy_model.RemainReward)

	earlyTokenInfo := economyModel.GetFoundation().GetFoundationInfo(economy_model.EarlyToken)

	if err = economy_model.MapMerge(g.Alloc, investorInfo); err != nil {
		return
	}

	if err = economy_model.MapMerge(g.Alloc, developerInfo); err != nil {
		return
	}

	if err = economy_model.MapMerge(g.Alloc, maintenanceInfo); err != nil {
		return
	}

	if err = economy_model.MapMerge(g.Alloc, remainRewardInfo); err != nil {
		return
	}

	if err = economy_model.MapMerge(g.Alloc, earlyTokenInfo); err != nil {
		return
	}

	return
}

func (g *Genesis) SetEarlyTokenContract() error {
	var earlyTokenContract contract.EarlyRewardContract
	if err := util.ParseJson(contract.EarlyRewardContractStr, &earlyTokenContract); err != nil {
		panic(err.Error())
	}

	balance, err := g.AccountStateProcessor.GetBalance(earlyTokenContract.Owner)
	if err != nil {
		log.Info("the account address is:", "address", earlyTokenContract.Owner.Hex())
		return err
	}

	if balance.Cmp(earlyTokenContract.NeedDIP) == -1 {
		return errors.New("the contract owner balance isn't enough")
	}

	// sub the money used for EDIP and put it in the contract.
	err = g.AccountStateProcessor.SetBalance(earlyTokenContract.Owner, big.NewInt(0).Sub(balance, earlyTokenContract.NeedDIP))
	if err != nil {
		return err
	}

	// todo be especially careful here
	//write earlyToken contract
	err = g.AccountStateProcessor.PutContract(contract.EarlyContractAddress, reflect.ValueOf(&earlyTokenContract))
	//storageErr = g.contractDB.PutContract(contract.EarlyContractAddress, []byte(contract.EarlyRewardContractStr))
	if err != nil {
		return err
	}

	return nil
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Prepare() (model.AbstractBlock, error) {
	block := g.ToBlock()
	/*	if block.Number() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}*/

	//write verifier state
	for _, v := range g.Verifiers {
		if err := g.AccountStateProcessor.NewAccountState(v); err != nil {
			return nil, err
		}
	}

	// write state，Must be placed after the initial verifier account, otherwise the amount of the verifier in alloc will be overwritten to 0
	for k, v := range g.Alloc {
		//log.Debug("add genesis balance", "addr", k.Hex(), "balance", v)
		if err := g.AccountStateProcessor.NewAccountState(k); err != nil {
			return nil, err
		}

		if err := g.AccountStateProcessor.SetBalance(k, v); err != nil {
			return nil, err
		}
	}

	err := g.SetEarlyTokenContract()
	if err != nil {
		return nil, err
	}

	/*	//todo delete after test get contract
		result,storageErr := g.accountStateProcessor.GetContract(contract.EarlyContractAddress)
		if storageErr != nil {
			return nil, storageErr
		}
		log.Info("the genesis contract is:","result",string(result))*/

	stateRoot, err := g.AccountStateProcessor.Finalise()
	if err != nil {
		return nil, err
	}

	block.SetStateRoot(stateRoot)

	if err = g.RegisterProcessor.PrepareRegisterDB(); err != nil {
		return nil, err
	}

	registerRoot := g.RegisterProcessor.Finalise()
	block.SetRegisterRoot(registerRoot)
	log.Info("set genesis registerDB successful", "root", registerRoot)

	return block, nil
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
/*func (g *Genesis) MustCommit() model.AbstractBlock {
	block, storageErr := g.Prepare()
	if storageErr != nil {
		panic(storageErr)
	}
	storageErr = g.Commit(block)
	if storageErr != nil {
		panic(storageErr)
	}
	return block
}*/

// GenesisBlockForTesting creates and writes a block in which addr has the given wei balance.
/*func GenesisBlockForTesting(addr common.Address, balance *big.Int) model.AbstractBlock {
	g := Genesis{Alloc: GenesisAlloc{addr: balance}}
	return g.MustCommit()
}*/

// DefaultGenesisBlock returns the Ethereum main net genesis block.
func DefaultGenesisBlock(chainDB chaindb.Database, accountStateProcessor state_processor.AccountStateProcessor, registerProcessor registerdb.RegisterProcessor, chainConf *chain_config.ChainConfig) *Genesis {
	log.Debug("call DefaultGenesisBlock")

	//read config file first
	if mGenesis := GenesisBlockFromFile(chainDB, accountStateProcessor); mGenesis != nil {
		return mGenesis
	}

	gTime, _ := time.Parse("2006-01-02 15:04:05", "2018-08-08 08:08:08")
	// use to reset test chain
	if chain_config.GetCurBootsEnv() == "test" {
		gTime, _ = time.Parse("2006-01-02 15:04:05", "2019-01-14 08:08:08")
	}
	return &Genesis{
		ChainDB:               chainDB,
		GasLimit:              chain_config.BlockGasLimit,
		AccountStateProcessor: accountStateProcessor,
		RegisterProcessor:     registerProcessor,
		Config:                chain_config.GetChainConfig(),
		//Nonce:      66,
		Timestamp:  big.NewInt(gTime.UnixNano()),
		ExtraData:  []byte("dipperin Genesis"),
		Difficulty: chain_config.GenesisDifficulty,
		Alloc:      map[common.Address]*big.Int{
						//for test
						common.HexToAddress("0x0000062493b705D52E4541e7Daa6343A8eD98d8dc15f"): big.NewInt(0).Mul(big.NewInt(1e8),big.NewInt(consts.DIP)),

			// corresponding private key:289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032
			/*			common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9"): big.NewInt(100 * consts.DIP),

						common.HexToAddress("0x0000b10d5b64AaF00CAF6FFf2dae7A38Cd5258a4e347"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x00000426EBf7eB1BB5CCCa4a0CC2413AbDb10cc86692"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000731306d72b6f3d8f480bC1E262A0f4c799769d14"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000B931313E887D3d51e9Af6c93Ed68117B1008FB5C"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000eacD02FA7a68965B317D609a583B578dF54e4E66"): big.NewInt(100 * consts.DIP),

						common.HexToAddress("0x00000BdE9CA0D03AFa040946C5bB274e4B3eBbD77CBB"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000aB662AB5617ec66Ce06c127DC779D1EEB5d7570c"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x000051Cc1d09F4d054C6d4171b0dD25e24A82e5664C6"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000C9bDe227FDEB9F62D6A377e081F27f8415fA61D1"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000cA84D03A2159f7779Fe712E43e5b2Cc49549d65C"): big.NewInt(100 * consts.DIP),

						common.HexToAddress("0x0000A7Ac94e8Da79505153896E201215a3E4a10F2a30"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000fdeB15BC5cf98F40c178B9148E53560bf9ACA891"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000Bed1a411136faA5Aece75474383707Feef4D65F8"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000408FAF7CA6289d7C45e0139adFc44F0841BdE521"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x00006e271bC6317d487345CB508DeDba36C2AE0EA8Ca"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000C3F6F2aC77A71fCfB33197eEA1A3e90d317BEebb"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000e899F84129b57ba2Da61718A3e8593dFe49B4BB4"): big.NewInt(100 * consts.DIP),
						common.HexToAddress("0x0000f01dA91C64eF6202c735e9362010196a556C7fc7"): big.NewInt(100 * consts.DIP),*/
		},
		Verifiers: VerifierAddress[:chainConf.VerifierNumber],
	}
}

type genesisCfgFile struct {
	Nonce uint64 `json:"nonce"`
	//Note       string           `json:"note"`
	Accounts   map[string]int64 `json:"accounts"`
	Timestamp  string           `json:"timestamp"`
	Difficulty string           `json:"difficulty" gencodec:"required"`
	Verifiers  []string         `json:"verifiers" gencodec:"required"`
	// todo add a foundation configuration
}

func GenesisBlockFromFile(chainDB chaindb.Database, accountStateProcessor state_processor.AccountStateProcessor) *Genesis {
	log.Debug("call GenesisBlockFromFile")

	gFPath := filepath.Join(util.HomeDir(), "softwares", "dipperin_deploy", "genesis.json")
	ge, e := ioutil.ReadFile(gFPath)
	if e != nil {
		return nil
	}
	log.Info("load genesis file", "path", gFPath)

	var info genesisCfgFile
	err := json.Unmarshal(ge, &info)
	if err != nil {
		log.Error("unmarshal genesisCfgFile failed", "storageErr", err)
		return nil
	}

	var gTime time.Time
	if gTime, err = time.Parse("2006-01-02 15:04:05", info.Timestamp); err != nil {
		gTime, _ = time.Parse("2006-01-02 15:04:05", "2018-08-08 08:08:08")
	}

	alloc := make(map[common.Address]*big.Int)
	for k, v := range info.Accounts {
		if v <= 0 {
			panic(fmt.Sprintf("genesis account balance wrong, %v: %v", k, v))
		}
		alloc[common.HexToAddress(k)] = big.NewInt(0).Mul(big.NewInt(v), big.NewInt(consts.DIP))
	}

	var verifiers []common.Address
	for _, v := range info.Verifiers {
		verifiers = append(verifiers, common.HexToAddress(v))
	}

	return &Genesis{
		ChainDB: chainDB,

		AccountStateProcessor: accountStateProcessor,
		Config:                chain_config.GetChainConfig(),
		Nonce:                 info.Nonce,
		Timestamp:             big.NewInt(gTime.UnixNano()),
		//ExtraData:             []byte(info.Note),
		Difficulty: common.HexToDiff(info.Difficulty),
		Alloc:      alloc,
		Verifiers:  verifiers,
	}
}
