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


package tests

import (
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/common"
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"math/big"
	"time"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/common/consts"
)

func NewGenesisEnv(chainDB chaindb.Database, stateStorage state_processor.StateStorage, accounts []Account) *GenesisEnv {
	g := &GenesisEnv{chainConf: chain_config.GetChainConfig()}

	g.initMiner()
	g.initVerifiers(accounts)
	g.initGenesis(chainDB, stateStorage)
	return g
}

type Account struct {
	Pk      *ecdsa.PrivateKey
	address common.Address
}

func NewAccount(pk *ecdsa.PrivateKey, address common.Address) *Account {
	return &Account{Pk: pk, address: address}
}

func (a *Account) Address() common.Address {
	if !a.address.IsEmpty() {
		return a.address
	}

	a.address = cs_crypto.GetNormalAddress(a.Pk.PublicKey)
	return a.address
}

func (a *Account) SignHash(hash []byte) ([]byte, error) {
	return crypto.Sign(hash, a.Pk)
}

type Accounts []Account

func (as Accounts) GetAddresses() (r []common.Address) {
	accs := ([]Account)(as)
	for _, a := range accs {
		r = append(r, a.Address())
	}
	return
}

// Initialize the default validator in which the external test calls its method to vote on the block
type GenesisEnv struct {
	genesis   *chain.Genesis
	gBlock    *model.Block
	chainConf *chain_config.ChainConfig

	defaultVerifiers []Account
	miner            Account
}

func (g *GenesisEnv) Miner() Account {
	return g.miner
}

func (g *GenesisEnv) DefaultVerifiers() []Account {
	return g.defaultVerifiers
}

func (g *GenesisEnv) initMiner() {
	g.miner = Account{Pk: crypto.HexToECDSAErrPanic("1e00ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")}
}

func (g *GenesisEnv) initGenesis(chainDB chaindb.Database, stateStorage state_processor.StateStorage) {
	genesisAccountStateProcessor, err := state_processor.MakeGenesisAccountStateProcessor(stateStorage)
	errPanic(err)
	genesisRegisterProcessor, err := registerdb.MakeGenesisRegisterProcessor(stateStorage)
	errPanic(err)

	gTime, _ := time.Parse("2006-01-02 15:04:05", "2019-06-06 08:08:08")
	g.genesis = &chain.Genesis{
		ChainDB:               chainDB,
		AccountStateProcessor: genesisAccountStateProcessor,
		RegisterProcessor:     genesisRegisterProcessor,
		Config:                chain_config.GetChainConfig(),
		Timestamp:             big.NewInt(gTime.UnixNano()),
		ExtraData:             []byte("dipperin Genesis"),
		Difficulty:            common.HexToDiff("0x1fffffff"),
		Alloc: map[common.Address]*big.Int{
			g.defaultVerifiers[0].Address(): big.NewInt(9999 * consts.DIP),
			g.defaultVerifiers[1].Address(): big.NewInt(9999 * consts.DIP),
			g.defaultVerifiers[2].Address(): big.NewInt(9999 * consts.DIP),
		},
		Verifiers: chain.VerifierAddress[:g.chainConf.VerifierNumber],
		GasLimit: chain_config.BlockGasLimit,
	}
	g.gBlock = g.genesis.ToBlock()
	_, _, err = chain.SetupGenesisBlock(g.genesis)
	errPanic(err)
}

func (g *GenesisEnv) initVerifiers(accounts []Account) {
	g.defaultVerifiers, _ = ChangeVerifierAddress(accounts)
}

// num is the number of votes
func (g *GenesisEnv) VoteBlock(num int, round uint64, b model.AbstractBlock) (result []model.AbstractVerification) {

	vLen := len(g.defaultVerifiers)
	for i := 0; i < num && i < vLen; i++ {
		verifier := g.defaultVerifiers[i]
		m, err := model.NewVoteMsgWithSign(b.Number(), round, b.Hash(), model.VoteMessage, verifier.SignHash, verifier.Address())
		if err != nil {
			panic(err)
		}
		result = append(result, m)
	}

	return
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
