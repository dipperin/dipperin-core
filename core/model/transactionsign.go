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

package model

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
)

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")
)

type sigCache struct {
	signer Signer
	from   common.Address
}

func (tx *Transaction) SignTx(priKey *ecdsa.PrivateKey, s Signer) (*Transaction, error) {
	wit := witness{HashKey: tx.wit.HashKey}
	h, err := s.GetSignHash(tx)
	if err != nil {
		return nil, err
	}

	//log.Debug("the tx hash is:","hash",h.Hex())
	sig, err := crypto.Sign(h[:], priKey)
	if err != nil {
		return nil, err
	}
	wit.R, wit.S, wit.V, err = s.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	tx.wit = wit

	return tx, nil
}

//todo set signer config
func MakeSigner(config *chain_config.ChainConfig, blockNumber uint64) Signer {
	var signer Signer
	switch {
	default:
		signer = DipperinSigner{chainId: config.ChainId}
	}
	return signer
}

// signer define different kind of method handling signature
type Signer interface {
	// Sender returns the sender address of the transaction.
	GetSender(tx *Transaction) (common.Address, error)
	//Get Sender Public Key
	GetSenderPublicKey(tx *Transaction) (*ecdsa.PublicKey, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	// Hash returns the VRFHash to be signed.
	GetSignHash(rtx *Transaction) (common.Hash, error)
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

// Mercury chainId is 1
// Venus chainId is 2

type DipperinSigner struct{ chainId *big.Int }

func NewSigner(chainId *big.Int) DipperinSigner {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return DipperinSigner{
		chainId: chainId,
	}
}

func (ds DipperinSigner) Equal(s2 Signer) bool {
	s, ok := s2.(DipperinSigner)
	return ok && ds.chainId.Cmp(s.chainId) == 0
}

// GetSignHash will return the VRFHash of the transaction with have the raw transaction data and chainId
func (ds DipperinSigner) GetSignHash(rtx *Transaction) (common.Hash, error) {
	//log.Debug("DipperinSigner GetSignHash","tx",rtx.data)
	//log.Debug("DipperinSigner GetSignHash","chainId",fs.chainId)
	res, err := rlpHash([]interface{}{rtx.data, ds.chainId})
	return res, err
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (ds DipperinSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 54})

	// Fixme mul 2 is EIP 155 modification, which is used to avoid transaction on ETC reply on ETH, can remove this modification
	if ds.chainId.Sign() != 0 {
		v.Add(v, big.NewInt(0).Mul(ds.chainId, big.NewInt(2)))
	}
	return r, s, v, nil
}

func (ds DipperinSigner) GetSender(tx *Transaction) (common.Address, error) {
	//different type use different address type
	hash, err := ds.GetSignHash(tx)
	if err != nil {
		return common.Address{}, err
	}

	//log.Health.Info("GetSender the tx wit r s v is:","r",tx.wit.R,"s",tx.wit.S,"v",tx.wit.V)
	//log.Health.Info("GetSender the ds chainId is:","chainId",ds.chainId)
	temp := big.NewInt(0).Sub(tx.wit.V, big.NewInt(0).Mul(ds.chainId, big.NewInt(2)))
	v := big.NewInt(0).Sub(temp, big.NewInt(54))
	//log.Health.Info("the calculated v is:","v",v)

	return recoverNormalSender(hash, tx.wit.R, tx.wit.S, v)
}

func (ds DipperinSigner) GetSenderPublicKey(tx *Transaction) (*ecdsa.PublicKey, error) {
	//different type use different address type

	emptyPk := ecdsa.PublicKey{}
	sigHash, err := ds.GetSignHash(tx)
	if err != nil {
		return &emptyPk, err
	}
	temp := big.NewInt(0).Sub(tx.wit.V, big.NewInt(0).Mul(ds.chainId, big.NewInt(2)))
	V := big.NewInt(0).Sub(temp, big.NewInt(54))
	R := tx.wit.R
	S := tx.wit.S

	if V.BitLen() > 8 {
		log.Health.Info("GetSenderPublicKey the error V is:","V",V)
		return &emptyPk, ErrInvalidSig
	}
	if !cs_crypto.ValidSigValue(R, S, V) {
		log.Health.Error("GetSenderPublicKey valid Signature Value error")
		return &emptyPk, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s, v := R.Bytes(), S.Bytes(), V.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	copy(sig[64:], v)
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		return &emptyPk, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return &emptyPk, errors.New("invalid public key")
	}
	pubKey := cs_crypto.ToECDSAPub(pub)
	return pubKey, nil
}

func recoverNormalSender(sigHash common.Hash, R, S, V *big.Int) (common.Address, error) {
	log.Health.Info("recoverNormalSender the r s v is:","r",R,"s",S,"v",V)
	if V.BitLen() > 8 {
		log.Health.Error("recoverNormalSender v bitLen is more than 8")
		return common.Address{}, ErrInvalidSig
	}
	if !cs_crypto.ValidSigValue(R, S, V) {
		log.Health.Error("recoverNormalSender valid signature error")
		return common.Address{}, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s, v := R.Bytes(), S.Bytes(), V.Bytes()
	sig := make([]byte, 65)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	copy(sig[64:], v)
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	pubkey := cs_crypto.ToECDSAPub(pub)
	var addr common.Address
	addr = cs_crypto.GetNormalAddress(*pubkey)
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	res := new(big.Int).Sub(v, big.NewInt(54))
	return res.Div(res, big.NewInt(2))
}
