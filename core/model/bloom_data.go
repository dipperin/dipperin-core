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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
)

type BloomBlockData struct {
	Header   *Header
	BloomRLP []byte
	// pre height vs
	PreVerification []AbstractVerification
	// cur height vs
	CurVerification []AbstractVerification
	//interlins
	Interlinks InterLink
}

func (data *BloomBlockData) EiRecoverToBlock(txPoolMap map[common.Hash]AbstractTransaction) (block *Block, err error) {
	var bloom iblt.Graphene
	// rlp decode bloomRLP
	if err = rlp.DecodeBytes(data.BloomRLP, &bloom); err != nil {
		return nil, err
	}

	//new method recover txs
	//possibleTxsBytes,err:=bloom.Recover(txPoolMap)
	//possibleTxs,err:=data.rebuildTxs(possibleTxsBytes)
	//old method
	invBloom := iblt.NewGraphene(bloom.InvBloomConfig(), bloom.BloomConfig())
	possibleTxsMap := make(map[common.Hash]AbstractTransaction)

	for k, v := range txPoolMap {
		if bloom.LookUp(k.Bytes()) {
			invBloom.InsertRLP(k, v)
			possibleTxsMap[k] = v
		}
	}

	tempGraphene := iblt.NewGraphene(bloom.InvBloomConfig(), bloom.BloomConfig())

	recovered, others, err := tempGraphene.InvBloom().Subtract(bloom.InvBloom(), invBloom.InvBloom()).ListRLP()
	if err != nil {
		log.Error("invBloom can't recover tx", "err", err)
		return nil, err
	}

	recoveredTxs, err := data.rebuildTxs(recovered)
	recoveredTxsMap := make(map[common.Hash]AbstractTransaction)
	for _, tx := range recoveredTxs {
		recoveredTxsMap[tx.CalTxId()] = tx
	}
	othersTxs, err := data.rebuildTxs(others)
	othersTxsMap := make(map[common.Hash]AbstractTransaction)
	for _, tx := range othersTxs {
		othersTxsMap[tx.CalTxId()] = tx
	}

	for k := range othersTxsMap {
		if _, ok := possibleTxsMap[k]; ok {
			delete(possibleTxsMap, k)
		}
	}
	for k, v := range recoveredTxsMap {
		possibleTxsMap[k] = v
	}

	var possibleTxs []*Transaction
	for _, v := range possibleTxsMap {
		possibleTxs = append(possibleTxs, v.(*Transaction))
	}
	//
	if err != nil {
		log.Error("recover txs err", err)
		return nil, err
	}

	if block = NewBlock(data.Header, possibleTxs, data.PreVerification); block == nil {
		return nil, errors.New("new block is nil")
	}

	return block, nil
}

func (data *BloomBlockData) rebuildTxs(recovered [][]byte) ([]*Transaction, error) {
	var txs []*Transaction

	for i, txRLP := range recovered {
		var tx Transaction
		if err := rlp.DecodeBytes(txRLP, &tx); err != nil {
			log.Error("rlp decode invBloom recovered data error", "tx index", i, "err", err)
			return nil, err
		}

		txs = append(txs, &tx)
	}

	return txs, nil
}

// Input TXs-list and Algorithms for Goroutine
func RunWorkMap(task iblt.Operation, txs []*Transaction) error {
	lock.Lock()
	defer lock.Unlock()

	iblt.WorkMap.SetOperation(task)
	raw := make([]interface{}, len(txs))
	after := make([]bool, len(txs))
	for index, tx := range txs {
		raw[index] = tx
	}
	result, err := iblt.WorkMap.StartWorks(raw)
	if err != nil {
		return err
	}

	//process result
	j := 0
pass:
	for {
		select {
		case r := <-result:
			if r.(bool) == true {
				after[j] = r.(bool)
				j++
			}
			if j == len(txs) {
				break pass
			}
		}
	}
	return nil
}

type MapWorkHybridEstimator struct {
	hybrid *iblt.HybridEstimator
}

func newMapWorkHybridEstimator(hy *iblt.HybridEstimator) *MapWorkHybridEstimator {
	return &MapWorkHybridEstimator{
		hybrid: hy,
	}
}

func (estimator MapWorkHybridEstimator) DoTask(i interface{}) interface{} {
	tx := i.(*Transaction)
	estimator.hybrid.EncodeByte(tx.CalTxId().Bytes())
	return true
}

type MapWorkInvBloom struct {
	graphene *iblt.Graphene
}

func newMapWorkInvBloom(g *iblt.Graphene) *MapWorkInvBloom {
	return &MapWorkInvBloom{
		graphene: g,
	}
}

func (bloom MapWorkInvBloom) DoTask(i interface{}) interface{} {
	tx := i.(*Transaction)
	bloom.graphene.InsertRLP(tx.CalTxId(), tx)
	bloom.graphene.Bloom().Digest(tx.CalTxId().Bytes())
	return true
}
