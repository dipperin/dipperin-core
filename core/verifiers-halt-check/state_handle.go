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

package verifiers_halt_check

import (
	"crypto/ecdsa"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./chain_reader_mock_test.go -package=verifiers_halt_check github.com/dipperin/dipperin-core/core/verifiers-halt-check NeedChainReaderFunction
type NeedChainReaderFunction interface {
	CurrentBlock() model.AbstractBlock
	GetSeenCommit(height uint64) []model.AbstractVerification
	GetCurrVerifiers() []common.Address
	SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error

	//AccountDB need ChainReader function
	GetBlockByNumber(number uint64) model.AbstractBlock
	GetVerifiers(round uint64) []common.Address

	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetLastChangePoint(block model.AbstractBlock) *uint64
	GetSlot(block model.AbstractBlock) *uint64

	BlockProcessor(root common.Hash) (*chain.BlockProcessor, error)
	BlockProcessorByNumber(num uint64) (*chain.BlockProcessor, error)

	BuildRegisterProcessor(preRoot common.Hash) (*registerdb.RegisterDB, error)
}

//go:generate mockgen -destination=./wallet_signer_mock_test.go -package=verifiers_halt_check github.com/dipperin/dipperin-core/core/verifiers-halt-check NeedWalletSigner
type NeedWalletSigner interface {
	GetAddress() common.Address
	SignHash(hash []byte) ([]byte, error)
	PublicKey() *ecdsa.PublicKey
	ValidSign(hash []byte, pubKey []byte, sign []byte) error
	Evaluate(account accountsbase.Account, seed []byte) (index [32]byte, proof []byte, err error)
}

type StateHandler struct {
	chainReader  NeedChainReaderFunction
	walletSigner NeedWalletSigner
	//need economy model for processing state
	economyModel economy_model.EconomyModel
	//need state storage for processing state
	//stateStorage state_processor.StateStorage
}

func MakeHaltCheckStateHandler(needChainReader NeedChainReaderFunction, walletSigner NeedWalletSigner, economyModel economy_model.EconomyModel) *StateHandler {
	return &StateHandler{
		chainReader:  needChainReader,
		walletSigner: walletSigner,
		economyModel: economyModel,
	}
}

func (haltCheckStateHandle *StateHandler) GenProposalConfig(voteType model.VoteMsgType) (ProposalGeneratorConfig, error) {
	curBlock := haltCheckStateHandle.chainReader.CurrentBlock()
	log.DLogger.Info("GenerateEmptyBlock", zap.Uint64("num", curBlock.Number()))
	account := accountsbase.Account{Address: haltCheckStateHandle.walletSigner.GetAddress()}

	seed, proof, err := haltCheckStateHandle.walletSigner.Evaluate(account, curBlock.Seed().Bytes())
	if err != nil {
		return ProposalGeneratorConfig{}, g_error.GenProposalConfigError
	}

	currentHeight := curBlock.Number()
	verifications := haltCheckStateHandle.chainReader.GetSeenCommit(currentHeight)
	config := ProposalGeneratorConfig{
		CurBlock:          curBlock,
		NewBlockSeed:      seed,
		NewBlockProof:     proof,
		LastVerifications: verifications,
		PubKey:            crypto.FromECDSAPub(haltCheckStateHandle.walletSigner.PublicKey()),
		SignHashFunc:      haltCheckStateHandle.walletSigner.SignHash,
		ProcessStateFunc:  haltCheckStateHandle.ProcessAccountAndRegisterState,
		VoteType:          voteType,
	}

	return config, nil
}

func (haltCheckStateHandle *StateHandler) ProcessAccountAndRegisterState(block model.AbstractBlock, preStateRoot, preRegisterRoot common.Hash) (stateRoot, registerRoot common.Hash, err error) {
	log.DLogger.Info("the preStateRoot is:", zap.String("preStateRoot", preStateRoot.Hex()))
	accountDB, err := haltCheckStateHandle.chainReader.BlockProcessor(preStateRoot)
	if err != nil {
		return common.Hash{}, common.Hash{}, err
	}

	//process account state
	if err = accountDB.ProcessExceptTxs(block, haltCheckStateHandle.economyModel, false); err != nil {
		log.DLogger.Error("process state except txs failed", zap.Error(err))
		return common.Hash{}, common.Hash{}, err
	}

	stateRoot, err = accountDB.Finalise()
	if err != nil {
		return common.Hash{}, common.Hash{}, err
	}

	//process register state
	log.DLogger.Info("the preRegisterRoot is:", zap.String("preRegisterRoot", preRegisterRoot.Hex()))
	registerDB, err := haltCheckStateHandle.chainReader.BuildRegisterProcessor(preRegisterRoot)
	if err != nil {
		return common.Hash{}, common.Hash{}, err
	}

	if err = registerDB.Process(block); err != nil {
		log.DLogger.Error("process register failed", zap.Error(err))
		return common.Hash{}, common.Hash{}, err
	}
	registerRoot = registerDB.Finalise()

	log.DLogger.Info("the calculated empty block root is:", zap.String("stateRoot", stateRoot.Hex()), zap.String("registerRoot", registerRoot.Hex()))

	return stateRoot, registerRoot, nil
}

func (haltCheckStateHandle *StateHandler) SaveFinalEmptyBlock(proposal ProposalMsg, votes map[common.Address]model.VoteMsg) error {
	// use boot node verifier vote and alive verifier votes received as the verifications of empty block
	verifications := make([]model.AbstractVerification, 0)
	verifications = append(verifications, &proposal.VoteMsg)
	for _, tmpValue := range votes {
		tmpVote := tmpValue
		verifications = append(verifications, &tmpVote)
	}

	log.DLogger.Info("save and Broadcast empty block", zap.Any("blockHash", proposal.EmptyBlock.Hash()))
	log.DLogger.Info("save and Broadcast verifications 3 ", zap.Any("verifications", verifications))

	log.DLogger.Info("save and Broadcast verifications", zap.Uint64("verifications", verifications[0].GetHeight()))
	err := haltCheckStateHandle.chainReader.SaveBlock(&proposal.EmptyBlock, verifications)
	if err != nil {
		log.DLogger.Info("verifier boot node save empty block failed", zap.Error(err))
		if err.Error() != "already have this block" {
			return err
		}
	}

	return nil
}
